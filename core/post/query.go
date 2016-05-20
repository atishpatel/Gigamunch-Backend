package post

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	mysqlDB  *sql.DB
	errSQLDB = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
)

const (
	datetimeFormat = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"
	sortByDate     = `SELECT post_id, item_id, gigachef_id,( 3959 * acos( cos( radians(%f) ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(%f) ) + sin( radians(%f) ) * sin( radians( latitude ) ) ) ) AS distance
                      FROM live_posts
											WHERE ready_datetime
											BETWEEN %s
											HAVING distance < %d
                      ORDER BY ready_datetime %s, distance
                      LIMIT %d , %d`
	selectClosedPosts = "SELECT post_id, gigachef_id FROM live_posts where close_datetime<'%s'"
	deleteStatement   = "DELETE FROM live_posts WHERE post_id=%d"
)

func insertLivePost(postID int64, post *Post) error {
	if mysqlDB == nil {
		connectSQL()
	}
	_, err := mysqlDB.Exec(
		`INSERT
		INTO live_posts
		(post_id, item_id, gigachef_id,close_datetime, ready_datetime, search_tags, is_order_now, is_experimental, is_baked_good, latitude, longitude)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		postID,
		post.ItemID,
		post.GigachefID,
		post.ClosingDateTime.UTC().Format(datetimeFormat),
		post.ReadyDateTime.UTC().Format(datetimeFormat),
		post.Title,
		post.IsOrderNow,
		0,
		0,
		post.GigachefAddress.Latitude,
		post.GigachefAddress.Longitude,
	)
	return err
}

func selectLivePosts(ctx context.Context, geopoint *types.GeoPoint, limit *types.Limit, radius int, readyDatetime time.Time, descending bool) ([]int64, []int64, []string, []float32, error) {
	if mysqlDB == nil {
		connectSQL()
	}
	var err error
	listLength := limit.End - limit.Start
	livePostQuery := getSortByDateQuery(geopoint.Latitude, geopoint.Longitude, radius,
		readyDatetime.UTC(), descending, limit)
	rows, err := mysqlDB.Query(livePostQuery)
	if err != nil {
		return nil, nil, nil, nil, errSQLDB.WithError(err)
	}
	defer handleCloser(ctx, rows)

	tmpPostIDs := make([]int64, listLength)
	tmpItemIDs := make([]int64, listLength)
	tmpDistances := make([]float32, listLength)
	tmpGigachefIDs := make([]string, listLength)
	actualReturnedRows := 0
	for i := 0; rows.Next(); i++ {
		_ = rows.Scan(&tmpPostIDs[i], &tmpItemIDs[i], &tmpGigachefIDs[i], &tmpDistances[i])
		actualReturnedRows++
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, nil, nil, errSQLDB.WithError(err).Wrap("failed while iterating rows")
	}
	// resize to actual meal returned size
	var postIDs []int64
	var itemIDs []int64
	var distances []float32
	var gigachefIDs []string
	if listLength != actualReturnedRows {
		postIDs = make([]int64, actualReturnedRows)
		itemIDs = make([]int64, actualReturnedRows)
		distances = make([]float32, actualReturnedRows)
		gigachefIDs = make([]string, actualReturnedRows)
		copy(postIDs, tmpPostIDs)
		copy(itemIDs, tmpItemIDs)
		copy(distances, tmpDistances)
		copy(gigachefIDs, tmpGigachefIDs)
	} else {
		postIDs = tmpPostIDs
		itemIDs = tmpItemIDs
		distances = tmpDistances
		gigachefIDs = tmpGigachefIDs
	}
	return postIDs, itemIDs, gigachefIDs, distances, nil
}

func getClosedPosts(ctx context.Context) ([]int64, []string, error) {
	if mysqlDB == nil {
		connectSQL()
	}
	statement := fmt.Sprintf(selectClosedPosts, time.Now().UTC().Add(1*time.Minute).Format(datetimeFormat))
	rows, err := mysqlDB.Query(statement)
	if err != nil {
		return nil, nil, errSQLDB.WithError(err).Wrapf("cannot query following statement: %s", statement)
	}
	defer handleCloser(ctx, rows)

	var postIDs []int64
	var chefIDs []string
	var tmpPostID int64
	var tmpChefID string
	for rows.Next() {
		_ = rows.Scan(&tmpPostID, &tmpChefID)
		postIDs = append(postIDs, tmpPostID)
		chefIDs = append(chefIDs, tmpChefID)
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
	}
	return postIDs, chefIDs, nil
}

func getSortByDateQuery(latitude float64, longitude float64, radius int, readyTime time.Time, descending bool, limit *types.Limit) string {
	var readyDatetimeOrder, readyWhere string
	if descending {
		readyWhere = "'2014-04-01 00:00:00' AND '" + readyTime.Format(time.RFC3339) + "'"
		readyDatetimeOrder = "DESC"
	} else {
		readyWhere = "'" + readyTime.Format(time.RFC3339) + "' AND '4000-04-01 00:00:00'"
		readyDatetimeOrder = "ASC"
	}
	return fmt.Sprintf(sortByDate, latitude, longitude, latitude, readyWhere, radius, readyDatetimeOrder, limit.Start, limit.End)
}

func removeLivePost(postID int64) error {
	if mysqlDB == nil {
		connectSQL()
	}
	_, err := mysqlDB.Exec(fmt.Sprintf(deleteStatement, postID))
	if err != nil {
		return errSQLDB.WithError(err).Wrap("cannot execute delete statement")
	}
	return nil
}

func connectSQL() {
	var err error
	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = "root@/gigamunch"
	} else {
		projectID := os.Getenv("PROJECTID")
		if projectID == "" {
			log.Fatal("PROJECTID env variable is not set.")
		}
		// MYSQL_CONNECTION: user:password@tcp([host]:3306)/dbname
		connectionString = fmt.Sprintf("root@cloudsql(%s:us-central1:gigasqldb)/gigamunch", projectID)
	}
	mysqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Couldn't connect to mysql database")
	}
}

type closer interface {
	Close() error
}

func handleCloser(ctx context.Context, c closer) {
	err := c.Close()
	if err != nil {
		utils.Errorf(ctx, "Error closing rows: %v", err)
	}
}

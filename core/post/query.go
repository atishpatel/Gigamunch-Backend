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

	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	mysqlDB *sql.DB
)

const (
	sortByDate = `SELECT post_id, gigachef_id,( 3959 * acos( cos( radians(%f) ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(%f) ) + sin( radians(%f) ) * sin( radians( latitude ) ) ) ) AS distance
                      FROM live_posts
											WHERE ready_datetime
											BETWEEN %s
											HAVING distance < %d
                      ORDER BY ready_datetime %s, distance
                      LIMIT %d , %d`
)

func insertLivePost(postID int64, post *Post) error {
	_, err := mysqlDB.Exec(
		`INSERT
		INTO live_posts
		(post_id, gigachef_id,close_datetime, ready_datetime, search_tags, is_experimental, is_baked_good, latitude, longitude)
		VALUES (?,?,?,?,?,?,?,?,?)`,
		postID, post.GigachefID, post.ClosingDateTime.Format(time.RFC3339), post.ReadyDateTime.Format(time.RFC3339), post.Title, 0, 0, post.Latitude, post.Longitude)
	return err
}

func selectLivePosts(ctx context.Context, geopoint *types.GeoPoint, limit *types.Limit, radius int, readyDatetime time.Time, descending bool) ([]int64, []string, []float32, error) {
	var err error
	listLength := limit.End - limit.Start
	livePostQuery := getSortByDateQuery(geopoint.Latitude, geopoint.Longitude, radius,
		readyDatetime, descending, limit.Start, limit.End)
	rows, err := mysqlDB.Query(livePostQuery)
	if err != nil {
		return nil, nil, nil, err //TODO change to external dep err
	}
	defer rows.Close()
	tmpPostIDs := make([]int64, listLength)
	tmpDistances := make([]float32, listLength)
	tmpGigachefIDs := make([]string, listLength)
	actualReturnedRows := 0
	for i := 0; rows.Next(); i++ {
		rows.Scan(&tmpPostIDs[i], &tmpGigachefIDs[i], &tmpDistances[i])
		actualReturnedRows++
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, nil, err // TODO change to external dep err
	}
	// resize to actual meal returned size
	var postIDs []int64
	var distances []float32
	var gigachefIDs []string
	if listLength != actualReturnedRows {
		postIDs = make([]int64, actualReturnedRows)
		distances = make([]float32, actualReturnedRows)
		gigachefIDs = make([]string, actualReturnedRows)
		copy(postIDs, tmpPostIDs)
		copy(distances, tmpDistances)
		copy(gigachefIDs, tmpGigachefIDs)
	} else {
		postIDs = tmpPostIDs
		distances = tmpDistances
		gigachefIDs = tmpGigachefIDs
	}
	return postIDs, gigachefIDs, distances, nil
}

func getSortByDateQuery(latitude float32, longitude float32, radius int, readyTime time.Time, descending bool, startLimit int, endLimit int) string {
	var readyDatetimeOrder, readyWhere string
	if descending {
		readyWhere = "'2014-04-01 00:00:00' AND " + readyTime.Format(time.RFC3339)
		readyDatetimeOrder = "DESC"
	} else {
		readyWhere = readyTime.Format(time.RFC3339) + " AND '4000-04-01 00:00:00'"
		readyDatetimeOrder = "ASC"
	}
	return fmt.Sprintf(sortByDate, latitude, longitude, latitude, readyWhere, radius, readyDatetimeOrder, startLimit, endLimit)
}

func init() {
	var err error
	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = "root@/gigamunch"
	} else {
		projectID := os.Getenv("PROJECTID")
		if projectID != "gigamunch-omninexus-dev" && projectID != "gigamunch-omninexus" {
			log.Fatalln("PROJECTID env variable not set")
		}
		connectionString = fmt.Sprintf("root@cloudsql(" + projectID + ":gigasql)/gigamunch")
	}
	mysqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Couldn't connect to mysql database")
	}
}

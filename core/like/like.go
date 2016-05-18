package like

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"

	"appengine"
)

const (
	datetimeFormat   = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"
	insertStatement  = "INSERT INTO `like` (created_datetime, user_id, item_id) VALUES ('%s','%s',%d)"
	increaseNumLikes = "INSERT INTO `num_likes` values (%d, 1) ON DUPLICATE KEY UPDATE num=num+1"
	decreaseNumLikes = "UPDATE `num_likes` SET num=num-1 where item_id=%d"
	deleteStatement  = "DELETE FROM `like` WHERE user_id='%s' AND item_id=%d"
	// selectByUserID   = "SELECT item_id FROM `like` WHERE user_id=? ORDER BY item_id ASC"
)

var (
	connectOnce = sync.Once{}
	mysqlDB     *sql.DB
	errSQLDB    = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	errBuffer   = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "An unknown error occured."}
)

// Client is the client for likes
type Client struct {
	ctx context.Context
}

// New returns a new Client for user likes
func New(ctx context.Context) *Client {
	connectOnce.Do(connectSQL)
	return &Client{ctx: ctx}
}

// Like likes an item for a user
func (c *Client) Like(userID string, itemID int64) error {
	if userID == "" {
		return nil
	}
	if itemID == 0 {
		return nil
	}
	st := fmt.Sprintf(insertStatement,
		time.Now().UTC().Format(datetimeFormat),
		userID,
		itemID)
	result, err := mysqlDB.Exec(st)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return nil
		}
		return errSQLDB.WithError(err).Wrapf("cannot execute insert statement(%s)", st)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to get rows affected")
	}
	switch rowsAffected {
	case 0:
		return nil // nothing liked
	case 1:
		break
	default:
		return errSQLDB.Wrapf("num rows affected(%d) is not 0 or 1", rowsAffected)
	}
	st = fmt.Sprintf(increaseNumLikes, itemID)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("cannot execute increaseNumLikes statement(%s)", st)
	}
	return nil
}

// Unlike unlikes an item for a user
func (c *Client) Unlike(userID string, itemID int64) error {
	if userID == "" {
		return nil
	}
	if itemID == 0 {
		return nil
	}
	st := fmt.Sprintf(deleteStatement, userID, itemID)
	result, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("cannot execute delete statement(%s)", st)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to get rows affected")
	}
	switch rowsAffected {
	case 0:
		return nil // nothing unliked
	case 1:
		break
	default:
		return errSQLDB.Wrapf("num rows affected(%d) is not 0 or 1", rowsAffected)
	}
	st = fmt.Sprintf(decreaseNumLikes, itemID)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("cannot execute delete statement(%s)", st)
	}
	return nil
}

// LikesItems returns an array that states if the muncher likes the item or not
func (c *Client) LikesItems(userID string, items []int64) ([]bool, []int, error) {
	likesItem := make([]bool, len(items))
	numLikes := make([]int, len(items))
	if len(items) == 0 {
		return likesItem, numLikes, nil
	}
	statement, err := buildLikesItemsStatement(userID, items)
	if err != nil {
		return likesItem, numLikes, errors.Wrap("failed to build like item statment", err)
	}
	rows, err := mysqlDB.Query(statement)
	if err != nil {
		return nil, nil, errSQLDB.WithError(err).Wrap("cannot query following statement: " + statement)
	}
	defer handleCloser(c.ctx, rows)

	var tmpItemID int64
	var tmpUserID string
	var tmpNumLike int
	for rows.Next() {
		err = rows.Scan(&tmpItemID, &tmpUserID, &tmpNumLike)
		if err != nil {
			return nil, nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
		}
		for i := range items {
			if items[i] == tmpItemID {
				if tmpUserID == "u" {
					likesItem[i] = true
				} else {
					numLikes[i] = tmpNumLike
				}
			}
		}
	}
	return likesItem, numLikes, nil
}

func buildLikesItemsStatement(userID string, items []int64) (string, error) {
	if len(items) == 0 {
		panic("items legth is 0")
	}
	// -- old
	// SELECT item_id, user_id, count(item_id) AS cnt FROM (( SELECT item_id, user_id FROM `like` WHERE (item_id=0 OR item_id=1)
	// ORDER BY user_id='user1' desc) AS s) group by item_id
	// -- new
	// SELECT item_id, user_id, 0 FROM `like` WHERE user_id='%s' and (%s)
	// UNION ALL SELECT item_id, '', num FROM `num_likes` WHERE (%s)
	var err error
	var buffer bytes.Buffer
	for i := range items[1:] {
		_, err = buffer.WriteString(fmt.Sprintf(" OR item_id=%d", items[i+1]))
		if err != nil {
			return "", errBuffer.WithError(err)
		}
	}
	itemIDStatement := fmt.Sprintf("item_id=%d %s", items[0], buffer.String())
	st := fmt.Sprintf("SELECT item_id, 'u', 0 FROM `like` WHERE user_id='%s' and (%s) UNION ALL SELECT item_id, '', num FROM `num_likes` WHERE (%s)",
		userID, itemIDStatement, itemIDStatement)
	return st, nil
}

// GetUserLikes returns an array of item ids the user likes
// func (c *Client) GetUserLikes(userID string) ([]int64, error) {
// 	rows, err := mysqlDB.Query(selectByUserID, userID)
// 	if err != nil {
// 		return nil, errSQLDB.WithError(err).Wrap("cannot query by user id")
// 	}
// 	defer handleCloser(c.ctx, rows)
// 	var id int64
// 	var ids []int64
// 	for rows.Next() {
// 		err = rows.Scan(&id)
// 		if err != nil {
// 			return nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
// 		}
// 		ids = append(ids, id)
// 	}
// 	return ids, nil
// }

func connectSQL() {
	var err error
	// TODO switch to config
	var connectionString string

	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = "root@/gigamunch"
	} else {
		projectID := os.Getenv("PROJECTID")
		if projectID == "" {
			log.Fatal("PROJECTID env variable is not set.")
		}
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

package like

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
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
	datetimeFormat  = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"
	insertStatement = "INSERT INTO `like` (created_datetime, user_id, item_id) VALUES (?,?,?)"
	deleteStatement = "DELETE FROM `like` WHERE user_id=? AND item_id=?"
	selectByUserID  = "SELECT item_id FROM `like` WHERE user_id=? ORDER BY item_id ASC"
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
	_, err := mysqlDB.Exec(insertStatement,
		time.Now().UTC().Format(datetimeFormat),
		userID,
		itemID)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("cannot execute insert statement")
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
	_, err := mysqlDB.Exec(deleteStatement,
		userID,
		itemID)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("cannot execute delete statement")
	}
	return nil
}

// LikesItems returns an array that states if the muncher likes the item or not
func (c *Client) LikesItems(userID string, items []int64) ([]bool, []int, error) {
	likesItem := make([]bool, len(items))
	numLikes := make([]int, len(items))
	if userID == "" {
		return likesItem, numLikes, nil
	}
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
				if tmpUserID == "" || tmpUserID == userID {
					likesItem[i] = true
				}
				numLikes[i] = tmpNumLike
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
	// SELECT item_id, '', COUNT(item_id) FROM `like` WHERE user_id='user1' and (item_id=0 OR item_id=1)
	// UNION ALL
	// SELECT item_id, user_id, COUNT(item_id) FROM `like` WHERE item_id=0 OR item_id=1 GROUP BY item_id;
	var err error
	var buffer bytes.Buffer
	for i := range items {
		if i != 0 {
			_, err = buffer.WriteString(fmt.Sprintf(" OR item_id=%d", items[i]))
			if err != nil {
				return "", errBuffer.WithError(err)
			}
		}
	}
	itemIDStatement := fmt.Sprintf("item_id=%d %s", items[0], buffer.String())
	st := fmt.Sprintf("SELECT item_id, '', COUNT(item_id) FROM `like` WHERE user_id='%s' and (%s) UNION ALL SELECT item_id, user_id, COUNT(item_id) FROM `like` WHERE %s GROUP BY item_id",
		userID, itemIDStatement, itemIDStatement)
	return st, nil
}

// GetUserLikes returns an array of item ids the user likes
func (c *Client) GetUserLikes(userID string) ([]int64, error) {
	rows, err := mysqlDB.Query(selectByUserID, userID)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("cannot query by user id")
	}
	defer handleCloser(c.ctx, rows)
	var id int64
	var ids []int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
		}
		ids = append(ids, id)
	}
	return ids, nil
}

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

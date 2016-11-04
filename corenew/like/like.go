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

	"google.golang.org/appengine"
)

const (
	datetimeFormat                     = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"
	insertStatement                    = "INSERT INTO `likes` (created_datetime, user_id, item_id, menu_id, cook_id) VALUES ('%s','%s',%d, %d, '%s')"
	deleteStatement                    = "DELETE FROM `likes` WHERE user_id='%s' AND item_id=%d"
	selectNumLikesStatement            = "SELECT item_id, COUNT(item_id) FROM `likes` WHERE %s GROUP BY item_id"
	selectNumLikesWithMenuIDStatement  = "SELECT item_id, menu_id, user_id, COUNT(item_id) FROM (SELECT user_id, item_id, menu_id FROM `likes` WHERE %s ORDER BY CASE WHEN user_id='%s' THEN 1 ELSE 2 END) as l GROUP BY item_id"
	selectNumLikesAndHasLikedStatement = "SELECT item_id, user_id, COUNT(item_id) FROM (SELECT user_id, item_id FROM `likes` WHERE %s ORDER BY CASE WHEN user_id='%s' THEN 1 ELSE 2 END) as l GROUP BY item_id"
	selectNumCookLikesStatement        = "SELECT COUNT(cook_id) FROM `likes` WHERE cook_id='%s'"
	// selectNumMenuLikesStatement        = "SELECT COUNT(menu_id) FROM `likes` WHERE menu_id=%d"
	// selectByUserID   = "SELECT item_id FROM `like` WHERE user_id=? ORDER BY item_id ASC"
)

var (
	connectOnce         = sync.Once{}
	mysqlDB             *sql.DB
	errSQLDB            = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	errBuffer           = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "An unknown error occured."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
)

// Client is the client for likes
type Client struct {
	ctx context.Context
}

// New returns a new Client for user likes
func New(ctx context.Context) *Client {
	connectOnce.Do(func() {
		connectSQL(ctx)
	})
	return &Client{ctx: ctx}
}

// Like likes an item for a user
func (c *Client) Like(userID string, itemID, menuID int64, cookID string) error {
	if userID == "" {
		return nil
	}
	if itemID == 0 {
		return nil
	}
	st := fmt.Sprintf(insertStatement,
		time.Now().UTC().Format(datetimeFormat), userID, itemID, menuID, cookID)
	_, err := mysqlDB.Exec(st)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return nil
		}
		return errSQLDB.WithError(err).Wrapf("cannot execute insert statement(%s)", st)
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
	_, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("cannot execute delete statement(%s)", st)
	}
	return nil
}

// GetNumCookLikes returns the number of likes for the cookID.
func (c *Client) GetNumCookLikes(cookID string) (int32, error) {
	if cookID == "" {
		return 0, nil
	}
	// create statement
	st := fmt.Sprintf(selectNumCookLikesStatement, cookID)
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return 0, errSQLDB.WithError(err).Wrap("cannot query following statement: " + st)
	}
	defer handleCloser(c.ctx, rows)
	var numLikes int32
	for rows.Next() {
		err = rows.Scan(&numLikes)
		if err != nil {
			return 0, errSQLDB.WithError(err).Wrap("cannot scan rows")
		}
	}
	return numLikes, nil
}

// GetNumLikes returns the number of likes for each item in the array
func (c *Client) GetNumLikes(items []int64) ([]int, error) {
	numLikes := make([]int, len(items))
	if len(items) == 0 {
		return numLikes, nil
	}
	// create statement
	statement, err := buildGetNumLikesStatement(items)
	if err != nil {
		return nil, errors.Wrap("failed to build get num likes statment", err)
	}
	rows, err := mysqlDB.Query(statement)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("cannot query following statement: " + statement)
	}
	defer handleCloser(c.ctx, rows)
	var tmpItemID int64
	var tmpNumLike int
	for rows.Next() {
		err = rows.Scan(&tmpItemID, &tmpNumLike)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
		}
		for i := range items {
			if items[i] == tmpItemID {
				numLikes[i] = tmpNumLike
			}
		}
	}
	return numLikes, nil
}

func buildGetNumLikesStatement(items []int64) (string, error) {
	if len(items) == 0 {
		return "", errInvalidParameter.Wrap("items length is 0")
	}
	var err error
	var buffer bytes.Buffer
	for i := range items[1:] {
		_, err = buffer.WriteString(fmt.Sprintf(" OR item_id=%d", items[i+1]))
		if err != nil {
			return "", errBuffer.WithError(err)
		}
	}
	itemIDStatement := fmt.Sprintf("item_id=%d %s", items[0], buffer.String())
	st := fmt.Sprintf(selectNumLikesStatement, itemIDStatement)
	return st, nil
}

// GetNumLikesWithMenuID returns the likesItem, numLikes, menuID, error.
func (c *Client) GetNumLikesWithMenuID(userID string, items []int64) ([]bool, []int32, []int64, error) {
	likesItem := make([]bool, len(items))
	numLikes := make([]int32, len(items))
	menuIDs := make([]int64, len(items))
	if len(items) == 0 {
		return likesItem, numLikes, menuIDs, nil
	}
	// create statement
	st, err := buildLikesWithMenuIDStatement(userID, items)
	if err != nil {
		return likesItem, numLikes, menuIDs, errors.Wrap("failed to build likes item statement", err)
	}
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return likesItem, numLikes, menuIDs, errSQLDB.WithError(err).Wrap("cannot query following statement: " + st)
	}
	defer handleCloser(c.ctx, rows)
	var tmpMenuID int64
	var tmpItemID int64
	var tmpUserID string
	var tmpNumLike int32
	for rows.Next() {
		err = rows.Scan(&tmpItemID, &tmpMenuID, &tmpUserID, &numLikes)
		if err != nil {
			return likesItem, numLikes, menuIDs, errSQLDB.WithError(err).Wrap("cannot scan rows")
		}
		for i := range items {
			if items[i] == tmpItemID {
				if tmpUserID == userID {
					likesItem[i] = true
				}
				numLikes[i] = tmpNumLike
				menuIDs[i] = tmpMenuID
			}
		}
	}
	return likesItem, numLikes, menuIDs, nil
}

func buildLikesWithMenuIDStatement(userID string, items []int64) (string, error) {
	if len(items) == 0 {
		return "", errInvalidParameter.Wrap("items length is 0")
	}
	var err error
	var buffer bytes.Buffer
	for i := range items[1:] {
		_, err = buffer.WriteString(fmt.Sprintf(" OR item_id=%d", items[i+1]))
		if err != nil {
			return "", errBuffer.WithError(err)
		}
	}
	itemIDStatement := fmt.Sprintf("item_id=%d %s", items[0], buffer.String())
	st := fmt.Sprintf(selectNumLikesWithMenuIDStatement, itemIDStatement, userID)
	return st, nil
}

// LikesItems returns an array that states if the muncher likes the item or not
func (c *Client) LikesItems(userID string, items []int64) ([]bool, []int32, error) {
	likesItem := make([]bool, len(items))
	numLikes := make([]int32, len(items))
	if len(items) == 0 {
		return likesItem, numLikes, nil
	}
	statement, err := buildLikesItemsStatement(userID, items)
	if err != nil {
		return likesItem, numLikes, errors.Wrap("failed to build likes item statement", err)
	}
	rows, err := mysqlDB.Query(statement)
	if err != nil {
		return nil, nil, errSQLDB.WithError(err).Wrap("cannot query following statement: " + statement)
	}
	defer handleCloser(c.ctx, rows)

	var tmpItemID int64
	var tmpUserID string
	var tmpNumLike int32
	for rows.Next() {
		err = rows.Scan(&tmpItemID, &tmpUserID, &tmpNumLike)
		if err != nil {
			return nil, nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
		}
		for i := range items {
			if items[i] == tmpItemID {
				if tmpUserID == userID {
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
		return "", errInvalidParameter.Wrap("items length is 0")
	}
	var err error
	var buffer bytes.Buffer
	for i := range items[1:] {
		_, err = buffer.WriteString(fmt.Sprintf(" OR item_id=%d", items[i+1]))
		if err != nil {
			return "", errBuffer.WithError(err)
		}
	}
	itemIDStatement := fmt.Sprintf("item_id=%d %s", items[0], buffer.String())
	st := fmt.Sprintf(selectNumLikesAndHasLikedStatement, itemIDStatement, userID)
	return st, nil
}

func connectSQL(ctx context.Context) {
	var err error
	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = "root@/gigamunch"
	} else {
		connectionString = os.Getenv("MYSQL_CONNECTION")
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

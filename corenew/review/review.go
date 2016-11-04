package review

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
	mysql "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	// datetimeFormat        = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"
	insertStatement       = "INSERT INTO review (cook_id,eater_id,eater_name,eater_photo_url,inquiry_id,item_id,item_name,item_photo_url,menu_id,rating,text) VALUES ('%s','%s','%s','%s',%d,%d,'%s','%s',%d,%d,'%s')"
	updateStatement       = "UPDATE review SET eater_name='%s', eater_photo_url='%s', rating=%d, text='%s', edited_datetime=NOW(), is_edited=1 WHERE id=%d"
	updateCookResponse    = "UPDATE review SET has_response=1, response_created_datetime=NOW(), response_text='%s' WHERE id=%d"
	selectReviewStatement = "SELECT id,cook_id,eater_id,eater_name,eater_photo_url,inquiry_id,item_id,item_name,item_photo_url,menu_id,created_datetime,rating,text,is_edited,edited_datetime,has_response,response_created_datetime,response_text FROM review WHERE id=%d %s"
	// TODO change to fn(created_datetime, item_id)
	selectReviewByCookID = "SELECT id,eater_id,eater_name,eater_photo_url,inquiry_id,item_id,item_name,item_photo_url,menu_id,created_datetime,rating,text,is_edited,edited_datetime,has_response,response_created_datetime,response_text FROM review WHERE cook_id='%s' ORDER BY created_datetime DESC LIMIT %d,%d"
	// selectCookReviews     = "SELECT "
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

// Post posts or updates a review.
func (c *Client) Post(user *types.User, id int64, cookID string, inquiryID, itemID int64, itemName, itemPhotoURL string, menuID int64, rating int32, text string) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errInvalidParameter.WithMessage("A rating has to be between 1 star and 5 stars.")
	}
	isNewReview := id == 0
	now := time.Now()
	review := &Review{
		CreatedDateTime: now,
		EditedDateTime:  now,
		ID:              id,
		CookID:          cookID,
		InquiryID:       inquiryID,
		ItemID:          itemID,
		ItemName:        itemName,
		ItemPhotoURL:    itemPhotoURL,
		MenuID:          menuID,
		Rating:          rating,
		Text:            getEscapedString(text),
	}
	if isNewReview {
		// insert review
		st := fmt.Sprintf(insertStatement, cookID, user.ID, user.Name, user.PhotoURL, inquiryID, itemID, itemName, itemPhotoURL, menuID, rating, text)
		results, err := mysqlDB.Exec(st)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrapf("failed execute %s", st)
		}
		id, err = results.LastInsertId()
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to results.LastInsertId()")
		}
		review.ID = id
		inquiryC := inquiry.New(c.ctx)
		err = inquiryC.SetReviewID(inquiryID, id)
		if err != nil {
			return nil, errors.Wrap("failed to inquiry.SetReviewID", err)
		}
	} else {
		// update review
		st := fmt.Sprintf(updateStatement, user.Name, user.PhotoURL, rating, text, id)
		results, err := mysqlDB.Exec(st)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrapf("failed execute %s", st)
		}
		rowsEffected, err := results.RowsAffected()
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to results.RowsAffected()")
		}
		if rowsEffected == 0 {
			return nil, errInvalidParameter.WithMessage("Invalid ReviewID.")
		}
		review.IsEdited = true
		review.EditedDateTime = time.Now()
	}
	return review, nil
}

// PostResponse updates a review with the cook's response.
func (c *Client) PostResponse(user *types.User, id int64, text string) error {
	if id == 0 {
		return errInvalidParameter.WithMessage("ReviewID cannot be 0.")
	}
	// update review
	st := fmt.Sprintf(updateCookResponse, getEscapedString(text), id)
	results, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("failed execute %s", st)
	}
	rowsEffected, err := results.RowsAffected()
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to results.RowsAffected()")
	}
	if rowsEffected == 0 {
		return errInvalidParameter.WithMessage("Invalid ReviewID.")
	}
	return nil
}

// func (c *Client) GetMultiByCookID() ([]*Review, error) {}

// GetByCookID gets reviews for a cook.
func (c *Client) GetByCookID(cookID string, itemID int64, startIndex, endIndex int32) ([]*Review, error) {
	var reviews []*Review
	if cookID == "" {
		return reviews, nil
	}
	st := fmt.Sprintf(selectReviewByCookID, cookID, startIndex, endIndex)
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrapf("failed to execute %s", st)
	}
	defer handleCloser(c.ctx, rows)
	for rows.Next() {
		review := new(Review)
		review.CookID = cookID
		var createdNulltime mysql.NullTime
		var editedNulltime mysql.NullTime
		var responseCreatedNulltime mysql.NullTime
		var text sql.NullString
		var responseText sql.NullString
		err = rows.Scan(&review.ID, &review.EaterID, &review.EaterName, &review.EaterPhotoURL,
			&review.InquiryID, &review.ItemID, &review.ItemName, &review.ItemPhotoURL, &review.MenuID, &createdNulltime,
			&review.Rating, &text, &review.IsEdited, &editedNulltime,
			&review.HasResponse, &responseCreatedNulltime, &responseText)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
		}
		if createdNulltime.Valid {
			review.CreatedDateTime = createdNulltime.Time
		}
		if editedNulltime.Valid {
			review.EditedDateTime = editedNulltime.Time
		}
		if responseCreatedNulltime.Valid {
			review.ResponseCreatedDateTime = responseCreatedNulltime.Time
		}
		if text.Valid {
			review.Text = getUnescapedString(text.String)
		}
		if responseText.Valid {
			review.ResponseText = getUnescapedString(responseText.String)
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func getEscapedString(s string) string {
	return strings.Replace(s, "'", "\\'", -1)
}

func getUnescapedString(s string) string {
	return strings.Replace(s, "\\'", "'", -1)
}

func getSelectReviewStatement(ids []int64) (string, error) {
	if len(ids) == 0 {
		return "", errInvalidParameter
	}
	var err error
	var buffer bytes.Buffer
	for _, v := range ids[1:] {
		_, err = buffer.WriteString(fmt.Sprintf(" OR id=%d", v))
		if err != nil {
			return "", errBuffer.WithError(err)
		}
	}
	st := fmt.Sprintf(selectReviewStatement, ids[0], buffer.String())
	return st, nil
}

// GetMultiByID gets multiple reviews by their ids. Reviews in the array might be nil if they are not found.
func (c *Client) GetMultiByID(ids []int64) ([]*Review, error) {
	reviews := make([]*Review, len(ids))
	if len(ids) != 0 {
		st, err := getSelectReviewStatement(ids)
		if err != nil {
			return nil, err
		}
		rows, err := mysqlDB.Query(st)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrapf("failed to execute %s", st)
		}
		defer handleCloser(c.ctx, rows)

		for rows.Next() {
			review := new(Review)
			var createdNulltime mysql.NullTime
			var editedNulltime mysql.NullTime
			var responseCreatedNulltime mysql.NullTime
			var text sql.NullString
			var responseText sql.NullString
			err = rows.Scan(&review.ID, &review.CookID, &review.EaterID, &review.EaterName, &review.EaterPhotoURL,
				&review.InquiryID, &review.ItemID, &review.ItemName, &review.ItemPhotoURL, &review.MenuID, &createdNulltime,
				&review.Rating, &text, &review.IsEdited, &editedNulltime,
				&review.HasResponse, &responseCreatedNulltime, &responseText)
			if err != nil {
				return nil, errSQLDB.WithError(err).Wrap("cannot scan rows")
			}
			if createdNulltime.Valid {
				review.CreatedDateTime = createdNulltime.Time
			}
			if editedNulltime.Valid {
				review.EditedDateTime = editedNulltime.Time
			}
			if responseCreatedNulltime.Valid {
				review.ResponseCreatedDateTime = responseCreatedNulltime.Time
			}
			if text.Valid {
				review.Text = getUnescapedString(text.String)
			}
			if responseText.Valid {
				review.ResponseText = getUnescapedString(responseText.String)
			}
			for i := range ids {
				if ids[i] == review.ID {
					reviews[i] = review
					break
				}
			}
		}
	}
	for i := range reviews {
		if reviews[i] == nil {
			reviews[i] = new(Review)
		}
	}
	return reviews, nil
}

// Get gets a review.
func (c *Client) Get(id int64) (*Review, error) {
	if id == 0 {
		return new(Review), nil
	}
	reviews, err := c.GetMultiByID([]int64{id})
	if err != nil {
		return nil, err
	}
	return reviews[0], err
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

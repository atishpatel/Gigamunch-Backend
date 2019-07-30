package discount

import (
	"context"
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	"github.com/jmoiron/sqlx"
)

const (
	// DateFormat is the expected format of date in activity.
	DateFormat                   = "2006-01-02" // "Jan 2, 2006"
	selectDiscount               = "SELECT * FROM discount WHERE id=?"
	selectDiscountByUserID       = "SELECT * FROM discount WHERE user_id=?"
	selectUnusedDiscountByUserID = "SELECT * FROM discount WHERE date_used='0000-00-00' AND user_id=?"
	updateUsed                   = "UPDATE discount set date_used=? WHERE id=?"
	insertStatement              = "INSERT INTO discount (user_id,email,first_name,last_name,discount_amount,discount_percent) VALUES (:user_id,:email,:first_name,:last_name,:discount_amount,:discount_percent)"
)

var (
	errSQL              = errors.InternalServerError
	errInternal         = errors.InternalServerError
	errBadRequest       = errors.BadRequestError
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	db         common.DB
	sqlDB      *sqlx.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	if sqlC == nil {
		return nil, errInternal.Annotate("failed to get sql client")
	}
	if dbC == nil {
		return nil, fmt.Errorf("failed to get db")
	}
	if serverInfo == nil {
		return nil, errInternal.Annotate("failed to get server info")
	}
	return &Client{
		ctx:        ctx,
		log:        log,
		sqlDB:      sqlC,
		db:         dbC,
		serverInfo: serverInfo,
	}, nil
}

// Get gets a discount.
func (c *Client) Get(id int64) (*Discount, error) {
	var err error
	dst := &Discount{}
	err = c.sqlDB.GetContext(c.ctx, dst, selectDiscount, id)
	if err != nil {
		return nil, errSQL.WithError(err).Annotate("failed to selectDiscount")
	}
	return dst, nil
}

// GetAllForUser gets a list of discounts for a user.
func (c *Client) GetAllForUser(userID string) ([]*Discount, error) {
	dsts := []*Discount{}
	err := c.sqlDB.SelectContext(c.ctx, &dsts, selectDiscountByUserID, userID)
	if err != nil {
		return nil, errSQL.WithError(err).Annotate("failed to selectDiscountByUserID")
	}
	return dsts, nil
}

// GetUnusedUserDiscount gets most relevant unused discount.
func (c *Client) GetUnusedUserDiscount(userID string) (*Discount, error) {
	dsts := []*Discount{}
	err := c.sqlDB.SelectContext(c.ctx, &dsts, selectUnusedDiscountByUserID, userID)
	if err != nil {
		return nil, errSQL.WithError(err).Annotate("failed to selectUnusedDiscountByUserID")
	}
	if len(dsts) < 1 {
		return nil, nil
	}
	return dsts[0], nil
}

// Used marks a discount as used.
func (c *Client) Used(id int64, date *time.Time) error {
	if date == nil {
		now := time.Now()
		date = &now
	}
	var err error
	_, err = c.sqlDB.ExecContext(c.ctx, updateUsed, date.Format(DateFormat), id)
	if err != nil {
		return errSQL.WithError(err).Annotate("failed to updateUsed")
	}
	return nil
}

// CreateReq is a request for create
type CreateReq struct {
	UserID          string  `json:"user_id" db:"user_id"`
	Email           string  `json:"email" db:"email"`
	FirstName       string  `json:"first_name" db:"first_name"`
	LastName        string  `json:"last_name" db:"last_name"`
	DiscountAmount  float32 `json:"discount_amount" db:"discount_amount"`
	DiscountPercent int8    `json:"discount_percent" db:"discount_percent"`
}

func (req *CreateReq) validate() error {
	if req.UserID == "" {
		return errBadRequest.WithMessage("UserID cannot be empty.")
	}
	if req.DiscountAmount < .01 && req.DiscountPercent == 0 {
		return errBadRequest.WithMessage("There must be a discount.")
	}
	if req.DiscountPercent > 100 {
		return errBadRequest.WithMessage("Invalid discount.")
	}
	return nil
}

// Create creates a discount for a user.
func (c *Client) Create(req *CreateReq) error {
	err := req.validate()
	if err != nil {
		return err
	}
	_, err = c.sqlDB.NamedExecContext(c.ctx, insertStatement, req)
	if err != nil {
		return errSQL.WithError(err).Annotate("failed to insertStatement")
	}
	return nil
}

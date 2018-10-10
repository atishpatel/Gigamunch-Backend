package activity

import (
	"context"
	"fmt"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
)

const (
	dateFormat                 = "2006-01-02" // "Jan 2, 2006"
	selectActivityStatement    = "SELECT * FROM activity WHERE date=? AND email=?"
	selectAllActivityStatement = "SELECT * FROM activity ORDER BY date DESC LIMIT ?"
	updateRefundedStatement    = "UPDATE activity SET refunded_dt=NOW(),refunded=1,refund_transaction_id=?,refunded_amount=? WHERE date=? AND email=?"
)

// Errors
var (
	errDatastore  = errors.InternalServerError
	errInternal   = errors.InternalServerError
	errBadRequest = errors.BadRequestError
	errSQLDB      = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
)

// Client is a client for manipulating activity.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	sqlDB      *sqlx.DB
	db         common.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	// if sqlC == nil {
	// 	return nil, fmt.Errorf("sqlDB cannot be nil for sub")
	// }
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

// Get gets an activity.
func (c *Client) Get(date time.Time, email string) (*Activity, error) {
	act := &Activity{}
	err := c.sqlDB.GetContext(c.ctx, act, selectActivityStatement, date.Format(dateFormat), email)
	if err != nil {
		return nil, errors.Annotate(err, "failed to selectActivity")
	}
	return act, nil
}

// GetAll gets a list of activity.
func (c *Client) GetAll(limit int) ([]*Activity, error) {
	if limit <= 0 {
		limit = 1000
	}
	acts := []*Activity{}
	err := c.sqlDB.SelectContext(c.ctx, &acts, selectAllActivityStatement, limit)
	if err != nil {
		return nil, errors.Annotate(err, "failed to selectAllActivityStatement")
	}
	return acts, nil
}

// Create creates an activity entry.
func (c *Client) Create(date time.Time, email string, servings, vegServings int8, amount float32, paymentMethodToken, customerID string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Setup(date, email, servings, vegServings, amount, 0, paymentMethodToken, customerID)
}

// Process processes an actvity.
func (c *Client) Process(date time.Time, email string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Process(date, email)
}

// Discount gives an discount to an actvitiy.
func (c *Client) Discount(date time.Time, email string, discountAmount float32, discountPercent int8, overrideDiscounts bool) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Discount(date, email, discountAmount, discountPercent, overrideDiscounts)
}

// Refund refunds and skips a subscriber.
func (c *Client) Refund(date time.Time, email string, amount float32, precent int32) error {
	if amount > 0 && precent > 0 {
		return errBadRequest.WithMessage("Only amount or percent can be specified.")
	}
	// Get activity
	act, err := c.Get(date, email)
	if err != nil {
		return errors.Wrap("failed to Get", err)
	}
	if precent > 0 {
		amount = act.Amount * float32(precent) / 100
	}
	// Refund card
	paymentC, err := payment.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to payment.New")
	}
	rID, err := paymentC.RefundSale(act.TransactionID, amount)
	if err != nil {
		return errors.Wrap("failed to payment.RefundSale", err)
	}
	// log
	c.log.Refund(act.UserID, act.Email, act.Date, act.Amount, amount, rID)
	c.log.Infof(c.ctx, "Refunding Customer(%s) on transaction(%s): refundID(%s)", act.CustomerID, act.TransactionID, rID)
	// Update actvity
	_, err = c.sqlDB.ExecContext(c.ctx, updateRefundedStatement, rID, amount, date.Format(dateFormat), email)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateRefundedAndSkipSubLogStatement")
	}
	return nil
}

// RefundAndSkip refunds and skips a subscriber.
func (c *Client) RefundAndSkip(date time.Time, email string, amount float32, precent int32) error {
	// Refund
	err := c.Refund(date, email, amount, precent)
	if err != nil {
		return errors.Annotate(err, "failed to Refund")
	}
	// Skip
	err = c.Skip(date, email, "Refunded")
	if err != nil {
		return errors.Annotate(err, "failed to Skip")
	}
	return nil
}

// Skip skips a subscriber for an activity.
func (c *Client) Skip(date time.Time, email, reason string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Skip(date, email, reason)
}

// Unskip unskips a subscriber for an activity.
func (c *Client) Unskip(date time.Time, email string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Unskip(date, email)
}

// Paid sets activity to paid.
// func (c *Client) Paid() error {
// 	// TODO: Reimplement
// 	suboldC := subold.NewWithLogging(c.ctx, c.log)
// 	return nil
// }

// ChangeServings for a day
func (c *Client) ChangeServings(date time.Time, email string, servings int8, amount float32) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.ChangeServings(date, email, servings, amount)
}

// Free f
func (c *Client) Free(date time.Time, email string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Free(date, email)
}

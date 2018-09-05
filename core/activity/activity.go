package activity

import (
	"context"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
)

var ()

// Errors
var (
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
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
	// if dbC == nil {
	// 	return nil, fmt.Errorf("db cannot be nil for sub")
	// }
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

// Create creates an activity entry.
func (c *Client) Create(date time.Time, email string, servings int8, amount float32, paymentMethodToken, customerID string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Setup(date, email, servings, amount, 0, paymentMethodToken, customerID)
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

// RefundAndSkip refunds and skips a subscriber.
func (c *Client) RefundAndSkip(date time.Time, email string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.RefundAndSkip(date, email)
}

// Refund refunds a subscriber.
// func (c *Client) Refund() error {
// 	// TODO: Reimplement
// 	suboldC := subold.NewWithLogging(c.ctx, c.log)
// 	return nil
// }

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

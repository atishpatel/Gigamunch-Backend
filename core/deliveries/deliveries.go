package deliveries

import (
	"context"
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	// import driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	standAppEngine bool
	sql            *sqlx.DB
	db             common.DB
	projID         string
)

var (
	errSQL             = errors.InternalServerError
	errInternal        = errors.InternalServerError
	errInvalidArgument = errors.BadRequestError
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx context.Context
	log *logging.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client) (*Client, error) {
	var err error
	if standAppEngine {
		err = setup(ctx)
		if err != nil {
			return nil, err
		}
	}
	if sql == nil {
		return nil, errInternal.Annotate("setup not called")
	}
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx: ctx,
		log: log,
	}, nil
}

// Get gets a subscriber.
func (c *Client) Get(date, driverEmail string) ([]*Delivery, error) {
	if date == "" {
		return nil, errInvalidArgument.WithMessage("Date cannot be empty")
	}
	deliveries := []*Delivery{}
	var err error
	if driverEmail == "" {
		err = sql.SelectContext(c.ctx, deliveries, "SELECT * FROM deliveries WHERE date=?", date)
	} else {
		err = sql.SelectContext(c.ctx, deliveries, "SELECT * FROM deliveries WHERE date=? AND driver_email=?", date, driverEmail)
	}
	if err != nil {
		return nil, errSQL.WithError(err).Annotate("failed to select deliveries")
	}
	return deliveries, nil
}

// Update updates a subscriber.
func (c *Client) Update(deliveries []*Delivery) error {
	if len(deliveries) == 0 {
		return nil
	}
	// build statement
	statement := "INSERT INTO deliveries (date,sub_email,updated_dt,driver_id,driver_email,driver_name,sub_id,delivery_order,success,fail) VALUES "
	vals := []interface{}{}
	now := time.Now()
	for _, delivery := range deliveries {
		statement += "(?,?,?,?,?,?,?,?,?,?),"
		if delivery.DriverID == 0 {
			delivery.DriverID = -1
		}
		if delivery.SubID == 0 {
			delivery.SubID = -1
		}
		vals = append(vals, delivery.Date, delivery.SubEmail, now, delivery.DriverID, delivery.DriverEmail, delivery.SubID, delivery.Order, delivery.Success, delivery.Fail)
	}
	// remove end comma
	statement = statement[:len(statement)-2]
	statement += " ON DUPLICATE KEY UPDATE driver_email=VALUES(driver_email),driver_name=VALUES(driver_name),driver_id=VALUES(driver_id),delivery_order=VALUES(delivery_order),updated_dt=VALUES(updated_dt),success=VALUES(success),fail=VALUES(fail)"
	// execute statement
	_, err := sql.ExecContext(c.ctx, statement, vals...)
	if err != nil {
		return errSQL.WithError(err).Annotate("failed to update deliveries")
	}
	return nil
}

// TODO success func and fail func

// Setup sets up the logging package.
func Setup(ctx context.Context, standardAppEngine bool, projectID string, sqlC *sqlx.DB, dbC common.DB) error {
	var err error
	standAppEngine = standardAppEngine
	if !standAppEngine {
		err = setup(ctx)
		if err != nil {
			return err
		}
	}
	if sqlC == nil {
		return fmt.Errorf("sqlDB cannot be nil for sub")
	}
	sql = sqlC
	if dbC == nil {
		return fmt.Errorf("db cannot be nil for sub")
	}
	db = dbC
	projID = projectID
	return nil
}

func setup(ctx context.Context) error {
	return nil
}

package sub

import (
	"context"
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	"github.com/jmoiron/sqlx"
)

var (
	standAppEngine bool
	sql            *sqlx.DB
	db             common.DB
	projID         string
)

var (
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx context.Context
	log *logging.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context) (*Client, error) {
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
	log, ok := ctx.Value(common.LoggingKey).(*logging.Client)
	if !ok {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx: ctx,
		log: log,
	}, nil
}

// Get gets a subscriber.
func (c *Client) Get(id int64) (*Subscriber, error) {
	key := db.IDKey(c.ctx, kind, id)
	sub := new(Subscriber)
	err := db.Get(c.ctx, key, sub)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to get")
	}
	return sub, nil
}

// GetActive gets all active subscribers.
func (c *Client) GetActive(start, limit int) ([]*Subscriber, error) {
	var subs []*Subscriber
	_, err := db.QueryFilter(c.ctx, kind, start, limit, "Active=", true, subs)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryFilter")
	}
	return subs, nil
}

// Update updates a subscriber.
func (c *Client) Update(sub *Subscriber) error {
	key := db.IDKey(c.ctx, kind, sub.ID)
	_, err := db.Put(c.ctx, key, sub)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to put")
	}
	return nil
}

// SetupActivity updates a subscriber.
func (c *Client) SetupActivity(date time.Time) error {
	subs, err := c.GetActive(0, 1000)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to put")
	}
	// TODO: implement
	_ = subs
	return nil
}

// func (c *Client) Deactivate(id string) error {

// }

// func (c *Client) UpdatePaymentToken(subEmail string, paymentMethodToken string) error {
// }

// func (c *Client) ChangeServings(subEmail string, servings int8, vegetarian bool) error {
// 	// rememver to change tags
// }

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

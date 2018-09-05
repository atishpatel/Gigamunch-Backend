package sub

import (
	"context"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	"github.com/jmoiron/sqlx"
)

var (
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
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

// Get gets a subscriber.
func (c *Client) Get(id int64) (*Subscriber, error) {
	key := c.db.IDKey(c.ctx, kind, id)
	sub := new(Subscriber)
	err := c.db.Get(c.ctx, key, sub)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to get")
	}
	return sub, nil
}

// GetActive gets all active subscribers.
func (c *Client) GetActive(start, limit int) ([]*Subscriber, error) {
	var subs []*Subscriber
	_, err := c.db.QueryFilter(c.ctx, kind, start, limit, "Active=", true, subs)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryFilter")
	}
	return subs, nil
}

// GetByEmail gets a subscriber by email.

// func (c *Client) get(id int64, email string) (*Subscriber, error)

// GetByPhoneNumber gets a subscriber by phone number.

// GetHasSubscribed returns a list of all subscribers ever.

// ChangeServingsPermanently

// UpdatePaymentToken

// Update updates a subscriber.
func (c *Client) Update(sub *Subscriber) error {
	// TODO: log change

	key := c.db.IDKey(c.ctx, kind, sub.ID)
	_, err := c.db.Put(c.ctx, key, sub)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to put")
	}
	return nil
}

// Activate

// Deactivate deactivates an account
func (c *Client) Deactivate() error {
	// TODO: implement
	return nil
}

// Create

// CreateFromOldSub

// SetupActivities updates a subscriber.
func (c *Client) SetupActivities(date time.Time) error {
	subs, err := c.GetActive(0, 10000)
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

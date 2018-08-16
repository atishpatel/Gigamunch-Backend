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
	var err error
	if serverInfo.IsStandardAppEngine {
		err = setup(ctx)
		if err != nil {
			return nil, err
		}
	}
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
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

// Update updates a subscriber.
func (c *Client) Update(sub *Subscriber) error {
	key := c.db.IDKey(c.ctx, kind, sub.ID)
	_, err := c.db.Put(c.ctx, key, sub)
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

func setup(ctx context.Context) error {
	return nil
}

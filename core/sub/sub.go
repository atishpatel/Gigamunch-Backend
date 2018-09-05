package sub

import (
	"context"
	"regexp"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"

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
	// if sqlC == nil {
	// 	return nil, errInternal.Annotate("failed to get sql client")
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
// func (c *Client) GetByPhoneNumber() error {
// 	// TODO: implement
// 	return nil
// }

// GetHasSubscribed returns a list of all subscribers ever.
// func (c *Client) GetHasSubscribed() error {
// 	// TODO: implement
// 	return nil
// }

// ChangeServingsPermanently changes a subscriber's servings permanently.
func (c *Client) ChangeServingsPermanently(email string, servings int8, vegetarian bool) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.ChangeServingsPermanently(email, servings, vegetarian, c.serverInfo)
}

// UpdatePaymentToken updates a user payment method token.
func (c *Client) UpdatePaymentToken(email, paymentMethodToken string) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.UpdatePaymentToken(email, paymentMethodToken)
}

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

// Activate activates an account.
// func (c *Client) Activate() error {
// 	// TODO: implement
// 	suboldC := subold.NewWithLogging(c.ctx, c.log)
// 	return nil
// }

// Deactivate deactivates an account
func (c *Client) Deactivate(email string) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Cancel(email, c.log, c.serverInfo)
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
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.SetupSubLogs(date)
}

// IncrementPageCount is when a user just leaves their email.
func (c *Client) IncrementPageCount(email string, referralPageOpens int, referredPageOpens int) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	s, err := suboldC.GetSubscriber(email)
	if err != nil {
		return err
	}
	s.ReferralPageOpens += referralPageOpens
	s.ReferredPageOpens += referredPageOpens
	err = subold.Put(c.ctx, email, s)
	if err != nil {
		errDatastore.WithError(err).Annotate("failed to subold.put")
	}
	return nil
}

// GetCleanPhoneNumber takes a raw phone number and formats it to clean phone number.
func GetCleanPhoneNumber(rawNumber string) string {
	reg := regexp.MustCompile("[^0-9]+")
	cleanNumber := reg.ReplaceAllString(rawNumber, "")
	if len(cleanNumber) < 10 {
		return cleanNumber
	}
	cleanNumber = cleanNumber[len(cleanNumber)-10:]
	cleanNumber = cleanNumber[:3] + "-" + cleanNumber[3:6] + "-" + cleanNumber[6:]
	return cleanNumber
}

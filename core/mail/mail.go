package mail

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/drip-go"
	"google.golang.org/appengine/urlfetch"
)

// Tag is a tag applied to email subscribers.
type Tag string

func (t Tag) String() string {
	return string(t)
}

const (
	// LeftWebsiteEmail if they left email on website.
	LeftWebsiteEmail Tag = "LEFT_WEBSITE_EMAIL"
	// Customer if they are a customer.
	Customer Tag = "CUSTOMER"
	// Vegetarian if they are a vegetarian.
	Vegetarian Tag = "VEGETARIAN"
	// NonVegetarian if they a non-vegetarian.
	NonVegetarian Tag = "NON_VEGETARIAN"
	// TwoServings if they are 2 servings.
	TwoServings Tag = "TWO_SERVINGS"
	// FourServings if they are 4 servings.
	FourServings Tag = "FOUR_SERVINGS"
)

var (
	standAppEngine bool
	key            string
	acctID         string
)

var (
	errBadRequest = errors.BadRequestError
	errDrip       = errors.InternalServerError
	errInternal   = errors.InternalServerError
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx   context.Context
	log   *logging.Client
	dripC *drip.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context) (*Client, error) {
	var err error
	dripClient, err := drip.New(key, acctID)
	if err != nil {
		return nil, errInternal.Annotate("failed to get drip client")
	}
	if key == "" {
		return nil, errInternal.Annotate("setup not called or key is empty")
	}
	if standAppEngine {
		dripClient.HTTPClient = urlfetch.Client(ctx)
	}
	log, ok := ctx.Value(common.LoggingKey).(*logging.Client)
	if !ok {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx:   ctx,
		log:   log,
		dripC: dripClient,
	}, nil
}

// UserFields contain all the possible fields a user can have.
type UserFields struct {
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	FirstDeliveryDate time.Time `json:"first_delivery_date"`
}

// UpdateUser updates the user custom fields.
func (c *Client) UpdateUser(req *UserFields) error {
	resp, err := c.dripC.FetchSubscriber(req.Email)
	if err != nil {
		return errDrip.WithError(err).Annotate("failed to drip.FetchSubscriber")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.FetchSubscriber")
	}
	if len(resp.Subscribers) != 1 {
		return errBadRequest.Annotate("failed to find subscriber")
	}
	sub := resp.Subscribers[0]
	if req.FirstName != "" {
		sub.CustomFields["FIRST_NAME"] = req.FirstName
	}
	if req.LastName != "" {
		sub.CustomFields["LAST_NAME"] = req.LastName
	}
	if !req.FirstDeliveryDate.IsZero() {
		sub.CustomFields["FIRST_DELIVERY_DATE"] = dateString(req.FirstDeliveryDate)
	}
	return nil
}

// AddTag adds a tag to a customer. This often triggers a workflow.
func (c *Client) AddTag(email string, tag Tag) error {
	req := &drip.TagsReq{
		Tags: []drip.TagReq{
			drip.TagReq{
				Email: email,
				Tag:   tag.String(),
			},
		},
	}
	resp, err := c.dripC.TagSubscriber(req)
	if err != nil {
		return errDrip.WithError(err).Annotate("failed to drip.TagSubscriber")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.TagSubscriber")
	}
	return nil
}

// RemoveTag removes a tag from a customer. This often triggers a workflow.
func (c *Client) RemoveTag(email string, tag Tag) error {
	req := &drip.TagReq{
		Email: email,
		Tag:   tag.String(),
	}
	resp, err := c.dripC.RemoveSubscriberTag(req)
	if err != nil {
		return errDrip.WithError(err).Annotate("failed to drip.TagSubscriber")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.TagSubscriber")
	}
	return nil
}

// Setup sets up the logging package.
func Setup(ctx context.Context, standardAppEngine bool, apiKey, accountID string) error {
	standAppEngine = standardAppEngine
	key = apiKey
	acctID = accountID
	return nil
}

func dateString(t time.Time) string {
	var numString string
	var suffix string
	if t.Day() == 1 || t.Day() == 21 || t.Day() == 31 {
		suffix = "st"
	} else if t.Day() == 2 || t.Day() == 22 {
		suffix = "nd"
	} else if t.Day() == 3 || t.Day() == 23 {
		suffix = "rd"
	} else {
		suffix = "th"
	}
	numString = fmt.Sprintf("%d%s", t.Day(), suffix)
	return fmt.Sprintf("%s, %s %s", t.Weekday().String(), t.Month().String(), numString)
}

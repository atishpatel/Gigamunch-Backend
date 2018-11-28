package mail

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gopkg.in/mailgun/mailgun-go.v1"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/drip-go"
)

// TODO: Logging

var (
	standAppEngine      bool
	dripSubAPIKey       string
	dripSubAcctID       string
	dripMarketingAPIKey string
	dripMarketingAcctID string
	mailgunAPIKey       string
	mailgunPublicAPIKey string
	projID              string
)

var (
	errBadRequest = errors.BadRequestError
	errDrip       = errors.InternalServerError
	errInternal   = errors.InternalServerError
)

// Tag is a tag applied to email subscribers.
type Tag string

func (t Tag) String() string {
	return string(t)
}

const (
	// ==================
	// For Marketing drip
	// ==================

	// LeftWebsiteEmail if they left email on website.
	LeftWebsiteEmail Tag = "LEFT_WEBSITE_EMAIL"
	// ==================
	// Both drip
	// ==================

	// Subscribed is applied when the user first subscribed and is never removed.
	Subscribed Tag = "HAS_SUBSCRIBED"
	// Subscriber is applied when the user subscribers or when a user reactivates the account. It is removed when the user unsubscribes.
	Subscriber Tag = "SUBSCRIBER"
	// Deactivated is applied when the user suspends their account. It is removed when the user reactivates their account.
	Deactivated Tag = "DEACTIVATED"
	// ==================
	// For Sub drip
	// ==================

	// NonVegAndVeg if they are a half non-veg and half veg subscriber.
	NonVegAndVeg Tag = "NON_VEG_AND_VEG"
	// Vegetarian if they are a vegetarian.
	Vegetarian Tag = "VEGETARIAN"
	// NonVegetarian if they are a non-vegetarian.
	NonVegetarian Tag = "NON_VEGETARIAN"
	// TwoServings if they are 2 servings.
	TwoServings Tag = "TWO_SERVINGS"
	// FourServings if they are 4 servings.
	FourServings Tag = "FOUR_SERVINGS"
	// Gifted if tehy were given Gigamunch as a gift
	Gifted Tag = "GIFTED"
	// Dev if they are development server subscriber.
	Dev Tag = "DEV"
)

// GetPreviewEmailTag returns the tag that needs to be added to get the preview email based on date provided. Date should be date the person is recieving their meal.
func GetPreviewEmailTag(t time.Time) Tag {
	return Tag(t.Format("01/02/2006") + "_PREVIEW_EMAIL")
}

// GetCultureEmailTag returns the tag that needs to be added to get the culture email based on date provided. Date should be date the person is recieving their first meal.
func GetCultureEmailTag(t time.Time) Tag {
	return Tag(t.Format("01/02/2006") + "_CULTURE_EMAIL")
}

// GetReceivedJourneyTag returns the tag that needs to be added to customer based on how many meals they have.
func GetReceivedJourneyTag(numJourneys int) Tag {
	return Tag(fmt.Sprintf("RECEIVED_%d_JOURNEY", numJourneys))
}

// Client is a client for manipulating subscribers.
type Client struct {
	ctx            context.Context
	log            *logging.Client
	dripSubC       *drip.Client
	dripMarketingC *drip.Client
	mailgunC       mailgun.Mailgun
	serverInfo     *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, serverInfo *common.ServerInfo) (*Client, error) {
	var err error
	if dripSubAPIKey == "" {
		cnfg := config.GetConfig(ctx)
		dripSubAPIKey = cnfg.DripAPIKey
		dripSubAcctID = cnfg.DripAccountID
		dripMarketingAPIKey = cnfg.DripMarketingAPIKey
		dripMarketingAcctID = cnfg.DripMarketingAccountID
		mailgunAPIKey = cnfg.MailgunAPIKey
		mailgunPublicAPIKey = cnfg.MailgunPublicAPIKey
	}
	dripSubClient, err := drip.New(dripSubAPIKey, dripSubAcctID)
	if err != nil {
		return nil, errInternal.WithError(err).Annotate("failed to get drip client")
	}
	dripMarketingClient, err := drip.New(dripMarketingAPIKey, dripMarketingAcctID)
	if err != nil {
		return nil, errInternal.WithError(err).Annotate("failed to get drip client")
	}
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx:            ctx,
		log:            log,
		dripSubC:       dripSubClient,
		dripMarketingC: dripMarketingClient,
		serverInfo:     serverInfo,
	}, nil
}

// LeftEmail is when user leaves an email.
func (c *Client) LeftEmail(email, firstName, lastName string) error {
	req := &UserFields{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		AddTags:   []Tag{LeftWebsiteEmail},
	}
	return c.updateUser(req, c.dripMarketingC)
}

// SubActivated is when a subscriber account is activated.
func (c *Client) SubActivated(req *UserFields) error {
	var err error
	req.AddTags = append(req.AddTags, Subscriber, Subscribed)
	req.RemoveTags = append(req.RemoveTags, Deactivated)
	// For Sub Drip account
	err = c.updateUser(req, c.dripSubC)
	if err != nil {
		return err
	}
	// For Marketing Drip account
	err = c.updateUser(req, c.dripMarketingC)
	return err
}

// SubDeactivated is when a subscriber account is deactivated.
func (c *Client) SubDeactivated(req *UserFields) error {
	var err error
	req.AddTags = append(req.AddTags, Deactivated)
	req.RemoveTags = append(req.RemoveTags, Subscriber)
	// For Sub Drip account
	err = c.updateUser(req, c.dripSubC)
	if err != nil {
		return err
	}
	// For Marketing Drip account
	err = c.updateUser(req, c.dripMarketingC)
	return err
}

// UserFields contain all the possible fields a user can have.
type UserFields struct {
	Email             string    `json:"email"`
	NewEmail          string    `json:"new_email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	FirstDeliveryDate time.Time `json:"first_delivery_date"`
	GifterName        string    `json:"gifter_name"`
	GifterEmail       string    `json:"gifter_email"`
	VegServings       int8      `json:"veg_servings"`
	NonVegServings    int8      `json:"non_veg_servings"`
	AddTags           []Tag     `json:"add_tags"`
	RemoveTags        []Tag     `json:"remove_tags"`
}

// UpdateUser updates the user custom fields.
func (c *Client) UpdateUser(req *UserFields) error {
	return c.updateUser(req, c.dripSubC)
}

func (c *Client) updateUser(req *UserFields, dripClient *drip.Client) error {
	if ignoreEmail(req.Email) {
		return nil
	}
	sub := drip.UpdateSubscriber{
		Email:        req.Email,
		NewEmail:     req.NewEmail,
		CustomFields: make(map[string]string),
	}
	if req.FirstName != "" {
		sub.CustomFields["FIRST_NAME"] = req.FirstName
	}
	if req.LastName != "" {
		sub.CustomFields["LAST_NAME"] = req.LastName
	}
	if !req.FirstDeliveryDate.IsZero() {
		sub.CustomFields["FIRST_DELIVERY_DATE"] = DateString(req.FirstDeliveryDate)
	}
	if req.GifterName != "" {
		sub.CustomFields["GIFTER_NAME"] = req.GifterName
	}
	if req.GifterEmail != "" {
		sub.CustomFields["GIFTER_EMAIL"] = req.GifterEmail
	}
	if req.VegServings > 0 && req.NonVegServings > 0 {
		req.AddTags = append(req.AddTags, Vegetarian, NonVegetarian, NonVegAndVeg)
	} else if req.VegServings > 0 {
		req.AddTags = append(req.AddTags, Vegetarian)
		req.RemoveTags = append(req.RemoveTags, NonVegetarian, NonVegAndVeg)
	} else if req.NonVegServings > 0 {
		req.AddTags = append(req.AddTags, NonVegetarian)
		req.RemoveTags = append(req.RemoveTags, Vegetarian, NonVegAndVeg)
	}
	if len(req.AddTags) > 0 {
		for _, v := range req.AddTags {
			sub.Tags = append(sub.Tags, v.String())
		}
	}
	if len(req.RemoveTags) > 0 {
		for _, v := range req.RemoveTags {
			sub.RemoveTags = append(sub.RemoveTags, v.String())
		}
	}
	// Add Dev tag to customers in non-prod env.
	if !common.IsProd(c.serverInfo.ProjectID) {
		sub.Tags = append(sub.Tags, Dev.String())
	}
	dripReq := &drip.UpdateSubscribersReq{
		Subscribers: []drip.UpdateSubscriber{
			sub,
		},
	}
	resp, err := dripClient.UpdateSubscriber(dripReq)
	if err != nil {
		if strings.Contains(err.Error(), "<html>") {
			err = fmt.Errorf("drip returned an html page")
		}
		return errDrip.WithError(err).Annotate("failed to drip.FetchSubscriber")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.FetchSubscriber")
	}
	return nil
}

// AddTag adds a tag to a customer. This often triggers a workflow.
func (c *Client) AddTag(email string, tag Tag) error {
	if ignoreEmail(email) {
		return nil
	}
	req := &drip.TagsReq{
		Tags: []drip.TagReq{
			drip.TagReq{
				Email: email,
				Tag:   tag.String(),
			},
		},
	}
	resp, err := c.dripSubC.TagSubscriber(req)
	if err != nil {
		if strings.Contains(err.Error(), "<html>") {
			err = fmt.Errorf("drip returned an html page")
		}
		return errDrip.WithError(err).Annotate("failed to drip.TagSubscriber")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.TagSubscriber")
	}
	return nil
}

// RemoveTag removes a tag from a customer. This often triggers a workflow.
func (c *Client) RemoveTag(email string, tag Tag) error {
	if ignoreEmail(email) {
		return nil
	}
	req := &drip.TagReq{
		Email: email,
		Tag:   tag.String(),
	}
	resp, err := c.dripSubC.RemoveSubscriberTag(req)
	if err != nil {
		if strings.Contains(err.Error(), "<html>") {
			err = fmt.Errorf("drip returned an html page")
		}
		return errDrip.WithError(err).Annotate("failed to drip.TagSubscriber")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.TagSubscriber")
	}
	return nil
}

// AddBatchTags adds tags to emails. This often triggers a workflow.
func (c *Client) AddBatchTags(emails []string, tags []Tag) error {
	// TODO: batch limit is 1000 emails so update to split those into two request
	tagsString := make([]string, len(tags))
	for i, tag := range tags {
		tagsString[i] = tag.String()
	}
	subs := make([]drip.UpdateSubscriber, len(emails))
	i := 0
	for _, email := range emails {
		if !ignoreEmail(email) {
			subs[i].Email = email
			subs[i].Tags = tagsString
			i++
		}
	}
	if i == 0 || len(tags) == 0 {
		return nil
	}
	req := &drip.UpdateBatchSubscribersReq{
		Batches: []drip.SubscribersBatch{
			drip.SubscribersBatch{
				Subscribers: subs[:i],
			},
		},
	}
	resp, err := c.dripSubC.UpdateBatchSubscribers(req)
	if err != nil {
		if strings.Contains(err.Error(), "<html>") {
			err = fmt.Errorf("drip returned an html page")
		}
		return errDrip.WithError(err).Annotate("failed to drip.UpdateBatchSubscribers")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.UpdateBatchSubscribers")
	}
	c.log.Infof(c.ctx, "drip resp: code(%d) errors: %+v", resp.StatusCode, resp.Errors)
	return nil
}

// Send sends a plain text email.
func (c *Client) Send(from, subject, message string, to ...string) error {
	msg := mailgun.NewMessage(from, subject, message, to...)
	_, _, err := c.mailgunC.Send(msg)
	if err != nil {
		return errInternal.WithError(err).Wrap("failed to send mailgun email")
	}
	return nil
}

func ignoreEmail(email string) bool {
	if strings.Contains(email, "@test.com") || strings.Contains(email, "@apartment.com") {
		return true
	}
	return false
}

// DateString formates the date into a "Monday, January 1st" format.
func DateString(t time.Time) string {
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

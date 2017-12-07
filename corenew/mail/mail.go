package mail

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	drip "github.com/atishpatel/drip-go"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/appengine/urlfetch"
)

const (
	welcomeEmailSubject = "Welcome to Gigamunch! 🎉"
	welcomeEmailText    = `Hey %s!

We're excited to show you how delicious the world truly is, one unforgettable dinner at a time! 😋 

We work with local international cooks and prepare a dinner package that will make you feel like you actually went to that cook’s country! Every Monday we’ll deliver a dinner package that includes a culture guide to discover the story behind the meal, a playlist to listen to the music from that country, and of course, a delicious ready-to-eat dinner. 

If you have any dietary restrictions (vegetarian, peanut allergy, etc) please let us know here, and we'll do our best to accommodate your needs. 😊

You will get your first dinner package delivered %s! And guess what? It's totally free! 🤑 It’s our way of saying thank you for giving our delicious journey a whirl! We hope you’ll love it as much as we loved creating it (we think you will 😉).

See you soon!
- The Gigamunch Team`
	welcomeEmailHTML = `Hey %s!

We're excited to show you how delicious the world truly is, one unforgettable dinner at a time! 😋 

We work with local international cooks and prepare a dinner package that will make you feel like you actually went to that cook’s country! Every Monday we’ll deliver a dinner package that includes a culture guide to discover the story behind the meal, a playlist to listen to the music from that country, and of course, a delicious ready-to-eat dinner. 

If you have any dietary restrictions (vegetarian, peanut allergy, etc) please let us know here, and we'll do our best to accommodate your needs. 😊

You will get your first dinner package delivered %s! And guess what? It's totally free! 🤑 It’s our way of saying thank you for giving our delicious journey a whirl! We hope you’ll love it as much as we loved creating it (we think you will 😉).

See you soon!
- The Gigamunch Team`
	introEmailSubject = "Hey thanks for the interest"
	introEmailText    = `Hello! Welcome to GIGAMUNCH! 🎉 I’m Enis, and I’m the CEO and a co-founder of Gigamunch. I want to personally thank you for checking us out and seeing what we’re about.
 
We’re all about great international food with great people. We came together from all different walks of life and made a company to reflect that. Our grassroots movement is catching on so fast that we’re about to hit our capacity. 🙌
 
While there’s still a seat at the dinner table, I’d like to personally invite you to join in the Gigamunch family. We’re so confident you’ll love this experience, I’d like to send you a offer to try it for free. 🤑 You have nothing to lose, just a great experience to gain! Welcome to Gigamunch! 😊

Warm Regards,
Enis
https://eatgigamunch.com/`
	introEmailHTML = `Hello! Welcome to GIGAMUNCH! 🎉 I’m Enis, and I’m the CEO and a co-founder of Gigamunch. I want to personally thank you for checking us out and seeing what we’re about.
 
We’re all about great international food with great people. We came together from all different walks of life and made a company to reflect that. Our grassroots movement is catching on so fast that we’re about to hit our capacity. 🙌
 
While there’s still a seat at the dinner table, I’d like to personally invite you to join in the Gigamunch family. We’re so confident you’ll love this experience, I’d like to send you a offer to try it for free. 🤑 You have nothing to lose, just a great experience to gain! Welcome to Gigamunch! 😊

Warm Regards,
Enis
https://eatgigamunch.com/`
)

var (
	welcomeSender = &User{Name: "Gigamunch", Email: "hello@gigamunchapp.com"}
	introSender   = &User{Name: "Enis", Email: "enis@gigamunchapp.com"}
	sendGridKey   string
	key           string
	acctID        string
)

var (
	errDrip       = errors.InternalServerError
	errBadRequest = errors.BadRequestError
)

// Tag is a tag applied to email subscribers.
type Tag string

func (t Tag) String() string {
	return string(t)
}

const (
	// LeftWebsiteEmail if they left email on website.
	LeftWebsiteEmail Tag = "LEFT_WEBSITE_EMAIL"
	// Customer if they are a customer and is removed when they unsubscribe.
	Customer Tag = "CUSTOMER"
	// Subscribed is applied when a someone subscribers and is never removed.
	Subscribed Tag = "HAS_SUBSCRIBED"
	// Canceled is applied when a subscribers cancels.
	Canceled Tag = "CANCELED"
	// Vegetarian if they are a vegetarian.
	Vegetarian Tag = "VEGETARIAN"
	// NonVegetarian if they a non-vegetarian.
	NonVegetarian Tag = "NON_VEGETARIAN"
	// TwoServings if they are 2 servings.
	TwoServings Tag = "TWO_SERVINGS"
	// FourServings if they are 4 servings.
	FourServings Tag = "FOUR_SERVINGS"
	// Dev if they are development server customers.
	Dev          Tag = "DEV"
	ignoreDomain     = "@test.com"
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

// Client is the client for this package.
type Client struct {
	ctx      context.Context
	sgClient *sendgrid.Client
	dripC    *drip.Client
}

// New returns a new Client.
func New(ctx context.Context) *Client {
	if sendGridKey == "" || key == "" || acctID == "" {
		mailConfig := config.GetMailConfig(ctx)
		sendGridKey = mailConfig.SendGridKey
		key = mailConfig.DripAPIKey
		acctID = mailConfig.DripAccountID
	}
	httpClient := urlfetch.Client(ctx)
	sendgrid.DefaultClient = &rest.Client{HTTPClient: httpClient}
	dripClient, err := drip.New(key, acctID)
	if err != nil {
		utils.Criticalf(ctx, "failed to get drip client: %+v", err)
	} else {
		dripClient.HTTPClient = httpClient
	}
	return &Client{
		ctx:      ctx,
		sgClient: sendgrid.NewSendClient(sendGridKey),
		dripC:    dripClient,
	}
}

// User contains has a name and email.
type User struct {
	Name  string
	Email string
}

// GetName returns a name.
func (u *User) GetName() string {
	return u.Name
}

// GetEmail returns an email.
func (u *User) GetEmail() string {
	return u.Email
}

// UserInterface is used for all mail.
type UserInterface interface {
	GetName() string
	GetEmail() string
}

// SendIntroEmail sends the email for people who just left emails.
func (c *Client) SendIntroEmail(to UserInterface) error {
	return sendEmail(c.sgClient, introSender, to, introEmailSubject, introEmailText, introEmailHTML)
}

// WelcomeEmailInterface is the interface for sending WelcomeEmails.
type WelcomeEmailInterface interface {
	UserInterface
	GetFirstDinnerDate() time.Time
}

// SendWelcomeEmail sends the email for new subscribers.
func (c *Client) SendWelcomeEmail(to WelcomeEmailInterface) error {
	return sendEmail(c.sgClient, welcomeSender, to, welcomeEmailSubject, fmt.Sprintf(welcomeEmailText, to.GetName(), DateString(to.GetFirstDinnerDate())), fmt.Sprintf(welcomeEmailHTML, to.GetName(), to.GetFirstDinnerDate()))
}

func sendEmail(sgClient *sendgrid.Client, from, to UserInterface, subject, emailText, emailHTML string) error {
	fromEmail := mail.NewEmail(from.GetName(), from.GetEmail())
	toEmail := mail.NewEmail(to.GetName(), to.GetEmail())
	email := mail.NewSingleEmail(fromEmail, subject, toEmail, emailText, emailHTML)
	email.Content = email.Content[:1] // only support text content
	resp, err := sgClient.Send(email)
	if err != nil {
		return err // TODO
	}
	if (resp.StatusCode / 100) != 2 {
		return fmt.Errorf("some error. status code: %d", resp.StatusCode) // TODO
	}
	return nil
}

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

// UserFields contain all the possible fields a user can have.
type UserFields struct {
	Email             string `json:"email"`
	Name              string `json:"name"`
	FirstName         string
	LastName          string
	FirstDeliveryDate time.Time `json:"first_delivery_date"`
	AddTags           []Tag     `json:"add_tags"`
	RemoveTags        []Tag     `json:"remove_tags"`
}

// UpdateUser updates the user custom fields.
func (c *Client) UpdateUser(req *UserFields, projID string) error {
	if strings.Contains(req.Email, ignoreDomain) {
		return nil
	}
	// resp, err := c.dripC.FetchSubscriber(req.Email)
	// if err != nil {
	// 	return errDrip.WithError(err).Annotate("failed to drip.FetchSubscriber")
	// }
	// if len(resp.Errors) > 0 {
	// 	return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.FetchSubscriber")
	// }
	// if len(resp.Subscribers) != 1 {
	// 	return errBadRequest.Annotate("failed to find subscriber")
	// }
	sub := drip.UpdateSubscriber{
		Email:        req.Email,
		CustomFields: make(map[string]string),
	}
	var firstName, lastName string
	if req.FirstName == "" {
		firstName, lastName = splitName(req.Name)
	} else {
		firstName = req.FirstName
		lastName = req.LastName
	}
	if firstName != "" {
		sub.CustomFields["FIRST_NAME"] = firstName
	}
	if lastName != "" {
		sub.CustomFields["LAST_NAME"] = lastName
	}
	if !req.FirstDeliveryDate.IsZero() {
		sub.CustomFields["FIRST_DELIVERY_DATE"] = DateString(req.FirstDeliveryDate)
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
	if !common.IsProd(projID) {
		sub.Tags = append(sub.Tags, Dev.String())
	}
	dripReq := &drip.UpdateSubscribersReq{
		Subscribers: []drip.UpdateSubscriber{
			sub,
		},
	}
	resp, err := c.dripC.UpdateSubscriber(dripReq)
	if err != nil {
		return errDrip.WithError(err).Annotate("failed to drip.FetchSubscriber")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.FetchSubscriber")
	}
	return nil
}

func splitName(name string) (string, string) {
	first := ""
	last := ""
	name = strings.Title(strings.TrimSpace(name))
	lastSpace := strings.LastIndex(name, " ")
	if lastSpace == -1 {
		first = name
	} else {
		first = name[:lastSpace]
		last = name[lastSpace:]
	}
	return first, last
}

// AddTag adds a tag to a customer. This often triggers a workflow.
func (c *Client) AddTag(email string, tag Tag) error {
	if strings.Contains(email, ignoreDomain) {
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
	if strings.Contains(email, ignoreDomain) {
		return nil
	}
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

// AddBatchTags adds tags to emails. This often triggers a workflow.
func (c *Client) AddBatchTags(emails []string, tags []Tag) error {
	// TODO:
	tagsString := make([]string, len(tags))
	for i, tag := range tags {
		tagsString[i] = tag.String()
	}
	subs := make([]drip.UpdateSubscriber, len(emails))
	i := 0
	for _, email := range emails {
		if !strings.Contains(email, ignoreDomain) {
			subs[i].Email = email
			subs[i].Tags = tagsString
			i++
		}
	}
	if len(subs) == 0 || len(tags) == 0 {
		return nil
	}
	req := &drip.UpdateBatchSubscribersReq{
		Batches: []drip.SubscribersBatch{
			drip.SubscribersBatch{
				Subscribers: subs[:i],
			},
		},
	}
	resp, err := c.dripC.UpdateBatchSubscribers(req)
	if err != nil {
		return errDrip.WithError(err).Annotate("failed to drip.UpdateBatchSubscribers")
	}
	if len(resp.Errors) > 0 {
		return errDrip.WithError(resp.Errors[0]).Annotate("failed to drip.UpdateBatchSubscribers")
	}
	return nil
}

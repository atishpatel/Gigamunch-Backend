package mail

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/config"
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
https://gigamunchapp.com/`
	introEmailHTML = `Hello! Welcome to GIGAMUNCH! 🎉 I’m Enis, and I’m the CEO and a co-founder of Gigamunch. I want to personally thank you for checking us out and seeing what we’re about.
 
We’re all about great international food with great people. We came together from all different walks of life and made a company to reflect that. Our grassroots movement is catching on so fast that we’re about to hit our capacity. 🙌
 
While there’s still a seat at the dinner table, I’d like to personally invite you to join in the Gigamunch family. We’re so confident you’ll love this experience, I’d like to send you a offer to try it for free. 🤑 You have nothing to lose, just a great experience to gain! Welcome to Gigamunch! 😊

Warm Regards,
Enis
https://gigamunchapp.com/`
)

var (
	welcomeSender = &User{Name: "Gigamunch", Email: "hello@gigamunchapp.com"}
	introSender   = &User{Name: "Enis", Email: "enis@gigamunchapp.com"}
	sendGridKey   string
)

// Client is the client for this package.
type Client struct {
	ctx      context.Context
	sgClient *sendgrid.Client
}

// New returns a new Client.
func New(ctx context.Context) *Client {
	if sendGridKey == "" {
		mailConfig := config.GetMailConfig(ctx)
		sendGridKey = mailConfig.SendGridKey
	}
	sendgrid.DefaultClient = &rest.Client{HTTPClient: urlfetch.Client(ctx)}
	return &Client{
		ctx:      ctx,
		sgClient: sendgrid.NewSendClient(sendGridKey),
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
	return sendEmail(c.sgClient, welcomeSender, to, welcomeEmailSubject, fmt.Sprintf(welcomeEmailText, to.GetName(), dateString(to.GetFirstDinnerDate())), fmt.Sprintf(welcomeEmailHTML, to.GetName(), to.GetFirstDinnerDate()))
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

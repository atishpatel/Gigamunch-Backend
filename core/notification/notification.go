package notification

import (
	"fmt"
	"sync"

	twilio "github.com/atishpatel/twiliogo"
	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"

	"google.golang.org/appengine/mail"
	"google.golang.org/appengine/urlfetch"
)

var (
	onceConfig   = sync.Once{}
	twilioConfig config.TwilioConfig
	from         []string
	errTwilio    = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with twilio sms."}
	errEmail     = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with sending email."}
	errFakeInput = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Input is invalid."}
)

type Client struct {
	ctx     context.Context
	twilioC *twilio.TwilioClient
}

func New(ctx context.Context) *Client {
	onceConfig.Do(func() {
		twilioConfig = config.GetTwilioConfig(ctx)
		from = twilioConfig.PhoneNumbers
	})
	return &Client{
		ctx:     ctx,
		twilioC: getTwilioClient(ctx),
	}
}

func (c *Client) SendEmail(to, subject, message string) error {
	msg := &mail.Message{
		Sender:  "support@gigamunchapp.com",
		To:      []string{to},
		Subject: subject,
		Body:    message,
	}
	if err := mail.Send(c.ctx, msg); err != nil {
		return errEmail.WithError(err).Wrap("error sending email")
	}
	return nil
}

// SendSMS sends an sms to the user
func (c *Client) SendSMS(to, message string) error {
	_, err := twilio.NewMessage(c.twilioC, getFromNumber(to), to, twilio.Body(message))
	if err != nil {
		if twilioErr, ok := err.(*twilio.TwilioError); ok {
			switch twilioErr.Code {
			case 21211:
				fallthrough
			case 21612:
				fallthrough
			case 21614:
				return errFakeInput.WithMessage(fmt.Sprintf("Failed to send sms because %s", twilioErr.Message))
			}
		}
		return errTwilio.WithError(err).Wrap("error sending message via twilio")
	}
	return nil
}

func getTwilioClient(ctx context.Context) *twilio.TwilioClient {
	client := twilio.NewClient(twilioConfig.AccountSID, twilioConfig.AuthToken)
	client.HTTPClient = urlfetch.Client(ctx)
	return client
}

func getFromNumber(to string) string {
	return from[0]
}

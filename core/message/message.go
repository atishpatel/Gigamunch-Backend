package message

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	twilio "github.com/atishpatel/twiliogo"
)

var (
	twilioConfig config.TwilioConfig
	errTwilio    = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with Twilio. Please try again in a few minutes."}
	errFakeInput = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Input is invalid."}
)

// Client is the client for ip messaging, sms, and email.
type Client struct {
	ctx     context.Context
	twilioC *twilio.TwilioClient
}

// New creates a message client.
func New(ctx context.Context) *Client {
	if twilioConfig.AccountSID == "" {
		twilioConfig = config.GetTwilioConfig(ctx)
	}
	c := &Client{
		ctx: ctx,
	}
	c.twilioC = twilio.NewClient(twilioConfig.AccountSID, twilioConfig.AuthToken)
	httpClient := urlfetch.Client(ctx)
	c.twilioC.HTTPClient = httpClient
	return c
}

// SendAdminSMS sends an sms,
func (c *Client) SendAdminSMS(to, message string) error {
	_, err := twilio.NewMessageFromService(c.twilioC, twilioConfig.AdminServiceSID, to, twilio.Body(message))
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

// SendDeliverySMS sends an sms,
func (c *Client) SendDeliverySMS(to, message string) error {
	_, err := twilio.NewMessageFromService(c.twilioC, twilioConfig.DeliveryServiceSID, to, twilio.Body(message))
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

// SendBagSMS sends an sms,
func (c *Client) SendBagSMS(to, message string) error {
	_, err := twilio.NewMessageFromService(c.twilioC, twilioConfig.BagServiceSID, to, twilio.Body(message))
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

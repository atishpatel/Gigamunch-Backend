package message

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	twilio "github.com/atishpatel/twiliogo"
	"golang.org/x/net/context"

	"google.golang.org/appengine/mail"
	"google.golang.org/appengine/urlfetch"
)

const (
	channelAttr = `{"cook_id":"%s","cook_name":"%s","cook_image":"%s","eater_id":"%s","eater_name":"%s","eater_image":"%s","inquiry_id":"%s","inquiry_status":"%s","cook_action":"%s","eater_action":"%s","item_id":"%s",item_name":"%s","item_image":"%s"}`
	userAttr    = `{"id":"%s","name":"%s","image":"%s"}`
)

var (
	onceConfig          = sync.Once{}
	twilioConfig        config.TwilioConfig
	from                []string
	serviceSID          string
	inquiryBotSID       string
	inquiryStatusBotSID string
	gigamunchBotSID     string
	errTwilio           = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with twilio sms."}
	errEmail            = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with sending email."}
	errFakeInput        = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Input is invalid."}
)

type Client struct {
	ctx       context.Context
	twilioC   *twilio.TwilioClient
	twilioIPC *twilio.TwilioIPMessagingClient
}

func New(ctx context.Context) *Client {
	onceConfig.Do(func() {
		twilioConfig = config.GetTwilioConfig(ctx)
		serviceSID = twilioConfig.IPMessagingSID
		from = twilioConfig.PhoneNumbers
	})
	c := &Client{
		ctx: ctx,
	}
	c.twilioC, c.twilioIPC = getTwilioClients(ctx, twilioConfig.AccountSID, twilioConfig.AuthToken)
	return c
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

// UserInfo contains the info for an message user.
type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// InquiryInfo contains info attached to an Inquiry for message.
type InquiryInfo struct {
	ID          int64  `json:"id,string"`
	Status      string `json:"status"`
	CookAction  string `json:"cook_action"`
	EaterAction string `json:"eater_action"`
	ItemID      int64  `json:"item_id"`
	ItemName    string `json:"item_name"`
	ItemImage   string `json:"item_image"`
}

// returns friendlyName, uniqueName, attributes
func getChannelNamesAndAttr(c *UserInfo, e *UserInfo, i *InquiryInfo) (string, string, string) {
	inqID := strconv.FormatInt(i.ID, 10)
	itemID := strconv.FormatInt(i.ItemID, 10)
	friendlyName := fmt.Sprintf("%s<;>%s", c.Name, e.Name)
	uniqueName := fmt.Sprintf("%s<;>%s", c.ID, e.ID)
	attr := fmt.Sprintf(channelAttr, c.ID, c.Name, c.Image, e.ID, e.Name, e.Image, inqID, i.Status, i.CookAction, i.EaterAction, itemID, i.ItemName, i.ItemImage)
	return friendlyName, uniqueName, attr
}

// returns friendlyName, uniqueName, attributes
func getUserAttr(u *UserInfo) string {
	attr := fmt.Sprintf(userAttr, u.ID, u.Name, u.Image)
	return attr
}

// UpdateChannel creates or updates a channel
func (c *Client) UpdateChannel(cookInfo *UserInfo, eaterInfo *UserInfo, inquiryInfo *InquiryInfo) error {
	friendlyName, uniqueName, attributes := getChannelNamesAndAttr(cookInfo, eaterInfo, inquiryInfo)
	channel, err := twilio.GetIPChannel(c.twilioIPC, serviceSID, uniqueName)
	isNew := false
	if err != nil {
		tErr, ok := err.(*twilio.TwilioError)
		if !ok || tErr.Code != twilio.CodeNotFound {
			return errTwilio.WithError(err).Wrap("failed to GetIPChannel")
		}
		isNew = true
	}
	if isNew {
		_, err = twilio.NewIPChannel(c.twilioIPC, serviceSID, friendlyName, uniqueName, false, attributes)
		if err != nil {
			return errTwilio.WithError(err)
		}
		// add users to channel
		err = addUserToChannel(c.twilioIPC, channel.Sid, cookInfo)
		if err != nil {
			return err
		}
		err = addUserToChannel(c.twilioIPC, channel.Sid, eaterInfo)
		if err != nil {
			return err
		}
		err = addBotsToChannel(c.twilioIPC, channel.Sid)
		if err != nil {
			return errors.Wrap("failed to addBotsToChannel", err)
		}
	} else {
		_, err = twilio.UpdateIPChannel(c.twilioIPC, serviceSID, friendlyName, uniqueName, uniqueName, false, attributes)
	}
	if err != nil {
		return errTwilio.WithError(err)
	}
	return nil
}

func addBotsToChannel(twilioIPC *twilio.TwilioIPMessagingClient, channelSID string) error {
	var err error
	if inquiryBotSID == "" {
		// get bot SIDs
		var u *twilio.IPUser
		u, err = twilio.GetIPUser(twilioIPC, serviceSID, "GigamunchBot")
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to get GigamunchBot")
		}
		gigamunchBotSID = u.Sid
		u, err = twilio.GetIPUser(twilioIPC, serviceSID, "InquiryBot")
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to get InquiryBot")
		}
		inquiryBotSID = u.Sid
		u, err = twilio.GetIPUser(twilioIPC, serviceSID, "InquiryStatusBot")
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to get InquiryStatusBot")
		}
		inquiryStatusBotSID = u.Sid
	}
	_, err = twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, gigamunchBotSID, "")
	if err != nil {
		return errTwilio.WithError(err).Wrapf("failed to add GigamunchBot to channel(%s)", channelSID)
	}
	_, err = twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, inquiryBotSID, "")
	if err != nil {
		return errTwilio.WithError(err).Wrapf("failed to add InquiryBot to channel(%s)", channelSID)
	}
	_, err = twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, inquiryStatusBotSID, "")
	if err != nil {
		return errTwilio.WithError(err).Wrapf("failed to add InquiryStatusBot to channel(%s)", channelSID)
	}
	return nil
}

func addUserToChannel(twilioIPC *twilio.TwilioIPMessagingClient, channelSID string, userInfo *UserInfo) error {
	doesntExist := false
	_, err := twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, userInfo.ID, "")
	if err != nil {
		tErr, ok := err.(*twilio.TwilioError)
		if !ok || tErr.Code != twilio.CodeNotFound {
			return errTwilio.WithError(err).Wrap("failed to AddIPMemberToChannel")
		}
		doesntExist = true
	}
	if doesntExist {
		_, err = createUser(twilioIPC, userInfo)
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to createUser")
		}
		return addUserToChannel(twilioIPC, channelSID, userInfo)
	}
	return nil
}

func createUser(twilioIPC *twilio.TwilioIPMessagingClient, userInfo *UserInfo) (*twilio.IPUser, error) {
	attributes := getUserAttr(userInfo)
	return twilio.NewIPUser(twilioIPC, serviceSID, userInfo.ID, "", attributes)
}

// func (c *Client) AddInquiry() error {}

// func (c *Client) UpdateInquiryStatus() error {}

// func (c *Client)

// UpdateUser creates or updates a user
func (c *Client) UpdateUser(userInfo *UserInfo) error {
	// TODO change if twilio updates their API
	doesntExist := false
	user, err := twilio.GetIPUser(c.twilioIPC, serviceSID, userInfo.ID)
	if err != nil {
		tErr, ok := err.(*twilio.TwilioError)
		if !ok || tErr.Code != twilio.CodeNotFound {
			return errTwilio.WithError(err).Wrap("failed to AddIPMemberToChannel")
		}
		doesntExist = true
	}
	if doesntExist {
		_, err = createUser(c.twilioIPC, userInfo)
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to createUser")
		}
	} else {
		attributes := getUserAttr(userInfo)
		_, err := twilio.UpdateIPUser(c.twilioIPC, serviceSID, user.Sid, userInfo.ID, "", attributes)
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to createUser")
		}
	}
	return nil
}

func getTwilioClients(ctx context.Context, accountSID, authToken string) (*twilio.TwilioClient, *twilio.TwilioIPMessagingClient) {
	client := twilio.NewClient(twilioConfig.AccountSID, twilioConfig.AuthToken)
	httpClient := urlfetch.Client(ctx)
	client.HTTPClient = httpClient
	ipClient := twilio.NewIPMessagingClient(accountSID, authToken)
	ipClient.HTTPClient = httpClient
	return client, ipClient
}

func getFromNumber(to string) string {
	return from[0]
}

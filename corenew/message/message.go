package message

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	twilio "github.com/atishpatel/twiliogo"
	jwt "gopkg.in/dgrijalva/jwt-go.v2"
)

const (
	channelAttr       = `{"cook_sid":"%s","cook_id":"%s","cook_name":"%s","cook_image":"%s","eater_sid":"%s","eater_id":"%s","eater_name":"%s","eater_image":"%s","inquiry_id":"%d","inquiry_state":"%s","cook_action":"%s","eater_action":"%s","item_id":"%d","item_name":"%s","item_image":"%s"}`
	userAttr          = `{"id":"%s","name":"%s","image":"%s"}`
	inquiryAttr       = `{"inquiry_id":"%d","inquiry_state":"%s","cook_action":"%s","eater_action":"%s","item_id":"%d","item_name":"%s","item_image":"%s","price":"%f","is_delivery":"%t","servings":"%d","exchange_time":"%d"}`
	inquiryStatusAttr = `{"inquiry_id":"%d","inquiry_state":"%s","cook_action":"%s","eater_action":"%s","item_id":"%d","item_name":"%s","item_image":"%s","title":"%s","message":"%s"}`
	// InquiryBotID is the Unique Name for the InquiryBot.
	InquiryBotID = "InquiryBot"
	// InquiryStatusBotID is the Unique Name for the InquiryStatusBot.
	InquiryStatusBotID = "InquiryStatusBot"
	// GigamunchBotID is the Unique Name for the GigamunchBot.
	GigamunchBotID = "GigamunchBot"
)

var (
	twilioConfig       config.TwilioConfig
	from               []string
	serviceSID         string
	errInvalidParamter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errInternal        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Something went wrong with the server."}
	errTwilio          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with twilio."}
	// errEmail           = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with sending email."}
	errFakeInput = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Input is invalid."}
)

// Client is the client for ip messaging, sms, and email.
type Client struct {
	ctx       context.Context
	twilioC   *twilio.TwilioClient
	twilioIPC *twilio.TwilioIPMessagingClient
}

// New creates a message client.
func New(ctx context.Context) *Client {
	if serviceSID == "" {
		twilioConfig = config.GetTwilioConfig(ctx)
		serviceSID = twilioConfig.IPMessagingSID
		from = twilioConfig.PhoneNumbers
	}
	c := &Client{
		ctx: ctx,
	}
	c.twilioC, c.twilioIPC = getTwilioClients(ctx, twilioConfig.AccountSID, twilioConfig.AuthToken, twilioConfig.KeySID, twilioConfig.KeyAuthToken)
	return c
}

// func (c *Client) SendEmail(to, subject, message string) error {
// 	msg := &mail.Message{
// 		Sender:  "support@gigamunchapp.com",
// 		To:      []string{to},
// 		Subject: subject,
// 		Body:    message,
// 	}
// 	if err := mail.Send(c.ctx, msg); err != nil {
// 		return errEmail.WithError(err).Wrap("error sending email")
// 	}
// 	return nil
// }

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
	ID           int64     `json:"id,string"`
	State        string    `json:"state"`
	CookAction   string    `json:"cook_action"`
	EaterAction  string    `json:"eater_action"`
	ItemID       int64     `json:"item_id"`
	ItemName     string    `json:"item_name"`
	ItemImage    string    `json:"item_image"`
	Price        float32   `json:"price"`
	IsDelivery   bool      `json:"is_delivery"`
	Servings     int32     `json:"servings"`
	ExchangeTime time.Time `json:"exchange_time"`
}

func getChannelUniqueName(cookID, eaterID string) string {
	return fmt.Sprintf("%s<;>%s", cookID, eaterID)
}

// returns friendlyName, uniqueName, attributes
func getChannelNamesAndAttr(cSID string, c *UserInfo, eSID string, e *UserInfo, i *InquiryInfo) (string, string, string) {
	friendlyName := fmt.Sprintf("%s<;>%s", c.Name, e.Name)
	uniqueName := getChannelUniqueName(c.ID, e.ID)
	attr := fmt.Sprintf(channelAttr, cSID, c.ID, c.Name, c.Image, eSID, e.ID, e.Name, e.Image, i.ID, i.State, i.CookAction, i.EaterAction, i.ItemID, i.ItemName, i.ItemImage)
	return friendlyName, uniqueName, attr
}

// returns friendlyName, uniqueName, attributes
func getUserAttr(u *UserInfo) string {
	attr := fmt.Sprintf(userAttr, u.ID, u.Name, u.Image)
	return attr
}

// UpdateChannel creates or updates a channel
func (c *Client) UpdateChannel(cookInfo *UserInfo, eaterInfo *UserInfo, inquiryInfo *InquiryInfo) error {
	if cookInfo == nil || eaterInfo == nil {
		return errInvalidParamter.Wrap("cookInfo and eaterInfo cannot be nil")
	}
	if inquiryInfo == nil {
		inquiryInfo = new(InquiryInfo)
	}
	eaterTUser, err := createUserIfNotExist(c.twilioIPC, eaterInfo)
	if err != nil {
		return errTwilio.WithError(err).Wrap("failed to createUserIfNotExist")
	}
	cookTUser, err := createUserIfNotExist(c.twilioIPC, cookInfo)
	if err != nil {
		return errTwilio.WithError(err).Wrap("failed to createUserIfNotExist")
	}
	friendlyName, uniqueName, attributes := getChannelNamesAndAttr(cookTUser.Sid, cookInfo, eaterTUser.Sid, eaterInfo, inquiryInfo)
	channel, err := twilio.GetIPChannel(c.twilioIPC, serviceSID, uniqueName)
	isNew := false
	if err != nil {
		tErr, ok := err.(*twilio.TwilioError)
		if !ok || tErr.Code != twilio.CodeNotFound {
			return errTwilio.WithError(err).Wrap("failed to twilio.GetIPChannel")
		}
		isNew = true
	}
	if isNew {
		channel, err = twilio.NewIPChannel(c.twilioIPC, serviceSID, friendlyName, uniqueName, false, attributes)
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to twilio.NewIPChannel")
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
		_, err = twilio.UpdateIPChannel(c.twilioIPC, serviceSID, uniqueName, friendlyName, uniqueName, false, attributes)
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to twilio.UpdateIPChannel")
		}
	}
	return nil
}

func addBotsToChannel(twilioIPC *twilio.TwilioIPMessagingClient, channelSID string) error {
	if channelSID == "" {
		return errInvalidParamter.Wrap("channelSID cannot be an empty string")
	}
	var err error
	_, err = twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, GigamunchBotID, "")
	if err != nil {
		return errTwilio.WithError(err).Wrapf("failed to add GigamunchBot to channel(%s)", channelSID)
	}
	_, err = twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, InquiryBotID, "")
	if err != nil {
		return errTwilio.WithError(err).Wrapf("failed to add InquiryBot to channel(%s)", channelSID)
	}
	_, err = twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, InquiryStatusBotID, "")
	if err != nil {
		return errTwilio.WithError(err).Wrapf("failed to add InquiryStatusBot to channel(%s)", channelSID)
	}
	return nil
}

func addUserToChannel(twilioIPC *twilio.TwilioIPMessagingClient, channelSID string, userInfo *UserInfo) error {
	if channelSID == "" {
		return errInvalidParamter.Wrap("channelSID cannot be an empty string")
	}
	if userInfo == nil || userInfo.ID == "" {
		return errInvalidParamter.Wrap("User ID cannot be an empty string")
	}
	doesntExist := false
	_, err := twilio.AddIPMemberToChannel(twilioIPC, serviceSID, channelSID, userInfo.ID, "")
	if err != nil {
		tErr, ok := err.(*twilio.TwilioError)
		if !ok || tErr.Code != twilio.CodeNotFound {
			return errTwilio.WithError(err).Wrap("failed to twilio.AddIPMemberToChannel")
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
	user, err := twilio.NewIPUser(twilioIPC, serviceSID, userInfo.ID, userInfo.Name+":", "", attributes)
	if err != nil {
		return nil, errTwilio.WithError(err).Wrap("failed to twilio.NewIPUser")
	}
	return user, nil
}

func createUserIfNotExist(twilioIPC *twilio.TwilioIPMessagingClient, userInfo *UserInfo) (*twilio.IPUser, error) {
	user, err := twilio.GetIPUser(twilioIPC, serviceSID, userInfo.ID)
	if err != nil {
		tErr, ok := err.(*twilio.TwilioError)
		if !ok || tErr.Code != twilio.CodeNotFound {
			return nil, errTwilio.WithError(err).Wrap("failed to twilio.AddIPMemberToChannel")
		}
		user, err = createUser(twilioIPC, userInfo)
		if err != nil {
			return nil, errors.Wrap("failed to createUser", err)
		}
	}
	return user, nil
}

func getInquiryBodyAndAttributes(inqI *InquiryInfo) (string, string) {
	body := fmt.Sprintf("There is an update about your request for %s", inqI.ItemName)
	attr := fmt.Sprintf(inquiryAttr, inqI.ID, inqI.State, inqI.CookAction, inqI.EaterAction, inqI.ItemID, inqI.ItemName, inqI.ItemImage, inqI.Price, inqI.IsDelivery, inqI.Servings, inqI.ExchangeTime.Unix())
	return body, attr
}

// SendInquiryBotMessage sends a inquiry message from the InquiryBot to the apporiate channel. This function calls UpdateChannel. Returns MessageID, error.
func (c *Client) SendInquiryBotMessage(cookInfo *UserInfo, eaterInfo *UserInfo, inquiryInfo *InquiryInfo) (string, error) {
	err := c.UpdateChannel(cookInfo, eaterInfo, inquiryInfo)
	if err != nil {
		return "", errors.Wrap("failed to message.UpdateChannel", err)
	}
	if cookInfo == nil || eaterInfo == nil || inquiryInfo == nil {
		return "", errInvalidParamter.Wrap("inquiryInfo, cookInfo, and eaterInfo cannot be nil")
	}
	channelUniqueName := getChannelUniqueName(cookInfo.ID, eaterInfo.ID)
	// get channel
	channel, err := twilio.GetIPChannel(c.twilioIPC, serviceSID, channelUniqueName)
	if err != nil {
		return "", errTwilio.WithError(err).Wrap("failed to twilio.GetIPChannel")
	}
	body, attr := getInquiryBodyAndAttributes(inquiryInfo)
	msg, err := twilio.SendIPMessageToChannel(c.twilioIPC, serviceSID, channel.Sid, InquiryBotID, body, attr)
	if err != nil {
		return "", errTwilio.WithError(err).Wrap("failed to twilio.GetIPChannel")
	}
	return msg.Sid, nil
}

func getInquiryStatusBodyAndAttributes(c *Client, cookInfo *UserInfo, eaterInfo *UserInfo, inqI *InquiryInfo) (string, string) {
	var message string
	var title string
	const canceled = "Canceled"
	switch inqI.State {
	case "Accepted":
		if inqI.CookAction == "Accepted" {
			message = fmt.Sprintf("%s just accepted the request for %s", cookInfo.Name, inqI.ItemName)
		} else {
			message = fmt.Sprintf("%s just accepted the request for %s", eaterInfo.Name, inqI.ItemName)
		}
		title = "Request Accepted"
	case "Declined":
		if inqI.CookAction == "Declined" {
			message = fmt.Sprintf("%s just declined the request for %s", cookInfo.Name, inqI.ItemName)
		}
		title = "Request Declined"
	case "TimedOut":
		message = fmt.Sprintf("%s couldn't fulfill the request for %s", cookInfo.Name, inqI.ItemName)
		title = "Request Timed Out"
	case canceled:
		if inqI.CookAction == canceled {
			message = fmt.Sprintf("%s just canceled the request for %s", cookInfo.Name, inqI.ItemName)
		} else if inqI.EaterAction == canceled {
			message = fmt.Sprintf("%s just canceled the request for %s", eaterInfo.Name, inqI.ItemName)
		}
		title = "Request Canceled"
	}
	if message == "" {
		_ = c.SendSMS("9316445311", fmt.Sprintf("Unknown Action cook(%s)/eater(%s) for State(%s) for inquiry state bot message updated!!!", inqI.CookAction, inqI.EaterAction, inqI.State))
		message = fmt.Sprintf("There is an update about your request for %s", inqI.ItemName)
	}
	if title == "" {
		_ = c.SendSMS("9316445311", fmt.Sprintf("Unknown title for Action cook(%s)/eater(%s) for State(%s) for inquiry state bot message updated!!!", inqI.CookAction, inqI.EaterAction, inqI.State))
		title = "-"
	}
	body := message
	attr := fmt.Sprintf(inquiryStatusAttr, inqI.ID, inqI.State, inqI.CookAction, inqI.EaterAction, inqI.ItemID, inqI.ItemName, inqI.ItemImage, title, message)
	return body, attr
}

// UpdateInquiryStatus sends a inquiry message from the InquiryStatusBot to the apporiate channel, and updates the InquiryBot Message. This function calls UpdateChannel.
func (c *Client) UpdateInquiryStatus(messageSID string, cookInfo *UserInfo, eaterInfo *UserInfo, inquiryInfo *InquiryInfo) error {
	err := c.UpdateChannel(cookInfo, eaterInfo, inquiryInfo)
	if err != nil {
		return errors.Wrap("failed to message.UpdateChannel", err)
	}
	if messageSID == "" || cookInfo == nil || eaterInfo == nil || inquiryInfo == nil {
		return errInvalidParamter.Wrap("messageSID, inquiryInfo, cookInfo, and eaterInfo cannot be nil")
	}
	channelUniqueName := getChannelUniqueName(cookInfo.ID, eaterInfo.ID)
	// get channel
	channel, err := twilio.GetIPChannel(c.twilioIPC, serviceSID, channelUniqueName)
	if err != nil {
		return errTwilio.WithError(err).Wrap("failed to twilio.GetIPChannel")
	}
	// send message
	body, attr := getInquiryStatusBodyAndAttributes(c, cookInfo, eaterInfo, inquiryInfo)
	_, err = twilio.SendIPMessageToChannel(c.twilioIPC, serviceSID, channel.Sid, InquiryStatusBotID, body, attr)
	if err != nil {
		return errTwilio.WithError(err).Wrap("failed to twilio.GetIPChannel")
	}
	// update InquiryBot message
	body, attr = getInquiryBodyAndAttributes(inquiryInfo)
	_, err = twilio.UpdateIPMessage(c.twilioIPC, serviceSID, channel.Sid, messageSID, body, attr)
	if err != nil {
		return errTwilio.WithError(err).Wrap("failed to twilio.UpdateIPMessage")
	}
	return nil
}

// func (c *Client) SendGigamunchBotMessage() error {}

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
		_, err := twilio.UpdateIPUser(c.twilioIPC, serviceSID, user.Sid, userInfo.ID, userInfo.Name+":", "", attributes)
		if err != nil {
			return errTwilio.WithError(err).Wrap("failed to createUser")
		}
	}
	return nil
}

// GetToken gets a messaging token for the user.
func (c *Client) GetToken(userInfo *UserInfo, deviceID string) (string, error) {
	_, err := createUserIfNotExist(c.twilioIPC, userInfo)
	if err != nil {
		return "", err
	}
	endpointID := fmt.Sprintf("Gigamunch:%s:%s", userInfo.ID, deviceID)
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	nowUnix := time.Now().Unix()
	jwtToken.Header["cty"] = "twilio-fpa;v=1"
	jwtToken.Claims["jti"] = fmt.Sprintf("%s-%d", twilioConfig.KeySID, nowUnix)
	jwtToken.Claims["iss"] = twilioConfig.KeySID
	jwtToken.Claims["sub"] = twilioConfig.AccountSID
	jwtToken.Claims["exp"] = nowUnix + 7200 // 2 hours
	ipMessaging := make(map[string]string, 3)
	ipMessaging["service_sid"] = serviceSID
	ipMessaging["endpoint_id"] = endpointID
	ipMessaging["push_credential_sid"] = twilioConfig.PushCredentialSID
	grants := make(map[string]interface{}, 2)
	grants["identity"] = userInfo.ID
	grants["ip_messaging"] = ipMessaging
	jwtToken.Claims["grants"] = grants
	tkn, err := jwtToken.SignedString([]byte(twilioConfig.KeyAuthToken))
	if err != nil {
		return "", errInternal.WithError(err).Wrap("failed to jwt.SignedString")
	}
	return tkn, nil
}

// GetChannelInfoResp is the response for GetChannelInfo.
type GetChannelInfoResp struct {
	EaterID string `json:"eater_id"`
	CookID  string `json:"cook_id"`
}

// GetChannelInfo returns the Cook and Eater ids.
func (c *Client) GetChannelInfo(channelSID string) (*GetChannelInfoResp, error) {
	channel, err := twilio.GetIPChannel(c.twilioIPC, serviceSID, channelSID)
	if err != nil {
		return nil, errTwilio.WithError(err).Wrap("failed to twilio.GetIPChannel")
	}
	resp := new(GetChannelInfoResp)
	err = json.Unmarshal([]byte(channel.Attributes), resp)
	if err != nil {
		return nil, errInternal.WithError(err).Wrapf("failed to json.Unmarshal channel(%s) attributes: %s", channelSID, channel.Attributes)
	}
	return resp, nil
}

// GetUserInfo returns the user info for a UserSID.
func (c *Client) GetUserInfo(userSID string) (*UserInfo, error) {
	user, err := twilio.GetIPUser(c.twilioIPC, serviceSID, userSID)
	if err != nil {
		return nil, errTwilio.WithError(err).Wrap("failed to twilio.GetIPUser")
	}
	userInfo := new(UserInfo)
	err = json.Unmarshal([]byte(user.Attributes), userInfo)
	if err != nil {
		return nil, errInternal.WithError(err).Wrapf("failed to json.Unmarshal userSID(%s) attributes: %s", userSID, user.Attributes)
	}
	return userInfo, nil
}

func getTwilioClients(ctx context.Context, accountSID, authToken, keySID, apiSecret string) (*twilio.TwilioClient, *twilio.TwilioIPMessagingClient) {
	client := twilio.NewClient(accountSID, authToken)
	httpClient := urlfetch.Client(ctx)
	client.HTTPClient = httpClient
	ipClient := twilio.NewIPMessagingClient(keySID, apiSecret)
	ipClient.HTTPClient = httpClient
	return client, ipClient
}

func getFromNumber(to string) string {
	if len(from) == 0 {
		return "14243484448"
	}
	return from[0]
}

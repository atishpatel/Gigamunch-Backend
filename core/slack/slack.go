package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/nlopes/slack"
	"google.golang.org/appengine/urlfetch"
)

const (
	gigabotToken            = "xoxb-4964967729-442626391122-LZn8tjVu2d1QnAufFTVxqpnr"
	customerSupportChannel  = "#customer-support"
	subscriberUpdateChannel = "#subscriber-pulse"
)

var (
	errSlack    = errors.InternalServerError.WithMessage("Error with Slack.")
	errInternal = errors.InternalServerError
)

// Client is a client.
type Client struct {
	ctx         context.Context
	log         *logging.Client
	serverInfo  *common.ServerInfo
	slackClient *slack.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	if serverInfo == nil {
		return nil, errInternal.Annotate("failed to get server info")
	}
	var ops slack.Option
	if serverInfo.IsStandardAppEngine {
		httpClient := urlfetch.Client(ctx)
		ops = slack.OptionHTTPClient(httpClient)
	}
	slackClient := slack.New(gigabotToken, ops)
	return &Client{
		ctx:         ctx,
		log:         log,
		serverInfo:  serverInfo,
		slackClient: slackClient,
	}, nil
}

func (c *Client) sendMessage(channel, message string) error {
	ops := []slack.MsgOption{
		slack.MsgOptionText(message, true),
	}
	_, _, _, err := c.slackClient.SendMessageContext(c.ctx, channel, ops...)
	if err != nil {
		return errSlack.WithError(err).Annotate("failed to slack.SendMessageContext")
	}
	return nil
}

func (c *Client) sendAttachmentsMessage(channel, message string, attachments []slack.Attachment) error {
	ops := []slack.MsgOption{
		slack.MsgOptionText(message, true),
		slack.MsgOptionAttachments(attachments...),
	}
	_, _, _, err := c.slackClient.SendMessageContext(c.ctx, channel, ops...)
	if err != nil {
		return errSlack.WithError(err).Annotate("failed to slack.SendMessageContext")
	}
	return nil
}

func (c *Client) SendCustomerSupportMessage(message string) error {
	return c.sendMessage(customerSupportChannel, message)
}

func (c *Client) SendMissedSubscriber(email, name, reference string, campaigns []common.Campaign, address common.Address) error {
	message := "Missed out on subscriber. Out of zone."
	attachment := slack.Attachment{
		Fallback: fmt.Sprintf("%s out of zone", email),
		Title:    fmt.Sprintf("%s - <https://eatgigamunch.com/subscriber/%s|%s>", name, email, email),
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Email",
				Value: fmt.Sprintf("<https://eatgigamunch.com/subscriber/%s|%s>", email, email),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Name",
				Value: name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Reference",
				Value: reference,
				Short: false,
			},
			slack.AttachmentField{
				Title: "Address",
				Value: address.StringNoAPT(),
				Short: false,
			},
		},
		Color: "warning",
	}
	for _, campaign := range campaigns {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Campaign",
			Value: fmt.Sprintf("Source:%s\nMedium:%s\nCampaign:%s\n", campaign.Source, campaign.Medium, campaign.Campaign),
		})
	}
	attachments := []slack.Attachment{
		attachment,
	}

	return c.sendAttachmentsMessage(subscriberUpdateChannel, message, attachments)
}

func (c *Client) SendNewSignup(email, name, reference string, campaigns []common.Campaign) error {
	message := "New Subscriber!!!"
	attachment := slack.Attachment{
		Fallback: fmt.Sprintf("%s just signed up", email),
		Title:    fmt.Sprintf("%s - <https://eatgigamunch.com/subscriber/%s|%s>", name, email, email),
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Email",
				Value: fmt.Sprintf("<https://eatgigamunch.com/subscriber/%s|%s>", email, email),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Name",
				Value: name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Reference",
				Value: reference,
				Short: false,
			},
		},
		Color: "green",
	}
	for _, campaign := range campaigns {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Campaign",
			Value: fmt.Sprintf("Source:%s\nMedium:%s\nCampaign:%s\n", campaign.Source, campaign.Medium, campaign.Campaign),
		})
	}
	attachments := []slack.Attachment{
		attachment,
	}

	return c.sendAttachmentsMessage(subscriberUpdateChannel, message, attachments)
}

func (c *Client) SendDeactivate(email, name, reason string, daysActive int) error {
	message := "Subscriber Deactivated"
	attachment := slack.Attachment{
		Fallback: fmt.Sprintf("%s just deactivated", email),
		Title:    fmt.Sprintf("%s - <https://eatgigamunch.com/subscriber/%s|%s>", name, email, email),
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Email",
				Value: fmt.Sprintf("<https://eatgigamunch.com/subscriber/%s|%s>", email, email),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Name",
				Value: name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Reason",
				Value: reason,
				Short: false,
			},
			slack.AttachmentField{
				Title: "Active Period",
				Value: fmt.Sprintf("%.2f weeks\n%d days", float32(daysActive)/7.0, daysActive),
				Short: true,
			},
		},
		Color: "danger",
	}
	attachments := []slack.Attachment{
		attachment,
	}

	return c.sendAttachmentsMessage(subscriberUpdateChannel, message, attachments)
}

type WebhookEvent struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type WebhookRequest struct {
	Challenge string       `json:"challenge"`
	Event     WebhookEvent `json:"event"`
	Type      string       `json:"type"`
}

func (c *Client) HandleWebhook(req *WebhookRequest) error {
	var err error
	if req.Type == "event_callback" && strings.Contains(req.Event.Text, "@") && strings.Contains(req.Event.Text, "<mailto:") {
		// add subscriber page link if there is an email
		email := req.Event.Text[strings.Index(req.Event.Text, "<mailto:")+8 : strings.Index(req.Event.Text, "|")]
		txt := fmt.Sprintf("https://eatgigamunch.com/admin/subscriber/%s", email)
		err = c.SendCustomerSupportMessage(txt)
		if err != nil {
			return errSlack.WithError(err).Annotate("failed to SendCustoemrSupportMessage")
		}
	}
	return nil
}

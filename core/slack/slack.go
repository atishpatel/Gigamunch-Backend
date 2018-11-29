package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/nlopes/slack"
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
	slackClient := slack.New(gigabotToken)
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

// SendCustomerSupportMessage sends message to slack #customer-support.
func (c *Client) SendCustomerSupportMessage(message string) error {
	return c.sendMessage(customerSupportChannel, message)
}

// SendMissedSubscriber sends message to slack #subscriber-pulse.
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
		Color: "#F57F17",
	}
	for _, campaign := range campaigns {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Campaign",
			Value: fmt.Sprintf("Source: %s\nMedium: %s\nCampaign: %s\n", campaign.Source, campaign.Medium, campaign.Campaign),
			Short: true,
		})
	}
	attachments := []slack.Attachment{
		attachment,
	}

	return c.sendAttachmentsMessage(subscriberUpdateChannel, message, attachments)
}

// SendNewSignup sends message to slack #subscriber-pulse.
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
		Color: "#1B5E20",
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

// SendDeactivate sends message to slack #subscriber-pulse.
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
		Color: "#B71C1C",
	}
	attachments := []slack.Attachment{
		attachment,
	}

	return c.sendAttachmentsMessage(subscriberUpdateChannel, message, attachments)
}

// WebhookEvent is an event in WebhookRequest.
type WebhookEvent struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// WebhookRequest is a request body for slack webhook.
type WebhookRequest struct {
	Challenge string       `json:"challenge"`
	Event     WebhookEvent `json:"event"`
	Type      string       `json:"type"`
}

// HandleWebhook handles slack webhook.
func (c *Client) HandleWebhook(req *WebhookRequest) error {
	var err error
	channelInfo, err := c.slackClient.GetChannelInfo(req.Event.Channel)
	if err != nil {
		return errSlack.WithError(err).Annotate("failed to slack.GetChannelInfo")
	}
	c.log.Infof(c.ctx, "channel name: %s", channelInfo.Name)
	if req.Type == "event_callback" {
		if ("#"+channelInfo.Name) == customerSupportChannel && strings.Contains(req.Event.Text, "@") && strings.Contains(req.Event.Text, "<mailto:") {
			// add subscriber link if there is an email
			email := req.Event.Text[strings.Index(req.Event.Text, "<mailto:")+8 : strings.Index(req.Event.Text, "|")]
			txt := fmt.Sprintf("https://eatgigamunch.com/admin/subscriber/%s", email)
			err = c.SendCustomerSupportMessage(txt)
			if err != nil {
				return errSlack.WithError(err).Annotate("failed to SendCustoemrSupportMessage")
			}
		}
	}
	return nil
}

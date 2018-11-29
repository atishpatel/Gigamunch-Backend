package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/slack"

	pbcommon "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/utils"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// SlackResp is a response to slack.
type SlackResp struct {
	Challenge string `json:"challenge"`
}

// GetError completes Response interface.
func (s *SlackResp) GetError() *pbcommon.Error {
	return nil
}

// Slack is a webhook for slack messages.
func (s *server) Slack(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := &slack.WebhookRequest{}
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &SlackResp{
		Challenge: req.Challenge,
	}
	slackC, err := slack.NewClient(ctx, log, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to slack.NewClient")
	}
	err = slackC.HandleWebhook(req)
	if err != nil {
		return errors.Annotate(err, "failed to slack.HandleWebhook")
	}
	return resp
}

// TwilioSMS is a webhook for twilio messages.
func (s *server) TwilioSMS(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	err := r.ParseForm()
	log.Infof(ctx, "req body: %s err: %s", r.Form, err)
	from := r.FormValue("From")
	from = sub.GetCleanPhoneNumber(from)
	body := r.FormValue("Body")
	var name, email, id string
	subC := subold.NewWithLogging(ctx, log)

	messageC := message.New(ctx)
	if message.EmployeeNumbers.IsEmployee(from) {
		// From Gigamunch to Customer
		splitBody := strings.Split(body, "::")
		if len(splitBody) < 2 {
			return nil
		}
		to := sub.GetCleanPhoneNumber(splitBody[0])
		body = splitBody[1]
		subs, err := subC.GetSubscribersByPhoneNumber(to)
		if err != nil {
			log.Errorf(ctx, "failed to sub.GetSubscribersByPhoneNumber: %v", err)
		}
		if len(subs) > 0 {
			email = subs[0].Email
			id = subs[0].ID
		}
		err = messageC.SendDeliverySMS(to, body)
		if err != nil {
			err = messageC.SendDeliverySMS(from, fmt.Sprintf("failed to send sms to %s. Err: %+v", email, err))
			if err != nil {
				utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
			}
		}
		if email != "" {
			payload := &logging.MessagePayload{
				Platform: "SMS",
				Body:     body,
				From:     "Gigamunch",
				To:       to,
			}
			log.SubMessage(id, email, payload)
		}

	} else {
		// From Customer to Gigamunch
		subs, err := subC.GetSubscribersByPhoneNumber(from)
		if err != nil {
			log.Errorf(ctx, "failed to sub.GetSubscribersByPhoneNumber: %v", err)
		}
		if len(subs) > 0 {
			name = subs[0].FirstName + " " + subs[0].LastName
			email = subs[0].Email
			id = subs[0].ID
		}
		// notify customer support agent
		err = messageC.SendDeliverySMS(message.EmployeeNumbers.CustomerSupport(), fmt.Sprintf("Customer Message:\nNumber: %s\nName: %s\nEmail: %s\nBody: %s", from, name, email, body))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to customer support. Err: %+v", err)
		}
		// log
		if email != "" {
			payload := &logging.MessagePayload{
				Platform: "SMS",
				Body:     body,
				From:     sub.GetCleanPhoneNumber(from),
				To:       "Gigamunch",
			}
			log.SubMessage(id, email, payload)
		}
		// check if rating
		if email != "" && (strings.Contains(body, "-") || strings.Contains(body, "star") || strings.Contains(body, "rate") || strings.Contains(body, "rating") || len(body) < 5) {
			regexpNumber := regexp.MustCompile("[0-9]+")
			potentialRatings := regexpNumber.FindAllString(body, -1)
			if len(potentialRatings) > 0 {
				rating, _ := strconv.ParseInt(potentialRatings[0], 10, 8)
				payload := &logging.RatingPayload{
					// TODO: add culture
					Rating:   int8(rating),
					Comments: body,
				}
				log.SubRating(id, email, payload)
			}
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	return nil
}

// TypeformSkip is the webhook for the skip typeform.
func (s *server) TypeformSkip(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	var req typefromWebhookRequest
	// payload, err := ioutil.ReadAll(r.Body)
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&req)
	if err != nil {
		log.Errorf(ctx, "failed to read body: %+v", err)
	}
	logging.Infof(ctx, "decoded req: %+v", req)

	// var email string
	var reason string
	email := req.FormResponse.Hidden.ID

	// get pass reason
	answers := req.FormResponse.Answers
	for _, answer := range answers {
		if answer.Type == "choice" {
			reason = answer.Choice.Label + answer.Choice.Other
		}
		if email == "" && answer.Type == "email" {
			email = answer.Email
		}
	}
	// check if subscriber email is legit
	if email == "" {
		utils.Criticalf(ctx, "failed to get subscriber email from typeform: %+v", err)
		return errBadRequest.WithError(err).Annotate("failed to get subscriber email from typeform")
	}
	email = strings.Replace(email, " ", "+", -1) // replaces space with + in emails
	ctx = context.WithValue(ctx, common.ContextUserEmail, email)
	suboldC := subold.NewWithLogging(ctx, log)
	subscriber, err := suboldC.GetSubscriber(email)
	if err != nil {
		utils.Criticalf(ctx, "failed to find subscriber: %s, they're probably not in our system: %+v", email, err)
		return nil
	}
	if !subscriber.IsSubscribed {
		utils.Criticalf(ctx, "user %s isn't  subscriber and tried to skip.", subscriber.Email)
		return nil
	}

	date := req.FormResponse.SubmittedAt.Add(time.Hour * -5)
	skipDate := date
	for skipDate.Weekday() != time.Monday {
		skipDate = skipDate.Add(time.Hour * 24)
	}

	if date.Weekday() == time.Monday || date.Weekday() == time.Sunday {
		// check for phone number if yes, text them, if no, text me
		if subscriber.PhoneNumber == "" {
			messageC := message.New(ctx)
			err = messageC.SendAdminSMS("6155454989", fmt.Sprintf("What up Chris. Looking fresh today. Nice. 🤠 %s just tried to skip, but it's too late. ", subscriber.Email))
			if err != nil {
				utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
				return nil
			}
			return nil
		}

		messageC := message.New(ctx)
		err = messageC.SendAdminSMS(subscriber.PhoneNumber, fmt.Sprintf("Hey %s, this is Chris from Gigamunch. It looks like you tried to pass a Gigamunch dinner, but we've already started making your meal. 🙊 You need to submit your pass form before Sunday in order for it to work. You will still receive a dinner this Monday.  Feel free to respond directly if you have any questions.", subscriber.FirstName))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to %s. Err: %+v", subscriber.Email, err)
			return errInternalError.WithError(err).Annotatef("failed to send subscriber sms: %s", subscriber.Email)
		}
		return nil
	}
	//if it's Tuesday - Saturday, skip them
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err)
	}
	err = activityC.Skip(skipDate, email, reason)
	if err != nil {
		err = errors.GetErrorWithCode(err).Annotate("failed to sub.Skip")
		utils.Criticalf(ctx, "Typeform webhook: %+v", err)
		return errors.GetErrorWithCode(err)
	}
	return nil
}

type typefromWebhookRequest struct {
	EventID      string           `json:"event_id,omitempty"`
	EventType    string           `json:"event_type,omitempty"`
	FormResponse typeformResponse `json:"form_response,omitempty"`
}

type typeformResponse struct {
	FormID      string           `json:"form_id,omitempty"`
	Token       string           `json:"token,omitempty"`
	SubmittedAt time.Time        `json:"submitted_at,omitempty"`
	Hidden      hiddenField      `json:"hidden,omitempty"`
	Answers     []typeformAnswer `json:"answers,omitempty"`
}

type hiddenField struct {
	ID string `json:"id,omitempty"`
}

type typeformAnswer struct {
	Type   string         `json:"type,omitempty"`
	Text   string         `json:"text,omitempty"`
	Email  string         `json:"email,omitempty"`
	Choice typeformChoice `json:"choice,omitempty"`
	Field  typeformField  `json:"field,omitempty"`
}

type typeformField struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

type typeformChoice struct {
	Label string `json:"label,omitempty"`
	Other string `json:"other,omitempty"`
}

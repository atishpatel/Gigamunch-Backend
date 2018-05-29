package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/utils"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func setupWebhooksHandlers() {
	http.HandleFunc("/admin/webhook/typeform-skip", handler(TypeformSkip))
}

// TypeformSkip is the webhook for the skip typeform.
func TypeformSkip(ctx context.Context, r *http.Request, log *logging.Client) Response {
	var err error
	var req TypefromWebhookRequest
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
	subC := sub.New(ctx)
	subscriber, err := subC.GetSubscriber(email)
	if err != nil {
		utils.Criticalf(ctx, "failed to get subscriber Err: %+v", err)
		return errBadRequest.WithError(err).Annotate("failed to get subscriber")
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


	subC := sub.New(ctx)
	subscriber, err := subC.GetSubscriber(email)
	if err != nil {
		utils.Criticalf(ctx, "failed to find subscriber: %s, they're probably not in our system: %+v", email, err)
		return nil
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
	err = subC.Skip(skipDate, email)
	if err != nil {
		err = errors.GetErrorWithCode(err).Annotate("failed to sub.Skip")
		utils.Criticalf(ctx, "Typeform webhook: %+v", err)
		return errors.GetErrorWithCode(err)
	}

	// TODO: logging caused error, will look into later
	_ = reason
	// log.SubSkip(skipDate.Format(time.RFC3339), 0, email, reason)
	return nil
}

type TypefromWebhookRequest struct {
	EventID      string           `json:"event_id,omitempty"`
	EventType    string           `json:"event_type,omitempty"`
	FormResponse TypeformResponse `json:"form_response,omitempty"`
}

type TypeformResponse struct {
	FormID      string           `json:"form_id,omitempty"`
	Token       string           `json:"token,omitempty"`
	SubmittedAt time.Time        `json:"submitted_at,omitempty"`
	Hidden      HiddenField      `json:"hidden,omitempty"`
	Answers     []TypeformAnswer `json:"answers,omitempty"`
}

type HiddenField struct {
	ID string `json:"id,omitempty"`
}

type TypeformAnswer struct {
	Type   string         `json:"type,omitempty"`
	Text   string         `json:"text,omitempty"`
	Email  string         `json:"email,omitempty"`
	Choice TypeformChoice `json:"choice,omitempty"`
	Field  TypeformField  `json:"field,omitempty"`
}

type TypeformField struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

type TypeformChoice struct {
	Label string `json:"label,omitempty"`
	Other string `json:"other,omitempty"`
}

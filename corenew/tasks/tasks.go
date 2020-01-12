package tasks

import (
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine/taskqueue"
)

const (
	UpdateDripQueue          = "update-drip"
	UpdateDripURL            = "/task/update-drip"
	SendSMSQueue             = "send-sms"
	SendSMSURL               = "/admin/task/SendSMS"
	ProcessSubscriptionQueue = "process-subscription"
	ProcessSubscriptionURL   = "/admin/task/ProcessActivity"
	// SendEmailQueue           = "send-email"
	// SendEmailURL             = "/send-email"
)

var (
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter for tasks."}
	errTasks            = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Something went wrong with tasks."}
	errParse            = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Bad request."}
)

// Client is a client for Tasks.
type Client struct {
	ctx context.Context
}

// New returns a new Client
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// ProcessSubscriptionParams are the parms for ProcessSubscription.
type ProcessSubscriptionParams struct {
	SubEmail string // Depecrated
	UserID   string
	Date     time.Time
}

// AddProcessSubscription adds a process subscription at specified time.
func (c *Client) AddProcessSubscription(at time.Time, req *ProcessSubscriptionParams) error {
	if (req.UserID == "" && req.SubEmail == "") || req.Date.IsZero() {
		return errInvalidParameter.Wrapf("expected(recieved): email(%s) date(%s)", req.SubEmail, req.Date.String())
	}
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("sub_email", req.SubEmail)
	v.Set("user_id", req.UserID)
	d, err := req.Date.MarshalText()
	if err != nil {
		return errInvalidParameter.WithError(err).Wrap("failed to req.Date.MarshalText")
	}
	v.Set("date", string(d))
	task := &taskqueue.Task{
		Path:    ProcessSubscriptionURL,
		Payload: []byte(v.Encode()),
		Header:  h,
		Method:  "POST",
		ETA:     at,
	}
	_, err = taskqueue.Add(c.ctx, task, ProcessSubscriptionQueue)
	if err != nil {
		return errTasks.WithError(err).Wrapf("failed to task.Add. Task: %v", task)
	}
	return nil
}

// UpdateDripParams are the parms for UpdateDrip.
type UpdateDripParams struct {
	UserID string
	Email  string
}

// AddUpdateDrip adds a process subscription at specified time.
func (c *Client) AddUpdateDrip(at time.Time, req *UpdateDripParams) error {
	var err error
	if req.Email == "" && req.UserID == "" {
		return errInvalidParameter.Wrapf("expected(recieved): email(%s) date(%s)", req.Email)
	}
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("email", req.Email)
	v.Set("user_id", req.UserID)
	task := &taskqueue.Task{
		Path:    UpdateDripURL,
		Payload: []byte(v.Encode()),
		Header:  h,
		Method:  "POST",
		ETA:     at,
	}
	_, err = taskqueue.Add(c.ctx, task, UpdateDripQueue)
	if err != nil {
		return errTasks.WithError(err).Wrapf("failed to task.Add. Task: %v", task)
	}
	utils.Infof(c.ctx, "added tasks update drip at %s for %s", at, req.Email)
	return nil
}

// ParseUpdateDripRequest parses an UpdateDripRequest from a task request.
func ParseUpdateDripRequest(req *http.Request) (*UpdateDripParams, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, errParse.WithError(err).Wrap("failed to parse from from request")
	}
	parms := new(UpdateDripParams)
	parms.Email = req.FormValue("email")
	parms.UserID = req.FormValue("user_id")
	if parms.Email == "" && parms.UserID == "" {
		return nil, errParse.Wrapf("Invalid request for UpdateDrip. SubEmail: %s", parms.Email)
	}
	return parms, nil
}

// ParseProcessSubscriptionRequest parses an ProcessSubscriptionRequest from a task request.
func ParseProcessSubscriptionRequest(req *http.Request) (*ProcessSubscriptionParams, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, errParse.WithError(err).Wrap("failed to parse from from request")
	}
	parms := new(ProcessSubscriptionParams)
	parms.SubEmail = req.FormValue("sub_email")
	parms.UserID = req.FormValue("user_id")
	if parms.SubEmail == "" && parms.UserID == "" {
		return nil, errParse.Wrapf("Invalid request for ProcessSubscription. SubEmail: %s", parms.SubEmail)
	}
	dateString := req.FormValue("date")
	err = parms.Date.UnmarshalText([]byte(dateString))
	if err != nil {
		return nil, errParse.WithError(err).Wrap("failed to parse date")
	}
	return parms, nil
}

// SendEmailParams are the parms for SendEmail.
// type SendEmailParams struct {
// 	Email string
// 	Type  string
// }

// // AddSendEmail adds a email to send at specified time.
// func (c *Client) AddSendEmail(at time.Time, req *SendEmailParams) error {
// 	if req.Email == "" {
// 		return errInvalidParameter.Wrapf("expected(recieved): email(%s)", req.Email)
// 	}
// 	h := make(http.Header)
// 	h.Set("Content-Type", "application/x-www-form-urlencoded")
// 	v := url.Values{}
// 	v.Set("email", req.Email)
// 	v.Set("type", req.Type)
// 	task := &taskqueue.Task{
// 		Path:    SendEmailURL,
// 		Payload: []byte(v.Encode()),
// 		Header:  h,
// 		Method:  "POST",
// 		ETA:     at,
// 	}
// 	_, err := taskqueue.Add(c.ctx, task, SendEmailQueue)
// 	if err != nil {
// 		return errTasks.WithError(err).Wrapf("failed to task.Add. Task: %v", task)
// 	}
// 	return nil
// }

// // ParseSendEmailRequest parses an SendEmailRequest from a task request.
// func ParseSendEmailRequest(req *http.Request) (*SendEmailParams, error) {
// 	err := req.ParseForm()
// 	if err != nil {
// 		return nil, errParse.WithError(err).Wrap("failed to parse from from request")
// 	}
// 	parms := new(SendEmailParams)
// 	parms.Email = req.FormValue("email")
// 	parms.Type = req.FormValue("type")
// 	if parms.Email == "" {
// 		return nil, errParse.Wrapf("Invalid request for SendEmail. Email: %s", parms.Email)
// 	}
// 	return parms, nil
// }

type SendSMSParam struct {
	Number  string
	Email   string
	Message string
}

// AddSendSMS sends an sms at certain time and at a rate.
func (c *Client) AddSendSMS(req *SendSMSParam, at time.Time) error {
	var err error
	if req.Message == "" {
		return errInvalidParameter.Wrap("bad param for SendSMS")
	}
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("email", req.Email)
	v.Set("number", req.Number)
	v.Set("message", req.Message)
	task := &taskqueue.Task{
		Path:    SendSMSURL,
		Payload: []byte(v.Encode()),
		Header:  h,
		Method:  "POST",
		ETA:     at,
	}
	_, err = taskqueue.Add(c.ctx, task, SendSMSQueue)
	if err != nil {
		return errTasks.WithError(err).Wrapf("failed to task.Add. Task: %v", task)
	}
	utils.Infof(c.ctx, "added tasks SendSMS at %s for %s - %s", at, req.Email, req.Number)
	return nil
}

// ParseSendSMSParam parses an SendSMSParam from a task request.
func ParseSendSMSParam(req *http.Request) (*SendSMSParam, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, errParse.WithError(err).Wrap("failed to parse from from request")
	}
	parms := new(SendSMSParam)
	parms.Email = req.FormValue("email")
	parms.Number = req.FormValue("number")
	parms.Message = req.FormValue("message")
	if parms.Number == "" && parms.Email == "" {
		return nil, errParse.Wrapf("Invalid request for SendSMS")
	}
	return parms, nil
}

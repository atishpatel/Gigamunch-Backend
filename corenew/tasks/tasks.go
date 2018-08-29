package tasks

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine/taskqueue"
)

const (
	UpdateDripQueue          = "update-drip"
	UpdateDripURL            = "/task/update-drip"
	ProcessInquiryQueue      = "process-inquiry"
	ProcessInquiryURL        = "/process-inquiry"
	ProcessSubscriptionQueue = "process-subscription"
	ProcessSubscriptionURL   = "/process-subscription"
	SendEmailQueue           = "send-email"
	SendEmailURL             = "/send-email"
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
	SubEmail string
	Date     time.Time
}

// AddProcessSubscription adds a process subscription at specified time.
func (c *Client) AddProcessSubscription(at time.Time, req *ProcessSubscriptionParams) error {
	if req.SubEmail == "" || req.Date.IsZero() {
		return errInvalidParameter.Wrapf("expected(recieved): email(%s) date(%s)", req.SubEmail, req.Date.String())
	}
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("sub_email", req.SubEmail)
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
	Email string
}

// AddUpdateDrip adds a process subscription at specified time.
func (c *Client) AddUpdateDrip(at time.Time, req *UpdateDripParams) error {
	var err error
	if req.Email == "" {
		return errInvalidParameter.Wrapf("expected(recieved): email(%s) date(%s)", req.Email)
	}
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("email", req.Email)
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
	if parms.Email == "" {
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
	if parms.SubEmail == "" {
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
type SendEmailParams struct {
	Email string
	Type  string
}

// AddSendEmail adds a email to send at specified time.
func (c *Client) AddSendEmail(at time.Time, req *SendEmailParams) error {
	if req.Email == "" {
		return errInvalidParameter.Wrapf("expected(recieved): email(%s)", req.Email)
	}
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("email", req.Email)
	v.Set("type", req.Type)
	task := &taskqueue.Task{
		Path:    SendEmailURL,
		Payload: []byte(v.Encode()),
		Header:  h,
		Method:  "POST",
		ETA:     at,
	}
	_, err := taskqueue.Add(c.ctx, task, SendEmailQueue)
	if err != nil {
		return errTasks.WithError(err).Wrapf("failed to task.Add. Task: %v", task)
	}
	return nil
}

// ParseSendEmailRequest parses an SendEmailRequest from a task request.
func ParseSendEmailRequest(req *http.Request) (*SendEmailParams, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, errParse.WithError(err).Wrap("failed to parse from from request")
	}
	parms := new(SendEmailParams)
	parms.Email = req.FormValue("email")
	parms.Type = req.FormValue("type")
	if parms.Email == "" {
		return nil, errParse.Wrapf("Invalid request for SendEmail. Email: %s", parms.Email)
	}
	return parms, nil
}

// AddProcessInquiry adds a process inquiry at specified time.
func (c *Client) AddProcessInquiry(inquiryID int64, at time.Time) error {
	if inquiryID == 0 {
		return errInvalidParameter.Wrap("inquriyID cannot be 0")
	}
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("inquiry_id", strconv.FormatInt(inquiryID, 10))
	task := &taskqueue.Task{
		Path:    ProcessInquiryURL,
		Payload: []byte(v.Encode()),
		Header:  h,
		Method:  "POST",
		ETA:     at,
	}
	_, err := taskqueue.Add(c.ctx, task, ProcessInquiryQueue)
	if err != nil {
		return errTasks.WithError(err).Wrapf("failed to task.Add. Task: %v", task)
	}
	return nil
}

// ParseInquiryID parses an inquiryID from a task request.
func ParseInquiryID(req *http.Request) (int64, error) {
	err := req.ParseForm()
	if err != nil {
		return 0, errParse.WithError(err).Wrap("failed to parse from from request")
	}
	inquiryIDString := req.FormValue("inquiry_id")
	if inquiryIDString == "" {
		return 0, errParse.Wrapf("Invalid task for close inquiry. inquiryID: %s", inquiryIDString)
	}
	inquiryID, err := strconv.ParseInt(inquiryIDString, 10, 64)
	if err != nil {
		return 0, errParse.WithError(err).Wrapf("Failed to parse inquiryID(%s). Err: ", inquiryIDString, err)
	}
	return inquiryID, nil
}

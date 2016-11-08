package tasks

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"

	"google.golang.org/appengine/taskqueue"
)

const (
	// NotifyEaterQueue    = "notify-eater"
	// NotifyEaterURL      = "/notify-eater"
	ProcessInquiryQueue = "process-inquiry"
	ProcessInquiryURL   = "/process-inquiry"
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

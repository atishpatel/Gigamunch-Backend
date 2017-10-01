package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"

	sdlogging "cloud.google.com/go/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	aelog "google.golang.org/appengine/log"
)

const (
	kind = "logging"
)

// Label of log.
type Label string

const (
	// Paid Label.
	Paid = Label("paid")
	// Decline Label.
	Decline = Label("decline")
	// Refund Label.
	Refund = Label("refund")
	// Forgiven Label.
	Forgiven = Label("forgiven")
	// DeliveryStarted Label.
	DeliveryStarted = Label("delivery_started")
	// DeliveryEnded Label.
	DeliveryEnded = Label("delivery_ended")
	// Delivered Label.
	Delivered = Label("delivered")
	// Update Label.
	Update = Label("update")
	// Cancel Label.
	Cancel = Label("cancel")
	// Signup Label.
	Signup = Label("signup")
	// Activate Label.
	Activate = Label("activate")
	// Deactivate Label.
	Deactivate = Label("deactivate")
)

// Type of log.
type Type string

func (t *Type) isNil() bool {
	return string(*t) == ""
}

const (
	// Unknown type.
	Unknown = Type("unknown")
	// Request type.
	Request = Type("request")
	// Sale type.
	Sale = Type("sale")
	// Activity type.
	Activity = Type("activity")
	// Subscriber Type.
	Subscriber = Type("subscriber")
	// Delivery Type.
	Delivery = Type("delivery")
	// System type.
	System = Type("system")
	// CultureExecution Type.
	// CultureExecution = Type("culture_execution")
	// Notification Type.
	// Notification = Type("notification")
)

const (
	// Default means the log entry has no assigned severity level.
	Default = sdlogging.Default
	// Debug means debug or trace information.
	Debug = sdlogging.Debug
	// Info means routine information, such as ongoing status or performance.
	Info = sdlogging.Info
	// Warning means events that might cause problems.
	Warning = sdlogging.Warning
	// Error means events that are likely to cause problems.
	Error = sdlogging.Error
	// Critical means events that cause more severe problems or brief outages.
	Critical = sdlogging.Critical
	// Alert means a person must take an action immediately.
	Alert = sdlogging.Alert
	// Emergency means one or more systems are unusable.
	// Emergency = sdlogging.Emergency

)

var (
	projID   string
	loggerID string
	sdClient *sdlogging.Client
	db       common.DB
)

var (
	// Errors
	errFailedToEncodeJSON = errors.ErrorWithCode{
		Code:    errors.CodeBadRequest,
		Message: "Failed to log.",
		Detail:  "failed to encode log json",
	}
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
)

// Infof logs info.
func Infof(ctx context.Context, format string, args ...interface{}) {
	aelog.Infof(ctx, format, args...)
}

// Errorf logs error messages.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	aelog.Errorf(ctx, format, args...)
}

// Debugf logs debug messages. Only logs on development servers.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	aelog.Debugf(ctx, format, args...)
}

// Client is a logging client.
type Client struct {
	ctx      context.Context
	sdLogger *sdlogging.Logger
	path     string
}

// NewClient returns a new logging client.
func NewClient(ctx context.Context, path string) (*Client, error) {
	if db == nil {
		return nil, errInternal.Annotate("setup not called")
	}
	var sdLogger *sdlogging.Logger
	if sdClient != nil {
		sdLogger = sdClient.Logger(loggerID)
	}
	return &Client{
		ctx:      ctx,
		sdLogger: sdLogger,
		path:     path,
	}, nil
}

// Debugf logs debug messages. Only logs on development servers.
func (c *Client) Debugf(ctx context.Context, format string, args ...interface{}) {
	Debugf(ctx, format, args...)
}

// GetLogs gets logs.
func (c *Client) GetLogs(start, limit int) ([]*Entry, error) {
	var dst []*Entry
	keys, err := db.Query(c.ctx, kind, start, limit, "-Timestamp", dst)
	if err != nil {
		return nil, errors.Annotate(err, "failed to db.Query")
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	return dst, nil
}

// GetLog gets a log.
func (c *Client) GetLog(id int64) (*Entry, error) {
	var entry *Entry
	key := db.IDKey(c.ctx, kind, id)
	err := db.Get(c.ctx, key, entry)
	if err != nil {
		return nil, errors.Annotate(err, "failed to db.Get")
	}
	entry.ID = id
	return entry, nil
}

// SaleEntry is a sales entry.
type SaleEntry struct {
	Amount float32 `json:"amount"`
}

// LogPaid is when a transaction is paid.
func (c *Client) LogPaid(e *SaleEntry) {

}

// LogRefund is when a transaction is refunded.
func (c *Client) LogRefund(e *SaleEntry) {

}

// LogDeclined is when a transaction is declined.
func (c *Client) LogDeclined(e *SaleEntry) {

}

// LogForgiven is when a transaction is forgiven.
func (c *Client) LogForgiven(e *SaleEntry) {

}

// SubEntry is a subscriber entry.
type SubEntry struct {
	ID string `json:"id"`
}

func (c *Client) LogSubActivate(e *SubEntry) {

}

func (c *Client) LogSubDeactivate(e *SubEntry) {

}

func (c *Client) LogSubCancel(e *SubEntry) {

}

func (c *Client) LogSubSignup(e *SubEntry) {

}

func (c *Client) LogSubUpdate(e *SubEntry) {

}

// ActivityEntry is a Activity entry.
type ActivityEntry struct {
	ActionUserID   string    `json:"action_user_id"`
	ActionUserName string    `json:"action_user_name"`
	Date           time.Time `json:"date"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
}

// LogSkip logs a skip.
func (c *Client) LogSkip(e *ActivityEntry) {

}

// LogUnskip logs a unskip.
func (c *Client) LogUnskip(e *ActivityEntry) {

}

// LogServingsChanged logs a servings change.
func (c *Client) LogServingsChanged(e *ActivityEntry) {

}

// SystemEntry is a System entry.
type SystemEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// LogActivitySetup is a log of when the cron job for activity setup runs.
func (c *Client) LogActivitySetup(e *SystemEntry) {

}

// ErrorEntry is an error entry assocted with LogRequestError.
type ErrorEntry struct {
	Request http.Request
	errors.ErrorWithCode
}

// LogRequestError is used to log an error at the end of a request.
func (c *Client) LogRequestError(r *http.Request, ewc errors.ErrorWithCode) {
	errEntry := &ErrorEntry{
		Request:       *r,
		ErrorWithCode: ewc,
	}
	e := &Entry{
		Type:     Request,
		Severity: Error,
		Path:     r.URL.Path,
	}
	err := e.setPayload(errEntry)
	if err != nil {
		Errorf(c.ctx, "failed to setPayload: %+v", err)
	}
	c.Log(e)
}

// Entry is a log entry.
type Entry struct {
	ID        int64              `json:"id" datastore:",noindex"`
	Type      Type               `json:"type" datastore:",index"`
	Severity  sdlogging.Severity `json:"serverity" datastore:",noindex"`
	Path      string             `json:"path" datastore:",noindex"`
	Labels    []Label            `json:"labels" datastore:",noindex"`
	LogName   string             `json:"log_name" datastore:",noindex"`
	Timestamp time.Time          `json:"timestamp" datastore:",index"`
	Payload   string             `json:"payload" datastore:",noindex"`
}

func (e *Entry) setPayload(payload interface{}) error {
	if payload == nil {
		return errFailedToEncodeJSON.Annotate("payload is empty")
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return errFailedToEncodeJSON.WithError(err)
	}
	e.Payload = string(b)
	return nil
}

// Log logs a random entry.
func (c *Client) Log(e *Entry) {
	if e.Type.isNil() {
		e.Type = Unknown
	}
	e.LogName = projID + "/" + loggerID
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	if e.Path == "" {
		e.Path = c.path
	}
	key := db.IncompleteKey(c.ctx, kind)
	_, err := db.Put(c.ctx, key, e)
	if err != nil {
		Errorf(c.ctx, "failed to log entry(%+v) error: %+v", e, err)
	}
}

// Setup sets up the logging package.
func Setup(ctx context.Context, projectID, logID string, httpClient *http.Client, dbC common.DB) error {
	projID = projectID
	loggerID = logID
	if dbC == nil {
		return fmt.Errorf("db cannot be nil for logging")
	}
	db = dbC
	// var ops option.ClientOption
	// if httpClient != nil {
	// 	ops = option.WithHTTPClient(httpClient)
	// }
	// var err error
	// sdClient, err = sdlogging.NewClient(ctx, projectID, ops)
	// if err != nil {
	// 	return err
	// }
	return nil
}

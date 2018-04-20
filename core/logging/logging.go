package logging

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"

	sdlogging "cloud.google.com/go/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"google.golang.org/api/option"
	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const (
	kind = "logging"
)

// Label of log.
type Label string

const (
	// ======================
	// Admin Actions
	// ======================

	// Delivered Label.
	Delivered = Label("delivered")
	// ======================
	// User or Admin Actions
	// ======================

	// Update Label.
	Update = Label("update")
	// Cancel Label.
	Cancel = Label("cancel")
	// Activate Label.
	Activate = Label("activate")
	// Deactivate Label.
	Deactivate = Label("deactivate")
	// Refund Label.
	Refund = Label("refund")
	// Forgiven Label.
	Forgiven = Label("forgiven")
	// ======================
	// User Actions
	// ======================

	// ======================
	// System Actions
	// ======================

	// Paid Label.
	Paid = Label("paid")
	// Decline Label.
	Decline = Label("decline")
)

// Type of log.
type Type string

func (t *Type) isNil() bool {
	return string(*t) == ""
}

const (
	// Unknown type.
	Unknown = Type("unknown")
	// AdminAction type.
	AdminAction = Type("admin_action")
	// UserAction Type.
	UserAction = Type("user_action")
	// System type.
	System = Type("system")
	// Error type.
	Error = Type("error")
)

const (
	// SeverityInfo means routine information, such as ongoing status or performance.
	SeverityInfo = sdlogging.Info
	// SeverityWarning means events that might cause problems.
	SeverityWarning = sdlogging.Warning
	// SeverityError means events that are likely to cause problems.
	SeverityError = sdlogging.Error
	// SeverityCritical means events that cause more severe problems or brief outages.
	SeverityCritical = sdlogging.Critical
	// SeverityAlert means a person must take an action immediately.
	SeverityAlert = sdlogging.Alert
)

var (
	standAppEngine bool
	projID         string
	loggerID       string
	sdClient       *sdlogging.Client
	db             common.DB
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
	if standAppEngine {
		httpClient := urlfetch.Client(ctx)
		setup(ctx, httpClient)
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

// Errorf logs error messages.
func (c *Client) Errorf(ctx context.Context, format string, args ...interface{}) {
	Errorf(ctx, format, args...)
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

// GetUserLogs gets logs with UserID.
func (c *Client) GetUserLogs(userID int64, start, limit int) ([]*Entry, error) {
	var dst []*Entry
	keys, err := db.QueryFilterOrdered(c.ctx, kind, start, limit, "-Timestamp", "UserID=", userID, dst)
	if err != nil {
		return nil, errors.Annotate(err, "failed to db.QueryFilterOrdered")
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

// SalePayload is a sales payload.
type SalePayload struct {
	Amount     float32 `json:"amount"`
	AmountPaid float32
}

// LogPaid is when a transaction is paid.
func (c *Client) LogPaid(e *SalePayload) {

}

// LogRefund is when a transaction is refunded.
func (c *Client) LogRefund(e *SalePayload) {

}

// LogDeclined is when a transaction is declined.
func (c *Client) LogDeclined(e *SalePayload) {

}

// LogForgiven is when a transaction is forgiven.
func (c *Client) LogForgiven(e *SalePayload) {

}

// SubPayload is a subscriber entry.
type SubPayload struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
}

func (c *Client) LogSubActivate(e *SubPayload) {

}

func (c *Client) LogSubDeactivate(e *SubPayload) {

}

func (c *Client) LogSubCancel(e *SubPayload) {

}

func (c *Client) LogSubUpdate(e *SubPayload) {

}

// ActivityPayload is a Activity entry.
type ActivityPayload struct {
	ActionUserID   string    `json:"action_user_id"`
	ActionUserName string    `json:"action_user_name"`
	Date           time.Time `json:"date"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
}

// LogSkip logs a skip.
func (c *Client) LogSkip(e *ActivityPayload) {

}

// LogUnskip logs a unskip.
func (c *Client) LogUnskip(e *ActivityPayload) {

}

// LogServingsChanged logs a servings change.
func (c *Client) LogServingsChanged(e *ActivityPayload) {

}

// SystemPayload is a System payload.
type SystemPayload struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// LogActivitySetup is a log of when the cron job for activity setup runs.
func (c *Client) LogActivitySetup(e *SystemPayload) {

}

// ErrorPayload is an error entry assocted with LogRequestError.
type ErrorPayload struct {
	Request http.Request // TODO: change this
	errors.ErrorWithCode
}

// LogRequestError is used to log an error at the end of a request.
// TODO: log body?
func (c *Client) LogRequestError(r *http.Request, ewc errors.ErrorWithCode) {
	e := &Entry{
		Type:     Error,
		Severity: SeverityError,
		Path:     r.URL.Path,
		ErrorPayload: ErrorPayload{
			Request:       *r,
			ErrorWithCode: ewc,
		},
	}
	c.Log(e)
}

// Entry is a log entry.
type Entry struct {
	ID        int64              `json:"id" datastore:",noindex"`
	Type      Type               `json:"type" datastore:",index"`
	UserID    int64              `json:"user_id" datastore:",index"`
	Severity  sdlogging.Severity `json:"serverity" datastore:",noindex"`
	Path      string             `json:"path" datastore:",noindex"`
	Labels    []Label            `json:"labels" datastore:",noindex"`
	LogName   string             `json:"log_name" datastore:",noindex"`
	Timestamp time.Time          `json:"timestamp" datastore:",index"`
	// ActivityPayload ActivityPayload    `json:"activity_payload" datastore:",noindex"`
	ErrorPayload ErrorPayload `json:"error_payload" datastore:",noindex"`
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
func Setup(ctx context.Context, standardAppEngine bool, projectID, logID string, httpClient *http.Client, dbC common.DB) error {
	projID = projectID
	loggerID = logID
	standAppEngine = standardAppEngine
	if dbC == nil {
		return fmt.Errorf("db cannot be nil for logging")
	}
	db = dbC
	if !standAppEngine {
		setup(ctx, httpClient)
	}
	return nil
}

func setup(ctx context.Context, httpClient *http.Client) error {
	var ops option.ClientOption
	if httpClient != nil {
		ops = option.WithHTTPClient(httpClient)
	}
	var err error
	sdClient, err = sdlogging.NewClient(ctx, projID, ops)
	if err != nil {
		return err
	}
	return nil
}

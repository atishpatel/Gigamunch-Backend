package logging

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-test/deep"

	"github.com/atishpatel/Gigamunch-Backend/core/common"

	sdreporting "cloud.google.com/go/errorreporting"

	sdlogging "cloud.google.com/go/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	aelog "google.golang.org/appengine/log"
)

const (
	kind = "logging"
)

// Action of log.
type Action string

const (
	// ======================
	// Admin only Actions
	// ======================

	// Delivered Action.
	Delivered = Action("delivered")
	// ExecutionUpdate Action.
	ExecutionUpdate = Action("execution_update")
	// ======================
	// User or Admin Actions
	// ======================

	// FailedSkip Action.
	FailedSkip = Action("failed_skip")
	// Skip Action.
	Skip = Action("skip")
	// Unskip Action.
	Unskip = Action("unskip")
	// Update Action.
	Update = Action("update")
	// Cancel Action.
	Cancel = Action("cancel")
	// Activate Action.
	Activate = Action("activate")
	// Deactivate Action.
	Deactivate = Action("deactivate")
	// Refund Action.
	Refund = Action("refund")
	// Forgiven Action.
	Forgiven = Action("forgiven")
	// CardUpdated Action.
	CardUpdated = Action("card_updated")
	// ServingsChanged Action.
	ServingsChanged = Action("servings_changed")
	// ServingsChangedPermanently Action.
	ServingsChangedPermanently = Action("servings_changed_permanently")
	// ======================
	// User only Actions
	// ======================

	// Message Action.
	Message = Action("message")
	// Rating Action.
	Rating = Action("rating")
	// ======================
	// System Actions
	// ======================

	// Paid Action.
	Paid = Action("paid")
	// Decline Action.
	Decline = Action("decline")
)

// Type of log.
type Type string

func (t *Type) isNil() bool {
	return string(*t) == ""
}

const (
	// Unknown type.
	Unknown = Type("unknown")
	// Subscriber type.
	Subscriber = Type("subscriber")
	// Execution type.
	Execution = Type("execution")
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
	// Errors
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

// Client is a logging client.
type Client struct {
	ctx        context.Context
	sdReporter *sdreporting.Client
	loggerID   string
	path       string
	db         common.DB
	serverInfo *common.ServerInfo
}

// NewClient returns a new logging client.
func NewClient(ctx context.Context, loggerID string, path string, dbC common.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if dbC == nil {
		return nil, errInternal.Annotate("dbC cannot be nil")
	}
	sdReporter, err := sdreporting.NewClient(context.Background(), serverInfo.ProjectID, sdreporting.Config{
		ServiceName: loggerID,
		OnError: func(err error) {
			Errorf(ctx, "failed to log to sdReporter", err)
		},
	})
	if err != nil {
		Errorf(ctx, "failed to create sdReporter", err)
	}
	return &Client{
		ctx:        ctx,
		sdReporter: sdReporter,
		loggerID:   loggerID,
		path:       path,
		db:         dbC,
		serverInfo: serverInfo,
	}, nil
}

// SetContext sets the context.
func (c *Client) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// Infof logs info.
func (c *Client) Infof(ctx context.Context, format string, args ...interface{}) {
	Infof(ctx, format, args...)
}

// Errorf logs error messages.
func (c *Client) Errorf(ctx context.Context, format string, args ...interface{}) {
	if c.sdReporter != nil {
		c.sdReporter.Report(sdreporting.Entry{
			Error: fmt.Errorf(format, args),
		})
	}
	Errorf(ctx, format, args...)
}

// Criticalf logs error messages that is critical and needs to be alerted and address immediately.
func (c *Client) Criticalf(ctx context.Context, format string, args ...interface{}) {
	if c.sdReporter != nil {
		c.sdReporter.Report(sdreporting.Entry{
			Error: fmt.Errorf("CRITICAL: "+format, args),
		})
	}
	aelog.Criticalf(ctx, format, args...)
}

// GetAll gets logs.
func (c *Client) GetAll(start, limit int) ([]*Entry, error) {
	var dst []*Entry
	keys, err := c.db.Query(c.ctx, kind, start, limit, "-Timestamp", &dst)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to db.QueryFilterOrdered")
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	return dst, nil
}

// GetAllByID gets logs with UserID.
func (c *Client) GetAllByID(userID string, start, limit int) ([]*Entry, error) {
	var dst []*Entry
	keys, err := c.db.QueryFilterOrdered(c.ctx, kind, start, limit, "-Timestamp", "UserIDString=", userID, &dst)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to db.QueryFilterOrdered")
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	return dst, nil
}

// GetAllByEmail gets logs with UserEmail.
func (c *Client) GetAllByEmail(userEmail string, start, limit int) ([]*Entry, error) {
	var dst []*Entry
	keys, err := c.db.QueryFilterOrdered(c.ctx, kind, start, limit, "-Timestamp", "UserEmail=", userEmail, &dst)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to db.QueryFilterOrdered")
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	return dst, nil
}

// GetAllByExecution gets logs by ExecutionID.
func (c *Client) GetAllByExecution(executionID int64) ([]*Entry, error) {
	var dst []*Entry
	keys, err := c.db.QueryFilter(c.ctx, kind, 0, 1000, "ExecutionID=", executionID, &dst)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to db.QueryFilter")
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	sort.Slice(dst, func(i, j int) bool {
		return dst[j].Timestamp.Before(dst[i].Timestamp)
	})
	return dst, nil
}

// Get gets a log.
func (c *Client) Get(id int64) (*Entry, error) {
	var entry Entry
	key := c.db.IDKey(c.ctx, kind, id)
	err := c.db.Get(c.ctx, key, &entry)
	if err != nil {
		return nil, errors.Annotate(err, "failed to db.Get")
	}
	entry.ID = id
	return &entry, nil
}

// SalePayload is a sales payload.
type SalePayload struct {
	Date           string  `json:"date,omitempty" datastore:",omitempty,noindex"`
	AmountDue      float32 `json:"amount_due,omitempty" datastore:",omitempty,noindex"`
	AmountPaid     float32 `json:"amount_paid,omitempty" datastore:",omitempty,noindex"`
	AmountDeclined float32 `json:"amount_declined,omitempty" datastore:",omitempty,noindex"`
	AmountRefunded float32 `json:"amount_refunded,omitempty" datastore:",omitempty,noindex"`
	AmountForgiven float32 `json:"amount_forgiven,omitempty" datastore:",omitempty,noindex"`
	TransactionID  string  `json:"transaction_id,omitempty" datastore:",omitempty,noindex"`
}

// Paid is when a transaction is paid.
func (c *Client) Paid(userID string, userEmail, date string, amountDue, amountPaid float32, transactionID string) {
	e := &Entry{
		Type:         Subscriber,
		Action:       Paid,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Paid for " + date,
			Description: fmt.Sprintf("%s successfully paid %.2f for %s", userEmail, amountPaid, date),
		},
		SalePayload: SalePayload{
			Date:          date,
			AmountDue:     amountDue,
			AmountPaid:    amountPaid,
			TransactionID: transactionID,
		},
	}
	c.Log(e)
}

// Refund is when a transaction is refunded.
func (c *Client) Refund(userID string, userEmail, date string, amountDue, amountRefunded float32, transactionID string) {
	e := &Entry{
		Type:         Subscriber,
		Action:       Refund,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Refunded for " + date,
			Description: fmt.Sprintf("%s was refunded %.2f", userEmail, amountRefunded),
		},
		SalePayload: SalePayload{
			Date:           date,
			AmountDue:      amountDue,
			AmountRefunded: amountRefunded,
			TransactionID:  transactionID,
		},
	}
	c.Log(e)
}

// Forgiven is when a transaction is forgiven.
func (c *Client) Forgiven(userID string, userEmail, date string, amountDue, amountForgiven float32) {
	e := &Entry{
		Type:         Subscriber,
		Action:       Forgiven,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Forgiven for " + date,
			Description: fmt.Sprintf("%s was forgiven for %.2f", userEmail, amountForgiven),
		},
		SalePayload: SalePayload{
			Date:           date,
			AmountDue:      amountDue,
			AmountForgiven: amountForgiven,
		},
	}
	c.Log(e)
}

// CardDeclined is when a transaction is declined.
func (c *Client) CardDeclined(userID string, userEmail, date string, amountDue, amountDeclined float32, transactionID string) {
	e := &Entry{
		Type: Subscriber,
		// TODO: Action?
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Card declined for " + date,
			Description: fmt.Sprintf("%s's card was declined for %.2f", userEmail, amountDeclined),
		},
		SalePayload: SalePayload{
			Date:           date,
			AmountDue:      amountDue,
			AmountDeclined: amountDeclined,
			TransactionID:  transactionID,
		},
	}
	c.Log(e)
}

// CreditCardPayload is the payload related to CreditCards.
type CreditCardPayload struct {
	OldPaymentMethodToken string `json:"old_payment_method_token,omitempty" datastore:",omitempty,noindex"`
	NewPaymentMethodToken string `json:"new_payment_method_token,omitempty" datastore:",omitempty,noindex"`
}

// SubCardUpdated is when a credit card is Updated.
func (c *Client) SubCardUpdated(userID string, userEmail, oldPaymentMethodToken, newPaymentMethodToken string) {
	e := &Entry{
		Type:         Subscriber,
		Action:       CardUpdated,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Changed Credit Card",
			Description: fmt.Sprintf("Changed card from %s to %s", oldPaymentMethodToken, newPaymentMethodToken),
		},
		CreditCardPayload: CreditCardPayload{
			OldPaymentMethodToken: oldPaymentMethodToken,
			NewPaymentMethodToken: newPaymentMethodToken,
		},
	}
	c.Log(e)
}

// SubUpdated is when a subscriber account is updated.
func (c *Client) SubUpdated(userID string, userEmail string, oldSub, newSub interface{}) {
	desc := getDiff(oldSub, newSub)
	e := &Entry{
		Type:         Subscriber,
		Action:       Update,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Updated Subscriber Info",
			Description: desc,
		},
	}
	c.Log(e)
}

// func (c *Client) SubActivate() {

// }

type SubDeactivatedPayload struct {
	Rason      string `json:"reason,omitempty" datastore:",omitempty,noindex"`
	DaysActive int    `json:"days_active,omitempty" datastore:",omitempty,noindex"`
}

func (c *Client) SubDeactivated(userID, userEmail, reason string, daysActive int) {
	desc := fmt.Sprintf("Rason: %s. Deactivated after %d weeks (%d days).", reason, daysActive/7, daysActive)
	e := &Entry{
		Type:         Subscriber,
		Action:       Deactivate,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Deactivated Subscriber",
			Description: desc,
		},
		SubDeactivatedPayload: SubDeactivatedPayload{
			Rason:      reason,
			DaysActive: daysActive,
		},
	}
	c.Log(e)
}

// MessagePayload is realted to a subscriber message interaction.
type MessagePayload struct {
	Platform string `json:"platform,omitempty" datastore:",omitempty,noindex"`
	Subject  string `json:"subject,omitempty" datastore:",omitempty,noindex"`
	Body     string `json:"body,omitempty" datastore:",omitempty,noindex"`
	From     string `json:"from,omitempty" datastore:",omitempty,noindex"`
	To       string `json:"to,omitempty" datastore:",omitempty,noindex"`
}

// SubMessage is a message interaction with a sub.
func (c *Client) SubMessage(userID string, userEmail string, payload *MessagePayload) {
	e := &Entry{
		Type:         Subscriber,
		Action:       Message,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Message from(" + payload.From + ") to(" + payload.To + ")",
			Description: "Body: '" + payload.Body + "'",
		},
		MessagePayload: *payload,
	}
	c.Log(e)
}

// RatingPayload is realted to a subscriber rating.
type RatingPayload struct {
	Rating      int8   `json:"rating,omitempty" datastore:",omitempty,noindex"`
	CultureDate string `json:"culture_date,omitempty" datastore:",omitempty,noindex"`
	Culture     string `json:"culture,omitempty" datastore:",omitempty,noindex"`
	Comments    string `json:"comments,omitempty" datastore:",omitempty,noindex"`
}

// SubRating is a message interaction with a sub.
func (c *Client) SubRating(userID string, userEmail string, payload *RatingPayload) {
	e := &Entry{
		Type:         Subscriber,
		Action:       Rating,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       fmt.Sprintf("Rating %d", payload.Rating),
			Description: "Comments: '" + payload.Comments + "'",
		},
		RatingPayload: *payload,
	}
	c.Log(e)
}

// ServingsChangedPayload is a ServingsChanged entry.
type ServingsChangedPayload struct {
	Date              string `json:"date,omniempty" datastore:",omitempty,noindex"`
	OldNonVegServings int8   `json:"old_non_veg_servings,omitempty" datastore:",omitempty,noindex"`
	NewNonVegServings int8   `json:"new_non_veg_servings,omitempty" datastore:",omitempty,noindex"`
	OldVegServings    int8   `json:"old_veg_servings,omitempty" datastore:",omitempty,noindex"`
	NewVegServings    int8   `json:"new_veg_servings,omitempty" datastore:",omitempty,noindex"`
}

// SubServingsChangedPermanently logs a servings change.
func (c *Client) SubServingsChangedPermanently(userID string, userEmail string, oldNonVegServings, newNonVegServings, oldVegServings, newVegServings int8) {
	e := &Entry{
		Type:         Subscriber,
		Action:       ServingsChangedPermanently,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Servings changed permanently",
			Description: fmt.Sprintf("Servings changed from %d to %d non-veg and %d to %d veg", oldNonVegServings, newNonVegServings, oldVegServings, newVegServings),
		},
		ServingsChangedPayload: ServingsChangedPayload{
			OldNonVegServings: oldNonVegServings,
			NewNonVegServings: newNonVegServings,
			OldVegServings:    oldVegServings,
			NewVegServings:    newVegServings,
		},
	}
	c.Log(e)
}

// SubServingsChanged logs a servings change.
func (c *Client) SubServingsChanged(userID string, userEmail string, date string, oldNonVegServings, newNonVegServings, oldVegServings, newVegServings int8) {
	e := &Entry{
		Type:         Subscriber,
		Action:       ServingsChangedPermanently,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Servings changed for " + date,
			Description: fmt.Sprintf("Servings changed from %d to %d non-veg and %d to %d veg", oldNonVegServings, newNonVegServings, oldVegServings, newVegServings),
		},
		ServingsChangedPayload: ServingsChangedPayload{
			Date:              date,
			OldNonVegServings: oldNonVegServings,
			NewNonVegServings: newNonVegServings,
			OldVegServings:    oldVegServings,
			NewVegServings:    newVegServings,
		},
	}
	c.Log(e)
}

// SkipPayload is a Skip entry.
type SkipPayload struct {
	UserID             int64  `json:"-" datastore:",omitempty,noindex"`
	UserIDString       string `json:"user_id,omitempty" datastore:",omitempty,noindex"`
	UserEmail          string `json:"user_email,omitempty" datastore:",omitempty,noindex"`
	Reason             string `json:"reason,omitempty" datastore:",omitempty,noindex"`
	ActionUserID       int64  `json:"-" datastore:",omitempty,noindex"`
	ActionUserIDString string `json:"action_user_id,omitempty" datastore:",omitempty,noindex"`
	ActionUserEmail    string `json:"action_user_email,omitempty" datastore:",omitempty,noindex"`
	Date               string `json:"date,omitempty" datastore:",omitempty,noindex"`
}

// SubSkip logs a skip.
func (c *Client) SubSkip(date string, userID string, userEmail, reason string) {
	actionUserEmail := c.getStringFromCtx(common.ContextUserEmail)
	e := &Entry{
		Type:         Subscriber,
		Action:       Skip,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Skip for " + date,
			Description: fmt.Sprintf("%s was skipped for %s by %s because '%s'", userEmail, date, actionUserEmail, reason),
		},
		SkipPayload: SkipPayload{
			Date:               date,
			UserIDString:       userID,
			UserEmail:          userEmail,
			Reason:             reason,
			ActionUserIDString: c.getStringFromCtx(common.ContextUserID),
			ActionUserEmail:    actionUserEmail,
		},
	}
	c.Log(e)
}

// SubUnskip logs a unskip.
func (c *Client) SubUnskip(date string, userID string, userEmail string) {
	actionUserEmail := c.getStringFromCtx(common.ContextUserEmail)
	e := &Entry{
		Type:         Subscriber,
		Action:       Unskip,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Unskip for " + date,
			Description: fmt.Sprintf("%s was unskipped for %s by %s", userEmail, date, actionUserEmail),
		},
		SkipPayload: SkipPayload{
			Date:               date,
			UserIDString:       userID,
			UserEmail:          userEmail,
			ActionUserIDString: c.getStringFromCtx(common.ContextUserID),
			ActionUserEmail:    actionUserEmail,
		},
	}
	c.Log(e)
}

// SubFailedSkip logs a skip.
func (c *Client) SubFailedSkip(date string, userID string, userEmail, reason string) {
	actionUserEmail := c.getStringFromCtx(common.ContextUserEmail)
	e := &Entry{
		Type:         Subscriber,
		Action:       FailedSkip,
		Severity:     SeverityInfo,
		UserIDString: userID,
		UserEmail:    userEmail,
		BasicPayload: BasicPayload{
			Title:       "Failed to skip for " + date,
			Description: fmt.Sprintf("%s failed to skip for %s by %s", userEmail, date, actionUserEmail),
		},
		SkipPayload: SkipPayload{
			Date:               date,
			UserIDString:       userID,
			UserEmail:          userEmail,
			Reason:             reason,
			ActionUserIDString: c.getStringFromCtx(common.ContextUserID),
			ActionUserEmail:    actionUserEmail,
		},
	}
	c.Log(e)
}

// // ActivitySetupPayload is a System payload.
type ActivitySetupPayload struct {
	BasicPayload `json:"basic_payload,omitempty" datastore:",omitempty,noindex"`
	Date         string `json:"date,omitempty" datastore:",omitempty,noindex"`
	NumSetup     int    `json:"num_setup,omitempty" datastore:",omitempty,noindex"`
}

// // ActivitySetup is a log of when the cron job for activity setup runs or admin runs activity setup.
// func (c *Client) ActivitySetup(date string, numSetup int) {
// 	e := &Entry{
// 		Type:     System,
// 		Severity: SeverityInfo,
// 		BasicPayload: BasicPayload{
// 			Title:       date,
// 			Description: fmt.Sprintf("Activity setup for %s", date),
// 		},
// 		ActivitySetupPayload: ActivitySetupPayload{
// 			Date:     date,
// 			NumSetup: numSetup,
// 		},
// 	}
// 	c.Log(e)
// }

// ExecutionUpdate is a log of when an execution is udated.
func (c *Client) ExecutionUpdate(executionID int64, oldExe interface{}, newExe interface{}) {
	desc := getDiff(oldExe, newExe)
	e := &Entry{
		Type:        Execution,
		Action:      ExecutionUpdate,
		Severity:    SeverityInfo,
		ExecutionID: executionID,
		BasicPayload: BasicPayload{
			Title:       "Execution Updated",
			Description: desc,
		},
	}
	c.Log(e)
}

// ErrorPayload is an error entry assocted with RequestError.
type ErrorPayload struct {
	Method        string               `json:"method,omitempty" datastore:",omitempty,noindex"`
	URL           string               `json:"url,omitempty" datastore:",omitempty,noindex"`
	Proto         string               `json:"proto,omitempty" datastore:",omitempty,noindex"`
	ContentLength int64                `json:"content_length,omitempty" datastore:",omitempty,noindex"`
	Host          string               `json:"host,omitempty" datastore:",omitempty,noindex"`
	Error         errors.ErrorWithCode `json:"error,omitempty" datastore:",omitempty,noindex"`
}

// RequestError is used to log an error at the end of a request.
func (c *Client) RequestError(r *http.Request, ewc errors.ErrorWithCode) {
	e := &Entry{
		Type:     Error,
		Severity: SeverityError,
		Path:     r.URL.Path,
		ErrorPayload: ErrorPayload{
			Method:        r.Method,
			URL:           r.URL.String(),
			Proto:         r.Proto,
			ContentLength: r.ContentLength,
			Host:          r.Host,
			Error:         ewc,
		},
	}
	c.Log(e)
	if c.sdReporter != nil {
		if ewc.Code >= 400 && ewc.Code < 500 {
			c.sdReporter.Report(sdreporting.Entry{
				Error: fmt.Errorf("req bad: %v", ewc),
				Req:   r,
				User:  e.UserEmail,
			})
		} else {
			c.sdReporter.Report(sdreporting.Entry{
				Error: fmt.Errorf("req err: %v", ewc),
				Req:   r,
				User:  e.UserEmail,
			})
		}
	}
}

// BasicPayload is in every payload.
type BasicPayload struct {
	Title       string `json:"title" datastore:",omitempty,noindex`
	Description string `json:"description" datastore:",omitempty,noindex`
}

// Entry is a log entry.
type Entry struct {
	ID                     int64                  `json:"id,omitempty" datastore:",noindex"`
	Type                   Type                   `json:"type,omitempty" datastore:",index"`
	ExecutionID            int64                  `json:"execution_id" datastore:",index"`
	Action                 Action                 `json:"action,omitempty" datastore:",index"`
	ActionUserIDString     string                 `json:"action_user_id,omitempty" datastore:",index"`
	ActionUserEmail        string                 `json:"action_user_email,omitempty" datastore:",index"`
	UserIDString           string                 `json:"user_id,omitempty" datastore:",index"`
	UserEmail              string                 `json:"user_email,omitempty" datastore:",index"`
	Severity               sdlogging.Severity     `json:"serverity,omitempty" datastore:",noindex"`
	Path                   string                 `json:"path,omitempty" datastore:",noindex"`
	LogName                string                 `json:"log_name,omitempty" datastore:",noindex"`
	Timestamp              time.Time              `json:"timestamp,omitempty" datastore:",index"`
	BasicPayload           BasicPayload           `json:"basic_payload,omitempty" datastore:",noindex"`
	ErrorPayload           ErrorPayload           `json:"error_payload,omitempty" datastore:",omitempty,noindex"`
	SkipPayload            SkipPayload            `json:"skip_payload,omitempty" datastore:",omitempty,noindex"`
	CreditCardPayload      CreditCardPayload      `json:"credit_card_payload,omitempty" datastore:",omitempty,noindex"`
	SalePayload            SalePayload            `json:"sale_payload,omitempty" datastore:",omitempty,noindex"`
	ServingsChangedPayload ServingsChangedPayload `json:"servings_changed_payload,omitempty" datastore:",omitempty,noindex"`
	MessagePayload         MessagePayload         `json:"message_payload,omitempty" datastore:",omitempty,noindex"`
	RatingPayload          RatingPayload          `json:"rating_payload,omitempty" datastore:",omitempty,noindex"`
	SubDeactivatedPayload  SubDeactivatedPayload  `json:"sub_deactivated_payload,omitempty" datastore:",omitempty,noindex"`
	// deprecated
	UserID               int64                `json:"-" datastore:",omitempty,noindex"`                                // deprecated
	ActionUserID         int64                `json:"-" datastore:",omitempty,noindex"`                                // deprecated
	SubUpdatedPayload    SubUpdatedPayload    `json:"sub_updated_payload,omitempty" datastore:",omitempty,noindex"`    // deprecated
	ActivitySetupPayload ActivitySetupPayload `json:"activity_setup_payload,omitempty" datastore:",omitempty,noindex"` // deprecated
}

// Log logs a random entry.
func (c *Client) Log(e *Entry) {
	if e.Type.isNil() {
		e.Type = Unknown
	}
	e.LogName = c.serverInfo.ProjectID + "/" + c.loggerID
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	if e.Path == "" {
		e.Path = c.path
	}
	if e.ActionUserIDString == "" {
		e.ActionUserIDString = c.getStringFromCtx(common.ContextUserID)
	}
	if e.ActionUserEmail == "" {
		e.ActionUserEmail = c.getStringFromCtx(common.ContextUserEmail)
	}
	key := c.db.IncompleteKey(c.ctx, kind)
	_, err := c.db.Put(c.ctx, key, e)
	if err != nil {
		Errorf(c.ctx, "failed to log entry(%+v) error: %+v", e, err)
	}
}

func (c *Client) getStringFromCtx(key interface{}) string {
	v := c.ctx.Value(key)
	if v == nil {
		return ""
	}
	return v.(string)
}

func getDiff(od interface{}, nw interface{}) string {
	diffs := deep.Equal(od, nw)
	desc := ""
	replacer := strings.NewReplacer("!=", "->", ".slice", "", " slice", "")
	for _, diff := range diffs {
		desc += replacer.Replace(diff) + ";;;\n"
	}
	return desc
}

// SubUpdatedPayload is the payload related to SubUpdated.
// DEPECRATED
type SubUpdatedPayload struct {
	OldEmail          string `json:"old_email,omitempty" datastore:",omitempty,noindex"`
	Email             string `json:"email,omitempty" datastore:",omitempty,noindex"`
	OldFirstName      string `json:"old_first_name,omitempty" datastore:",omitempty,noindex"`
	FirstName         string `json:"first_name,omitempty" datastore:",omitempty,noindex"`
	OldLastName       string `json:"old_last_name,omitempty" datastore:",omitempty,noindex"`
	LastName          string `json:"last_name,omitempty" datastore:",omitempty,noindex"`
	OldAddress        string `json:"old_address,omitempty" datastore:",omitempty,noindex"`
	Address           string `json:"address,omitempty" datastore:",omitempty,noindex"`
	OldRawPhoneNumber string `json:"old_raw_phone_number,omitempty" datastore:",omitempty,noindex"`
	RawPhoneNumber    string `json:"raw_phone_number,omitempty" datastore:",omitempty,noindex"`
	OldPhoneNumber    string `json:"old_phone_number,omitempty" datastore:",omitempty,noindex"`
	PhoneNumber       string `json:"phone_number,omitempty" datastore:",omitempty,noindex"`
	OldDeliveryNotes  string `json:"old_delivery_tip,omitempty" datastore:",omitempty,noindex"`
	DeliveryNotes     string `json:"delivery_tip,omitempty" datastore:",omitempty,noindex"`
}

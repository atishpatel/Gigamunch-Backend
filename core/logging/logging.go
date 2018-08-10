package logging

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

// Action of log.
type Action string

const (
	// ======================
	// Admin only Actions
	// ======================

	// Delivered Action.
	Delivered = Action("delivered")
	// ======================
	// User or Admin Actions
	// ======================

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

// GetAll gets logs.
func (c *Client) GetAll(start, limit int) ([]*Entry, error) {
	var dst []*Entry
	keys, err := db.Query(c.ctx, kind, start, limit, "-Timestamp", &dst)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to db.QueryFilterOrdered")
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	return dst, nil
}

// GetAllByID gets logs with UserID.
func (c *Client) GetAllByID(userID int64, start, limit int) ([]*Entry, error) {
	var dst []*Entry
	keys, err := db.QueryFilterOrdered(c.ctx, kind, start, limit, "-Timestamp", "UserID=", userID, &dst)
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
	keys, err := db.QueryFilterOrdered(c.ctx, kind, start, limit, "-Timestamp", "UserEmail=", userEmail, &dst)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to db.QueryFilterOrdered")
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	return dst, nil
}

// Get gets a log.
func (c *Client) Get(id int64) (*Entry, error) {
	var entry Entry
	key := db.IDKey(c.ctx, kind, id)
	err := db.Get(c.ctx, key, &entry)
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
func (c *Client) Paid(userID int64, userEmail, date string, amountDue, amountPaid float32, transactionID string) {
	e := &Entry{
		Type:      Subscriber,
		Action:    Paid,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
func (c *Client) Refund(userID int64, userEmail, date string, amountDue, amountRefunded float32, transactionID string) {
	e := &Entry{
		Type:      Subscriber,
		Action:    Refund,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
func (c *Client) Forgiven(userID int64, userEmail, date string, amountDue, amountForgiven float32) {
	e := &Entry{
		Type:      Subscriber,
		Action:    Forgiven,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
func (c *Client) CardDeclined(userID int64, userEmail, date string, amountDue, amountDeclined float32, transactionID string) {
	e := &Entry{
		Type: Subscriber,
		// TODO: Action?
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
func (c *Client) SubCardUpdated(userID int64, userEmail, oldPaymentMethodToken, newPaymentMethodToken string) {
	e := &Entry{
		Type:      Subscriber,
		Action:    CardUpdated,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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

// SubUpdatedPayload is the payload related to SubUpdated.
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

// SubUpdated is when a subscriber account is updated.
func (c *Client) SubUpdated(userID int64, userEmail string, payload *SubUpdatedPayload) {
	desc := "Changed the following: "
	if payload.OldEmail != payload.Email {
		desc += "Email(" + payload.OldEmail + " -> " + payload.Email + "); "
	}
	if payload.OldFirstName != payload.FirstName {
		desc += "FirstName(" + payload.OldFirstName + " -> " + payload.FirstName + "); "
	}
	if payload.OldLastName != payload.LastName {
		desc += "LastName(" + payload.OldLastName + " -> " + payload.LastName + "); "
	}
	if payload.OldAddress != payload.Address {
		desc += "Address(" + payload.OldAddress + " -> " + payload.Address + "); "
	}
	if payload.OldRawPhoneNumber != payload.RawPhoneNumber {
		desc += "RawPhoneNumber(" + payload.OldRawPhoneNumber + " -> " + payload.RawPhoneNumber + "); "
	}
	if payload.OldPhoneNumber != payload.PhoneNumber {
		desc += "PhoneNumber(" + payload.OldPhoneNumber + " -> " + payload.PhoneNumber + "); "
	}
	if payload.OldDeliveryNotes != payload.DeliveryNotes {
		desc += "DeliveryNotes(" + payload.OldDeliveryNotes + " -> " + payload.DeliveryNotes + "); "
	}
	e := &Entry{
		Type:      Subscriber,
		Action:    Update,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
		BasicPayload: BasicPayload{
			Title:       "Updated Subscriber Info",
			Description: desc,
		},
		SubUpdatedPayload: *payload,
	}
	c.Log(e)
}

// func (c *Client) SubActivate() {

// }

// func (c *Client) SubDeactivate() {

// }

// MessagePayload is realted to a subscriber message interaction.
type MessagePayload struct {
	Platform string `json:"platform,omitempty" datastore:",omitempty,noindex"`
	Subject  string `json:"subject,omitempty" datastore:",omitempty,noindex"`
	Body     string `json:"body,omitempty" datastore:",omitempty,noindex"`
	From     string `json:"from,omitempty" datastore:",omitempty,noindex"`
	To       string `json:"to,omitempty" datastore:",omitempty,noindex"`
}

// SubMessage is a message interaction with a sub.
func (c *Client) SubMessage(userID int64, userEmail string, payload *MessagePayload) {
	e := &Entry{
		Type:      Subscriber,
		Action:    Message,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
func (c *Client) SubRating(userID int64, userEmail string, payload *RatingPayload) {
	e := &Entry{
		Type:      Subscriber,
		Action:    Rating,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
func (c *Client) SubServingsChangedPermanently(userID int64, userEmail string, oldNonVegServings, newNonVegServings, oldVegServings, newVegServings int8) {
	e := &Entry{
		Type:      Subscriber,
		Action:    ServingsChangedPermanently,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
func (c *Client) SubServingsChanged(userID int64, userEmail string, date string, oldNonVegServings, newNonVegServings, oldVegServings, newVegServings int8) {
	e := &Entry{
		Type:      Subscriber,
		Action:    ServingsChangedPermanently,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
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
	UserID          int64  `json:"user_id,omitempty" datastore:",omitempty,noindex"`
	UserEmail       string `json:"user_email,omitempty" datastore:",omitempty,noindex"`
	Reason          string `json:"reason,omitempty" datastore:",omitempty,noindex"`
	ActionUserID    int64  `json:"action_user_id,omitempty" datastore:",omitempty,noindex"`
	ActionUserEmail string `json:"action_user_email,omitempty" datastore:",omitempty,noindex"`
	Date            string `json:"date,omitempty" datastore:",omitempty,noindex"`
}

// SubSkip logs a skip.
func (c *Client) SubSkip(date string, userID int64, userEmail, reason string) {
	actionUserEmail := c.ctx.Value(common.ContextUserEmail).(string)
	e := &Entry{
		Type:      Subscriber,
		Action:    Skip,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
		BasicPayload: BasicPayload{
			Title:       "Skip for " + date,
			Description: fmt.Sprintf("%s was skipped for %s by %s because '%s'", userEmail, date, actionUserEmail, reason),
		},
		SkipPayload: SkipPayload{
			Date:            date,
			UserID:          userID,
			UserEmail:       userEmail,
			Reason:          reason,
			ActionUserID:    c.ctx.Value(common.ContextUserID).(int64),
			ActionUserEmail: actionUserEmail,
		},
	}
	c.Log(e)
}

// SubUnskip logs a unskip.
func (c *Client) SubUnskip(date string, userID int64, userEmail string) {
	actionUserEmail := c.ctx.Value(common.ContextUserEmail).(string)
	e := &Entry{
		Type:      Subscriber,
		Action:    Unskip,
		Severity:  SeverityInfo,
		UserID:    userID,
		UserEmail: userEmail,
		BasicPayload: BasicPayload{
			Title:       "Unskip for " + date,
			Description: fmt.Sprintf("%s was unskipped for %s by %s", userEmail, date, actionUserEmail),
		},
		SkipPayload: SkipPayload{
			Date:            date,
			UserID:          userID,
			UserEmail:       userEmail,
			ActionUserID:    c.ctx.Value(common.ContextUserID).(int64),
			ActionUserEmail: actionUserEmail,
		},
	}
	c.Log(e)
}

// ActivitySetupPayload is a System payload.
type ActivitySetupPayload struct {
	BasicPayload `json:"basic_payload,omitempty" datastore:",omitempty,noindex"`
	Date         string `json:"date,omitempty" datastore:",omitempty,noindex"`
	NumSetup     int    `json:"num_setup,omitempty" datastore:",omitempty,noindex"`
}

// ActivitySetup is a log of when the cron job for activity setup runs or admin runs activity setup.
func (c *Client) ActivitySetup(date string, numSetup int) {
	e := &Entry{
		Type:     System,
		Severity: SeverityInfo,
		BasicPayload: BasicPayload{
			Title:       date,
			Description: fmt.Sprintf("Activity setup for %s", date),
		},
		ActivitySetupPayload: ActivitySetupPayload{
			Date:     date,
			NumSetup: numSetup,
		},
	}
	c.Log(e)
}

// ErrorPayload is an error entry assocted with RequestError.
type ErrorPayload struct {
	Method        string               `json:"method,omitempty" datastore:",omitempty,noindex"`
	URL           string               `json:"url,omitempty" datastore:",omitempty,noindex"`
	Proto         string               `json:"proto,omitempty" datastore:",omitempty,noindex"`
	Header        http.Header          `json:"header,omitempty" datastore:",omitempty,noindex"`
	ContentLength int64                `json:"content_length,omitempty" datastore:",omitempty,noindex"`
	Host          string               `json:"host,omitempty" datastore:",omitempty,noindex"`
	Form          url.Values           `json:"form,omitempty" datastore:",omitempty,noindex"`
	Error         errors.ErrorWithCode `json:"error,omitempty" datastore:",omitempty,noindex"`
}

// RequestError is used to log an error at the end of a request.
// TODO: log body?
func (c *Client) RequestError(r *http.Request, ewc errors.ErrorWithCode, userID int64, userEmail string) {
	e := &Entry{
		Type:      Error,
		Severity:  SeverityError,
		Path:      r.URL.Path,
		UserID:    userID,
		UserEmail: userEmail,
		ErrorPayload: ErrorPayload{
			Method:        r.Method,
			URL:           r.URL.String(),
			Header:        r.Header,
			Proto:         r.Proto,
			ContentLength: r.ContentLength,
			Host:          r.Host,
			Form:          r.Form,
			Error:         ewc,
		},
	}
	c.Log(e)
}

// BasicPayload is in every payload.
type BasicPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Entry is a log entry.
type Entry struct {
	ID                     int64                  `json:"id,omitempty" datastore:",noindex"`
	Type                   Type                   `json:"type,omitempty" datastore:",index"`
	Action                 Action                 `json:"action,omitempty" datastore:",index"`
	ActionUserID           int64                  `json:"action_user_id,omitempty" datastore:",index"`
	ActionUserEmail        string                 `json:"action_user_email,omitempty" datastore:",index"`
	UserID                 int64                  `json:"user_id,omitempty" datastore:",index"`
	UserEmail              string                 `json:"user_email,omitempty" datastore:",index"`
	Severity               sdlogging.Severity     `json:"serverity,omitempty" datastore:",noindex"`
	Path                   string                 `json:"path,omitempty" datastore:",noindex"`
	LogName                string                 `json:"log_name,omitempty" datastore:",noindex"`
	Timestamp              time.Time              `json:"timestamp,omitempty" datastore:",index"`
	BasicPayload           BasicPayload           `json:"basic_payload,omitempty" datastore:",noindex"`
	ErrorPayload           ErrorPayload           `json:"error_payload,omitempty" datastore:",omitempty,noindex"`
	ActivitySetupPayload   ActivitySetupPayload   `json:"activity_setup_payload,omitempty" datastore:",omitempty,noindex"`
	SkipPayload            SkipPayload            `json:"skip_payload,omitempty" datastore:",omitempty,noindex"`
	CreditCardPayload      CreditCardPayload      `json:"credit_card_payload,omitempty" datastore:",omitempty,noindex"`
	SalePayload            SalePayload            `json:"sale_payload,omitempty" datastore:",omitempty,noindex"`
	ServingsChangedPayload ServingsChangedPayload `json:"servings_changed_payload,omitempty" datastore:",omitempty,noindex"`
	SubUpdatedPayload      SubUpdatedPayload      `json:"sub_updated_payload,omitempty" datastore:",omitempty,noindex"`
	MessagePayload         MessagePayload         `json:"message_payload,omitempty" datastore:",omitempty,noindex"`
	RatingPayload          RatingPayload          `json:"rating_payload,omitempty" datastore:",omitempty,noindex"`
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
	if e.ActionUserID == 0 {
		tmp := c.ctx.Value(common.ContextUserID)
		if tmp != nil {
			e.ActionUserID = tmp.(int64)
		}
	}
	if e.ActionUserEmail == "" {
		tmp := c.ctx.Value(common.ContextUserEmail)
		if tmp != nil {
			e.ActionUserEmail = tmp.(string)
		}
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

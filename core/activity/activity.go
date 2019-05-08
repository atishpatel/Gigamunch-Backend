package activity

import (
	"context"
	"fmt"
	"strings"
	"time"

	// mysql driver
	mysql "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
)

const (
	datetimeFormat = "2006-01-02 15:04:05" // "Jan 2, 2006 at 3:04pm (MST)"
	// DateFormat is the expected format of date in activity.
	DateFormat                               = "2006-01-02" // "Jan 2, 2006"
	selectActivityByEmailStatement           = "SELECT * FROM activity WHERE date=? AND email=?"
	selectActivityByUserIDStatement          = "SELECT * FROM activity WHERE date=? AND user_id=?"
	selectAllActivityStatement               = "SELECT * FROM activity ORDER BY date DESC LIMIT ?"
	selectActivityForUserStatement           = "SELECT * FROM activity WHERE user_id=? ORDER BY date DESC"
	selectActivityAfterDateForUserStatement  = "SELECT * FROM activity WHERE user_id=? AND date>=? ORDER BY date DESC"
	selectActivityBeforeDateForUserStatement = "SELECT * FROM activity WHERE user_id=? AND date<=? ORDER BY date DESC"
	updateRefundedStatement                  = "UPDATE activity SET refunded_dt=NOW(),refunded=1,refund_transaction_id=?,refunded_amount=? WHERE date=? AND email=?"
	skipStatement                            = "UPDATE activity SET skip=1 WHERE date=? AND email=?"
	unskipStatement                          = "UPDATE activity SET skip=0,active=1 WHERE date=? AND email=?"
	insertStatement                          = "INSERT INTO activity (date,user_id,email,first_name,last_name,location,addr_apt,addr_string,zip,lat,`long`,active,skip,servings,veg_servings,first,amount,discount_amount,discount_percent,payment_provider,payment_method_token,customer_id) VALUES (:date,:user_id,:email,:first_name,:last_name,:location,:addr_apt,:addr_string,:zip,:lat,:long,:active,:skip,:servings,:veg_servings,:first,:amount,:discount_amount,:discount_percent,:payment_provider,:payment_method_token,:customer_id)"
	deleteFutureStatment                     = "DELETE from activity WHERE date>? AND user_id=? AND paid=0"
	updatePaidStatement                      = "UPDATE activity SET amount_paid=?,paid=1,paid_dt=?,transaction_id=? WHERE date=? AND user_id=?"
	// TODO: switch to user id
	deleteFutureEmailStatment = "DELETE from activity WHERE date>? AND email=? AND paid=0"
)

// Errors
var (
	errInternal       = errors.InternalServerError
	errBadRequest     = errors.BadRequestError
	errSQLDB          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	errDuplicateEntry = errors.ErrorWithCode{Code: errors.CodeBadRequest, Message: "Duplicate entry."}
)

// Client is a client for manipulating activity.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	sqlDB      *sqlx.DB
	db         common.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	if sqlC == nil {
		return nil, fmt.Errorf("sqlDB cannot be nil")
	}
	if dbC == nil {
		return nil, fmt.Errorf("failed to get db")
	}
	if serverInfo == nil {
		return nil, errInternal.Annotate("failed to get server info")
	}
	return &Client{
		ctx:        ctx,
		log:        log,
		sqlDB:      sqlC,
		db:         dbC,
		serverInfo: serverInfo,
	}, nil
}

// Get gets an activity.
func (c *Client) Get(date time.Time, idOrEmail string) (*Activity, error) {
	var err error
	act := &Activity{}
	if strings.Contains(idOrEmail, "@") {
		err = c.sqlDB.GetContext(c.ctx, act, selectActivityByEmailStatement, date.Format(DateFormat), idOrEmail)
	} else {
		err = c.sqlDB.GetContext(c.ctx, act, selectActivityByUserIDStatement, date.Format(DateFormat), idOrEmail)
	}
	if err != nil {
		// TODO: Wrap with errSQL
		return nil, errSQLDB.WithError(err).Annotate("failed to selectActivity")
	}
	return act, nil
}

// GetAll gets a list of activity.
func (c *Client) GetAll(limit int) ([]*Activity, error) {
	if limit <= 0 {
		limit = 1000
	}
	acts := []*Activity{}
	err := c.sqlDB.SelectContext(c.ctx, &acts, selectAllActivityStatement, limit)
	if err != nil {
		return nil, errSQLDB.WithError(err).Annotate("failed to selectAllActivityStatement")
	}
	return acts, nil
}

// GetAllForUser gets a list of activity for a user.
func (c *Client) GetAllForUser(userID string) ([]*Activity, error) {
	acts := []*Activity{}
	err := c.sqlDB.SelectContext(c.ctx, &acts, selectActivityForUserStatement, userID)
	if err != nil {
		return nil, errSQLDB.WithError(err).Annotate("failed to selectActivityForUserStatement")
	}
	return acts, nil
}

// GetAfterDateForUser gets a list of activity for a user.
func (c *Client) GetAfterDateForUser(date time.Time, userID string) ([]*Activity, error) {
	acts := []*Activity{}
	err := c.sqlDB.SelectContext(c.ctx, &acts, selectActivityAfterDateForUserStatement, userID, date.Format(DateFormat))
	if err != nil {
		return nil, errSQLDB.WithError(err).Annotate("failed to selectActivityAfterDateForUserStatement")
	}
	return acts, nil
}

// GetBeforeDateForUser gets a list of activity for a user.
func (c *Client) GetBeforeDateForUser(date time.Time, userID string) ([]*Activity, error) {
	acts := []*Activity{}
	err := c.sqlDB.SelectContext(c.ctx, &acts, selectActivityBeforeDateForUserStatement, userID, date.Format(DateFormat))
	if err != nil {
		return nil, errSQLDB.WithError(err).Annotate("failed to selectActivityBeforerDateForUserStatement")
	}
	return acts, nil
}

// TODO: GetForDate
//

// CreateReq is the request for Create.
type CreateReq struct {
	Date      string          `json:"date" db:"date"`
	UserID    string          `json:"user_id" db:"user_id"`
	Email     string          `json:"email" db:"email"`
	FirstName string          `json:"first_name" db:"first_name"`
	LastName  string          `json:"last_name" db:"last_name"`
	Location  common.Location `json:"location" db:"location"`
	// Address
	AddressAPT    string  `json:"address_apt" db:"addr_apt"`
	AddressString string  `json:"address_string" db:"addr_string"`
	Zip           string  `json:"zip" db:"zip"`
	Latitude      float64 `json:"latitude,string" db:"lat"`
	Longitude     float64 `json:"longitude,string" db:"long"`
	// Detail
	Active bool `json:"active" db:"active"`
	Skip   bool `json:"skip" db:"skip"`
	// Bag detail
	ServingsNonVegetarian int8 `json:"servings_non_vegetarian" db:"servings"`
	ServingsVegetarain    int8 `json:"servings_vegetarian" db:"veg_servings"`
	First                 bool `json:"first" db:"first"`
	// Payment
	Amount             float32                `json:"amount" db:"amount"`
	DiscountAmount     float32                `json:"discount_amount" db:"discount_amount"`
	DiscountPercent    int8                   `json:"discount_percent" db:"discount_percent"`
	PaymentProvider    common.PaymentProvider `json:"payment_provider" db:"payment_provider"`
	PaymentMethodToken string                 `json:"payment_method_token" db:"payment_method_token"`
	CustomerID         string                 `json:"customer_id" db:"customer_id"`
}

// SetAddress is a convenience function for setting address.
func (req *CreateReq) SetAddress(addr *common.Address) {
	req.AddressAPT = addr.APT
	req.AddressString = addr.StringNoAPT()
	req.Zip = addr.Zip
	req.Latitude = addr.Latitude
	req.Longitude = addr.Longitude
}

func (req *CreateReq) validate() error {
	if req.Date == "" {
		return errBadRequest.WithMessage("Date cannot be empty.")
	}
	if req.UserID == "" {
		return errBadRequest.WithMessage("UserID cannot be empty.")
	}
	if req.Email == "" {
		return errBadRequest.WithMessage("Email cannot be empty.")
	}
	if req.FirstName == "" {
		return errBadRequest.WithMessage("FirstName cannot be empty.")
	}
	if req.LastName == "" {
		return errBadRequest.WithMessage("LastName cannot be empty.")
	}
	if req.AddressString == "" {
		return errBadRequest.WithMessage("AddressString cannot be empty.")
	}
	if req.Zip == "" {
		return errBadRequest.WithMessage("Zip cannot be empty.")
	}
	if req.PaymentMethodToken == "" {
		return errBadRequest.WithMessage("PaymentMethodToken cannot be empty.")
	}
	if req.CustomerID == "" {
		return errBadRequest.WithMessage("CustomerID cannot be empty.")
	}
	if req.Amount < 0.001 {
		return errBadRequest.WithMessage("Amount cannot be empty.")
	}
	geopoint := common.GeoPoint{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	if !geopoint.Valid() {
		return errBadRequest.WithMessage("Geopoint is invalid.")
	}
	if req.ServingsNonVegetarian == 0 && req.ServingsVegetarain == 0 {
		return errBadRequest.WithMessage("Servings cannot be empty.")
	}
	return nil
}

// Create creates an activity entry.
func (c *Client) Create(req *CreateReq) error {
	err := req.validate()
	if err != nil {
		return err
	}
	_, err = c.sqlDB.NamedExecContext(c.ctx, insertStatement, req)
	if err != nil {
		if merr, ok := err.(*mysql.MySQLError); ok {
			if merr.Number == 1062 {
				return errDuplicateEntry.WithError(err).Wrap("activity already exists")
			}
		}
		return errSQLDB.WithError(err).Annotate("failed to insertStatement")
	}
	return nil
}

// DeleteFuture deletes activities for future subscriber that are unpaid.
func (c *Client) DeleteFuture(date time.Time, idOrEmail string) error {
	if idOrEmail == "" {
		return errBadRequest.Annotate("invalid user_id or email")
	}
	// Update actvity
	var err error
	if strings.Contains(idOrEmail, "@") {
		_, err = c.sqlDB.ExecContext(c.ctx, deleteFutureEmailStatment, date.Format(DateFormat), idOrEmail)
	} else {
		_, err = c.sqlDB.ExecContext(c.ctx, deleteFutureStatment, date.Format(DateFormat), idOrEmail)
	}
	if err != nil {
		return errSQLDB.WithError(err).Annotate("failed to execute deleteFutureStatment")
	}
	return nil
}

// Process processes an actvity.
// func (c *Client) Process(date time.Time, email string) error {
// 	// TODO: Reimplement
// 	suboldC := subold.NewWithLogging(c.ctx, c.log)
// 	return suboldC.Process(date, email)
// }

// Discount gives an discount to an actvitiy.
func (c *Client) Discount(date time.Time, email string, discountAmount float32, discountPercent int8, overrideDiscounts bool) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Discount(date, email, discountAmount, discountPercent, overrideDiscounts)
}

// Refund refunds and skips a subscriber.
func (c *Client) Refund(date time.Time, email string, amount float32, precent int32) error {
	if amount > 0 && precent > 0 {
		return errBadRequest.WithMessage("Only amount or percent can be specified.")
	}
	// Get activity
	act, err := c.Get(date, email)
	if err != nil {
		return errors.Wrap("failed to Get", err)
	}
	if act.Refunded {
		return errBadRequest.WithMessage("Activity is already refunded.")
	}
	if precent > 0 {
		amount = act.Amount * float32(precent) / 100
	}
	// Refund card
	paymentC, err := payment.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to payment.New")
	}
	rID, err := paymentC.RefundSale(act.TransactionID, amount)
	if err != nil {
		return errors.Wrap("failed to payment.RefundSale", err)
	}
	// log
	c.log.Refund(act.UserID, act.Email, act.Date, act.Amount, amount, rID)
	c.log.Infof(c.ctx, "Refunding Customer(%s) on transaction(%s): refundID(%s)", act.CustomerID, act.TransactionID, rID)
	// Update actvity
	_, err = c.sqlDB.ExecContext(c.ctx, updateRefundedStatement, rID, amount, date.Format(DateFormat), email)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateRefundedStatement")
	}
	return nil
}

// RefundAndSkip refunds and skips a subscriber.
func (c *Client) RefundAndSkip(date time.Time, email string, amount float32, precent int32) error {
	// Refund
	err := c.Refund(date, email, amount, precent)
	if err != nil {
		return errors.Annotate(err, "failed to Refund")
	}
	// Skip
	err = c.Skip(date, email, "Refunded")
	if err != nil {
		return errors.Annotate(err, "failed to Skip")
	}
	return nil
}

// Skip skips a subscriber for an activity.
func (c *Client) Skip(date time.Time, email, reason string) error {
	// TODO: Reimplement

	// TODO: Requires discount to be in a seperate table
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Skip(date, email, reason)
}

// Unskip unskips a subscriber for an activity.
func (c *Client) Unskip(date time.Time, email string) error {
	// TODO: Reimplement

	// TODO: Requires discount to be in a seperate table
	// _, err = c.sqlDB.ExecContext(c.ctx, unskipStatement, date.Format(DateFormat), email)
	// if err != nil {
	// 	return errSQLDB.WithError(err).Wrap("failed to execute unskipStatement")
	// }

	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Unskip(date, email)
}

// Paid sets activity to paid.
func (c *Client) Paid(date time.Time, userID string, amount float32, tID string) error {
	if userID == "" || date.IsZero() {
		return errBadRequest.Annotate("invalid UserID or Date")
	}
	act, err := c.Get(date, userID)
	if err != nil {
		return errBadRequest.WithError(err).Annotate("activity not found")
	}
	if act.Paid {
		return errBadRequest.WithMessage("Failed. User has already paid.")
	}
	_, err = c.sqlDB.ExecContext(c.ctx, updatePaidStatement, amount, time.Now().Format(datetimeFormat), tID, date.Format(DateFormat), userID)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updatePaidStatement")
	}
	c.log.Paid(userID, act.Email, date.Format(time.RFC3339), act.Amount, amount, tID)
	return nil
}

// ChangeServings for a day
func (c *Client) ChangeServings(date time.Time, email string, servings int8, amount float32) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.ChangeServings(date, email, servings, amount)
}

// Free f
func (c *Client) Free(date time.Time, email string) error {
	// TODO: Reimplement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Free(date, email)
}

func (c *Client) BatchUpdateActivityWithUserID(userIDs []string, emails []string) error {
	if len(userIDs) != len(emails) {
		return errBadRequest.WithMessage("UserIDs and Emails should be same length")
	}
	var err error
	for i := range userIDs {
		statement := "UPDATE activity set user_id=? where email=?"
		_, err = c.sqlDB.Exec(statement, userIDs[i], emails[i])
		if err != nil {
			return errSQLDB.WithError(err).Annotate("failed to run statment")
		}
	}
	return nil
}

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
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
)

const (
	datetimeFormat = "2006-01-02 15:04:05" // "Jan 2, 2006 at 3:04pm (MST)"
	// DateFormat is the expected format of date in activity.
	DateFormat = "2006-01-02" // "Jan 2, 2006"
	// Select
	selectActivityByEmailStatement           = "SELECT * FROM activity WHERE date=? AND email=?"
	selectActivityByUserIDStatement          = "SELECT * FROM activity WHERE date=? AND user_id=?"
	selectAllActivityStatement               = "SELECT * FROM activity ORDER BY date DESC LIMIT ?"
	selectActivityForUserStatement           = "SELECT * FROM activity WHERE user_id=? ORDER BY date DESC"
	selectActivityForUserStatementEmail      = "SELECT * FROM activity WHERE email=? ORDER BY date DESC"
	selectActivityAfterDateForUserStatement  = "SELECT * FROM activity WHERE user_id=? AND date>=? ORDER BY date DESC"
	selectActivityBeforeDateForUserStatement = "SELECT * FROM activity WHERE user_id=? AND date<=? ORDER BY date DESC"
	selectUnpaidSummaries                    = "SELECT min(date) as mn,max(date) as mx,user_id,email,first_name,last_name,sum(amount) as amount_due,count(user_id) as num_unpaid FROM activity WHERE date<NOW() AND discount_percent<>100 AND paid=0 AND skip=0 AND refunded=0 AND forgiven=0 GROUP BY user_id ORDER BY mx"
	selectUnpaidSummaryForUser               = "SELECT min(date) as mn,max(date) as mx,user_id,email,first_name,last_name,sum(amount) as amount_due,count(user_id) as num_unpaid FROM activity WHERE date<NOW() AND discount_percent<>100 AND paid=0 AND skip=0 AND refunded=0 AND forgiven=0 AND user_id=? GROUP BY user_id ORDER BY mx"
	// selectUnpaidActivities                   = "SELECT * FROM activity WHERE date<NOW() AND discount_percent<>100 AND paid=0 AND skip=0 AND refunded=0 AND forgiven=0 ORDER BY user_id,date"
	selectFirstActivity = "SELECT * FROM activity WHERE first=1 AND skip=0 AND user_id=?"
	// update and insert
	updateServingsStatement       = "UPDATE activity SET servings=?,veg_servings=?,amount=?,servings_changed=1 WHERE date=? AND user_id=?"
	updateFutureServingsStatement = "UPDATE activity SET servings=?,veg_servings=?,amount=? WHERE date>? AND user_id=? AND paid=0 AND servings_changed=0"

	insertStatement         = "INSERT INTO activity (date,user_id,email,first_name,last_name,location,addr_apt,addr_string,zip,lat,`long`,active,skip,servings,veg_servings,first,amount,discount_amount,discount_percent,payment_provider,payment_method_token,customer_id) VALUES (:date,:user_id,:email,:first_name,:last_name,:location,:addr_apt,:addr_string,:zip,:lat,:long,:active,:skip,:servings,:veg_servings,:first,:amount,:discount_amount,:discount_percent,:payment_provider,:payment_method_token,:customer_id)"
	updateRefundedStatement = "UPDATE activity SET refunded_dt=NOW(),refunded=1,refund_transaction_id=?,refunded_amount=? WHERE date=? AND user_id=?"
	skipStatement           = "UPDATE activity SET skip=1 WHERE date=? AND user_id=?"
	unskipStatement         = "UPDATE activity SET skip=0,active=1 WHERE date=? AND user_id=?"
	updateForgiveStatement  = "UPDATE activity SET forgiven=1 WHERE date=? AND user_id=?"
	updatePaidStatement     = "UPDATE activity SET amount_paid=?,discount_amount=?,discount_percent=?,paid=1,paid_dt=?,transaction_id=? WHERE date=? AND user_id=?"

	// delete
	deleteFutureStatment = "DELETE from activity WHERE date>? AND user_id=? AND paid=0"
	// TODO: switch to user id
	deleteFutureEmailStatment     = "DELETE from activity WHERE date>? AND email=? AND paid=0"
	deleteFutureUnskippedStatment = "DELETE from activity WHERE date>? AND user_id=? AND paid=0 AND skip=0"
)

// Errors
var (
	errInternal       = errors.InternalServerError
	errBadRequest     = errors.BadRequestError
	errSQLDB          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	errDuplicateEntry = errors.ErrorWithCode{Code: errors.CodeBadRequest, Message: "Duplicate entry."}
	errNotFound       = errors.NotFoundError
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
		if strings.Contains(err.Error(), "no rows in result") {
			return nil, errNotFound.WithError(err).Annotate("failed to selectActivity")
		}
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
	var err error
	if strings.Contains(userID, "@") {
		err = c.sqlDB.SelectContext(c.ctx, &acts, selectActivityForUserStatementEmail, userID)
	} else {
		err = c.sqlDB.SelectContext(c.ctx, &acts, selectActivityForUserStatement, userID)
	}
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

// GetUnpaidSummaries gets a list of activity.
func (c *Client) GetUnpaidSummaries() ([]*UnpaidSummary, error) {
	dist := []*UnpaidSummary{}
	err := c.sqlDB.SelectContext(c.ctx, &dist, selectUnpaidSummaries)
	if err != nil {
		return nil, errSQLDB.WithError(err).Annotate("failed to selectUnpaidSummaries")
	}
	return dist, nil
}

// GetUnpaidSummary gets a list of activity.
func (c *Client) GetUnpaidSummary(userID string) (*UnpaidSummary, error) {
	dist := []*UnpaidSummary{}
	err := c.sqlDB.SelectContext(c.ctx, &dist, selectUnpaidSummaryForUser, userID)
	if err != nil {
		return nil, errSQLDB.WithError(err).Annotate("failed to selectUnpaidSummaryForUser")
	}

	if len(dist) > 1 {
		return nil, errBadRequest.Annotate("dist len does not equal 1")
	}
	if len(dist) == 1 {
		return dist[0], nil
	}
	// no unpaid
	return &UnpaidSummary{}, nil
}

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
	// set if first meal or not
	act := &Activity{}
	err = c.sqlDB.GetContext(c.ctx, act, selectFirstActivity, req.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			req.First = true
		} else {
			return errSQLDB.WithError(err).Annotate("failed to selectFirstActivity")
		}
	}
	// insert
	_, err = c.sqlDB.NamedExecContext(c.ctx, insertStatement, req)
	if err != nil {
		if merr, ok := err.(*mysql.MySQLError); ok {
			if merr.Number == 1062 {
				return errDuplicateEntry.WithError(err).Wrap("activity already exists")
			}
		}
		return errSQLDB.WithError(err).Annotate("failed to insertStatement")
	}

	// Add to process activity
	d, _ := time.Parse(DateFormat, req.Date)
	dayBeforeActivity := d.Add(-24 * time.Hour)
	taskC := tasks.New(c.ctx)
	r := &tasks.ProcessSubscriptionParams{
		UserID:   req.UserID,
		SubEmail: req.Email,
		Date:     d,
	}
	err = taskC.AddProcessSubscription(dayBeforeActivity, r)
	if err != nil {
		return errors.Wrap("failed to tasks.AddProcessSubscription", err)
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

// DeleteFutureUnskipped deletes activities for future subscriber that are unpaid and upskipped.
func (c *Client) DeleteFutureUnskipped(date *time.Time, id string) error {
	if id == "" {
		return errBadRequest.Annotate("invalid user_id or email")
	}
	// Update actvity
	var err error
	_, err = c.sqlDB.ExecContext(c.ctx, deleteFutureUnskippedStatment, date.Format(DateFormat), id)
	if err != nil {
		return errSQLDB.WithError(err).Annotate("failed to execute deleteFutureUnskippedStatment")
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
func (c *Client) Refund(date time.Time, idOrEmail string, amount float32, precent int32) error {
	if amount > 0 && precent > 0 {
		return errBadRequest.WithMessage("Only amount or percent can be specified.")
	}
	// Get activity
	act, err := c.Get(date, idOrEmail)
	if err != nil {
		return errors.Wrap("failed to Get", err)
	}
	if !act.Paid {
		return errBadRequest.WithMessage("Activity has not paid.")
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
	_, err = c.sqlDB.ExecContext(c.ctx, updateRefundedStatement, rID, amount, date.Format(DateFormat), act.UserID)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateRefundedStatement")
	}
	return nil
}

// RefundAndSkip refunds and skips a subscriber.
func (c *Client) RefundAndSkip(date time.Time, idOrEmail string, amount float32, precent int32) error {
	// Refund
	err := c.Refund(date, idOrEmail, amount, precent)
	if err != nil {
		return errors.Annotate(err, "failed to Refund")
	}
	// Skip
	err = c.Skip(date, idOrEmail, "Refunded")
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

// Forgive forgives a subscriber for an activity.
func (c *Client) Forgive(date time.Time, userID string) error {
	if userID == "" || date.IsZero() {
		return errBadRequest.Annotate("invalid UserID or Date")
	}
	act, err := c.Get(date, userID)
	if err != nil {
		return errBadRequest.WithError(err).Annotate("activity not found")
	}
	_, err = c.sqlDB.ExecContext(c.ctx, updateForgiveStatement, date.Format(DateFormat), userID)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateForgiveStatement")
	}
	c.log.Forgiven(userID, act.Email, date.Format(time.RFC3339), act.Amount, act.AmountPaid-act.Amount)
	return nil
}

// Paid sets activity to paid.
func (c *Client) Paid(date time.Time, userID string, amount, discountAmount float32, discountPercent int8, tID string) error {
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
	_, err = c.sqlDB.ExecContext(c.ctx, updatePaidStatement, amount, discountAmount, discountPercent, time.Now().Format(datetimeFormat), tID, date.Format(DateFormat), userID)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updatePaidStatement")
	}
	c.log.Paid(userID, act.Email, date.Format(time.RFC3339), act.Amount, amount, tID)
	return nil
}

// ChangeServings for an activity.
func (c *Client) ChangeServings(date time.Time, id string, servingsNonVeg, servingsVeg int8, amount float32) error {
	if servingsNonVeg < 0 {
		return errBadRequest.WithMessage("Servings non-veg cannot be less than zero.")
	}
	if servingsVeg < 0 {
		return errBadRequest.WithMessage("Servings veg cannot be less than zero.")
	}
	if servingsNonVeg < 0 && servingsVeg < 0 {
		return errBadRequest.WithMessage("Servings non-veg and servings both cannot be less than zero.")
	}
	if amount < 0.01 {
		return errBadRequest.WithMessage("Amount cannot be less than 0.")
	}
	act, err := c.Get(date, id)
	if err != nil {
		return errors.Annotate(err, "failed to Get")
	}
	if act.Paid {
		return errBadRequest.WithMessage("Activity is already paid.")
	}
	if act.Skip {
		return errBadRequest.WithMessage("Activity is already skipped.")
	}
	_, err = c.sqlDB.ExecContext(c.ctx, updateServingsStatement, servingsNonVeg, servingsVeg, amount, date.Format(DateFormat), id)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateServingsStatement")
	}
	c.log.SubServingsChanged(id, act.Email, date.Format(DateFormat), act.ServingsNonVegetarian, servingsNonVeg, act.ServingsVegetarain, servingsVeg)
	return nil
}

// ChangeFutureServings for an activity.
func (c *Client) ChangeFutureServings(date time.Time, id string, servingsNonVeg, servingsVeg int8, amount float32) error {
	if servingsNonVeg < 0 {
		return errBadRequest.WithMessage("Servings non-veg cannot be less than zero.")
	}
	if servingsVeg < 0 {
		return errBadRequest.WithMessage("Servings veg cannot be less than zero.")
	}
	if servingsNonVeg < 0 && servingsVeg < 0 {
		return errBadRequest.WithMessage("Servings non-veg and servings both cannot be less than zero.")
	}
	if amount < 0.01 {
		return errBadRequest.WithMessage("Amount cannot be less than 0.")
	}
	_, err := c.sqlDB.ExecContext(c.ctx, updateFutureServingsStatement, servingsNonVeg, servingsVeg, amount, date.Format(DateFormat), id)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateFutureServingsStatement")
	}
	return nil
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

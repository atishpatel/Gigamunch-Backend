package sub

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	// driver for mysql
	"cloud.google.com/go/datastore"
	mysql "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"google.golang.org/appengine"
)

const (
	datetimeFormat                           = "2006-01-02 15:04:05" // "Jan 2, 2006 at 3:04pm (MST)"
	dateFormat                               = "2006-01-02"          // "Jan 2, 2006"
	insertSubLogStatement                    = "INSERT INTO `sub` (date,sub_email,servings,amount,delivery_time,payment_method_token,customer_id) VALUES ('%s','%s',%d,%f,%d,'%s','%s')"
	selectSubLogEmails                       = "SELECT DISTINCT sub_email from sub where date>? and date<?"
	selectSubLogStatement                    = "SELECT created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent,refunded FROM `sub` WHERE date='%s' AND sub_email='%s'"
	selectSubscriberSubLogsStatement         = "SELECT date,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent, refunded FROM `sub` WHERE sub_email=? ORDER BY date DESC"
	selectAllSubLogStatement                 = "SELECT date,sub_email,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent,refunded FROM `sub` ORDER BY date DESC LIMIT %d"
	selectUnpaidSubLogStatement              = "SELECT date,sub_email,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent,refunded FROM `sub` WHERE paid=0 AND free=0 AND skip=0 AND refunded=0 ORDER BY date DESC LIMIT %d"
	selectSubscriberUnpaidSubLogStatement    = "SELECT date,sub_email,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent,refunded FROM `sub` WHERE paid=0 AND free=0 AND skip=0 AND refunded=0 AND sub_email=?"
	selectSubLogFromDateStatement            = "SELECT date,sub_email,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent,refunded FROM `sub` WHERE date=?"
	selectSubLogFromSubscriberStatement      = "SELECT date,sub_email,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent,refunded FROM `sub` WHERE sub_email=? ORDER BY date DESC"
	selectSublogSummaryStatement             = "SELECT min(date) as mn,max(date),sub_email,count(sub_email),sum(skip),sum(paid),sum(refunded),sum(amount),sum(amount_paid),sum(discount_amount) FROM sub WHERE date>'2017-04-08' AND date<? GROUP BY sub_email ORDER BY mn"
	updatePaidSubLogStatement                = "UPDATE `sub` SET amount_paid=%f,paid=1,paid_datetime='%s',transaction_id='%s' WHERE date='%s' AND sub_email='%s'"
	updateSkipSubLogStatement                = "UPDATE `sub` SET skip=1 WHERE date='%s' AND sub_email='%s'"
	updateUnskipSubLogStatement              = "UPDATE `sub` SET skip=0 WHERE date='%s' AND sub_email='%s'"
	updateRefundedAndSkipSubLogStatement     = "UPDATE `sub` SET skip=1,refunded=1 WHERE date=? AND sub_email=?"
	updateFreeSubLogStatment                 = "UPDATE `sub` SET free=1 WHERE date='%s' AND sub_email='%s'"
	updateDiscountSubLogStatment             = "UPDATE `sub` SET discount_amount=?, discount_percent=? WHERE date=? AND sub_email=?"
	updateServingsSubLogStatement            = "UPDATE sub SET servings=?, amount=? WHERE date=? AND sub_email=?"
	updateServingsPermanentlySubLogStatement = "UPDATE sub SET servings=?, amount=? WHERE date>? AND sub_email=? AND servings=?"
	deleteSubLogStatment                     = "DELETE from `sub` WHERE date>? AND sub_email=? AND paid=0"
	updateUnpaidPayment                      = "UPDATE sub SET payment_method_token=? WHERE free=0 AND paid=0 AND skip=0 AND sub_email=?"
	updateEmailAddress                       = "UPDATE `sub` SET sub_email='%s' WHERE sub_email='%s'"
	// insertPromoCodeStatement     = "INSERT INTO `promo_code` (code,free_delivery,percent_off,amount_off,discount_cap,free_dish,buy_one_get_one_free,start_datetime,end_datetime,num_uses) VALUES ('%s',%t,%d,%f,%f,%t,%t,'%s','%s',%d)"
	// selectPromoCodesStatement    = "SELECT created_datetime,free_delivery,percent_off,amount_off,discount_cap,free_dish,buy_one_get_one_free,start_datetime,end_datetime,num_uses FROM `promo_code` WHERE code='%s'"
)

var (
	connectOnce = sync.Once{}
	mysqlDB     *sql.DB
	errSQLDB    = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	// errBuffer           = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "An unknown error occured."}
	errDatastore         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errDatastoreNotFound = errors.ErrorWithCode{Code: errors.CodeNotFound, Message: "Not found."}
	errInvalidParameter  = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errEntrySkipped      = errors.ErrorWithCode{Code: 401, Message: "Invalid parameter. Entry is skipped."}
	errNoSuchEntry       = errors.ErrorWithCode{Code: 4001, Message: "Invalid parameter."}
	errDuplicateEntry    = errors.ErrorWithCode{Code: 4000, Message: "Invalid parameter."}
	projID               string
)

// Client is the client fro this package.
type Client struct {
	ctx context.Context
	log *logging.Client
}

// New returns a new Client.
func New(ctx context.Context) *Client {
	connectOnce.Do(func() {
		connectSQL(ctx)
	})
	return &Client{ctx: ctx}
}

// NewWithLogging returns a new Client.
func NewWithLogging(ctx context.Context, log *logging.Client) *Client {
	connectOnce.Do(func() {
		connectSQL(ctx)
	})
	return &Client{
		ctx: ctx,
		log: log,
	}
}

// GetSubEmails gets a list of unique subscriber emails within the date range.
func (c *Client) GetSubEmails(from, to time.Time) ([]string, error) {
	rows, err := mysqlDB.Query(selectSubLogEmails, from.Format(dateFormat), to.Format(dateFormat))
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to run GetSubEmails")
	}
	var emails []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		emails = append(emails, email)
	}
	return emails, nil
}

// GetSubscriber returns a SubscriptionSignUp.
func (c *Client) GetSubscriber(email string) (*SubscriptionSignUp, error) {
	if email == "" {
		return nil, errInvalidParameter.Wrap("emails cannot be empty.")
	}
	subs, err := c.GetSubscribers([]string{email})
	if err == datastore.ErrNoSuchEntity {
		return nil, errDatastoreNotFound
	}
	if err != nil || len(subs) != 1 {
		return nil, errors.Wrap("failed to c.GetSubscribers", err)
	}
	return subs[0], nil
}

// GetSubscribers returns a list of SubscriptionSignUp.
func (c *Client) GetSubscribers(emails []string) ([]*SubscriptionSignUp, error) {
	if len(emails) == 0 {
		return nil, errInvalidParameter.Wrap("emails cannot be of length 0.")
	}
	for i := range emails {
		emails[i] = strings.TrimSpace(emails[i])
	}
	subs, err := getMulti(c.ctx, emails)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to getMulti")
	}
	return subs, nil
}

// GetSubscribersByPhoneNumber returns a SubscriptionSignUp.
func (c *Client) GetSubscribersByPhoneNumber(number string) ([]*SubscriptionSignUp, error) {
	if number == "" {
		return nil, errInvalidParameter.Wrap("number cannot be empty.")
	}
	cleanNumber := getCleanPhoneNumber(number)
	subs, err := getSubscribersByPhoneNumber(c.ctx, cleanNumber)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to getHasSubscribed")
	}
	return subs, nil
}

// GetHasSubscribed returns a list of all SubscriptionSignUp.
func (c *Client) GetHasSubscribed(date time.Time) ([]SubscriptionSignUp, error) {
	subs, err := getHasSubscribed(c.ctx, date)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to getHasSubscribed")
	}
	return subs, nil
}

// GetSublogSummaries gets a summary of SubLogs.
func (c *Client) GetSublogSummaries() ([]*SublogSummary, error) {
	rows, err := mysqlDB.Query(selectSublogSummaryStatement, time.Now().Format(dateFormat))
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query selectSublogSummaryStatement statement:")
	}
	defer handleCloser(c.ctx, rows)
	var subs []*SublogSummary
	for rows.Next() {
		sub := new(SublogSummary)
		var minDate mysql.NullTime
		var maxDate mysql.NullTime
		err = rows.Scan(&minDate, &maxDate, &sub.Email, &sub.NumTotal, &sub.NumSkip, &sub.NumPaid, &sub.NumRefunded, &sub.TotalAmount, &sub.TotalAmountPaid, &sub.TotalDiscountAmount)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if minDate.Valid {
			sub.MinDate = minDate.Time
		}
		if maxDate.Valid {
			sub.MaxDate = maxDate.Time
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// GetAll gets all the SubLogs.
func (c *Client) GetAll(limit int32) ([]*SubscriptionLog, error) {
	if limit <= 0 {
		limit = 1000
	}
	st := fmt.Sprintf(selectAllSubLogStatement, limit)
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query selectAllSubLogStatement statement:")
	}
	defer handleCloser(c.ctx, rows)
	var subLogs []*SubscriptionLog
	for rows.Next() {
		subLog := new(SubscriptionLog)
		var date mysql.NullTime
		var createdNulltime mysql.NullTime
		var paidNulltime mysql.NullTime
		err = rows.Scan(&date, &subLog.SubEmail, &createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent, &subLog.Refunded)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if date.Valid {
			subLog.Date = date.Time
		}
		if createdNulltime.Valid {
			subLog.CreatedDatetime = createdNulltime.Time
		}
		if paidNulltime.Valid {
			subLog.PaidDatetime = paidNulltime.Time
		}
		subLogs = append(subLogs, subLog)
	}
	return subLogs, nil
}

// GetUnpaidSublogs gets unpaid SubLogs.
func (c *Client) GetUnpaidSublogs(limit int32) ([]*SubscriptionLog, error) {
	if limit <= 0 {
		limit = 1000
	}
	st := fmt.Sprintf(selectUnpaidSubLogStatement, limit)
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query selectUnpaidSubLogStatement statement:")
	}
	defer handleCloser(c.ctx, rows)
	var subLogs []*SubscriptionLog
	for rows.Next() {
		subLog := new(SubscriptionLog)
		var date mysql.NullTime
		var createdNulltime mysql.NullTime
		var paidNulltime mysql.NullTime
		err = rows.Scan(&date, &subLog.SubEmail, &createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent, &subLog.Refunded)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if date.Valid {
			subLog.Date = date.Time
		}
		if createdNulltime.Valid {
			subLog.CreatedDatetime = createdNulltime.Time
		}
		if paidNulltime.Valid {
			subLog.PaidDatetime = paidNulltime.Time
		}
		subLogs = append(subLogs, subLog)
	}
	return subLogs, nil
}

// GetSubscriberUnpaidSublogs gets unpaid SubLogs.
func (c *Client) GetSubscriberUnpaidSublogs(email string) ([]*SubscriptionLog, error) {
	rows, err := mysqlDB.Query(selectSubscriberUnpaidSubLogStatement, email)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query selectSubscriberUnpaidSubLogStatement statement:")
	}
	defer handleCloser(c.ctx, rows)
	var subLogs []*SubscriptionLog
	for rows.Next() {
		subLog := new(SubscriptionLog)
		var date mysql.NullTime
		var createdNulltime mysql.NullTime
		var paidNulltime mysql.NullTime
		err = rows.Scan(&date, &subLog.SubEmail, &createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent, &subLog.Refunded)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if date.Valid {
			subLog.Date = date.Time
		}
		if createdNulltime.Valid {
			subLog.CreatedDatetime = createdNulltime.Time
		}
		if paidNulltime.Valid {
			subLog.PaidDatetime = paidNulltime.Time
		}
		subLogs = append(subLogs, subLog)
	}
	return subLogs, nil
}

// GetForDate gets all the SubLogs.
func (c *Client) GetForDate(date time.Time) ([]*SubscriptionLog, error) {
	rows, err := mysqlDB.Query(selectSubLogFromDateStatement, date.Format(dateFormat))
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query statement:" + selectSubLogFromDateStatement)
	}
	defer handleCloser(c.ctx, rows)
	var subLogs []*SubscriptionLog
	for rows.Next() {
		subLog := new(SubscriptionLog)
		var date mysql.NullTime
		var createdNulltime mysql.NullTime
		var paidNulltime mysql.NullTime
		err = rows.Scan(&date, &subLog.SubEmail, &createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent, &subLog.Refunded)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if date.Valid {
			subLog.Date = date.Time
		}
		if createdNulltime.Valid {
			subLog.CreatedDatetime = createdNulltime.Time
		}
		if paidNulltime.Valid {
			subLog.PaidDatetime = paidNulltime.Time
		}
		subLogs = append(subLogs, subLog)
	}
	return subLogs, nil
}

// GetSubscriberSublogs gets all the SubLogs of a subscriber.
func (c *Client) GetSubscriberSublogs(email string) ([]*SubscriptionLog, error) {
	rows, err := mysqlDB.Query(selectSubLogFromSubscriberStatement, email)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query statement:" + selectSubLogFromSubscriberStatement)
	}
	defer handleCloser(c.ctx, rows)
	var subLogs []*SubscriptionLog
	for rows.Next() {
		subLog := new(SubscriptionLog)
		var date mysql.NullTime
		var createdNulltime mysql.NullTime
		var paidNulltime mysql.NullTime
		err = rows.Scan(&date, &subLog.SubEmail, &createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent, &subLog.Refunded)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if date.Valid {
			subLog.Date = date.Time
		}
		if createdNulltime.Valid {
			subLog.CreatedDatetime = createdNulltime.Time
		}
		if paidNulltime.Valid {
			subLog.PaidDatetime = paidNulltime.Time
		}
		subLogs = append(subLogs, subLog)
	}
	return subLogs, nil
}

// Get gets a SubLog.
func (c *Client) Get(date time.Time, subEmail string) (*SubscriptionLog, error) {
	if date.IsZero() || subEmail == "" {
		return nil, errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) ", date, subEmail)
	}
	st := fmt.Sprintf(selectSubLogStatement, date.Format(dateFormat), subEmail)
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query selectSubLogStatement statement.")
	}
	defer handleCloser(c.ctx, rows)
	if !rows.Next() {
		return nil, errNoSuchEntry.Wrap("no such entry found")
	}
	subLog := &SubscriptionLog{
		Date:     date,
		SubEmail: subEmail,
	}
	var createdNulltime mysql.NullTime
	var paidNulltime mysql.NullTime
	err = rows.Scan(&createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent, &subLog.Refunded)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
	}
	if createdNulltime.Valid {
		subLog.CreatedDatetime = createdNulltime.Time
	}
	if paidNulltime.Valid {
		subLog.PaidDatetime = paidNulltime.Time
	}
	return subLog, nil
}

// GetSubscriberActivities gets a subscriber.
func (c *Client) GetSubscriberActivities(email string) ([]*SubscriptionLog, error) {
	if email == "" {
		return nil, errInvalidParameter.Wrapf("expected(actual): subEmail(%s) ", email)
	}
	rows, err := mysqlDB.Query(selectSubscriberSubLogsStatement, email)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query selectSubscriberSubLogsStatement statement.")
	}
	defer handleCloser(c.ctx, rows)
	var subLogs []*SubscriptionLog
	for rows.Next() {
		subLog := new(SubscriptionLog)
		subLog.SubEmail = email
		var date mysql.NullTime
		var createdNulltime mysql.NullTime
		var paidNulltime mysql.NullTime
		err = rows.Scan(&date, &createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent, &subLog.Refunded)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if date.Valid {
			subLog.Date = date.Time
		}
		if createdNulltime.Valid {
			subLog.CreatedDatetime = createdNulltime.Time
		}
		if paidNulltime.Valid {
			subLog.PaidDatetime = paidNulltime.Time
		}
		subLogs = append(subLogs, subLog)
	}
	return subLogs, nil
}

// Update updates all the subs info.
func (c *Client) Update(subs []*SubscriptionSignUp) error {
	if len(subs) == 0 {
		return nil
	}
	emails := make([]string, len(subs))
	for i, sub := range subs {
		emails[i] = sub.Email
	}
	oldSubs, err := getMulti(c.ctx, emails)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to getMulti")
	}
	err = putMulti(c.ctx, subs)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to putMulti")
	}
	logging.Infof(c.ctx, "after putMulti(%d) %s", len(subs), err)
	if c.log != nil {
		for i := range subs {
			payload := &logging.SubUpdatedPayload{
				OldEmail:          oldSubs[i].Email,
				Email:             subs[i].Email,
				OldFirstName:      oldSubs[i].FirstName,
				FirstName:         subs[i].FirstName,
				OldLastName:       oldSubs[i].LastName,
				LastName:          subs[i].LastName,
				OldAddress:        oldSubs[i].Address.String(),
				Address:           subs[i].Address.String(),
				OldRawPhoneNumber: oldSubs[i].RawPhoneNumber,
				RawPhoneNumber:    subs[i].RawPhoneNumber,
				OldPhoneNumber:    oldSubs[i].PhoneNumber,
				PhoneNumber:       subs[i].PhoneNumber,
				OldDeliveryNotes:  oldSubs[i].DeliveryTips,
				DeliveryNotes:     subs[i].DeliveryTips,
			}
			c.log.SubUpdated(0, subs[i].Email, payload)
		}
	}
	return nil
}

// Setup sets up a SubLog.
func (c *Client) Setup(date time.Time, subEmail string, servings int8, amount float32, deliveryTime int8, paymentMethodToken, customerID string) error {
	if date.IsZero() || subEmail == "" || servings == 0 || amount == 0 || paymentMethodToken == "" || customerID == "" {
		return errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) servings(%d) amount(%f) deliveryTime(%d) paymentMethodToken(%s) customerID(%s)", date, subEmail, servings, amount, deliveryTime, paymentMethodToken, customerID)
	}
	st := fmt.Sprintf(insertSubLogStatement, date.Format(dateFormat), subEmail, servings, amount, deliveryTime, paymentMethodToken, customerID)
	_, err := mysqlDB.Exec(st)
	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1062 {
			return errDuplicateEntry.WithError(err).Wrap("failed to execute insertSubLogStatement statement.")
		}
		return errSQLDB.WithError(err).Wrap("failed to execute insertSubLogStatement statement.")
	}
	return nil
}

func (c *Client) UpdatePaymentToken(subEmail string, paymentMethodToken string) error {
	if subEmail == "" || paymentMethodToken == "" {
		return errInvalidParameter.Wrapf("expected(actual): subEmail(%s) paymentMethodToken(%s)", subEmail, paymentMethodToken)
	}
	s, err := c.GetSubscriber(subEmail)
	if err != nil {
		return errors.Wrap("failed to sub.GetSubscriber", err)
	}
	oldPMT := s.PaymentMethodToken
	s.PaymentMethodToken = paymentMethodToken
	utils.Infof(c.ctx, "changing sub(%s)'s payment method token from old(%s) new(%s)", subEmail, oldPMT, paymentMethodToken)
	err = put(c.ctx, subEmail, s)
	if err != nil {
		return errors.Wrap("failed to put", err)
	}
	_, err = mysqlDB.Exec(updateUnpaidPayment, paymentMethodToken, subEmail)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateUnpaidPayment statement")
	}
	if c.log != nil {
		c.log.SubCardUpdated(0, subEmail, oldPMT, paymentMethodToken)
	}
	return nil
}

// ChangeServings inserts or updates a SubLog with a different serving amount.
func (c *Client) ChangeServings(date time.Time, subEmail string, servings int8, amount float32) error {
	// insert or update
	if date.IsZero() || subEmail == "" || servings < 1 || amount < 0.1 {
		return errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) servings(%f) amount(%s)", date, subEmail, servings, amount)
	}
	sl, err := c.Get(date, subEmail)
	if err != nil {
		if errors.GetErrorWithCode(err).Code != errNoSuchEntry.Code {
			return errors.Wrap("failed to sub.Get", err)
		}
		// insert
		var s *SubscriptionSignUp
		s, err = get(c.ctx, subEmail)
		if err != nil {
			return errDatastore.WithError(err).Wrap("failed to get")
		}
		serv := s.Servings + s.VegetarianServings
		err = c.Setup(date, subEmail, serv, s.WeeklyAmount, s.DeliveryTime, s.PaymentMethodToken, s.CustomerID)
		if err != nil {
			return errors.Wrap("failed to sub.Setup", err)
		}
	} else {
		if sl.Paid {
			return errEntrySkipped.Wrap("cannot give change servings to a week that is already paid")
		}
		if sl.Skip {
			return errEntrySkipped.Wrap("cannot give change servings to a week that is already skipped")
		}
	}
	_, err = mysqlDB.Exec(updateServingsSubLogStatement, servings, amount, date.Format(dateFormat), subEmail)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateServingsSubLogStatement statement")
	}
	return nil
}

// ChangeServingsPermanently changes a subscriber's servings permanently for all bags from now onwards.
func (c *Client) ChangeServingsPermanently(subEmail string, servings int8, vegetarian bool, serverInfo *common.ServerInfo) error {
	// insert or update
	if subEmail == "" || servings < 1 {
		return errInvalidParameter.Wrapf("expected(actual): subEmail(%s) servings(%f)", subEmail, servings)
	}
	s, err := c.GetSubscriber(subEmail)
	if err != nil {
		return errors.Wrap("failed to sub.GetSubscriber", err)
	}
	oldWeeklyAmount := s.WeeklyAmount
	oldServings := s.Servings
	oldVegServings := s.VegetarianServings
	weeklyAmount := DerivePrice(servings)
	var vegServings int8
	var nonvegServings int8
	if vegetarian {
		vegServings = servings
	} else {
		nonvegServings = servings
	}
	utils.Infof(c.ctx, "changing sub(%s)'s servings from nonveg(%d) veg(%d) amount(%2f) to nonveg(%d) veg(%d) amount(%2f)", subEmail, oldServings, oldVegServings, oldWeeklyAmount, nonvegServings, vegServings, weeklyAmount)
	s.WeeklyAmount = weeklyAmount
	s.Servings = nonvegServings
	s.VegetarianServings = vegServings
	err = put(c.ctx, subEmail, s)
	if err != nil {
		return errors.Wrap("failed to put", err)
	}
	// TODO don't update if past deadline date for Serving count
	_, err = mysqlDB.Exec(updateServingsPermanentlySubLogStatement, servings, s.WeeklyAmount, time.Now().Format(dateFormat), subEmail, oldServings+oldVegServings)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateServingsPermanentlySubLogStatement statement")
	}
	if c.log != nil {
		c.log.SubServingsChangedPermanently(0, subEmail, oldServings, nonvegServings, oldVegServings, vegServings)
	}

	mailC, err := mail.NewClient(c.ctx, c.log, serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to mail.NewClient")
	}
	mailReq := &mail.UserFields{
		Email:          subEmail,
		VegServings:    s.VegetarianServings,
		NonVegServings: s.Servings,
	}
	err = mailC.UpdateUser(mailReq)
	if err != nil {
		return errors.Annotate(err, "failed to mail.UpdateUser")
	}

	return nil
}

// Paid inserts or updates a SubLog to paid.
func (c *Client) Paid(date time.Time, subEmail string, amountPaid float32, transactionID string) error {
	// insert or update
	if date.IsZero() || subEmail == "" {
		return errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) amountPaid(%f) transactionID(%s)", date, subEmail, amountPaid, transactionID)
	}
	oldEntry, err := c.Get(date, subEmail)
	if err != nil {
		if errors.GetErrorWithCode(err).Code != errNoSuchEntry.Code {
			return errors.Wrap("failed to sub.Get", err)
		}
		// insert
		var s *SubscriptionSignUp
		s, err = get(c.ctx, subEmail)
		if err != nil {
			return errDatastore.WithError(err).Wrap("failed to get")
		}
		servings := s.Servings + s.VegetarianServings
		err = c.Setup(date, subEmail, servings, s.WeeklyAmount, s.DeliveryTime, s.PaymentMethodToken, s.CustomerID)
		if err != nil {
			return errors.Wrap("failed to sub.Setup", err)
		}
	}
	st := fmt.Sprintf(updatePaidSubLogStatement, amountPaid, time.Now().Format(datetimeFormat), transactionID, date.Format(dateFormat), subEmail)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updatePaidSubLogStatement statement.")
	}
	if c.log != nil {
		c.log.Paid(0, subEmail, date.Format(time.RFC3339), oldEntry.Amount, amountPaid, transactionID)
	}
	return nil
}

// Skip skips that subscription for that day.
func (c *Client) Skip(date time.Time, subEmail, reason string) error {
	s, err := c.GetSubscriber(subEmail)
	if err != nil {
		return errors.Wrap("failed to sub.GetSubscriber", err)
	}
	// insert or update
	sl, err := c.Get(date, subEmail)
	if err != nil {
		if errors.GetErrorWithCode(err).Code != errNoSuchEntry.Code {
			return errors.Wrap("failed to sub.Get", err)
		}
		// insert
		// var s *SubscriptionSignUp
		// s, err = get(c.ctx, subEmail)
		// if err != nil {
		// 	return errDatastore.WithError(err).Wrap("failed to get")
		// }
		servings := s.Servings + s.VegetarianServings
		err = c.Setup(date, subEmail, servings, s.WeeklyAmount, s.DeliveryTime, s.PaymentMethodToken, s.CustomerID)
		if err != nil {
			return errors.Wrap("failed to sub.Setup", err)
		}
	} else {
		if sl.Paid {
			return errInvalidParameter.WithMessage("Subscriber has already paid. Must refund instead.")
		}
		if s.IsSubscribed {
			// if there is a discount, move it to next week unskipped week.
			if sl.DiscountAmount > .01 || sl.DiscountPercent != 0 {
				nextWeek := date.Add(7 * 24 * time.Hour)
				for {
					err = c.Discount(nextWeek, subEmail, sl.DiscountAmount, sl.DiscountPercent, false)
					if err == nil {
						break
					}
					if errors.GetErrorWithCode(err).Code != errEntrySkipped.Code {
						return errors.Wrap("failed to Discount", err)
					}
					nextWeek = nextWeek.Add(7 * 24 * time.Hour)
				}
			}
		}
		// if first
		if sl.Free {
			err = c.Free(date.Add(7*24*time.Hour), subEmail)
			if err != nil {
				return errors.Wrap("failed to Free", err)
			}
		}
	}
	st := fmt.Sprintf(updateSkipSubLogStatement, date.Format(dateFormat), subEmail)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateSkipSubLogStatement statement.")
	}
	if c.log != nil {
		utils.Infof(c.ctx, "log not nil. logging skip")
		c.log.SubSkip(date.Format(time.RFC3339), 0, subEmail, reason)
	} else {
		utils.Infof(c.ctx, "log nil")
	}
	return nil
}

// Unskip Unskips that subscription for that day.
func (c *Client) Unskip(date time.Time, subEmail string) error {
	s, err := c.GetSubscriber(subEmail)
	if err != nil {
		return errors.Wrap("failed to sub.GetSubscriber", err)
	}
	// insert or update
	_, err = c.Get(date, subEmail)
	if err != nil {
		if errors.GetErrorWithCode(err).Code != errNoSuchEntry.Code {
			return errors.Wrap("failed to sub.Get", err)
		}
		// insert
		// var s *SubscriptionSignUp
		// s, err = get(c.ctx, subEmail)
		// if err != nil {
		// 	return errDatastore.WithError(err).Wrap("failed to get")
		// }
		servings := s.Servings + s.VegetarianServings
		err = c.Setup(date, subEmail, servings, s.WeeklyAmount, s.DeliveryTime, s.PaymentMethodToken, s.CustomerID)
		if err != nil {
			return errors.Wrap("failed to sub.Setup", err)
		}
	}
	st := fmt.Sprintf(updateUnskipSubLogStatement, date.Format(dateFormat), subEmail)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateUnskipSubLogStatement statement.")
	}
	if c.log != nil {
		utils.Infof(c.ctx, "log not nil. logging skip")
		c.log.SubUnskip(date.Format(time.RFC3339), 0, subEmail)
	} else {
		utils.Infof(c.ctx, "log nil")
	}
	return nil
}

// RefundAndSkip refunds and skips that subscription for that day.
func (c *Client) RefundAndSkip(date time.Time, subEmail string) error {
	// insert or update
	sl, err := c.Get(date, subEmail)
	if err != nil {
		return errors.Wrap("failed to sub.Get", err)
	}
	if !sl.Paid {
		return errInvalidParameter.WithMessage("Subscriber has not paid. Use skip instead.")
	}
	paymentC := payment.New(c.ctx)
	rID, err := paymentC.RefundSale(sl.TransactionID)
	if err != nil {
		return errors.Wrap("failed to payment.RefundSale", err)
	}
	utils.Infof(c.ctx, "Refunding Customer(%s) on transaction(%s): refundID(%s)", sl.CustomerID, sl.TransactionID, rID)
	r, err := mysqlDB.Exec(updateRefundedAndSkipSubLogStatement, date.Format(dateFormat), subEmail)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateRefundedAndSkipSubLogStatement")
	}
	numEffectedRows, err := r.RowsAffected()
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed get RowsAffected")
	}
	if numEffectedRows != 1 {
		return errSQLDB.WithError(err).Wrapf("num effected rows is not 1: %s", numEffectedRows)
	}
	return nil
}

// Discount inserts or updates a SubLog with a discount.
func (c *Client) Discount(date time.Time, subEmail string, discountAmount float32, discountPercent int8, overrideDiscounts bool) error {
	// insert or update
	if date.IsZero() || subEmail == "" || discountAmount < 0 || discountPercent < 0 || discountPercent > 100 {
		return errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) discountAmount(%f) discountPercent(%s)", date, subEmail, discountAmount, discountPercent)
	}
	s, err := c.GetSubscriber(subEmail)
	if err != nil {
		return errors.Wrap("failed to sub.GetSubscriber", err)
	}
	if !s.IsSubscribed {
		return errInvalidParameter.WithMessage(subEmail + " is no longer a subscriber. :(")
	}
	sl, err := c.Get(date, subEmail)
	if err != nil {
		if errors.GetErrorWithCode(err).Code != errNoSuchEntry.Code {
			return errors.Wrap("failed to sub.Get", err)
		}
		// insert
		var s *SubscriptionSignUp
		s, err = get(c.ctx, subEmail)
		if err != nil {
			return errDatastore.WithError(err).Wrap("failed to get")
		}
		servings := s.Servings + s.VegetarianServings
		err = c.Setup(date, subEmail, servings, s.WeeklyAmount, s.DeliveryTime, s.PaymentMethodToken, s.CustomerID)
		if err != nil {
			return errors.Wrap("failed to sub.Setup", err)
		}
	} else {
		if sl.Paid {
			return errEntrySkipped.Wrap("cannot give discount to a week that is already paid")
		}
		if sl.Skip {
			return errEntrySkipped.Wrap("cannot give discount to a week that is already skipped")
		}
		if sl.Free {
			return errEntrySkipped.Wrap("cannot give discount to a week that is already free (first meal)")
		}
		if !overrideDiscounts && (sl.DiscountAmount > .1 || sl.DiscountPercent != 0) {
			return errInvalidParameter.WithMessage("Cannot give discount because entry already has a discount!")
		}
	}
	_, err = mysqlDB.Exec(updateDiscountSubLogStatment, discountAmount, discountPercent, date.Format(dateFormat), subEmail)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateDiscountSubLogStatment statement")
	}
	return nil
}

// Free inserts or updates sub for that day to free.
func (c *Client) Free(date time.Time, subEmail string) error {
	// insert or update
	_, err := c.Get(date, subEmail)
	if err != nil {
		if errors.GetErrorWithCode(err).Code != errNoSuchEntry.Code {
			return errors.Wrap("failed to sub.Get", err)
		}
		// insert
		var s *SubscriptionSignUp
		s, err = get(c.ctx, subEmail)
		if err != nil {
			return errDatastore.WithError(err).Wrap("failed to get")
		}
		servings := s.Servings + s.VegetarianServings
		err = c.Setup(date, subEmail, servings, s.WeeklyAmount, s.DeliveryTime, s.PaymentMethodToken, s.CustomerID)
		if err != nil {
			return errors.Wrap("failed to sub.Setup", err)
		}
	}
	st := fmt.Sprintf(updateFreeSubLogStatment, date.Format(dateFormat), subEmail)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateFreeSubLogStatment statement: ")
	}
	return nil
}

// Cancel cancels an user's subscription.
func (c *Client) Cancel(subEmail string, log *logging.Client, serverInfo *common.ServerInfo) error {
	if subEmail == "" {
		return errInvalidParameter.Wrap("sub email cannot be empty.")
	}
	// remove any SubLog that are > now
	_, err := mysqlDB.Exec(deleteSubLogStatment, time.Now().Format(dateFormat), subEmail)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("failed to execute statement: %s", deleteSubLogStatment)
	}
	// change isSubscribed to false
	sub, err := get(c.ctx, subEmail)
	if err != nil {
		return errors.Wrap("failed to get sub", err)
	}
	if !sub.IsSubscribed {
		return errInvalidParameter.Wrapf("%s is already not subscribed.", subEmail)
	}
	sub.IsSubscribed = false
	sub.UnSubscribedDate = time.Now()
	err = put(c.ctx, subEmail, sub)
	if err != nil {
		return errors.Wrap("failed to put sub", err)
	}
	mailC, err := mail.NewClient(c.ctx, log, serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to mail.NewClient")
	}
	mailReq := &mail.UserFields{
		Email:     subEmail,
		FirstName: sub.FirstName,
		LastName:  sub.LastName,
	}
	err = mailC.SubDeactivated(mailReq)
	if err != nil {
		return errors.Annotate(err, "failed to mail.SubDeactivated")
	}
	return nil
}

// Process process a SubLog.
func (c *Client) Process(date time.Time, subEmail string) error {
	utils.Infof(c.ctx, "Processing Sub: date(%v) subEmail(%s)", date, subEmail)
	subLog, err := c.Get(date, subEmail)
	if err != nil {
		errCode := errors.GetErrorWithCode(err)
		if errCode.Code == errNoSuchEntry.Code {
			utils.Infof(c.ctx, "failed to sub.Get because user canceled: %+v", err)
			return nil
		}
		return errors.Wrap("failed to sub.Get", err)
	}
	// done if Skipped
	if subLog.Skip {
		utils.Infof(c.ctx, "Subscription is already finished. Skip(%v)", subLog.Skip)
		return nil
	}
	dayBeforeBox := subLog.Date.Add(-24 * time.Hour)
	if time.Now().Before(dayBeforeBox) {
		// too early to process
		r := &tasks.ProcessSubscriptionParams{
			SubEmail: subLog.SubEmail,
			Date:     subLog.Date,
		}
		taskC := tasks.New(c.ctx)
		err = taskC.AddProcessSubscription(dayBeforeBox, r)
		if err != nil {
			// TODO critical?
			return errors.Wrap("failed to tasks.AddProcessSubscription", err)
		}
		utils.Infof(c.ctx, "Too early to process Sub. now(%v) < dayBeforeBox(%v)", time.Now(), dayBeforeBox)
		return nil
	}
	taskC := tasks.New(c.ctx)
	r := &tasks.UpdateDripParams{
		Email: subLog.SubEmail,
	}
	err = taskC.AddUpdateDrip(dayBeforeBox, r)
	if err != nil {
		utils.Criticalf(c.ctx, "failed to tasks.AddUpdateDrip: %+v", err)
	}
	// done if Free, Paid
	if subLog.Free || subLog.Paid {
		utils.Infof(c.ctx, "Subscription is already finished. Free(%v) Paid(%v)", subLog.Free, subLog.Paid)
		return nil
	}
	// charge customer
	amount := subLog.Amount
	amount -= subLog.DiscountAmount
	amount -= (float32(subLog.DiscountPercent) / 100) * amount
	orderID := fmt.Sprintf("Gigamunch box for %s.", date.Format("01/02/2006"))
	var tID string
	if amount > 0.0 {
		paymentC := payment.New(c.ctx)
		saleReq := &payment.SaleReq{
			CustomerID:         subLog.CustomerID,
			Amount:             amount,
			PaymentMethodToken: subLog.PaymentMethodToken,
			OrderID:            orderID,
		}
		utils.Infof(c.ctx, "Charging Customer(%s) %f on card(%s)", subLog.CustomerID, amount, subLog.PaymentMethodToken)
		tID, err = paymentC.Sale(saleReq)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				// Dulicate transaction error because two customers have same card
				r := &tasks.ProcessSubscriptionParams{
					SubEmail: subLog.SubEmail,
					Date:     subLog.Date,
				}
				taskC := tasks.New(c.ctx)
				err = taskC.AddProcessSubscription(time.Now().Add(1*time.Hour), r)
				if err != nil {
					// TODO critical?
					return errors.Wrap("failed to tasks.AddProcessSubscription", err)
				}
				return nil
			}
			return errors.Wrap("failed to payment.Sale", err)
		}
	}
	// update TransactionID
	err = c.Paid(subLog.Date, subLog.SubEmail, amount, tID)
	if err != nil {
		// TODO
		return errors.Wrap("failed to sub.Paid", err)
	}
	return nil
}

// SetupSubLogs gets all the subscribers who are subscribed and adds them to the SubLog for the specified date.
func (c *Client) SetupSubLogs(date time.Time) error {
	// get all SubSignups
	dayName := date.Weekday().String()
	subs, err := getSubscribers(c.ctx, dayName)
	if err != nil {
		return errDatastore.WithError(err).Wrap("failed to getSubscribers")
	}
	utils.Infof(c.ctx, "adding %d subscribers to SubLog", len(subs))
	taskC := tasks.New(c.ctx)
	dayBeforeBox := date.Add(-2 * time.Hour) // TODO: change cron to timezone to make code easier to understand
	for _, v := range subs {
		if (!v.FirstBoxDate.IsZero() && v.FirstBoxDate.After(dayBeforeBox)) || (!v.SubscriptionDate.IsZero() && v.SubscriptionDate.After(dayBeforeBox)) {
			continue
		}
		// TODO instead of inserting all in this task, split it into many tasks?
		// insert into subLog
		amt := v.WeeklyAmount
		servings := v.Servings + v.VegetarianServings
		if amt < .01 { // TODO remove and just give error?
			switch servings {
			case 1:
				amt = 17
			case 2:
				amt = 15 * 2
			default:
				amt = 14 * float32(servings)
			}
		}
		err = c.Setup(date, v.Email, servings, amt, v.DeliveryTime, v.PaymentMethodToken, v.CustomerID)
		if err != nil {
			if errors.GetErrorWithCode(err).Code == errDuplicateEntry.Code {
				continue
			}
			return errors.Wrap("failed to sub.Setup", err)
		}
		// add to task queue
		r := &tasks.ProcessSubscriptionParams{
			SubEmail: v.Email,
			Date:     date,
		}
		err = taskC.AddProcessSubscription(dayBeforeBox, r)
		if err != nil {
			return errors.Wrap("failed to tasks.AddProcessSubscription", err)
		}
	}
	return nil
}

func connectSQL(ctx context.Context) {
	var err error
	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = "root@/gigamunch"
	} else {
		connectionString = os.Getenv("MYSQL_CONNECTION")
	}
	mysqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Couldn't connect to mysql database")
	}
}

type closer interface {
	Close() error
}

func handleCloser(ctx context.Context, c closer) {
	err := c.Close()
	if err != nil {
		utils.Errorf(ctx, "Error closing rows: %v", err)
	}
}

// DerivePrice returns the price for a set number of servings.
func DerivePrice(servings int8) float32 {
	switch servings {
	case 1:
		return 17 + 1.66
	case 2:
		return (16 * 2) + 3.12
	case 4:
		return (15 * 4) + 5.85
	default:
		return 15 * float32(servings) * 1.0975
	}
}

// UpdateEmail changes all instances of a customer email to new email.
func (c *Client) UpdateEmail(oldEmail, newEmail string) error {
	utils.Infof(c.ctx, "changed email from %s to %s in sql db", oldEmail, newEmail)
	st := fmt.Sprintf(updateEmailAddress, newEmail, oldEmail)
	_, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute updateEmailAddress statement")
	}
	return nil
}

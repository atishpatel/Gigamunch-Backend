package sub

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	// driver for mysql
	mysql "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"google.golang.org/appengine"
)

const (
	datetimeFormat                = "2006-01-02 15:04:05" // "Jan 2, 2006 at 3:04pm (MST)"
	dateFormat                    = "2006-01-02"          // "Jan 2, 2006"
	insertSubLogStatement         = "INSERT INTO `sub` (date,sub_email,servings,amount,delivery_time,payment_method_token,customer_id) VALUES ('%s','%s',%d,%f,%d,'%s','%s')"
	selectSubLogStatement         = "SELECT created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent FROM `sub` WHERE date='%s' AND sub_email='%s'"
	selectAllSubLogStatement      = "SELECT date,sub_email,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent FROM `sub` ORDER BY date DESC LIMIT %d"
	selectSubLogFromDateStatement = "SELECT date,sub_email,created_datetime,skip,servings,amount,amount_paid,paid,paid_datetime,delivery_time,payment_method_token,transaction_id,free,discount_amount,discount_percent,refunded FROM `sub` WHERE date=?"
	updatePaidSubLogStatement     = "UPDATE `sub` SET amount_paid=%f,paid=1,paid_datetime='%s',transaction_id='%s' WHERE date='%s' AND sub_email='%s'"
	updateSkipSubLogStatement     = "UPDATE `sub` SET skip=1 WHERE date='%s' AND sub_email='%s'"
	updateFreeSubLogStatment      = "UPDATE `sub` SET free=1 WHERE date='%s' AND sub_email='%s'"
	deleteSubLogStatment          = "DELETE from `sub` WHERE date>? AND sub_email=? AND paid=0"
	// insertPromoCodeStatement     = "INSERT INTO `promo_code` (code,free_delivery,percent_off,amount_off,discount_cap,free_dish,buy_one_get_one_free,start_datetime,end_datetime,num_uses) VALUES ('%s',%t,%d,%f,%f,%t,%t,'%s','%s',%d)"
	// selectPromoCodesStatement    = "SELECT created_datetime,free_delivery,percent_off,amount_off,discount_cap,free_dish,buy_one_get_one_free,start_datetime,end_datetime,num_uses FROM `promo_code` WHERE code='%s'"
)

var (
	connectOnce = sync.Once{}
	mysqlDB     *sql.DB
	errSQLDB    = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	// errBuffer           = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "An unknown error occured."}
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errNoSuchEntry      = errors.ErrorWithCode{Code: 4001, Message: "Invalid parameter."}
	errDuplicateEntry   = errors.ErrorWithCode{Code: 4000, Message: "Invalid parameter."}
)

// Client is the client fro this package.
type Client struct {
	ctx context.Context
}

// New returns a new Client.
func New(ctx context.Context) *Client {
	connectOnce.Do(func() {
		connectSQL(ctx)
	})
	return &Client{ctx: ctx}
}

// GetSubscribers returns a list of SubscriptionSignUp.
func (c *Client) GetSubscribers(emails []string) ([]*SubscriptionSignUp, error) {
	if len(emails) == 0 {
		return nil, errInvalidParameter.Wrap("emails cannot be of length 0.")
	}
	subs, err := getMulti(c.ctx, emails)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to getMulti")
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
		return nil, errSQLDB.WithError(err).Wrap("failed to query statement:" + st)
	}
	defer handleCloser(c.ctx, rows)
	var subLogs []*SubscriptionLog
	for rows.Next() {
		subLog := new(SubscriptionLog)
		var date mysql.NullTime
		var createdNulltime mysql.NullTime
		var paidNulltime mysql.NullTime
		err = rows.Scan(&date, &subLog.SubEmail, &createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent)
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

// Get gets a SubLog.
func (c *Client) Get(date time.Time, subEmail string) (*SubscriptionLog, error) {
	if date.IsZero() || subEmail == "" {
		return nil, errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) ", date, subEmail)
	}
	st := fmt.Sprintf(selectSubLogStatement, date.Format(dateFormat), subEmail)
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query statement:" + st)
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
	err = rows.Scan(&createdNulltime, &subLog.Skip, &subLog.Servings, &subLog.Amount, &subLog.AmountPaid, &subLog.Paid, &paidNulltime, &subLog.DeliveryTime, &subLog.PaymentMethodToken, &subLog.TransactionID, &subLog.Free, &subLog.DiscountAmount, &subLog.DiscountPercent)
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

// Setup sets up a SubLog.
func (c *Client) Setup(date time.Time, subEmail string, servings int8, amount float32, deliveryTime int8, paymentMethodToken, customerID string) error {
	if date.IsZero() || subEmail == "" || servings == 0 || amount == 0 || deliveryTime == 0 || paymentMethodToken == "" || customerID == "" {
		return errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) servings(%d) amount(%f) deliveryTime(%d) paymentMethodToken(%s) customerID(%s)", date, subEmail, servings, amount, deliveryTime, paymentMethodToken, customerID)
	}
	st := fmt.Sprintf(insertSubLogStatement, date.Format(dateFormat), subEmail, servings, amount, deliveryTime, paymentMethodToken, customerID)
	_, err := mysqlDB.Exec(st)
	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1062 {
			return errDuplicateEntry.WithError(err).Wrap("failed to execute statement: " + st)
		}
		return errSQLDB.WithError(err).Wrap("failed to execute statement: " + st)
	}
	return nil
}

// Paid inserts or updates a SubLog to paid.
func (c *Client) Paid(date time.Time, subEmail string, amountPaid float32, transactionID string) error {
	// insert or update
	if date.IsZero() || subEmail == "" || amountPaid < .01 || transactionID == "" {
		return errInvalidParameter.Wrapf("expected(actual): date(%v) subEmail(%s) amountPaid(%f) transactionID(%s)", date, subEmail, amountPaid, transactionID)
	}
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
	st := fmt.Sprintf(updatePaidSubLogStatement, amountPaid, time.Now().Format(datetimeFormat), transactionID, date.Format(dateFormat), subEmail)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute statement: " + st)
	}
	// TODO add insert if not in table
	return nil
}

// Skip skips that subscription for that day.
func (c *Client) Skip(date time.Time, subEmail string) error {
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
	st := fmt.Sprintf(updateSkipSubLogStatement, date.Format(dateFormat), subEmail)
	_, err = mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute statement: " + st)
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
		return errSQLDB.WithError(err).Wrap("failed to execute statement: " + st)
	}
	return nil
}

// Cancel cancels an user's subscription.
func (c *Client) Cancel(subEmail string) error {
	if subEmail == "" {
		return errInvalidParameter.Wrap("sub email cannot be empty.")
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
	// remove any SubLog that are > now
	_, err = mysqlDB.Exec(deleteSubLogStatment, time.Now().Format(dateFormat), subEmail)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("failed to execute statement: %s", deleteSubLogStatment)
	}
	return nil
}

// Process process a SubLog.
func (c *Client) Process(date time.Time, subEmail string) error {
	utils.Infof(c.ctx, "Processing Sub: date(%v) subEmail(%s)", date, subEmail)
	subLog, err := c.Get(date, subEmail)
	if err != nil {
		return errors.Wrap("failed to sub.Get", err)
	}
	if subLog.Free || subLog.Skip || subLog.Paid {
		// done
		utils.Infof(c.ctx, "Subscription is already finished. Free(%v) Skip(%v) Paid(%v)", subLog.Free, subLog.Skip, subLog.Paid)
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
	// charge customer
	amount := subLog.Amount
	amount -= subLog.DiscountAmount
	amount -= (float32(subLog.DiscountPercent) / 100) * amount
	orderID := fmt.Sprintf("Gigamunch box for %s.", date.Format("01/02/2006"))
	paymentC := payment.New(c.ctx)
	saleReq := &payment.SaleReq{
		CustomerID:         subLog.CustomerID,
		Amount:             amount,
		PaymentMethodToken: subLog.PaymentMethodToken,
		OrderID:            orderID,
	}
	utils.Infof(c.ctx, "Charging Customer(%s) %f on card(%s)", subLog.CustomerID, amount, subLog.PaymentMethodToken)
	tID, err := paymentC.Sale(saleReq)
	if err != nil {
		// TODO
		return errors.Wrap("failed to payment.Sale", err)
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
	dayBeforeBox := date.Add(-2 * time.Hour)
	for _, v := range subs {
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

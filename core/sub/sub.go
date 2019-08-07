package sub

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/discount"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/core/slack"
	paymentold "github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"

	"github.com/jmoiron/sqlx"
)

var (
	errDatastore             = errors.InternalServerError
	errNoSuchEntityDatastore = errors.NotFoundError
	errInternal              = errors.InternalServerError
	errInvalidParameter      = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	db         common.DB
	sqlDB      *sqlx.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	if sqlC == nil {
		return nil, errInternal.Annotate("failed to get sql client")
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

// Get gets a subscriber.
func (c *Client) Get(id string) (*subold.Subscriber, error) {
	return c.getByIDOrEmail(id)
}

// GetMulti gets a subscriber.
func (c *Client) GetMulti(ids []string) ([]*subold.Subscriber, error) {
	subs, err := c.getMulti(ids)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to getMulti")
	}
	return subs, nil
}

// GetByEmail gets a subscriber by email.
func (c *Client) GetByEmail(email string) (*subold.Subscriber, error) {
	sub, err := c.getByEmail(email)
	if err != nil {
		if err == c.db.ErrNoSuchEntity() {
			return nil, errNoSuchEntityDatastore.WithError(err).Annotate("failed to getByEmail")
		}
		return nil, errDatastore.WithError(err).Annotate("failed to getByEmail")
	}
	return sub, nil
}

// GetActive gets all active subscribers.
func (c *Client) GetActive(start, limit int) ([]*subold.Subscriber, error) {
	var subs []*subold.Subscriber
	_, err := c.db.QueryFilter(c.ctx, kind, start, limit, "Active=", true, &subs)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryFilter")
	}
	return subs, nil
}

// GetHasSubscribed returns a list of all Subscribers.
func (c *Client) GetHasSubscribed(start, limit int) ([]*subold.Subscriber, error) {
	subs, err := c.getHasSubscribed(start, limit)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to getHasSubscribed")
	}
	return subs, nil
}

// GetByEmail gets a subscriber by email.

// GetByPhoneNumber gets a subscriber by phone number.
// func (c *Client) GetByPhoneNumber() error {
// 	// TODO: implement
// 	return nil
// }

// GetHasSubscribed returns a list of all subscribers ever.
// func (c *Client) GetHasSubscribed() error {
// 	// TODO: implement
// 	return nil
// }

// ChangeServingsPermanently changes a subscriber's servings permanently.
func (c *Client) ChangeServingsPermanently(id string, servingsNonVeg, servingsVeg int8) error {
	if servingsNonVeg < 0 {
		return errInvalidParameter.WithMessage("Servings non-veg cannot be less than zero.")
	}
	if servingsVeg < 0 {
		return errInvalidParameter.WithMessage("Servings veg cannot be less than zero.")
	}
	if servingsNonVeg < 0 && servingsVeg < 0 {
		return errInvalidParameter.WithMessage("Servings non-veg and servings both cannot be less than zero.")
	}
	sub, err := c.getByIDOrEmail(id)
	if err != nil {
		return errors.Annotate(err, "failed to Get")
	}
	oldWeeklyAmount := sub.Amount
	oldServingsNonVeg := sub.ServingsNonVegetarian
	oldServingsVeg := sub.ServingsVegetarian

	if (servingsNonVeg + servingsVeg) != (sub.ServingsNonVegetarian + sub.ServingsVegetarian) {
		sub.Amount = DerivePrice(servingsNonVeg + servingsVeg)
	}
	sub.ServingsNonVegetarian = servingsNonVeg
	sub.ServingsVegetarian = servingsVeg
	err = c.put(sub.ID, sub)
	if err != nil {
		return errors.Annotate(err, "failed to put")
	}
	// update activities
	activityC, err := activity.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.ChangeFutureServings(time.Now(), sub.ID, sub.ServingsNonVegetarian, sub.ServingsVegetarian, sub.Amount)
	if err != nil {
		return errors.Annotate(err, "failed to activity.ChangeFutureServings")
	}
	// log
	c.log.SubServingsChangedPermanently(sub.ID, sub.Email(), oldServingsNonVeg, sub.ServingsNonVegetarian, oldServingsVeg, sub.ServingsVegetarian, oldWeeklyAmount, sub.Amount)
	// mail
	mailC, err := mail.NewClient(c.ctx, c.log, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to mail.NewClient")
	}
	mailReq := &mail.UserFields{
		Email:          sub.Email(),
		VegServings:    sub.ServingsVegetarian,
		NonVegServings: sub.ServingsNonVegetarian,
	}
	err = mailC.UpdateUser(mailReq)
	if err != nil {
		return errors.Annotate(err, "failed to mail.UpdateUser")
	}
	return nil
}

// ChangePlanDay changes a subscriber's plan day.
func (c *Client) ChangePlanDay(id string, planDay string, intervalStartPoint *time.Time) error {
	if planDay != time.Monday.String() && planDay != time.Thursday.String() {
		return errInvalidParameter.WithMessage("Invalid plan day.")
	}
	if intervalStartPoint == nil || intervalStartPoint.IsZero() {
		return errInvalidParameter.WithMessage("Invalid interval start date.")
	}
	sub, err := c.getByIDOrEmail(id)
	if err != nil {
		return errors.Annotate(err, "failed to Get")
	}
	oldSub := *sub
	sub.PlanWeekday = planDay
	sub.IntervalStartPoint = *intervalStartPoint
	for sub.IntervalStartPoint.Weekday().String() != sub.PlanWeekday {
		sub.IntervalStartPoint = sub.IntervalStartPoint.Add(12 * time.Hour)
	}
	sub.IntervalStartPoint = sub.IntervalStartPoint.Add(12 * time.Hour) // set to midday
	err = c.put(sub.ID, sub)
	if err != nil {
		return errors.Annotate(err, "failed to put")
	}
	// update activities
	activityC, err := activity.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.DeleteFutureUnskipped(intervalStartPoint, sub.ID)
	if err != nil {
		return errors.Annotate(err, "failed to activity.DeleteFutureUnskipped")
	}
	// log
	c.log.SubUpdated(sub.ID, sub.Email(), &oldSub, sub)
	// mail
	taskC := tasks.New(c.ctx)
	err = taskC.AddUpdateDrip(time.Now(), &tasks.UpdateDripParams{
		UserID: sub.ID,
	})
	if err != nil {
		return errors.Annotate(err, "failed to tasks.AddUpdateDrip")
	}
	// setup next activity
	act, err := activityC.Get(sub.IntervalStartPoint, sub.ID)
	if err != nil && errors.GetErrorWithCode(err).Code != errors.CodeNotFound {
		return errors.Annotate(err, "failed to activity.Get")
	}
	if (act == nil || act.Date == "") && sub.IntervalStartPoint.After(time.Now()) {
		// no activity so set up an activity
		err = c.SetupActivity(sub.IntervalStartPoint, sub.ID, true, 0, 0)
		if err != nil {
			return errors.Annotate(err, "failed to sub.SetupActivity")
		}
	}
	return nil
}

// UpdatePaymentToken updates a user payment method token.
func (c *Client) UpdatePaymentToken(email, paymentMethodToken string) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.UpdatePaymentToken(email, paymentMethodToken)
}

// Update updates a subscriber.
func (c *Client) Update(sub *subold.Subscriber) error {
	if sub.ID == "" {
		return errInvalidParameter.WithMessage("sub doesn't have an id")
	}
	subold, err := c.Get(sub.ID)
	if err != nil {
		return errors.Annotate(err, "failed to Get")
	}
	key := c.db.NameKey(c.ctx, kind, sub.ID)
	_, err = c.db.Put(c.ctx, key, sub)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to put")
	}
	c.log.SubUpdated(sub.ID, sub.Email(), subold, sub)
	return nil
}

// Activate activates an account.
func (c *Client) Activate(email string, firstBagDate time.Time) error {
	// TODO: implement
	// TODO: check if has payment method and such
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Activate(email, firstBagDate, c.log, c.serverInfo)
}

// Deactivate deactivates an account
func (c *Client) Deactivate(idOrEmail, reason string) error {

	if idOrEmail == "" {
		return errInvalidParameter.Annotate("id or email is empty")
	}
	// remove any Activity that are greater than now
	activityC, err := activity.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	// TODO: should time be now or now + couple of days?
	err = activityC.DeleteFuture(time.Now(), idOrEmail)
	if err != nil {
		return errors.Annotate(err, "failed to activity.DeleteFuture")
	}
	// change isSubscribed to false
	sub, err := c.getByIDOrEmail(idOrEmail)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to getByIDOrEmail")
	}
	if !sub.Active {
		return errInvalidParameter.WithMessage("User is already deactivated.").Annotate("not active")
	}
	sub.Active = false
	sub.DeactivatedDatetime = time.Now()
	err = c.Update(sub)
	if err != nil {
		return errors.Annotate(err, "failed to Update")
	}
	// check of unpaid activities
	outstandingCharges := false
	unpaidSummary, err := activityC.GetUnpaidSummary(sub.ID)
	if err != nil {
		c.log.Errorf(c.ctx, "failed to GetUnpaidSummary: %+v", err)
	} else {
		if unpaidSummary != nil && unpaidSummary.NumUnpaid != "0" && unpaidSummary.NumUnpaid != "" {
			outstandingCharges = true
		}
	}
	c.log.Infof(c.ctx, "Has outstanding Charges? %v", outstandingCharges)
	// update mail client
	mailC, err := mail.NewClient(c.ctx, c.log, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to mail.NewClient")
	}
	for _, emailPref := range sub.EmailPrefs {
		mailReq := &mail.UserFields{
			Email:     emailPref.Email,
			FirstName: emailPref.FirstName,
			LastName:  emailPref.LastName,
		}
		err = mailC.SubDeactivated(mailReq, outstandingCharges)
		if err != nil {
			return errors.Annotate(err, "failed to mail.SubDeactivated")
		}
	}
	daysActive := int(sub.DeactivatedDatetime.Sub(sub.ActivateDatetime) / (time.Hour * 24))
	// logging
	c.log.SubDeactivated(sub.ID, sub.Email(), reason, daysActive)
	// send to slack
	slackC, err := slack.NewClient(c.ctx, c.log, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to slack.NewClient")
	}
	if !strings.Contains(sub.Email(), "@test.com") {
		err = slackC.SendDeactivate(sub.Email(), sub.FullName(), reason, daysActive)
		if err != nil {
			c.log.Errorf(c.ctx, "failed to slack.SendDeactivate: %+v", err)
		}
	}
	return nil
}

// CreateReq is the request for Create.
type CreateReq struct {
	Email                 string            `json:"email"`
	FirstName             string            `json:"first_name"`
	LastName              string            `json:"last_name"`
	PhoneNumber           string            `json:"phone_number"`
	Address               common.Address    `json:"address"`
	DeliveryNotes         string            `json:"delivery_notes"`
	Reference             string            `json:"reference"`
	ReferenceEmail        string            `json:"reference_email"`
	PaymentMethodNonce    string            `json:"payment_method_nonce"`
	ServingsNonVegetarian int8              `json:"servings_non_vegetarian"`
	ServingsVegetarian    int8              `json:"servings_vegetarian"`
	FirstDeliveryDate     time.Time         `json:"first_delivery_date"`
	Campaigns             []common.Campaign `json:"campaigns"`
	DiscountAmount        float32           `json:"discount_amount"`
	DiscountPercent       int8              `json:"discount_precent"`
}

func (req *CreateReq) fix() {
	req.Email = strings.Replace(strings.ToLower(req.Email), " ", "", -1)
	req.PhoneNumber = strings.Replace(req.PhoneNumber, " ", "", -1)
	req.Address.Street = strings.Title(req.Address.Street)
	req.Address.City = strings.Title(req.Address.City)
	req.Address.Zip = strings.TrimSpace(req.Address.Zip)
}

func (req *CreateReq) validateBasic() error {
	if req.Email == "" {
		return errInvalidParameter.WithMessage("Email address cannot be empty.")
	}
	if !strings.Contains(req.Email, "@") {
		return errInvalidParameter.WithMessage("Invalid email address.")
	}
	if req.PaymentMethodNonce == "" {
		return errInvalidParameter.WithMessage("Invalid payment info.").Annotate("no payment nonce")
	}
	if req.FirstName == "" {
		return errInvalidParameter.WithMessage("First name must be provided.")
	}
	return nil
}

func (req *CreateReq) validateAll() error {
	if req.ServingsNonVegetarian+req.ServingsVegetarian < 2 {
		return errInvalidParameter.WithMessage("Incorrect servings amount")
	}
	if req.FirstDeliveryDate.Before(time.Now()) {
		return errInvalidParameter.WithMessage("Invalid first delivery date")
	}
	return nil
}

// Create creates an new subscriber. Error is returned if was a subscriber before.
func (c *Client) Create(req *CreateReq) (*subold.Subscriber, error) {
	// TODO: log
	err := req.validateBasic()
	if err != nil {
		return nil, err
	}
	req.fix()
	sub, err := c.getByEmail(req.Email)
	if err != nil {
		if err != c.db.ErrNoSuchEntity() {
			return nil, errDatastore.WithError(err).Annotate("failed to sub.getByEmail")
		}
		// not error
		sub = new(subold.Subscriber)
	}
	if sub.Active {
		return nil, errInvalidParameter.WithMessage("You already have a subscription! :)")
	}
	if !sub.SignUpDatetime.IsZero() {
		return nil, errInvalidParameter.WithMessage("You already have an account. Please reactivate your account by emailing hello@eatgigamunch.com!")
	}
	subSameAddress, err := c.getByAddress(&req.Address)
	if err != nil {
		return nil, errors.Annotate(err, "failed to s.getByAddress")
	}
	if subSameAddress != nil && subSameAddress.Email() != req.Email {
		// same address person
		return nil, errInvalidParameter.WithMessage("Someone in your household is already a subscriber. Please ask them to reactivate their account. If this is not the case, please email: hello@eatgigamunch.com")
	}
	// create payment customer and save payment info
	paymentC, err := payment.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return nil, errInternal.WithError(err).Annotate("failed to payment.NewClient")
	}
	paymentCustomerID, paymentMethodToken, err := paymentC.CreateCustomer(req.PaymentMethodNonce, req.Email, req.FirstName, req.LastName)
	if err != nil {
		return nil, errInternal.WithError(err).Annotate("failed to payment.Create")
	}
	// save subscriber info
	if sub.CreatedDatetime.IsZero() {
		sub.CreatedDatetime = time.Now()
	}
	sub.AddEmail(req.Email, req.FirstName, req.LastName, true)
	sub.AddPhoneNumber(req.PhoneNumber)
	sub.Address = req.Address
	sub.DeliveryNotes = req.DeliveryNotes
	sub.ServingsNonVegetarian = req.ServingsNonVegetarian
	sub.ServingsVegetarian = req.ServingsVegetarian
	sub.PaymentCustomerID = paymentCustomerID
	sub.PaymentMethodToken = paymentMethodToken
	sub.ReferenceText = req.Reference
	sub.ReferenceEmail = req.ReferenceEmail
	sub.AddCampains(req.Campaigns)
	// check if in service zone
	geofenceC, err := geofence.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return nil, errors.Annotate(err, "failed to geofence.NewClient")
	}
	inServiceZone, err := geofenceC.InServiceZone(&req.Address)
	if err != nil {
		return nil, errors.Annotate(err, "failed to geofence.InServiceZone")
	}
	if inServiceZone {
		// signup if in service zone
		err = req.validateAll()
		if err != nil {
			return nil, err
		}
		sub.PlanInterval = 7
		sub.PlanWeekday = req.FirstDeliveryDate.Weekday().String()
		sub.IntervalStartPoint = req.FirstDeliveryDate
		sub.ServingsNonVegetarian = req.ServingsNonVegetarian
		sub.ServingsVegetarian = req.ServingsVegetarian
		sub.Amount = DerivePrice(req.ServingsNonVegetarian + req.ServingsVegetarian)
		sub.Active = true
		now := time.Now()
		sub.SignUpDatetime = now
		sub.ActivateDatetime = now
	}

	err = c.put(sub.ID, sub)
	if err != nil {
		return nil, errDatastore.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Annotate("failed to db.Put")
	}
	// error if not in service zone
	slackC, err := slack.NewClient(c.ctx, c.log, c.serverInfo)
	if err != nil {
		return nil, errors.Annotate(err, "failed to slack.NewClient")
	}
	if !inServiceZone {
		if req.Address.Street == "" {
			return nil, errInvalidParameter.WithMessage("Please select an address.")
		}
		err = slackC.SendMissedSubscriber(req.Email, req.FirstName+" "+req.LastName, req.Reference+" - "+req.ReferenceEmail, req.Campaigns, req.Address)
		if err != nil {
			c.log.Errorf(c.ctx, "failed to slack.SendMissedSubscriber: %+v", err)
		}
		return nil, errInvalidParameter.WithMessage("Sorry, you are outside our delivery range! We'll let you know as soon as we are in your area!")
	}

	if !strings.Contains(req.Email, "@test.com") {
		err = slackC.SendNewSignup(req.Email, req.FirstName+" "+req.LastName, req.Reference+" - "+req.ReferenceEmail, req.Campaigns)
		if err != nil {
			c.log.Errorf(c.ctx, "failed to slack.SendNewSignup: %+v", err)
		}
	}

	// setup activity
	activityC, err := activity.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		c.log.Errorf(c.ctx, "%+v", errors.Annotate(err, "failed to activity.NewClient"))
	} else {
		createReq := &activity.CreateReq{
			Date:                  sub.IntervalStartPoint.Format(activity.DateFormat),
			UserID:                sub.ID,
			Email:                 sub.Email(),
			FirstName:             sub.FirstName(),
			LastName:              sub.LastName(),
			Location:              sub.Location,
			Active:                true,
			Skip:                  false,
			First:                 true,
			ServingsNonVegetarian: sub.ServingsNonVegetarian,
			ServingsVegetarain:    sub.ServingsVegetarian,
			Amount:                sub.Amount,
			// DiscountAmount:        req.DiscountAmount, // switched to new discount system
			// DiscountPercent:       req.DiscountPercent,
			PaymentProvider:    sub.PaymentProvider,
			PaymentMethodToken: sub.PaymentMethodToken,
			CustomerID:         sub.PaymentCustomerID,
		}
		createReq.SetAddress(&req.Address)
		err = activityC.Create(createReq)
		if err != nil {
			c.log.Errorf(c.ctx, "%+v", errors.Annotate(err, "failed to activity.Create"))
		}
	}
	// create discount
	discountC, err := discount.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		c.log.Errorf(c.ctx, "%+v", errors.Annotate(err, "failed to discount.NewClient"))
	} else {
		createReq := &discount.CreateReq{
			UserID:          sub.ID,
			Email:           sub.Email(),
			FirstName:       sub.FirstName(),
			LastName:        sub.LastName(),
			DiscountAmount:  req.DiscountAmount,
			DiscountPercent: req.DiscountPercent,
		}
		err = discountC.Create(createReq)
		if err != nil {
			c.log.Errorf(c.ctx, "%+v", errors.Annotate(err, "failed to discount.Create"))
		}
	}

	// Add to mail service
	taskC := tasks.New(c.ctx)
	// add to update drip
	err = taskC.AddUpdateDrip(time.Now(), &tasks.UpdateDripParams{UserID: sub.ID, Email: req.Email})
	if err != nil {
		c.log.Errorf(c.ctx, "failed to task.AddUpdateDrip: %+v", err)
	}
	return sub, nil
}

// SetupActivity sets up an activity for a subscriber.
func (c *Client) SetupActivity(date time.Time, userIDOrEmail string, active bool, discountAmount float32, discountPrecent int8) error {
	if date.IsZero() || userIDOrEmail == "" {
		return errInvalidParameter.Annotate("date or user id is nil")
	}
	sub, err := c.getByIDOrEmail(userIDOrEmail)
	if err != nil {
		return errors.Annotate(err, "failed to sub.GetByIDOrEmail")
	}
	// setup activity
	activityC, err := activity.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}

	createReq := &activity.CreateReq{
		Date:                  date.Format(activity.DateFormat),
		UserID:                sub.ID,
		Email:                 sub.Email(),
		FirstName:             sub.FirstName(),
		LastName:              sub.LastName(),
		Location:              sub.Location,
		Active:                active,
		Skip:                  false,
		ServingsNonVegetarian: sub.ServingsNonVegetarian,
		ServingsVegetarain:    sub.ServingsVegetarian,
		Amount:                sub.Amount,
		DiscountAmount:        discountAmount,
		DiscountPercent:       discountPrecent,
		PaymentProvider:       sub.PaymentProvider,
		PaymentMethodToken:    sub.PaymentMethodToken,
		CustomerID:            sub.PaymentCustomerID,
		// First: TODO: check activities if any are not skipped
	}
	createReq.SetAddress(&sub.Address)
	err = activityC.Create(createReq)
	if err != nil {
		return errors.Annotate(err, "failed to activity.Create")
	}
	return nil
}

// func (c *Client) Discount()

// SetupActivities updates a subscriber.
func (c *Client) SetupActivities(date time.Time) error {
	// subs, err := c.GetActive(0, 10000)
	// if err != nil {
	// 	return errDatastore.WithError(err).Annotate("failed to put")
	// }
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.SetupSubLogs(date)
}

// ProcessActivity processes an activity.
func (c *Client) ProcessActivity(date time.Time, userIDOrEmail string) error {
	c.log.Infof(c.ctx, "Processing Sub: date(%v) userIDOrEmail(%s)", date, userIDOrEmail)
	// get sub
	sub, err := c.getByIDOrEmail(userIDOrEmail)
	if err != nil {
		if err == c.db.ErrNoSuchEntity() {
			return errInvalidParameter.WithError(err).Annotate("cannot find sub")
		}
		return errDatastore.WithError(err).Annotate("faield to getByIDOrEmail")
	}
	// check if deactivated
	if !sub.Active && (sub.DeactivatedDatetime.IsZero() || sub.DeactivatedDatetime.Before(date)) {
		// user is deactivated
		c.log.Infof(c.ctx, "did not to proess: sub was deactivated before date")
		return nil
	}
	// get activity

	activityC, err := activity.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	act, err := activityC.Get(date, sub.ID)
	if err != nil {
		return errors.Annotate(err, "failed to activity.Get")
	}
	c.log.Infof(c.ctx, "act: %+v", act)
	if act.Skip {
		c.log.Infof(c.ctx, "user is skipped")
		return nil
	}
	// check if should delay processing
	taskC := tasks.New(c.ctx)
	dayBeforeBox := act.DateParsed().Add(-24 * time.Hour)
	if time.Now().Before(dayBeforeBox) {
		// too early to process
		r := &tasks.ProcessSubscriptionParams{
			UserID:   sub.ID,
			SubEmail: sub.Email(),
			Date:     act.DateParsed(),
		}
		err = taskC.AddProcessSubscription(dayBeforeBox, r)
		if err != nil {
			return errors.Wrap("failed to tasks.AddProcessSubscription", err)
		}
		c.log.Infof(c.ctx, "Too early to process Sub. now(%v) < dayBeforeBox(%v)", time.Now(), dayBeforeBox)
		return nil
	}
	r := &tasks.UpdateDripParams{
		Email:  sub.Email(),
		UserID: sub.ID,
	}
	err = taskC.AddUpdateDrip(dayBeforeBox, r)
	if err != nil {
		c.log.Errorf(c.ctx, "failed to tasks.AddUpdateDrip at %s: %+v", dayBeforeBox, err)
	}
	// done if paid
	if act.Paid {
		c.log.Infof(c.ctx, "Subscription is already finished. Paid(%v)", act.Paid)
		return nil
	}

	amount := act.Amount
	var discountAmount float32
	var discountPercent int8
	var discnt *discount.Discount
	discountC, err := discount.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to discount.NewClient")
	}

	if act.DiscountAmount > 0.0 || act.DiscountPercent > 0 {
		// handle applied discount or old discount system
		discountAmount = act.DiscountAmount
		discountPercent = act.DiscountPercent

	} else {
		// new discount system
		// get unused discount
		discnt, err = discountC.GetUnusedUserDiscount(sub.ID)
		if err != nil {
			return errors.Annotate(err, "failed to discount.NewClient")
		}
		if discnt != nil && !discnt.IsUsed() {
			c.log.Infof(c.ctx, "Using discount: %+v", discnt)
			discountAmount = discnt.DiscountAmount
			discountPercent = discnt.DiscountPercent
		}
	}
	amount -= discountAmount
	amount -= (float32(discountPercent) / 100) * amount
	// charge
	orderID := fmt.Sprintf("Gigamunch dinner for %s.", date.Format("01/02/2006"))
	var tID string
	if amount > 0.0 {
		paymentC := paymentold.New(c.ctx)
		saleReq := &paymentold.SaleReq{
			CustomerID:         act.CustomerID,
			Amount:             amount,
			PaymentMethodToken: act.PaymentMethodToken,
			OrderID:            orderID,
		}
		c.log.Infof(c.ctx, "Charging Customer(%s) %f on card(%s)", act.CustomerID, amount, act.PaymentMethodToken)
		tID, err = paymentC.Sale(saleReq)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				// Dulicate transaction error because two customers have same card
				r := &tasks.ProcessSubscriptionParams{
					UserID:   sub.ID,
					SubEmail: sub.Email(),
					Date:     act.DateParsed(),
				}
				err = taskC.AddProcessSubscription(time.Now().Add(1*time.Hour), r)
				if err != nil {
					return errors.Annotate(err, "failed to tasks.AddProcessSubscription")
				}
				return nil
			}
			return errors.Annotate(err, "failed to payment.Sale")
		}
		c.log.Infof(c.ctx, "Charge successful Customer(%s) %f TransactionID(%s)", act.CustomerID, amount, tID)
		c.log.Paid(sub.ID, sub.Email(), act.Date, act.Amount, amount, tID)
	}
	// update TransactionID
	err = activityC.Paid(act.DateParsed(), act.UserID, amount, discountAmount, discountPercent, tID)
	if err != nil {
		c.log.Criticalf(c.ctx, "user paid but didn't get marked as paid: %+v", err)
		return errors.Annotate(err, "user paid but didn't get marked as paid")
	}
	// mark as used
	if discnt != nil && !discnt.IsUsed() {
		t := act.DateParsed()
		err = discountC.Used(discnt.ID, &t)
		if err != nil {
			c.log.Criticalf(c.ctx, "user paid but didn't get marked as paid: %+v", err)
			return errors.Annotate(err, "failed to marked used discount as used")
		}
	}

	return nil
}

// IncrementPageCount is when a user just leaves their email.
func (c *Client) IncrementPageCount(email string, referralPageOpens int, referredPageOpens int) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	s, err := suboldC.GetSubscriber(email)
	if err != nil {
		return err
	}
	s.ReferralPageOpens += referralPageOpens
	s.ReferredPageOpens += referredPageOpens
	err = subold.Put(c.ctx, email, s)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to subold.put")
	}
	return nil
}

// GetCleanPhoneNumber takes a raw phone number and formats it to clean phone number.
func GetCleanPhoneNumber(rawNumber string) string {
	reg := regexp.MustCompile("[^0-9]+")
	cleanNumber := reg.ReplaceAllString(rawNumber, "")
	if len(cleanNumber) < 10 {
		return cleanNumber
	}
	cleanNumber = cleanNumber[len(cleanNumber)-10:]
	cleanNumber = cleanNumber[:3] + "-" + cleanNumber[3:6] + "-" + cleanNumber[6:]
	return cleanNumber
}

// DerivePrice returns the price for a set number of servings.
func DerivePrice(servings int8) float32 {
	switch servings {
	case 1:
		return 17 + 1.66
	case 2:
		return (16.5 * 2) + 3.22
	case 4:
		return (15.25 * 4) + 5.95
	default:
		return 15.25 * float32(servings) * 1.0975
	}
}

func (c *Client) BatchUpdateActivityWithUserID(start, limit int) error {
	subs, err := c.getAll(start, limit)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to getAll")
	}
	if len(subs) == 0 {
		return nil
	}
	max := int(start + limit)
	if len(subs) < max {
		max = len(subs)
	}
	c.log.Infof(c.ctx, "updating %s subs", len(subs))
	subs = subs[0 : max-1]
	userIDs := make([]string, len(subs))
	emails := make([]string, len(subs))
	for i := range subs {
		userIDs[i] = subs[i].ID
		emails[i] = subs[i].Email()
	}
	activityC, err := activity.NewClient(c.ctx, c.log, c.db, c.sqlDB, c.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.BatchUpdateActivityWithUserID(userIDs, emails)
	if err != nil {
		return errors.Annotate(err, "failed to activity.BatchUpdateActivityWithUserID")
	}
	return nil
}

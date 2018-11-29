package sub

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/core/slack"
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
	// if sqlC == nil {
	// 	return nil, errInternal.Annotate("failed to get sql client")
	// }
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
	key := c.db.NameKey(c.ctx, kind, id)
	sub := new(subold.Subscriber)
	err := c.db.Get(c.ctx, key, sub)
	if err != nil {
		if err == c.db.ErrNoSuchEntity() {
			return nil, errNoSuchEntityDatastore.WithError(err).Annotate("failed to get")
		}
		return nil, errDatastore.WithError(err).Annotate("failed to get")
	}
	return sub, nil
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
func (c *Client) ChangeServingsPermanently(email string, servings int8, vegetarian bool) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.ChangeServingsPermanently(email, servings, vegetarian, c.serverInfo)
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
func (c *Client) Deactivate(email string) error {
	// TODO: implement
	suboldC := subold.NewWithLogging(c.ctx, c.log)
	return suboldC.Cancel(email, c.log, c.serverInfo)
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
	if subSameAddress != nil {
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
		sub.Amount = subold.DerivePrice(req.ServingsNonVegetarian + req.ServingsVegetarian)
		sub.Active = true
		now := time.Now()
		sub.SignUpDatetime = now
		sub.ActivateDatetime = now
	}

	err = c.put("", sub)
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
			Date:      sub.IntervalStartPoint.Format(activity.DateFormat),
			UserID:    sub.ID,
			Email:     sub.Email(),
			FirstName: sub.FirstName(),
			LastName:  sub.LastName(),
			Location:  sub.Location,
			Active:    true,
			Skip:      false,
			First:     true,
			ServingsNonVegetarian: sub.ServingsNonVegetarian,
			ServingsVegetarain:    sub.ServingsVegetarian,
			Amount:                sub.Amount,
			DiscountAmount:        req.DiscountAmount,
			DiscountPercent:       req.DiscountPercent,
			PaymentProvider:       sub.PaymentProvider,
			PaymentMethodToken:    sub.PaymentMethodToken,
			CustomerID:            sub.PaymentCustomerID,
		}
		createReq.SetAddress(&req.Address)
		err = activityC.Create(createReq)
		if err != nil {
			c.log.Errorf(c.ctx, "%+v", errors.Annotate(err, "failed to activity.Create"))
		}
	}

	// Add to mail service
	taskC := tasks.New(c.ctx)
	// add to update drip
	err = taskC.AddUpdateDrip(time.Now(), &tasks.UpdateDripParams{Email: req.Email})
	if err != nil {
		c.log.Errorf(c.ctx, "failed to task.AddUpdateDrip: %+v", err)
	}
	// add to task queue
	err = taskC.AddProcessSubscription(sub.IntervalStartPoint.Add(-24*time.Hour), &tasks.ProcessSubscriptionParams{
		SubEmail: req.Email,
		Date:     sub.IntervalStartPoint,
	})
	if err != nil {
		c.log.Errorf(c.ctx, "failed to task.AddProcessSubscription: %+v", err)
	}
	return sub, nil
}

// SetupActivity sets up an activity for a subscriber.
func (c *Client) SetupActivity(date time.Time, userIDOrEmail string, active bool, discountAmount float32, discountPrecent int8) error {
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
		Date:      date.Format(activity.DateFormat),
		UserID:    sub.ID,
		Email:     sub.Email(),
		FirstName: sub.FirstName(),
		LastName:  sub.LastName(),
		Location:  sub.Location,
		Active:    active,
		Skip:      false,
		ServingsNonVegetarian: sub.ServingsNonVegetarian,
		ServingsVegetarain:    sub.ServingsVegetarian,
		Amount:                sub.Amount,
		DiscountAmount:        discountAmount,
		DiscountPercent:       discountPrecent,
		PaymentProvider:       sub.PaymentProvider,
		PaymentMethodToken:    sub.PaymentMethodToken,
		CustomerID:            sub.PaymentCustomerID,
	}
	createReq.SetAddress(&sub.Address)
	err = activityC.Create(createReq)
	if err != nil {
		return errors.Annotate(err, "failed to activity.Create")
	}
	return nil
}

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

func (c *Client) BatchUpdateActivityWithUserID(start, limit int32) error {
	subs, err := c.getAll()
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to getAll")
	}
	max := int(start + limit)
	if len(subs) < max {
		max = len(subs)
	}
	subs = subs[start : max-1]
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

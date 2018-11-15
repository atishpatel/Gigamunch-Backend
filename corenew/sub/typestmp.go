package sub

import (
	"fmt"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

const kindSubscriber = "Subscriber"

// Subscriber is a subscriber.
type Subscriber struct {
	CreatedDatetime time.Time       `json:"created_datetime" datastore:",noindex"`
	SignUpDatetime  time.Time       `json:"sign_up_datetime" datastore:",index"`
	ID              string          `json:"id" datastore:",noindex"`
	AuthID          string          `json:"auth_id"`
	Location        common.Location `json:"location" datastore:",noindex"`
	PhotoURL        string          `json:"photo_url" datastore:",noindex"`
	// Pref
	EmailPrefs []EmailPref `json:"email_prefs"`
	PhonePrefs []PhonePref `json:"phone_prefs"`
	// Account
	PaymentProvider     common.PaymentProvider `json:"payment_provider" datastore:",noindex"`
	PaymentCustomerID   string                 `json:"payment_customer_id"`
	PaymentMethodToken  string                 `json:"payment_method_token"`
	Active              bool                   `json:"active"`
	ActivateDatetime    time.Time              `json:"activate_datetime"`
	DeactivatedDatetime time.Time              `json:"deactivated_datetime"`
	Address             common.Address         `json:"address" datastore:",index"`
	DeliveryNotes       string                 `json:"delivery_notes" datastore:",noindex"`
	// Plan
	ServingsNonVegetarian int8      `json:"servings_non_vegetarian" datastore:",noindex"`
	ServingsVegetarian    int8      `json:"servings_vegetarian" datastore:",noindex"`
	PlanInterval          int8      `json:"plan_interval" datastore:",noindex"`
	PlanWeekday           string    `json:"plan_weekday" datastore:",index"`
	IntervalStartPoint    time.Time `json:"interval_start_point" datastore:",noindex"`
	Amount                float32   `json:"amount" datastore:",noindex"`
	FoodPref              FoodPref  `json:"food_pref"`
	// Gift
	NumGiftDinners     int       `json:"num_gift_dinners" datastore:",noindex"`
	GiftRevealDatetime time.Time `json:"gift_reveal_datetime"`
	// Marketing
	ReferralPageOpens int               `json:"referral_page_opens" datastore:",noindex"`
	ReferredPageOpens int               `json:"referred_page_opens" datastore:",noindex"`
	ReferrerUserID    int64             `json:"referrer_user_id"`
	ReferenceEmail    string            `json:"reference_email"`
	ReferenceText     string            `json:"reference_text" datastore:",noindex"`
	Campaigns         []common.Campaign `json:"campaigns"`
}

// FoodPref are pref for food.
type FoodPref struct {
	NoPork bool `json:"no_pork"`
	NoBeef bool `json:"no_beef"`
}

// EmailPref is a pref for an email.
type EmailPref struct {
	Default   bool   `json:"default" datastore:",noindex"`
	FirstName string `json:"first_name" datastore:",index"`
	LastName  string `json:"last_name" datastore:",index"`
	Email     string `json:"email" datastore:",index"`
}

// PhonePref is a pref for a phone.
type PhonePref struct {
	Number             string `json:"number" datastore:",index"`
	RawNumber          string `json:"raw_number" datastore:",index"`
	DisableBagReminder bool   `json:"disable_bag_reminder" datastore:",noindex"`
	DisableDelivered   bool   `json:"disable_delivered" datastore:",noindex"`
	DisableReview      bool   `json:"disable_review" datastore:",noindex"`
}

func (sub *Subscriber) Email() string {
	var v string
	for _, emailPref := range sub.EmailPrefs {
		v = emailPref.Email
		if emailPref.Default {
			break
		}
	}
	return v
}

func (sub *Subscriber) FirstName() string {
	var v string
	for _, emailPref := range sub.EmailPrefs {
		v = emailPref.FirstName
		if emailPref.Default {
			break
		}
	}
	return v
}

func (sub *Subscriber) LastName() string {
	var v string
	for _, emailPref := range sub.EmailPrefs {
		v = emailPref.LastName
		if emailPref.Default {
			break
		}
	}
	return v
}

func (sub *Subscriber) AddEmail(email, firstName, lastName string, defaultEmail bool) {
	for i := range sub.EmailPrefs {
		if sub.EmailPrefs[i].Email == email {
			// already exists
			sub.EmailPrefs[i].FirstName = firstName
			sub.EmailPrefs[i].LastName = lastName
			sub.EmailPrefs[i].Default = defaultEmail
			return
		}
	}
	sub.EmailPrefs = append(sub.EmailPrefs, EmailPref{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Default:   defaultEmail,
	})
}

func (sub *Subscriber) AddPhoneNumber(rawNumber string) {
	cleanNumber := getCleanPhoneNumber(rawNumber)
	for i := range sub.PhonePrefs {
		if sub.PhonePrefs[i].RawNumber == rawNumber {
			// already exists
			sub.PhonePrefs[i].Number = cleanNumber
			return
		}
	}
	sub.PhonePrefs = append(sub.PhonePrefs, PhonePref{
		RawNumber: rawNumber,
		Number:    cleanNumber,
	})
}

func (sub *Subscriber) AddCampains(campains []common.Campaign) {
	for _, c := range campains {
		found := false
		for _, loggedC := range sub.Campaigns {
			if loggedC.Campaign != c.Campaign {
				continue
			}
			diff := c.Timestamp.Sub(loggedC.Timestamp)
			if diff < 0 {
				diff *= -1
			}
			if diff < time.Hour {
				found = true
			}
		}
		if !found {
			sub.Campaigns = append(sub.Campaigns, c)
		}
	}
}

func (e *Subscriber) GetSubscriptionSignUp() *SubscriptionSignUp {
	sOld := &SubscriptionSignUp{
		ID:                 e.ID,
		Date:               e.CreatedDatetime,
		CustomerID:         e.PaymentCustomerID,
		IsSubscribed:       e.Active,
		SubscriptionDate:   e.ActivateDatetime,
		UnSubscribedDate:   e.DeactivatedDatetime,
		FirstBoxDate:       e.IntervalStartPoint,
		Servings:           e.ServingsNonVegetarian,
		VegetarianServings: e.ServingsVegetarian,
		SubscriptionDay:    e.PlanWeekday,
		WeeklyAmount:       e.Amount,
		PaymentMethodToken: e.PaymentMethodToken,
		Reference:          e.ReferenceText,
		ReferenceEmail:     e.ReferenceEmail,
		DeliveryTips:       e.DeliveryNotes,
		NumGiftDinners:     e.NumGiftDinners,
		GiftRevealDate:     e.GiftRevealDatetime,
		ReferralPageOpens:  e.ReferralPageOpens,
		ReferredPageOpens:  e.ReferredPageOpens,
		Campaigns:          e.Campaigns,
	}
	for _, emailPref := range e.EmailPrefs {
		if emailPref.Default {
			sOld.Email = emailPref.Email
			sOld.FirstName = strings.Title(emailPref.FirstName)
			sOld.LastName = strings.Title(emailPref.LastName)
			sOld.Name = emailPref.FirstName + " " + emailPref.LastName
			break
		}
	}
	for _, phonePref := range e.PhonePrefs {
		sOld.PhoneNumber = phonePref.Number
		sOld.RawPhoneNumber = phonePref.RawNumber
		sOld.BagReminderSMS = !phonePref.DisableBagReminder
		break
	}
	sOld.Address = types.Address{
		APT:     e.Address.APT,
		Street:  e.Address.Street,
		City:    e.Address.City,
		State:   e.Address.State,
		Zip:     e.Address.Zip,
		Country: e.Address.Country,
		GeoPoint: types.GeoPoint{
			Latitude:  e.Address.GeoPoint.Latitude,
			Longitude: e.Address.GeoPoint.Longitude,
		},
	}
	return sOld
}

func (e *SubscriptionSignUp) GetSubscriber() *Subscriber {
	snew := &Subscriber{
		ID:                    e.ID,
		CreatedDatetime:       e.Date,
		PaymentCustomerID:     e.CustomerID,
		Active:                e.IsSubscribed,
		ActivateDatetime:      e.SubscriptionDate,
		SignUpDatetime:        e.SubscriptionDate,
		DeactivatedDatetime:   e.UnSubscribedDate,
		ServingsNonVegetarian: e.Servings,
		ServingsVegetarian:    e.VegetarianServings,
		PlanInterval:          7,
		PlanWeekday:           e.SubscriptionDay,
		Amount:                e.WeeklyAmount,
		PaymentMethodToken:    e.PaymentMethodToken,
		ReferenceText:         e.Reference,
		ReferenceEmail:        e.ReferenceEmail,
		DeliveryNotes:         e.DeliveryTips,
		NumGiftDinners:        e.NumGiftDinners,
		GiftRevealDatetime:    e.GiftRevealDate,
		ReferralPageOpens:     e.ReferralPageOpens,
		ReferredPageOpens:     e.ReferredPageOpens,
		Campaigns:             e.Campaigns,
	}
	lastMonday := time.Now()
	if e.SubscriptionDay == "" {
		e.SubscriptionDay = time.Monday.String()
	}
	count := 0
	for lastMonday.Weekday().String() != e.SubscriptionDay {
		lastMonday = lastMonday.Add(-1 * 24 * time.Hour)
		count++
		if count > 8 {
			break
		}
	}
	snew.IntervalStartPoint = lastMonday
	snew.EmailPrefs = []EmailPref{
		EmailPref{
			Default:   true,
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Email:     e.Email,
		},
	}
	snew.PhonePrefs = []PhonePref{
		PhonePref{
			Number:             e.PhoneNumber,
			RawNumber:          e.RawPhoneNumber,
			DisableBagReminder: !e.BagReminderSMS,
		},
	}

	snew.Address = common.Address{
		APT:     e.Address.APT,
		Street:  strings.Title(e.Address.Street),
		City:    strings.Title(e.Address.City),
		State:   e.Address.State,
		Zip:     strings.TrimSpace(e.Address.Zip),
		Country: e.Address.Country,
		GeoPoint: common.GeoPoint{
			Latitude:  e.Address.GeoPoint.Latitude,
			Longitude: e.Address.GeoPoint.Longitude,
		},
	}
	return snew
}

func (c *Client) BatchSubscriptionSignUpToSubscriber(start, limit int64) error {
	getHasSubscribedPointer := func(ctx context.Context, date time.Time) ([]*SubscriptionSignUp, error) {
		query := datastore.NewQuery(kindSubscriptionSignUp).
			Filter("SubscriptionDate>", 0).
			Filter("SubscriptionDate<", date).
			Limit(1000)
		var results []*SubscriptionSignUp
		_, err := query.GetAll(ctx, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	}

	subsold, err := getHasSubscribedPointer(c.ctx, time.Now().Add(100*24*time.Hour))
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to getHasSubscribed")
	}
	if limit > int64(len(subsold)) {
		limit = int64(len(subsold))
	}
	subsold = subsold[start:limit]
	subsnew := make([]*Subscriber, len(subsold))
	keys := make([]*datastore.Key, len(subsold))
	for i := range subsold {
		if subsold[i].ID == "" {
			id := ksuid.New().String()
			keys[i] = datastore.NewKey(c.ctx, kindSubscriber, id, 0, nil)
		} else {
			keys[i] = datastore.NewKey(c.ctx, kindSubscriber, subsold[i].ID, 0, nil)
		}
		subsnew[i] = subsold[i].GetSubscriber()
	}
	keys, err = datastore.PutMulti(c.ctx, keys, subsnew)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to putMulti Subscribers")
	}
	emails := make([]string, len(keys))
	for i := range keys {
		subsold[i].ID = keys[i].StringID()
		subsnew[i].ID = keys[i].StringID()
		emails[i] = subsold[i].Email
	}
	_, err = datastore.PutMulti(c.ctx, keys, subsnew)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to putMulti Subscribers with ID")
	}
	putMulti := func(ctx context.Context, subs []*SubscriptionSignUp) error {
		keys := make([]*datastore.Key, len(subs))
		for i := range subs {
			keys[i] = datastore.NewKey(ctx, kindSubscriptionSignUp, subs[i].Email, 0, nil)
		}
		_, err := datastore.PutMulti(ctx, keys, subs)
		return err
	}
	err = putMulti(c.ctx, subsold)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to putMulti Subscribers")
	}
	if c.log != nil {
		c.log.Infof(c.ctx, "emails: %+v", emails)
	}

	return nil
}

func get(ctx context.Context, id string) (*SubscriptionSignUp, error) {
	query := datastore.NewQuery(kindSubscriber).
		Filter("EmailPrefs.Email=", id)

	var results []*Subscriber
	_, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("failed to find sub by email: length is 0")
	}
	return results[0].GetSubscriptionSignUp(), nil
}

func getMulti(ctx context.Context, ids []string) ([]*SubscriptionSignUp, error) {
	query := datastore.NewQuery(kindSubscriber).
		Filter("EmailPrefs.Email>", "")

	var results []*Subscriber
	_, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}

	dst := make([]*SubscriptionSignUp, len(ids))
	for i, email := range ids {
		for _, sub := range results {
			if sub.Email() == email {
				dst[i] = sub.GetSubscriptionSignUp()
				break
			}
		}
	}
	return dst, nil
}

// getSubscribersByPhoneNumber returns the subscribers via phone number.
func getSubscribersByPhoneNumber(ctx context.Context, number string) ([]*SubscriptionSignUp, error) {
	query := datastore.NewQuery(kindSubscriber).
		Filter("PhonePrefs.Number=", number)

	var results []*Subscriber
	_, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	dst := make([]*SubscriptionSignUp, len(results))
	for i := range results {
		dst[i] = results[i].GetSubscriptionSignUp()
	}
	return dst, nil
}

// getSubscribers returns the list of Subscribers for that day.
func getSubscribers(ctx context.Context) ([]SubscriptionSignUp, error) {
	query := datastore.NewQuery(kindSubscriber).
		Filter("Active=", true)

	var results []*Subscriber
	_, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	dst := make([]SubscriptionSignUp, len(results))
	for i := range results {
		dst[i] = *results[i].GetSubscriptionSignUp()
	}
	return dst, nil
}

// getSubscribersForWeekday returns the list of Subscribers for that day.
func getSubscribersForWeekday(ctx context.Context, subDay string) ([]SubscriptionSignUp, error) {
	query := datastore.NewQuery(kindSubscriber).
		Filter("Active=", true).
		Filter("PlanWeekday=", subDay)

	var results []*Subscriber
	_, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	dst := make([]SubscriptionSignUp, len(results))
	for i := range results {
		dst[i] = *results[i].GetSubscriptionSignUp()
	}
	return dst, nil
}

// getHasSubscribed returns the list of all Subscribers
func getHasSubscribed(ctx context.Context, date time.Time) ([]SubscriptionSignUp, error) {
	yearsAgo := time.Now().Add(-1 * 100 * 365 * 24 * time.Hour)
	query := datastore.NewQuery(kindSubscriber).
		Filter("SignUpDatetime>", yearsAgo)

	var results []*Subscriber
	_, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	dst := make([]SubscriptionSignUp, len(results))
	for i := range results {
		dst[i] = *results[i].GetSubscriptionSignUp()
	}
	return dst, nil
}

func put(ctx context.Context, email string, i *SubscriptionSignUp) error {
	var err error
	sub := i.GetSubscriber()
	if sub.ID == "" {
		sub.ID = ksuid.New().String()
	}
	key := datastore.NewKey(ctx, kindSubscriber, i.ID, 0, nil)
	_, err = datastore.Put(ctx, key, sub)
	if err != nil {
		return err
	}
	if i.ID != sub.ID {
		i.ID = sub.ID
		key = datastore.NewKey(ctx, kindSubscriptionSignUp, i.ID, 0, nil)
		_, err = datastore.Put(ctx, key, i)
		if err != nil {
			return err
		}
	}
	return oldput(ctx, email, i)
}

func Put(ctx context.Context, email string, i *SubscriptionSignUp) error {
	return put(ctx, email, i)
}

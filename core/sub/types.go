package sub

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
)

const kind = "Subscriber"

// Campaign is a campaign a subscriber was a part of.
type Campaign struct {
	Source    string    `json:"source"`
	Medium    string    `json:"medium"`
	Campaign  string    `json:"campaign"`
	Term      string    `json:"term" datastore:",noindex"`
	Content   string    `json:"content" datastore:",noindex"`
	Timestamp time.Time `json:"timestamp"`
}

// Subscriber is a subscriber.
type Subscriber struct {
	CreatedDatetime time.Time       `json:"created_datetime" datastore:",noindex"`
	SignUpDate      time.Time       `json:"sign_up_date" datastore:",noindex"`
	ID              int64           `json:"id" datastore:",noindex"`
	Email           string          `json:"email"`
	AuthID          string          `json:"auth_id"`
	Location        common.Location `json:"location" datastore:",noindex"`
	FirstName       string          `json:"first_name" datastore:",noindex"`
	LastName        string          `json:"last_name" datastore:",noindex"`
	PhotoURL        string          `json:"photo_url" datastore:",noindex"`
	// Pref
	EmailPrefs []EmailPref `json:"email_prefs" datastore:",noindex"`
	PhonePrefs []PhonePref `json:"phone_prefs" datastore:",noindex"`
	// Account
	PaymentProvider    common.PaymentProvider `json:"payment_provider" datastore:",noindex"`
	PaymentCustomerID  string                 `json:"payment_customer_id" datastore:",noindex"`
	PaymentMethodToken string                 `json:"payment_method_token" datastore:",noindex"`
	Active             bool                   `json:"active" datastore:",index"`
	ActivateDate       time.Time              `json:"activate_date" datastore:",noindex"`
	DeactivatedDate    time.Time              `json:"deactivated_date" datastore:",noindex"`
	Address            common.Address         `json:"address" datastore:",noindex"`
	DeliveryNotes      string                 `json:"delivery_notes" datastore:",noindex"`
	// Plan
	Servings           int8      `json:"servings" datastore:",noindex"`
	VegetarianServings int8      `json:"vegetarian_servings" datastore:",noindex"`
	PlanInterval       int8      `json:"plan_interval" datastore:",noindex"`
	IntervalStartPoint time.Time `json:"interval_start_point" datastore:",noindex"`
	Amount             float32   `json:"amount" datastore:",noindex"`
	FoodPref           FoodPref  `json:"food_pref" datastore:",noindex"`
	// Gift
	NumGiftDinners int       `json:"num_gift_dinners" datastore:",noindex"`
	GiftRevealDate time.Time `json:"gift_reveal_date"`
	// Marketing
	ReferralPageOpens int        `json:"referral_page_opens" datastore:",noindex"`
	ReferredPageOpens int        `json:"referred_page_opens" datastore:",noindex"`
	ReferrerUserID    int64      `json:"referrer_user_id" datastore:",noindex"`
	ReferenceEmail    string     `json:"reference_email" datastore:",noindex"`
	ReferenceText     string     `json:"reference_text" datastore:",noindex"`
	Campaigns         []Campaign `json:"campaigns"`
}

// FoodPref are pref for food.
type FoodPref struct {
	NoPork bool `json:"no_pork" datastore:",noindex"`
	NoBeef bool `json:"no_beef" datastore:",noindex"`
}

// EmailPref is a pref for an email.
type EmailPref struct {
	Default   bool   `json:"default" datastore:",noindex"`
	FirstName string `json:"first_name" datastore:",noindex"`
	LastName  string `json:"last_name" datastore:",noindex"`
	Email     string `json:"email" datastore:",index"`
}

// PhonePref is a pref for a phone.
type PhonePref struct {
	Number             string `json:"number" datastore:",index"`
	RawNumber          string `json:"raw_number" datastore:",index"`
	DisableBagReminder bool   `json:"disable_bag_reminder" datastore:",noindex"`
	DisableDelivered   bool   `json:"disable_delivered" datastore:",noindex"`
}

package sub

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
)

const kind = "Subscriber"

// CampaignCategory are catogeries of campaigns the subscriber came from.
type CampaignCategory string

const (
	// Referral is a referral signup.
	Referral CampaignCategory = "referal"
	// Facebook is a Facebook signup.
	Facebook CampaignCategory = "facebook"
)

// Subscriber is a subscriber.
type Subscriber struct {
	CreatedDatetime time.Time       `json:"created_datetime" datastore:",noindex"`
	ID              int64           `json:"id" datastore:",noindex"`
	AuthID          string          `json:"auth_id" datastore:",index"`
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
	DeactivatedDate    time.Time              `json:"deactivated_date" datastore:",noindex"`
	Address            common.Address         `json:"address" datastore:",noindex"`
	DeliveryNotes      string                 `json:"delivery_notes" datastore:",noindex"`
	// Plan
	Servings           int8      `json:"servings" datastore:",noindex"`
	PlanInterval       int8      `json:"plan_interval" datastore:",noindex"`
	IntervalStartPoint time.Time `json:"interval_start_point" datastore:",noindex"`
	Amount             float32   `json:"amount" datastore:",noindex"`
	DinnerTime         int8      `json:"dinner_time" datastore:",noindex"`
	FoodPref           bool      `json:"food_pref" datastore:",noindex"`
	// Trial
	SubscriptionDate time.Time `json:"subscription_date" datastore:",noindex"`
	// TrialStartDate   time.Time `json:"trail_start_date" datastore:",noindex"`
	// TrialEndDate     time.Time `json:"trail_end_date" datastore:",noindex"`
	// Marketing
	NumReferals      int32            `json:"num_referals" datastore:",noindex"`
	ReferrerUserID   int64            `json:"referrer_user_id" datastore:",noindex"`
	ReferenceText    string           `json:"reference_text" datastore:",noindex"`
	CampaignCategory CampaignCategory `json:"campaign_category" datastore:",noindex"`
	CampaignTags     []string         `json:"campaign_tags" datastore:",noindex"`
}

// FoodPref are pref for food.
type FoodPref struct {
	Vegetarian    bool `json:"vegetarian" datastore:",noindex"`
	EatsChicken   bool `json:"eats_chicken" datastore:",noindex"`
	EatsPork      bool `json:"eats_pork" datastore:",noindex"`
	EatsBeef      bool `json:"eats_beef" datastore:",noindex"`
	EatsLamb      bool `json:"eats_lamb" datastore:",noindex"`
	EatsFish      bool `json:"eats_fish" datastore:",noindex"`
	EatsShellfish bool `json:"eats_shellfish" datastore:",noindex"`
}

// EmailPref is a pref for an email.
type EmailPref struct {
	Default             bool   `json:"default" datastore:",noindex"`
	FirstName           string `json:"first_name" datastore:",noindex"`
	LastName            string `json:"last_name" datastore:",noindex"`
	Email               string `json:"email" datastore:",index"`
	DisableCultureEmail bool   `json:"disable_culture_email" datastore:",noindex"`
}

// PhonePref is a pref for a phone.
type PhonePref struct {
	Number             string `json:"number" datastore:",index"`
	DisableBagReminder bool   `json:"disable_bag_reminder" datastore:",noindex"`
	DisableDelivered   bool   `json:"disable_delivered" datastore:",noindex"`
}

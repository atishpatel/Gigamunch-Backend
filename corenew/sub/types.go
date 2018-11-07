package sub

import (
	"regexp"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

// Campaign is a campaign a subscriber was a part of.
type Campaign struct {
	Source    string    `json:"source"`
	Medium    string    `json:"medium"`
	Campaign  string    `json:"campaign"`
	Term      string    `json:"term" datastore:",noindex"`
	Content   string    `json:"content" datastore:",noindex"`
	Timestamp time.Time `json:"timestamp"`
}

type SubscriptionSignUp struct {
	Email              string        `json:"email"`
	Date               time.Time     `json:"date"` // CreatedDate
	Name               string        `json:"name"`
	FirstName          string        `json:"first_name"`
	LastName           string        `json:"last_name"`
	Address            types.Address `json:"address"`
	CustomerID         string        `json:"customer_id"`
	SubscriptionIDs    []string      `json:"subscription_id"`    // depecrated
	FirstPaymentDate   time.Time     `json:"first_payment_date"` // depecrated
	IsSubscribed       bool          `json:"is_subscribed"`
	SubscriptionDate   time.Time     `json:"subscription_date"`
	UnSubscribedDate   time.Time     `json:"unsubscribed_date"`
	FirstBoxDate       time.Time     `json:"first_box_date"`
	Servings           int8          `json:"servings"`
	VegetarianServings int8          `json:"vegetarian_servings"`
	DeliveryTime       int8          `json:"delivery_time"` // depecrated
	SubscriptionDay    string        `json:"subscription_day"`
	WeeklyAmount       float32       `json:"weekly_amount"`
	PaymentMethodToken string        `json:"payment_method_token"`
	Reference          string        `json:"reference" datastore:",noindex"`
	PhoneNumber        string        `json:"phone_number"`
	RawPhoneNumber     string        `json:"raw_phone_number"`
	DeliveryTips       string        `json:"delivery_tips"`
	BagReminderSMS     bool          `json:"bag_reminder_sms" datastore:",noindex"`
	// gift
	NumGiftDinners int       `json:"num_gift_dinners" datastore:",noindex"`
	ReferenceEmail string    `json:"reference_email"`
	GiftRevealDate time.Time `json:"gift_reveal_date"`
	// stats
	ReferralPageOpens int `json:"referral_page_opens" datastore:",noindex"`
	ReferredPageOpens int `json:"referred_page_opens" datastore:",noindex"`
	GiftPageOpens     int `json:"gift_page_opens" datastore:",noindex"`
	GiftedPageOpens   int `json:"gifted_page_opens" datastore:",noindex"`
	// Campaign
	Campaigns []Campaign `json:"campaigns"`
}

// GetName returns the name of subscriber.
func (s *SubscriptionSignUp) GetName() string {
	return s.Name
}

// GetEmail returns the email of subscriber.
func (s *SubscriptionSignUp) GetEmail() string {
	return s.Email
}

// GetFirstDinnerDate returns the first dinner for the subscriber.
func (s *SubscriptionSignUp) GetFirstDinnerDate() time.Time {
	return s.FirstBoxDate
}

// UpdatePhoneNumber takes a raw number and updates the PhoneNumber.
// Note: Make sure to log old raw number to subscriber.
func (s *SubscriptionSignUp) UpdatePhoneNumber(rawNumber string) {
	s.RawPhoneNumber = rawNumber
	s.PhoneNumber = getCleanPhoneNumber(rawNumber)
}

// getCleanPhoneNumber takes a raw phone number and formats it to clean phone number.
func getCleanPhoneNumber(rawNumber string) string {
	reg := regexp.MustCompile("[^0-9]+")
	cleanNumber := reg.ReplaceAllString(rawNumber, "")
	if len(cleanNumber) < 10 {
		return cleanNumber
	}
	cleanNumber = cleanNumber[len(cleanNumber)-10:]
	cleanNumber = cleanNumber[:3] + "-" + cleanNumber[3:6] + "-" + cleanNumber[6:]
	return cleanNumber
}

// SubscriptionLog is an activity done by a sub.
type SubscriptionLog struct {
	Date               time.Time `json:"date"`      // Primary Key
	SubEmail           string    `json:"sub_email"` // Primary Key
	CreatedDatetime    time.Time `json:"created_datetime"`
	Skip               bool      `json:"skip"`
	Servings           int8      `json:"servings"`
	VegServings        int8      `json:"veg_servings"`
	Amount             float32   `json:"amount"`
	AmountPaid         float32   `json:"amount_paid"`
	Paid               bool      `json:"paid"`
	PaidDatetime       time.Time `json:"paid_datetime"`
	DeliveryTime       int8      `json:"delivery_time"` // depecreated
	PaymentMethodToken string    `json:"payment_method_token"`
	TransactionID      string    `json:"transaction_id"`
	Free               bool      `json:"free"`
	DiscountAmount     float32   `json:"discount_amount"`
	DiscountPercent    int8      `json:"discount_percent"`
	CustomerID         string    `json:"customer_id"`
	Refunded           bool      `json:"refunded"`
	RefundedAmount     float32   `json:"refunded_amount"`
}

// SublogSummary is a summary of sublogs for a email;
type SublogSummary struct {
	MinDate             time.Time `json:"min_date,omitempty"`
	MaxDate             time.Time `json:"max_date,omitempty"`
	Email               string    `json:"email,omitempty"`
	NumTotal            int       `json:"num_total,omitempty"`
	NumSkip             int       `json:"num_skip,omitempty"`
	NumPaid             int       `json:"num_paid,omitempty"`
	NumRefunded         int       `json:"num_refunded,omitempty"`
	TotalAmount         float32   `json:"total_amount,omitempty"`
	TotalAmountPaid     float32   `json:"total_amount_paid,omitempty"`
	TotalDiscountAmount float32   `json:"total_discount_amount,omitempty"`
	TotalRefundedAmount float32   `json:"total_refunded_amount"`
	TotalVegServings    int       `json:"total_veg_servings"`
	TotalNonVegServings int       `json:"total_non_veg_servings"`
	Amount              float32   `json:"amount"`
}

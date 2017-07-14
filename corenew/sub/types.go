package sub

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

type SubscriptionSignUp struct {
	Email              string        `json:"email"`
	Date               time.Time     `json:"date"` // CreatedDate
	Name               string        `json:"name"`
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
	DeliveryTime       int8          `json:"delivery_time"`
	SubscriptionDay    string        `json:"subscription_day"`
	WeeklyAmount       float32       `json:"weekly_amount"`
	PaymentMethodToken string        `json:"payment_method_token"`
	Reference          string        `json:"reference" datastore:",noindex"`
	PhoneNumber        string        `json:"phone_number"`
	DeliveryTips       string        `json:"delivery_tips"`
	BagReminderSMS     bool          `json:"bag_reminder_sms" datastore:",noindex"`
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

type SubscriptionLog struct {
	Date               time.Time `json:"date"`      // Primary Key
	SubEmail           string    `json:"sub_email"` // Primary Key
	CreatedDatetime    time.Time `json:"created_datetime"`
	Skip               bool      `json:"skip"`
	Servings           int8      `json:"servings"`
	Amount             float32   `json:"amount"`
	AmountPaid         float32   `json:"amount_paid"`
	Paid               bool      `json:"paid"`
	PaidDatetime       time.Time `json:"paid_datetime"`
	DeliveryTime       int8      `json:"delivery_time"`
	PaymentMethodToken string    `json:"payment_method_token"`
	TransactionID      string    `json:"transaction_id"`
	Free               bool      `json:"free"`
	DiscountAmount     float32   `json:"discount_amount"`
	DiscountPercent    int8      `json:"discount_percent"`
	CustomerID         string    `json:"customer_id"`
	Refunded           bool      `json:"refunded"`
}

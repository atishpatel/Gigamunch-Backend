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
	SubscriptionIDs    []string      `json:"subscription_id"` // depecrated
	IsSubscribed       bool          `json:"is_subscribed"`
	SubscriptionDate   time.Time     `json:"subscription_date"`
	FirstPaymentDate   time.Time     `json:"first_payment_date"` // depecrated
	FirstBoxDate       time.Time     `json:"first_box_date"`
	Servings           int8          `json:"servings"`
	DeliveryTime       int8          `json:"delivery_time"`
	SubscriptionDay    string        `json:"subscription_day"`
	WeeklyAmount       float32       `json:"weekly_amount"`
	PaymentMethodToken string        `json:"payment_method_token"`
}

type SubscriptionLog struct {
	Date               time.Time // Primary Key
	SubEmail           string    // Primary Key
	CreatedDatetime    time.Time
	Skip               bool
	Servings           int8
	Amount             float32
	AmountPaid         float32
	Paid               bool
	PaidDatetime       time.Time
	DeliveryTime       int8
	PaymentMethodToken string
	TransactionID      string
	Free               bool
	DiscountAmount     float32
	DiscountPercent    int8
	CustomerID         string
}

package inquiry

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	kindInquiry = "Inquiry"
)

// State is the possible states for an Inquiry.
var State = struct{ Pending, Accepted, Declined, TimedOut, Canceled, RefundRequested, Refunded, Fulfilled, Paid string }{
	Pending:         "Pending",
	Accepted:        "Accepted",
	Declined:        "Declined",
	TimedOut:        "TimedOut",
	Canceled:        "Canceled",
	RefundRequested: "RefundRequested",
	Refunded:        "Refunded",
	Fulfilled:       "Fulfilled",
	Paid:            "Paid",
}

// EaterAction is the possible states for an Eater in regards to an Inquiry.
var EaterAction = struct{ Pending, Accepted, Canceled, RefundRequested string }{
	Pending:         "Pending",
	Accepted:        "Accepted",
	Canceled:        "Canceled",
	RefundRequested: "RefundRequested",
}

// CookAction is the possible states for a Cook in regards to an Inquiry.
var CookAction = struct{ Pending, Accepted, Declined, Canceled, RefundAccepted, RefundDeclined string }{
	Pending:        "Pending",
	Accepted:       "Accepted",
	Declined:       "Declined",
	Canceled:       "Canceled",
	RefundAccepted: "RefundAccepted",
	RefundDeclined: "RefundDeclined",
}

// PaymentInfo is the payment information related to an Inquiry.
type PaymentInfo struct {
	CookPricePerServing     float32 `json:"cook_price_per_serving" datastore:",noindex"`
	PricePerServing         float32 `json:"price_per_serving" datastore:",noindex"`
	CookPrice               float32 `json:"cook_price" datastore:",noindex"`
	ExchangePrice           float32 `json:"exchange_price" datastore:",noindex"`
	TaxPrice                float32 `json:"tax_price" datastore:",noindex"`
	ServiceFee              float32 `json:"service_fee" datastore:",noindex"`
	GigamunchToCook         float32 `json:"gigamunch_to_cook" datastore:",noindex"`
	AmountOff               float32 `json:"amount_off" datastore:",noindex"`
	AmountOffDelivery       float32 `json:"amount_off_delivery" datastore:",noindex"`
	TotalPrice              float32 `json:"total_price" datastore:",noindex"`
	TotalPriceWithAmountOff float32 `json:"total_price_with_amount_off" datastore:",noindex"`
}

// ItemInfo contains information about the Item related to an Inquiry.
type ItemInfo struct {
	Name               string               `json:"name" datastore:",noindex"`
	Description        string               `json:"description" datastore:",noindex"`
	Photos             []string             `json:"photos" datastore:",noindex"`
	Ingredients        []string             `json:"ingredients" datastore:",noindex"`
	DietaryConcerns    item.DietaryConcerns `json:"dietary_concerns" datastore:",noindex"`
	ServingDescription string               `json:"serving_description" datastore:",noindex"`
}

// ExchangePlanInfo is the basic information need for an exchange plan.
type ExchangePlanInfo struct {
	EaterAddress types.Address `json:"eater_address" datastore:",noindex"`
	CookAddress  types.Address `json:"cook_address" datastore:",noindex"`
	Distance     float32       `json:"distance" datastore:",noindex"`
	Duration     int64         `json:"duration" datastore:",noindex"`
}

// Promo is the information related to the promo
type Promo struct {
	Code             string  `json:"code" datastore:",noindex"`
	FreeDelivery     bool    `json:"free_delivery" datastore:",noindex"`
	FreeDish         bool    `json:"free_dish" datastore:",noindex"`
	PercentOff       int32   `json:"precent_off" datastore:",noindex"`
	AmountOff        float32 `json:"amount_off" datastore:",noindex"`
	BuyOneGetOneFree bool    `json:"buy_one_get_one_free" datastore:",noindex"`
	DiscountCap      float32 `json:"discount_cap" datastore:",noindex"`
}

// Inquiry is an Inquiry made about an Item by an Eater.
type Inquiry struct {
	ID              int64     `json:"id,string" datastore:",noindex"`
	CreatedDateTime time.Time `json:"created_datetime" datastore:",index"`
	CookID          string    `json:"cook_id" datastore:",index"`
	CookName        string    `json:"cook_name" datastore:",noindex"`
	CookEmail       string    `json:"cook_email" datastore:",noindex"`
	EaterID         string    `json:"eater_id" datastore:",index"`
	EaterPhotoURL   string    `json:"eater_photo_url" datastore:",noindex"`
	EaterName       string    `json:"eater_name" datastore:",noindex"`
	ReviewID        int64     `json:"review_id,string" datastore:",noindex"`
	ItemID          int64     `json:"item_id,string" datastore:",index"`
	Item            ItemInfo  `json:"item" datastore:",noindex"`
	MarkedAsDone    bool      `json:"marked_as_done" datastore:",index"`
	State           string    `json:"state" datastore:",noindex"`
	EaterAction     string    `json:"eater_action" datastore:",noindex"`
	CookAction      string    `json:"cook_action" datastore:",noindex"`
	Issue           bool      `json:"issue" datastore:",noindex"`

	// Braintree info
	BTTransactionID                string `json:"bt_transaction_id" datastore:",index"`
	BTRefundTransactionID          string `json:"bt_refund_transaction_id" datastore:",noindex"`
	BTGigamunchToCookTransactionID string `json:"bt_gigamunch_to_cook_transaction_id" datastore:",noindex"`

	// Promo code
	Promo Promo `json:"promo" datastore:",noindex"`

	// Exchange Info
	ExpectedExchangeDateTime time.Time `json:"expected_exchange_datetime" datastore:",index"`

	Servings         int32                `json:"servings" datastore:",noindex"`
	PaymentInfo      PaymentInfo          `json:"payment_info" datastore:",noindex"`
	ExchangeMethod   types.ExchangeMethod `json:"exchange_method" datastore:",noindex"`
	ExchangePlanInfo ExchangePlanInfo     `json:"exchange_plan_info" datastore:",noindex"`

	// Message Info
	MessageID string `json:"message_id" datastore:",noindex"`
}

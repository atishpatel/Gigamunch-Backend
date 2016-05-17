package order

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindOrder is an order made by a Gigamuncher on a item
	kindOrder = "Order"
)

// State is the possible states the order can be in
var State = struct{ Canceled, Pending, Issues, Refunded, Paid string }{
	Canceled: "Canceled",
	Refunded: "Refunded",
	Pending:  "Pending",
	Issues:   "Issues",
	Paid:     "Paid",
}

// PaymentInfo is the payment information related to an order
type PaymentInfo struct {
	BTTransactionID string  `json:"bt_transaction_id" database:",index"`
	Price           float32 `json:"price" datastore:",noindex"`
	ExchangePrice   float32 `json:"exchange_price" datastore:",noindex"`
	GigaFee         float32 `json:"giga_fee" datastore:",noindex"`
	TaxPrice        float32 `json:"tax_price" datastore:",noindex"`
	TotalPrice      float32 `json:"total_price" datastore:",noindex"`
}

// exchangePlanInfo is the basic information need for an exchange plan
type exchangePlanInfo struct {
	GigamuncherAddress types.Address `json:"gigamuncher_address" datastore:",noindex"`
	GigachefAddress    types.Address `json:"gigachef_address" datastore:",noindex"`
	Distance           float32       `json:"distance" datastore:",noindex"`
	Duration           int64         `json:"duration" datastore:",noindex"`
}

// BasicOrderIDs contains all the associated IDs
type BasicOrderIDs struct {
	GigachefID    string `json:"gigachef_id" datastore:",index"`
	GigamuncherID string `json:"gigamuncher_id" datastore:",index"`
	ReviewID      int64  `json:"review_id" datastore:",noindex"`
	PostID        int64  `json:"post_id" datastore:",noindex"`
	ItemID        int64  `json:"item_id" datastore:",noindex"`
}

// PostInfo contains info related to a post
type PostInfo struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
}

// Order is the order made by a gigamuncher for a post
type Order struct {
	CreatedDateTime          time.Time             `json:"created_datetime"`
	ExpectedExchangeDateTime time.Time             `json:"expected_exchange_datetime"`
	PostCloseDateTime        time.Time             `json:"post_close_datetime"`
	PostReadyDateTime        time.Time             `json:"post_ready_datetime"`
	GigachefHasCompleted     bool                  `json:"gigachef_has_completed" datastore:",noindex"`
	State                    string                `json:"state" datastore:",index"`
	BTRefundTransactionID    string                `json:"bt_refund_transaction_id" datastore:",noindex"`
	ZendeskIssueID           int64                 `json:"zendesk_issue_id" datastore:",noindex"`
	GigachefCanceled         bool                  `json:"gigachef_canceled" datastore:",noindex"`
	GigamuncherCanceled      bool                  `json:"gigamuncher_canceled" datastore:",noindex"`
	BasicOrderIDs                                  // embedded
	PostTitle                string                `json:"post_title" datastore:",noindex"`
	PostPhotoURL             string                `json:"post_photo_url" datastore:",noindex"`
	GigamuncherPhotoURL      string                `json:"gigamuncher_photo_url" datastore:",noindex"`
	GigamuncherName          string                `json:"gigamuncher_name" datastore:",noindex"`
	PricePerServing          float32               `json:"price_per_serving" datastore:",noindex"`
	ChefPricePerServing      float32               `json:"chef_price_per_serving" datastore:",noindex"`
	Servings                 int32                 `json:"servings" datastore:",noindex"`
	PaymentInfo              PaymentInfo           `json:"payment_info" datastore:",noindex"`
	ExchangeMethod           types.ExchangeMethods `json:"exchange_method" datastore:",noindex"`
	ExchangePlanInfo         exchangePlanInfo      `json:"exchange_plan_info" datastore:",noindex"`
}

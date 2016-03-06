package order

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindOrder is an order made by a Gigamuncher on a item
	kindOrder = "Order"
)

// PaymentInfo is the payment information related to an order
type PaymentInfo struct {
	Price         float32 `json:"price" datastore:",noindex"`
	ExchangePrice float32 `json:"exchange_price" datastore:",noindex"`
	TaxPrice      float32 `json:"tax_price" datastore:",noindex"`
	TotalPrice    float32 `json:"total_price" datastore:",noindex"`
}

// BasicExchangePlan is the basic information need for an exchange plan
type BasicExchangePlan struct {
	Cost                    float64       `json:"cost" datastore:",noindex"`
	GigamuncherAddress      types.Address `json:"gigamuncher_address" datastore:",noindex"`
	GigachefAddress         types.Address `json:"gigachef_address" datastore:",noindex"`
	DistanceFromGigachef    float64       `json:"distance_from_gigachef" datastore:",noindex"`
	DistanceFromGigamuncher float64       `json:"distance_from_gigamuncher" datastore:",noindex"`
	StartDateTime           time.Time     `json:"start_datetime" datastore:",noindex"`
	EndDateTime             time.Time     `json:"end_datetime" datastore:",noindex"`
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
	CreatedDateTime          time.Time         `json:"created_datetime" datastore:",noindex"`
	ExpectedExchangeDataTime time.Time         `json:"expected_exchange_datatime" datastore:",index"`
	GigachefHasCompleted     bool              `json:"gigachef_has_completed" datastore:",noindex"`
	Paid                     bool              `json:"paid" datastore:",noindex"`
	Refunded                 bool              `json:"refunded" datastore:",noindex"`
	GigachefCanceled         bool              `json:"gigachef_canceled" datastore:",noindex"`
	GigamuncherCanceled      bool              `json:"gigamuncher_canceled" datastore:",noindex"`
	BasicOrderIDs                              // embedded
	PostTitle                string            `json:"post_title" datastore:",noindex"`
	PostPhotoURL             string            `json:"post_photo_url" datastore:",noindex"`
	GigamuncherPhotoURL      string            `json:"gigamuncher_photo_url" datastore:",noindex"`
	GigamuncherName          string            `json:"gigamuncher_name" datastore:",noindex"`
	PaymentInfo              PaymentInfo       `json:"payment_info" datastore:",noindex"`
	ExchangeMethod           int64             `json:"exchange_method" datastore:",noindex"`
	BasicExchangePlan        BasicExchangePlan `json:"basic_exchange_plan" datastore:",noindex"`
}

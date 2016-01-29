package types

import "time"

// PaymentInfo is the payment information related to an order
type PaymentInfo struct {
	Tip           float64 `json:"tip" datastore:",noindex"`
	MealPrice     float64 `json:"meal_price" datastore:",noindex"`
	ExchangePrice float64 `json:"exchange_price" datastore:",noindex"`
	TotalPrice    float64 `json:"total_price" datastore:",noindex"`
}

// BasicExchangePlan is the basic information need for the exchange plan
type BasicExchangePlan struct {
	Cost                    float64   `json:"cost" datastore:",noindex"`
	GigamuncherAddress      Address   `json:"gigamuncher_address" datastore:",noindex"`
	GigachefAddress         Address   `json:"gigachef_address" datastore:",noindex"`
	ExchangeLocation        Address   `json:"exchange_location" datstore:",noindex"`
	DistanceFromGigachef    float64   `json:"distance_from_gigachef" datastore:",noindex"`
	DistanceFromGigamuncher float64   `json:"distance_from_gigamuncher" datastore:",noindex"`
	StartDateTime           time.Time `json:"start_datetime" datastore:",noindex"`
	EndDateTime             time.Time `json:"end_datetime" datastore:",noindex"`
}

// Order is the payment made by a gigamuncher to a gigachef
type Order struct {
	CreatedDateTime          time.Time         `json:"created_datetime" datastore:",noindex"`
	ExpectedExchangeDataTime time.Time         `json:"expected_exchange_datatime" datastore:",index"`
	Delivered                bool              `json:"delivered" datastore:",index"`
	GigachefEmail            string            `json:"gigachef_email" datastore:",index"`
	GigamuncherEmail         string            `json:"gigamuncher_email" datastore:",index"`
	PaymentInfo              PaymentInfo       `json:"payment_info" datastore:",noindex"`
	ExchangeMethod           int64             `json:"exchange_method" datastore:",noindex"`
	BasicExchangePlan        BasicExchangePlan `json:"basic_exchange_plan" datastore:",noindex"`
}

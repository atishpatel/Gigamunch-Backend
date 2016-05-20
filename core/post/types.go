package post

import (
	"time"

	"gitlab.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindPost is a post of an item that was made
	kindPost = "Post"
)

// OrderPost is an order for a post
type OrderPost struct {
	OrderID             int64                 `json:"order_id" datastore:",noindex"`
	GigamuncherID       string                `json:"gigamuncher_id" datastore:",noindex"`
	ExchangeMethod      types.ExchangeMethods `json:"exchange_method" datastore:",noindex"`
	GigamuncherGeopoint types.GeoPoint        `json:"gigamuncher_geopoint" datastore:",noindex"`
	Servings            int32                 `json:"servings" datastore:",noindex"`
}

// GigachefDelivery contains all the information related to gigachef doing delivery
type GigachefDelivery struct {
	Radius        int32   `json:"radius" datastore:",noindex"`
	MaxDuration   int64   `json:"max_duration" datastore:",noindex"`
	TotalDuration int64   `json:"total_duration" datastore:",noindex"`
	Price         float32 `json:"price" datastore:",noindex"`
}

// Post is a public post created by the Gigachef
type Post struct {
	types.BaseItem                                 // embedded
	ItemID                   int64                 `json:"item_id" datastore:",noindex"`
	BTSubMerchantID          string                `json:"-" datastore:",noindex"`
	Title                    string                `json:"title" datastore:",noindex"`
	IsOrderNow               bool                  `json:"is_order_now" datastore:",noindex"`
	GigachefCanceled         bool                  `json:"gigachef_canceled" datastore:",noindex"`
	ClosingDateTime          time.Time             `json:"closing_datetime" datastore:",index"`
	ReadyDateTime            time.Time             `json:"ready_datetime" datastore:",index"`
	ServingsOffered          int32                 `json:"servings_offered" datastore:",noindex"`
	ChefPricePerServing      float32               `json:"chef_price_per_serving" datastore:",noindex"`
	PricePerServing          float32               `json:"price_per_serving" datastore:",noindex"`
	TaxPercentage            float32               `json:"tax_percentage" datastore:",noindex"`
	EstimatedPreperationTime int64                 `json:"estimated_preperation_time" datastore:",noindex"`
	TotalGigachefRevenue     float32               `json:"total_gigachef_revenue" datastore:",noindex"`
	NumServingsOrdered       int32                 `json:"num_servings_ordered" datastore:",noindex"`
	Orders                   []OrderPost           `json:"orders" datastore:",noindex"`
	AvailableExchangeMethods types.ExchangeMethods `json:"available_exchange_methods" datastore:",noindex"`
	GigachefDelivery         GigachefDelivery      `json:"gigachef_delivery" datastore:",noindex"`
	GigachefAddress          types.Address         `json:"gigachef_address" datastore:",noindex"`
}

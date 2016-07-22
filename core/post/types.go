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
	GigamuncherName     string                `json:"gigamuncher_name" datastore:",noindex"`
	GigamuncherPhotoURL string                `json:"gigamuncher_photo_url" datastore:",noindex"`
	GigamuncherAddress  types.Address         `json:"gigamuncher_address" datastore:",noindex"`
	ExchangeWindowIndex int32                 `json:"exchange_window_index" datastore:",noindex"`
	ExchangeTime        time.Time             `json:"exchange_time" datastore:",noindex"`
	ExchangeMethod      types.ExchangeMethods `json:"exchange_method" datastore:",noindex"`
	Servings            int32                 `json:"servings" datastore:",noindex"`
}

// GigachefDelivery contains all the information related to gigachef doing delivery
type GigachefDelivery struct {
	Radius    int32   `json:"radius" datastore:",noindex"`
	BasePrice float32 `json:"base_price" datastore:",noindex"`
}

// ExchangeTimeSegment is the time range where the exchange can be made
type ExchangeTimeSegment struct {
	StartDateTime            time.Time             `json:"start_datetime" datastore:",noindex"`
	EndDateTime              time.Time             `json:"end_datetime" datastore:",noindex"`
	AvailableExchangeMethods types.ExchangeMethods `json:"available_exchange_methods" datastore:",noindex"`
}

// Post is a public post created by the Gigachef
type Post struct {
	types.BaseItem                            // embedded
	ItemID              int64                 `json:"item_id" datastore:",noindex"`
	BTSubMerchantID     string                `json:"-" datastore:",noindex"`
	Title               string                `json:"title" datastore:",noindex"`
	GigachefCanceled    bool                  `json:"gigachef_canceled" datastore:",noindex"`
	ClosingDateTime     time.Time             `json:"closing_datetime" datastore:",index"`
	ExchangeTimes       []ExchangeTimeSegment `json:"exchange_times" datastore:",noindex"`
	ServingsOffered     int32                 `json:"servings_offered" datastore:",noindex"`
	ChefPricePerServing float32               `json:"chef_price_per_serving" datastore:",noindex"`
	PricePerServing     float32               `json:"price_per_serving" datastore:",noindex"`
	TaxPercentage       float32               `json:"tax_percentage" datastore:",noindex"`
	NumServingsOrdered  int32                 `json:"num_servings_ordered" datastore:",noindex"`
	Orders              []OrderPost           `json:"orders" datastore:",noindex"`
	GigachefDelivery    GigachefDelivery      `json:"gigachef_delivery" datastore:",noindex"`
	GigachefAddress     types.Address         `json:"gigachef_address" datastore:",noindex"`
}

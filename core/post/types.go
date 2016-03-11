package post

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindPost is a post of an item that was made
	kindPost = "Post"
)

type PostOrder struct {
	OrderID       int64  `json:"order_id" datastore:",noindex"`
	GigamuncherID string `json:"gigamuncher_id" datastore:",noindex"`
}

// Post is a public post created by the Gigachef
type Post struct {
	types.BaseItem                                          // embedded
	ItemID                   int64                          `json:"item_id" datastore:",noindex"`
	Title                    string                         `json:"title" datastore:",noindex"`
	IsOrderNow               bool                           `json:"is_order_now" datastore:",noindex"`
	GigachefCanceled         bool                           `json:"gigachef_canceled" datastore:",noindex"`
	ClosingDateTime          time.Time                      `json:"closing_datetime" datastore:",index"`
	ReadyDateTime            time.Time                      `json:"ready_datetime" datastore:",index"`
	ServingsOffered          int                            `json:"servings_offered" datastore:",noindex"`
	PricePerServing          float32                        `json:"price_per_serving" datastore:",noindex"`
	EstimatedPreperationTime int64                          `json:"estimated_preperation_time" datastore:",noindex"`
	TotalRevenue             float32                        `json:"total_revenue" datastore:",noindex"`
	NumOrders                int                            `json:"num_orders" datastoer:",noindex"`
	Orders                   []PostOrder                    `json:"orders" datastore:",noindex"`
	AvaliableExchangeMethods types.AvaliableExchangeMethods `json:"avaliable_exchange_methods" datastore:",noindex"`
	GigachefDeliveryRadius   int                            `json:"gigachef_delivery_radius" datastore:",noindex"`
	types.Address                                           // embedded
}

// Validate validates the BaseItem properties.
// The form is valid if errors.Errors.HasErrors() == false.
func (post *Post) Validate() errors.Errors {
	var multipleErrors errors.Errors
	if post.ClosingDateTime.After(post.ReadyDateTime) {
		multipleErrors.AddError(fmt.Errorf("ClosingDateTime cannot be after ReadyDateTime"))
	}
	if post.ServingsOffered < 1 {
		multipleErrors.AddError(fmt.Errorf("ServingsOffered need to be greater than 0"))
	}
	if post.PricePerServing < 1 {
		multipleErrors.AddError(fmt.Errorf("PricePerServing must be more than $1.00"))
	}
	multipleErrors = append(multipleErrors, post.BaseItem.Validate())
	return multipleErrors
}

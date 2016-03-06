package item

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// KindItem is an item in a Gigachef's kitchen
	KindItem = "Item"
)

// Item is the template Gigachefs can use to post a meal
type Item struct {
	types.BaseItem              // embedded
	Title             string    `json:"title" datastore:",index"`
	LastUsedDataTime  time.Time `json:"lastused_datetime" datastore:",index"`
	NumPostsCreated   int       `json:"num_posts_created" datastore:",noindex"`
	NumTotalOrders    int       `json:"num_total_orders" datastore:",noindex"`
	AverageItemRating float32   `json:"average_item_rating" datastore:",index"`
	NumRatings        int       `json:"num_ratings" datastore:",noindex"`
}

// func (item *Item) Validate() errors.Errors {
// 	var multipleErrors errors.Errors
//
// }

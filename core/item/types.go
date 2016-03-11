package item

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindItem is an item in a Gigachef's kitchen
	kindItem = "Item"
)

// Item is the template Gigachefs can use to post a meal
type Item struct {
	types.BaseItem              // embedded
	Title             string    `json:"title"`
	LastUsedDateTime  time.Time `json:"last_used_datetime"`
	NumPostsCreated   int       `json:"num_posts_created"`
	NumTotalOrders    int       `json:"num_total_orders"`
	AverageItemRating float32   `json:"average_item_rating"`
	NumRatings        int       `json:"num_ratings"`
}

// func (item *Item) Validate() errors.Errors {
// 	var multipleErrors errors.Errors
//
// }

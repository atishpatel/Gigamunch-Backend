package item

import "time"

const kindItem = "MenuItem"

// Item is an item on a menu.
type Item struct {
	CreatedDateTime time.Time `json:"created_datetime" datastore:",noindex"`
	Active          bool      `json:"active" datastore:",noindex"`
	// IsCatering          bool            `json:"is_catering" datastore:",noindex"`
	CookID              string          `json:"cook_id" datastore:",index"`
	MenuID              int64           `json:"menu_id,string" datastore:",index"`
	Title               string          `json:"title" datastore:",noindex"`
	Description         string          `json:"description" datastore:",noindex"`
	DietaryConcerns     DietaryConcerns `json:"dietary_concerns" datastore:",noindex"`
	Ingredients         []string        `json:"ingredients" datastore:",noindex"`
	Photos              []string        `json:"photos" datastore:",noindex"`
	CookPricePerServing float32         `json:"cook_price_per_serving" datastore:",noindex"`
	MaxServings         int32           `json:"max_servings" datastore:",noindex"`
	MinServings         int32           `json:"min_servings" datastore:",noindex"`
	// Stats
	NumServingsSold  int32   `json:"num_servings_sold" datastore:",noindex"`
	NumOrdersSold    int32   `json:"num_orders_sold" datastore:",noindex"`
	TotalCookRevenue float32 `json:"total_cook_revenue" datastore:",noindex"`
}

// DietaryConcerns is a list of booleans that address dietary concerns.
type DietaryConcerns int32

package menu

import "time"

const kindMenu = "Menu"

// Menu is a menu of items from a cook.
type Menu struct {
	CreatedDateTime time.Time `json:"created_datetime" datastore:",noindex"`
	EditedDateTime  time.Time `json:"edited_datetime" datastore:",noindex"`
	CookID          string    `json:"cook_id" datastore:",index"`
	Name            string    `json:"name" datastore:",noindex"`
	Color           string    `json:"color" datastore:",noindex"`
}

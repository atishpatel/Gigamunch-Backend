package review

import "time"

const (
	kindReview = "Review"
)

// Review is a review from a Gigamuncher about a post. It also contains a
// response from the Gigachef
type Review struct {
	CreatedDataTime     time.Time      `json:"created_datetime"`
	IsEdited            bool           `json:"is_edited" datastore:",noindex"`
	EditedDateTime      time.Time      `json:"edited_datetime" datastore:",noindex"`
	GigachefID          string         `json:"gigachef_id"`
	GigamuncherID       string         `json:"gigamuncher_id"`
	GigamuncherPhotoURL string         `json:"gigamuncher_photo_url" datastore:",noindex"`
	ItemID              int64          `json:"item_id"`
	OrderID             int64          `json:"order_id" datastore:",noindex"`
	Post                reviewPost     `json:"post" datastore:",noindex"`
	Rating              int            `json:"rating" datastore:",noindex"`
	Text                string         `json:"text" datastore:",noindex"`
	HasResponse         bool           `json:"has_response" datastore:",noindex"`
	Response            reviewResponse `json:"repsonse" datastore:",noindex"`
}

type reviewPost struct {
	ID       int64  `json:"id" datastore:",noindex"`
	Title    string `json:"title" datastore:",noindex"`
	PhotoURL string `json:"photo_url" datastore:",noindex"`
}

type reviewResponse struct {
	CreatedDateTime time.Time `json:"created_datetime" datastore:",noindex"`
	Text            string    `json:"text" datastore:",noindex"`
}

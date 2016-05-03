package gigachef

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindGigachef is used for the basic Gigachef account info
	kindGigachef = "Gigachef"
)

type GigachefRating struct {
	AverageRating       float32 `json:"average_rating" datastore:",index"`
	NumRatings          int     `json:"num_ratings" datastore:",index"`
	NumOneStarRatings   int     `json:"num_one_star_ratings" datastore:",noindex"`
	NumTwoStarRatings   int     `json:"num_two_star_ratings" datastore:",noindex"`
	NumThreeStarRatings int     `json:"num_three_star_ratings" datastore:",noindex"`
	NumFourStarRatings  int     `json:"num_four_star_ratings" datastore:",noindex"`
	NumFiveStarRatings  int     `json:"num_five_star_ratings" datastore:",noindex"`
}

func (r *GigachefRating) updateAvg() {
	if r.NumRatings == 0 {
		r.AverageRating = 0
		return
	}
	totalStars := r.NumOneStarRatings + r.NumTwoStarRatings + r.NumThreeStarRatings + r.NumFourStarRatings + r.NumFiveStarRatings
	r.AverageRating = float32(totalStars) / float32(r.NumRatings)
}

// AddRating adds a rating and updates average rating
func (r *GigachefRating) AddRating(rating int) {
	switch rating {
	case 1:
		r.NumOneStarRatings++
	case 2:
		r.NumTwoStarRatings++
	case 3:
		r.NumThreeStarRatings++
	case 4:
		r.NumFourStarRatings++
	case 5:
		r.NumFiveStarRatings++
	}
	r.NumRatings++
	r.updateAvg()
}

// RemoveRating removes a rating and updates average rating
func (r *GigachefRating) RemoveRating(rating int) {
	switch rating {
	case 1:
		r.NumOneStarRatings--
	case 2:
		r.NumTwoStarRatings--
	case 3:
		r.NumThreeStarRatings--
	case 4:
		r.NumFourStarRatings--
	case 5:
		r.NumFiveStarRatings--
	}
	r.NumRatings--
	r.updateAvg()
}

// Gigachef contains the basic Gigachef account info
type Gigachef struct {
	CreatedDatetime   time.Time     `json:"created_datetime" datastore:",noindex"`
	HasCarInsurance   bool          `json:"has_car_insurance" datastore:",noindex"`
	types.UserDetail                //embedded
	PhoneNumber       string        `json:"phone_number" datastore:",noindex"`
	Address           types.Address `json:"address" datastore:",noindex"`
	DeliveryRange     int32         `json:"delivery_range" datastore:",noindex"`
	SendWeeklySummary bool          `json:"send_weekly_summary" datastore:",noindex"`
	UseEmailOverSMS   bool          `json:"use_email_over_sms" datastore:",noindex"`
	GigachefRating                  // embedded
	NumPosts          int           `json:"num_posts" datastore:",noindex"`
	NumOrders         int           `json:"num_orders" datastore:",noindex"`
	NumFollowers      int           `json:"num_followers" datastore:",index"`
	KitchenPhotoURLs  []string      `json:"kitchen_photo_urls" datastore:",noindex"`
	BTSubMerchantID   string        `json:"-" datastore:",index"`
}

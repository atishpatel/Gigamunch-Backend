package gigachef

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindGigachef is used for the basic Gigachef account info
	kindGigachef = "Gigachef"
)

// Rating is all the rating info for a gigachef
type Rating struct {
	AverageRating       float32 `json:"average_rating" datastore:",index"`
	NumRatings          int     `json:"num_ratings" datastore:",index"`
	NumOneStarRatings   int     `json:"num_one_star_ratings" datastore:",noindex"`
	NumTwoStarRatings   int     `json:"num_two_star_ratings" datastore:",noindex"`
	NumThreeStarRatings int     `json:"num_three_star_ratings" datastore:",noindex"`
	NumFourStarRatings  int     `json:"num_four_star_ratings" datastore:",noindex"`
	NumFiveStarRatings  int     `json:"num_five_star_ratings" datastore:",noindex"`
}

func (r *Rating) updateAvg() {
	if r.NumRatings == 0 {
		r.AverageRating = 0
		return
	}
	totalStars := r.NumOneStarRatings + (r.NumTwoStarRatings * 2) + (r.NumThreeStarRatings * 3) + (r.NumFourStarRatings * 4) + (r.NumFiveStarRatings * 5)
	r.AverageRating = float32(totalStars) / float32(r.NumRatings)
}

// addRating adds a rating and updates average rating
func (r *Rating) addRating(rating int) {
	r.changeRating(rating, 1)
}

// removeRating removes a rating and updates average rating
func (r *Rating) removeRating(rating int) {
	r.changeRating(rating, -1)
}

func (r *Rating) changeRating(rating, value int) {
	switch rating {
	case 1:
		r.NumOneStarRatings += value
	case 2:
		r.NumTwoStarRatings += value
	case 3:
		r.NumThreeStarRatings += value
	case 4:
		r.NumFourStarRatings += value
	case 5:
		r.NumFiveStarRatings += value
	}
	r.NumRatings += value
	r.updateAvg()
}

// Gigachef contains the basic Gigachef account info
type Gigachef struct {
	CreatedDatetime   time.Time     `json:"created_datetime" datastore:",noindex"`
	HasCarInsurance   bool          `json:"has_car_insurance" datastore:",noindex"`
	types.UserDetail                //embedded
	Bio               string        `json:"bio" datastore:",noindex"`
	PhoneNumber       string        `json:"phone_number" datastore:",noindex"`
	Address           types.Address `json:"address" datastore:",noindex"`
	DeliveryRange     int32         `json:"delivery_range" datastore:",noindex"`
	SendWeeklySummary bool          `json:"send_weekly_summary" datastore:",noindex"`
	UseEmailOverSMS   bool          `json:"use_email_over_sms" datastore:",noindex"`
	Rating                          // embedded
	NumPosts          int           `json:"num_posts" datastore:",noindex"`
	NumOrders         int           `json:"num_orders" datastore:",noindex"`
	NumFollowers      int           `json:"num_followers" datastore:",index"`
	KitchenPhotoURLs  []string      `json:"kitchen_photo_urls" datastore:",noindex"`
	SubMerchantStatus string        `json:"sub_merchant_status" datastore:",noindex"`
	BTSubMerchantID   string        `json:"-" datastore:",index"`
	Application       bool          `json:"application" datastore:",noindex"`
	KitchenInspection bool          `json:"kitchen_inspection" datastore:",noindex"`
	BackgroundCheck   bool          `json:"background_check" datastore:",noindex"`
	FoodHandlerCard   bool          `json:"food_handler_card" datastore:",noindex"`
	PayoutMethod      bool          `json:"payout_method" datastore:",noindex"`
	Verified          bool          `json:"verified" datastore:",noindex"`
}

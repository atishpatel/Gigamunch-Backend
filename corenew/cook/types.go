package cook

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const kindCook = "Cook"

// Rating is all the rating info for a cook
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

// Cook contains all the information related to a cook.
type Cook struct {
	ID               string        `json:"id" datastore:",noindex"`
	CreatedDatetime  time.Time     `json:"created_datetime" datastore:",index"`
	types.UserDetail               // embedded
	Bio              string        `json:"bio" datastore:",noindex"`
	PhoneNumber      string        `json:"phone_number" datastore:",noindex"`
	Address          types.Address `json:"address" datastore:",noindex"`
	DeliveryRange    int32         `json:"delivery_range" datastore:",noindex"`
	// TODO add WeekSchedule and ScheduleModifications
	WeekSchedule          []WeekSchedule          `json:"week_schedule" datastore:",noindex"`
	ScheduleModifications []ScheduleModifications `json:"schedule_modifications" datastore:",noindex"`
	Rating                                        // embedded
	SubMerchantStatus     string                  `json:"sub_merchant_status" datastore:",noindex"`
	BTSubMerchantID       string                  `json:"-" datastore:",index"`
	KitchenPhotoURLs      []string                `json:"kitchen_photo_urls" datastore:",noindex"`
	KitchenInspection     bool                    `json:"kitchen_inspection" datastore:",noindex"`
	BackgroundCheck       bool                    `json:"background_check" datastore:",noindex"`
	FoodHandlerCard       bool                    `json:"food_handler_card" datastore:",noindex"`
	Verified              bool                    `json:"verified" datastore:",index"`
	// SocialMedia
	InstagramID string `json:"instagram_id" datastore:",noindex"`
	TwitterID   string `json:"twitter_id" datastore:",noindex"`
	SnapchatID  string `json:"snapchat_id" datastore:",noindex"`
	// Stats
	AverageResponseTime  int `json:"average_response_time" datastore:",noindex"`
	NumResponses         int `json:"num_responses" datastore:",noindex"`
	NumAcceptedInquiries int `json:"num_accepted_inquiries" datastore:",noindex"`
	NumDeclinedInquiries int `json:"num_declined_inquiries" datastore:",noindex"`
	NumIgnoredInquiries  int `json:"num_ignored_inquiries" datastore:",noindex"`
}

// WeekSchedule is used to make a cook's week's schedule.
type WeekSchedule struct {
	DayOfWeek int32 `json:"day_of_week" datastore:",noindex"`
	StartTime int32 `json:"start_time" datastore:",noindex"`
	EndTime   int32 `json:"end_time" datastore:",noindex"`
}

// ScheduleModifications is used to add modifications to a cook's week's schedule.
type ScheduleModifications struct {
	StartDateTime time.Time `json:"start_datetime" datastore:",noindex"`
	EndDateTime   time.Time `json:"end_datetime" datastore:",noindex"`
	Available     bool      `json:"available" datastore:",noindex"`
}

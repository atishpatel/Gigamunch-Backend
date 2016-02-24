package types

import "time"

type GigachefRating struct {
	AverageRating       float64 `json:"average_rating" datastore:",index"`
	NumRatings          int64   `json:"num_ratings" datastore:",index"`
	NumOneStarRatings   int64   `json:"num_one_star_ratings" datastore:",noindex"`
	NumTwoStarRatings   int64   `json:"num_two_star_ratings" datastore:",noindex"`
	NumThreeStarRatings int64   `json:"num_three_star_ratings" datastore:",noindex"`
	NumFourStarRatings  int64   `json:"num_four_star_ratings" datastore:",noindex"`
	NumFiveStarRatings  int64   `json:"num_five_star_ratings" datastore:",noindex"`
}

// Gigachef contains the basic Gigachef account info
type Gigachef struct {
	IsVerified       bool     `json:"is_verified" datastore:",noindex"`
	UserDetail                //embedded
	Address          Address  `json:"address" datastore:",noindex"`
	GigachefRating            // embedded
	GeneralPhotoURLs []string `json:"general_photo_urls" datastore:",noindex"`
	KitchenPhotoURLs []string `json:"kitchen_photo_urls" datastore:",noindex"`
}

type Photo struct {
	CreatedDataTime time.Time `json:"created_datetime" datastore:",noindex"`
	PhotoURL        string    `json:"photo_url" datastore:",noindex"`
}

// Gallery contains picture uploaded by a Gigachef
type Gallery struct {
	Photos []Photo `json:"photos" datastore:",noindex"`
}

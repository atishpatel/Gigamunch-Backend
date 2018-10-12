package execution

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
)

const (
	// DateFormat is the format used by date.
	DateFormat = "2006-01-02"
	// Kind is a datastore kind.
	Kind = "Execution"
)

// Execution is an execution of a culture.
type Execution struct {
	ID              int64           `json:"id,omitempty" datastore:",noindex"`
	Date            string          `json:"date,omitempty" datastore:",index"`
	Location        common.Location `json:"location,omitempty"`
	Publish         bool            `json:"publish,omitempty"`
	CreatedDatetime time.Time       `json:"created_datetime,omitempty"`
	// Info
	Culture     Culture     `json:"culture,omitempty"`
	Content     Content     `json:"content,omitempty"`
	CultureCook CultureCook `json:"culture_cook,omitempty"`
	Dishes      []Dish      `json:"dishes,omitempty"`
	// Diet
	HasPork    bool `json:"has_pork,omitempty"`
	HasBeef    bool `json:"has_beef,omitempty"`
	HasChicken bool `json:"has_chicken,omitempty"`
	// HasWeirdMeat    bool `json:"has_weird_meat,omitempty"`
	// HasFish         bool `json:"has_fish,omitempty"`
	// HasOtherSeafood bool `json:"has_other_seafood,omitempty"`
}

// Content is a collection of urls pointing to content realted to the execution.
type Content struct {
	HeroImageURL             string `json:"hero_image_url,omitempty" datastore:",noindex"`
	CookImageURL             string `json:"cook_image_url,omitempty" datastore:",noindex"`
	HandsPlateNonVegImageURL string `json:"hands_plate_non_veg_image_url,omitempty" datastore:",noindex"`
	HandsPlateVegImageURL    string `json:"hands_plate_veg_image_url,omitempty" datastore:",noindex"`
	DinnerImageURL           string `json:"dinner_image_url,omitempty" datastore:",noindex"`
	SpotifyURL               string `json:"spotify_url,omitempty" datastore:",noindex"`
	YoutubeURL               string `json:"youtube_url,omitempty" datastore:",noindex"`
}

// Culture is the culture in a culture execution.
type Culture struct {
	Country     string `json:"country,omitempty"`
	City        string `json:"city,omitempty"`
	Description string `json:"description,omitempty" datastore:",noindex"`
	Nationality string `json:"nationality,omitempty" datastore:",noindex"`
	Greeting    string `json:"greeting,omitempty" datastore:",noindex"`
	FlagEmoji   string `json:"flag_emoji,omitempty" datastore:",noindex"`
}

// Dish is a dish in a culture execution.
type Dish struct {
	Number             int32  `json:"number,omitempty" datastore:",noindex"`
	Color              string `json:"color,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty" datastore:",noindex"`
	Ingredients        string `json:"ingredients,omitempty" datastore:",noindex"`
	IsForVegetarian    bool   `json:"is_for_vegetarian,omitempty"`
	IsForNonVegetarian bool   `json:"is_for_non_vegetarian,omitempty"`
}

// CultureCook is the culture cook for a culture execution.
type CultureCook struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Story     string `json:"story,omitempty" datastore:",noindex"`
}

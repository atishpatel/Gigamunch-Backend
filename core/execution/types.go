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
	Publish         bool            `json:"publish"`
	CreatedDatetime time.Time       `json:"created_datetime,omitempty" datastore:",noindex"`
	// Info
	Country     Country     `json:"country,omitempty"`
	Content     Content     `json:"content,omitempty"`
	CultureCook CultureCook `json:"culture_cook,omitempty"`
	Dishes      []Dish      `json:"dishes,omitempty"`
	// Diet
	HasPork bool `json:"has_pork,omitempty"`
	HasBeef bool `json:"has_beef,omitempty"`
}

// Content is a collection of urls pointing to content realted to the execution.
type Content struct {
	HeroImageURL       string `json:"hero_image_url,omitempty" datastore:",noindex"`
	CookImageURL       string `json:"cook_image_url,omitempty" datastore:",noindex"`
	HandsPlateImageURL string `json:"hands_plate_image_url,omitempty" datastore:",noindex"`
	DinnerImageURL     string `json:"dinner_image_url,omitempty" datastore:",noindex"`
	SpotifyURL         string `json:"spotify_url,omitempty" datastore:",noindex"`
	YoutubeURL         string `json:"youtube_url,omitempty" datastore:",noindex"`
}

// Country is the country in a culture execution.
type Country struct {
	Country     string `json:"country,omitempty"`
	City        string `json:"city,omitempty"`
	Description string `json:"description,omitempty" datastore:",noindex"`
	// ??
	Adjective string `json:"adjective,omitempty" datastore:",noindex"`
	Hello     string `json:"hello,omitempty" datastore:",noindex"`
	FlagEmoji string `json:"flag_emoji,omitempty" datastore:",noindex"`
}

// Dish is a dish in a culture execution.
type Dish struct {
	Number             int      `json:"number,omitempty" datastore:",noindex"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty" datastore:",noindex"`
	Ingredients        []string `json:"ingredients,omitempty" datastore:",noindex"`
	IsForVegetarian    bool     `json:"is_for_vegetarian,omitempty"`
	IsForNonVegetarian bool     `json:"is_for_non_vegetarian,omitempty"`
}

// CultureCook is the culture cook for a culture execution.
type CultureCook struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Story     string `json:"story,omitempty" datastore:",noindex"`
}

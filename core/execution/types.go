package execution

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
)

// DateFormat is the format used by date.
const DateFormat = "2006-01-02"

// CultureExecution is an execution of a culture.
type CultureExecution struct {
	ID              int64           `json:"id,omitempty" datastore:",noindex"`
	Date            string          `json:"date,omitempty" datastore:",index"`
	Location        common.Location `json:"location,omitempty"`
	CreatedDatetime time.Time       `json:"created_datetime,omitempty" datastore:",noindex"`
	// Info
	Country     Country     `json:"country,omitempty"`
	CultureCook CultureCook `json:"culture_cook,omitempty"`
	Dishes      []Dish      `json:"dishes,omitempty"`
	Music       Music       `json:"music,omitempty" datastore:",noindex"`
	// Diet
	HasPork bool `json:"has_pork,omitempty"`
}

// Country is the country in a culture execution.
type Country struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty" datastore:",noindex"`
	Location    string
	FlagEmoji   string `json:"flag_emoji,omitempty" datastore:",noindex"`
}

// Dish is a dish in a culture execution.
type Dish struct {
	Number             int      `json:"number,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty" datastore:",noindex"`
	Ingredients        []string `json:"ingredients,omitempty" datastore:",noindex"`
	IsForVegetarian    bool     `json:"is_for_vegetarian,omitempty"`
	IsForNonVegetarian bool     `json:"is_for_non_vegetarian,omitempty"`
}

// Music contains all the music info.
type Music struct {
	SpotifyURL string `json:"spotify_url,omitempty"`
	YoutubeURL string `json:"youtube_url,omitempty"`
}

// CultureCook is the culture cook for a culture execution.
type CultureCook struct {
	Name     string `json:"name,omitempty"`
	Bio      string `json:"bio,omitempty"`
	PhotoURL string `json:"photo_url,omitempty"`
}

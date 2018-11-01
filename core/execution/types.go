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
	Date            string          `json:"date" datastore:",index"`
	Location        common.Location `json:"location"`
	Publish         bool            `json:"publish"`
	CreatedDatetime time.Time       `json:"created_datetime"`
	// Info
	Culture       Culture       `json:"culture"`
	Content       Content       `json:"content"`
	CultureCook   CultureCook   `json:"culture_cook"`
	CultureGuide  CultureGuide  `json:"culture_guide"`
	Dishes        []Dish        `json:"dishes"`
	Stickers      []Sticker     `json:"stickers"`
	Notifications Notifications `json:"notifications"`
	// Diet
	HasPork    bool `json:"has_pork"`
	HasBeef    bool `json:"has_beef"`
	HasChicken bool `json:"has_chicken"`
}

// Notifications are notifications the subscribers gets.
type Notifications struct {
	DeliverySMS string `json:"delivery_sms" datastore:",noindex"`
	RatingSMS   string `json:"rating_sms" datastore:",noindex"`
}

// InfoBox is the infobox in a culture guide.
type InfoBox struct {
	Title   string `json:"title" datastore:",noindex"`
	Text    string `json:"text" datastore:",noindex"`
	Caption string `json:"caption" datastore:",noindex"`
	Image   string `json:"image" datastore:",noindex"`
}

// CultureGuide is content related to the culture guide.
type CultureGuide struct {
	InfoBoxes                    []InfoBox `json:"info_boxes" datastore:",noindex"`
	DinnerInstructions           string    `json:"dinner_instructions" datastore:",noindex"`
	MainColor                    string    `json:"main_color" datastore:",noindex"`
	FontName                     string    `json:"font_name" datastore:",noindex"`
	FontStyle                    string    `json:"font_style" datastore:",noindex"`
	FontCaps                     bool      `json:"font_caps" datastore:",noindex"`
	FontNamePostScript           string    `json:"font_name_post_script" datastore:",noindex"`
	VegetarianDinnerInstructions string    `json:"vegetarian_dinner_instructions" datastore:",noindex"`
}

// Content is a collection of urls pointing to content realted to the execution.
type Content struct {
	LandscapeImageURL        string `json:"landscape_image_url" datastore:",noindex"`
	CookImageURL             string `json:"cook_image_url" datastore:",noindex"`
	HandsPlateNonVegImageURL string `json:"hands_plate_non_veg_image_url" datastore:",noindex"`
	HandsPlateVegImageURL    string `json:"hands_plate_veg_image_url" datastore:",noindex"`
	DinnerNonVegImageURL     string `json:"dinner_non_veg_image_url" datastore:",noindex"`
	DinnerVegImageURL        string `json:"dinner_veg_image_url" datastore:",noindex"`
	SpotifyURL               string `json:"spotify_url" datastore:",noindex"`
	YoutubeURL               string `json:"youtube_url" datastore:",noindex"`
	FontURL                  string `json:"font_url" datastore:",noindex"`
}

// Culture is the culture in a culture execution.
type Culture struct {
	Country            string `json:"country"`
	City               string `json:"city"`
	DescriptionPreview string `json:"description_preview" datastore:",noindex"`
	Description        string `json:"description" datastore:",noindex"`
	Nationality        string `json:"nationality" datastore:",noindex"`
	Greeting           string `json:"greeting" datastore:",noindex"`
	FlagEmoji          string `json:"flag_emoji" datastore:",noindex"`
}

// Sticker are reheat stickers for dishes.
type Sticker struct {
	Name                   string `json:"name"`
	Ingredients            string `json:"ingredients"`
	ExtraInstructions      string `json:"extra_instructions"`
	ReheatOption1          string `json:"reheat_option_1" datastore:",noindex"`
	ReheatOption2          string `json:"reheat_option_2" datastore:",noindex"`
	ReheatOption1Preferred bool   `json:"reheat_option_1_preferred" datastore:",noindex"`
	ReheatTime1            string `json:"reheat_time_1" datastore:",noindex"`
	ReheatTime2            string `json:"reheat_time_2" datastore:",noindex"`
	ReheatInstructions1    string `json:"reheat_instructions_1" datastore:",noindex"`
	ReheatInstructions2    string `json:"reheat_instructions_2" datastore:",noindex"`
	EatingTemperature      string `json:"eating_temperature" datastore:",noindex"`
}

// Dish is a dish in a culture execution.
type Dish struct {
	Number             int32  `json:"number" datastore:",noindex"`
	Color              string `json:"color"`
	Name               string `json:"name"`
	DescriptionPreview string `json:"description_preview" datastore:",noindex"`
	Description        string `json:"description" datastore:",noindex"`
	Ingredients        string `json:"ingredients" datastore:",noindex"`
	IsForVegetarian    bool   `json:"is_for_vegetarian"`
	IsForNonVegetarian bool   `json:"is_for_non_vegetarian"`
	IsOnMainPlate      bool   `json:"is_on_main_plate" datastore:",noindex"`
	ImageURL           string `json:"image_url" datastore:",noindex"`
}

// QandA are questions and answers with the culture cook.
type QandA struct {
	Question string `json:"question" datastore:",noindex"`
	Answer   string `json:"answer" datastore:",noindex"`
}

// CultureCook is the culture cook for a culture execution.
type CultureCook struct {
	FirstName    string  `json:"first_name,omitempty"`
	LastName     string  `json:"last_name,omitempty"`
	Story        string  `json:"story,omitempty" datastore:",noindex"`
	StoryPreview string  `json:"story_preview" datastore:",noindex"`
	QandA        []QandA `json:"q_and_a" datastore:",noindex"`
}

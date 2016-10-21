package review

import "time"

// Review is a review from an Eater about an item. It also contains a
// response from the Cook
type Review struct {
	ID              int64     `json:"id,string"`
	CookID          string    `json:"cook_id"`
	EaterID         string    `json:"eater_id"`
	EaterName       string    `json:"eater_name"`
	EaterPhotoURL   string    `json:"eater_photo_url"`
	InquiryID       int64     `json:"inquiry_id,string"`
	ItemID          int64     `json:"item_id,string"`
	ItemName        string    `json:"item_name"`
	ItemPhotoURL    string    `json:"item_photo_url"`
	MenuID          int64     `json:"menu_id,string"`
	CreatedDateTime time.Time `json:"created_datetime"`
	// eater review
	Rating         int32     `json:"rating"`
	Text           string    `json:"text"`
	IsEdited       bool      `json:"is_edited"`
	EditedDateTime time.Time `json:"edited_datetime"`
	// cook response
	HasResponse             bool      `json:"has_response"`
	ResponseCreatedDateTime time.Time `json:"response_created_datetime"`
	ResponseText            string    `json:"response_text"`
}

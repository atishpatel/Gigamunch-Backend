package gigamuncher

import "github.com/atishpatel/Gigamunch-Backend/core/review"

/*
 * This file is for types shared between multiple files.
 */

// Review is a review
type Review struct {
	ID                  int            `json:"id"`
	CreatedDataTime     int            `json:"created_datetime"`
	IsEdited            bool           `json:"is_edited"`
	EditedDateTime      int            `json:"edited_datetime"`
	GigachefID          string         `json:"gigachef_id"`
	GigamuncherID       string         `json:"gigamuncher_id"`
	GigamuncherPhotoURL string         `json:"gigamuncher_photo_url"`
	ItemID              int64          `json:"item_id"`
	OrderID             int64          `json:"order_id"`
	Post                reviewPost     `json:"post"`
	Rating              int            `json:"rating"`
	Text                string         `json:"text"`
	HasResponse         bool           `json:"has_response"`
	Response            reviewResponse `json:"repsonse"`
}

type reviewPost struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
}

type reviewResponse struct {
	CreatedDateTime int    `json:"created_datetime"`
	Text            string `json:"text"`
}

func (r *Review) Set(id int, review *review.Review) {
	r.CreatedDataTime = int(review.CreatedDataTime.Unix())
	r.EditedDateTime = int(review.EditedDateTime.Unix())
	r.Response.CreatedDateTime = int(review.Response.CreatedDateTime.Unix())
	r.ID = id
	r.IsEdited = review.IsEdited
	r.GigachefID = review.GigachefID
	r.GigamuncherID = review.GigamuncherID
	r.GigamuncherPhotoURL = review.GigamuncherPhotoURL
	r.ItemID = review.ItemID
	r.OrderID = review.OrderID
	r.Post.ID = review.Post.ID
	r.Post.Title = review.Post.Title
	r.Post.PhotoURL = review.Post.PhotoURL
	r.Rating = review.Rating
	r.Text = review.Text
	r.HasResponse = review.HasResponse
	r.Response.Text = review.Response.Text
}

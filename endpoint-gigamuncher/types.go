package gigamuncher

import (
	"encoding/json"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/review"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

/*
 * This file is for types shared between multiple files.
 */

// GigatokenOnlyReq is a request with only a gigatoken input
type GigatokenOnlyReq struct {
	Gigatoken string `json:"gigatoken"`
}

func (req *GigatokenOnlyReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *GigatokenOnlyReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// ErrorOnlyResp is a response with only an error with code
type ErrorOnlyResp struct {
	Err errors.ErrorWithCode `json:"err"`
}

// Review is a review
type Review struct {
	ID                  json.Number    `json:"id"`
	ID64                int64          `json:"-"`
	CreatedDataTime     int            `json:"created_datetime"`
	IsEdited            bool           `json:"is_edited"`
	EditedDateTime      int            `json:"edited_datetime"`
	GigachefID          string         `json:"gigachef_id"`
	GigamuncherID       string         `json:"gigamuncher_id"`
	GigamuncherPhotoURL string         `json:"gigamuncher_photo_url"`
	ItemID              json.Number    `json:"item_id"`
	ItemID64            int64          `json:"-"`
	OrderID             json.Number    `json:"order_id"`
	OrderID64           int64          `json:"-"`
	Post                reviewPost     `json:"post"`
	Rating              int            `json:"rating"`
	Text                string         `json:"text"`
	HasResponse         bool           `json:"has_response"`
	Response            reviewResponse `json:"repsonse"`
}

type reviewPost struct {
	ID       json.Number `json:"id"`
	ID64     int64       `json:"-"`
	Title    string      `json:"title"`
	PhotoURL string      `json:"photo_url"`
}

type reviewResponse struct {
	CreatedDateTime int    `json:"created_datetime"`
	Text            string `json:"text"`
}

func (r *Review) set(review *review.Resp) {
	r.ID = itojn(review.ID)
	r.ID64 = review.ID
	r.CreatedDataTime = ttoi(review.CreatedDataTime)
	r.IsEdited = review.IsEdited
	r.EditedDateTime = ttoi(review.EditedDateTime)
	r.GigachefID = review.GigachefID
	r.GigamuncherID = review.GigamuncherID
	r.GigamuncherPhotoURL = review.GigamuncherPhotoURL
	r.ItemID = itojn(review.ItemID)
	r.ItemID64 = review.ItemID
	r.OrderID = itojn(review.OrderID)
	r.OrderID64 = review.OrderID
	r.Post.ID = itojn(review.Post.ID)
	r.Post.ID64 = review.Post.ID
	r.Post.Title = review.Post.Title
	r.Post.PhotoURL = review.Post.PhotoURL
	r.Rating = review.Rating
	r.Text = review.Text
	r.HasResponse = review.HasResponse
	r.Response.CreatedDateTime = ttoi(review.Response.CreatedDateTime)
	r.Response.Text = review.Response.Text
}

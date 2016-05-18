package gigamuncher

import (
	"encoding/json"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/review"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// PostReviewReq is the request for posting a review
type PostReviewReq struct {
	Gigatoken  string      `json:"gigatoken"`
	OrderID    json.Number `json:"order_id,omitempty"`
	OrderID64  int64       `json:"-"`
	ReviewID   json.Number `json:"review_id,omitempty"`
	ReviewID64 int64       `json:"-"`
	Rating     int         `json:"rating"`
	RatingText string      `json:"rating_text"`
}

// Gigatoken returns the Gigatoken string
func (req *PostReviewReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *PostReviewReq) valid() error {
	if req.Rating < 1 || req.Rating > 5 {
		return fmt.Errorf("Rating is out of range.")
	}
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	var err error
	req.OrderID64, err = req.OrderID.Int64()
	if err != nil {
		return fmt.Errorf("error with OrderID: %v", err)
	}
	if req.ReviewID.String() == "" {
		req.ReviewID64 = 0
	} else {
		req.ReviewID64, err = req.ReviewID.Int64()
		if err != nil {
			return fmt.Errorf("error with ReviewID: %v", err)
		}
	}
	return nil
}

// PostReviewResp is the response to posting a review
// returns: review_id, gigatoken, error
type PostReviewResp struct {
	Review Review               `json:"review"`
	Err    errors.ErrorWithCode `json:"err"`
}

// PostReview is an endpoint that creates or updates a review
func (service *Service) PostReview(ctx context.Context, req *PostReviewReq) (*PostReviewResp, error) {
	resp := new(PostReviewResp)
	defer handleResp(ctx, "PostReview", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	reviewC := review.New(ctx)
	review, err := reviewC.PostReview(user, req.ReviewID64, req.Rating, req.RatingText, req.OrderID64)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Review.set(review)
	return resp, nil
}

// GetReviewsReq is the request for getting a list of reviews
type GetReviewsReq struct {
	GigachefID string      `json:"gigachef_id"`
	StartLimit int         `json:"start_limit"`
	EndLimit   int         `json:"end_limit"`
	ItemID     json.Number `json:"item_id"`
	ItemID64   int64       `json:"-"`
}

// valid validates a req
func (req *GetReviewsReq) valid() error {
	if req.StartLimit < 0 || req.EndLimit < 0 {
		return fmt.Errorf("Limit is out of range.")
	}
	if req.StartLimit >= req.EndLimit {
		return fmt.Errorf("StartLimit cannot be greater than or equal to EndLimit.")
	}
	if req.GigachefID == "" {
		return fmt.Errorf("GigachefID cannot be empty")
	}
	var err error
	req.ItemID64, err = req.ItemID.Int64()
	if err != nil {
		return fmt.Errorf("error with ItemID: %v", err)
	}
	return nil
}

// GetReviewsResp is the response to getting a list of reviews
// returns: []reviews, error
type GetReviewsResp struct {
	Reviews []Review             `json:"reviews,omitempty"`
	Err     errors.ErrorWithCode `json:"err"`
}

// GetReviews gets reviews for a Gigachef. If a ItemID is provided, the reviews
// are resorted with a formula that weighs review post date and item relevence
func (service *Service) GetReviews(ctx context.Context, req *GetReviewsReq) (*GetReviewsResp, error) {
	resp := new(GetReviewsResp)
	defer handleResp(ctx, "GetReviews", resp.Err)
	err := req.valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	limit := &types.Limit{
		Start: req.StartLimit,
		End:   req.EndLimit,
	}
	reviewC := review.New(ctx)
	reviews, err := reviewC.GetReviews(req.GigachefID, limit, req.ItemID64)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	for i := range reviews {
		r := Review{}
		r.set(&reviews[i])
		resp.Reviews = append(resp.Reviews, r)
	}
	return resp, nil
}

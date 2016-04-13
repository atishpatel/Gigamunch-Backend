package gigamuncher

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/review"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// PostReviewReq is the request for posting a review
type PostReviewReq struct {
	GigaToken  string `json:"giga_token"`
	OrderID    int    `json:"order_id"`
	ReviewID   int    `json:"review_id"`
	Rating     int    `json:"rating"`
	RatingText string `json:"rating_text"`
}

// Gigatoken returns the GigaToken string
func (req *PostReviewReq) Gigatoken() string {
	return req.GigaToken
}

// Valid validates a req
func (req *PostReviewReq) Valid() error {
	if req.Rating < 1 || req.Rating > 5 {
		return fmt.Errorf("Rating is out of range.")
	}
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty.")
	}
	if req.OrderID == 0 {
		return fmt.Errorf("OrderID is 0.")
	}
	return nil
}

// PostReviewResp is the response to posting a review
// returns: review_id, gigatoken, error
type PostReviewResp struct {
	ReviewID int                  `json:"review_id"`
	Err      errors.ErrorWithCode `json:"err"`
}

// PostReview is an endpoint that creates or updates a review
func (service *Service) PostReview(ctx context.Context, req *PostReviewReq) (*PostReviewResp, error) {
	resp := new(PostReviewResp)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	reviewID, err := review.PostReview(ctx, user, int64(req.ReviewID), req.Rating, req.RatingText, int64(req.OrderID))
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.ReviewID = int(reviewID)
	return resp, nil
}

// GetReviewsReq is the request for getting a list of reviews
type GetReviewsReq struct {
	GigachefID string `json:"gigachef_id"`
	StartLimit int    `json:"start_limit"`
	EndLimit   int    `json:"end_limit"`
	ItemID     int    `json:"item_id"`
}

// Valid validates a req
func (req *GetReviewsReq) Valid() error {
	if req.StartLimit < 0 || req.EndLimit < 0 {
		return fmt.Errorf("Limit is out of range.")
	}
	if req.StartLimit <= req.EndLimit {
		return fmt.Errorf("StartLimit cannot be less than or equal to EndLimit.")
	}
	if (req.EndLimit - req.StartLimit) > 40 {
		return fmt.Errorf("StartLimit cannot be less than or equal to EndLimit.")
	}
	if req.GigachefID == "" {
		return fmt.Errorf("GigachefID cannot be empty")
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
	err := req.Valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	limit := &types.Limit{
		Start: req.StartLimit,
		End:   req.EndLimit,
	}
	reviewIDs, reviews, err := review.GetReviews(ctx, req.GigachefID, limit, int64(req.ItemID))
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	for i := range reviewIDs {
		r := Review{}
		r.Set(int(reviewIDs[i]), &reviews[i])
		resp.Reviews = append(resp.Reviews, r)
	}
	return resp, nil
}

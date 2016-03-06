package gigamuncher

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/review"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// PostReviewReq is the request for posting a review
type PostReviewReq struct {
	GigaToken  string `json:"giga_token"`
	OrderID    int    `json:"order_id"`
	Rating     int    `json:"rating"`
	RatingText string `json:"rating_text"`
	ReviewID   int    `json:"review_id"`
}

// Valid validates a req
func (req *PostReviewReq) Valid() error {
	if req.Rating < 1 || req.Rating > 5 {
		return fmt.Errorf("Rating is out of range.")
	}
	if req.GigaToken == "" {
		return fmt.Errorf("Gtoken is empty.")
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
	Error    errors.ErrorWithCode `json:"error"`
}

// PostReview is an endpoint that creates or updates a review
func (service *Service) PostReview(ctx context.Context, req *PostReviewReq) (resp *PostReviewResp, respErr error) {
	err := req.Valid()
	if err != nil {
		resp.Error = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	user, err := auth.GetUserFromToken(ctx, req.GigaToken)
	if err != nil {
		resp.Error = errors.GetErrorWithCode(err)
		return resp, nil
	}
	reviewID, err := review.PostReview(ctx, user, int64(req.ReviewID), req.Rating, req.RatingText, int64(req.OrderID))
	if err != nil {
		resp.Error = errors.GetErrorWithCode(err)
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
	Reviews []Review             `json:"reviews"`
	Error   errors.ErrorWithCode `json:"error"`
}

// GetReviews gets reviews for a Gigachef. If a ItemID is provided, the reviews
// are resorted with a formula that weighs review post date and item relevence
func (service *Service) GetReviews(ctx context.Context, req *GetReviewsReq) (resp *GetReviewsResp, respErr error) {
	err := req.Valid()
	if err != nil {
		resp.Error = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	limit := &types.Limit{
		Start: req.StartLimit,
		End:   req.EndLimit,
	}
	reviewIDs, reviews, err := review.GetReviews(ctx, req.GigachefID, limit, int64(req.ItemID))
	if err != nil {
		resp.Error = errors.GetErrorWithCode(err)
		return resp, nil
	}
	for i := range reviewIDs {
		resp.Reviews = append(resp.Reviews, Review{ID: int(reviewIDs[i]), Review: reviews[i]})
	}
	return resp, nil
}

// Review is a review
type Review struct {
	ID            int `json:"id"`
	review.Review     // embedded
}

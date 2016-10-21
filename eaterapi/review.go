package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/review"
)

func (s *service) PostReview(ctx context.Context, req *pb.PostReviewRequest) (*pb.PostReviewResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.PostReviewResponse)
	defer handleResp(ctx, "PostReview", resp.Error)

	return resp, nil
}

func (s *service) GetReviews(ctx context.Context, req *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetReviewsResponse)
	defer handleResp(ctx, "GetReviews", resp.Error)
	validateErr := validateGetReviewsRequest(req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	reviewC := review.New(ctx)
	reviews, err := reviewC.GetByCookID(req.CookId, req.ItemId, req.StartIndex, req.EndIndex)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to reviewC.GetByCookID")
		return resp, nil
	}
	resp.Reviews = getPBReviews(reviews)
	return resp, nil
}

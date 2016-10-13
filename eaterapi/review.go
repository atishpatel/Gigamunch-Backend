package main

import (
	"golang.org/x/net/context"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
)

func (s *service) PostReview(ctx context.Context, req *pb.PostReviewRequest) (*pb.PostReviewResponse, error) {
	resp := new(pb.PostReviewResponse)
	defer handleResp(ctx, "PostReview", resp.Error)

	return resp, nil
}

func (s *service) GetReviews(ctx context.Context, req *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	resp := new(pb.GetReviewsResponse)
	defer handleResp(ctx, "GetReviews", resp.Error)

	return resp, nil
}

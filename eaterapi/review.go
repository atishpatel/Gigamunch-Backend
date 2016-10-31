package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/review"
)

func (s *service) PostReview(ctx context.Context, req *pb.PostReviewRequest) (*pb.PostReviewResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.PostReviewResponse)
	defer handleResp(ctx, "PostReview", resp.Error)
	user, validateErr := validatePostReviewRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inq, err := inquiryC.Get(user, req.InquiryId)
	if err != nil {
		resp.Error = getGRPCError(err, "")
		return resp, nil
	}
	itemC := item.New(ctx)
	itm, err := itemC.Get(inq.ItemID)
	if err != nil {
		resp.Error = getGRPCError(err, "")
		return resp, nil
	}
	reviewC := review.New(ctx)
	var photoURL string
	if len(inq.Item.Photos) > 0 {
		photoURL = inq.Item.Photos[0]
	}
	rvw, err := reviewC.Post(user, req.ReviewId, inq.CookID, inq.ID, inq.ItemID, inq.Item.Name, photoURL, itm.MenuID, req.Rating, req.Text)
	if err != nil {
		resp.Error = getGRPCError(err, "")
		return resp, nil
	}
	resp.Review = getPBReview(rvw)
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

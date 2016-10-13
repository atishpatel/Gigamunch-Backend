package main

import (
	"golang.org/x/net/context"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
)

func (s *service) GetItem(ctx context.Context, id *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	resp := new(pb.GetItemResponse)
	defer handleResp(ctx, "GetItem", resp.Error)

	return resp, nil
}

func (s *service) GetFeed(ctx context.Context, req *pb.GetFeedRequest) (*pb.GetFeedResponse, error) {
	resp := new(pb.GetFeedResponse)
	defer handleResp(ctx, "GetFeed", resp.Error)

	return resp, nil
}

func (s *service) LikeItem(ctx context.Context, req *pb.LikeItemRequest) (*pb.ErrorOnlyResponse, error) {
	resp := new(pb.ErrorOnlyResponse)
	defer handleResp(ctx, "LikeItem", resp.Error)

	return resp, nil
}

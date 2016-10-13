package main

import (
	"golang.org/x/net/context"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
)

func (s *service) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	resp := new(pb.SignInResponse)
	defer handleResp(ctx, "SignIn", resp.Error)

	resp.Gigatoken = "hey chris"

	return resp, nil
}

func (s *service) SignOut(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.ErrorOnlyResponse, error) {
	resp := new(pb.ErrorOnlyResponse)
	defer handleResp(ctx, "SignOut", resp.Error)

	return resp, nil
}

func (s *service) RefreshToken(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.RefreshTokenResponse, error) {
	resp := new(pb.RefreshTokenResponse)
	defer handleResp(ctx, "RefreshToken", resp.Error)

	return resp, nil
}

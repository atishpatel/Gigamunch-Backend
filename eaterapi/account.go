package main

import (
	"golang.org/x/net/context"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func (s *service) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	resp := new(pb.SignInResponse)
	defer handleResp(ctx, "SignIn", resp.Error)

	user, gigatoken, err := auth.GetSessionWithGToken(ctx, req.Gtoken)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to auth.GetSessionWithGToken")
	}
	eaterC := eater.New(ctx)
	_, err = eaterC.Update(user)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to eater.Update")
		return resp, nil
	}
	resp.Gigatoken = gigatoken
	return resp, nil
}

func (s *service) SignOut(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.ErrorOnlyResponse, error) {
	resp := new(pb.ErrorOnlyResponse)
	defer handleResp(ctx, "SignOut", resp.Error)
	if req.Gigatoken == "" {
		return resp, nil
	}
	err := auth.DeleteSessionToken(ctx, req.Gigatoken)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to auth.DeleteSessionToken")
	}
	return resp, nil
}

func (s *service) RefreshToken(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.RefreshTokenResponse, error) {
	resp := new(pb.RefreshTokenResponse)
	defer handleResp(ctx, "RefreshToken", resp.Error)
	if req.Gigatoken == "" {
		resp.Error = &pb.Error{Code: errors.CodeInvalidParameter, Message: "Gigatoken cannot be empty."}
		return resp, nil
	}
	newToken, err := auth.RefreshToken(ctx, req.Gigatoken)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to auth.RefreshToken")
		return resp, nil
	}
	resp.Gigatoken = newToken
	return resp, nil
}

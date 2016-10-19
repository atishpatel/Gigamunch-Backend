package main

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
)

var (
	errInvalidGigatoken = &pb.Error{Code: errors.CodeInvalidParameter, Message: "Gigatoken cannot be empty."}
	// errInvalidParameter = &pb.Error{Code: errors.CodeInvalidParameter, Message: "A parameter is invalid."}
)

func validateGigatokenOnlyReq(req *pb.GigatokenOnlyRequest) *pb.Error {
	if req.Gigatoken == "" {
		return errInvalidGigatoken
	}
	return nil
}

func validateGigatokenAndGetUser(ctx context.Context, gigatoken string) (*types.User, *pb.Error) {
	if gigatoken == "" {
		return nil, errInvalidGigatoken
	}
	user, err := auth.GetUserFromToken(ctx, gigatoken)
	if err != nil {
		return nil, getGRPCError(err, "failed to auth.GetUserFromToken")
	}
	return user, nil
}

func validateSelectAddressRequest(ctx context.Context, req *pb.SelectAddressRequest) (*types.User, *pb.Error) {
	if req.Address == nil || req.Address.Street == "" {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "Address is invalid."}
	}
	return validateGigatokenAndGetUser(ctx, req.Gigatoken)
}
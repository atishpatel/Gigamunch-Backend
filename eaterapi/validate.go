package main

import (
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/golang/protobuf/ptypes"

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

func validateLikeItemRequest(ctx context.Context, req *pb.LikeItemRequest) (*types.User, *pb.Error) {
	if req.ItemId == 0 {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "ItemID cannot be 0."}
	}
	if req.MenuId == 0 {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "MenuID cannot be 0."}
	}
	if req.CookId == "" {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "CookID cannot be empty."}
	}
	return validateGigatokenAndGetUser(ctx, req.Gigatoken)
}

func validateRegisterNotificationTokenRequest(ctx context.Context, req *pb.RegisterNotificationTokenRequest) (*types.User, *pb.Error) {
	if req.NotificationToken == "" {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "NotificationToken cannot be empty."}
	}
	return validateGigatokenAndGetUser(ctx, req.Gigatoken)
}

func validateGetItemRequest(req *pb.GetItemRequest) *pb.Error {
	if req.ItemId == 0 {
		return &pb.Error{Code: errors.CodeInvalidParameter, Message: "ItemID cannot be 0."}
	}
	return nil
}

func validatePostReviewRequest(ctx context.Context, req *pb.PostReviewRequest) (*types.User, *pb.Error) {
	if req.InquiryId == 0 {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "InquiryID cannot be 0."}
	}
	if req.Rating == 0 {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "Rating cannot be 0."}
	}
	return validateGigatokenAndGetUser(ctx, req.Gigatoken)
}

func validateGetReviewsRequest(req *pb.GetReviewsRequest) *pb.Error {
	if req.CookId == "" {
		return &pb.Error{Code: errors.CodeInvalidParameter, Message: "CookID cannot be empty."}
	}
	if req.StartIndex < 0 {
		return &pb.Error{Code: errors.CodeInvalidParameter, Message: "StartIndex must be greater than 0."}
	}
	if req.EndIndex <= req.StartIndex {
		return &pb.Error{Code: errors.CodeInvalidParameter, Message: "EndIndex must be greater than StartIndex."}
	}
	return nil
}

func validateMakeInquiryRequest(ctx context.Context, req *pb.MakeInquiryRequest) (*types.User, time.Time, *pb.Error) {
	var exchangeTime time.Time
	exchangeTime, exchangeErr := ptypes.Timestamp(req.ExchangeTime)
	if exchangeErr != nil {
		return nil, exchangeTime, &pb.Error{Code: errors.CodeInvalidParameter, Message: "Exchange Time is invalid.", Detail: exchangeErr.Error()}
	}
	if req.Address == nil || req.Address.Latitude == 0 || req.Address.Longitude == 0 {
		return nil, exchangeTime, &pb.Error{Code: errors.CodeInvalidParameter, Message: "Invalid address."}
	}
	if req.BraintreeNonce == "" {
		return nil, exchangeTime, &pb.Error{Code: errors.CodeInvalidParameter, Message: "BraintreeNonce cannot be empty."}
	}
	if req.Servings == 0 {
		return nil, exchangeTime, &pb.Error{Code: errors.CodeInvalidParameter, Message: "Servings cannot be 0."}
	}
	user, err := validateGigatokenAndGetUser(ctx, req.Gigatoken)
	return user, exchangeTime, err
}

func validateIDReq(ctx context.Context, id int64, gigatoken string) (*types.User, *pb.Error) {
	if id == 0 {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "ID cannot be 0."}
	}
	return validateGigatokenAndGetUser(ctx, gigatoken)
}

func validateGetInquiriesRequest(ctx context.Context, req *pb.GetInquiriesRequest) (*types.User, *pb.Error) {
	if req.StartIndex < 0 {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "StartIndex must be greater than 0."}
	}
	if req.EndIndex <= req.StartIndex {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "EndIndex must be greater than StartIndex."}
	}
	return validateGigatokenAndGetUser(ctx, req.Gigatoken)
}

func validateGetMessageTokenRequest(ctx context.Context, req *pb.GetMessageTokenRequest) (*types.User, *pb.Error) {
	if req.DeviceId == "" {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "DeviceID cannot be empty."}
	}
	return validateGigatokenAndGetUser(ctx, req.Gigatoken)
}

func validateCreateMessageChannelRequest(ctx context.Context, req *pb.CreateMessageChannelRequest) (*types.User, *pb.Error) {
	if req.CookId == "" {
		return nil, &pb.Error{Code: errors.CodeInvalidParameter, Message: "CookID cannot be empty."}
	}
	return validateGigatokenAndGetUser(ctx, req.Gigatoken)
}

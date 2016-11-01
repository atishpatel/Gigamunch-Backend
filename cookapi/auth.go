package main

import (
	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// RefreshTokenResp is the output form the RefreshToken endpoint
type RefreshTokenResp struct {
	Gigatoken string               `json:"gigatoken"`
	Err       errors.ErrorWithCode `json:"err"`
}

// RefreshToken refreshs a token. A token should be refreshed at least ever hour.
func (service *Service) RefreshToken(ctx context.Context, req *GigatokenReq) (*RefreshTokenResp, error) {
	resp := new(RefreshTokenResp)
	defer handleResp(ctx, "RefreshToken", resp.Err)
	var err error
	err = req.valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	newToken, err := auth.RefreshToken(ctx, req.Gigatoken)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Gigatoken = newToken
	return resp, nil
}

// GetMessageTokenReq is the request for GetMessageToken.
type GetMessageTokenReq struct {
	GigatokenReq
	DeviceID string `json:"device_id"`
}

// GetMessageTokenResp is the output form GetMessageToken.
type GetMessageTokenResp struct {
	Token string               `json:"token"`
	Err   errors.ErrorWithCode `json:"err"`
}

// GetMessageToken gets a token for messaging.
func (service *Service) GetMessageToken(ctx context.Context, req *GetMessageTokenReq) (*GetMessageTokenResp, error) {
	resp := new(GetMessageTokenResp)
	defer handleResp(ctx, "GetMessageToken", resp.Err)
	var err error
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	messageC := message.New(ctx)
	userInfo := &message.UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Image: user.PhotoURL,
	}
	tkn, err := messageC.GetToken(userInfo, req.DeviceID)
	if err != nil {
		resp.Err = errors.Wrap("failed to message.GetToken", err)
		return resp, nil
	}
	resp.Token = tkn
	return resp, nil
}

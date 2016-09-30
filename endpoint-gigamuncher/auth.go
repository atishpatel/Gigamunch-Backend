package gigamuncher

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/gigamuncher"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// SignInReq is the input required to sign in
type SignInReq struct {
	Gtoken string `json:"gtoken"`
}

// valid validates a req
func (req *SignInReq) valid() error {
	if req.Gtoken == "" {
		return fmt.Errorf("Gtoken is empty.")
	}
	return nil
}

// SignInResp is the response to signing in
// returns: gigatoken, error
type SignInResp struct {
	Gigatoken string               `json:"gigatoken"`
	Err       errors.ErrorWithCode `json:"err"`
}

// SignIn is an endpoint that signs in a user
func (service *Service) SignIn(ctx context.Context, req *SignInReq) (*SignInResp, error) {
	resp := new(SignInResp)
	defer handleResp(ctx, "SignIn", resp.Err)
	var err error
	err = req.valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	user, gigaToken, err := auth.GetSessionWithGToken(ctx, req.Gtoken)
	resp.Gigatoken = gigaToken
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	err = gigamuncher.SaveUserInfo(ctx, user, nil)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
	}
	return resp, nil
}

// SignOut is an endpoint that signs out a user
func (service *Service) SignOut(ctx context.Context, req *GigatokenOnlyReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SignOut", resp.Err)
	var err error
	if req.Gigatoken == "" {
		return resp, nil
	}
	err = auth.DeleteSessionToken(ctx, req.Gigatoken)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
	}
	return resp, nil
}

// RefreshTokenResp is the output form the RefreshToken endpoint
type RefreshTokenResp struct {
	Gigatoken string               `json:"gigatoken"`
	Err       errors.ErrorWithCode `json:"err"`
}

// RefreshToken refreshs a token. A token should be refreshed at least ever hour.
func (service *Service) RefreshToken(ctx context.Context, req *GigatokenOnlyReq) (*RefreshTokenResp, error) {
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

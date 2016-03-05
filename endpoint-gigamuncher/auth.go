package gigamuncher

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// SignInReq is the input required to sign in
type SignInReq struct {
	Gtoken string `json:"gtoken"`
}

// Valid validates a req
func (req *SignInReq) Valid() error {
	if req.Gtoken == "" {
		return fmt.Errorf("Gtoken is empty.")
	}
	return nil
}

// SignInResp is the response to signing in
// returns: gigatoken, error
type SignInResp struct {
	GigaToken string               `json:"gigatoken"`
	Error     errors.ErrorWithCode `json:"error"`
}

// SignIn is an endpoint that signs in a user
func (service *Service) SignIn(ctx context.Context, req *SignInReq) (resp *SignInResp, respErr error) {
	var err error
	err = req.Valid()
	if err != nil {
		resp.Error = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	_, gigaToken, err := auth.GetSessionWithGToken(ctx, req.Gtoken)
	resp.GigaToken = gigaToken
	if err != nil {
		resp.Error = errors.GetErrorWithCode(err)
	}
	return resp, nil
}

// SignOutReq is the input required to sign out
type SignOutReq struct {
	GigaToken string `json:"gigatoken"`
}

// SignOutResp is the response to signing out
// returns: error
type SignOutResp struct {
	Error errors.ErrorWithCode `json:"error"`
}

// SignOut is an endpoint that signs out a user
func (service *Service) SignOut(ctx context.Context, req *SignOutReq) (resp *SignOutResp, respErr error) {
	var err error
	if req.GigaToken == "" {
		return resp, nil
	}
	err = auth.DeleteSessionToken(ctx, req.GigaToken)
	if err != nil {
		resp.Error = errors.GetErrorWithCode(err)
	}
	return resp, nil
}

// RefreshTokenReq is the input requred to refresh a token
type RefreshTokenReq struct {
	GigaToken string `json:"gigatoken"`
}

// Valid validates a req
func (req *RefreshTokenReq) Valid() error {
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty.")
	}
	return nil
}

// RefreshTokenResp is the output form the RefreshToken endpoint
type RefreshTokenResp struct {
	GigaToken string               `json:"gigatoken"`
	Error     errors.ErrorWithCode `json:"error"`
}

// RefreshToken refreshs a token. A token should be refreshed at least ever hour.
func (service *Service) RefreshToken(ctx context.Context, req *RefreshTokenReq) (resp *RefreshTokenResp, respErr error) {
	var err error
	err = req.Valid()
	if err != nil {
		resp.Error = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	newToken, err := auth.RefreshToken(ctx, req.GigaToken)
	if err != nil {
		resp.Error = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.GigaToken = newToken
	return resp, nil
}

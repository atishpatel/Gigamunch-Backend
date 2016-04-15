package gigamuncher

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/gigamuncher"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
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
	Err       errors.ErrorWithCode `json:"err"`
}

// SignIn is an endpoint that signs in a user
func (service *Service) SignIn(ctx context.Context, req *SignInReq) (*SignInResp, error) {
	resp := new(SignInResp)
	defer func() {
		if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "SignIn err: ", resp.Err)
		}
	}()
	var err error
	err = req.Valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	user, gigaToken, err := auth.GetSessionWithGToken(ctx, req.Gtoken)
	resp.GigaToken = gigaToken
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
	}
	err = gigamuncher.SaveUserInfo(ctx, user, nil)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
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
	Err errors.ErrorWithCode `json:"err"`
}

// SignOut is an endpoint that signs out a user
func (service *Service) SignOut(ctx context.Context, req *SignOutReq) (*SignOutResp, error) {
	resp := new(SignOutResp)
	defer func() {
		if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "SignOut err: ", resp.Err)
		}
	}()
	var err error
	if req.GigaToken == "" {
		return resp, nil
	}
	err = auth.DeleteSessionToken(ctx, req.GigaToken)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
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
	Err       errors.ErrorWithCode `json:"err"`
}

// RefreshToken refreshs a token. A token should be refreshed at least ever hour.
func (service *Service) RefreshToken(ctx context.Context, req *RefreshTokenReq) (*RefreshTokenResp, error) {
	resp := new(RefreshTokenResp)
	defer func() {
		if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "RefreshToken err: ", resp.Err)
		}
	}()
	var err error
	err = req.Valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	newToken, err := auth.RefreshToken(ctx, req.GigaToken)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.GigaToken = newToken
	return resp, nil
}

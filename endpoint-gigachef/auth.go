package gigachef

import (
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/auth"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

// RefreshTokenReq is the input requred to refresh a token
type RefreshTokenReq struct {
	Gigatoken string `json:"gigatoken"`
}

// Valid validates a req
func (req *RefreshTokenReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// RefreshTokenResp is the output form the RefreshToken endpoint
type RefreshTokenResp struct {
	Gigatoken string               `json:"gigatoken"`
	Err       errors.ErrorWithCode `json:"err"`
}

// RefreshToken refreshs a token. A token should be refreshed at least ever hour.
func (service *Service) RefreshToken(ctx context.Context, req *RefreshTokenReq) (*RefreshTokenResp, error) {
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

package gigachef

import (
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

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

package gigachef

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

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

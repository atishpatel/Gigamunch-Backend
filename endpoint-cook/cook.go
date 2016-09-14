package cookendpoint

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// CookResp is the output response with a Cook and error.
type CookResp struct {
	Cook cook.Cook            `json:"cook"`
	Err  errors.ErrorWithCode `json:"err"`
}

// GetCook is an endpoint that get the chef info.
func (service *Service) GetCook(ctx context.Context, req *GigatokenOnlyReq) (*CookResp, error) {
	resp := new(CookResp)
	defer handleResp(ctx, "GetCook", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	cookC := cook.New(ctx)
	cook, err := cookC.Get(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to get cook")
		return resp, nil
	}
	resp.Cook = *cook
	return resp, nil
}

package main

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// CookResp is the output response with a Cook and error.
type CookResp struct {
	Cook cook.Cook            `json:"cook"`
	Err  errors.ErrorWithCode `json:"err"`
}

// GetCook is an endpoint that get the chef info.
func (service *Service) GetCook(ctx context.Context, req *GigatokenReq) (*CookResp, error) {
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

// UpdateCookReq is a request for UpdateCook.
type UpdateCookReq struct {
	Gigatoken string    `json:"gigatoken"`
	Cook      cook.Cook `json:"cook"`
}

func (req *UpdateCookReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *UpdateCookReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// UpdateCook updates the cook profile.
func (service *Service) UpdateCook(ctx context.Context, req *UpdateCookReq) (*CookResp, error) {
	resp := new(CookResp)
	defer handleResp(ctx, "UpdateCook", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	cookC := cook.New(ctx)
	cook, err := cookC.Update(user, &req.Cook.Address, req.Cook.PhoneNumber, req.Cook.Bio, req.Cook.DeliveryRange,
		req.Cook.WeekSchedule, req.Cook.InstagramID, req.Cook.TwitterID, req.Cook.SnapchatID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to update cook")
		return resp, nil
	}
	resp.Cook = *cook
	return resp, nil
}

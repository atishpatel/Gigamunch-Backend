package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
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
		resp.Err = errors.Wrap("failed to update cook", err)
		return resp, nil
	}
	resp.Cook = *cook
	return resp, nil
}

// GetSubMerchantResp is a resp for GetSubMerchant
type GetSubMerchantResp struct {
	SubMerchant SubMerchant          `json:"sub_merchant"`
	Err         errors.ErrorWithCode `json:"err"`
}

// GetSubMerchant gets a submerchant.
func (service *Service) GetSubMerchant(ctx context.Context, req *GigatokenReq) (*GetSubMerchantResp, error) {
	resp := new(GetSubMerchantResp)
	defer handleResp(ctx, "GetSubMerchant", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	cookC := cook.New(ctx)
	cook, err := cookC.Get(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot get cook(%d)", user.ID)
		return resp, nil
	}
	paymentC := payment.New(ctx)
	sm, err := paymentC.GetSubMerchant(cook.BTSubMerchantID)
	if err != nil {
		resp.SubMerchant.Address = cook.Address
		resp.SubMerchant.Email = cook.Email
		nameArray := strings.Split(cook.Name, " ")
		switch len(nameArray) {
		case 3:
			resp.SubMerchant.FirstName = nameArray[0]
			resp.SubMerchant.LastName = nameArray[2]
		case 2:
			resp.SubMerchant.LastName = nameArray[1]
			fallthrough
		case 1:
			resp.SubMerchant.FirstName = nameArray[0]
		}
		utils.Infof(ctx, "cannot get sub-merchant(%s): err: %v", cook.BTSubMerchantID, err)
		return resp, nil
	}
	resp.SubMerchant.SubMerchantInfo = *sm
	return resp, nil
}

// UpdateSubMerchantReq updates sub-merchant payment info
type UpdateSubMerchantReq struct {
	GigatokenReq
	SubMerchant SubMerchant `json:"sub_merchant"`
}

// UpdateSubMerchant creates or updates sub-merchant info
func (service *Service) UpdateSubMerchant(ctx context.Context, req *UpdateSubMerchantReq) (*CookResp, error) {
	resp := new(CookResp)
	defer handleResp(ctx, "UpdateSubMerchant", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	cookC := cook.New(ctx)
	cook, err := cookC.Get(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot get cook(%d)", user.ID)
		return resp, nil
	}
	paymentC := payment.New(ctx)
	req.SubMerchant.ID = cook.BTSubMerchantID
	err = paymentC.UpdateSubMerchant(user, &req.SubMerchant.SubMerchantInfo)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot update sub-merchant(%d)", cook.BTSubMerchantID)
		return resp, nil
	}
	cook.SubMerchantStatus = "active"
	resp.Cook = *cook
	return resp, nil
}

// FinishOnboardingReq is a request for finishing onboarding for a cook.
type FinishOnboardingReq struct {
	GigatokenReq
	Cook        cook.Cook   `json:"cook"`
	SubMerchant SubMerchant `json:"sub_merchant"`
}

// FinishOnboarding is the function to finish the onboarding process
func (service *Service) FinishOnboarding(ctx context.Context, req *FinishOnboardingReq) (*RefreshTokenResp, error) {
	resp := new(RefreshTokenResp)
	defer handleResp(ctx, "FinishOnboarding", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}

	// update cook
	cookC := cook.New(ctx)
	cook, err := cookC.Update(user, &req.Cook.Address, req.Cook.PhoneNumber, req.Cook.Bio, req.Cook.DeliveryRange,
		req.Cook.WeekSchedule, req.Cook.InstagramID, req.Cook.TwitterID, req.Cook.SnapchatID)
	if err != nil {
		resp.Err = errors.Wrap("failed to update cook", err)
		return resp, nil
	}
	req.SubMerchant.ID = cook.BTSubMerchantID
	req.SubMerchant.Email = user.Email

	// create or update submerchant info
	paymentC := payment.New(ctx)
	if req.SubMerchant.AccountNumber != "" {
		// update submerchant with real banking info
		err = paymentC.UpdateSubMerchant(user, &req.SubMerchant.SubMerchantInfo)
		if err != nil {
			resp.Err = errors.Wrap("failed to paymentC.UpdateSubMerchant", err)
			return resp, nil
		}
	} else {
		// create submerchant with fake banking stuff
		if !user.HasSubMerchantID() {
			err = paymentC.CreateFakeSubMerchant(user, cook.BTSubMerchantID)
			if err != nil {
				resp.Err = errors.Wrap("failed to paymentC.CreateFakeSubMerchant", err)
				return resp, nil
			}
		}
	}

	// set onboard to true
	user.SetOnboard(true)
	err = auth.SaveUser(ctx, user)
	if err != nil {
		resp.Err = errors.Wrap("failed to set user to onboarded", err)
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

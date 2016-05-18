package gigachef

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/application"
	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// SubmitApplicationReq is the input request needed for SubmitApplication.
type SubmitApplicationReq struct {
	Gigatoken   string                      `json:"gigatoken"`
	Application application.ChefApplication `json:"application"`
}

// gigatoken returns the Gigatoken string
func (req *SubmitApplicationReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *SubmitApplicationReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	if req.Application.Email == "" {
		return fmt.Errorf("Application email is empty")
	}
	return nil
}

// SubmitApplicationResp is the output response for SubmitApplication.
type SubmitApplicationResp struct {
	Application application.ChefApplication `json:"application"`
	Err         errors.ErrorWithCode        `json:"err"`
}

// SubmitApplication is an endpoint that submits or updates a chef application.
func (service *Service) SubmitApplication(ctx context.Context, req *SubmitApplicationReq) (*SubmitApplicationResp, error) {
	resp := new(SubmitApplicationResp)
	defer handleResp(ctx, "SubmitApplication", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if req.Application.Address.Country == "" {
		req.Application.Address.Country = "USA"
	}
	chefApplication, err := application.SubmitApplication(ctx, user, &req.Application)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Application = *chefApplication
	return resp, nil
}

// GetApplicationReq is the input request needed for GetApplication.
type GetApplicationReq struct {
	Gigatoken string `json:"gigatoken"`
}

// gigatoken returns the Gigatoken string
func (req *GetApplicationReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *GetApplicationReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// GetApplicationResp is the output response for GetApplication.
type GetApplicationResp struct {
	Application application.ChefApplication `json:"application"`
	Err         errors.ErrorWithCode        `json:"err"`
}

// GetApplication is an endpoint that gets a chef application.
func (service *Service) GetApplication(ctx context.Context, req *GetApplicationReq) (*GetApplicationResp, error) {
	resp := new(GetApplicationResp)
	defer handleResp(ctx, "GetApplication", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	chefApplication, err := application.GetApplication(ctx, user)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Application = *chefApplication
	return resp, nil
}

// UpdateSubMerchantReq updates sub-merchant payment info
type UpdateSubMerchantReq struct {
	Gigatoken     string        `json:"gigatoken"`
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	DateOfBirth   string        `json:"date_of_birth"`
	AccountNumber string        `json:"account_number"`
	RoutingNumber string        `json:"routing_number"`
	Address       types.Address `json:"address"`
}

// gigatoken returns the Gigatoken string
func (req *UpdateSubMerchantReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *UpdateSubMerchantReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// UpdateSubMerchant creates or updates sub-merchant info
func (service *Service) UpdateSubMerchant(ctx context.Context, req *UpdateSubMerchantReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "UpdateSubMerchant", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	chefC := gigachef.New(ctx)
	chef, err := chefC.Get(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot get chef(%d)", user.ID)
		return resp, nil
	}
	if req.Address.Country == "" {
		req.Address.Country = "USA"
	}
	updateMerchantReq := &payment.UpdateSubMerchantReq{
		User:          *user,
		ID:            chef.BTSubMerchantID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         user.Email,
		DateOfBirth:   req.DateOfBirth,
		AccountNumber: req.AccountNumber,
		RoutingNumber: req.RoutingNumber,
		Address:       req.Address,
	}
	paymentC := payment.New(ctx)
	_, err = paymentC.UpdateSubMerchant(updateMerchantReq)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot update sub-merchant(%d)", chef.BTSubMerchantID)
		return resp, nil
	}
	return resp, nil
}

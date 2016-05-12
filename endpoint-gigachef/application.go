package gigachef

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/application"
	"github.com/atishpatel/Gigamunch-Backend/errors"
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

package gigachef

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/application"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

// SubmitApplicationReq is the input request needed for SubmitApplication.
type SubmitApplicationReq struct {
	GigaToken              string        `json:"gigatoken"`
	Name                   string        `json:"name"`
	Email                  string        `json:"email"`
	PhoneNumber            string        `json:"phone_number"`
	Address                types.Address `json:"address"`
	AttendedCulinarySchool bool          `json:"attended_culinary_school"`
	WorkedAtResturant      bool          `json:"worked_at_resturant"`
	PostFrequency          int           `json:"post_frequency"`
}

// Gigatoken returns the GigaToken string
func (req *SubmitApplicationReq) Gigatoken() string {
	return req.GigaToken
}

// Valid validates a req
func (req *SubmitApplicationReq) Valid() error {
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty.")
	}
	if req.Email == "" {
		return fmt.Errorf("Email is incorrect.")
	}
	if len(req.PhoneNumber) < 7 {
		return fmt.Errorf("Phone number is invalid.")
	}
	return nil
}

// SubmitApplicationResp is the output response for SubmitApplication.
type SubmitApplicationResp struct {
	ApplicationProgress int                  `json:"application_progress"`
	Err                 errors.ErrorWithCode `json:"err"`
}

// SubmitApplication is an endpoint that submits or updates a chef application.
func (service *Service) SubmitApplication(ctx context.Context, req *SubmitApplicationReq) (*SubmitApplicationResp, error) {
	resp := new(SubmitApplicationResp)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	chefApplication := &application.ChefApplication{
		Name:                   req.Name,
		Email:                  req.Email,
		PhoneNumber:            req.PhoneNumber,
		AttendedCulinarySchool: req.AttendedCulinarySchool,
		Address:                req.Address,
		WorkedAtResturant:      req.WorkedAtResturant,
		PostFrequency:          req.PostFrequency,
	}
	chefApplication, err = application.SubmitApplication(ctx, user, chefApplication)
	if err != nil {
		utils.Debugf(ctx, "err: %+v", err)
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.ApplicationProgress = chefApplication.ApplicationProgress
	return resp, nil
}

// GetApplicationReq is the input request needed for GetApplication.
type GetApplicationReq struct {
	GigaToken string `json:"gigatoken"`
}

// Gigatoken returns the GigaToken string
func (req *GetApplicationReq) Gigatoken() string {
	return req.GigaToken
}

// Valid validates a req
func (req *GetApplicationReq) Valid() error {
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty.")
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
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	chefApplication, err := application.GetApplication(ctx, user)
	if err != nil {
		utils.Debugf(ctx, "err: %+v", err)
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Application = *chefApplication
	return resp, nil
}

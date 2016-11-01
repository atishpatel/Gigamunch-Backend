package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

type coder interface {
	GetCode() int32
}

func handleResp(ctx context.Context, fnName string, resp coder) {
	code := resp.GetCode()
	if code == errors.CodeInvalidParameter {
		utils.Warningf(ctx, "%s invalid parameter: %+v", fnName, resp)
		return
	} else if code != 0 {
		utils.Errorf(ctx, "%s err: %+v", fnName, resp)
	}
}

type validatableTokenReq interface {
	gigatoken() string
	valid() error
}

func validateRequestAndGetUser(ctx context.Context, req validatableTokenReq) (*types.User, error) {
	err := req.valid()
	if err != nil {
		return nil, errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}.Wrap("failed to validate request")
	}
	user, err := auth.GetUserFromToken(ctx, req.gigatoken())
	return user, err
}

// Service is the REST API Endpoint exposed to Gigamunchers
type Service struct{}

func init() {
	api, err := endpoints.RegisterService(&Service{}, "cookservice", "v1", "An endpoint service for cooks.", true)
	if err != nil {
		log.Fatalf("Failed to register service: %#v", err)
	}

	register := func(orig, name, method, path, desc string) {
		m := api.MethodByName(orig)
		if m == nil {
			log.Fatalf("Missing method %s", orig)
		}
		i := m.Info()
		i.Name, i.HTTPMethod, i.Path, i.Desc = name, method, path, desc
	}
	// Page stuff
	register("FinishOnboarding", "finishOnboarding", "POST", "cookservice/finishOnboarding", "Updates cook and submerchant")
	// Refresh stuff
	register("RefreshToken", "refreshToken", "POST", "cookservice/refreshToken", "Refresh a token.")
	// Cook stuff
	register("GetCook", "getCook", "GET", "cookservice/getCook", "Get the cook info.")
	register("UpdateCook", "updateCook", "POST", "cookservice/updateCook", "Update cook information.")
	// Submerchant stuff
	register("UpdateSubMerchant", "updateSubMerchant", "POST", "gigachefservice/updateSubMerchant", "Update or create sub-merchant.")
	register("GetSubMerchant", "getSubMerchant", "GET", "gigachefservice/getSubMerchant", "Get the sub merchant info.")
	// Item stuff
	register("SaveItem", "saveItem", "POST", "cookservice/saveItem", "Save an item.")
	register("GetItem", "getItem", "GET", "cookservice/getItem", "Get an item.")
	register("ActivateItem", "activateItem", "POST", "cookservice/activateItem", "Activate an item.")
	register("DeactivateItem", "deactivateItem", "POST", "cookservice/deactivateItem", "Deactivate an item.")
	// Menu stuff
	register("GetMenus", "getMenus", "GET", "cookservice/getMenus", "Gets the menus for a cook.")
	register("SaveMenu", "saveMenu", "POST", "cookservice/saveMenu", "Save a menu.")
	// Inquiry stuffffffffff
	register("GetMessageToken", "getMessageToken", "GET", "cookservice/getMessageToken", "Gets the a token for messaging.")
	register("GetInquiries", "GetInquiries", "GET", "cookservice/GetInquiries", "GetInquiries gets a cook's inquiries.")
	register("GetInquiry", "GetInquiry", "GET", "cookservice/GetInquiry", "GetInquiry gets a cook's inquiry.")
	register("AcceptInquiry", "AcceptInquiry", "POST", "cookservice/acceptInquiry", "AcceptInquiry accepts an inquiry for a cook.")
	register("DeclineInquiry", "DeclineInquiry", "POST", "cookservice/declineInquiry", "DeclineInquiry declines an inquiry for a cook.")
	endpoints.HandleHTTP()
}

package cookapi

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
	GetCode() int
}

func handleResp(ctx context.Context, fnName string, resp coder) {
	code := resp.GetCode()
	if code != 0 && code != errors.CodeInvalidParameter {
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
	// refresh stuff
	register("RefreshToken", "refreshToken", "POST", "cookservice/refreshToken", "Refresh a token.")
	// cook stuff
	register("GetCook", "getCook", "GET", "cookservice/getCook", "Get the cook info.")
	// register("UpdateProfile", "updateProfile", "POST", "gigachefservice/updateProfile", "Update chef profile.")
	// register("UpdateSubMerchant", "updateSubMerchant", "POST", "gigachefservice/updateSubMerchant", "Update or create sub-merchant.")
	// register("GetSubMerchant", "getSubMerchant", "GET", "gigachefservice/getSubMerchant", "Get the sub merchant info.")
	// // item stuff
	register("SaveItem", "saveItem", "POST", "cookservice/saveItem", "Save a item.")
	// register("GetItem", "getItem", "GET", "gigachefservice/getItem", "Get a item.")
	register("GetMenus", "getMenus", "GET", "cookservice/getMenus", "Gets the menus for a cook.")
	// // post stuff
	// register("PublishPost", "publishPost", "POST", "gigachefservice/publishPost", "Publish a post.")
	// register("GetPosts", "getPosts", "GET", "gigachefservice/getPosts", "Get a chef's posts.")

	endpoints.HandleHTTP()
}

package gigachef

import (
	"log"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"gitlab.com/atishpatel/Gigamunch-Backend/auth"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"
)

func itot(i int) time.Time {
	return time.Unix(int64(i), 0)
}

func ttoi(t time.Time) int {
	return int(t.Unix())
}

func itos(i int64) string {
	return strconv.FormatInt(i, 10)
}

func stoi(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

type coder interface {
	GetCode() int
}

func handleResp(ctx context.Context, fnName string, resp coder) {
	code := resp.GetCode()
	if code != 0 && code != errors.CodeInvalidParameter {
		utils.Errorf(ctx, "%s err: ", fnName, resp)
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
	api, err := endpoints.RegisterService(&Service{}, "gigachefservice", "v1", "An endpoint service for Gigachefs.", true)
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
	// TODO update to GET
	// refresh stuff
	register("RefreshToken", "refreshToken", "POST", "gigachefservice/refreshToken", "Refresh a token.")
	// application stuff
	register("GetGigachef", "getGigachef", "GET", "gigachefservice/getGigachef", "Get the chef info.")
	register("UpdateProfile", "updateProfile", "POST", "gigachefservice/updateProfile", "Update chef profile.")
	register("UpdateSubMerchant", "updateSubMerchant", "POST", "gigachefservice/updateSubMerchant", "Update or create sub-merchant.")
	// item stuff
	register("SaveItem", "saveItem", "POST", "gigachefservice/saveItem", "Save a item.")
	register("GetItem", "getItem", "POST", "gigachefservice/getItem", "Get a item.")
	register("GetItems", "getItems", "POST", "gigachefservice/getItems", "Gets items sorted by last used.")
	// post stuff
	register("PostPost", "postPost", "POST", "gigachefservice/postPost", "Post a post.")
	endpoints.HandleHTTP()
}

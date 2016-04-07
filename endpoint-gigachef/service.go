package gigachef

import (
	"log"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

type validatableTokenReq interface {
	Gigatoken() string
	Valid() error
}

func validateRequestAndGetUser(ctx context.Context, req validatableTokenReq) (*types.User, error) {
	err := req.Valid()
	if err != nil {
		return nil, errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
	}
	user, err := auth.GetUserFromToken(ctx, req.Gigatoken())
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
	// Register course stuff
	register("GetApplication", "getApplication", "POST", "gigachefservice/getApplication", "Get chef application.")
	register("SubmitApplication", "submitApplication", "POST", "gigachefservice/submitApplication", "Apply to be a chef.")
	register("SaveItem", "saveItem", "POST", "gigachefservice/saveItem", "Save a item.")
	register("GetItem", "getItem", "POST", "gigachefservice/getItem", "Get a item.")
	register("GetItems", "getItems", "POST", "gigachefservice/getItems", "Gets items sorted by last used.")
	register("PostPost", "postPost", "POST", "gigachefservice/postPost", "Post a post.")
	endpoints.HandleHTTP()
}

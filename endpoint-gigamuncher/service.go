package gigamuncher

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
	api, err := endpoints.RegisterService(&Service{}, "gigamuncherservice", "v1", "An endpoint service for Gigamunchers.", true)
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
	register("SignIn", "signIn", "POST", "gigamuncherservice/signIn", "Sign in a user using a gtoken.")
	register("SignOut", "signOut", "POST", "gigamuncherservice/signOut", "Sign out a user.")
	register("RefreshToken", "refreshToken", "POST", "gigamuncherservice/refreshToken", "Refresh a token.")

	register("GetPost", "getPost", "POST", "gigamuncherservice/getPost", "Get post details.")
	register("GetLivePosts", "getLivePosts", "POST", "gigamuncherservice/getLivePosts", "Get live posts.")

	register("PostReview", "postReview", "POST", "gigamuncherservice/postReview", "Post a review.")
	register("GetReviews", "getReviews", "POST", "gigamuncherservice/getReviews", "Get reviews.")
	endpoints.HandleHTTP()
}

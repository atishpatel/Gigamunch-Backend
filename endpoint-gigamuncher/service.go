package gigamuncher

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

func ttoi(t time.Time) int {
	return int(t.Unix())
}

func itojn(i int64) json.Number {
	return json.Number(strconv.FormatInt(i, 10))
}

func itos(i int64) string {
	return strconv.FormatInt(i, 10)
}

func stoi(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ftos64(f float64) string {
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func stof64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

type validatableTokenReq interface {
	gigatoken() string
	valid() error
}

func validateRequestAndGetUser(ctx context.Context, req validatableTokenReq) (*types.User, error) {
	err := req.valid()
	if err != nil {
		return nil, errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
	}
	user, err := auth.GetUserFromToken(ctx, req.gigatoken())
	return user, err
}

type coder interface {
	GetCode() int32
}

func handleResp(ctx context.Context, fnName string, resp coder) {
	code := resp.GetCode()
	if code != 0 && code != errors.CodeInvalidParameter {
		utils.Errorf(ctx, "%s err: ", fnName, resp)
	}
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
	// login
	register("SignIn", "signIn", "POST", "gigamuncherservice/signIn", "Sign in a user using a gtoken.")
	register("SignOut", "signOut", "POST", "gigamuncherservice/signOut", "Sign out a user.")
	register("RefreshToken", "refreshToken", "POST", "gigamuncherservice/refreshToken", "Refresh a token.")
	// address
	register("GetAddresses", "getAddresses", "GET", "gigamuncherservice/getAddresses", "Get the muncher's addresses.")
	register("SelectAddress", "selectAddress", "POST", "gigamuncherservice/selectAddress", "Select an Address.")
	// post
	register("GetPost", "getPost", "GET", "gigamuncherservice/getPost", "Get post details.")
	register("GetLivePosts", "getLivePosts", "GET", "gigamuncherservice/getLivePosts", "Get live posts.")
	// like
	register("LikeItem", "likeItem", "POST", "gigamuncherservice/likeItem", "Like an item.")
	register("UnlikeItem", "unlikeItem", "POST", "gigamuncherservice/unlikeItem", "Unlike an item.")
	// order
	register("GetBraintreeToken", "getBraintreeToken", "GET", "gigamuncherservice/getBraintreeToken", "Get a braintreeToken.")
	register("MakeOrder", "makeOrder", "POST", "gigamuncherservice/makeOrder", "Make an order.")
	register("CancelOrder", "cancelOrder", "POST", "gigamuncherservice/cancelOrder", "Cancel an order.")
	register("GetOrders", "getOrders", "GET", "gigamuncherservice/getOrders", "Gets the orders for a muncher.")
	register("GetOrder", "getOrder", "GET", "gigamuncherservice/getOrder", "Get an order.")
	// review
	register("PostReview", "postReview", "POST", "gigamuncherservice/postReview", "Post a review.")
	register("GetReviews", "getReviews", "GET", "gigamuncherservice/getReviews", "Get reviews.")
	endpoints.HandleHTTP()
}

package gigamuncher

import (
	"log"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
)

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
	// register("GetLiveMeals", "getLiveMeals", "POST", "gigamuncherservice/getLiveMeals", "Get live meals")
	register("SignIn", "signIn", "POST", "gigamuncherservice/signIn", "Sign in a user using a gtoken.")
	register("SignOut", "signOut", "POST", "gigamuncherservice/signOut", "Sign out a user.")
	register("RefreshToken", "refreshToken", "POST", "gigamuncherservice/refreshToken", "Refresh a token.")
	register("PostReview", "postReview", "POST", "gigamuncherservice/postReview", "Post a review.")
	register("GetReviews", "getReviews", "POST", "gigamuncherservice/getReviews", "Get reviews.")
	endpoints.HandleHTTP()
}

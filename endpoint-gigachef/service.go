package gigachef

import (
	"log"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
)

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
	// register("PostMeal", "postMeal", "POST", "gigachefservice/postMeal", "Post a meal")
	endpoints.HandleHTTP()
}

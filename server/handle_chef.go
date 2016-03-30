package server

import (
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	chefIndexPage []byte
)

func handleGigachefApp(w http.ResponseWriter, req *http.Request) {
	var err error
	var page []byte
	if appengine.IsDevAppServer() {
		page, err = ioutil.ReadFile("chef/app/index.html")
	} else {
		page = chefIndexPage
	}

	if err != nil {
		ctx := appengine.NewContext(req)
		utils.Errorf(ctx, "Error reading login page: %+v", err)
	}
	w.Write(page)
}

func middlewareLoggedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user := CurrentUser(w, req)
		if user == nil {
			http.Redirect(w, req, loginURL, http.StatusTemporaryRedirect)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func init() {
	var err error
	// TODO switch to template with footer and stuff in different page
	chefIndexPage, err = ioutil.ReadFile("chef/app/index.html")
	if err != nil {
		log.Fatal("chef/app/index.html not found")
	}
}

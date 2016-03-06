package server

import (
	"io/ioutil"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine"
)

func handleChefApplication(w http.ResponseWriter, req *http.Request) {
	page, err := ioutil.ReadFile("app/application.html")
	if err != nil {
		ctx := appengine.NewContext(req)
		utils.Errorf(ctx, "Error reading login page: %+v", err)
	}
	w.Write(page)
}

func middlewareVerifiedChef(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authToken := CurrentUser(w, req)
		if authToken == nil {
			http.Redirect(w, req, loginURL, http.StatusTemporaryRedirect)
			return
		} else if !authToken.User.IsVerifiedChef() {
			http.Redirect(w, req, chefApplicationURL, http.StatusTemporaryRedirect)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func handleChefHome(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("chef home"))
}

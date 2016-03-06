package server

import (
	// golang

	"io/ioutil"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/utils"

	// appengine
	"google.golang.org/appengine"

	"github.com/justinas/alice"
)

func init() {
	http.HandleFunc(baseLoginURL, handleLogin)
	http.HandleFunc(signOutURL, handleSignout)

	// chef stuff
	http.HandleFunc(chefApplicationURL, handleChefApplication)
	// verfied chef stuff
	verifiedChefChain := alice.New(middlewareVerifiedChef)
	http.Handle(chefHomeURL, verifiedChefChain.ThenFunc(handleChefHome))
	// admin stuff
	adminChain := alice.New(middlewareAdmin)
	http.Handle(adminHomeURL, adminChain.ThenFunc(handleAdminHome))
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	authToken := CurrentUser(w, req)
	if authToken != nil {
		utils.Debugf(appengine.NewContext(req), "User: %+v", authToken.User)
		if authToken.User.IsVerifiedChef() {
			http.Redirect(w, req, chefHomeURL, http.StatusTemporaryRedirect)
		} else if authToken.User.IsAdmin() {
			http.Redirect(w, req, adminHomeURL, http.StatusTemporaryRedirect)
		} else {
			http.Redirect(w, req, chefApplicationURL, http.StatusTemporaryRedirect)
		}
	}
	removeCookies(w)
	page, err := ioutil.ReadFile("app/login.html")
	if err != nil {
		ctx := appengine.NewContext(req)
		utils.Errorf(ctx, "Error reading login page: %+v", err)
	}
	w.Write(page)
}

func handleSignout(w http.ResponseWriter, req *http.Request) {
	removeCookies(w)
	http.Redirect(w, req, homeURL, http.StatusTemporaryRedirect)
}

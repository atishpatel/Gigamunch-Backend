package server

import (
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/utils"
)

func init() {
	r := httprouter.New()

	// chef stuff
	// loggedInChain := alice.New(middlewareLoggedIn)
	// r.Handler("GET", baseGigachefURL+"/*path", loggedInChain.ThenFunc(handleGigachefApp))

	r.GET(baseLoginURL, handleLogin)
	r.GET(signOutURL, handleSignout)
	r.POST("/upload", handleUpload)

	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	r.NotFound = http.HandlerFunc(handle404)
	http.Handle("/", r)
}

func handle404(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("GIGA 404 page. :()"))
}

func handleLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	user := CurrentUser(w, req)
	if user != nil {
		http.Redirect(w, req, baseGigachefURL, http.StatusTemporaryRedirect)
	}
	removeCookies(w)
	page, err := ioutil.ReadFile("app/login.html")
	if err != nil {
		utils.Errorf(ctx, "Error reading login page: %+v", err)
	}
	w.Write(page)
}

func handleSignout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	removeCookies(w)
	http.Redirect(w, req, homeURL, http.StatusTemporaryRedirect)
}

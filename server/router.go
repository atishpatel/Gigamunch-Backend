package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := httprouter.New()

	r.GET(baseLoginURL, handleLogin)
	r.GET(signOutURL, handleSignout)

	// r.POST("/upload", handleUpload)
	r.GET("/get-upload-url", hangleGetUploadURL)

	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	r.NotFound = http.HandlerFunc(handle404)
	http.HandleFunc("/upload", handleUpload)
	http.Handle("/", r)
}

func handle404(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("GIGA 404 page. :()"))
}

func handleLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := CurrentUser(w, req)
	if user != nil {
		http.Redirect(w, req, baseGigachefURL, http.StatusTemporaryRedirect)
	}
	removeCookies(w)
	http.Redirect(w, req, "/becomechef", http.StatusTemporaryRedirect)
}

func handleSignout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	removeCookies(w)
	http.Redirect(w, req, homeURL, http.StatusTemporaryRedirect)
}

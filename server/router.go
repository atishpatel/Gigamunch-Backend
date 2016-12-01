package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := httprouter.New()

	r.GET(baseLoginURL, handleLogin)
	r.GET(signOutURL, handleSignout)

	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	r.NotFound = http.HandlerFunc(handle404)
	http.HandleFunc("/get-upload-url", handleGetUploadURL)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/get-feed", handleGetFeed)
	http.HandleFunc("/get-item", handleGetItem)
	http.Handle("/", r)
}

func handle404(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("GIGA 404 page. :()"))
}

func handleLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := CurrentUser(w, req)
	if user != nil {
		http.Redirect(w, req, baseCookURL, http.StatusTemporaryRedirect)
		return
	}
	removeCookies(w)
	http.Redirect(w, req, "/becomechef", http.StatusTemporaryRedirect)
}

func handleSignout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	removeCookies(w)
	http.Redirect(w, req, homeURL, http.StatusTemporaryRedirect)
}

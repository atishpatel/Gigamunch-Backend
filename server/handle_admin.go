package server

import "net/http"

func middlewareAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authToken := CurrentUser(w, req)
		if authToken == nil {
			http.Redirect(w, req, loginURL, http.StatusTemporaryRedirect)
			return
		} else if !authToken.User.IsAdmin() {
			http.Redirect(w, req, chefApplicationURL, http.StatusTemporaryRedirect)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func handleAdminHome(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("admin home"))
}

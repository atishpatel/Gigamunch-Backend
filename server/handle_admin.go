package server

import "net/http"

func middlewareAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user := CurrentUser(w, req)
		if user == nil {
			http.Redirect(w, req, loginURL, http.StatusTemporaryRedirect)
			return
		} else if !user.IsAdmin() {
			http.Redirect(w, req, baseGigachefURL, http.StatusTemporaryRedirect)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func handleAdminHome(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("admin home"))
}

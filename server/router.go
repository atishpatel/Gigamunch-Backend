package server

import (
	// golang

	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"

	// appengine
	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc(types.LoginURL, handleLogin)
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	user, err := auth.CurrentUser(w, req)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	utils.Infof(ctx, "Got the following user: %+v", user)
	w.Write([]byte("Token generated"))
}

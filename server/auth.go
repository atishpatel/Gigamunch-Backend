package server

import (
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine"
)

func removeCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   sessionTokenCookieName,
		MaxAge: -1,
		Secure: true,
	})
	http.SetCookie(w, &http.Cookie{Name: gitkitCookieName, MaxAge: -1, Secure: true})
}

func saveAuthCookie(w http.ResponseWriter, authToken *types.AuthToken) {
	http.SetCookie(w, &http.Cookie{
		Name:   sessionTokenCookieName,
		Value:  authToken.JWTString,
		MaxAge: int(authToken.Expire.Sub(time.Now()).Seconds()),
		Secure: true,
	})
}

// CurrentUser extracts the user information stored in current session.
//
// If there is no existing session, identity toolkit token is checked. If the
// token is valid, an AuthToken is returned
//
// If any error happens, nil is returned.
// TODO: When to refresh?
func CurrentUser(w http.ResponseWriter, req *http.Request) *types.AuthToken {
	ctx := appengine.NewContext(req)
	var err error
	sessionTokenCookie, err := req.Cookie(sessionTokenCookieName)
	// doesn't have a session cookie
	if err == http.ErrNoCookie {
		gTokenCookie, err := req.Cookie(gitkitCookieName)
		if err != nil || gTokenCookie.Value == "" {
			return nil
		}
		user, authToken, err := auth.GetSessionWithGToken(ctx, gTokenCookie.Value)
		if err != nil {
			errWithCode := errors.GetErrorWithCode(err)
			if errWithCode.Code == errors.CodeSignOut {
				removeCookies(w)
			} else {
				utils.Errorf(ctx, "Error getting user form gtoken: %+v", err)
			}
			return nil
		}
		saveAuthCookie(w, authToken)
		http.SetCookie(w, &http.Cookie{Name: gitkitCookieName, MaxAge: -1})
		return authToken
	} else if err != nil {
		utils.Errorf(ctx, "Error getting session cookie: %+v", err)
		return nil
	}
	// has session cookie
	authToken := &types.AuthToken{
		JWTString: sessionTokenCookie.Value,
	}
	authToken, err = auth.GetUserFromToken(ctx, sessionTokenCookie.Value)
	if err != nil {
		errWithCode, ok := err.(errors.ErrorWithCode)
		if ok && (errWithCode.Code == errors.CodeSignOut) {
			removeCookies(w)
		} else {
			utils.Errorf(ctx, "Error getting user from session: %+v", err)
		}
		return nil
	}
	saveAuthCookie(w, authToken)
	return authToken
}

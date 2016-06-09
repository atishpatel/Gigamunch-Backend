package server

import (
	"net/http"
	"time"

	"gitlab.com/atishpatel/Gigamunch-Backend/auth"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/account"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"

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

func saveAuthCookie(w http.ResponseWriter, authToken string) {
	cookie := &http.Cookie{
		Name:   sessionTokenCookieName,
		Value:  authToken,
		MaxAge: int(auth.GetExpTime().Sub(time.Now()).Seconds()),
		Secure: true,
	}
	if appengine.IsDevAppServer() {
		cookie.Secure = false
	}
	http.SetCookie(w, cookie)
}

func getTokenFromReq(req *http.Request) (string, error) {
	cookie, err := req.Cookie(sessionTokenCookieName)
	if err != nil {
		return "", err
	}
	token := cookie.Value
	if token == "" {
		return "", errInvalidParameter.WithMessage("Token is empty.")
	}
	return token, nil
}

func getUserFromCookie(req *http.Request) (*types.User, error) {
	token, err := getTokenFromReq(req)
	if err != nil {
		return nil, err
	}
	ctx := appengine.NewContext(req)
	// get user
	return auth.GetUserFromToken(ctx, token)
}

// CurrentUser extracts the user information stored in current session.
//
// If there is no existing session, identity toolkit token is checked. If the
// token is valid, an AuthToken is returned
//
// If any error happens, nil is returned.
// TODO: When to refresh?
func CurrentUser(w http.ResponseWriter, req *http.Request) *types.User {
	ctx := appengine.NewContext(req)
	var err error
	sessionTokenCookie, err := req.Cookie(sessionTokenCookieName)
	// doesn't have a session cookie
	if err == http.ErrNoCookie {
		var gTokenCookie *http.Cookie
		gTokenCookie, err = req.Cookie(gitkitCookieName)

		if err != nil || gTokenCookie.Value == "" {
			return nil
		}
		var user *types.User
		var authToken string
		user, authToken, err = auth.GetSessionWithGToken(ctx, gTokenCookie.Value)
		if err != nil {
			errWithCode := errors.GetErrorWithCode(err)
			if errWithCode.Code == errors.CodeSignOut {
				removeCookies(w)
			} else {
				utils.Errorf(ctx, "Error getting user form gtoken: %+v", err)
			}
			return nil
		}
		// if user isn't a chef make them a chef
		if !user.IsChef() {
			user.SetChef(true)
			err = auth.SaveUser(ctx, user)
			if err != nil {
				utils.Errorf(ctx, "Error auth.SaveUser", err)
				// TODO redirect?
			}
			err = account.SaveUserInfo(ctx, user, nil, "")
			if err != nil {
				utils.Errorf(ctx, "Error account.SavingUserInfo", err)
				return nil
			}
			authToken, err = auth.RefreshToken(ctx, authToken)
			if err != nil {
				utils.Errorf(ctx, "failed to refresh token ", err)
				return nil
			}
		}
		saveAuthCookie(w, authToken)
		gTokenCookie.MaxAge = 120
		http.SetCookie(w, gTokenCookie)
		return user
	} else if err != nil {
		utils.Errorf(ctx, "Error getting session cookie: %+v", err)
		return nil
	}
	// has session cookie
	user, err := auth.GetUserFromToken(ctx, sessionTokenCookie.Value)
	if err != nil {
		errWithCode, ok := err.(errors.ErrorWithCode)
		if ok && (errWithCode.Code == errors.CodeSignOut) {
			removeCookies(w)
		} else {
			utils.Errorf(ctx, "Error getting user from session: %+v", err)
		}
		return nil
	}
	// saveAuthCookie(w, authToken)
	return user
}

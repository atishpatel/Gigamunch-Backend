package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/identity-toolkit-go-client/gitkit"
	"golang.org/x/net/context"

	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/databases/session"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	gitkitClient    *gitkit.Client
	gitkitAudiences []string
)

// GetUserFromSessionID returns channel that returns types.User if user session is found
// otherwise, it return nil
func GetUserFromSessionID(ctx context.Context, sessionID string) <-chan *types.User {
	//TODO fix with redis stuff, check if sessionID is valid
	return session.GetUserSession(ctx, sessionID)
}

// GetUserAndSessionFromGToken returns a user with a valid sessionID if the gtoken is valid
// otherwise, it returns nil, ""
func GetUserAndSessionFromGToken(ctx context.Context, tokenString string) (string, *types.User) {
	if tokenString == "" {
		return "", nil
	}
	token, err := gitkitClient.ValidateToken(ctx, tokenString, gitkitAudiences)
	if err != nil {
		utils.Errorf(ctx, "Invalid gitkit token %s: %+v", tokenString, err)
		return "", nil
	}
	if time.Now().Sub(token.IssueAt) > 15*time.Minute {
		utils.Infof(ctx, "Token %s is too old. Issued at: %s", tokenString, token.IssueAt)
		return "", nil
	}
	// get user info from gitkit servers
	gitkitUser, err := gitkitClient.UserByLocalID(ctx, token.LocalID)
	if err != nil {
		utils.Errorf(ctx, "Failed to fetch user info from gitkit servers %s(%s): %+v", token.Email, token.LocalID, err)
		return "", nil
	}

	sessionID, user := createSession(ctx, gitkitUser)

	return sessionID, user
}

func createSession(ctx context.Context, gitkitUser *gitkit.User) (string, *types.User) {
	// only set name like that if it's the first time the user in logging in
	user := &types.User{
		Email: gitkitUser.Email,
	}

	// TODO(Atish): fix setting permissions for session
	// chefChan := getBasicChefInfo(u.Email)
	// muncherChan := getBasicMuncherInfo(u.Email)

	// TODO(Atish): If first time logging in, download PhotoURL to server

	sessionID := utils.GetUUID()
	errChan := SaveUserSession(ctx, sessionID, user)
	err := <-errChan
	if err != nil {
		utils.Errorf(ctx, "Failed to save user session", err)
		return "", nil
	}
	return sessionID, user
}

// SaveUserSession saves the tpyes.User value with the sessionID key
func SaveUserSession(ctx context.Context, sessionID string, user *types.User) <-chan error {
	return session.SaveUserSession(ctx, sessionID, user)
}

// CurrentUser extracts the user information stored in current session.
//
// If there is no existing session, identity toolkit token is checked. If the
// token is valid, a new session is created.
//
// If any error happens, nil is returned.
func CurrentUser(req *http.Request, w http.ResponseWriter) *types.User {
	ctx := appengine.NewContext(req)
	sessionCookie, err := req.Cookie(types.SessionCookieName)
	// doesn't have a session cookie
	if err == http.ErrNoCookie {
		ts := gitkitClient.TokenFromRequest(req)
		if ts == "" {
			return nil
		}
		sessionID, user := GetUserAndSessionFromGToken(ctx, ts)
		if user != nil {
			http.SetCookie(w, &http.Cookie{Name: types.SessionCookieName, Value: sessionID, MaxAge: 86400 * 100})
			http.SetCookie(w, &http.Cookie{Name: types.GitkitCookieName, MaxAge: 1})
		}
		return user
	} else if err != nil {
		utils.Errorf(ctx, "Error getting session cookie: %+v", err)
	}
	sessionID := sessionCookie.Value
	userChan := GetUserFromSessionID(ctx, sessionID)
	return <-userChan
}

func init() {
	config := config.GetGitkitConfig()
	// setup gitkit
	c := &gitkit.Config{
		WidgetURL: types.GitkitURL,
	}
	if appengine.IsDevAppServer() {
		c.GoogleAppCredentialsPath = config.GoogleAppCredentialsPath
	}
	var err error
	gitkitClient, err = gitkit.New(context.Background(), c)
	if err != nil {
		log.Fatal(err)
	}
	gitkitAudiences = []string{config.ClientID}
}

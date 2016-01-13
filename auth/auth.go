package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/identity-toolkit-go-client/gitkit"
	"golang.org/x/net/context"

	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	sessionName = "GIGASID"
)

var (
	gitkitClient *gitkit.Client
	serverConfig *config.Config
)

// GetUserFromSessionID returns a User object if user session is found
// otherwise, it return nil
func GetUserFromSessionID(sessionID string) *types.User {
	//TODO fix with redis stuff, check if sessionID is valid
	return nil
}

func GetUserAndSessionFromGToken(ctx context.Context, tokenString string) (*types.User, string) {
	if tokenString == "" {
		return nil, ""
	}
	token, err := gitkitClient.ValidateToken(ctx, tokenString, []string{serverConfig.ClientID})
	if err != nil {
		utils.Errorf(ctx, "Invalid gitkit token %s: %+v", tokenString, err)
		return nil, ""
	}
	if time.Now().Sub(token.IssueAt) > 15*time.Minute {
		utils.Infof(ctx, "Token %s is too old. Issued at: %s", tokenString, token.IssueAt)
		return nil, ""
	}
	// get user info from gitkit servers
	//u, err := gitkitClient.UserByLocalID(ctx, token.LocalID)
	if err != nil {
		utils.Errorf(ctx, "Failed to fetch user info from gitkit servers %s(%s): %+v", token.Email, token.LocalID, err)
		return nil, ""
	}

	// TODO(Atish): check if can use getmulti on redis
	// chefChan := getBasicChefInfo(u.Email)
	// muncherChan := getBasicMuncherInfo(u.Email)

	sessionID := CreateSession(ctx)

	return nil, sessionID
}

func createSession(ctx context.Context) string {
	//TODO: Probably takes in a user and returns a sessionID
	return utils.GetUUID()
}

// CurrentUser extracts the user information stored in current session.
//
// If there is no existing session, identity toolkit token is checked. If the
// token is valid, a new session is created.
//
// If any error happens, nil is returned.
func CurrentUser(req *http.Request, w http.ResponseWriter) *types.User {
	ctx := appengine.NewContext(req)
	sessionCookie, err := req.Cookie(sessionName)
	// doesn't have a session cookie
	if err == http.ErrNoCookie {
		ts := gitkitClient.TokenFromRequest(req)
		if ts == "" {
			return nil
		}
		user, sessionID := getUserFromGToken(ctx, ts)
		// TODO fix
		utils.Infof(ctx, "tmp to help golang build", sessionID)
		return user
	} else if err != nil {
		utils.Errorf(ctx, "Error getting user: %+v", err)
	}
	sessionID := sessionCookie.Value
	return GetUserFromSessionID(sessionID)
}

func init() {
	serverConfig = config.GetConfig()
	// setup gitkit
	c := &gitkit.Config{
		WidgetURL: types.GitkitURL,
	}
	if appengine.IsDevAppServer() {
		c.GoogleAppCredentialsPath = serverConfig.GoogleAppCredentialsPath
	}
	var err error
	gitkitClient, err = gitkit.New(context.Background(), c)
	if err != nil {
		log.Fatal(err)
	}

	// databaseConfig := databases.Config{
	// 	RedisSessionDBIP:       config.RedisSessionServerIP,
	// 	RedisSessionDBPassword: config.RedisSessionServerPassword
	//	RedisSessionDBPoolSize: 10}
	// db := databases.CreateDatabase(databaseConfig)

}

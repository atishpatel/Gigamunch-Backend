package auth

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/google/identity-toolkit-go-client/gitkit"
	"golang.org/x/net/context"
	jwt "gopkg.in/dgrijalva/jwt-go.v2"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/core/account"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	gitkitClient    *gitkit.Client
	gitkitAudiences []string
	jwtKey          []byte
)

func GetUserFromGToken(ctx context.Context, gTokenString string, authToken *types.AuthToken) error {
	if gitkitClient == nil {
		getConfig(ctx)
	}
	if gTokenString == "" {
		return errors.ErrInvalidToken
	}
	token, err := gitkitClient.ValidateToken(ctx, gTokenString, gitkitAudiences)
	if err != nil {
		return errors.ErrInvalidToken.WithArgs("Invalid gitkit token " + err.Error())
	}
	if time.Now().Sub(token.IssueAt) > 15*time.Minute {
		return errors.ErrInvalidToken.WithArgs("token too old")
	}
	// get user info from gitkit servers
	gitkitUser, err := gitkitClient.UserByLocalID(ctx, token.LocalID)
	if err != nil {
		return errors.ErrExternalDependencyFail.WithArgs("gitkit", "failed to fetch user from gitkit server", err)
	}
	return createSessionToken(ctx, gitkitUser, authToken)
}

func GetUserFromToken(ctx context.Context, authToken *types.AuthToken) error {
	if gitkitClient == nil {
		getConfig(ctx)
	}
	token, err := jwt.Parse(authToken.JWTString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		// Token is invalid
		return errors.ErrInvalidToken.WithArgs(err)
	}
	// Token is good
	err = updateAuthToken(ctx, authToken, token)
	if err != nil {
		return errors.ErrInvalidToken.WithArgs(err)
	}
	return nil
}

func updateAuthToken(ctx context.Context, authToken *types.AuthToken, token *jwt.Token) error {
	var err error
	extractClaims(authToken, token)
	// if issuedTime is > 60 minutes, token info gets refreshed
	if time.Now().Sub(authToken.ITA) > 60*time.Minute {
		err = updateToken(ctx, authToken)
		if err != nil {
			return err
		}
		err = insertClaims(authToken, token)
		return err
	}
	// else data is new enough
	return nil
}

func extractClaims(authToken *types.AuthToken, token *jwt.Token) {
	authToken.User = types.User{
		UserID:      token.Claims["user_id"].(string),
		Name:        token.Claims["name"].(string),
		PhotoURL:    token.Claims["photo_url"].(string),
		Permissions: int32(token.Claims["perm"].(float64)),
	}
	authToken.JTI = int32(token.Claims["jti"].(float64))
	authToken.ITA = time.Unix(int64(token.Claims["ita"].(float64)), 0)
	authToken.Expire = time.Unix(int64(token.Claims["exp"].(float64)), 0)

}

func insertClaims(authToken *types.AuthToken, token *jwt.Token) error {
	token.Claims["user_id"] = authToken.User.UserID
	token.Claims["name"] = authToken.User.Name
	token.Claims["photo_url"] = authToken.User.PhotoURL
	token.Claims["perm"] = int32(authToken.User.Permissions)
	token.Claims["jti"] = authToken.JTI
	token.Claims["ita"] = int(authToken.ITA.Unix())
	token.Claims["exp"] = int(authToken.Expire.Unix())
	var err error
	authToken.JWTString, err = token.SignedString(jwtKey)
	return err
}

func createSessionToken(ctx context.Context, gitkitUser *gitkit.User, authToken *types.AuthToken) error {
	var err error
	userSessions := &UserSessions{}
	err = getUserSessions(ctx, gitkitUser.LocalID, userSessions)
	firstTime := false
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			firstTime = true
		} else {
			return errors.ErrDatastore.WithArgs("get", KindUserSessions, gitkitUser.LocalID, err)
		}
	}

	if firstTime {
		// create gigamuncher kind
		gigamuncher := &types.Gigamuncher{
			UserDetail: types.UserDetail{
				Name:       gitkitUser.DisplayName,
				Email:      gitkitUser.Email,
				PhotoURL:   gitkitUser.PhotoURL,
				ProviderID: gitkitUser.ProviderID,
			},
		}
		errChan := make(chan error, 1)
		account.PutGigamuncher(ctx, gitkitUser.LocalID, gigamuncher, errChan)
		err, ok := <-errChan
		utils.Debugf(ctx, "got this err and ok", err, ok)
		if !ok || err != nil {
			return errors.ErrDatastore.WithArgs("put", types.KindGigamuncher, gitkitUser.LocalID, err)
		}
		utils.Debugf(ctx, "success putting gigamuncher")
		// TODO remove. only tmp thing
		gigachef := &types.Gigachef{
			UserDetail: types.UserDetail{
				Name:       gitkitUser.DisplayName,
				Email:      gitkitUser.Email,
				PhotoURL:   gitkitUser.PhotoURL,
				ProviderID: gitkitUser.ProviderID,
			},
		}
		errChan = make(chan error, 1)
		account.PutGigachef(ctx, gitkitUser.LocalID, gigachef, errChan)
		err, ok = <-errChan
		if !ok || err != nil {
			return errors.ErrDatastore.WithArgs("put", types.KindGigachef, gitkitUser.LocalID, err)
		}
		// TODO end
		// create UserSessions kind
		userSessions.User = types.User{
			UserID:      gitkitUser.LocalID,
			Name:        gitkitUser.DisplayName,
			PhotoURL:    gitkitUser.PhotoURL,
			Permissions: 0,
		}
	}
	// create the token
	authToken.User = userSessions.User
	authToken.ITA = getITATime()
	authToken.Expire = getExpTime()
	authToken.JTI = getNewJTI()
	err = addUserToken(ctx, authToken, userSessions)
	if err != nil {
		return errors.ErrDatastore.WithArgs("put", KindUserSessions, authToken.User.UserID, err)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	err = insertClaims(authToken, token)
	utils.Debugf(ctx, "jwt string", authToken.JWTString)
	return err
}

func getExpTime() time.Time {
	return time.Now().UTC().Add(time.Hour * 24 * 90)
}

func getNewJTI() int32 {
	return rand.Int31()
}

func getITATime() time.Time {
	return time.Now().UTC()
}

// CurrentUser extracts the user information stored in current session.
//
// If there is no existing session, identity toolkit token is checked. If the
// token is valid, a new session is created.
//
// If any error happens, nil is returned.
func CurrentUser(w http.ResponseWriter, req *http.Request) (*types.User, error) {
	ctx := appengine.NewContext(req)
	if gitkitClient == nil {
		getConfig(ctx)
	}
	var err error
	sessionTokenCookie, err := req.Cookie(types.SessionTokenCookieName)
	// doesn't have a session cookie
	if err == http.ErrNoCookie {
		gTokenString := gitkitClient.TokenFromRequest(req)
		if gTokenString == "" {
			return nil, fmt.Errorf("No gtoken")
		}
		authToken := new(types.AuthToken)
		err = GetUserFromGToken(ctx, gTokenString, authToken)
		if err != nil {
			return nil, err
		}
		http.SetCookie(w, &http.Cookie{Name: types.SessionTokenCookieName,
			Value:  authToken.JWTString,
			MaxAge: int(authToken.Expire.Sub(time.Now()).Seconds())})
		http.SetCookie(w, &http.Cookie{Name: types.GitkitCookieName, MaxAge: -1})
		return &authToken.User, nil
	} else if err != nil {
		utils.Errorf(ctx, "Error getting session cookie: %+v", err)
	}
	authToken := &types.AuthToken{
		JWTString: sessionTokenCookie.Value,
	}
	err = GetUserFromToken(ctx, authToken)
	if err != nil {
		return nil, err
	}
	return &authToken.User, nil
}

func getConfig(ctx context.Context) {
	config := config.GetGitkitConfig(ctx)
	// setup gitkit
	c := &gitkit.Config{
		WidgetURL: types.LoginURL,
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
	jwtKey = []byte(config.JWTSecret)
	rand.Seed(time.Now().Unix())
}

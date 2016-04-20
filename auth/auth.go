package auth

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/google/identity-toolkit-go-client/gitkit"
	"golang.org/x/net/context"
	jwt "gopkg.in/dgrijalva/jwt-go.v2"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	gitkitClient    *gitkit.Client
	gitkitAudiences []string
	jwtKey          []byte
)

var (
	errInvalidGToken = errors.ErrorWithCode{Code: errors.CodeSignOut, Message: "Invalid GITKIT token."}
	errInvalidToken  = errors.ErrorWithCode{Code: errors.CodeSignOut, Message: "Invalid Gigatoken."}
	errTokenExpired  = errors.ErrorWithCode{Code: errors.CodeSignOut, Message: "Invalid Gigatoken. Token is expired."}
	errDatastore     = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// DeleteSessionToken removes the session token from valid tokens
func DeleteSessionToken(ctx context.Context, JWTString string) error {
	token, err := getAuthTokenFromString(ctx, JWTString)
	if err != nil {
		return err
	}
	userSessions := &UserSessions{}
	err = getUserSessions(ctx, token.User.ID, userSessions)
	if err != nil {
		return errDatastore.WithError(err)
	}
	for i := len(userSessions.TokenIDs) - 1; i >= 0; i-- {
		if userSessions.TokenIDs[i].JTI == token.JTI {
			// UserSession token should be removed
			userSessions.TokenIDs = append(userSessions.TokenIDs[:i], userSessions.TokenIDs[i+1:]...)
			break
		}
	}
	err = putUserSessions(ctx, token.User.ID, userSessions)
	if err != nil {
		return err
	}
	return nil
}

// SaveUser saves user information such as permissions
func SaveUser(ctx context.Context, user *types.User) error {
	var err error
	userSessions := &UserSessions{}
	err = getUserSessions(ctx, user.ID, userSessions)
	if err != nil {
		return errDatastore.WithError(err)
	}
	userSessions.User = *user
	err = putUserSessions(ctx, user.ID, userSessions)
	if err != nil {
		return err
	}
	return nil
}

// GetSessionWithGToken takes a GITKIT token string.
// Returns: User, JWTString,error
func GetSessionWithGToken(ctx context.Context, gTokenString string) (*types.User, string, error) {
	if gitkitClient == nil {
		getConfig(ctx)
	}
	if gTokenString == "" {
		return nil, "", errInvalidGToken
	}
	gtoken, err := gitkitClient.ValidateToken(ctx, gTokenString, gitkitAudiences)
	if err != nil || time.Now().Sub(gtoken.IssueAt) > 15*time.Minute {
		return nil, "", errInvalidGToken
	}
	// get user info from gitkit servers
	gitkitUser, err := gitkitClient.UserByLocalID(ctx, gtoken.LocalID)
	if err != nil {
		return nil, "", errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "failed to fetch user from gitkit server"}.WithError(err)
	}
	token, err := createSessionToken(ctx, gitkitUser)
	jwtString, err := token.JWTString()
	if err != nil {
		return nil, "", errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Failed to encode user."}
	}
	return &token.User, jwtString, nil
}

func createSessionToken(ctx context.Context, gitkitUser *gitkit.User) (*Token, error) {
	var err error
	userSessions := &UserSessions{}
	err = getUserSessions(ctx, gitkitUser.LocalID, userSessions)
	firstTime := false
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			firstTime = true
		} else {
			return nil, errDatastore.WithError(err)
		}
	}
	if firstTime {
		// create UserSessions kind
		userSessions.User = types.User{
			ID:          gitkitUser.LocalID,
			Name:        gitkitUser.DisplayName,
			Email:       gitkitUser.Email,
			ProviderID:  gitkitUser.ProviderID,
			PhotoURL:    gitkitUser.PhotoURL,
			Permissions: 0,
		}
	}
	// create the token
	token := &Token{
		User:   userSessions.User,
		ITA:    getITATime(),
		JTI:    getNewJTI(),
		Expire: GetExpTime(),
	}
	userSessions.TokenIDs = append(userSessions.TokenIDs, TokenID{JTI: token.JTI, Expire: token.Expire})
	err = putUserSessions(ctx, token.User.ID, userSessions)
	if err != nil {
		return nil, errors.ErrorWithCode{
			Code:    errors.CodeInternalServerErr,
			Message: fmt.Sprintf("error put UserSession(%s) err: %+v", gitkitUser.LocalID, err),
		}
	}
	return token, nil
}

// GetUserFromToken validates a token and extracts the user from it.
// Returns: User, error
func GetUserFromToken(ctx context.Context, JWTString string) (*types.User, error) {
	token, err := getAuthTokenFromString(ctx, JWTString)
	if err != nil {
		return nil, err
	}
	if token.IsExpired() {
		return nil, errTokenExpired
	}
	// TODO if issue time is old, check in database
	// log a "token miss"
	if token.IsOld() {
		userSessions := &UserSessions{}
		err := getUserSessions(ctx, token.User.ID, userSessions)
		if err != nil {
			// error doesn't matter. They should just call RefreshToken
			return nil, errInvalidToken.WithMessage("Datastore error.").WithError(err)
		}
		return &userSessions.User, nil
	}
	return &token.User, nil
}

func getAuthTokenFromString(ctx context.Context, JWTString string) (*Token, error) {
	if jwtKey == nil {
		getConfig(ctx)
	}
	jwtToken, err := jwt.Parse(JWTString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil || !jwtToken.Valid {
		// Token is invalid
		return nil, errInvalidToken
	}
	token, err := extractClaims(jwtToken)
	if err != nil {
		return nil, errInvalidToken
	}
	return token, nil
}

// RefreshToken takes a token string and returns a new token with updated claims
// and expiration time. It also, invalidates the old token in an hour.
func RefreshToken(ctx context.Context, JWTString string) (string, error) {
	// TODO test
	token, err := getAuthTokenFromString(ctx, JWTString)
	if err != nil {
		return "", err
	}
	userSessions := &UserSessions{}
	err = getUserSessions(ctx, token.User.ID, userSessions)
	if err != nil {
		return "", errDatastore.WithError(err)
	}
	found := false
	needPut := false
	for i := len(userSessions.TokenIDs) - 1; i >= 0; i-- {
		if time.Now().After(userSessions.TokenIDs[i].Expire) {
			// UserSession token should be removed
			userSessions.TokenIDs = append(userSessions.TokenIDs[:i], userSessions.TokenIDs[i+1:]...)
			needPut = true
		}
	}
	for i := len(userSessions.TokenIDs) - 1; i >= 0; i-- {
		if userSessions.TokenIDs[i].JTI == token.JTI {
			found = true
			if userSessions.TokenIDs[i].UpdatedToJTI != 0 {
				token.JTI = userSessions.TokenIDs[i].UpdatedToJTI
				token.Expire = userSessions.TokenIDs[i].UpdateToExpire
			} else {
				newJTI := getNewJTI()
				token.JTI = newJTI
				userSessions.TokenIDs[i].UpdatedToJTI = newJTI
				expTime := GetExpTime()
				token.Expire = expTime
				userSessions.TokenIDs[i].UpdateToExpire = expTime
				userSessions.TokenIDs = append(userSessions.TokenIDs, TokenID{JTI: newJTI, Expire: expTime})
				userSessions.TokenIDs[i].Expire = time.Now().UTC().Add(time.Hour * 1)
				needPut = true
			}
			token.ITA = getITATime()
			break
		}
	}
	if needPut {
		err = putUserSessions(ctx, token.User.ID, userSessions)
		if err != nil {
			return "", errDatastore.WithError(err)
		}
	}
	if !found {
		return "", errInvalidToken
	}
	token.User = userSessions.User
	jwtString, err := token.JWTString()
	if err != nil {
		return "", err
	}
	return jwtString, nil
}

// GetExpTime returns the time a token should expire from now
func GetExpTime() time.Time {
	return time.Now().UTC().Add(time.Hour * 24 * 60)
}

func getNewJTI() int32 {
	return rand.Int31()
}

func getITATime() time.Time {
	return time.Now().UTC()
}

func getJWTToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

func extractClaims(jwtToken *jwt.Token) (*Token, error) {
	getStringClaim := func(name string, ok bool) (string, bool) {
		if ok {
			tmp, ok2 := jwtToken.Claims[name].(string)
			return tmp, ok2
		}
		return "", ok
	}
	getInt32Claim := func(name string, ok bool) (int32, bool) {
		if ok {
			tmp, ok2 := jwtToken.Claims[name].(float64)
			if ok2 {
				return int32(tmp), ok2
			}
		}
		return 0, ok
	}
	getTimeClaim := func(name string, ok bool) (time.Time, bool) {
		if ok {
			tmp, ok2 := jwtToken.Claims[name].(float64)
			if ok2 {
				return time.Unix(int64(tmp), 0), ok2
			}
		}
		return time.Now(), ok
	}
	var userID, name, email, providerID, photoURL string
	var permissions, jti int32
	var ita, expire time.Time
	ok := true
	userID, ok = getStringClaim("id", ok)
	name, ok = getStringClaim("name", ok)
	email, ok = getStringClaim("email", ok)
	providerID, ok = getStringClaim("provider_id", ok)
	photoURL, ok = getStringClaim("photo_url", ok)
	permissions, ok = getInt32Claim("perm", ok)
	jti, ok = getInt32Claim("jti", ok)
	ita, ok = getTimeClaim("ita", ok)
	expire, ok = getTimeClaim("exp", ok)
	if !ok {
		return nil, errInvalidToken
	}
	token := new(Token)
	token.User = types.User{
		ID:          userID,
		Name:        name,
		Email:       email,
		ProviderID:  providerID,
		PhotoURL:    photoURL,
		Permissions: permissions,
	}
	token.JTI = jti
	token.ITA = ita
	token.Expire = expire
	return token, nil
}

func getConfig(ctx context.Context) {
	config := config.GetGitkitConfig(ctx)
	// setup gitkit
	c := &gitkit.Config{
		WidgetURL: loginURL,
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
}

func init() {
	rand.Seed(time.Now().Unix())
}

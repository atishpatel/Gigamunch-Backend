package auth

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// loginURL is the url to login
	loginURL = "/login"
	// kindUserSessions is a datastore kind for user information
	kindUserSessions = "UserSessions"
)

// UserSessions is stored in the database to indicate valid user sessions
type UserSessions struct {
	User     types.User `datastore:",noindex"`
	TokenIDs []TokenID  `datastore:",noindex"`
}

// TokenID has unique ids and exp time for tokens
type TokenID struct {
	JTI            int32     `datastore:",noindex"`
	UpdatedToJTI   int32     `datastore:",noindex"`
	Expire         time.Time `datastore:",noindex"`
	UpdateToExpire time.Time `datastore:",noindex"`
}

// Token contains the the user and the user's token string
type Token struct {
	User   types.User
	JTI    int32
	ITA    time.Time
	Expire time.Time
}

// IsExpired returns true if the token is expired
func (token *Token) IsExpired() bool {
	return time.Now().After(token.Expire)
}

// IsOld returns true if the token issue time is >60 minutes
func (token *Token) IsOld() bool {
	return time.Now().Sub(token.ITA) > 1*time.Minute // TODO switch back to 60
}

// JWTString returns a signed JSON Web Token string
func (token *Token) JWTString() (string, error) {
	jwtToken := getJWTToken()
	jwtToken.Claims["id"] = token.User.ID
	jwtToken.Claims["name"] = token.User.Name
	jwtToken.Claims["email"] = token.User.Email
	jwtToken.Claims["provider_id"] = token.User.ProviderID
	jwtToken.Claims["photo_url"] = token.User.PhotoURL
	jwtToken.Claims["perm"] = int32(token.User.Permissions)
	jwtToken.Claims["jti"] = token.JTI
	jwtToken.Claims["ita"] = int(token.ITA.Unix())
	jwtToken.Claims["exp"] = int(token.Expire.Unix())
	jwtString, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		return "", errors.ErrorWithCode{
			Code:    errors.CodeInternalServerErr,
			Message: fmt.Sprintf("Error inserting token claims."),
		}.WithError(err)
	}
	return jwtString, nil
}

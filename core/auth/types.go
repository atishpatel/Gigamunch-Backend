package auth

import (
	"context"
	"fmt"
	"time"

	jwt "gopkg.in/dgrijalva/jwt-go.v2"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
)

const kind = "User"

// UserSessions is stored in the database to indicate valid user sessions
type UserSessions struct {
	Provider string      `datastore:",noindex"`
	User     common.User `datastore:",index"`
	TokenIDs []TokenID   `datastore:",noindex"`
}

// TokenID has unique ids and exp time for tokens
type TokenID struct {
	OriginalIAT    time.Time `datastore:",noindex"`
	JTI            int32     `datastore:",noindex"`
	UpdatedToJTI   int32     `datastore:",noindex"`
	Expire         time.Time `datastore:",noindex"`
	UpdateToExpire time.Time `datastore:",noindex"`
}

// Token contains the the user and the user's token string
type Token struct {
	User   common.User
	JTI    int32
	IAT    time.Time
	Expire time.Time
}

// IsExpired returns true if the token is expired
func (token *Token) IsExpired() bool {
	return time.Now().After(token.Expire)
}

// IsOld returns true if the token issue time is >60 minutes
func (token *Token) IsOld() bool {
	return time.Since(token.IAT) > 60*time.Minute
}

// JWTString returns a signed JSON Web Token string
func (token *Token) JWTString() (string, error) {
	jwtToken := getJWTToken()
	jwtToken.Claims["id"] = token.User.ID
	jwtToken.Claims["auth_id"] = token.User.AuthID
	jwtToken.Claims["first_name"] = token.User.FirstName
	jwtToken.Claims["last_name"] = token.User.LastName
	jwtToken.Claims["email"] = token.User.Email
	jwtToken.Claims["photo_url"] = token.User.PhotoURL
	jwtToken.Claims["perm"] = int32(token.User.Permissions)
	jwtToken.Claims["jti"] = token.JTI
	jwtToken.Claims["iat"] = int(token.IAT.Unix())
	jwtToken.Claims["exp"] = int(token.Expire.Unix())
	jwtString, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		return "", errInternal.WithMessage("Error inserting token claims.").WithError(err)
	}
	return jwtString, nil
}

func newToken(ctx context.Context, JWTString string) (*Token, error) {
	jwtToken, err := jwt.Parse(JWTString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil || !jwtToken.Valid {
		// Token is invalid
		return nil, errInvalidToken.WithError(err).Annotate("jwtToken is not valid")
	}
	token, err := extractClaims(jwtToken)
	if err != nil {
		return nil, errInvalidToken.Annotate("failed to extract token claims")
	}
	return token, nil
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
	getInt64Claim := func(name string, ok bool) (int64, bool) {
		if ok {
			tmp, ok2 := jwtToken.Claims[name].(float64)
			if ok2 {
				return int64(tmp), ok2
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
	var authID, firstName, lastName, email, photoURL string
	var permissions, jti int32
	var userID int64
	var iat, expire time.Time
	ok := true
	userID, ok = getInt64Claim("id", ok)
	authID, ok = getStringClaim("auth_id", ok)
	firstName, ok = getStringClaim("first_name", ok)
	lastName, ok = getStringClaim("last_name", ok)
	email, ok = getStringClaim("email", ok)
	photoURL, ok = getStringClaim("photo_url", ok)
	permissions, ok = getInt32Claim("perm", ok)
	jti, ok = getInt32Claim("jti", ok)
	iat, ok = getTimeClaim("iat", ok)
	expire, ok = getTimeClaim("exp", ok)
	if !ok {
		return nil, errInvalidToken.Wrap("failed to extract claims from token")
	}
	token := new(Token)
	token.User = common.User{
		ID:          userID,
		AuthID:      authID,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhotoURL:    photoURL,
		Permissions: permissions,
	}
	token.JTI = jti
	token.IAT = iat
	token.Expire = expire
	return token, nil
}

func getJWTToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

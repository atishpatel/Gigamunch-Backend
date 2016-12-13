package auth

import (
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

var (
	errNotAdmin = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
)

func VerifyCook(ctx context.Context, user *types.User, userID string) error {
	var err error
	if !user.IsAdmin() {
		utils.Errorf(ctx, "no admin user(%v) attemted to verify cook.", *user)
		return errNotAdmin
	}
	userSessions := &UserSessions{}
	err = getUserSessions(ctx, userID, userSessions)
	if err != nil {
		return err
	}
	userSessions.User.SetVerifiedChef(true)
	err = putUserSessions(ctx, userID, userSessions)
	return err
}

// GetGigatokenViaAdmin returns a valid Gigatoken for a user if user is an admin.s
func GetGigatokenViaAdmin(ctx context.Context, user *types.User, userID string) (string, error) {
	if !user.IsAdmin() {
		return "", errNotAdmin
	}
	userSessions := new(UserSessions)
	err := getUserSessions(ctx, userID, userSessions)
	if err != nil {
		return "", errDatastore.WithError(err).Wrapf("failed to getUserSession for userID(%s)", userID)
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
		return "", errDatastore.WithError(err).Wrapf("failed to putUserSession for userID(%s)", userID)
	}
	jwtString, err := token.JWTString()
	if err != nil {
		return "", errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Failed to encode user."}
	}
	return jwtString, nil
}

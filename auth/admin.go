package auth

import (
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

func VerifyChef(ctx context.Context, user *types.User, userID string) error {
	var err error
	if !user.IsAdmin() {
		utils.Errorf(ctx, "no admin user(%v) attemted to verify chef.", *user)
		return errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
	}
	userSessions := &UserSessions{}
	err = getUserSessions(ctx, userID, userSessions)
	if err != nil {
		return err
	}
	userSessions.User.SetVerifiedChef(true)
	err = putUserSessions(ctx, userID, userSessions)
	if err != nil {
		return err
	}
	return nil
}

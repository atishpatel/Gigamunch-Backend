package application

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/account"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"github.com/atishpatel/Gigamunch-Backend/utils/map"
)

// GetApplications gets all applications
func GetApplications(ctx context.Context, user *types.User) ([]*ChefApplication, error) {
	if !user.IsAdmin() {
		utils.Errorf(ctx, "user(%v) attemted to do an admin task.", *user)
		return nil, errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
	}
	chefApplications, err := getAll(ctx)
	return chefApplications, err
}

// GetApplication gets a chef application
func GetApplication(ctx context.Context, user *types.User) (*ChefApplication, error) {
	var err error
	chefApplication := new(ChefApplication)
	err = get(ctx, user.ID, chefApplication)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			chefApplication.Name = user.Name
			chefApplication.Email = user.Email
			return chefApplication, nil
		}
		return nil, err
	}
	return chefApplication, nil
}

// SubmitApplication saves a ChefApplication.
// A token should be refreshed if this function is called.
func SubmitApplication(ctx context.Context, user *types.User, chefApplication *ChefApplication) (*ChefApplication, error) {
	var err error
	chefApplicationEntity := &ChefApplication{}
	err = get(ctx, user.ID, chefApplicationEntity)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}

	if err != nil && err == datastore.ErrNoSuchEntity {
		chefApplication.ApplicationProgress = 1
		if chefApplication.Address.String() != chefApplicationEntity.Address.String() {
			err = maps.GetGeopointFromAddress(ctx, &chefApplication.Address)
			if err != nil {
				return nil, err
			}
		} else {
			chefApplication.Address = chefApplicationEntity.Address
		}
	} else {
		chefApplication.ApplicationProgress = chefApplicationEntity.ApplicationProgress
	}
	chefApplication.UserID = user.ID
	chefApplication.LastUpdatedDateTime = time.Now().UTC()
	if chefApplicationEntity.CreatedDateTime.IsZero() {
		chefApplication.CreatedDateTime = time.Now().UTC()
	}
	err = put(ctx, user.ID, chefApplication)
	if err != nil {
		return nil, err
	}
	err = setAuthUserAndPerm(ctx, user, chefApplication.Name, chefApplication.Email)
	if err != nil {
		return nil, err
	}
	err = account.SaveUserInfo(ctx, user, &chefApplication.Address, chefApplication.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return chefApplication, nil
}

func setAuthUserAndPerm(ctx context.Context, user *types.User, name string, email string) error {
	userChanged := false
	if !user.IsChef() {
		user.SetChef(true)
		userChanged = true
	}
	if !user.HasAddress() {
		user.SetAddress(true)
		userChanged = true
	}
	if user.Name != name {
		user.Name = name
		userChanged = true
	}
	if user.Email != email {
		user.Email = email
		userChanged = true
	}
	if userChanged {
		return auth.SaveUser(ctx, user)
	}
	return nil
}

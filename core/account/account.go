package account

import (
	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/core/gigamuncher"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

// func UpdateUserDetails(authtoken, userdetails)
// func SaveAddress(authToken, address)
// update gigamuncher address
// if user.IsChef(), update gigachef
// set user permissions in auth

// SaveUserInfo updates the user info all the places it is stored.
// A token should be refreshed after making this call.
func SaveUserInfo(ctx context.Context, user *types.User, address *types.Address, phoneNumber string) error {
	var err error
	permChanged := false
	if address != nil && !user.HasAddress() {
		user.SetAddress(true)
		permChanged = true
	}
	if permChanged {
		err = auth.SaveUser(ctx, user)
	}
	chefErrChan := make(chan error, 1)
	if user.IsChef() {
		// update chef info
		go func() {
			chefErrChan <- gigachef.SaveUserInfo(ctx, user, address, phoneNumber)
		}()
	}
	// update gigamuncher info
	err = gigamuncher.SaveUserInfo(ctx, user, address)
	if err != nil {
		return err
	}
	if user.IsChef() {
		err = <-chefErrChan
		if err != nil {
			return err
		}
	}
	return nil
}

package gigachef

import (
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"

	"appengine/datastore"
)

// SaveUserInfo is to save a user's info. Only exposed for account package.
// Please use the account package's func instead of this one.
func SaveUserInfo(ctx context.Context, user *types.User, address *types.Address, phoneNumber string) error {
	var err error
	gigachef := new(Gigachef)
	err = get(ctx, user.ID, gigachef)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}
	gigachef.Name = user.Name
	gigachef.PhotoURL = user.PhotoURL
	gigachef.Email = user.Email
	gigachef.ProviderID = user.ProviderID
	gigachef.Address = *address
	gigachef.PhoneNumber = phoneNumber
	err = put(ctx, user.ID, gigachef)
	if err != nil {
		return err
	}
	return nil
}

func GetAddress(ctx context.Context, userID string) (*types.Address, error) {
	var err error
	gigachef := &Gigachef{}
	err = get(ctx, userID, gigachef)
	if err != nil {
		return nil, err
	}
	return &gigachef.Address, nil
}

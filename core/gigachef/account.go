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
	if err != nil && err.Error() != datastore.ErrNoSuchEntity.Error() {
		return err
	}
	gigachef.Name = user.Name
	gigachef.PhotoURL = user.PhotoURL
	gigachef.Email = user.Email
	gigachef.ProviderID = user.ProviderID
	if address != nil {
		gigachef.Address = *address
	}
	if phoneNumber != "" {
		gigachef.PhoneNumber = phoneNumber
	}
	err = put(ctx, user.ID, gigachef)
	if err != nil {
		return err
	}
	return nil
}

// GetDeliveryInfo gets info related to Gigachef's delivery
// returns: address, gigachefDeliveryRange, error
func GetDeliveryInfo(ctx context.Context, user *types.User) (*types.Address, int, error) {
	var err error
	gigachef := new(Gigachef)
	err = get(ctx, user.ID, gigachef)
	if err != nil {
		return nil, 0, errDatastore.WithError(err)
	}
	return &gigachef.Address, gigachef.DeliveryRange, nil
}

// IncrementNumPost increases NumPosts for a gigachef by 1
func IncrementNumPost(ctx context.Context, user *types.User) error {
	var err error
	gigachef := new(Gigachef)
	err = get(ctx, user.ID, gigachef)
	if err != nil {
		return errDatastore.WithError(err)
	}
	gigachef.NumPosts++
	err = put(ctx, user.ID, gigachef)
	if err != nil {
		return errDatastore.WithError(err)
	}
	return nil
}

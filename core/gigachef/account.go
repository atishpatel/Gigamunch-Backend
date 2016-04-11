package gigachef

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

// SaveUserInfo is to save a user's info. Only exposed for account package.
// Please use the account package's func instead of this one.
func SaveUserInfo(ctx context.Context, user *types.User, address *types.Address, phoneNumber string) error {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}
	chef.Name = user.Name
	chef.PhotoURL = user.PhotoURL
	chef.Email = user.Email
	chef.ProviderID = user.ProviderID
	if address != nil {
		chef.Address = *address
	}
	if phoneNumber != "" {
		chef.PhoneNumber = phoneNumber
	}
	err = put(ctx, user.ID, chef)
	if err != nil {
		return err
	}
	return nil
}

// GetDeliveryInfo gets info related to Gigachef's delivery
// returns: address, gigachefDeliveryRange, error
func GetDeliveryInfo(ctx context.Context, user *types.User) (*types.Address, int, error) {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil {
		return nil, 0, errDatastore.WithError(err)
	}
	return &chef.Address, chef.DeliveryRange, nil
}

// IncrementNumPost increases NumPosts for a Gigachef by 1
func IncrementNumPost(ctx context.Context, user *types.User) error {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil {
		return errDatastore.WithError(err)
	}
	chef.NumPosts++
	err = put(ctx, user.ID, chef)
	if err != nil {
		return errDatastore.WithError(err)
	}
	return nil
}

// GetInfo returns nonsensitive Gigachef details
func GetInfo(ctx context.Context, id string) (*Gigachef, error) {
	// TODO switch so it's returns a 'nonsensitive' gigachef info
	chef := new(Gigachef)
	err := get(ctx, id, chef)
	if err != nil {
		return nil, errDatastore.WithError(err)
	}
	return chef, nil
}

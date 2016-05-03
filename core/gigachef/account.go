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
		return errDatastore.WithError(err).Wrap("cannot get gigachef")
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
		return errDatastore.WithError(err).Wrap("cannot put gigachef")
	}
	return nil
}

// GetInfo returns nonsensitive Gigachef details
func GetInfo(ctx context.Context, id string) (*Gigachef, error) {
	// TODO switch so it's returns a 'nonsensitive' gigachef info
	chef := new(Gigachef)
	err := get(ctx, id, chef)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	return chef, nil
}

// PostInfoResp contains information related to a post
type PostInfoResp struct {
	Address         types.Address `json:"address"`
	DeliveryRange   int32         `json:"delivery_range"`
	BTSubMerchantID string        `json:"bt_sub_merchant_id"`
}

// GetPostInfo returns info related to a post
// returns: *PostInfoResp, error
func GetPostInfo(ctx context.Context, user *types.User) (*PostInfoResp, error) {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	postInfo := &PostInfoResp{
		Address:         chef.Address,
		DeliveryRange:   chef.DeliveryRange,
		BTSubMerchantID: chef.BTSubMerchantID,
	}
	return postInfo, nil
}

// IncrementNumPost increases NumPosts for a Gigachef by 1
func IncrementNumPost(ctx context.Context, user *types.User) error {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	chef.NumPosts++
	err = put(ctx, user.ID, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot put gigachef")
	}
	return nil
}

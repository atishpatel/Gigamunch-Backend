package gigamuncher

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

// SaveUserInfo is to save a user's info. Only exposed for account package.
// Please use the account package's func instead of this one.
func SaveUserInfo(ctx context.Context, user *types.User, address *types.Address) error {
	var err error
	gigamuncher := new(Gigamuncher)
	err = get(ctx, user.ID, gigamuncher)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}
	gigamuncher.Name = user.Name
	gigamuncher.PhotoURL = user.PhotoURL
	gigamuncher.Email = user.Email
	gigamuncher.ProviderID = user.ProviderID
	addresses := []Addresses{Addresses{LastUsedDataTime: time.Now().UTC(), Address: *address}}
	gigamuncher.Addresses = append(addresses, gigamuncher.Addresses...)
	err = put(ctx, user.ID, gigamuncher)
	if err != nil {
		return err
	}
	return nil
}

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
	changed := false
	gigamuncher := new(Gigamuncher)
	err = get(ctx, user.ID, gigamuncher)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}
	if gigamuncher.Name != user.Name {
		gigamuncher.Name = user.Name
		changed = true
	}
	if gigamuncher.PhotoURL == user.PhotoURL {
		gigamuncher.PhotoURL = user.PhotoURL
		changed = true
	}
	if gigamuncher.Email != user.Email {
		gigamuncher.Email = user.Email
		changed = true
	}
	if gigamuncher.ProviderID != user.ProviderID {
		gigamuncher.ProviderID = user.ProviderID
		changed = true
	}
	if address != nil {
		addresses := []Addresses{Addresses{LastUsedDataTime: time.Now().UTC(), Address: *address}}
		gigamuncher.Addresses = append(addresses, gigamuncher.Addresses...)
		changed = true
	}
	if changed {
		err = put(ctx, user.ID, gigamuncher)
		if err != nil {
			return err
		}
	}
	return nil
}

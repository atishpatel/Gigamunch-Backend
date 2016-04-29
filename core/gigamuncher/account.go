package gigamuncher

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	errDatastore = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// SaveUserInfo is to save a user's info. Only exposed for account package.
// Please use the account package's func instead of this one.
func SaveUserInfo(ctx context.Context, user *types.User, address *types.Address) error {
	var err error
	changed := false
	muncher := new(Gigamuncher)
	err = get(ctx, user.ID, muncher)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return errDatastore.WithError(err).Wrap("cannot save gigamuncher info because cannot get gigamuncher")
	}
	if muncher.Name != user.Name {
		muncher.Name = user.Name
		changed = true
	}
	if muncher.PhotoURL == user.PhotoURL {
		muncher.PhotoURL = user.PhotoURL
		changed = true
	}
	if muncher.Email != user.Email {
		muncher.Email = user.Email
		changed = true
	}
	if muncher.ProviderID != user.ProviderID {
		muncher.ProviderID = user.ProviderID
		changed = true
	}
	if address != nil {
		addresses := []Addresses{Addresses{LastUsedDataTime: time.Now().UTC(), Address: *address}}
		muncher.Addresses = append(addresses, muncher.Addresses...)
		changed = true
	}
	if changed {
		err = put(ctx, user.ID, muncher)
		if err != nil {
			return errDatastore.WithError(err).Wrap("cannot save gigamuncher info because cannot put gigamuncher")
		}
	}
	return nil
}

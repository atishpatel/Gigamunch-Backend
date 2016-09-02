package gigamuncher

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/maps"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	errDatastore = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// Client is a client for gigamuncher
type Client struct {
	ctx context.Context
}

// New returns a client for gigamuncher
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

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
	if muncher.CreatedDatetime.IsZero() {
		muncher.CreatedDatetime = time.Now()
		changed = true
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
		addresses := []Addresses{Addresses{AddedDateTime: time.Now().UTC(), Address: *address}}
		muncher.Addresses = append(addresses, muncher.Addresses...)
		changed = true
	}
	if muncher.BTCustomerID == "" {
		tmpID := user.ID
		for len(tmpID) <= 36 {
			tmpID += tmpID
		}
		muncher.BTCustomerID = tmpID[:36]
	}
	if changed {
		err = put(ctx, user.ID, muncher)
		if err != nil {
			return errDatastore.WithError(err).Wrap("cannot save gigamuncher info because cannot put gigamuncher")
		}
	}
	return nil
}

// GetBTCustomerID gets the Braintree customer id
func GetBTCustomerID(ctx context.Context, muncherID string) (string, error) {
	muncher := new(Gigamuncher)
	err := get(ctx, muncherID, muncher)
	if err != nil {
		return "", errDatastore.WithError(err).Wrap("cannot get gigamuncher")
	}
	return muncher.BTCustomerID, nil
}

// SaveBTCustomerID saves the Braintree customer id
func SaveBTCustomerID(ctx context.Context, muncherID, btCustomerID string) error {
	muncher := new(Gigamuncher)
	err := get(ctx, muncherID, muncher)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot get gigamuncher")
	}
	if muncher.BTCustomerID != btCustomerID {
		muncher.BTCustomerID = btCustomerID
		err := put(ctx, muncherID, muncher)
		if err != nil {
			return errDatastore.WithError(err).Wrap("cannot put gigamuncher")
		}
	}
	return nil
}

// GetAddresses returns the addresses for a gigamuncher
func (c *Client) GetAddresses(muncherID string) ([]Addresses, error) {
	muncher := new(Gigamuncher)
	err := get(c.ctx, muncherID, muncher)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get gigamuncher")
	}
	return muncher.Addresses, nil
}

// SelectAddress adds and selects the address as main address
func (c *Client) SelectAddress(muncherID string, address *types.Address) (*Addresses, error) {
	muncher := new(Gigamuncher)
	err := get(c.ctx, muncherID, muncher)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get gigamuncher")
	}
	var found bool
	var a *Addresses
	for i := range muncher.Addresses {
		// check if address already exists
		if address.String() == muncher.Addresses[i].Address.String() {
			muncher.Addresses[i].Selected = true
			a = &muncher.Addresses[i]
			found = true
		} else {
			muncher.Addresses[i].Selected = false
		}
	}
	// add address if not found
	if !found {
		if address.Longitude == 0 && address.Latitude == 0 {
			err = maps.GetGeopointFromAddress(c.ctx, address)
			if err != nil {
				return nil, errors.Wrap("failed to maps.GetGeopointFromAddress", err)
			}
		}
		a = &Addresses{
			Address:       *address,
			AddedDateTime: time.Now(),
			Selected:      true,
		}
		muncher.Addresses = append(muncher.Addresses, *a)
	}
	err = put(c.ctx, muncherID, muncher)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot put gigamuncher")
	}
	return a, nil
}

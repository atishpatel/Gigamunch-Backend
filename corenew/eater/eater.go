package eater

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/core/maps"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// Client is a client for the eater package.
type Client struct {
	ctx context.Context
}

// New returns a new eater Client.
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// Get gets the eater.
func (c *Client) Get(id string) (*Eater, error) {
	eater, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get eater(%s)", id)
	}
	return eater, nil
}

// GetBTCustomerID gets the eater Braintree customer id
func GetBTCustomerID(ctx context.Context, id string) (string, error) {
	eater, err := get(ctx, id)
	if err != nil {
		return "", errDatastore.WithError(err).Wrapf("cannot get eater(%s)", id)
	}
	return eater.BTCustomerID, nil
}

// GetAddresses gets the eater Addresses.
func (c *Client) GetAddresses(id string) ([]Address, error) {
	eater, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get eater(%s)", id)
	}
	return eater.Addresses, nil
}

// Update updates a eater's info.
func (c *Client) Update(user *types.User) (*Eater, error) {
	eater, err := get(c.ctx, user.ID)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, errDatastore.WithError(err).Wrapf("failed to get eater(%s)", user.ID)
	}
	if eater.CreatedDatetime.IsZero() {
		eater.CreatedDatetime = time.Now()
	}
	eater.ID = user.ID
	eater.Name = user.Name
	eater.PhotoURL = user.PhotoURL
	eater.Email = user.Email
	eater.ProviderID = user.ProviderID
	if eater.BTCustomerID == "" {
		tmpID := user.ID
		for len(tmpID) <= 36 {
			tmpID += tmpID
		}
		eater.BTCustomerID = tmpID[:36]
	}
	err = put(c.ctx, eater.ID, eater)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to put eater(%s)", eater.ID)
	}
	return eater, nil
}

// SelectAddress adds and selects the address as main address
func (c *Client) SelectAddress(id string, address *types.Address) (*Address, error) {
	eater, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot get eater(%s)", id)
	}
	var found bool
	var a *Address
	for i := range eater.Addresses {
		// check if address already exists
		if address.String() == eater.Addresses[i].Address.String() {
			eater.Addresses[i].Selected = true
			a = &eater.Addresses[i]
			found = true
		} else {
			eater.Addresses[i].Selected = false
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
		a = &Address{
			Address:       *address,
			AddedDateTime: time.Now(),
			Selected:      true,
		}
		eater.Addresses = append(eater.Addresses, *a)
	}
	err = put(c.ctx, id, eater)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot put eater(%s)", id)
	}
	return a, nil
}

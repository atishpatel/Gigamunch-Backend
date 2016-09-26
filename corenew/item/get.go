package item

import (
	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	errInvalidParameter   = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errUnauthorizedAccess = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have permission to item."}
)

// Client is a client for the cook package.
type Client struct {
	ctx context.Context
}

// New returns a new cook Client.
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// Get gets the item.
func (c *Client) Get(id int64) (*Item, error) {
	item, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get item(%d)", id)
	}
	return item, nil
}

// GetAllByCook returns an array of items of the cookID
func (c *Client) GetAllByCook(cookID string) ([]int64, []Item, error) {
	ids, items, err := getCookItems(c.ctx, cookID)
	if err != nil {
		return nil, nil, errDatastore.WithError(err).Wrap("cannot getCookItems")
	}
	return ids, items, nil
}

// GetMulti returns an array of items.
func (c *Client) GetMulti(ids []int64) ([]Item, error) {
	items, err := getMulti(c.ctx, ids)
	if err != nil {
		return nil, errors.Wrap("failed to getMulti items", err)
	}
	return items, nil
}

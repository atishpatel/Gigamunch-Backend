package menu

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// Client is a client for the cook package.
type Client struct {
	ctx context.Context
}

// New returns a new cook Client.
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// Get gets the Menu.
func (c *Client) Get(id int64) (*Menu, error) {
	menu, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get Menu(%d)", id)
	}
	return menu, nil
}

// GetCookMenus returns an array of all the Menus for a Cook.
func (c *Client) GetCookMenus(cookID string) ([]int64, []Menu, error) {
	ids, menus, err := getCookMenus(c.ctx, cookID)
	if err != nil {
		return nil, nil, errDatastore.WithError(err).Wrap("cannot getCookMenus")
	}
	return ids, menus, nil
}

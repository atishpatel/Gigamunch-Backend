package menu

import (
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	errInvalidParameter   = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errUnauthorizedAccess = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have permission to menu."}
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

// GetMulti returns an array of Menus.
func (c *Client) GetMulti(ids []int64) ([]Menu, error) {
	menus, err := getMulti(c.ctx, ids)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to getMulti menu ids: %v", ids)
	}
	return menus, nil
}

// GetCookMenus returns an array of all the Menus for a Cook.
func (c *Client) GetCookMenus(cookID string) ([]int64, []Menu, error) {
	ids, menus, err := getCookMenus(c.ctx, cookID)
	if err != nil {
		return nil, nil, errDatastore.WithError(err).Wrap("cannot getCookMenus")
	}
	return ids, menus, nil
}

// Save saves the Menu.
func (c *Client) Save(user *types.User, id int64, cookID, name string, color Color) (int64, *Menu, error) {
	if cookID == "" {
		return 0, nil, errInvalidParameter.Wrap("cook id cannot be empty")
	}
	var menu *Menu
	var err error
	if id == 0 {
		// create new menu
		menu = &Menu{
			CreatedDateTime: time.Now(),
			CookID:          cookID,
		}
	} else {
		// get menu
		menu, err = get(c.ctx, id)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to get Menu(%d)", id)
		}
	}
	if !user.IsAdmin() && menu.CookID != cookID {
		return 0, nil, errUnauthorizedAccess.Wrapf("CookID(%s) does not have permission to change menu(%d)", cookID, id)
	}
	menu.CookID = cookID
	menu.Name = name
	if color.isZero() {
		menu.Color = NewColor()
	} else {
		menu.Color = color
	}
	if id == 0 {
		// create menu
		id, err = putIncomplete(c.ctx, menu)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to putIncomplete Menu(%v)", *menu)
		}
	} else {
		// update menu
		err = put(c.ctx, id, menu)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to put Menu(%d)", id)
		}
	}
	return id, menu, nil
}
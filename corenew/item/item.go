package item

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	errInvalidParameter   = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errSQLDB              = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	errUnauthorizedAccess = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have permission to item."}
)

// Client is a client for the cook package.
type Client struct {
	ctx context.Context
}

// New returns a new cook Client.
func New(ctx context.Context) *Client {
	connectOnce.Do(func() {
		connectSQL(ctx)
	})
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

// func (c *Client) GetActiveItems(startIndex, endIndex int, long, lat float64) ([]int64, []Item, error) {
// return ids, items, nil
// }

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

// Save saves the item.
func (c *Client) Save(user *types.User, id, menuID int64, cookID, title, desc string,
	dietaryConcerns DietaryConcerns, ingredients, photos []string,
	cookPricePerServing float32, minServings, maxServings int32) (int64, *Item, error) {
	if cookID == "" {
		return 0, nil, errInvalidParameter.Wrap("cook id cannot be empty")
	}
	if menuID == 0 {
		return 0, nil, errInvalidParameter.Wrap("menu id cannot be 0")
	}
	var item *Item
	var err error
	if id == 0 {
		// create is a new item
		item = &Item{
			CreatedDateTime: time.Now(),
			CookID:          cookID,
		}
	} else {
		// get item
		item, err = get(c.ctx, id)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to get item(%d)", id)
		}
	}
	if !user.IsAdmin() && item.CookID != cookID {
		return 0, nil, errUnauthorizedAccess.Wrapf("CookID(%s) does not have permission to change item(%d)", cookID, id)
	}
	item.MenuID = menuID
	item.CookID = cookID
	item.Title = title
	item.Description = desc
	item.DietaryConcerns = dietaryConcerns
	item.Ingredients = ingredients
	item.Photos = photos
	item.CookPricePerServing = cookPricePerServing
	item.MinServings = minServings
	item.MaxServings = maxServings

	if id == 0 {
		// TODO activate item if cook is a verifedCook
		// insertOrUpdateActiveItem(id, item, cook.Address.Latitude, cook.Address.Longitude)
		// item.Active = true
		// create item
		id, err = putIncomplete(c.ctx, item)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to putIncomplete Item(%v)", *item)
		}
	} else {
		// TODO update in sql active items if cook is a verifedCook
		// insertOrUpdateActiveItem(id, item, cook.Address.Latitude, cook.Address.Longitude)
		// item.Active = true
		// update item
		err = put(c.ctx, id, item)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to put Item(%d)", id)
		}
	}
	return id, item, nil
}

// Activate activates an item
func (c *Client) Activate(user *types.User, id int64) error {
	item, err := get(c.ctx, id)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return errInvalidParameter.WithError(err).Wrapf("item(%d) does not exist", id)
		}
		return errDatastore.WithError(err).Wrapf("failed to get item(%d)", id)
	}
	if !user.IsAdmin() && user.ID != item.CookID {
		return errUnauthorizedAccess.Wrapf("CookID(%s) does not have permission to change item(%d)", item.CookID, id)
	}
	if !user.IsAdmin() && !user.IsVerifiedCook() {
		return errUnauthorizedAccess.WithMessage("Only verifed cook are allowed to activate items. Get verfied asap! :)")
	}
	cookC := cook.New(c.ctx)
	cook, err := cookC.Get(item.CookID)
	if err != nil {
		return errors.Wrap("failed to cookC.Get", err)
	}
	err = insertOrUpdateActiveItem(c.ctx, id, item, cook.Address.Latitude, cook.Address.Longitude)
	if err != nil {
		return errors.Wrap("failed to insertOrUpdateActiveItem", err)
	}
	return nil
}

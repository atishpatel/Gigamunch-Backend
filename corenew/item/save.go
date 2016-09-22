package item

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

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
		// create item
		id, err = putIncomplete(c.ctx, item)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to putIncomplete Item(%v)", *item)
		}
	} else {
		// update item
		err = put(c.ctx, id, item)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to put Item(%d)", id)
		}
	}
	return id, item, nil
}

package menu

import "time"

// Save saves the Menu.
func (c *Client) Save(id int64, cookID, name string, color Color) (int64, *Menu, error) {
	if cookID == "" {
		return 0, nil, errInvalidParameter.Wrap("cook id cannot be empty")
	}
	var menu *Menu
	var err error
	if id == 0 {
		// create new menu
		menu = &Menu{CreatedDateTime: time.Now()}
	} else {
		// get menu
		menu, err = get(c.ctx, id)
		if err != nil {
			return 0, nil, errDatastore.WithError(err).Wrapf("failed to get Menu(%d)", id)
		}
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

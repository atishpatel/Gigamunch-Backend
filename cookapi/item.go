package cookapi

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// Menu is a Menu.
type Menu struct {
	ID int64 `json:"id,string"`
	menu.Menu
}

// SaveItemReq is the request for SaveItem.
type SaveItemReq struct {
	Gigatoken string `json:"gigatoken"`
	Item      Item   `json:"item"`
	Menu      Menu   `json:"menu"`
}

func (req *SaveItemReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *SaveItemReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// ItemOnlyResp is a response with only an Item and err.
type ItemOnlyResp struct {
	Item Item                 `json:"item"`
	Err  errors.ErrorWithCode `json:"err"`
}

// SaveItem creates or updates an Item. If the menu does not exist, it creates the menu
func (service *Service) SaveItem(ctx context.Context, req *SaveItemReq) (*ItemOnlyResp, error) {
	resp := new(ItemOnlyResp)
	defer handleResp(ctx, "SaveItem", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	menuID := req.Item.MenuID
	if menuID == 0 {
		// create new menu
		menuC := menu.New(ctx)
		menuID, _, err = menuC.Save(menuID, user.ID, req.Menu.Name, req.Menu.Color)
		if err != nil {
			resp.Err = errors.Wrap("failed to menuC.Save", err)
			return resp, nil
		}
	}

	itemC := item.New(ctx)
	id, item, err := itemC.Save(req.Item.ID, menuID, user.ID, req.Item.Title, req.Item.Description,
		req.Item.DietaryConcerns, req.Item.Ingredients, req.Item.Photos,
		req.Item.CookPricePerServing, req.Item.MinServings, req.Item.MaxServings)
	if err != nil {
		resp.Err = errors.Wrap("failed to itemC.Save", err)
		return resp, nil
	}
	resp.Item.ID = id
	resp.Item.Item = *item
	return resp, nil
}

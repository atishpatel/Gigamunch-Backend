package cookapi

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

// MenuWithItems is a Menu with all it's Items.
type MenuWithItems struct {
	ID int64 `json:"id,string"`
	menu.Menu
	Items []Item `json:"items"`
}

// GetMenusResp is the response for GetMenusResp.
type GetMenusResp struct {
	Menus []MenuWithItems      `json:"menus"`
	Err   errors.ErrorWithCode `json:"err"`
}

// GetMenus gets the Menus for a cook.
func (service *Service) GetMenus(ctx context.Context, req *GigatokenOnlyReq) (*GetMenusResp, error) {
	resp := new(GetMenusResp)
	defer handleResp(ctx, "GetMenus", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	menuC := menu.New(ctx)
	menuIDs, menus, err := menuC.GetCookMenus(user.ID)
	if err != nil {
		resp.Err = errors.Wrap("failed to menuC.GetCookMenus", err)
		return resp, nil
	}
	itemC := item.New(ctx)
	itemIDs, items, err := itemC.GetAllByCook(user.ID)
	if err != nil {
		resp.Err = errors.Wrap("failed to itemC.GetCookItems", err)
		return resp, nil
	}
	for i := range menuIDs {
		menu := MenuWithItems{ID: menuIDs[i], Menu: menus[i]}
		for j := len(itemIDs); j >= 0; j-- {
			if items[j].MenuID == menu.ID {
				item := Item{
					ID:   itemIDs[j],
					Item: items[j],
				}
				menu.Items = append(menu.Items, item)
				if j == 0 {
					itemIDs = itemIDs[j+1:]
					items = items[j+1:]
				} else {
					itemIDs = append(itemIDs[:j], itemIDs[j+1:]...)
					items = append(items[:j], items[j+1:]...)
				}
			}
		}
		resp.Menus = append(resp.Menus, menu)
	}
	if len(itemIDs) != 0 {
		// shouldn't every happen
		utils.Errorf(ctx, "GetMenus: cook(%s) has items(%v) that have no menu", user.ID, itemIDs)
	}
	return resp, nil
}

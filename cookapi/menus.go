package cookapi

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

// GetMenusResp is the response for GetMenusResp.
type GetMenusResp struct {
	Menus []MenuWithItems      `json:"menus"`
	Err   errors.ErrorWithCode `json:"err"`
}

// GetMenus gets the Menus for a cook.
func (service *Service) GetMenus(ctx context.Context, req *GigatokenReq) (*GetMenusResp, error) {
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
		menu := MenuWithItems{
			Menu: Menu{
				ID:   menuIDs[i],
				Menu: menus[i],
			},
		}
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

// // SaveMenuReq is the request for SaveMenu.
// type SaveMenuReq struct {
// 	GigatokenReq
// 	Menu Menu `json:"menu"`
// }

// type SaveMenuResp struct {
// 	Menu MenuWithItems `json:"menu"`
// 	ErrorOnlyResp
// }

// // SaveMenu creates or updates an Menu.
// func (service *Service) SaveMenu(ctx context.Context, req *SaveMenuReq) (*SaveMenuResp, error) {
// 	resp := new(SaveMenuResp)
// 	defer handleResp(ctx, "SaveMenu", resp.Err)
// 	user, err := validateRequestAndGetUser(ctx, req)
// 	if err != nil {
// 		resp.Err = errors.GetErrorWithCode(err)
// 		return resp, nil
// 	}

// 	// create new menu
// 	menuC := menu.New(ctx)
// 	menuID, menu, err := menuC.Save(user, req.Menu.ID, req.Menu.CookID, req.Menu.Name, req.Menu.Color)
// 	if err != nil {
// 		resp.Err = errors.Wrap("failed to menuC.Save", err)
// 		return resp, nil
// 	}

// 	return resp, nil
// }

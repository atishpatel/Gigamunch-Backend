package main

import (
	"golang.org/x/net/context"

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
	menus, err := menuC.GetCookMenus(user.ID)
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
	// get likes
	likes := getLikes(ctx, itemIDs)
	// set menus and items
	for i := range menus {
		menu := MenuWithItems{
			Menu: menus[i],
		}
		for j := len(itemIDs) - 1; j >= 0; j-- {
			if items[j].MenuID == menu.ID {
				item := Item{
					ID:       itemIDs[j],
					Item:     items[j],
					NumLikes: likes[j],
				}
				menu.Items = append(menu.Items, item)
				if j == 0 {
					itemIDs = itemIDs[j+1:]
					items = items[j+1:]
					likes = likes[j+1:]
				} else {
					itemIDs = append(itemIDs[:j], itemIDs[j+1:]...)
					items = append(items[:j], items[j+1:]...)
					likes = append(likes[:j], likes[j+1:]...)
				}
			}
		}
		resp.Menus = append(resp.Menus, menu)
	}
	if len(itemIDs) != 0 {
		// shouldn't every happen
		utils.Criticalf(ctx, "GetMenus: cook(%s) has items(%v) that have no menu", user.ID, itemIDs)
	}
	return resp, nil
}

// SaveMenuReq is the request for SaveMenu.
type SaveMenuReq struct {
	GigatokenReq
	Menu Menu `json:"menu"`
}

// MenuResp has a menu and error code
type MenuResp struct {
	Menu Menu `json:"menu"`
	ErrorOnlyResp
}

// SaveMenu creates or updates an Menu.
func (service *Service) SaveMenu(ctx context.Context, req *SaveMenuReq) (*MenuResp, error) {
	resp := new(MenuResp)
	defer handleResp(ctx, "SaveMenu", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}

	if req.Menu.CookID == "" {
		req.Menu.CookID = user.ID
	}

	// save menu or create new menu
	menuC := menu.New(ctx)
	menu, err := menuC.Save(user, req.Menu.ID, req.Menu.CookID, req.Menu.Name, req.Menu.Color)
	if err != nil {
		resp.Err = errors.Wrap("failed to menuC.Save", err)
		return resp, nil
	}
	resp.Menu.Menu = *menu
	return resp, nil
}

package main

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/like"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

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

// ItemResp is a response with only an Item and err.
type ItemResp struct {
	Item Item                 `json:"item"`
	Err  errors.ErrorWithCode `json:"err"`
}

// SaveItem creates or updates an Item. If the menu does not exist, it creates the menu
func (service *Service) SaveItem(ctx context.Context, req *SaveItemReq) (*ItemResp, error) {
	resp := new(ItemResp)
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
		menuID, _, err = menuC.Save(user, menuID, user.ID, req.Menu.Name, req.Menu.Color)
		if err != nil {
			resp.Err = errors.Wrap("failed to menuC.Save", err)
			return resp, nil
		}
	}

	itemC := item.New(ctx)
	id, item, err := itemC.Save(user, req.Item.ID, menuID, user.ID, req.Item.Title, req.Item.Description,
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

// GetItem gets an Item.
func (service *Service) GetItem(ctx context.Context, req *IDReq) (*ItemResp, error) {
	resp := new(ItemResp)
	defer handleResp(ctx, "GetItem", resp.Err)
	// user, err := validateRequestAndGetUser(ctx, req)
	// if err != nil {
	// 	resp.Err = errors.GetErrorWithCode(err)
	// 	return resp, nil
	// }
	itemC := item.New(ctx)
	item, err := itemC.Get(req.ID)
	if err != nil {
		resp.Err = errors.Wrap("failed to itemC.Get", err)
		return resp, nil
	}
	// get likes for item
	ids := []int64{req.ID}
	likes := getLikes(ctx, ids)
	// set item
	resp.Item.ID = req.ID
	resp.Item.Item = *item
	resp.Item.NumLikes = likes[0]
	return resp, nil
}

func getLikes(ctx context.Context, ids []int64) []int {
	likeC := like.New(ctx)
	likes, err := likeC.GetNumLikes(ids)
	if err != nil {
		utils.Errorf(ctx, "failed to likeC.GetNumLikes: %v", err)
		likes = make([]int, len(ids))
	}
	return likes
}

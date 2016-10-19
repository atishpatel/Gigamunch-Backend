package main

import (
	"golang.org/x/net/context"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/like"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
)

func (s *service) GetItem(ctx context.Context, id *pb.GetItemRequest) (resp *pb.GetItemResponse, unusedErr error) {
	defer handleResp(ctx, "GetItem", resp.Error)

	return
}

func (s *service) GetFeed(ctx context.Context, req *pb.GetFeedRequest) (resp *pb.GetFeedResponse, unusedErr error) {
	defer handleResp(ctx, "GetFeed", resp.Error)
	itemC := item.New(ctx)
	itemIDs, menuIDs, cookIDs, err := itemC.GetActiveItemIDs(req.StartIndex, req.EndIndex, req.Latitude, req.Longitude)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to itemC.GetActiveItemIDs")
		return
	}
	// get items
	var items []item.Item
	itemsErrChan := make(chan error, 1)
	go func() {
		var goErr error
		items, goErr = itemC.GetMulti(itemIDs)
		itemsErrChan <- goErr
	}()
	// get menus
	var menus map[int64]*menu.Menu
	menusErrChan := make(chan error, 1)
	go func() {
		var goErr error
		menuC := menu.New(ctx)
		menus, goErr = menuC.GetMulti(menuIDs)
		menusErrChan <- goErr
	}()
	// get cooks
	var cooks map[string]*cook.Cook
	cooksErrChan := make(chan error, 1)
	go func() {
		var goErr error
		cookC := cook.New(ctx)
		cooks, goErr = cookC.GetMulti(cookIDs)
		cooksErrChan <- goErr
	}()
	// get likes
	var likes []bool
	var numLikes []int32
	likeErrChan := make(chan error, 1)
	go func() {
		// get user if there
		var userID string
		if req.Gigatoken != "" {
			user, _ := auth.GetUserFromToken(ctx, req.Gigatoken)
			if user != nil {
				userID = user.ID
			}
		}
		var goErr error
		likeC := like.New(ctx)
		likes, numLikes, goErr = likeC.LikesItems(userID, itemIDs)
		likeErrChan <- goErr
	}()
	// handle errors
	err = processErrorChans(itemsErrChan, menusErrChan, cooksErrChan, likeErrChan)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to item.GetMulti or menu.GetMulti or cook.GetMulti")
		return
	}
	// get menu order
	menuOrder := make([]int64, len(menus))
	index := 0
	for i := range items {
		found := false
		for _, v := range menuOrder {
			if v == items[i].MenuID {
				found = true
			}
		}
		if !found {
			menuOrder[index] = items[i].MenuID
			index++
		}
	}
	// set menus
	resp.Menus = make([]*pb.Menu, len(menus))
	for i, v := range menuOrder {
		// TODO add cookDistance and exchangeoptions
		resp.Menus[i] = getPBMenu(menus[v], cooks[menus[v].CookID], 0, nil)
		menuID := menus[v].ID
		for i := range items {
			if items[i].MenuID == menuID {
				resp.Menus[i].Items = append(resp.Menus[i].Items, getPBBaseItem(&items[i], numLikes[i], likes[i]))
			}
		}
	}
	return
}

func (s *service) LikeItem(ctx context.Context, req *pb.LikeItemRequest) (resp *pb.ErrorOnlyResponse, unusedErr error) {
	defer handleResp(ctx, "LikeItem", resp.Error)
	user, validateErr := validateLikeItemRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return
	}
	likeC := like.New(ctx)
	var err error
	if req.Like {
		err = likeC.Like(user.ID, req.ItemId, req.MenuId, req.CookId)
	} else {
		err = likeC.Unlike(user.ID, req.ItemId)
	}
	if err != nil {
		resp.Error = getGRPCError(err, "failed to like or unlike")
		return
	}
	return
}

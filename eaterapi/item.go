package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/like"
	"github.com/atishpatel/Gigamunch-Backend/corenew/maps"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/corenew/review"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

func (s *service) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetItemResponse)
	defer handleResp(ctx, "GetItem", resp.Error)
	validateErr := validateGetItemRequest(req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	likeC := like.New(ctx)
	// get item
	itemC := item.New(ctx)
	item, err := itemC.Get(req.ItemId)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to item.Get")
		return resp, nil
	}
	// get cook
	var c *cook.Cook
	cooksErrChan := make(chan error, 1)
	go func() {
		var goErr error
		cookC := cook.New(ctx)
		c, goErr = cookC.Get(item.CookID)
		cooksErrChan <- goErr
	}()
	// get reviews
	var reviews []*review.Review
	reviewErrChan := make(chan error, 1)
	go func() {
		var goErr error
		reviewC := review.New(ctx)
		reviews, goErr = reviewC.GetByCookID(item.CookID, item.ID, 0, 5)
		reviewErrChan <- goErr
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
		likes, numLikes, goErr = likeC.LikesItems(userID, []int64{item.ID})
		likeErrChan <- goErr
	}()
	// get cook likes
	var cookLikes int32
	cookLikeErrChan := make(chan error, 1)
	go func() {
		var goErr error
		cookLikes, goErr = likeC.GetNumCookLikes(item.CookID)
		cookLikeErrChan <- goErr
	}()
	// handle errors
	err = processErrorChans(cooksErrChan, reviewErrChan, likeErrChan, cookLikeErrChan)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to cook.Get or like.LikesItems")
		return resp, nil
	}
	eaterPoint := types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	cookPoint := c.Address.GeoPoint
	// get distance
	distance, _, err := maps.GetDistance(ctx, cookPoint, eaterPoint)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to maps.GetDistance")
		return resp, nil
	}
	// get exchangeoptions
	ems := types.GetExchangeMethods(cookPoint, c.DeliveryRange, c.DeliveryPrice, eaterPoint)
	resp.Item = getPBItem(item, numLikes[0], likes[0], c, distance, ems, cookLikes, reviews)
	return resp, nil
}

func (s *service) GetFeed(ctx context.Context, req *pb.GetFeedRequest) (*pb.GetFeedResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetFeedResponse)
	defer handleResp(ctx, "GetFeed", resp.Error)
	itemC := item.New(ctx)
	itemIDs, menuIDs, cookIDs, err := itemC.GetActiveItemIDs(req.StartIndex, req.EndIndex, req.Latitude, req.Longitude)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to itemC.GetActiveItemIDs")
		return resp, nil
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
		return resp, nil
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
	eaterPoint := types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	resp.Menus = make([]*pb.Menu, len(menus))
	for i, v := range menuOrder {
		c := cooks[menus[v].CookID] // cook for this menu
		// get exchangeoptions and distance
		ems := types.GetExchangeMethods(c.Address.GeoPoint, c.DeliveryRange, c.DeliveryPrice, eaterPoint)
		distance := eaterPoint.GreatCircleDistance(c.Address.GeoPoint)
		resp.Menus[i] = getPBMenu(menus[v], c, distance, ems)
		menuID := menus[v].ID
		for i := range items {
			if items[i].MenuID == menuID {
				resp.Menus[i].Items = append(resp.Menus[i].Items, getPBBaseItem(&items[i], numLikes[i], likes[i]))
			}
		}
	}
	return resp, nil
}

func (s *service) LikeItem(ctx context.Context, req *pb.LikeItemRequest) (*pb.ErrorOnlyResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.ErrorOnlyResponse)
	defer handleResp(ctx, "LikeItem", resp.Error)
	user, validateErr := validateLikeItemRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
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
		return resp, nil
	}
	return resp, nil
}

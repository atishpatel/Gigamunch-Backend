package main

import (
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	itemold "github.com/atishpatel/Gigamunch-Backend/core/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	itemnew "github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	weekSchedule = []cook.WeekSchedule{
		cook.WeekSchedule{
			DayOfWeek: 0,
			StartTime: 6 * 60 * 60,
			EndTime:   86399,
		},
		cook.WeekSchedule{
			DayOfWeek: 1,
			StartTime: 6 * 60 * 60,
			EndTime:   (8 * 60 * 60) - 1,
		},
		cook.WeekSchedule{
			DayOfWeek: 1,
			StartTime: 17 * 60 * 60,
			EndTime:   86399,
		},
		cook.WeekSchedule{
			DayOfWeek: 2,
			StartTime: 6 * 60 * 60,
			EndTime:   (8 * 60 * 60) - 1,
		},
		cook.WeekSchedule{
			DayOfWeek: 2,
			StartTime: 17 * 60 * 60,
			EndTime:   86399,
		},
		cook.WeekSchedule{
			DayOfWeek: 3,
			StartTime: 6 * 60 * 60,
			EndTime:   (8 * 60 * 60) - 1,
		},
		cook.WeekSchedule{
			DayOfWeek: 3,
			StartTime: 17 * 60 * 60,
			EndTime:   86399,
		},
		cook.WeekSchedule{
			DayOfWeek: 4,
			StartTime: 6 * 60 * 60,
			EndTime:   (8 * 60 * 60) - 1,
		},
		cook.WeekSchedule{
			DayOfWeek: 4,
			StartTime: 17 * 60 * 60,
			EndTime:   86399,
		},
		cook.WeekSchedule{
			DayOfWeek: 5,
			StartTime: 6 * 60 * 60,
			EndTime:   (8 * 60 * 60) - 1,
		},
		cook.WeekSchedule{
			DayOfWeek: 5,
			StartTime: 17 * 60 * 60,
			EndTime:   86399,
		},
		cook.WeekSchedule{
			DayOfWeek: 6,
			StartTime: 6 * 60 * 60,
			EndTime:   86399,
		},
	}
)

func (service *Service) CopyOverCooks(ctx context.Context, req *GigatokenReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "CopyOverCooks", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err.Message = "User is not admin."
		return resp, nil
	}
	kindOld := "Gigachef"
	kindNew := "Cook"
	query := datastore.NewQuery(kindOld).Limit(1000)
	var results []gigachef.Gigachef
	oldKeys, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	cooks := make([]cook.Cook, len(oldKeys))
	newKeys := make([]*datastore.Key, len(oldKeys))
	for i := range oldKeys {
		newKeys[i] = datastore.NewKey(ctx, kindNew, oldKeys[i].StringID(), 0, nil)
		cooks[i].CreatedDatetime = results[i].CreatedDatetime
		cooks[i].UserDetail = results[i].UserDetail
		cooks[i].UserDetail.PhotoURL = strings.Replace(cooks[i].UserDetail.PhotoURL, "s96-c", "s250-c", 0)
		cooks[i].Bio = results[i].Bio
		cooks[i].PhoneNumber = results[i].PhoneNumber
		cooks[i].Address = results[i].Address
		cooks[i].DeliveryRange = results[i].DeliveryRange
		cooks[i].WeekSchedule = weekSchedule
		cooks[i].Rating.AverageRating = results[i].Rating.AverageRating
		cooks[i].Rating.NumRatings = results[i].Rating.NumRatings
		cooks[i].Rating.NumOneStarRatings = results[i].Rating.NumOneStarRatings
		cooks[i].Rating.NumTwoStarRatings = results[i].Rating.NumTwoStarRatings
		cooks[i].Rating.NumThreeStarRatings = results[i].Rating.NumThreeStarRatings
		cooks[i].Rating.NumFourStarRatings = results[i].Rating.NumFourStarRatings
		cooks[i].Rating.NumFiveStarRatings = results[i].Rating.NumFiveStarRatings
		cooks[i].SubMerchantStatus = results[i].SubMerchantStatus
		cooks[i].KitchenPhotoURLs = results[i].KitchenPhotoURLs
		cooks[i].KitchenInspection = results[i].KitchenInspection
		cooks[i].BackgroundCheck = results[i].BackgroundCheck
		cooks[i].FoodHandlerCard = results[i].FoodHandlerCard
		cooks[i].Verified = results[i].Verified

	}

	_, err = datastore.PutMulti(ctx, newKeys, cooks)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *Service) CopyOverItems(ctx context.Context, req *GigatokenReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "CopyOverItems", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err.Message = "User is not admin."
		return resp, nil
	}
	kindOld := "Item"
	kindNew := "MenuItem"

	query := datastore.NewQuery(kindOld).Limit(1000)
	var results []itemold.Item
	oldKeys, err := query.GetAll(ctx, &results)
	if err != nil {
		resp.Err.WithError(err)
		return resp, nil
	}
	menuC := menu.New(ctx)

	cookMenuIDs := make(map[string]int64)
	items := make([]itemnew.Item, len(oldKeys))
	newKeys := make([]*datastore.Key, len(oldKeys))
	for i := range oldKeys {
		newKeys[i] = datastore.NewKey(ctx, kindNew, "", oldKeys[i].IntID(), nil)
		items[i].CreatedDateTime = results[i].CreatedDateTime
		items[i].CookID = results[i].GigachefID
		menuID, ok := cookMenuIDs[results[i].GigachefID]
		if !ok {
			menuID, _, err = menuC.Save(user, 0, results[i].GigachefID, "", menu.NewColor())
			if err != nil {
				return nil, err
			}
			cookMenuIDs[results[i].GigachefID] = menuID
		}
		items[i].MenuID = menuID
		items[i].Title = results[i].Title
		items[i].Description = results[i].Description
		items[i].Ingredients = results[i].Ingredients
		items[i].Photos = results[i].Photos
		items[i].MaxServings = 20
		items[i].MinServings = 1
	}

	_, err = datastore.PutMulti(ctx, newKeys, items)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

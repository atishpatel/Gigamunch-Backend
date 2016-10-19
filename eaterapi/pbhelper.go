package main

import (
	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
)

func getPBMenu(menu *menu.Menu, cook *cook.Cook, cookDistance float64, exchangeOptions []*pb.ExchangeOption) *pb.Menu {
	return &pb.Menu{
		Id:    menu.ID,
		Name:  menu.Name,
		Color: menu.Color.HexValue(),
		Cook:  getPBBaseCook(cook, cookDisatance, totalCookLikes),
	}
}

func getPBBaseCook(cook *cook.Cook, distance float64, exchangeOptions []*pb.ExchangeOption) *pb.BaseCook {
	return &pb.BaseCook{
		Id:              cook.ID,
		Name:            cook.Name,
		Image:           cook.PhotoURL,
		NumRatings:      cook.NumRatings,
		Rating:          cook.Rating,
		Distance:        distance,
		ExchangeOptions: exchangeOptions,
	}
}

func getPBCook(cook *cook.Cook, distance float64, exchangeOptions []*pb.ExchangeOption, totalLikes int32) *pb.Cook {
	c := &pb.Cook{
		Id:              cook.ID,
		Name:            cook.Name,
		Image:           cook.PhotoURL,
		Distance:        distance,
		ExchangeOptions: exchangeOptions,
		Latitude:        cook.Address.Latitude,
		Longitude:       cook.Address.Longitude,
		KitchenImages:   cook.KitchenPhotoURLs,
		Bio:             cook.Bio,
		InstagramHandle: cook.InstagramID,
		RatingStats: &pb.CookRatingStats{
			NumRatings:          cook.NumRatings,
			Rating:              cook.AverageRating,
			NumOneStarRatings:   cook.NumOneStarRatings,
			NumTwoStarRatings:   cook.NumTwoStarRatings,
			NumThreeStarRatings: cook.NumThreeStarRatings,
			NumFourStarRatings:  cook.NumFourStarRatings,
			NumFiveStarRatings:  cook.NumFiveStarRatings,
			TotalLikes:          totalLikes,
		},
	}
	c.Availability = &pb.Availability{}
	for _, a := range cook.WeekSchedule {
		c.Availability.WeekSchedule = append(c.Availability.WeekSchedule, &pb.AvailabilityWindow{
			DayOfWeek: a.DayOfWeek,
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
		})
	}
	for _, m := range cook.ScheduleModifications {
		c.Availability.ScheduleModifications = append(c.Availability.ScheduleModifications, &pb.ScheduleModification{
			StartDatetime: m.StartDateTime,
			EndDatetime:   m.EndDateTime,
			IsAvailable:   m.Available,
		})
	}
	return c
}

func getPBItem(item *item.Item, numLikes int32, hasLiked bool) *pb.Item {
	return &pb.Item{}
}

func getPBBaseItem(item *item.Item, numLikes int32, hasLiked bool) *pb.BaseItem {
	return &pb.BaseItem{}
}

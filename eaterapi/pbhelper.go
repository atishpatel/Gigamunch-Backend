package main

import (
	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/corenew/review"

	"github.com/golang/protobuf/ptypes"
)

func getPBMenu(menu *menu.Menu, cook *cook.Cook, cookDistance float32, exchangeOptions []*pb.ExchangeOption) *pb.Menu {
	return &pb.Menu{
		Id:    menu.ID,
		Name:  menu.Name,
		Color: menu.Color.HexValue(),
		Cook:  getPBBaseCook(cook, cookDistance, exchangeOptions),
	}
}

func getPBBaseCook(cook *cook.Cook, distance float32, exchangeOptions []*pb.ExchangeOption) *pb.BaseCook {
	return &pb.BaseCook{
		Id:              cook.ID,
		Name:            cook.Name,
		Image:           cook.PhotoURL,
		NumRatings:      cook.NumRatings,
		Rating:          cook.AverageRating,
		Distance:        distance,
		ExchangeOptions: exchangeOptions,
	}
}

func getPBCook(cook *cook.Cook, distance float32, exchangeOptions []*pb.ExchangeOption, totalLikes int32) *pb.Cook {
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
		startTimeStamp, _ := ptypes.TimestampProto(m.StartDateTime)
		endTimeStamp, _ := ptypes.TimestampProto(m.EndDateTime)
		c.Availability.ScheduleModifications = append(c.Availability.ScheduleModifications, &pb.ScheduleModification{
			StartDatetime: startTimeStamp,
			EndDatetime:   endTimeStamp,
			IsAvailable:   m.Available,
		})
	}
	return c
}

// TODO add reviews
func getPBItem(item *item.Item, numLikes int32, hasLiked bool, cook *cook.Cook, distance float32, exchangeOptions []*pb.ExchangeOption, cookLikes int32, reviews []*review.Review) *pb.Item {
	i := &pb.Item{
		Id:              item.ID,
		MenuId:          item.MenuID,
		Name:            item.Title,
		Description:     item.Description,
		DietaryConcerns: int64(item.DietaryConcerns),
		PricePerServing: item.CookPricePerServing * 1.2, // TODO move out of here
		MinServings:     item.MinServings,
		MaxServings:     item.MaxServings,
		Images:          item.Photos,
		NumOrdersSold:   item.NumOrdersSold,
		NumServingsSold: item.NumServingsSold,
		NumLikes:        numLikes,
		HasLiked:        hasLiked,
		Ingredients:     item.Ingredients,
		PriceInfo: &pb.PriceInfo{
			CookPricePerServing:  item.CookPricePerServing,
			ServiceFeePercentage: .2,
			TaxPercentage:        7.25,
		},
		Cook: getPBCook(cook, distance, exchangeOptions, cookLikes),
	}
	i.CreatedDatetime, _ = ptypes.TimestampProto(item.CreatedDateTime)
	i.Reviews = getPBReviews(reviews)
	return i
}

func getPBBaseItem(item *item.Item, numLikes int32, hasLiked bool) *pb.BaseItem {
	i := &pb.BaseItem{
		Id:              item.ID,
		MenuId:          item.MenuID,
		Name:            item.Title,
		Description:     item.Description,
		DietaryConcerns: int64(item.DietaryConcerns),
		PricePerServing: item.CookPricePerServing * 1.2, // TODO move out of here
		MinServings:     item.MinServings,
		MaxServings:     item.MaxServings,
		Images:          item.Photos,
		NumOrdersSold:   item.NumOrdersSold,
		NumServingsSold: item.NumServingsSold,
		NumLikes:        numLikes,
		HasLiked:        hasLiked,
	}
	i.CreatedDatetime, _ = ptypes.TimestampProto(item.CreatedDateTime)
	return i
}

func getPBReviews(reviews []*review.Review) []*pb.Review {
	r := make([]*pb.Review, len(reviews))
	for i := range reviews {
		r[i] = getPBReview(reviews[i])
	}
	return r
}

func getPBReview(review *review.Review) *pb.Review {
	createdTimestamp, _ := ptypes.TimestampProto(review.CreatedDateTime)
	editedTimestamp, _ := ptypes.TimestampProto(review.EditedDateTime)
	responseCreatedTimestamp, _ := ptypes.TimestampProto(review.ResponseCreatedDateTime)
	return &pb.Review{
		Id:                      review.ID,
		CookId:                  review.CookID,
		EaterId:                 review.EaterID,
		InquiryId:               review.InquiryID,
		ItemId:                  review.ItemID,
		EaterName:               review.EaterName,
		EaterImage:              review.EaterPhotoURL,
		CreatedDatetime:         createdTimestamp,
		IsEdited:                review.IsEdited,
		EditedDatetime:          editedTimestamp,
		Rating:                  review.Rating,
		Text:                    review.Text,
		HasResponse:             review.HasResponse,
		ResponseCreatedDatetime: responseCreatedTimestamp,
		ResponseText:            review.ResponseText,
		ItemImage:               review.ItemPhotoURL,
		ItemName:                review.ItemName,
	}
}

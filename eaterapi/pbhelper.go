package main

import (
	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/review"
	"github.com/atishpatel/Gigamunch-Backend/types"

	"github.com/atishpatel/Gigamunch-Backend/corenew/promo"
	"github.com/golang/protobuf/ptypes"
)

func getPBPromo(code *promo.Code) *pb.Promo {
	return &pb.Promo{
		FreeDelivery:     code.FreeDelivery,
		FreeDish:         code.FreeDish,
		PercentOff:       code.PercentOff,
		AmountOff:        code.AmountOff,
		BuyOneGetOneFree: code.BuyOneGetOneFree,
		DiscountCap:      code.DiscountCap,
	}
}

func getPBMenu(menu *menu.Menu, cook *cook.Cook, cookDistance float32, ems []types.ExchangeMethodWithPrice) *pb.Menu {
	return &pb.Menu{
		Id:    menu.ID,
		Name:  menu.Name,
		Color: menu.Color.HexValue(),
		Cook:  getPBBaseCook(cook, cookDistance, ems),
	}
}

func getPBBaseCook(cook *cook.Cook, distance float32, ems []types.ExchangeMethodWithPrice) *pb.BaseCook {
	return &pb.BaseCook{
		Id:              cook.ID,
		Name:            cook.Name,
		Image:           cook.PhotoURL,
		NumRatings:      cook.NumRatings,
		Rating:          cook.AverageRating,
		Distance:        distance,
		ExchangeOptions: getPBExchangeOptions(ems),
	}
}

func getPBCook(cook *cook.Cook, distance float32, ems []types.ExchangeMethodWithPrice, totalLikes int32) *pb.Cook {
	c := &pb.Cook{
		Id:              cook.ID,
		Name:            cook.Name,
		Image:           cook.PhotoURL,
		Distance:        distance,
		ExchangeOptions: getPBExchangeOptions(ems),
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

func getPBItem(item *item.Item, numLikes int32, hasLiked bool, cook *cook.Cook, distance float32, ems []types.ExchangeMethodWithPrice, cookLikes int32, reviews []*review.Review) *pb.Item {
	pricePerServing := payment.GetPricePerServing(item.CookPricePerServing)
	i := &pb.Item{
		Id:              item.ID,
		MenuId:          item.MenuID,
		Name:            item.Title,
		Description:     item.Description,
		DietaryConcerns: int64(item.DietaryConcerns),
		PricePerServing: pricePerServing,
		ServingSize:     item.ServingDescription,
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
			ServiceFeePrice:      pricePerServing - item.CookPricePerServing,
			ServiceFeePercentage: pricePerServing / item.CookPricePerServing,
			TaxPercentage:        payment.GetTaxPercentage(cook.Address.Latitude, cook.Address.Longitude),
		},
		Cook: getPBCook(cook, distance, ems, cookLikes),
	}
	i.CreatedDatetime, _ = ptypes.TimestampProto(item.CreatedDateTime)
	i.Reviews = getPBReviews(reviews)
	return i
}

func getPBBaseItem(item *item.Item, numLikes int32, hasLiked bool) *pb.BaseItem {
	pricePerServing := payment.GetPricePerServing(item.CookPricePerServing)
	i := &pb.BaseItem{
		Id:              item.ID,
		MenuId:          item.MenuID,
		Name:            item.Title,
		Description:     item.Description,
		DietaryConcerns: int64(item.DietaryConcerns),
		PricePerServing: pricePerServing,
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

func getPBInquiry(inq *inquiry.Inquiry, cookName, cookImage string, menuID int64, numLikes int32, hasLiked bool) *pb.Inquiry {
	createdDatetime, _ := ptypes.TimestampProto(inq.CreatedDateTime)
	expectedExchangeDatetime, _ := ptypes.TimestampProto(inq.ExpectedExchangeDateTime)
	return &pb.Inquiry{
		Id:                    inq.ID,
		CookId:                inq.CookID,
		EaterId:               inq.EaterID,
		ReviewId:              inq.ReviewID,
		ItemId:                inq.ItemID,
		BtTransactionId:       inq.BTTransactionID,
		BtRefundTransactionId: inq.BTRefundTransactionID,
		CreatedDatetime:       createdDatetime,
		Item: &pb.InquiryItem{
			Name:            inq.Item.Name,
			Description:     inq.Item.Description,
			DietaryConcerns: int64(inq.Item.DietaryConcerns),
			Images:          inq.Item.Photos,
			Ingredients:     inq.Item.Ingredients,
			MenuId:          menuID,
			NumLikes:        numLikes,
			HasLiked:        hasLiked,
		},
		Promo: &pb.Promo{
			FreeDelivery:     inq.Promo.FreeDelivery,
			FreeDish:         inq.Promo.FreeDish,
			PercentOff:       inq.Promo.PercentOff,
			AmountOff:        inq.Promo.AmountOff,
			BuyOneGetOneFree: inq.Promo.BuyOneGetOneFree,
			DiscountCap:      inq.Promo.DiscountCap,
		},
		ExpectedExchangeDatetime: expectedExchangeDatetime,
		EaterImage:               inq.EaterPhotoURL,
		EaterName:                inq.EaterName,
		CookName:                 cookName,
		CookImage:                cookImage,
		Servings:                 inq.Servings,
		TotalPriceInfo: &pb.TotalPriceInfo{
			CookPricePerServing:     inq.PaymentInfo.CookPricePerServing,
			PricePerServing:         inq.PaymentInfo.PricePerServing,
			TotalCookPrice:          inq.PaymentInfo.CookPrice,
			ExchangePrice:           inq.PaymentInfo.ExchangePrice,
			TaxPrice:                inq.PaymentInfo.TaxPrice,
			ServiceFeePrice:         inq.PaymentInfo.ServiceFee,
			AmountOff:               inq.PaymentInfo.AmountOff,
			AmountOffDelivery:       inq.PaymentInfo.AmountOffDelivery,
			TotalPrice:              inq.PaymentInfo.TotalPrice,
			TotalPriceWithAmountOff: inq.PaymentInfo.TotalPriceWithAmountOff,
		},
		ExchangePlanInfo: &pb.ExchangePlanInfo{
			MethodName:   inq.ExchangeMethod.String(),
			IsDelivery:   inq.ExchangeMethod.Delivery(),
			EaterAddress: getPBAddress(&inq.ExchangePlanInfo.EaterAddress, false),
			CookAddress:  getPBAddress(&inq.ExchangePlanInfo.CookAddress, false),
			Distance:     inq.ExchangePlanInfo.Distance,
			Duration:     int32(inq.ExchangePlanInfo.Duration),
		},
		State:       inq.State,
		EaterAction: inq.EaterAction,
		CookAction:  inq.CookAction,
		HasIssue:    inq.Issue,
	}
}

func getPBInquiries(inqs []*inquiry.Inquiry, cookNames, cookImages []string, menuIDs []int64, numLikes []int32, hasLiked []bool) []*pb.Inquiry {
	l := len(inqs)
	is := make([]*pb.Inquiry, l)
	if len(cookNames) != l {
		cookNames = make([]string, l)
	}
	if len(cookImages) != l {
		cookImages = make([]string, l)
	}
	for i := range inqs {
		is[i] = getPBInquiry(inqs[i], cookNames[i], cookImages[i], menuIDs[i], numLikes[i], hasLiked[i])
	}
	return is
}

func getPBAddress(addr *types.Address, selected bool) *pb.Address {
	return &pb.Address{
		Country:    addr.Country,
		State:      addr.State,
		City:       addr.City,
		Zip:        addr.Zip,
		Street:     addr.Street,
		UnitNumber: addr.APT,
		Latitude:   addr.Latitude,
		Longitude:  addr.Longitude,
		IsSelected: selected,
	}
}

func getAddress(addr *pb.Address) *types.Address {
	return &types.Address{
		Country: addr.Country,
		State:   addr.State,
		City:    addr.City,
		Zip:     addr.Zip,
		Street:  addr.Street,
		APT:     addr.UnitNumber,
		GeoPoint: types.GeoPoint{
			Latitude:  addr.Latitude,
			Longitude: addr.Longitude,
		},
	}
}

func getPBExchangeOptions(ems []types.ExchangeMethodWithPrice) []*pb.ExchangeOption {
	pbems := make([]*pb.ExchangeOption, len(ems))
	for i := range ems {
		pbems[i] = &pb.ExchangeOption{
			Id:         ems[i].ID(),
			Name:       ems[i].String(),
			IsDelivery: ems[i].Delivery(),
			Price:      ems[i].Price,
		}
	}
	return pbems
}

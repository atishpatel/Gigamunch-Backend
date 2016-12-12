package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/like"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/review"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

func (s *service) MakeInquiry(ctx context.Context, req *pb.MakeInquiryRequest) (*pb.MakeInquiryResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.MakeInquiryResponse)
	defer handleResp(ctx, "MakeInquiry", resp.Error)
	user, exchangeTime, validateErr := validateMakeInquiryRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	exchangeMethod := types.ExchangeMethod(req.ExchangeId)
	inquiryC := inquiry.New(ctx)
	inq, err := inquiryC.Make(req.ItemId, req.BraintreeNonce, user.ID, getAddress(req.Address), req.Servings, exchangeMethod, exchangeTime, req.PromoCode)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to inquiry.Make")
		return resp, nil
	}
	// get cook name and photo
	cookC := cook.New(ctx)
	names, photos, err := cookC.GetMultiNamesAndPhotos([]string{inq.CookID})
	if err != nil {
		utils.Errorf(ctx, "failed to cook.GetMultiNamesAndPhotos. Err: %v", err)
		names = make([]string, 1)
		photos = make([]string, 1)
	}
	likeC := like.New(ctx)
	hasLiked, numLikes, menuIDs, err := likeC.GetNumLikesWithMenuID(user.ID, []int64{inq.ItemID})
	if err != nil {
		utils.Errorf(ctx, "failed to likeC.GetNumLikesWithMenuID. Err: %v", err)
		hasLiked = make([]bool, 1)
		numLikes = make([]int32, 1)
		menuIDs = make([]int64, 1)
	}
	resp.Inquiry = getPBInquiry(inq, names[0], photos[0], menuIDs[0], numLikes[0], hasLiked[0])
	return resp, nil
}

func (s *service) GetInquiries(ctx context.Context, req *pb.GetInquiriesRequest) (*pb.GetInquiriesResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetInquiriesResponse)
	defer handleResp(ctx, "GetInquiries", resp.Error)
	user, validateErr := validateGetInquiriesRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inqs, err := inquiryC.GetByEaterID(user.ID, int(req.StartIndex), int(req.EndIndex))
	if err != nil {
		resp.Error = getGRPCError(err, "")
		return resp, nil
	}
	// get cook names and photos
	cookC := cook.New(ctx)
	cookIDs := make([]string, len(inqs))
	for i := range inqs {
		cookIDs[i] = inqs[i].CookID
	}
	names, photos, err := cookC.GetMultiNamesAndPhotos(cookIDs)
	if err != nil {
		utils.Errorf(ctx, "failed to cook.GetMultiNamesAndPhotos. Err: %v", err)
		names = make([]string, len(inqs))
		photos = make([]string, len(inqs))
	}
	// get likes and menuID
	itemIDs := make([]int64, len(inqs))
	for i := range inqs {
		itemIDs[i] = inqs[i].ItemID
	}
	likeC := like.New(ctx)
	hasLiked, numLikes, menuIDs, err := likeC.GetNumLikesWithMenuID(user.ID, itemIDs)
	if err != nil {
		utils.Errorf(ctx, "failed to likeC.GetNumLikesWithMenuID. Err: %v", err)
	}
	resp.Inquiries = getPBInquiries(inqs, names, photos, menuIDs, numLikes, hasLiked)
	return resp, nil
}

func (s *service) GetInquiry(ctx context.Context, req *pb.GetInquiryRequest) (*pb.GetInquiryResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetInquiryResponse)
	defer handleResp(ctx, "GetInquiry", resp.Error)
	user, validateErr := validateIDReq(ctx, req.InquiryId, req.Gigatoken)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inq, err := inquiryC.Get(user, req.InquiryId)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to inquiry.Get")
		return resp, nil
	}
	// get cook name and photo
	cookC := cook.New(ctx)
	names, photos, err := cookC.GetMultiNamesAndPhotos([]string{inq.CookID})
	if err != nil {
		utils.Errorf(ctx, "failed to cook.GetMultiNamesAndPhotos. Err: %v", err)
		names = make([]string, 1)
		photos = make([]string, 1)
	}
	// get likes
	likeC := like.New(ctx)
	hasLiked, numLikes, menuIDs, err := likeC.GetNumLikesWithMenuID(user.ID, []int64{inq.ItemID})
	if err != nil {
		utils.Errorf(ctx, "failed to likeC.GetNumLikesWithMenuID. Err: %v", err)
		hasLiked = make([]bool, 1)
		numLikes = make([]int32, 1)
		menuIDs = make([]int64, 1)
	}
	// get review
	reviewC := review.New(ctx)
	rvw, err := reviewC.Get(inq.ReviewID)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to review.Get")
	} else {
		resp.Review = getPBReview(rvw)
	}
	resp.Inquiry = getPBInquiry(inq, names[0], photos[0], menuIDs[0], numLikes[0], hasLiked[0])
	return resp, nil
}

func (s *service) CancelInquiry(ctx context.Context, req *pb.CancelInquiryRequest) (*pb.CancelInquiryResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.CancelInquiryResponse)
	defer handleResp(ctx, "CancelInquiry", resp.Error)
	user, validateErr := validateIDReq(ctx, req.InquiryId, req.Gigatoken)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inq, err := inquiryC.EaterCancel(user, req.InquiryId)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to inquiry.EaterCancel")
		return resp, nil
	}
	// get cook name and photo
	cookC := cook.New(ctx)
	names, photos, err := cookC.GetMultiNamesAndPhotos([]string{inq.CookID})
	if err != nil {
		utils.Errorf(ctx, "failed to cook.GetMultiNamesAndPhotos. Err: %v", err)
		names = make([]string, 1)
		photos = make([]string, 1)
	}
	// get likes
	likeC := like.New(ctx)
	hasLiked, numLikes, menuIDs, err := likeC.GetNumLikesWithMenuID(user.ID, []int64{inq.ItemID})
	if err != nil {
		utils.Errorf(ctx, "failed to likeC.GetNumLikesWithMenuID. Err: %v", err)
		hasLiked = make([]bool, 1)
		numLikes = make([]int32, 1)
		menuIDs = make([]int64, 1)
	}
	resp.Inquiry = getPBInquiry(inq, names[0], photos[0], menuIDs[0], numLikes[0], hasLiked[0])
	return resp, nil
}

func (s *service) GetBraintreeToken(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.GetBraintreeTokenResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetBraintreeTokenResponse)
	defer handleResp(ctx, "GetBraintreeToken", resp.Error)
	user, validateErr := validateGigatokenAndGetUser(ctx, req.Gigatoken)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	eaterC := eater.New(ctx)
	customerID, err := eaterC.GetBTCustomerID(user.ID)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to eater.GetBTCustomerID")
		return resp, nil
	}
	paymentC := payment.New(ctx)
	token, err := paymentC.GenerateToken(customerID)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to payment.GenerateToken")
		return resp, nil
	}
	resp.BraintreeToken = token
	return resp, nil
}

func (s *service) CheckDeliveryAddresses(ctx context.Context, req *pb.CheckDeliveryAddressesRequest) (*pb.CheckDeliveryAddressesResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.CheckDeliveryAddressesResponse)
	defer handleResp(ctx, "CheckDeliveryAddresses", resp.Error)
	cookC := cook.New(ctx)
	c, err := cookC.Get(req.CookId)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to cookC.Get")
		return resp, nil
	}
	resp.Addresses = make([]*pb.DeliveryAddress, len(req.Addresses))
	for i := range req.Addresses {
		eaterPoint := types.GeoPoint{Latitude: req.Addresses[i].Latitude, Longitude: req.Addresses[i].Longitude}
		exchangeOptions := types.GetExchangeMethods(c.Address.GeoPoint, c.DeliveryRange, c.DeliveryPrice, eaterPoint)
		var cheapestExchangeOption *types.ExchangeMethodWithPrice
		for _, v := range exchangeOptions {
			if v.Delivery() {
				if cheapestExchangeOption == nil || v.Price < cheapestExchangeOption.Price {
					cheapestExchangeOption = &v
				}
			}
		}
		resp.Addresses[i] = &pb.DeliveryAddress{
			Address: req.Addresses[i],
		}
		if cheapestExchangeOption != nil {
			resp.Addresses[i].Available = true
			resp.Addresses[i].ExchangeOption = &pb.ExchangeOption{
				Id:         cheapestExchangeOption.ID(),
				Name:       cheapestExchangeOption.String(),
				IsDelivery: cheapestExchangeOption.Delivery(),
				Price:      cheapestExchangeOption.Price,
			}
		}
	}
	return resp, nil
}

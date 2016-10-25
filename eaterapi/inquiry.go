package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
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
	// TODO figure out exchange info
	exchangeMethod := types.ExchangeMethods(req.ExchangeId)
	var exchangePrice float32

	inquiryC := inquiry.New(ctx)
	inq, err := inquiryC.Make(req.ItemId, req.BraintreeNonce, user.ID, getAddress(req.Address), req.Servings, exchangeMethod, exchangeTime, exchangePrice)
	if err != nil {
		resp.Error = getGRPCError(err, "")
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
	resp.Inquiry = getPBInquiry(inq, names[0], photos[0])
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
	resp.Inquiry = getPBInquiries(inqs, names, photos)
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
	resp.Inquiry = getPBInquiry(inq, names[0], photos[0])
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
	resp.Inquiry = getPBInquiry(inq, names[0], photos[0])
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
	return resp, nil
}

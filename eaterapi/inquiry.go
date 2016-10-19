package main

import (
	"golang.org/x/net/context"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
)

func (s *service) MakeInquiry(ctx context.Context, req *pb.MakeInquiryRequest) (*pb.MakeInquiryResponse, error) {
	resp := new(pb.MakeInquiryResponse)
	defer handleResp(ctx, "MakeInquiry", resp.Error)

	return resp, nil
}

func (s *service) GetInquiries(ctx context.Context, req *pb.GetInquiriesRequest) (*pb.GetInquiriesResponse, error) {
	resp := new(pb.GetInquiriesResponse)
	defer handleResp(ctx, "GetInquiries", resp.Error)

	return resp, nil
}

func (s *service) GetInquiry(ctx context.Context, req *pb.GetInquiryRequest) (*pb.GetInquiryResponse, error) {
	resp := new(pb.GetInquiryResponse)
	defer handleResp(ctx, "GetInquiry", resp.Error)

	return resp, nil

}

func (s *service) CancelInquiry(ctx context.Context, req *pb.CancelInquiryRequest) (*pb.CancelInquiryResponse, error) {
	resp := new(pb.CancelInquiryResponse)
	defer handleResp(ctx, "CancelInquiry", resp.Error)

	return resp, nil
}

func (s *service) GetBraintreeToken(ctx context.Context, req *pb.GigatokenOnlyRequest) (resp *pb.GetBraintreeTokenResponse, unusedErr error) {
	defer handleResp(ctx, "GetBraintreeToken", resp.Error)
	user, validateErr := validateGigatokenAndGetUser(ctx, req.Gigatoken)
	if validateErr != nil {
		resp.Error = validateErr
		return
	}
	eaterC := eater.New(ctx)
	customerID, err := eaterC.GetBTCustomerID(user.ID)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to eater.GetBTCustomerID")
		return
	}
	paymentC := payment.New(ctx)
	token, err := paymentC.GenerateToken(customerID)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to payment.GenerateToken")
		return
	}
	resp.BraintreeToken = token
	return resp, nil
}

func (s *service) CheckDeliveryAddresses(ctx context.Context, req *pb.CheckDeliveryAddressesRequest) (resp *pb.CheckDeliveryAddressesResponse, unusedErr error) {
	defer handleResp(ctx, "CheckDeliveryAddresses", resp.Error)

	return
}

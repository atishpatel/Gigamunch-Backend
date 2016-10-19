package main

import (
	"golang.org/x/net/context"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

func (s *service) SignIn(ctx context.Context, req *pb.SignInRequest) (resp *pb.SignInResponse, unusedErr error) {
	defer handleResp(ctx, "SignIn", resp.Error)
	user, gigatoken, err := auth.GetSessionWithGToken(ctx, req.Gtoken)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to auth.GetSessionWithGToken")
	}
	eaterC := eater.New(ctx)
	_, err = eaterC.Update(user)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to eater.Update")
		return resp, nil
	}
	resp.Gigatoken = gigatoken
	return resp, nil
}

func (s *service) SignOut(ctx context.Context, req *pb.GigatokenOnlyRequest) (resp *pb.ErrorOnlyResponse, unusedErr error) {
	defer handleResp(ctx, "SignOut", resp.Error)
	if req.Gigatoken == "" {
		return resp, nil
	}
	err := auth.DeleteSessionToken(ctx, req.Gigatoken)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to auth.DeleteSessionToken")
	}
	return resp, nil
}

func (s *service) RefreshToken(ctx context.Context, req *pb.GigatokenOnlyRequest) (resp *pb.RefreshTokenResponse, unusedErr error) {
	defer handleResp(ctx, "RefreshToken", resp.Error)
	if validateErr := validateGigatokenOnlyReq(req); validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	newToken, err := auth.RefreshToken(ctx, req.Gigatoken)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to auth.RefreshToken")
		return resp, nil
	}
	resp.Gigatoken = newToken
	return resp, nil
}

func (s *service) GetAddresses(ctx context.Context, req *pb.GigatokenOnlyRequest) (resp *pb.GetAddressesResponse, unusedErr error) {
	defer handleResp(ctx, "GetAddresses", resp.Error)
	user, validateErr := validateGigatokenAndGetUser(ctx, req.Gigatoken)
	if validateErr != nil {
		resp.Error = validateErr
		return
	}
	eaterC := eater.New(ctx)
	addresses, err := eaterC.GetAddresses(user.ID)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to eater.GetAddresses")
		return
	}
	resp.Addresses = make([]*pb.Address, len(addresses))
	for i := range addresses {
		resp.Addresses[i] = &pb.Address{
			Country:    addresses[i].Country,
			State:      addresses[i].State,
			City:       addresses[i].City,
			Zip:        addresses[i].Zip,
			Street:     addresses[i].Street,
			UnitNumber: addresses[i].APT,
			Latitude:   addresses[i].Latitude,
			Longitude:  addresses[i].Longitude,
			IsSelected: addresses[i].Selected,
		}
	}
	return
}

func (s *service) SelectAddress(ctx context.Context, req *pb.SelectAddressRequest) (resp *pb.SelectAddressResponse, unusedErr error) {
	defer handleResp(ctx, "SelectAddress", resp.Error)
	user, validateErr := validateSelectAddressRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return
	}
	eaterC := eater.New(ctx)
	address, err := eaterC.SelectAddress(user.ID, &types.Address{
		Country: req.Address.Country,
		State:   req.Address.State,
		City:    req.Address.City,
		Zip:     req.Address.Zip,
		Street:  req.Address.Street,
		APT:     req.Address.UnitNumber,
		GeoPoint: types.GeoPoint{
			Latitude:  req.Address.Latitude,
			Longitude: req.Address.Longitude,
		},
	})
	if err != nil {
		resp.Error = getGRPCError(err, "failed to eater.SelectAddress")
		return
	}
	resp.Address = &pb.Address{
		Country:    address.Country,
		State:      address.State,
		City:       address.City,
		Zip:        address.Zip,
		Street:     address.Street,
		UnitNumber: address.APT,
		Latitude:   address.Latitude,
		Longitude:  address.Longitude,
		IsSelected: address.Selected,
	}
	return
}

func (s *service) RegisterNotificationToken(ctx context.Context, req *pb.RegisterNotificationTokenRequest) (resp *pb.ErrorOnlyResponse, unusedErr error) {
	defer handleResp(ctx, "RegisterNotificationToken", resp.Error)

	return
}

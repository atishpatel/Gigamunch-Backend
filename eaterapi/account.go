package main

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

func (s *service) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.SignInResponse)
	defer handleResp(ctx, "SignIn", resp.Error)
	user, gigatoken, err := auth.GetSessionWithGToken(ctx, req.Gtoken)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to auth.GetSessionWithGToken")
		return resp, nil
	}
	if user.Name == "" && req.Name != "" {
		user.Name = req.Name
		err = auth.SaveUser(ctx, user)
		if err != nil {
			resp.Error = getGRPCError(err, "failed to auth.SaveUser")
			return resp, nil
		}
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

func (s *service) SignOut(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.ErrorOnlyResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.ErrorOnlyResponse)
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

func (s *service) RefreshToken(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.RefreshTokenResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.RefreshTokenResponse)
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

func (s *service) GetAddresses(ctx context.Context, req *pb.GigatokenOnlyRequest) (*pb.GetAddressesResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetAddressesResponse)
	defer handleResp(ctx, "GetAddresses", resp.Error)
	user, validateErr := validateGigatokenAndGetUser(ctx, req.Gigatoken)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	eaterC := eater.New(ctx)
	addresses, err := eaterC.GetAddresses(user.ID)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to eater.GetAddresses")
		return resp, nil
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
	return resp, nil
}

func (s *service) SelectAddress(ctx context.Context, req *pb.SelectAddressRequest) (*pb.SelectAddressResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.SelectAddressResponse)
	defer handleResp(ctx, "SelectAddress", resp.Error)
	user, validateErr := validateSelectAddressRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
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
		return resp, nil
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
	return resp, nil
}

func (s *service) RegisterNotificationToken(ctx context.Context, req *pb.RegisterNotificationTokenRequest) (*pb.ErrorOnlyResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.ErrorOnlyResponse)
	defer handleResp(ctx, "RegisterNotificationToken", resp.Error)
	user, validateErr := validateRegisterNotificationTokenRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	eaterC := eater.New(ctx)
	err := eaterC.RegisterNotificationToken(user.ID, req.NotificationToken)
	if err != nil {
		getGRPCError(err, "failed to eater.RegisterNotificationToken")
	}
	return resp, nil
}

func (s *service) GetMessageToken(ctx context.Context, req *pb.GetMessageTokenRequest) (*pb.GetMessageTokenResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.GetMessageTokenResponse)
	defer handleResp(ctx, "GetMessageToken", resp.Error)
	user, validateErr := validateGetMessageTokenRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	messageC := message.New(ctx)
	userInfo := &message.UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Image: user.PhotoURL,
	}
	tkn, err := messageC.GetToken(userInfo, req.DeviceId)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to message.GetToken")
		return resp, nil
	}
	resp.Token = tkn
	return resp, nil
}

func (s *service) CreateMessageChannel(ctx context.Context, req *pb.CreateMessageChannelRequest) (*pb.ErrorOnlyResponse, error) {
	ctx = appengine.BackgroundContext()
	resp := new(pb.ErrorOnlyResponse)
	defer handleResp(ctx, "CreateMessageChannel", resp.Error)
	user, validateErr := validateCreateMessageChannelRequest(ctx, req)
	if validateErr != nil {
		resp.Error = validateErr
		return resp, nil
	}
	cookC := cook.New(ctx)
	ck, err := cookC.Get(req.CookId)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to cook.Get")
		return resp, nil
	}
	cookInfo := &message.UserInfo{
		ID:    ck.ID,
		Name:  ck.Name,
		Image: ck.PhotoURL,
	}
	eaterInfo := &message.UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Image: user.PhotoURL,
	}
	messageC := message.New(ctx)
	err = messageC.UpdateChannel(cookInfo, eaterInfo, nil)
	if err != nil {
		resp.Error = getGRPCError(err, "failed to message.UpdateChannel")
		return resp, nil
	}
	return resp, nil
}

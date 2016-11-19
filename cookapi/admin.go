package main

import (
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

// CreateFakeGigatokenReq is the request for CreateFakeGigatoken.
type CreateFakeGigatokenReq struct {
	GigatokenReq
	UserID string `json:"user_id"`
}

// CreateFakeGigatoken creates a fake Gigatoken if user is an admin.
func (service *Service) CreateFakeGigatoken(ctx context.Context, req *CreateFakeGigatokenReq) (*RefreshTokenResp, error) {
	resp := new(RefreshTokenResp)
	defer handleResp(ctx, "CreateFakeGigatoken", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	gigatoken, err := auth.GetGigatokenViaAdmin(ctx, user, req.UserID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to auth.GetGigatokenViaAdmin")
		return resp, nil
	}
	resp.Gigatoken = gigatoken
	return resp, nil
}

// AddToProcessInquiryQueueReq is the request for AddToProcessInquiryQueue.
type AddToProcessInquiryQueueReq struct {
	IDReq
	Hours int `json:"hours"`
}

// AddToProcessInquiryQueue adds a process to the inquiry queue.
func (service *Service) AddToProcessInquiryQueue(ctx context.Context, req *AddToProcessInquiryQueueReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "AddToProcessInquiryQueue", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	at := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	tasksC := tasks.New(ctx)
	err = tasksC.AddProcessInquiry(req.ID, at)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to tasks.AddProcessInquiry")
		return resp, nil
	}

	return resp, nil
}

// CreateFakeSubmerchantReq is the request for CreateFakeSubmerchant.
type CreateFakeSubmerchantReq struct {
	GigatokenReq
	CookID string `json:"cook_id"`
}

// CreateFakeSubmerchant creates a fake submerchant for cook.
func (service *Service) CreateFakeSubmerchant(ctx context.Context, req *CreateFakeSubmerchantReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "CreateFakeSubmerchant", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	cookC := cook.New(ctx)
	ck, err := cookC.Get(req.CookID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to cook.Get")
		return resp, nil
	}
	if ck.BTSubMerchantID == "" {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "BTSubMerchantID is empty."}
		return resp, nil
	}
	ckUser := &types.User{
		ID:       ck.ID,
		Name:     ck.Name,
		Email:    ck.Email,
		PhotoURL: ck.PhotoURL,
	}
	paymentC := payment.New(ctx)
	err = paymentC.CreateFakeSubMerchant(ckUser, ck.BTSubMerchantID)
	if err != nil {
		resp.Err = errors.Wrap("failed to paymentC.CreateFakeSubMerchant", err)
		return resp, nil
	}
	return resp, nil
}

// SendSMSReq is the request for CreateFakeSubmerchant.
type SendSMSReq struct {
	GigatokenReq
	Number  string `json:"number"`
	Message string `json:"message"`
}

// SendSMS sends an sms from Gigamunch to number.
func (service *Service) SendSMS(ctx context.Context, req *SendSMSReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SendSMS", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	messageC := message.New(ctx)
	err = messageC.SendSMS(req.Number, req.Message)
	if err != nil {
		resp.Err = errors.Wrap("failed to message.SendSMS", err)
		return resp, nil
	}
	return resp, nil
}

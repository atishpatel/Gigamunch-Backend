package main

import (
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/promo"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
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

// CreatePromoCodeReq is the request for CreatePromoCode.
type CreatePromoCodeReq struct {
	GigatokenReq
	promo.Code
}

// CreatePromoCode creates a Promo Code.
func (service *Service) CreatePromoCode(ctx context.Context, req *CreatePromoCodeReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "CreatePromoCode", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	promoC := promo.New(ctx)
	err = promoC.InsertCode(user, &req.Code)
	if err != nil {
		resp.Err = errors.Wrap("failed to promo.InsertCode", err)
		return resp, nil
	}
	return resp, nil
}

// SubLogReq is the request for SubLogStuff.
type SubLogReq struct {
	GigatokenReq
	SubEmail string    `json:"sub_email"`
	Date     time.Time `json:"date"`
}

// SetupSubLogs runs sub.SetupSubLogs.
func (service *Service) SetupSubLogs(ctx context.Context, req *SubLogReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SetupSubLogs", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subC := sub.New(ctx)
	err = subC.SetupSubLogs(req.Date)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.SetupSubLogs")
		return resp, nil
	}
	return resp, nil
}

// ProcessSubLog runs sub.Process.
func (service *Service) ProcessSubLog(ctx context.Context, req *SubLogReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "ProcessSubLog", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subC := sub.New(ctx)
	err = subC.Process(req.Date, req.SubEmail)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.Process")
		return resp, nil
	}
	return resp, nil
}

// GetSubEmailsResp is a resp for GetSubEmails.
type GetSubEmailsResp struct {
	SubEmails []string `json:"sub_emails"`
	ErrorOnlyResp
}

// GetSubEmails returns a list of SubEmails that can be skipped from the last month.
func (service *Service) GetSubEmails(ctx context.Context, req *GigatokenReq) (*GetSubEmailsResp, error) {
	resp := new(GetSubEmailsResp)
	defer handleResp(ctx, "GetSubEmails", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	from := time.Now().Add(-14 * 24 * time.Hour)
	to := time.Now().Add(14 * 24 * time.Hour)
	subC := sub.New(ctx)
	resp.SubEmails, err = subC.GetSubEmails(from, to)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.GetSubEmails")
		return resp, nil
	}
	return resp, nil
}

// SkipSubLog runs sub.Skip.
func (service *Service) SkipSubLog(ctx context.Context, req *SubLogReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SkipSubLog", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subC := sub.New(ctx)
	err = subC.Skip(req.Date, req.SubEmail)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.Skip")
		return resp, nil
	}
	return resp, nil
}

// FreeSubLog runs sub.Free.
func (service *Service) FreeSubLog(ctx context.Context, req *SubLogReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "FreeSubLog", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subC := sub.New(ctx)
	err = subC.Free(req.Date, req.SubEmail)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.Free")
		return resp, nil
	}
	return resp, nil
}

// SubReq is the request for Sub stuff.
type SubReq struct {
	GigatokenReq
	SubEmail string `json:"sub_email"`
}

// CancelSub cancels a subscriber's subscription.
func (service *Service) CancelSub(ctx context.Context, req *SubReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "CancelSub", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subC := sub.New(ctx)
	err = subC.Cancel(req.SubEmail)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.Cancel")
		return resp, nil
	}
	return resp, nil
}

// GetSubLogsResp is a resp for GetSubLogs.
type GetSubLogsResp struct {
	SubLogs []*sub.SubscriptionLog `json:"sublogs"`
	ErrorOnlyResp
}

// GetSubLogs gets all the SubLogs.
func (service *Service) GetSubLogs(ctx context.Context, req *GigatokenReq) (*GetSubLogsResp, error) {
	resp := new(GetSubLogsResp)
	defer handleResp(ctx, "GetSubLogs", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	subC := sub.New(ctx)
	subLogs, err := subC.GetAll(1000)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.GetAll")
		return resp, nil
	}
	resp.SubLogs = subLogs
	return resp, nil
}

// SubLog is the sublog with subscriber's name and address.
type SubLog struct {
	Date         time.Time `json:"date"`
	Servings     int8      `json:"servings"`
	DeliveryTime int8      `json:"delivery_time"`
	sub.SubscriptionLog
	sub.SubscriptionSignUp
}

// GetSubLogsForDateResp is a resp for GetSubLogsForDate.
type GetSubLogsForDateResp struct {
	SubLogs []SubLog `json:"sublogs"`
	ErrorOnlyResp
}

// GetSubLogsForDate gets all the SubLogs for a date.
func (service *Service) GetSubLogsForDate(ctx context.Context, req *DateReq) (*GetSubLogsForDateResp, error) {
	resp := new(GetSubLogsForDateResp)
	defer handleResp(ctx, "GetSubLogsForDate", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	subC := sub.New(ctx)
	subLogs, err := subC.GetForDate(req.Date)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.GetForDate")
		return resp, nil
	}
	if len(subLogs) != 0 {
		subEmails := make([]string, len(subLogs))
		for i := range subLogs {
			subEmails[i] = subLogs[i].SubEmail
		}
		subs, err := subC.GetSubscribers(subEmails)
		if err != nil {
			resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.GetSubscribers")
			return resp, nil
		}
		resp.SubLogs = make([]SubLog, len(subLogs))
		for i := range subLogs {
			for j := range subs {
				if subLogs[i].SubEmail == subs[j].Email {
					resp.SubLogs[i].SubscriptionLog = *subLogs[i]
					resp.SubLogs[i].SubscriptionSignUp = *subs[j]
					resp.SubLogs[i].Date = subLogs[i].Date
					resp.SubLogs[i].Servings = subs[i].Servings
					resp.SubLogs[i].DeliveryTime = subs[i].DeliveryTime
				}
			}
		}
	}
	return resp, nil
}

// AddToProcessSubscriptionQueueReq is a request for AddToProcessSubscriptionQueue.
type AddToProcessSubscriptionQueueReq struct {
	SubLogReq
	Hours int `json:"hours"`
}

// AddToProcessSubscriptionQueue adds a process to the subscription queue.
func (service *Service) AddToProcessSubscriptionQueue(ctx context.Context, req *AddToProcessSubscriptionQueueReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "AddToProcessSubscriptionQueue", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subR := &tasks.ProcessSubscriptionParams{
		SubEmail: req.SubEmail,
		Date:     req.Date,
	}
	at := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	tasksC := tasks.New(ctx)
	err = tasksC.AddProcessSubscription(at, subR)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to tasks.AddProcessSubscription")
		return resp, nil
	}

	return resp, nil
}

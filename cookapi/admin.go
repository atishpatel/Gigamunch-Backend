package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	authnew "github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	subnew "github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/promo"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func getUserFromRequest(ctx context.Context, req validatableTokenReq) (*common.User, error) {
	log, serverInfo, db, sqlDB, err := setupAll(ctx, "/cookapi")
	if err != nil {
		return nil, errors.Annotate(err, "failed to setupLoggingAndServerInfo")
	}
	authC, err := authnew.NewClient(ctx, log, db, sqlDB, serverInfo)
	if err != nil {
		return nil, errors.Annotate(err, "failed to get auth.NewClient")
	}
	user, err := authC.Verify(req.gigatoken())
	if err != nil {

		return nil, errors.Annotate(err, "failed to get auth.Verify")
	}
	return user, err
}

// AddToProcessInquiryQueueReq is the request for AddToProcessInquiryQueue.
type AddToProcessInquiryQueueReq struct {
	EmailReq
	Hours int `json:"hours"`
}

// AddToProcessInquiryQueue adds a process to the inquiry queue.
func (service *Service) AddToProcessInquiryQueue(ctx context.Context, req *AddToProcessInquiryQueueReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "AddToProcessInquiryQueue", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	at := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	tasksC := tasks.New(ctx)
	err = tasksC.AddUpdateDrip(at, &tasks.UpdateDripParams{Email: req.Email})
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to tasks.AddProcessInquiry")
		return resp, nil
	}

	return resp, nil
}

// CreateFakeSubmerchantReq is the request for CreateFakeSubmerchant.
// type CreateFakeSubmerchantReq struct {
// 	GigatokenReq
// 	CookID string `json:"cook_id"`
// }

// // CreateFakeSubmerchant creates a fake submerchant for cook.
// func (service *Service) CreateFakeSubmerchant(ctx context.Context, req *CreateFakeSubmerchantReq) (*ErrorOnlyResp, error) {
// 	resp := new(ErrorOnlyResp)
// 	defer handleResp(ctx, "CreateFakeSubmerchant", resp.Err)
// 	user, err := getUserFromRequest(ctx, req)
// 	if err != nil {
// 		resp.Err = errors.GetErrorWithCode(err)
// 		return resp, nil
// 	}
// 	if !user.Admin {
// 		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
// 		return resp, nil
// 	}
// 	cookC := cook.New(ctx)
// 	ck, err := cookC.Get(req.CookID)
// 	if err != nil {
// 		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to cook.Get")
// 		return resp, nil
// 	}
// 	if ck.BTSubMerchantID == "" {
// 		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "BTSubMerchantID is empty."}
// 		return resp, nil
// 	}
// 	ckUser := &types.User{
// 		ID:       ck.ID,
// 		Name:     ck.Name,
// 		Email:    ck.Email,
// 		PhotoURL: ck.PhotoURL,
// 	}
// 	paymentC := payment.New(ctx)
// 	err = paymentC.CreateFakeSubMerchant(ckUser, ck.BTSubMerchantID)
// 	if err != nil {
// 		resp.Err = errors.Wrap("failed to paymentC.CreateFakeSubMerchant", err)
// 		return resp, nil
// 	}
// 	return resp, nil
// }

// SendSMSReq is the request for CreateFakeSubmerchant.
type SendSMSReq struct {
	GigatokenReq
	Numbers string `json:"numbers"`
	Message string `json:"message"`
}

// SendSMS sends an sms from Gigamunch to number.
func (service *Service) SendSMS(ctx context.Context, req *SendSMSReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SendSMS", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	messageC := message.New(ctx)
	var numbers []string
	if strings.Contains(req.Numbers, ",") {
		numbers = strings.Split(req.Numbers, ",")
	} else {
		numbers = []string{req.Numbers}
	}
	for _, n := range numbers {
		err = messageC.SendDeliverySMS(n, req.Message)
		if err != nil {
			resp.Err = errors.Wrap("failed to message.SendSMS", err)
			return resp, nil
		}
	}
	return resp, nil
}

// SendCustomerSMSReq is the request for CreateFakeSubmerchant.
type SendCustomerSMSReq struct {
	GigatokenReq
	Emails  string `json:"emails"`
	Message string `json:"message"`
}

// SendCustomerSMS sends an CustomerSMS from Gigamunch to number.
func (service *Service) SendCustomerSMS(ctx context.Context, req *SendCustomerSMSReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SendCustomerSMS", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	log, _, _, _, _ := setupAll(ctx, "/cookapi/api/SendCusomterSMS")
	dilm := "{{name}}"
	// if !strings.Contains(req.Message, dilm) {
	// 	resp.Err = errors.BadRequestError.WithMessage("Message requires {{name}}.")
	// 	return resp, nil
	// }
	messageC := message.New(ctx)
	var emails []string
	if strings.Contains(req.Emails, ",") {
		emails = strings.Split(req.Emails, ",")
	} else {
		emails = []string{req.Emails}
	}
	subC := subold.New(ctx)
	subs, err := subC.GetSubscribers(emails)
	if err != nil {
		resp.Err = errors.Wrap("failed to subold.GetSubscribers", err)
		return resp, nil
	}
	for _, s := range subs {
		if s.PhoneNumber == "" {
			continue
		}
		name := s.FirstName
		if name == "" {
			name = s.Name
		}
		name = strings.Title(name)
		msg := strings.Replace(req.Message, dilm, name, -1)
		err = messageC.SendDeliverySMS(s.PhoneNumber, msg)
		if err != nil {
			resp.Err = errors.Wrap("failed to message.SendSMS", err)
			return resp, nil
		}
		// log
		if log != nil {
			payload := &logging.MessagePayload{
				Platform: "SMS",
				Body:     msg,
				From:     "Gigamunch",
				To:       s.PhoneNumber,
			}
			log.SubMessage(s.ID, s.Email, payload)
		}
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

// SetupSubLogs runs subold.SetupSubLogs.
func (service *Service) SetupSubLogs(ctx context.Context, req *DateReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SetupSubLogs", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	if req.Date.Before(time.Now()) {
		resp.Err = errors.BadRequestError.WithMessage("Date is before now.")
		return resp, nil
	}
	log, serverInfo, db, sqlDB, err := setupAll(ctx, "/cookapi/setupsublogs")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	subC, err := subnew.NewClient(ctx, log, db, sqlDB, serverInfo)
	err = subC.SetupActivities(req.Date)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subnew.SetupActivities")
		return resp, nil
	}
	return resp, nil
}

// ProcessSubLog runs subnew.Process.
func (service *Service) ProcessSubLog(ctx context.Context, req *SubLogReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "ProcessSubLog", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	log, serverInfo, db, sqlDB, err := setupAll(ctx, "/cookapi/processsublogs")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	subC, err := subnew.NewClient(ctx, log, db, sqlDB, serverInfo)
	err = subC.ProcessActivity(req.Date, req.SubEmail)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to activity.Process")
		return resp, nil
	}
	return resp, nil
}

// GetSubEmailsResp is a resp for GetSubEmails.
type GetSubEmailsResp struct {
	SubEmails   []string                     `json:"sub_emails"`
	Subscribers []*subold.SubscriptionSignUp `json:"subscribers"`
	ErrorOnlyResp
}

// GetSubEmails returns a list of SubEmails that can be skipped from the last month.
func (service *Service) GetSubEmails(ctx context.Context, req *GigatokenReq) (*GetSubEmailsResp, error) {
	resp := new(GetSubEmailsResp)
	defer handleResp(ctx, "GetSubEmails", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	from := time.Now().Add(-7 * 24 * time.Hour)
	to := time.Now().Add(14 * 24 * time.Hour)
	subC := subold.New(ctx)
	resp.SubEmails, err = subC.GetSubEmails(from, to)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetSubEmails")
		return resp, nil
	}
	resp.Subscribers, err = subC.GetSubscribers(resp.SubEmails)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetSubscribers")
		return resp, nil
	}
	return resp, nil
}

// SkipSubLog runs subnew.Skip.
func (service *Service) SkipSubLog(ctx context.Context, req *SubLogReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "SkipSubLog", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	log, serverInfo, db, sqlDB, err := setupAll(ctx, "/cookapi/skipsublogs")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	activityC, err := activity.NewClient(ctx, log, db, sqlDB, serverInfo)
	err = activityC.Skip(req.Date, req.SubEmail, "Admin skip.")

	// subC := subold.New(ctx)
	// err = subC.Skip(req.Date, req.SubEmail, "Admin skip.")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to activity.Skip")
		return resp, nil
	}
	return resp, nil
}

type DiscountSubLogReq struct {
	SubLogReq
	Amount           float32 `json:"amount"`
	Percent          int8    `json:"percent"`
	OverrideDiscount bool    `json:"override_discount"`
}

// DiscountSubLog gives a discount to a user for a specific week.
func (service *Service) DiscountSubLog(ctx context.Context, req *DiscountSubLogReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "DiscountSubLog", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	log, serverInfo, db, sqlDB, err := setupAll(ctx, "/cookapi/setupsublogs")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}

	// discountC, err := discount.NewClient(ctx, log, db, sqlDB, serverInfo)
	// err = discountC.Discount(req.Date, req.SubEmail, req.Amount, req.Percent, req.OverrideDiscount)

	activityC, err := activity.NewClient(ctx, log, db, sqlDB, serverInfo)
	err = activityC.Discount(req.Date, req.SubEmail, req.Amount, req.Percent, req.OverrideDiscount)

	// subC := subold.New(ctx)
	// err = subC.Discount(req.Date, req.SubEmail, req.Amount, req.Percent, req.OverrideDiscount)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to activity.Discount")
		return resp, nil
	}
	return resp, nil
}

// ChangeServingsPermanentlyReq is a request for ChangeServingsPermanently.
type ChangeServingsPermanentlyReq struct {
	EmailReq
	Servings   int8 `json:"servings"`
	Vegetarian bool `json:"vegetarian"`
}

// ChangeServingsPermanently gives a changes the serving count for a user permanently.
func (service *Service) ChangeServingsPermanently(ctx context.Context, req *ChangeServingsPermanentlyReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "ChangeServingsPermanently", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	log, serverInfo, _, _, err := setupAll(ctx, "/cookapi/ChangeServingsPermanently")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	// subC, err := subnew.NewClient(ctx, log, db, sqlDB, serverInfo)
	// err = subC.ChangeServingsPermanently(req.Email, req.Servings, req.Vegetarian)
	subC := subold.NewWithLogging(ctx, log)
	err = subC.ChangeServingsPermanently(req.Email, req.Servings, req.Vegetarian, serverInfo)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subnew.ChangeServingsPermanently")
		return resp, nil
	}
	return resp, nil
}

// UpdatePaymentMethodTokenReq is a request for UpdatePaymentMethodToken.
type UpdatePaymentMethodTokenReq struct {
	EmailReq
	PaymentMethodToken string `json:"payment_method_token"`
}

// UpdatePaymentMethodToken updates a uesr payment token.
func (service *Service) UpdatePaymentMethodToken(ctx context.Context, req *UpdatePaymentMethodTokenReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "UpdatePaymentMethodToken", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	log, serverInfo, db, sqlDB, err := setupAll(ctx, "/cookapi/setupsublogs")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	subC, err := subnew.NewClient(ctx, log, db, sqlDB, serverInfo)
	err = subC.UpdatePaymentToken(req.Email, req.PaymentMethodToken)
	// subC := subold.New(ctx)
	// err = subC.UpdatePaymentToken(req.Email, req.PaymentMethodToken)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subnew.UpdatePaymentToken")
		return resp, nil
	}
	return resp, nil
}

// ChangeServingsForDateReq is a request for ChangeServingForDate.
type ChangeServingsForDateReq struct {
	SubLogReq
	Servings int8 `json:"servings"`
}

// ChangeServingsForDate gives a changes the serving count for a user for a specific week.
func (service *Service) ChangeServingsForDate(ctx context.Context, req *ChangeServingsForDateReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "ChangeServingsForDateSubLog", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	log, _, _, _, err := setupAll(ctx, "/cookapi/setupsublogs")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	// activityC, err := activity.NewClient(ctx, log, db, sqlDB, serverInfo)
	// err = activityC.ChangeServings(req.Date, req.SubEmail, req.Servings, subold.DerivePrice(req.Servings))

	subC := subold.NewWithLogging(ctx, log)
	err = subC.ChangeServings(req.Date, req.SubEmail, req.Servings, subold.DerivePrice(req.Servings))
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to activity.ChangeServings")
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
func (service *Service) CancelSub(ctx context.Context, req *EmailReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "CancelSub", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	log, serverInfo, db, sqlDB, err := setupAll(ctx, "cookapi/CancelSub")
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	subC, err := subnew.NewClient(ctx, log, db, sqlDB, serverInfo)
	err = subC.Deactivate(req.Email, "")
	// subC := subold.New(ctx)
	// err = subC.Cancel(req.Email, log, serverInfo)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subnew.Deactivate")
		return resp, nil
	}
	return resp, nil
}

// GetSubLogsResp is a resp for GetSubLogs.
type GetSubLogsResp struct {
	SubLogs []*subold.SubscriptionLog `json:"sublogs"`
	ErrorOnlyResp
}

// GetSubLogs gets all the SubLogs.
func (service *Service) GetSubLogs(ctx context.Context, req *GigatokenReq) (*GetSubLogsResp, error) {
	resp := new(GetSubLogsResp)
	defer handleResp(ctx, "GetSubLogs", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	subC := subold.New(ctx)
	subLogs, err := subC.GetAll(3000)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetAll")
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
	CustomerID   string    `json:"customer_id"`
	subold.SubscriptionLog
	subold.SubscriptionSignUp
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
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}
	subC := subold.New(ctx)
	subLogs, err := subC.GetForDate(req.Date)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetForDate")
		return resp, nil
	}
	if len(subLogs) != 0 {
		subEmails := make([]string, len(subLogs))
		for i := range subLogs {
			subEmails[i] = subLogs[i].SubEmail
		}
		subs, err := subC.GetSubscribers(subEmails)
		if err != nil {
			resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetSubscribers")
			return resp, nil
		}
		resp.SubLogs = make([]SubLog, len(subLogs))
		for i := range subLogs {
			resp.SubLogs[i].SubscriptionLog = *subLogs[i]
			resp.SubLogs[i].Date = subLogs[i].Date
			resp.SubLogs[i].DeliveryTime = subLogs[i].DeliveryTime
			resp.SubLogs[i].Servings = subLogs[i].Servings
			logging.Infof(ctx, "email:%s", subLogs[i].SubEmail)
			for j := range subs {
				if subLogs[i].SubEmail == subs[j].Email {
					resp.SubLogs[i].SubscriptionSignUp = *subs[j]
					resp.SubLogs[i].CustomerID = subs[j].CustomerID
				}
			}
		}
	}
	return resp, nil
}

// AddToProcessSubscriptionQueueReq is a request for AddToProcessSubscriptionQueue.
type AddToProcessSubscriptionQueueReq struct {
	SubLogReq
	Emails []string `json:"emails"`
	Hours  int      `json:"hours"`
}

// AddToProcessSubscriptionQueue adds a process to the subscription queue.
func (service *Service) AddToProcessSubscriptionQueue(ctx context.Context, req *AddToProcessSubscriptionQueueReq) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, "AddToProcessSubscriptionQueue", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subR := &tasks.ProcessSubscriptionParams{
		SubEmail: req.SubEmail,
		Date:     req.Date,
	}
	at := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	tasksC := tasks.New(ctx)
	if req.SubEmail != "" {
		err = tasksC.AddProcessSubscription(at, subR)
		if err != nil {
			resp.Err = errors.GetErrorWithCode(err).Wrap("failed to tasks.AddProcessSubscription")
			return resp, nil
		}
	}

	for _, email := range req.Emails {
		subR := &tasks.ProcessSubscriptionParams{
			SubEmail: email,
			Date:     req.Date,
		}
		at := time.Now().Add(time.Duration(req.Hours) * time.Hour)
		err = tasksC.AddProcessSubscription(at, subR)
		if err != nil {
			resp.Err = errors.GetErrorWithCode(err).Wrap("failed to tasks.AddProcessSubscription")
			return resp, nil
		}
	}

	return resp, nil
}

// SendEmailReq is a request for Email.
type SendEmailReq struct {
	GigatokenReq
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	FirstDinnerDate time.Time `json:"first_dinner_date"`
}

// GetName returns the name.
func (e *SendEmailReq) GetName() string {
	return e.Name
}

// GetEmail returns the email.
func (e *SendEmailReq) GetEmail() string {
	return e.Email
}

// GetFirstDinnerDate returns the first dinner for the subscriber.
func (e *SendEmailReq) GetFirstDinnerDate() time.Time {
	return e.FirstDinnerDate
}

func setupAll(ctx context.Context, path string) (*logging.Client, *common.ServerInfo, common.DB, *sqlx.DB, error) {
	dbC, err := db.NewClient(ctx, projectID, nil)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get database client: %+v", err)
	}
	// Setup logging
	serverInfo := &common.ServerInfo{
		ProjectID:           projectID,
		IsStandardAppEngine: true,
	}
	log, err := logging.NewClient(ctx, "admin", path, dbC, serverInfo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	sqlConnectionString := os.Getenv("MYSQL_CONNECTION")
	if sqlConnectionString == "" {
		return nil, nil, nil, nil, fmt.Errorf(`You need to set the environment variable "MYSQL_CONNECTION"`)
	}
	sqlDB, err := sqlx.Connect("mysql", sqlConnectionString+"?collation=utf8mb4_general_ci&parseTime=true")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get sql database client: %+v", err)
	}
	return log, serverInfo, dbC, sqlDB, nil
}

func setupLog(ctx context.Context, path string) (*logging.Client, *common.ServerInfo, error) {
	dbC, err := db.NewClient(ctx, projectID, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get database client: %+v", err)
	}
	// Setup logging
	serverInfo := &common.ServerInfo{
		ProjectID:           projectID,
		IsStandardAppEngine: true,
	}
	log, err := logging.NewClient(ctx, "admin", path, dbC, serverInfo)
	if err != nil {
		return nil, nil, err
	}
	return log, serverInfo, nil
}

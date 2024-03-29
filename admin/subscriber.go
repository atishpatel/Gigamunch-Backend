package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/maps"
	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	pbcommon "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

// SendCustomerSMS sends an CustomerSMS from Gigamunch to number.
func (s *server) SendCustomerSMS(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.SendCustomerSMSReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	tasksC := tasks.New(ctx)
	var errs []error
	for _, v := range req.Emails {
		params := &tasks.SendSMSParam{
			Message: req.Message,
		}
		if strings.Contains(v, "@") {
			params.Email = v
		} else if strings.Contains(v, "-") {
			params.Number = v
		} else {
			params.Email = v
		}
		err = tasksC.AddSendSMS(params, time.Now())
		if err != nil {
			errs = append(errs, errors.Annotate(err, "failed to message.SendSMS To("+v+")"))
			continue
		}
	}
	if len(errs) >= 1 {
		return errors.GetErrorWithCode(errs[0]).Annotatef("errors count: %d", len(errs))
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// GetSubscriber gets all info about a subscriber from their email address
func (s *server) GetSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetSubscriberReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	email := req.Email

	subC := subold.New(ctx)
	subscriber, err := subC.GetSubscriber(email)

	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get subscriber")
	}

	resp := &pb.GetSubscriberResp{
		Subscriber: pbSubscriber(subscriber),
	}

	return resp
}

// GetSubscriberV2 gets all info about a subscriber.
func (s *server) GetSubscriberV2(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.UserIDReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	subscriber, err := subC.Get(req.ID)
	if err != nil {
		return errors.Annotate(err, "failed to get subscriber")
	}

	sResp, err := serverhelper.PBSubscriber(subscriber)
	if err != nil {
		return errors.Annotate(err, "failed to encode")
	}

	resp := &pb.GetSubscriberRespV2{
		Subscriber: sResp,
	}

	return resp
}

// GetHasSubscribed gets all subscribers.
func (s *server) GetHasSubscribed(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetHasSubscribedReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	var t time.Time

	subC := subold.New(ctx)
	subscribers, err := subC.GetHasSubscribed(t)

	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get all subscribers")
	}

	resp := &pb.GetHasSubscribedResp{
		Subscribers: pbSubscribers(subscribers),
	}

	return resp
}

// GetHasSubscribedV2 gets all subscribers.
func (s *server) GetHasSubscribedV2(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetHasSubscribedReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed sub.NewClient")
	}
	subscribers, err := subC.GetHasSubscribed(int(req.Start), int(req.Limit))
	if err != nil {
		return errors.Annotate(err, "failed to get all subscribers")
	}
	ss, err := serverhelper.PBSubscribers(subscribers)
	if err != nil {
		return errors.Annotate(err, "failed to PBSubscribers")
	}
	resp := &pb.GetHasSubscribedRespV2{
		Subscribers: ss,
	}

	return resp
}

// ChangeSubscriberServings change a subscriber's servings.
func (s *server) ChangeSubscriberServings(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.ChangeSubscriberServingsReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	err = subC.ChangeServingsPermanently(req.ID, int8(req.ServingsNonVeg), int8(req.ServingsVeg))
	if err != nil {
		return errors.Annotate(err, "failed to sub.ChangeServingsPermanently")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// UpdateAddress activates a subscriber account.
func (s *server) UpdateAddress(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.UpdateAddressReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	subscriber, err := subC.Get(req.UserID)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Get")
	}
	address, err := maps.GetAddress(ctx, req.Address.FullAddress, req.Address.Apt)
	if err != nil {
		return errors.Annotate(err, "failed to sub.GetAddress")
	}
	subscriber.DeliveryNotes = req.DeliveryNotes
	subscriber.Address = *address
	err = subC.Update(subscriber)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Activate")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// ActivateSubscriber activates a subscriber account.
func (s *server) ActivateSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.ActivateSubscriberReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	firstBagDate := getDatetime(req.FirstBagDate)

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	err = subC.Activate(req.Email, firstBagDate)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Activate")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// DeactivateSubscriber activates a subscriber account.
func (s *server) DeactivateSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.DeactivateSubscriberReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	err = subC.Deactivate(req.ID, req.Reason)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Deactivate")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// UpdateDrip updates drip.
func (s *server) UpdateDrip(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	// TODO: Rename to AddToUpdateDripQueue
	var err error
	req := new(pb.UpdateDripReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	at := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	tasksC := tasks.New(ctx)
	for _, email := range req.Emails {
		err = tasksC.AddUpdateDrip(at, &tasks.UpdateDripParams{Email: email})
		if err != nil {
			return errors.Annotate(err, "failed to tasks.AddUpdateDrip")
		}
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// ChangeSubscriberPlanDay change a subscriber's PlanDay.
func (s *server) ChangeSubscriberPlanDay(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.ChangeSubscriberPlanDayReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	activitySwitchDate := getDatetime(req.ActivitySwitchDate)
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	err = subC.ChangePlanDay(req.ID, req.NewPlanDay, &activitySwitchDate)
	if err != nil {
		return errors.Annotate(err, "failed to sub.ChangePlanDay")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// ReplaceSubscriberEmail replaces a subscriber's old email with a new email.
func (s *server) ReplaceSubscriberEmail(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.ReplaceSubscriberEmailReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.OldEmail == "" {
		return errors.Annotate(errors.BadRequestError, "old email cannot be empty")
	}
	if req.NewEmail == "" {
		req.NewEmail = req.OldEmail
	}
	req.NewEmail = strings.ToLower(req.NewEmail)
	subnewC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	snew, err := subnewC.GetByEmail(req.OldEmail)
	if err != nil {
		return errors.Annotate(err, "failed to sub.GetByEmail")
	}
	for i := range snew.EmailPrefs {
		if snew.EmailPrefs[i].Email == req.OldEmail {
			snew.EmailPrefs[i].Email = req.NewEmail
			snew.EmailPrefs[i].FirstName = req.FirstName
			snew.EmailPrefs[i].LastName = req.LastName
		}
	}
	err = subnewC.Update(snew)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Update")
	}
	subC := subold.New(ctx)
	err = subC.UpdateEmail(req.OldEmail, req.NewEmail)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to subold.UpdateEmail")
	}
	mailC, err := mail.NewClient(ctx, log, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to mail.NewClient")
	}
	err = mailC.UpdateUser(&mail.UserFields{
		Email:    req.OldEmail,
		NewEmail: req.NewEmail,
	})
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to mail.Update")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

func pbSubscribers(subscribers []subold.SubscriptionSignUp) []*pb.Subscriber {
	sls := make([]*pb.Subscriber, len(subscribers))
	for i := range subscribers {
		sls[i] = pbSubscriber(&subscribers[i])
	}
	return sls
}

func pbSubscriber(subscriber *subold.SubscriptionSignUp) *pb.Subscriber {
	return &pb.Subscriber{
		Email:              subscriber.Email,
		Date:               subscriber.Date.Format(time.RFC3339),
		Name:               subscriber.Name,
		FirstName:          subscriber.FirstName,
		LastName:           subscriber.LastName,
		Address:            pbAddress(&subscriber.Address),
		CustomerID:         subscriber.CustomerID,
		SubscriptionIDs:    subscriber.SubscriptionIDs,
		FirstPaymentDate:   subscriber.FirstPaymentDate.Format(time.RFC3339),
		IsSubscribed:       subscriber.IsSubscribed,
		SubscriptionDate:   subscriber.SubscriptionDate.Format(time.RFC3339),
		UnsubscribedDate:   subscriber.UnSubscribedDate.Format(time.RFC3339),
		FirstBoxDate:       subscriber.FirstBoxDate.Format(time.RFC3339),
		Servings:           int32(subscriber.Servings),
		VegetarianServings: int32(subscriber.VegetarianServings),
		DeliveryTime:       int32(subscriber.DeliveryTime),
		SubscriptionDay:    subscriber.SubscriptionDay,
		WeeklyAmount:       subscriber.WeeklyAmount,
		PaymentMethodToken: subscriber.PaymentMethodToken,
		Reference:          subscriber.Reference,
		PhoneNumber:        subscriber.PhoneNumber,
		DeliveryTips:       subscriber.DeliveryTips,
		BagReminderSMS:     subscriber.BagReminderSMS,
		// gift
		NumGiftDinners: int32(subscriber.NumGiftDinners),
		ReferenceEmail: subscriber.ReferenceEmail,
		GiftRevealDate: subscriber.GiftRevealDate.Format(time.RFC3339),
		// stats
		ReferralPageOpens: int32(subscriber.ReferralPageOpens),
		ReferredPageOpens: int32(subscriber.ReferredPageOpens),
		GiftPageOpens:     int32(subscriber.GiftPageOpens),
		GiftedPageOpens:   int32(subscriber.GiftedPageOpens),
	}
}

func pbAddress(address *types.Address) *pbcommon.Address {
	return &pbcommon.Address{
		Apt:         address.APT,
		Street:      address.Street,
		City:        address.City,
		State:       address.State,
		Zip:         address.Zip,
		Country:     address.Country,
		Latitude:    address.Latitude,
		Longitude:   address.Longitude,
		FullAddress: address.String(),
	}
}

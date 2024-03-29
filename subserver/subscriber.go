package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbsub"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/maps"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"

	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// ChangeSubscriberServings change a subscriber's servings.
func (s *server) ChangeSubscriberServings(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.ChangeSubscriberServingsReq)

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
	err = subC.ChangeServingsPermanently(user.ID, int8(req.ServingsNonVeg), int8(req.ServingsVeg))
	if err != nil {
		return errors.Annotate(err, "failed to sub.ChangeServingsPermanently")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

// ActivateSubscriber activates a subscriber account.
func (s *server) ActivateSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.ActivateSubscriberReq)

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
	subscriber, err := subC.Get(user.ID)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Get")
	}

	planDay := subscriber.PlanWeekday
	intervalStartDate := time.Now().Add(81 * time.Hour)
	for intervalStartDate.Weekday().String() != planDay {
		intervalStartDate = intervalStartDate.Add(time.Hour * 24)
	}
	firstBagDate := intervalStartDate
	if req.FirstBagDate != "" {
		firstBagDate = getDatetime(req.FirstBagDate)
	}

	err = subC.Activate(subscriber.ID, firstBagDate)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Activate")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

// DeactivateSubscriber deactivates the subscriber.
func (s *server) DeactivateSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.DeactivateSubscriberReq)
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
	subscriber, err := subC.Get(user.ID)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Get")
	}
	err = subC.Deactivate(subscriber.ID, req.Reason)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Deactivate")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

// UpdateSubscriber updates the subscriber.
func (s *server) UpdateSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.UpdateSubscriberReq)
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
	subscriber, err := subC.Get(user.ID)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Get")
	}
	if req.Address.FullAddress == "" {
		a := req.Address
		req.Address.FullAddress = fmt.Sprintf("%s, %s, %s %s, %s", a.Street, a.City, a.State, a.Zip, a.Country)
	}
	address, err := maps.GetAddress(ctx, req.Address.FullAddress, req.Address.Apt)
	if err != nil {
		return errors.Annotate(err, "failed to sub.GetAddress")
	}
	subscriber.DeliveryNotes = req.DeliveryNotes
	subscriber.EmailPrefs[0].FirstName = req.FirstName
	subscriber.EmailPrefs[0].LastName = req.LastName
	subscriber.PhonePrefs = make([]subold.PhonePref, 0)
	subscriber.AddPhoneNumber(req.PhoneNumber)
	subscriber.Address = *address
	err = subC.Update(subscriber)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Activate")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

// UpdatePayment updates the Payment.
func (s *server) UpdatePayment(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.UpdatePaymentReq)
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
	subscriber, err := subC.Get(user.ID)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Get")
	}
	paymentC := payment.New(ctx)
	paymentReq := &payment.CreateCustomerReq{
		CustomerID: subscriber.PaymentCustomerID,
		FirstName:  subscriber.FirstName(),
		LastName:   subscriber.LastName(),
		Email:      subscriber.Email(),
		Nonce:      req.PaymentMethodNonce,
	}

	paymenttkn, err := paymentC.CreateCustomer(paymentReq)
	if err != nil {
		return errors.Wrap("failed to payment.CreateCustomer", err)
	}

	err = subC.UpdatePaymentToken(subscriber.Email(), paymenttkn)
	if err != nil {
		return errors.Wrap("failed to subold.UpdatePaymentToken", err)
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

// ChangePlanDay updates the plan day for a subscriber.
func (s *server) ChangePlanDay(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.ChangePlanDayReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.NewPlanDay != time.Monday.String() && req.NewPlanDay != time.Thursday.String() {
		return errors.BadRequestError.Annotate("Plan day must be Monday or Thursday.")
	}
	intervalStartDate := time.Now().Add(81 * time.Hour)
	for intervalStartDate.Weekday().String() != req.NewPlanDay {
		intervalStartDate = intervalStartDate.Add(time.Hour * 24)
	}
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}

	err = subC.ChangePlanDay(user.ID, req.NewPlanDay, &intervalStartDate)
	if err != nil {
		return errors.Annotate(err, "failed to sub.ChangePlanDay")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

package main

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbsub"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetUserSummary gets the user summary.
func (s *server) GetUserSummary(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.GetUserSummaryReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pbsub.GetUserSummaryResp{}
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	var subscriber *subold.Subscriber
	if user != nil {
		resp.IsLoggedIn = true
		subscriber, err = subC.GetByEmail(user.Email)
		if err != nil {
			ewc := errors.GetErrorWithCode(err)
			if ewc.Code != errors.CodeNotFound {
				return errors.GetErrorWithCode(err).Annotate("failed to Sub.GetByEmail")
			}
		} else {
			resp.IsActive = subscriber.Active
			resp.HasSubscribed = !subscriber.SignUpDatetime.IsZero()
			// TODO: add probation info
			resp.OnProbation = false
		}
	}
	return resp
}

// GetAccountInfo gets the user's account info.
func (s *server) GetAccountInfo(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.GetAccountInfoReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pbsub.GetAccountInfoResp{}
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	var subscriber *subold.Subscriber
	if user != nil {
		subscriber, err = subC.GetByEmail(user.Email)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get subscriber")
		}
		resp.Address, err = serverhelper.PBAddress(&subscriber.Address)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to PBAddress")
		}
		resp.EmailPrefs, err = serverhelper.PBEmailPrefs(subscriber.EmailPrefs)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to PBEmailPrefs")
		}
		resp.PhonePrefs, err = serverhelper.PBPhonePrefs(subscriber.PhonePrefs)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to PBPhonePrefs")
		}
		resp.PaymentInfo = &pbsub.PaymentInfo{
			CardNumberPreview: "Visa card ending in ****1234",
			CardType:          "Visa",
		}
	}
	return resp
}

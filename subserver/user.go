package main

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/sub"

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
	if user.ID == "" {
		subscriber, err = subC.Get(user.ID)
	} else {
		subscriber, err = subC.GetByEmail(user.Email)
	}
	if err != nil {
		ewc := errors.GetErrorWithCode(err)
		if ewc.Code != errors.CodeNotFound {
			return errors.GetErrorWithCode(err).Annotate("failed to get all sub")
		}
	} else {
		resp.IsActive = subscriber.Active
		// if subscriber.PaymentMethodToken != "" {
		// resp.HasPayment = true
		// }
	}
	// TODO: add probation info
	return resp
}

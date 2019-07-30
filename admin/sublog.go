package main

import (
	"context"
	"net/http"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetUnpaidSublogs gets a list of unpaid sublogs log.
func (s *server) GetUnpaidSublogs(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetUnpaidSublogsReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	subC := subold.New(ctx)
	sublogs, err := subC.GetUnpaidSublogs(req.Limit)
	if err != nil {
		return errors.Annotate(err, "failed to subold.GetUnpaidSublogs")
	}
	resp := &pb.GetUnpaidSublogsResp{
		Sublogs: pbSublogs(sublogs),
	}
	return resp
}

func (s *server) GetSubscriberSublogs(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetSubscriberSublogsReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	email := req.Email
	subC := subold.New(ctx)
	sublogs, err := subC.GetSubscriberSublogs(email)

	if err != nil {
		return errors.Annotate(err, "failed to subold.GetSubscriberSublogs")
	}
	resp := &pb.GetSubscriberSublogsResp{
		Sublogs: pbSublogs(sublogs),
	}
	return resp
}

func pbSublogs(sublogs []*subold.SubscriptionLog) []*pb.Sublog {
	sls := make([]*pb.Sublog, len(sublogs))
	for i := range sublogs {
		sls[i] = pbSublog(sublogs[i])
	}
	return sls
}

func pbSublog(sublog *subold.SubscriptionLog) *pb.Sublog {
	return &pb.Sublog{
		Date:               sublog.Date.Format(time.RFC3339),
		SubEmail:           sublog.SubEmail,
		CreatedDatetime:    sublog.CreatedDatetime.Format(time.RFC3339),
		Skip:               sublog.Skip,
		Servings:           int32(sublog.Servings),
		VegServings:        int32(sublog.VegServings),
		Amount:             sublog.Amount,
		AmountPaid:         sublog.AmountPaid,
		Paid:               sublog.Paid,
		PaidDatetime:       sublog.PaidDatetime.Format(time.RFC3339),
		PaymentMethodToken: sublog.PaymentMethodToken,
		TransactionID:      sublog.TransactionID,
		Free:               sublog.Free,
		DiscountAmount:     sublog.DiscountAmount,
		DiscountPercent:    int32(sublog.DiscountPercent),
		CustomerID:         sublog.CustomerID,
		Refunded:           sublog.Refunded,
		RefundedAmount:     sublog.RefundedAmount,
	}
}

package admin

import (
	"context"
	"net/http"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// RefundAndSkipSublog runs subold.RefundAndSkip.
func (s *server) RefundAndSkipSublog(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.RefundAndSkipSublogReq)
	var err error

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	err = activityC.RefundAndSkip(getDatetime(req.Date), req.Email)
	if err != nil {
		return errors.Annotate(err, "failed to activity.RefundAndSkip")
	}
	return nil
}

// ProcessSublog runs sub.Process.
func (s *server) ProcessSublog(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.ProcessSublogsReq)
	var err error

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	date := getDatetime(req.Date)
	subC := subold.NewWithLogging(ctx, log)
	err = subC.Process(date, req.Email)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to sub.Process")
	}
	resp := &pb.ProcessSublogsResp{}
	return resp
}

// GetUnpaidSublogs gets a list of unpaid sublogs log.
func (s *server) GetUnpaidSublogs(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetUnpaidSublogsReq)
	var err error

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
	req := new(pb.GetSubscriberSublogsReq)
	var err error

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
		TransactionId:      sublog.TransactionID,
		Free:               sublog.Free,
		DiscountAmount:     sublog.DiscountAmount,
		DiscountPercent:    int32(sublog.DiscountPercent),
		CustomerId:         sublog.CustomerID,
		Refunded:           sublog.Refunded,
	}
}

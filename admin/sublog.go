package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/gorilla/schema"
)

// GetUnpaidSublogs gets a list of unpaid sublogs log.
func GetUnpaidSublogs(ctx context.Context, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetUnpaidSublogsReq)
	var err error
	// decode request
	if r.Method == "GET" {
		decoder := schema.NewDecoder()
		err := decoder.Decode(req, r.URL.Query())
		if err != nil {
			return failedToDecode(err)
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			return failedToDecode(err)
		}
		defer closeRequestBody(r)
	}
	logging.Infof(ctx, "Request: %+v", req)
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

package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"

	"github.com/atishpatel/Gigamunch-Backend/core/deliveries"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/gorilla/schema"
)

// GetDeliveries gets deliveries for a date.
func GetDeliveries(ctx context.Context, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetDeliveriesReq)
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

	deliveriesC, err := deliveries.NewClient(ctx, log)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotatef("failed to deliveries.NewClient")
	}
	delivs, err := deliveriesC.Get(req.Date, req.DriverEmail)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to deliveries.Get")
	}
	subEmails := make([]string, len(delivs))
	for i := range delivs {
		subEmails[i] = delivs[i].SubEmail
	}
	subC := subold.New(ctx)
	subs, err := subC.GetSubscribers(subEmails)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to sub.GetSubscribers")
	}
	t, _ := time.Parse("2006-01-02", req.Date)
	sublogs, err := subC.GetForDate(t)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to sub.GetForDate")
	}
	firstMap := make(map[string]bool)
	for _, s := range sublogs {
		if s.Free {
			firstMap[s.SubEmail] = true
		}
	}

	resp := &pb.GetDeliveriesResp{
		Deliveries: pbDeliveries(delivs, subs, firstMap),
	}
	return resp
}

// UpdateDeliveries updates deliveries.
func UpdateDeliveries(ctx context.Context, r *http.Request, log *logging.Client) Response {
	req := new(pb.UpdateDeliveriesReq)
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
	deliveriesC, err := deliveries.NewClient(ctx, log)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotatef("failed to deliveries.NewClient")
	}
	err = deliveriesC.Update(deliveriesPB(req.Deliveries))
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to update deliveries")
	}
	return nil
}

func pbDeliveries(delivs []*deliveries.Delivery, subs []*subold.SubscriptionSignUp, firstMap map[string]bool) []*pb.Delivery {
	ds := make([]*pb.Delivery, len(delivs))
	for i := range delivs {
		ds[i] = pbDelivery(delivs[i], subs[i], firstMap)
	}
	return ds
}

func pbDelivery(d *deliveries.Delivery, sub *subold.SubscriptionSignUp, firstMap map[string]bool) *pb.Delivery {
	vegetarian := false
	if sub.VegetarianServings > 0 {
		vegetarian = true
	}
	return &pb.Delivery{
		Date:          d.Date,
		DriverName:    d.DriverName,
		DriverEmail:   d.DriverEmail,
		SubEmail:      d.SubEmail,
		Order:         int32(d.Order),
		Success:       d.Success,
		Fail:          d.Fail,
		SubName:       sub.Name,
		PhoneNumber:   sub.PhoneNumber,
		Address:       pbAddress(&sub.Address),
		DeliveryNotes: sub.DeliveryTips,
		Servings:      int32(sub.Servings + sub.VegetarianServings),
		Vegetarian:    vegetarian,
		First:         firstMap[sub.Email],
	}
}

func deliveriesPB(ds []*pb.Delivery) []*deliveries.Delivery {
	delivs := make([]*deliveries.Delivery, len(ds))
	for i := range ds {
		delivs[i] = deliveryPB(ds[i])
	}
	return delivs
}

func deliveryPB(d *pb.Delivery) *deliveries.Delivery {
	return &deliveries.Delivery{
		Date:        d.Date,
		DriverName:  d.DriverName,
		DriverEmail: d.DriverEmail,
		SubEmail:    d.SubEmail,
		Order:       int(d.Order),
		Success:     d.Success,
		Fail:        d.Fail,
	}
}

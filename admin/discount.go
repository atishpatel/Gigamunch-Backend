package main

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/core/discount"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// DiscountSubscriber gives a subscriber a discount.
func (s *server) DiscountSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.DiscountSubscriberReq)

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
		return errors.Annotate(err, "failed to get subscriber")
	}

	discountC, err := discount.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to discount.NewClient")
	}
	err = discountC.Create(&discount.CreateReq{
		UserID:          subscriber.ID,
		Email:           subscriber.Email(),
		FirstName:       subscriber.FirstName(),
		LastName:        subscriber.LastName(),
		DiscountAmount:  req.DiscountAmount,
		DiscountPercent: int8(req.DiscountPercent),
	})
	if err != nil {
		return errors.Annotate(err, "failed to discount.Create")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// GetSubscriberDiscounts gets all the discounts for a subscriber.
func (s *server) GetSubscriberDiscounts(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
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

	discountC, err := discount.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to discount.NewClient")
	}
	discs, err := discountC.GetAllForUser(subscriber.ID)
	if err != nil {
		return errors.Annotate(err, "failed to discount.GetAllForUser")
	}
	discsResp, err := serverhelper.PBDiscounts(discs)
	if err != nil {
		return errors.Annotate(err, "failed to serverhelper.PBDiscounts")
	}
	resp := &pb.GetSubscriberDiscountsResp{
		Discounts: discsResp,
	}
	return resp
}

package main

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// SkipActivity gets a log.
func (s *server) SkipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.SkipActivityReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.Skip(getDatetime(req.Date), req.Email, "Admin skip.")
	if err != nil {
		return errors.Annotate(err, "failed to activity.SkipActivity")
	}
	resp := &pb.SkipActivityResp{}
	return resp
}

// UnskipActivity gets a log.
func (s *server) UnskipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.UnskipActivityReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.Unskip(getDatetime(req.Date), req.Email)
	if err != nil {
		return errors.Annotate(err, "failed to activity.UnskipActivity")
	}
	resp := &pb.UnskipActivityResp{}
	return resp
}

// RefundAndSkipActivity gets a log.
func (s *server) RefundAndSkipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.RefundAndSkipActivityReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.RefundAndSkip(getDatetime(req.Date), req.Email, req.Amount, req.Percent)
	if err != nil {
		return errors.Annotate(err, "failed to activity.RefundAndSkipActivity")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// RefundActivity gets a log.
func (s *server) RefundActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.RefundActivityReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.Refund(getDatetime(req.Date), req.Email, req.Amount, req.Percent)
	if err != nil {
		return errors.Annotate(err, "failed to activity.RefundActivity")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

package main

import (
	"context"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

// GetSubscriberActivities gets all activties for a subscriber.
func (s *server) GetSubscriberActivities(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.UserIDReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	activities, err := activityC.GetAllForUser(req.ID)
	if err != nil {
		return errors.Annotate(err, "failed to get activities")
	}

	aResp, err := serverhelper.PBActivities(activities)
	if err != nil {
		return errors.Annotate(err, "failed to PBActivities")
	}

	resp := &pb.GetSubscriberActivitiesResp{
		Activities: aResp,
	}

	return resp
}

// SkipActivity skips an activity.
func (s *server) SkipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.SkipActivityReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.ID == "" {
		req.ID = req.Email
	}
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.Skip(getDatetime(req.Date), req.ID, "Admin skip.")
	if err != nil {
		return errors.Annotate(err, "failed to activity.SkipActivity")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// UnskipActivity unskips an activity.
func (s *server) UnskipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.UnskipActivityReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.ID == "" {
		req.ID = req.Email
	}
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to activity.NewClient")
	}
	err = activityC.Unskip(getDatetime(req.Date), req.ID)
	if err != nil {
		return errors.Annotate(err, "failed to activity.UnskipActivity")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// SetupActivities setups activities
func (s *server) SetupActivities(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.SetupActivitiesReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	date := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	err = subC.SetupActivities(date)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.SetupActivities(Date:%v). Err:%+v", date, err)
		return errors.Annotate(err, "failed to sub.SetupActivities")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// SetupActivity setup an activity.
func (s *server) SetupActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.SetupActivityReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	date := getDatetime(req.Date)
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	err = subC.SetupActivity(date, req.ID, true, 0, 0)
	if err != nil {
		return errors.Annotate(err, "failed to sub.SetupActivity")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// RefundAndSkipActivity refunds and skips an activity.
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

// RefundActivity refunds an activity.
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
	date := getDatetime(req.Date)
	for _, e := range req.Emails {
		err = activityC.Refund(date, e, req.Amount, req.Percent)
		if err != nil {
			return errors.Annotate(err, "failed to activity.RefundActivity")
		}
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// ChangeActivityServings changes the servings for an activity.
func (s *server) ChangeActivityServings(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.ChangeActivityServingsReq)
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
	date := getDatetime(req.Date)
	err = activityC.ChangeServings(date, req.ID, int8(req.ServingsNonVeg), int8(req.ServingsVeg), sub.DerivePrice(int8(req.ServingsNonVeg+req.ServingsVeg)))
	if err != nil {
		return errors.Annotate(err, "failed to activity.ChangeServings")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}

// ProcessActivity processes an activity.
func (s *server) ProcessActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.ProcessActivityReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.ID == "" {
		req.ID = req.Email
	}

	date := getDatetime(req.Date)
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to sub.NewClient")
	}
	err = subC.ProcessActivity(date, req.ID)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to sub.Process")
	}
	resp := &pb.ProcessActivityResp{}
	return resp
}

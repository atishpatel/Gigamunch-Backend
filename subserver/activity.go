package main

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbsub"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func (s *server) SkipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.DateReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to activity.NewClient")
	}
	t := getDatetime(req.Date)
	err = activityC.Skip(t, user.Email, "")
	if err != nil {
		return errors.Annotate(err, "failed to activity.Skip")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

func (s *server) UnskipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.DateReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to activity.NewClient")
	}
	t := getDatetime(req.Date)
	err = activityC.Unskip(t, user.Email)
	if err != nil {
		return errors.Annotate(err, "failed to activity.UnSskip")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

// ChangeActivityServings changes the servings for an activity.
func (s *server) ChangeActivityServings(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.ChangeActivityServingsReq)
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
	err = activityC.ChangeServings(date, user.ID, int8(req.ServingsNonVeg), int8(req.ServingsVeg), sub.DerivePrice(int8(req.ServingsNonVeg+req.ServingsVeg)))
	if err != nil {
		return errors.Annotate(err, "failed to activity.ChangeServings")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

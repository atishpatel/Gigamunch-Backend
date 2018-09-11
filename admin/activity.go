package admin

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// SkipActivity gets a log.
func (s *server) SkipActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.SkipActivityReq)
	var err error

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
	req := new(pb.UnskipActivityReq)
	var err error

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

package main

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbsub"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// ChangeSubscriberServings change a subscriber's servings.
func (s *server) ChangeSubscriberServings(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pbsub.ChangeSubscriberServingsReq)

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
	err = subC.ChangeServingsPermanently(req.ID, int8(req.ServingsNonVeg), int8(req.ServingsVeg))
	if err != nil {
		return errors.Annotate(err, "failed to sub.ChangeServingsPermanently")
	}
	resp := &pbsub.ErrorOnlyResp{}
	return resp
}

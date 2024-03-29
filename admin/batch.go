package main

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
)

// UpdateSubs updates Subs for subscribers.
func (s *server) UpdateSubs(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pbadmin.GetExecutionsReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	log.Infof(ctx, "updating subs")
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	err = subC.BatchUpdateActivityWithUserID(int(req.Start), int(req.Limit))
	if err != nil {
		return errors.Annotate(err, "failed to sub.BatchUpdateActivityWithUserID")
	}
	return nil
}

// MigrateToNewSubscribersStruct migrates subscribers to new struct.
func (s *server) MigrateToNewSubscribersStruct(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {

	return nil
}

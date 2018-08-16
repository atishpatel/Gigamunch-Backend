package admin

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// Login takes a Firebase Token and returns an auth token.
func (s *server) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.TokenOnlyReq)
	var err error

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	authC, err := auth.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to get auth.NewClient")
	}
	_, token, err := authC.GetFromFBToken(ctx, req.Token)
	if err != nil {
		return errors.Annotate(err, "failed to get auth.GetFromFBToken")
	}
	resp := &pb.TokenOnlyResp{
		Token: token,
	}
	return resp
}

// Refresh refreshs an auth token.
func (s *server) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.TokenOnlyReq)
	var err error

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	authC, err := auth.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to get auth.NewClient")
	}
	token, err := authC.Refresh(ctx, req.Token)
	if err != nil {
		return errors.Annotate(err, "failed to get auth.Refresh")
	}
	resp := &pb.TokenOnlyResp{
		Token: token,
	}
	return resp
}

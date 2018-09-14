package admin

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// MakeAdmin makes a user an admin.
func (s *server) MakeAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.MakeAdminReq)
	var err error
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	authC, err := auth.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get auth.NewClient")
	}
	err = authC.MakeAdmin(req.Email)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to auth.MakeAdmin")
	}
	return nil
}

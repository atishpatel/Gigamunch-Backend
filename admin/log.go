package admin

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	pbcommon "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetLog gets a log.
func (s *server) GetLog(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetLogReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	l, err := log.Get(req.Id)
	if err != nil {
		return errors.Annotate(err, "failed to log.GetLogs")
	}
	resp := &pb.GetLogResp{
		Log: pbLog(l),
	}
	return resp
}

// GetLogs gets logs.
func (s *server) GetLogs(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetLogsReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	logs, err := log.GetAll(int(req.Start), int(req.Limit))
	if err != nil {
		return errors.Annotate(err, "failed to log.GetLogs")
	}
	resp := &pb.GetLogsResp{}
	resp.Logs = make([]*pbcommon.Log, len(logs))
	for i := range logs {
		resp.Logs[i] = pbLog(logs[i])
	}
	return resp
}

// GetLogsByEmail gets logs by email.
func (s *server) GetLogsByEmail(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetLogsByEmailReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	logs, err := log.GetAllByEmail(req.Email, int(req.Start), int(req.Limit))
	if err != nil {
		return errors.Annotate(err, "failed to log.GetUserLogsByEmail")
	}
	resp := &pb.GetLogsByEmailResp{}
	resp.Logs = make([]*pbcommon.Log, len(logs))
	for i := range logs {
		resp.Logs[i] = pbLog(logs[i])
	}
	return resp
}

func pbLog(l *logging.Entry) *pbcommon.Log {
	return &pbcommon.Log{
		Id:              l.ID,
		LogName:         l.LogName,
		Timestamp:       l.Timestamp.String(),
		Type:            string(l.Type),
		Action:          string(l.Action),
		Path:            l.Path,
		Severity:        int32(l.Severity),
		UserId:          l.UserID,
		UserEmail:       l.UserEmail,
		ActionUserId:    l.ActionUserID,
		ActionUserEmail: l.ActionUserEmail,
		BasicPayload: &pbcommon.BasicPayload{
			Title:       l.BasicPayload.Title,
			Description: l.BasicPayload.Description,
		},
	}
}

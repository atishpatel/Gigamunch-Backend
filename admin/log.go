package admin

import (
	"context"
	"encoding/json"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	shared "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/gorilla/schema"
)

// GetLog gets a log.
func GetLog(ctx context.Context, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetLogReq)
	var err error
	// decode request
	if r.Method == "GET" {
		decoder := schema.NewDecoder()
		err := decoder.Decode(req, r.URL.Query())
		if err != nil {
			return failedToDecode(err)
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			return failedToDecode(err)
		}
		defer closeRequestBody(r)
	}
	logging.Infof(ctx, "Request: %+v", req)
	// end decode request
	l, err := log.GetLog(req.Id)
	if err != nil {
		return errors.Annotate(err, "failed to log.GetLogs")
	}
	resp := &pb.GetLogResp{
		Log: pbLog(l),
	}
	return resp
}

// GetLogs gets logs.
func GetLogs(ctx context.Context, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetLogsReq)
	var err error
	// decode request
	if r.Method == "GET" {
		decoder := schema.NewDecoder()
		err := decoder.Decode(req, r.URL.Query())
		if err != nil {
			return failedToDecode(err)
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			return failedToDecode(err)
		}
		defer closeRequestBody(r)
	}
	logging.Infof(ctx, "Request: %+v", req)
	// end decode request
	logs, err := log.GetLogs(int(req.Start), int(req.Limit))
	if err != nil {
		return errors.Annotate(err, "failed to log.GetLogs")
	}
	resp := &pb.GetLogsResp{}
	resp.Logs = make([]*shared.Log, len(logs))
	for i := range logs {
		resp.Logs[i] = pbLog(logs[i])
	}
	return resp
}

// GetLogsByEmail gets logs by email.
func GetLogsByEmail(ctx context.Context, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetLogsByEmailReq)
	var err error
	// decode request
	if r.Method == "GET" {
		decoder := schema.NewDecoder()
		err := decoder.Decode(req, r.URL.Query())
		if err != nil {
			return failedToDecode(err)
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			return failedToDecode(err)
		}
		defer closeRequestBody(r)
	}
	logging.Infof(ctx, "Request: %+v", req)
	// end decode request
	logs, err := log.GetUserLogsByEmail(req.Email, int(req.Start), int(req.Limit))
	if err != nil {
		return errors.Annotate(err, "failed to log.GetUserLogsByEmail")
	}
	resp := &pb.GetLogsByEmailResp{}
	resp.Logs = make([]*shared.Log, len(logs))
	for i := range logs {
		resp.Logs[i] = pbLog(logs[i])
	}
	return resp
}

func pbLog(l *logging.Entry) *shared.Log {
	return &shared.Log{
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
		BasicPayload: &shared.BasicPayload{
			Title:       l.BasicPayload.Title,
			Description: l.BasicPayload.Description,
		},
	}
}

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

func pbLog(l *logging.Entry) *shared.Log {
	labels := make([]string, len(l.Labels))
	for i := range l.Labels {
		labels[i] = string(l.Labels[i])
	}
	return &shared.Log{
		Id:        l.ID,
		LogName:   l.LogName,
		Timestamp: l.Timestamp.String(),
		Type:      string(l.Type),
		// Path: l.Path,
		Labels:   labels,
		Severity: int32(l.Severity),
		Payload:  l.Payload,
	}
}

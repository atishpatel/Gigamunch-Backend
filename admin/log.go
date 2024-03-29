package main

import (
	"context"
	"net/http"
	"strings"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"
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

	l, err := log.Get(req.ID)
	if err != nil {
		return errors.Annotate(err, "failed to log.GetLogs")
	}
	resp := &pb.GetLogResp{}
	resp.Log, err = serverhelper.PBLog(l)
	if err != nil {
		return errors.Annotate(err, "failed to PBLog")
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
	resp.Logs, err = serverhelper.PBLogs(logs)
	if err != nil {
		return errors.Annotate(err, "failed to PBLogs")
	}
	return resp
}

// GetLogsByAction gets logs.
func (s *server) GetLogsByAction(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetLogsByActionReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	action := logging.Action(req.Action)
	logs, err := log.GetAllByAction(action, int(req.Start), int(req.Limit))
	if err != nil {
		return errors.Annotate(err, "failed to log.GetLogs")
	}
	resp := &pb.GetLogsResp{}
	resp.Logs, err = serverhelper.PBLogs(logs)
	if err != nil {
		return errors.Annotate(err, "failed to PBLogs")
	}
	return resp
}

// GetLogsForUser gets logs by email.
func (s *server) GetLogsForUser(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetLogsForUserReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	var logs []*logging.Entry
	if strings.Contains(req.ID, "@") {
		logs, err = log.GetAllByEmail(req.ID, int(req.Start), int(req.Limit))
	} else {
		logs, err = log.GetAllByID(req.ID, int(req.Start), int(req.Limit))
	}
	if err != nil {
		return errors.Annotate(err, "failed to log.GetAllBy")
	}
	resp := &pb.GetLogsResp{}
	resp.Logs, err = serverhelper.PBLogs(logs)
	if err != nil {
		return errors.Annotate(err, "failed to PBLogs")
	}
	return resp
}

// GetLogsByExecution gets logs by email.
func (s *server) GetLogsByExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetLogsByExecutionReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	logs, err := log.GetAllByExecution(req.ExecutionID)
	if err != nil {
		return errors.Annotate(err, "failed to log.GetAllByExecution")
	}
	resp := &pb.GetLogsResp{}
	resp.Logs, err = serverhelper.PBLogs(logs)
	if err != nil {
		return errors.Annotate(err, "failed to PBLogs")
	}
	return resp
}

package main

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbsub"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func getExecutionsResp(exes []*execution.Execution, activities []*activity.Activity) (*pbsub.GetExecutionsResp, error) {
	resp := &pbsub.GetExecutionsResp{}
	if exes == nil {
		return resp, nil
	}
	exespb, err := serverhelper.PBExecutions(exes)
	if err != nil {
		return nil, errors.GetErrorWithCode(err).Annotate("failed to PBExecutions")
	}
	var activitiespb []*pbcommon.Activity
	if activities != nil {
		activitiespb, err = serverhelper.PBActivities(activities)
		if err != nil {
			return nil, errors.GetErrorWithCode(err).Annotate("failed to PBActivities")
		}
	}
	resp.ExecutionAndActivity = make([]*pbcommon.ExecutionAndActivity, len(exespb))
	hasActivtites := activitiespb != nil
	for i := range exespb {
		resp.ExecutionAndActivity[i] = &pbcommon.ExecutionAndActivity{
			Execution: exespb[i],
		}
		if hasActivtites && exespb[i] != nil {
			for j := range activitiespb {
				if activitiespb[j] != nil && strings.Contains(activitiespb[j].Date, exespb[i].Date) {
					resp.ExecutionAndActivity[i].Activity = activitiespb[j]
					break
				}
			}
		}
	}
	if hasActivtites {
		resp.Activities = activitiespb
	}
	return resp, nil
}

// GetExecutions gets list of executions.
func (s *server) GetExecutions(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.GetExecutionsReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution.NewClient")
	}
	executions, err := exeC.GetAll(int(req.Start), int(req.Limit))
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get all executions")
	}
	var activities []*activity.Activity
	if user != nil {
		activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.NewClient")
		}
		activities, err = activityC.GetAllForUser(user.Email)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.GetAllForUser")
		}
	}
	resp, err := getExecutionsResp(executions, activities)
	if err != nil {
		return errors.GetErrorWithCode(err)
	}
	return resp
}

// GetExecutionsAfterDate gets list of executions after date.
func (s *server) GetExecutionsAfterDate(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.GetExecutionsDateReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution.NewClient")
	}
	executions, err := exeC.GetAfterDate(getDatetime(req.Date))
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution.GetAfterDate")
	}
	publishedExecutions := filterPublishedExecutions(executions)
	sort.Slice(publishedExecutions, func(i, j int) bool {
		if publishedExecutions[i].Date == "" {
			return false
		}
		if publishedExecutions[j].Date == "" {
			return true
		}
		return publishedExecutions[i].Date < publishedExecutions[j].Date
	})

	var activities []*activity.Activity
	if user != nil {
		activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.NewClient")
		}
		activities, err = activityC.GetAfterDateForUser(getDatetime(req.Date), user.Email)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.GetAfterDateForUser")
		}
	}
	resp, err := getExecutionsResp(publishedExecutions, activities)
	if err != nil {
		return errors.GetErrorWithCode(err)
	}
	return resp
}

// GetExecutionsBeforeDate gets list of executions before date.
func (s *server) GetExecutionsBeforeDate(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.GetExecutionsDateReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution.NewClient")
	}
	executions, err := exeC.GetBeforeDate(getDatetime(req.Date))
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution.GetBeforeDate")
	}
	publishedExecutions := filterPublishedExecutions(executions)
	sort.Slice(publishedExecutions, func(i, j int) bool {
		if publishedExecutions[i].Date == "" {
			return false
		}
		if publishedExecutions[j].Date == "" {
			return true
		}
		return publishedExecutions[i].Date < publishedExecutions[j].Date
	})

	var activities []*activity.Activity
	if user != nil {
		activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.NewClient")
		}
		activities, err = activityC.GetBeforeDateForUser(getDatetime(req.Date), user.Email)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.GetBeforeDateForUser")
		}
	}
	resp, err := getExecutionsResp(publishedExecutions, activities)
	if err != nil {
		return errors.GetErrorWithCode(err)
	}
	return resp
}

// filterPublishedExecutions takes in executions, returns only published executions
func filterPublishedExecutions(executions []*execution.Execution) []*execution.Execution {

	var publishedExecutions []*execution.Execution
	for _, execution := range executions {
		if execution.Publish {
			publishedExecutions = append(publishedExecutions, execution)
		}
	}
	return publishedExecutions
}

// GetExecution gets an execution.
func (s *server) GetExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client, user *common.User) Response {
	var err error
	req := new(pbsub.GetExecutionReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	// get execution
	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	exe, err := exeC.Get(req.IDOrDate)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution")
	}

	// get activity
	var act *activity.Activity
	var acts []*activity.Activity
	if user != nil && exe.Date != "" {
		activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.NewClient")
		}
		t, err := time.Parse(execution.DateFormat, exe.Date)
		if err != nil {
			return errInternalError.WithError(err).Annotate("failed to get time.Parse")
		}
		act, err = activityC.Get(t, user.Email)
		if err != nil {
			ewc := errors.GetErrorWithCode(err)
			if ewc.Code != errors.CodeNotFound {
				return ewc.Annotate("failed to get activity.GetBeforeDateForUser")
			}
		}
		if act != nil {
			acts = []*activity.Activity{act}
		}
	}
	exesResp, err := getExecutionsResp([]*execution.Execution{exe}, acts)
	if err != nil {
		return errors.GetErrorWithCode(err)
	}
	resp := &pbsub.GetExecutionResp{
		ExecutionAndActivity: exesResp.ExecutionAndActivity[0],
	}
	return resp
}

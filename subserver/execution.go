package main

import (
	"context"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/sub"

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
		if hasActivtites {
			for j := range activitiespb {
				if exespb[i].Date == activitiespb[j].Date {
					resp.ExecutionAndActivity[i].Activity = activitiespb[j]
					break
				}
			}
		}
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
		activities, err = activityC.GetAllForUser(user.ID)
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
	var activities []*activity.Activity
	if user != nil {
		activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.NewClient")
		}
		activities, err = activityC.GetAfterDateForUser(getDatetime(req.Date), user.ID)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.GetAfterDateForUser")
		}
	}
	resp, err := getExecutionsResp(executions, activities)
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
	var activities []*activity.Activity
	if user != nil {
		activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.NewClient")
		}
		activities, err = activityC.GetBeforeDateForUser(getDatetime(req.Date), user.ID)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.GetBeforeDateForUser")
		}
	}
	resp, err := getExecutionsResp(executions, activities)
	if err != nil {
		return errors.GetErrorWithCode(err)
	}
	return resp
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
	if user != nil && exe.Date != "" {
		activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.NewClient")
		}
		t, err := time.Parse(execution.DateFormat, exe.Date)
		if err != nil {
			return errInternalError.WithError(err).Annotate("failed to get time.Parse")
		}
		act, err = activityC.Get(t, user.ID)
		if err != nil {
			return errors.GetErrorWithCode(err).Annotate("failed to get activity.GetBeforeDateForUser")
		}
	}
	exesResp, err := getExecutionsResp([]*execution.Execution{exe}, []*activity.Activity{act})
	if err != nil {
		return errors.GetErrorWithCode(err)
	}
	resp := &pbsub.GetExecutionResp{
		ExecutionAndActivity: exesResp.ExecutionAndActivity[0],
	}
	return resp
}

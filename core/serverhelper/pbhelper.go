package serverhelper

import (
	"bytes"
	"encoding/json"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	errMarshalUnmarshal = errors.InternalServerError
)

func marshalUnmarshal(oldType interface{}, newType interface{}) error {
	// marshal old struct
	inputJSON, err := json.Marshal(oldType)
	if err != nil {
		return errMarshalUnmarshal.WithError(err).Annotate("failed to json.Marshal")
	}
	// fix empty datetime strings
	inputJSON = bytes.Replace(inputJSON, []byte("datetime\":\"\","), []byte("datetime\":\"0001-01-01T00:00:00Z\","), -1)
	// unmarshal into new struct
	err = json.Unmarshal(inputJSON, newType)
	if err != nil {
		return errMarshalUnmarshal.WithError(err).Annotate("failed to json.Unmarshal")
	}
	return nil
}

// PBLogs turns an Log into a protobuff Log.
func PBLogs(in []*logging.Entry) ([]*pbcommon.Log, error) {
	out := make([]*pbcommon.Log, len(in))
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(&in, &out)
	return out, err
}

// PBLog turns an Log into a protobuff Log.
func PBLog(in *logging.Entry) (*pbcommon.Log, error) {
	out := &pbcommon.Log{}
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(in, out)
	return out, err
}

// PBActivities turns an array of Activities into a protobuff array of Activities.
func PBActivities(in []*activity.Activity) ([]*pbcommon.Activity, error) {
	out := make([]*pbcommon.Activity, len(in))
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(&in, &out)
	return out, err
}

// PBExecutions turns an array of executions into a protobuff array of executions.
func PBExecutions(in []*execution.Execution) ([]*pbcommon.Execution, error) {
	out := make([]*pbcommon.Execution, len(in))
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(&in, &out)
	return out, err
}

// PBExecution turns an execution into a protobuff executions.
func PBExecution(in *execution.Execution) (*pbcommon.Execution, error) {
	out := &pbcommon.Execution{}
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(in, out)
	return out, err
}

// ExecutionFromPb turns pbcommon.Execution to execution.Execution
func ExecutionFromPb(in *pbcommon.Execution) (*execution.Execution, error) {
	out := &execution.Execution{}
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(in, out)
	return out, err
}

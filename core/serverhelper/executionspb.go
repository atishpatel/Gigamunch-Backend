package serverhelper

import (
	"encoding/json"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	errMarshalUnmarshal = errors.InternalServerError
)

func marshalUnmarshal(oldType interface{}, newType interface{}) error {
	inputJSON, err := json.Marshal(oldType)
	if err != nil {
		return errMarshalUnmarshal.WithError(err).Annotate("failed to json.Marshal")
	}
	err = json.Unmarshal(inputJSON, newType)
	if err != nil {
		return errMarshalUnmarshal.WithError(err).Annotate("failed to json.Unmarshal")
	}
	return nil
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
	// fix empty time strings
	var t time.Time
	if in.CreatedDatetime == "" {
		in.CreatedDatetime = t.Format(time.RFC3339Nano)
	}
	err := marshalUnmarshal(in, out)
	return out, err
}

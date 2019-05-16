package serverhelper

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"
	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
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

// **********************
// type -> PB
// **********************

// PBAddress turns an Address into a protobuff Address.
func PBAddress(in *common.Address) (*pbcommon.Address, error) {
	out := &pbcommon.Address{}
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(in, out)
	return out, err
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
	for _, a := range in {
		a.PaidDatetimeJSON = a.PaidDatetime.Time.Format(time.RFC3339)
		a.RefundedDatetimeJSON = a.RefundedDatetime.Time.Format(time.RFC3339)
		a.GiftFromUserIDJSON = a.GiftFromUserID.Int64
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

// PBSubscribers turns an array of subscribers into a protobuff array of subscribers.
func PBSubscribers(in []*subold.Subscriber) ([]*pbcommon.Subscriber, error) {
	out := make([]*pbcommon.Subscriber, len(in))
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(&in, &out)
	return out, err
}

// PBSubscriber turns an subscriber into a protobuff subscriber.
func PBSubscriber(in *subold.Subscriber) (*pbcommon.Subscriber, error) {
	out := &pbcommon.Subscriber{}
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(in, out)
	return out, err
}

// PBEmailPrefs turns an array of EmailPrefs into a protobuff array of EmailPrefs.
func PBEmailPrefs(in []subold.EmailPref) ([]*pbcommon.EmailPref, error) {
	out := make([]*pbcommon.EmailPref, len(in))
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(&in, &out)
	return out, err
}

// PBPhonePrefs turns an array of PhonePrefs into a protobuff array of PhonePrefs.
func PBPhonePrefs(in []subold.PhonePref) ([]*pbcommon.PhonePref, error) {
	out := make([]*pbcommon.PhonePref, len(in))
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(&in, &out)
	return out, err
}

// **********************
// PB -> type
// **********************

// ExecutionFromPb turns pbcommon.Execution to execution.Execution
func ExecutionFromPb(in *pbcommon.Execution) (*execution.Execution, error) {
	out := &execution.Execution{}
	if in == nil {
		return out, nil
	}
	err := marshalUnmarshal(in, out)
	return out, err
}

package serverhelper

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/gorilla/schema"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
)

// GetDatetime converts string to datetime
func GetDatetime(s string) time.Time {
	if len(s) == 10 {
		s += "T12:12:12.000Z"
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// FailedToDecode creates a ErrorOnlyResp
func FailedToDecode(err error) *pbcommon.ErrorOnlyResp {
	return &pbcommon.ErrorOnlyResp{
		Error: errors.BadRequestError.WithError(err).Annotate("failed to decode").SharedError(),
	}
}

func DecodeRequest(ctx context.Context, r *http.Request, v interface{}) error {
	if r.Method == "GET" {
		decoder := schema.NewDecoder()
		err := decoder.Decode(v, r.URL.Query())
		logging.Infof(ctx, "Query: %+v", r.URL.Query())
		if err != nil {
			return err
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		logging.Infof(ctx, "Body: %s", body)
		err = json.Unmarshal(body, v)
		if err != nil {
			return err
		}
	}
	return nil
}

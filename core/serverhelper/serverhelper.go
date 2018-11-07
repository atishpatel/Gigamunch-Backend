package serverhelper

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/errors"

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

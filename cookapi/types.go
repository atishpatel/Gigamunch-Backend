package cookapi

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// Item is an Item.
type Item struct {
	ID int64 `json:"id,string"`
	item.Item
}

// GigatokenOnlyReq is a request with only a gigatoken input
type GigatokenOnlyReq struct {
	Gigatoken string `json:"gigatoken"`
}

func (req *GigatokenOnlyReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *GigatokenOnlyReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// ErrorOnlyResp is a response with only an error with code
type ErrorOnlyResp struct {
	Err errors.ErrorWithCode `json:"err"`
}

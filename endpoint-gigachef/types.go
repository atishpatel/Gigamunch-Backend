package gigachef

import (
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
)

// BaseItem is the basic stuff in an Item
type BaseItem struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Ingredients      []string `json:"ingredients"`
	GeneralTags      []string `json:"general_tags"`
	CuisineTags      []string `json:"cuisine_tags"`
	DietaryNeedsTags []string `json:"dietary_needs_tags"`
	Photos           []string `json:"photos"`
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

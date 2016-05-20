package gigachef

import "gitlab.com/atishpatel/Gigamunch-Backend/errors"

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

// ErrorOnlyResp is a response with only an error with code
type ErrorOnlyResp struct {
	Err errors.ErrorWithCode `json:"err"`
}

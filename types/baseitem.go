package types

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

// BaseItem is the basic stuff in a Item and Post
type BaseItem struct {
	GigachefID       string    `json:"gigachef_id"`
	CreatedDateTime  time.Time `json:"created_datetime"`
	Subtitle         string    `json:"subtitle"`
	Description      string    `json:"description"`
	GeneralTags      []string  `json:"general_tags"`
	DietaryNeedsTags []string  `json:"dietary_needs_tags"`
	CuisineTags      []string  `json:"cuisine_tags"`
	Ingredients      []string  `json:"ingredients"`
	Photos           []string  `json:"photos"`
}

// Validate validates the BaseItem properties.
// The form is valid if errors.Errors.HasErrors() == false.
func (baseItem *BaseItem) Validate() errors.Errors {
	var multipleErrors errors.Errors
	if baseItem.GigachefID == "" {
		multipleErrors.AddError(fmt.Errorf("GigachefID is empty"))
	}
	if baseItem.CreatedDateTime.Year() < 3 {
		multipleErrors.AddError(fmt.Errorf("CreatedDateTime is not set"))
	}
	if len(baseItem.Description) > 10 || utils.ContainsBanWord(baseItem.Description) {
		multipleErrors.AddError(fmt.Errorf("Description is too short"))
	}
	if len(baseItem.Photos) == 0 {
		multipleErrors.AddError(fmt.Errorf("Photos must be more than zero"))
	}
	return multipleErrors
}

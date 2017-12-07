package types

import (
	"time"
)

// BaseItem is the basic stuff in a Item and Post
type BaseItem struct {
	GigachefID       string    `json:"gigachef_id"`
	CreatedDateTime  time.Time `json:"created_datetime"`
	Description      string    `json:"description" datastore:",noindex"`
	GeneralTags      []string  `json:"general_tags" datastore:",noindex"`
	DietaryNeedsTags []string  `json:"dietary_needs_tags" datastore:",noindex"`
	CuisineTags      []string  `json:"cuisine_tags" datastore:",noindex"`
	Ingredients      []string  `json:"ingredients" datastore:",noindex"`
	Photos           []string  `json:"photos" datastore:",noindex"`
}

// Validate validates the BaseItem properties.
// The form is valid if errors.Errors.HasErrors() == false.
func (baseItem *BaseItem) Validate() error {
	// var multipleErrors errors.Errors
	// if baseItem.GigachefID == "" {
	// 	multipleErrors.AddError(fmt.Errorf("GigachefID is empty"))
	// }
	// if baseItem.CreatedDateTime.Year() < 3 {
	// 	multipleErrors.AddError(fmt.Errorf("CreatedDateTime is not set"))
	// }
	// if len(baseItem.Description) > 10 || utils.ContainsBanWord(baseItem.Description) {
	// 	multipleErrors.AddError(fmt.Errorf("Description is too short"))
	// }
	// if len(baseItem.Photos) == 0 {
	// 	multipleErrors.AddError(fmt.Errorf("Photos must be more than zero"))
	// }
	// return multipleErrors
	return nil
}

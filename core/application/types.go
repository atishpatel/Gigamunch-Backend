package application

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindChefApplication is an application for a chef to become verfified
	kindChefApplication = "ChefApplication"
)

type ChefApplication struct {
	UserID                string    `json:"user_id" datastore:",noindex"`
	CreatedDateTime       time.Time `json:"created_datetime" datastore:",noindex"`
	LastUpdatedDateTime   time.Time `json:"last_updated_datetime" datastore:",index"`
	types.ChefApplication           // embedded
}

package application

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindChefApplication is an application for a chef to become verfified
	kindChefApplication = "ChefApplication"
)

// ChefApplication is the chef application to become a verfified Gigachef
type ChefApplication struct {
	UserID                 string        `json:"user_id" datastore:",noindex"`
	CreatedDateTime        time.Time     `json:"created_datetime" datastore:",noindex"`
	LastUpdatedDateTime    time.Time     `json:"last_updated_datetime" datastore:",index"`
	Name                   string        `json:"name" datastore:",noindex"`
	Email                  string        `json:"email" datastore:",index"`
	PhoneNumber            string        `json:"phone_number" datastore:",noindex"`
	Address                types.Address `json:"address" datastore:",noindex"`
	AttendedCulinarySchool bool          `json:"attended_culinary_school" datastore:",noindex"`
	WorkedAtResturant      bool          `json:"worked_at_resturant" datastore:",noindex"`
	PostFrequency          int           `json:"post_frequency" datastore:",noindex"`
	ApplicationProgress    int           `json:"application_progress" datastore:",noindex"` // TODO remove
}

package types

type ChefApplication struct {
	Name                   string  `json:"name" datastore:",noindex"`
	Email                  string  `json:"email" datastore:",index"`
	PhoneNumber            string  `json:"phone_number" datastore:",noindex"`
	Address                Address `json:"address" datastore:",noindex"`
	AttendedCulinarySchool bool    `json:"attended_culinary_school" datastore:",noindex"`
	WorkedAtResturant      bool    `json:"worked_at_resturant" datastore:",noindex"`
	PostFrequency          int     `json:"post_frequency" datastore:",noindex"`
	ApplicationProgress    int     `json:"application_progress" datastore:",noindex"`
}

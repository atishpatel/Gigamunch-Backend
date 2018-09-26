package common

// User is a user. xP
type User struct {
	ID int64 `json:"id,omitempty"`
	// Firebase's User ID
	AuthID    string `json:"auth_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	Admin     bool   `json:"admin,omitempty"`
	// ActiveSubscriber bool   `json:"active_subscriber,omitempty"`
}

package common

// User is a user. xP
type User struct {
	ID int64 `json:"id"`
	// Firebase's User ID
	AuthID    string `json:"auth_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PhotoURL  string `json:"photo_url"`
	Admin     bool   `json:"admin"`
}

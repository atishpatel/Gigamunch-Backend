package types

// User information in a session.
type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"` // TODO(Atish): remove photourl if we are using a naming convention
	// bit1 isChef | bit2 isVerifiedChef |
	Permissions uint32 `json:"permissions"`
}

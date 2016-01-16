package types

// User information in a session.
type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"` // TODO(Atish): remove photourl if we are using a naming convention
	// bit1 isChef | bit2 isVerifiedChef |
	Permissions uint32 `json:"permissions"`
}

// IsChef returns true if a user is a chef
func (user *User) IsChef() bool {
	return getKthBit(user.Permissions, 0)
}

// IsVerifiedChef returns true if a user is a verified chef
func (user *User) IsVerifiedChef() bool {
	return getKthBit(user.Permissions, 1)
}

func getKthBit(num uint32, k uint32) bool {
	return (num>>k)&1 == 1
}

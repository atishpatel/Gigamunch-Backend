package types

// User information in a session.
type User struct {
	Email string `json:"email"`
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

// UserDetail is the structure that is stored in the database for a chef's
// or muncher's details
type UserDetail struct {
	Name       string `json:"name" datastore:",noindex"`
	PhotoURL   string `json:"photo_url" datastore:",noindex"`
	ProviderID string `json:"provider_id" datastore:",noindex"`
}

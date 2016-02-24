package types

// User information in a session.
type User struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
	// bit1 | bit2 | bit3 | ...
	// isChef | isVerifiedChef | isAdmin | HasAddress | HasCreditCardInfo
	Permissions int32 `json:"permissions"`
}

// IsChef returns true if a user is a chef
func (user *User) IsChef() bool {
	return getKthBit(user.Permissions, 0)
}

// SetChef updates the permission of the user
func (user *User) SetChef(x bool) {
	user.Permissions = setKthBit(user.Permissions, 0, x)
}

// IsVerifiedChef returns true if a user is a verified chef
func (user *User) IsVerifiedChef() bool {
	return getKthBit(user.Permissions, 1)
}

// SetVerifiedChef updates the permission of the user
func (user *User) SetVerifiedChef(x bool) {
	user.Permissions = setKthBit(user.Permissions, 1, x)
}

// IsAdmin returns true if a user is an admin
func (user *User) IsAdmin() bool {
	return getKthBit(user.Permissions, 2)
}

// SetAdmin updates the permission of the user
func (user *User) SetAdmin(x bool) {
	user.Permissions = setKthBit(user.Permissions, 2, x)
}

// HasAddress returns true if a user has credit card info
func (user *User) HasAddress() bool {
	return getKthBit(user.Permissions, 4)
}

// SetAddress updates the permission of the user
func (user *User) SetAddress(x bool) {
	user.Permissions = setKthBit(user.Permissions, 4, x)
}

// HasCreditCardInfo returns true if a user has credit card info
func (user *User) HasCreditCardInfo() bool {
	return getKthBit(user.Permissions, 4)
}

// SetCreditCardInfo updates the permission of the user
func (user *User) SetCreditCardInfo(x bool) {
	user.Permissions = setKthBit(user.Permissions, 4, x)
}

func getKthBit(num int32, k uint32) bool {
	return (uint32(num)>>k)&1 == 1
}

func setKthBit(num int32, k uint32, x bool) int32 {
	if x {
		return int32(uint32(num) ^ ((1<<k)^uint32(num))&(1<<k))
	}
	return int32(uint32(num) ^ ((0<<k)^uint32(num))&(1<<k))
}

// UserDetail is the structure that is stored in the database for a chef's
// or muncher's details
type UserDetail struct {
	Name       string `json:"name" datastore:",noindex"`
	Email      string `json:"email" datastore:",noindex"`
	PhotoURL   string `json:"photo_url" datastore:",noindex"`
	ProviderID string `json:"provider_id" datastore:",noindex"`
}

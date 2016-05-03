package types

// User information in a session.
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProviderID string `json:"provider_id"`
	PhotoURL   string `json:"photo_url"`
	// bit32 | bit31 | bit30 | ...
	// isChef | isVerifiedChef | isAdmin | HasAddress | HasSubMerchantID
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

// HasSubMerchantID returns true if a user has a submerchant id
func (user *User) HasSubMerchantID() bool {
	return getKthBit(user.Permissions, 4)
}

// SetSubMerchantID updates the permission of the user
func (user *User) SetSubMerchantID(x bool) {
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

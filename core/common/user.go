package common

// User is a user. xP
type User struct {
	ID int64 `json:"id" datastore:",noindex"`
	// Firebase's User ID
	AuthID      string `json:"auth_id" datastore:",noindex"`
	FirstName   string `json:"first_name" datastore:",noindex"`
	LastName    string `json:"last_name" datastore:",noindex"`
	Email       string `json:"email" datastore:",noindex"`
	PhotoURL    string `json:"photo_url" datastore:",noindex"`
	Permissions int32  `json:"permissions" datastore:",noindex"`
}

// GetFullName returns the user's name.
func (user *User) GetFullName() string {
	return user.FirstName + " " + user.LastName
}

// GetEmail returns the user's email.
func (user *User) GetEmail() string {
	return user.Email
}

// HasCreditCard returns true if a user has a credit card on file.
func (user *User) HasCreditCard() bool {
	return getKthBit(user.Permissions, 0)
}

// SetHasCreditCard updates the permission of the user.
func (user *User) SetHasCreditCard(x bool) {
	user.Permissions = setKthBit(user.Permissions, 0, x)
}

// IsDriver returns true if a user is a driver.
func (user *User) IsDriver() bool {
	return getKthBit(user.Permissions, 1)
}

// SetIsDriver updates the permission of the user.
func (user *User) SetIsDriver(x bool) {
	user.Permissions = setKthBit(user.Permissions, 1, x)
}

// IsDriverAdmin returns true if a user is a driver admin.
func (user *User) IsDriverAdmin() bool {
	return getKthBit(user.Permissions, 2)
}

// SetIsDriverAdmin updates the permission of the user.
func (user *User) SetIsDriverAdmin(x bool) {
	user.Permissions = setKthBit(user.Permissions, 2, x)
}

// IsUserAdmin returns true if a user is a user admin.
func (user *User) IsUserAdmin() bool {
	return getKthBit(user.Permissions, 2) // TODO: change to 3 at new system stage
}

// SetUserAdmin updates the permission of the user.
func (user *User) SetUserAdmin(x bool) {
	user.Permissions = setKthBit(user.Permissions, 3, x)
}

// IsSystemsAdmin returns true if a user is a systems admin.
func (user *User) IsSystemsAdmin() bool {
	return getKthBit(user.Permissions, 2) // TODO: change to 3 or 4 at new system stage
}

// SetSystemsAdmin updates the permission of the user.
func (user *User) SetSystemsAdmin(x bool) {
	user.Permissions = setKthBit(user.Permissions, 3, x)
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

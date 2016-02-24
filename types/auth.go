package types

import "time"

// AuthToken contains the the user and the user's token string
type AuthToken struct {
	JWTString string
	User      User
	JTI       int32
	ITA       time.Time
	Expire    time.Time
}

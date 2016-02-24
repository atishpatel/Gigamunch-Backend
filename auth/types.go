package auth

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// KindUserSessions is a datastore kind for user information
	KindUserSessions = "UserSessions"
)

// UserSessions is stored in the database to indicate valid user sessions
type UserSessions struct {
	User     types.User `datastore:",noindex"`
	TokenIDs []TokenID  `datastore:",noindex"`
}

// TokenID has unique ids and exp time for tokens
type TokenID struct {
	JTI            int32     `datastore:",noindex"`
	UpdatedToJTI   int32     `datastore:",noindex"`
	Expire         time.Time `datastore:",noindex"`
	UpdateToExpire time.Time `datastore:",noindex"`
}

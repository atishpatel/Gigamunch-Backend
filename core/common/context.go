package common

type contextKey string

const (
	// ContextUserID is the key to the user id value in context.
	ContextUserID = contextKey("user_id")
	// ContextUserEmail is the key to the user email value in context.
	ContextUserEmail = contextKey("user_email")
)

func (c contextKey) String() string {
	return string(c)
}

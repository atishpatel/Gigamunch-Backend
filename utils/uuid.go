package utils

import uuid "github.com/satori/go.uuid"

// GetUUID gnerates a random string of given length
func GetUUID() string {
	return uuid.NewV4().String()
}

// IsValidUUID checks if the passed in text is a valid UUID
func IsValidUUID(text string) bool {
	// TODO(Atish) update len(text) == ?? to validate it
	if len(text) < 32 {
		return false
	}
	return true
}

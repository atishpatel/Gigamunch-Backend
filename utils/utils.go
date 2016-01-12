package utils

import (
	"bytes"

	uuid "github.com/satori/go.uuid"
)

// GetUUID gnerates a random string of given length
func GetUUID() string {
	return uuid.NewV4().String()
}

// IsValidUUID checks if the passed in text is a valid UUID
func IsValidUUID(text string) bool {
	if len(text) < 32 {
		return false
	}
	return true
}

// ConcatenateStrings fast way of concatenating strings
func ConcatenateStrings(args ...string) string {
	var buffer bytes.Buffer
	for _, value := range args {
		buffer.WriteString(value)
	}
	return buffer.String()
}

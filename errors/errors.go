package errors

import "errors"

var (
	// ErrInvalidUUID is returned when a UUID is invalid
	ErrInvalidUUID = errors.New("Invalid UUID")

	// ErrNilParamenter is returned when nil is passed as a parameter
	ErrNilParamenter = errors.New("Used nil as parameter")
)

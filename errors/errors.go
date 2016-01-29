package errors

import (
	"net/http"

	"github.com/docker/distribution/registry/api/errcode"
)

const (
	errGroup = "gigamunch"
)

var (
	// ErrInvalidUUID is returned when a UUID is invalid
	ErrInvalidUUID = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "INVALID_UUID",
		Message:        "The UUID used(%s) is invalid",
		Description:    "The UUID that was used is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	// ErrSessionNotFound is returned when a SessionNotFound is not found
	ErrSessionNotFound = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "SESSION_NOT_FOUND",
		Message:        "The Session(%s) was not found.",
		Description:    "The session was not found.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	// ErrInvalidParameter is returned when a parameter is invalid
	ErrInvalidParameter = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "INVALID_PARAMETER",
		Message:        "%v is not a valid parameter for %s.",
		Description:    "The parameter used is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	// ErrUnauthorizedAccess is returned when a user does not have access to content
	ErrUnauthorizedAccess = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "UNAUTHORIZED_ACCESS",
		Message:        "%s does not have access to %s(%v).",
		Description:    "The user is not allowed to access.",
		HTTPStatusCode: http.StatusUnauthorized,
	})

	// ErrDatastore is returned when there is an error with Datastore
	// Arguments: string(get,put, or delete) string(object name) string(user email) error(error)
	ErrDatastore = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "DATASTORE_ERROR",
		Message:        "Failed to %s %s for %v: %+v",
		Description:    "There was a problem with our database.",
		HTTPStatusCode: http.StatusInternalServerError,
	})
)

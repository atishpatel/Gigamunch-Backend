package errors

import (
	"net/http"

	"github.com/docker/distribution/registry/api/errcode"
)

const (
	errGroup                       = "gigamunch"
	ErrInvalidUUIDValue            = "INVALID_UUID"
	ErrSessionNotFoundValue        = "SESSION_NOT_FOUND"
	ErrExternalDependencyFailValue = "EXTERNAL_DEPENDENCY_FAIL"
	ErrInvalidTokenValue           = "INVALID_TOKEN"
	ErrInvalidParameterValue       = "INVALID_PARAMETER"
	ErrUnauthorizedAccessValue     = "UNAUTHORIZED_ACCESS"
	ErrDatastoreValue              = "DATASTORE_ERROR"
	ErrUnknownValue                = "UNKNOWN_INTERNAL"
)

var (
	// ErrInvalidUUID is returned when a UUID is invalid
	ErrInvalidUUID = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrInvalidUUIDValue,
		Message:        "The UUID used(%s) is invalid",
		Description:    "The UUID that was used is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	// ErrSessionNotFound is returned when a SessionNotFound is not found
	ErrSessionNotFound = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrSessionNotFoundValue,
		Message:        "The Session(%s) was not found.",
		Description:    "The session was not found.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	// ErrInvalidToken is returned when a token is invalid and the user should signout
	ErrInvalidToken = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrInvalidTokenValue,
		Message:        "The token is invalid. err: %+v",
		Description:    "The login information is invalid. Signing out.",
		HTTPStatusCode: 450,
	})

	// ErrExternalDependencyFail is returned when a token is invalid and the user should signout
	ErrExternalDependencyFail = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrExternalDependencyFailValue,
		Message:        "External dependency(%s) failed for reason(%s): %+v",
		Description:    "An external dependency failed.",
		HTTPStatusCode: 450,
	})

	// ErrInvalidParameter is returned when a parameter is invalid
	ErrInvalidParameter = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrInvalidParameterValue,
		Message:        "%v is not a valid parameter for %#v.",
		Description:    "The parameter used is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	// ErrUnauthorizedAccess is returned when a user does not have access to content
	ErrUnauthorizedAccess = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrUnauthorizedAccessValue,
		Message:        "%s does not have access to %s(%v).",
		Description:    "The user is not allowed to access.",
		HTTPStatusCode: http.StatusUnauthorized,
	})

	// ErrDatastore is returned when there is an error with Datastore
	// Arguments: string(get,put, or delete) string(object name) string(user id) error(error)
	ErrDatastore = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrDatastoreValue,
		Message:        "Failed to %s %s for %v: %+v",
		Description:    "There was a problem with our database.",
		HTTPStatusCode: http.StatusInternalServerError,
	})

	// ErrUnknown is returned when there is an unknown internal error
	// Arguments: error(error)
	ErrUnknown = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          ErrUnknownValue,
		Message:        "An unknown internal error occured: %+v",
		Description:    "There was a problem.",
		HTTPStatusCode: http.StatusInternalServerError,
	})
)

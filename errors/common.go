package errors

import "github.com/docker/distribution/registry/api/errcode"

// GetErrorWithStatusCode returns an error casted as a errcode.Error or
// converts a error to an ErrUnknown
func GetErrorWithStatusCode(err error) errcode.Error {
	switch err.(type) {
	case errcode.Error:
		return err.(errcode.Error)
	case errcode.ErrorCode:
		return err.(errcode.ErrorCode).WithArgs()
	default:
		return ErrUnknown.WithArgs(err)
	}
}

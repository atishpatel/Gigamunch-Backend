package errors

import "github.com/docker/distribution/registry/api/errcode"

// GetErrorWithStatusCode returns an error casted as a errcode.Error or
// converts a error to an ErrUnknown
func GetErrorWithStatusCode(err error) errcode.Error {
	_, ok := err.(errcode.Error)
	if ok {
		return err.(errcode.Error)
	}
	return ErrUnknown.WithArgs(err)
}

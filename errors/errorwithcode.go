package errors

import "fmt"

type ErrorWithCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error"`
}

func (err ErrorWithCode) Error() string {
	return fmt.Sprintf("Errorcode %d: %s: %+v", err.Code, err.Message, err.Err)
}

func (err ErrorWithCode) WithError(attachedErr error) ErrorWithCode {
	err.Err = attachedErr
	return err
}

func GetErrorWithCode(err error) ErrorWithCode {
	if errWithCode, ok := err.(ErrorWithCode); ok {
		return errWithCode
	}
	return ErrorWithCode{Code: CodeUnknownError, Message: err.Error()}
}

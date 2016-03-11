package errors

import "fmt"

type ErrorWithCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	err     error  `json:"err"` // err is private because of endpoints
}

func (err ErrorWithCode) Error() string {
	return fmt.Sprintf("Errorcode %d: %s: %+v", err.Code, err.Message, err.err)
}

func (err ErrorWithCode) WithError(attachedErr error) ErrorWithCode {
	err.err = attachedErr
	return err
}

func GetErrorWithCode(err error) ErrorWithCode {
	if errWithCode, ok := err.(ErrorWithCode); ok {
		return errWithCode
	}
	return ErrorWithCode{Code: CodeUnknownError, Message: err.Error()}
}

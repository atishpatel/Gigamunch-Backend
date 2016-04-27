package errors

import "fmt"

type ErrorWithCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (err ErrorWithCode) Error() string {
	return fmt.Sprintf("Errorcode %d:\tMessage %s:\tDetails %s", err.Code, err.Message, err.Details)
}

// WithMessage sets Message. Makes chaining easier.
func (err ErrorWithCode) WithMessage(msg string) ErrorWithCode {
	err.Message = msg
	return err
}

// WithError sets Details to error.Error()
func (err ErrorWithCode) WithError(e error) ErrorWithCode {
	err.Details = e.Error()
	return err
}

// Wrap adds details to the error by making error.Details = additionalDetail + oldDetails
func Wrap(additionalDetail string, e error) ErrorWithCode {
	err := GetErrorWithCode(e)
	err.Details = additionalDetail + err.Details
	return err
}

// GetErrorWithCode casts error as ErrorWithCode or creates an ErrorWithCode
// with CodeUnkownError as the code
func GetErrorWithCode(err error) ErrorWithCode {
	if errWithCode, ok := err.(ErrorWithCode); ok {
		return errWithCode
	}
	return ErrorWithCode{Code: CodeUnknownError, Message: err.Error()}
}

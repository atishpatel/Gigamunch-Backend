package errors

import "fmt"

// ErrorWithCode is used to pass errors with a code, message, and details.
// The code is used to identify the type of error.
// The message is a human readable message.
// The details is the error message itself.
type ErrorWithCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

// GetCode return the code
func (err ErrorWithCode) GetCode() int {
	return err.Code
}

func (err ErrorWithCode) Error() string {
	return fmt.Sprintf("Errorcode: %d| Message: %s| Detail: %s", err.Code, err.Message, err.Detail)
}

// WithMessage sets Message. Makes chaining easier.
func (err ErrorWithCode) WithMessage(msg string) ErrorWithCode {
	err.Message = msg
	return err
}

// WithError sets Details to error.Error()
func (err ErrorWithCode) WithError(e error) ErrorWithCode {
	return err.Wrap(e.Error())
}

//Wrap annotates the error by prepending string to error's details
func (err ErrorWithCode) Wrap(additionalDetail string) ErrorWithCode {
	if err.Detail == "" {
		err.Detail = additionalDetail
	} else {
		err.Detail = additionalDetail + ": " + err.Detail
	}
	return err
}

//Wrapf annotates the error by fromating string and prepending it to error's details
func (err ErrorWithCode) Wrapf(format string, args ...interface{}) ErrorWithCode {
	return err.Wrap(fmt.Sprintf(format, args))
}

// Equal returns if the two errors have the same Code
func (err ErrorWithCode) Equal(e error) bool {
	if e == nil {
		return false
	}
	if err.Code != GetErrorWithCode(e).Code {
		return false
	}
	return true
}

// Wrap annotates the error by prepending string to error's details
func Wrap(additionalDetail string, e error) ErrorWithCode {
	err := GetErrorWithCode(e)
	return err.Wrap(additionalDetail)
}

// GetErrorWithCode casts error as ErrorWithCode or creates an ErrorWithCode
// with CodeUnkownError as the code
func GetErrorWithCode(err error) ErrorWithCode {
	if errWithCode, ok := err.(ErrorWithCode); ok {
		return errWithCode
	}
	return ErrorWithCode{Code: CodeUnknownError, Message: err.Error()}
}

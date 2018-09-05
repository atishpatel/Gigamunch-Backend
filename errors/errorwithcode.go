package errors

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"
)

const (
	CodeSuccess                    = int32(shared.Code_Success)
	CodeBadRequest                 = 400
	CodeInvalidParameter           = 400
	CodeInvalidPromoCode           = 402
	CodeNotFound                   = 404
	CodePermissionDenied           = int32(shared.Code_PermissionDenied)
	CodeUnauthenticated            = int32(shared.Code_Unauthenticated)
	CodeUnauthorizedAccess         = 401
	CodeSignOut                    = 452
	CodeForbidden                  = 403
	CodeNotEnoughServingsLeft      = 453
	CodeDeliveryMethodNotAvaliable = 454
	CodeBTInvalidPaymentMethod     = 455
	CodeBTFailedToProcess          = 456
	CodeInternalServerErr          = 500
	CodeUnknownError               = 666
)

var (
	NoError               = ErrorWithCode{Code: CodeSuccess, Message: "Success."}
	Success               = NoError
	BadRequestError       = ErrorWithCode{Code: CodeBadRequest, Message: "Bad Request."}
	InternalServerError   = ErrorWithCode{Code: CodeInternalServerErr, Message: "Internal server error."}
	SignOutError          = ErrorWithCode{Code: CodeSignOut, Message: "Bad Request."}
	PermissionDeniedError = ErrorWithCode{Code: CodePermissionDenied, Message: "Permission Denied."}
	UnauthenticatedError  = ErrorWithCode{Code: CodeUnauthenticated, Message: "Unauthenticated user."}
	UnknownError          = ErrorWithCode{Code: CodeUnknownError, Message: "An unknown error occured."}
)

// ErrorWithCode is used to pass errors with a code, message, and details.
// The code is used to identify the type of error.
// The message is a human readable message.
// The details is the error message itself.
type ErrorWithCode struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

// IsNil checks if error is code success therefore not an error.
func (ewc ErrorWithCode) IsNil() bool {
	return ewc.Code == CodeSuccess
}

// GetCode return the code
func (ewc ErrorWithCode) GetCode() int32 {
	return ewc.Code
}

// GetHTTPCode return the HTTP code.
func (ewc ErrorWithCode) GetHTTPCode() int {

	return int(ewc.Code)
}

func (ewc ErrorWithCode) Error() string {
	return fmt.Sprintf("Code: %d |\n Message: %s |\n Detail: %s", ewc.Code, ewc.Message, ewc.Detail)
}

// WithMessage sets Message. Makes chaining easier.
func (ewc ErrorWithCode) WithMessage(msg string) ErrorWithCode {
	ewc.Message = msg
	return ewc
}

// WithError sets Details to error.Error()
func (ewc ErrorWithCode) WithError(e error) ErrorWithCode {
	if e == nil {
		return ewc
	}
	return ewc.Wrap(e.Error())
}

// Annotate annotates the error by prepending string to error's details.
func (ewc ErrorWithCode) Annotate(additionalDetail string) ErrorWithCode {
	return ewc.Wrap(additionalDetail)
}

//Wrap annotates the error by prepending string to error's details
func (ewc ErrorWithCode) Wrap(additionalDetail string) ErrorWithCode {
	if ewc.Detail == "" {
		ewc.Detail = additionalDetail
	} else {
		ewc.Detail = additionalDetail + ": " + ewc.Detail
	}
	return ewc
}

// Annotatef annotates the error by prepending string to error's details.
func (ewc ErrorWithCode) Annotatef(format string, args ...interface{}) ErrorWithCode {
	return ewc.Wrap(fmt.Sprintf(format, args...))
}

//Wrapf annotates the error by fromating string and prepending it to error's details
func (ewc ErrorWithCode) Wrapf(format string, args ...interface{}) ErrorWithCode {
	return ewc.Wrap(fmt.Sprintf(format, args...))
}

// Equal returns if the two errors have the same Code
func (ewc ErrorWithCode) Equal(e error) bool {
	if e == nil {
		return false
	}
	if ewc.Code != GetErrorWithCode(e).Code {
		return false
	}
	return true
}

// SharedError returns ErrorWithCode as a shared.Error.
func (ewc ErrorWithCode) SharedError() *shared.Error {
	return &shared.Error{
		Code:    shared.Code(ewc.Code),
		Message: ewc.Message,
		Detail:  ewc.Detail,
	}
}

// GetError returns ErrorWithCode as a shared.Error.
func (ewc ErrorWithCode) GetError() *shared.Error {
	return ewc.SharedError()
}

// Annotate annotates the error by prepending string to error's details.
func Annotate(e error, detail string) ErrorWithCode {
	ewc := GetErrorWithCode(e)
	return ewc.Annotate(detail)
}

// Wrap annotates the error by prepending string to error's details
func Wrap(additionalDetail string, e error) ErrorWithCode {
	ewc := GetErrorWithCode(e)
	return ewc.Wrap(additionalDetail)
}

// GetErrorWithCode casts error as ErrorWithCode or creates an ErrorWithCode
// with CodeUnkownError as the code
func GetErrorWithCode(ewc interface{}) ErrorWithCode {
	if errWithCode, ok := ewc.(ErrorWithCode); ok {
		return errWithCode
	}
	if sharedError, ok := ewc.(*shared.Error); ok {
		return ErrorWithCode{
			Code:    int32(sharedError.Code),
			Message: sharedError.Message,
			Detail:  sharedError.Detail,
		}
	}
	if sharedError, ok := ewc.(shared.Error); ok {
		return ErrorWithCode{
			Code:    int32(sharedError.Code),
			Message: sharedError.Message,
			Detail:  sharedError.Detail,
		}
	}
	if err, ok := ewc.(error); ok {
		return UnknownError.WithMessage(err.Error())
	}
	return UnknownError
}

// GetSharedError converts ErrorWithCode into a shared.Error.
func GetSharedError(ewc error) *shared.Error {
	if ewc, ok := ewc.(ErrorWithCode); ok {
		return ewc.SharedError()
	}
	return UnknownError.Annotate(ewc.Error()).SharedError()
}

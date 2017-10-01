package errors

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"
)

const (
	// OK is returned on success.
	OK shared.Code = shared.Code_Success

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled shared.Code = 1

	// Unknown error.  An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space.  Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown shared.Code = 2

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument shared.Code = 3

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	DeadlineExceeded shared.Code = 4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound shared.Code = 5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists shared.Code = 6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors).  It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied shared.Code = 7

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated shared.Code = 16

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted shared.Code = 8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed.  E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	FailedPrecondition shared.Code = 9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted shared.Code = 10

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange.  We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	OutOfRange shared.Code = 11

	// Internal errors.  Means some invariants expected by underlying
	// system has been broken.  If you see one of these errors,
	// something is very broken.
	Internal shared.Code = 13

	CodeSuccess                    = int32(shared.Code_Success)
	CodeBadRequest                 = 400
	CodeInvalidParameter           = 400
	CodeInvalidPromoCode           = 402
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
	return fmt.Sprintf("Errorcode: %d| Message: %s| Detail: %s", ewc.Code, ewc.Message, ewc.Detail)
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

package apperrors

import (
	"fmt"
	"runtime/debug"

	"github.com/google/uuid"
)

type Code int

const (
	InternalCode Code = iota + 1
	BadRequestCode
	NotFoundCode
	AlreadyExistsCode
	UnauthorizedCode
	ForbiddenCode
)

var _ error = (*Error)(nil)

type Error struct {
	Code       Code
	StackTrace string
	IncidentId string

	innerError error
	hiderError bool
	message    string
}

func (e *Error) Error() string {
	return e.error(false)
}

func (e *Error) SafeError() string {
	return e.error(true)
}

func (e *Error) Unwrap() error {
	return e.innerError
}

func (e *Error) error(safe bool) string {
	switch {
	case e.innerError == nil:
		return e.message
	case safe && e.hiderError:
		return e.message
	case e.message == "":
		return e.innerError.Error()
	default:
		return fmt.Sprintf("%s: %s", e.message, e.innerError.Error())
	}
}

func WrapInternal(err error) *Error {
	if err == nil {
		return nil
	}

	return Internal(err)
}

func Internal(err error) *Error {
	appErr := InternalWithoutStackTrace(err)
	appErr.StackTrace = string(debug.Stack())
	return appErr
}

func InternalWithoutStackTrace(err error) *Error {
	appErr := newHiddenError(err, InternalCode, "internal error")
	appErr.IncidentId = uuid.New().String()
	return appErr
}

func BadRequest(err error) *Error {
	return newWrappedError(err, BadRequestCode)
}

func BadRequestHidden(err error, message string) *Error {
	return newHiddenError(err, BadRequestCode, message)
}

func NotFound(subject, key string, value any) *Error {
	return newError(NotFoundCode, fmt.Sprintf("%s %s: %v not found", subject, key, value))
}

func AlreadyExists(subject, key string, value any) *Error {
	return newError(AlreadyExistsCode, fmt.Sprintf("%s %s: %v already extists", subject, key, value))
}

func Unauthorized(message string) *Error {
	return newError(UnauthorizedCode, message)
}

func UnauthorizedHidden(err error, message string) *Error {
	return newHiddenError(err, UnauthorizedCode, message)
}

func Forbidden(message string) *Error {
	return newError(ForbiddenCode, message)
}

func newError(code Code, message string) *Error {
	return &Error{
		Code:    code,
		message: message,
	}
}

func newWrappedError(err error, code Code) *Error {
	return &Error{
		Code:       code,
		innerError: err,
	}
}

func newHiddenError(err error, code Code, message string) *Error {
	return &Error{
		Code:       code,
		message:    message,
		innerError: err,
		hiderError: true,
	}
}
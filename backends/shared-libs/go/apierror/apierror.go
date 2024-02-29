package apierror

import (
	"fmt"

	"github.com/pkg/errors"
)

type APIError struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Desc)
}

func FromError(err error) APIError {
	return APIError{Code: CodeInternalServerError, Desc: err.Error()}
}

func New(code string) APIError {
	return APIError{
		Code: code,
	}
}

func WithDesc(code string, desc string) APIError {
	return APIError{
		Code: code,
		Desc: desc,
	}
}

// NOTE(raka) still in development
func GetCauseError(err error) (APIError, bool) {
	if f, ok := errors.Cause(err).(APIError); ok {
		return f, true
	}
	return APIError{}, false
}

const (
	CodeInternalServerError = "InternalServerError"
	CodeEntityNotFound      = "EntityNotFound"
	CodeInvalidRequest      = "InvalidRequest"
	CodeUnauthorized        = "Unauthorized"
	CodeForbidden           = "Forbidden"
)

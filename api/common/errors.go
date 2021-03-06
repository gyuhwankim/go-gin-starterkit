package common

import (
	"github.com/pkg/errors"
)

var (
	// ErrParsingFailed was occurred when request invalid syntax.
	ErrParsingFailed = errors.New("parsing is failed by invalid syntax")

	// ErrEntityNotFound was occurred when can't found entity.
	ErrEntityNotFound = errors.New("Entity was not found")

	// ErrAlreadyExistsEntity was occurred
	// when insert data that contain duplicated key.
	ErrAlreadyExistsEntity = errors.New("Already exists entity")

	// ErrInvalidUUID was occurred when uuid format was invalid.
	ErrInvalidUUID = errors.New("UUID was invalid")

	// ErrInvalidRequestPayload is occurred when payload is invalid.
	ErrInvalidRequestPayload = errors.New("Request payload is invalid")
)

// ErrorResponse is app response.
type ErrorResponse struct {
	Errors []APIError `json:"errors"`
}

// AddError add new error at ErrorResponse.
func (errResponse ErrorResponse) AddError(err error) ErrorResponse {
	errResponse.Errors = append(errResponse.Errors, APIError{
		Message: err.Error(),
	})

	return errResponse
}

// APIError is http error object.
type APIError struct {
	Message string `json:"message"`
}

// NewErrResp is return new error response.
func NewErrResp(err error) ErrorResponse {
	return ErrorResponse{
		Errors: []APIError{
			APIError{
				Message: err.Error(),
			},
		},
	}
}

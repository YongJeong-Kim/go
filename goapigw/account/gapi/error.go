package gapi

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Inner      error
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Inner.Error()
}

func (e *APIError) Unwrap() error {
	return e.Inner
}

func (e *APIError) ErrorWrap(err error, msg string) error {
	var ae *APIError
	errors.As(err, &ae)
	return ae
}

func errorResponse(err error) gin.H {
	var ae *APIError
	errors.As(err, &ae)

	return gin.H{
		"inner":       ae.Inner.Error(),
		"status_code": ae.StatusCode,
		"message":     ae.Message,
	}
}

func gErrorResponse(err error) map[string]any {
	var ae *APIError
	errors.As(err, &ae)
	return map[string]any{
		"inner":       ae.Inner.Error(),
		"status_code": ae.StatusCode,
		"message":     ae.Message,
	}
}

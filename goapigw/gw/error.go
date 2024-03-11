package main

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type CustomError struct {
	Inner      error
	StatusCode int
	Message    string
}

func (e *CustomError) Error() string {
	return e.Inner.Error()
}

func (e *CustomError) Unwrap() error {
	return e.Inner
}

func (e *CustomError) ErrorWrap(err error, msg string) error {
	var ce *CustomError
	errors.As(err, &ce)
	return ce
}

func errorResponse(err error) gin.H {
	var ce *CustomError
	errors.As(err, &ce)

	return gin.H{
		"inner":       ce.Inner,
		"status_code": ce.StatusCode,
		"message":     ce.Message,
	}
}

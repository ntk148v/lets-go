# Handling errors gracefully

[Source](https://hackernoon.com/golang-handling-errors-gracefully-8e27f1db729f)

## What is an error in go?

```golang
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
```

This definition tells us that all it takes to create an error is a simple string.

## Error flow

For the sake of simplicity and don’t repeat yourself principle, is desirable to take action on an error once at a single place.

Go built-in error doesn't provide a stack trace, might be [Go 2 feature](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling-overview.md) but not this time. Function should typically include relevant information about their errors, like `os.Open` returning the name of the file being opened. Returning the error unmodified produces a failure without any information about the sequence of operations that led to the error.

**github.com/pkg/errors** to the rescue.

Question: We just handled the error at the top layer? Perfect? Nah!

## A solution

Goals:

- Provide a good error stack trace.
- Log the error.
- Provide a contextual error information to the user when necessary.

Create an error type:

```golang
package error

import (
	"fmt"
	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	// NoType error
	NoType ErrorType = iota
	// BadRequest error
	BadRequest
	// NotFound error
	NotFound
)

type customError struct {
	errorType     ErrorType
	originalError error
	context       errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// New create a new customError
func (errorType ErrorType) New(msg string) error {
	return customError{errorType: errorType, originalError: errors.New(msg)}
}

// Newf creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrapf creates a new wrapped error with format message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

// Error returns the message of a customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// New creates a no type error
func New(msg string) error {
	return customError{errorType: NoType, originalError: errors.New(msg)}
}

// Newf creates a no type rerror with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf an error with format sring
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}
	return customError{errorType: NoType, originalError: wrappedError}
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(customError); ok {
		return customError{errorType: customErr.errorType, originalError: customErr.originalError, context: context}
	}
	return customError{errorType: NoType, originalError: err, context: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok || customErr.context != emptyContext {
		return map[string]string{"field": customErr.context.Field, "message": customErr.context.Message}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return NoType
}
```

# Handling errors gracefully

Source:

- <https://hackernoon.com/golang-handling-errors-gracefully-8e27f1db729f>
- <https://www.sobyte.net/post/2022-06/go-error/>

- [Handling errors gracefully](#handling-errors-gracefully)
	- [1. What is an error in go?](#1-what-is-an-error-in-go)
	- [2. Error flow](#2-error-flow)
	- [3. A solution to handle errors gracefully?](#3-a-solution-to-handle-errors-gracefully)
	- [4. Techniques and principles of error handling](#4-techniques-and-principles-of-error-handling)
		- [4.1. Using wrappers to avoid repetitive error judgments](#41-using-wrappers-to-avoid-repetitive-error-judgments)
		- [4.2. Error handling before Golang 1.13](#42-error-handling-before-golang-113)
		- [4.3. Error Handling in Golang 1.13](#43-error-handling-in-golang-113)

## 1. What is an error in go?

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
    Error() string
}
```

This definition tells us that all it takes to create an error is a simple string.

## 2. Error flow

For the sake of simplicity and don’t repeat yourself principle, is desirable to take action on an error once at a single place.

Go built-in error doesn't provide a stack trace, might be [Go 2 feature](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling-overview.md) but not this time. Function should typically include relevant information about their errors, like `os.Open` returning the name of the file being opened. Returning the error unmodified produces a failure without any information about the sequence of operations that led to the error.

**github.com/pkg/errors** to the rescue.

Question: We just handled the error at the top layer? Perfect? Nah!

## 3. A solution to handle errors gracefully?

Goals:

- Provide a good error stack trace.
- Log the error.
- Provide a contextual error information to the user when necessary.

Create an error type:

```go
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

## 4. Techniques and principles of error handling

### 4.1. Using wrappers to avoid repetitive error judgments

- Golang takes errors as return values, so you have to deal with them.

```go
_, err = fd.Write(p0[a:b])
if err != nil {
    return err
}
_, err = fd.Write(p1[c:d])
if err != nil {
    return err
}
_, err = fd.Write(p2[e:f])
if err != nil {
    return err
}
// and so on
// have to write 9 error judgments
```

- Wrap `io.Writer` in another layer:

```go
type errWriter struct {
    w   io.Writer
    err error
}

func (ew *errWriter) write(buf []byte) {
    if ew.err != nil {
        return
    }
    _, ew.err = ew.w.Write(buf)
}

ew := &errWriter{w: fd}
ew.write(p0[a:b])
ew.write(p1[c:d])
ew.write(p2[e:f])
// and so on
if ew.err != nil {
    return ew.err
}
// Drawback: Have no way to know which line the error is called on
```

### 4.2. Error handling before Golang 1.13

- Checking for errors: There are times when we need to handle the error differently depending on the type of error, for example, if the error is retryable (connection not connected).

```go
// Methed 1: Compare the error with a known value
var ErrNotFound = errors.New("not found")

if err == ErrNotFound {
    // something wasn't found
}

// Method 2: Determine the specific type of error
type NotFoundError struct {
    Name string
}

func (e *NotFoundError) Error() string { return e.Name + ": not found" }

if e, ok := err.(*NotFoundError); ok {
    // e.Name wasn't found
}
```

- Adding information

```go
// Construct a new error using the information from the previous error
// Using fmt.Errorf only keeps the text of the previous error
if err != nil {
    return fmt.Errorf("decompress %v: %v", name, err)
}

// If we want to keep all the information from the previous error,
// we can use the following.
type QueryError struct {
    Query string
    Err   error
}

if e, ok := err.(*QueryError); ok && e.Err == ErrPermission {
    // query failed because of a permission problem
}
```

### 4.3. Error Handling in Golang 1.13

- The underlying error can be returned by implementing the `Unwrap()` method. If `e1.Unwrap()` returns e2, we can say that e1 contains e2.
- Using Is and As to check for errors.

```go
// Similar to:
//   if err == ErrNotFound { … }
if errors.Is(err, ErrNotFound) {
    // something wasn't found
}

// Similar to:
//   if e, ok := err.(*QueryError); ok { … }
var e *QueryError
if errors.As(err, &e) {
    // err is a *QueryError, and e is set to the error's value
}
```

```go
type ErrorA struct {
    Msg string
}

func (e *ErrorA) Error() string {
    return e.Msg
}

type ErrorB struct {
    Msg string
    Err *ErrorA
}

func (e *ErrorB) Error() string {
    return e.Msg + e.Err.Msg
}

func (e *ErrorB) Unwrap() error {
    return e.Err
}

func main() {
    a := &ErrorA{"error a"}

    b := &ErrorB{"error b", a}

    if errors.Is(b, a) {
        log.Println("error b is a")
    }

    var tmpa *ErrorA
    if errors.As(b, &tmpa) {
        log.Println("error b as ErrorA")
    }
}
```

- Wrapping errors with %w

```go
type ErrorA struct {
    Msg string
}

func (e *ErrorA) Error() string {
    return e.Msg
}

func main() {
    a := &ErrorA{"error a"}

    b := fmt.Errorf("new error: %w", a)

    if errors.Is(b, a) {
        fmt.Println("error b is a")
    }

    var tmpa *ErrorA
    if errors.As(b, &tmpa) {
        fmt.Println("error b as ErrorA")
    }
}
```

- If you don't want to expose implementation details, don't wrap the error. Because exposing an error with details means that the caller is coupled to our code. This also violates the principle of abstraction.

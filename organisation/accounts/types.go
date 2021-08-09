package accounts

import (
	"errors"
	"time"
)

const (
	AccountsPath        = "/v1/organisation/accounts"
	DefaultTimeOutValue = time.Duration(10000 * time.Millisecond)
)

// Error Standard Types
var (
	// add api errror type per http status: NotFoundError; BadRequestError
	ApiHttpErrorType = errors.New("API Error")
	// HttpBadRequestErr     = errors.New("Http Bad Request error")
	// HttpNotFoundErr       = errors.New("Http Not Found error")
	// HttpInternalServerErr = errors.New("Http Internal Server error")
	UnmarshallingError    = errors.New("UnmarshallingError")
	BaseUrlParsingError   = errors.New("BaseUrlParsingError")
	PathParsingError      = errors.New("PathParsingError")
	BuildingRequestError  = errors.New("BuildingRequestError")
	ExecutingRequestError = errors.New("ExecutingRequestError")
	ResponseReadError     = errors.New("ResponseReadError")
	TimeoutError          = errors.New("TimeoutError")
)

// type ClientError struct {
// 	url          string
// 	httpCode     int
// 	errorMessage string
// 	err          error
// }

// Errors and package APIs
// A package which returns errors (and most do) should describe what properties of those errors programmers may rely on. A well-designed package will also avoid returning errors with properties that should not be relied upon.

// The simplest specification is to say that operations either succeed or fail, returning a nil or non-nil error value respectively. In many cases, no further information is needed.

// If we wish a function to return an identifiable error condition, such as "item not found," we might return an error wrapping a sentinel.

// https://blog.golang.org/go1.13-errors
// fmt.Errorf("%s | Path: %q returned %d with message %s", TimeoutError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
// func (e *ClientError) Unwrap() error        { return e.err }
// func (e *ClientError) URL() string          { return e.url }
// func (e *ClientError) Err() error           { return e.err }

// func (e *ClientError) Error() string {
// 	return fmt.Sprintf("%s | Path: %q returned %d with message %s", e.err, e.url, e.httpCode, e.errorMessage)
// }

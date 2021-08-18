package accounts

import (
	"errors"
	"fmt"
	"time"
)

const (
	DefaultTimeOutValue = time.Duration(10 * time.Second)
)

// Error Standard Types
var (
	ApiHttpErrorType      = errors.New("Error message returned by the API")
	UnmarshallingError    = errors.New("UnmarshallingError")
	ExecutingRequestError = errors.New("ExecutingRequestError")
	clientCreationError   = errors.New("Unable to create the client")
)

type apiErrorMessage struct {
	ErrorMessage string `json:"error_message"`
	Status       int
}

func isHttpCodeOK(httpCode int) bool {
	return httpCode >= 200 && httpCode < 300
}

func (e *apiErrorMessage) Error() error {

	if !isHttpCodeOK(e.Status) {
		return fmt.Errorf("%w | %d | %s", ApiHttpErrorType, e.Status, e.ErrorMessage)
	}

	return nil
}

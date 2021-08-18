package accounts

import (
	"errors"
	"fmt"
	"time"
)

// Default timeoutvalue defined within the client.
// You can override this value using the WithTimeout function
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

// apiErrorMessage contains the error message returned by the Form3 API and it's http code. This is used internally.
type apiErrorMessage struct {

	// ErrorMessage is the explanatory field added when API returns an error.
	ErrorMessage string `json:"error_message"`

	// Status is a field added within the client. It concerns the http status code from the client call and
	// is meant to help you track down any bug
	Status int
}

func isHttpCodeOK(httpCode int) bool {
	return httpCode >= 200 && httpCode < 300
}

// Error returns an error if this object has a Status not between 200 and 300
func (e *apiErrorMessage) Error() error {

	if !isHttpCodeOK(e.Status) {
		return fmt.Errorf("%w | %d | %s", ApiHttpErrorType, e.Status, e.ErrorMessage)
	}

	return nil
}

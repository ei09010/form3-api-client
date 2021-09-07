package accounts

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Default timeoutvalue defined within the client.
// You can override this value using the WithTimeout function
const (
	DefaultTimeOutValue = time.Duration(10 * time.Second)
)

// Error Standard Types
var (
	ApiHttpErrorType     = errors.New("Error message returned by the API")
	BuildingRequestError = errors.New("Error while building the request")
	ClientCreationError  = errors.New("Unable to create the client")
)

// apiCommonResult contains the error message returned by the Form3 API and it's http code. This is used internally.
type apiCommonResult struct {

	// ErrorMessage is the explanatory field added when API returns an error.
	ErrorMessage string `json:"error_message"`

	// Status is a field mapped from the http response status code. It concerns the http status code from the client call and
	// is meant to help you track down any bug
	Status int
}

func isHttpCodeOK(httpCode int) bool {
	return httpCode >= http.StatusOK && httpCode < http.StatusBadRequest
}

// Error returns an error if this object has a Status not between 200 and 300
func (e *apiCommonResult) Error() error {

	if !isHttpCodeOK(e.Status) {
		return fmt.Errorf("%w | %d | %s", ApiHttpErrorType, e.Status, e.ErrorMessage)
	}

	return nil
}

func addHeaders(customReq *http.Request) {
	customReq.Header.Set("Content-Type", "application/json")
	customReq.Header.Set("User-Agent", "form3-go rest-client/0.1 go1.17")
}

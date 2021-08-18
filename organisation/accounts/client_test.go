package accounts

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientSuccessCases(t *testing.T) {

	// Arrange
	errorCases := map[string]struct {
		baseURL              string
		expectedScheme       string
		expectedPath         string
		expectedHost         string
		argumentTimeoutValue time.Duration
		timeoutValue         time.Duration
	}{
		"Valid Url": {
			baseURL:              "http://localhost:8080",
			expectedScheme:       "http",
			expectedPath:         "",
			expectedHost:         "localhost:8080",
			argumentTimeoutValue: time.Duration(5 * time.Millisecond),
			timeoutValue:         time.Duration(5 * time.Millisecond),
		},
		"Valid Url and Zero Timeout Value": {
			baseURL:              "http://localhost:8080",
			expectedScheme:       "http",
			expectedPath:         "",
			expectedHost:         "localhost:8080",
			argumentTimeoutValue: time.Duration(0),
			timeoutValue:         DefaultTimeOutValue,
		},
		"Valid Url and Negative Timeout Value": {
			baseURL:              "http://localhost:8080",
			expectedScheme:       "http",
			expectedPath:         "",
			expectedHost:         "localhost:8080",
			argumentTimeoutValue: time.Duration(-1 * time.Millisecond),
			timeoutValue:         DefaultTimeOutValue,
		},
	}

	for _, tt := range errorCases {

		accountClient, err := NewClient(WithBaseURL(tt.baseURL), WithTimeout(time.Duration(tt.argumentTimeoutValue)))

		if err != nil {
			t.Errorf("Returned reponse: got %v want %v",
				accountClient, nil)
		}

		// Assert

		if accountClient.baseURL.Path != tt.expectedPath {
			t.Errorf("client returned path: got %s want %s",
				accountClient.baseURL.Path, tt.expectedPath)
		}

		if accountClient.baseURL.Host != tt.expectedHost {
			t.Errorf("client returned host: got %s want %s",
				accountClient.baseURL.Host, tt.expectedHost)
		}

		if accountClient.baseURL.Scheme != tt.expectedScheme {
			t.Errorf("client returned scheme: got %s want %s",
				accountClient.baseURL.Scheme, tt.expectedScheme)
		}

		if accountClient.timeout != tt.timeoutValue {
			t.Errorf("client returned timeout: got %s want %s",
				accountClient.timeout, tt.timeoutValue)
		}
	}

}

func TestClientErrorCases(t *testing.T) {

	// Arrange
	errorCases := map[string]struct {
		argumentUrl          string
		expectedErrorMessage string
		expectedErrorType    error
		baseURL              string
		expectedHttpStatus   int
		timeoutValue         time.Duration
	}{
		"Empty Base Url": {
			expectedErrorMessage: `parse "": empty url`,
			expectedErrorType:    clientCreationError,
			baseURL:              "",
			expectedHttpStatus:   http.StatusBadRequest,
			timeoutValue:         time.Duration(1 * time.Millisecond),
		},
		"Invalid Base Url": {
			expectedErrorMessage: `parse "wrongURL": invalid URI for request`,
			expectedErrorType:    clientCreationError,
			baseURL:              "wrongURL",
			expectedHttpStatus:   http.StatusBadRequest,
			timeoutValue:         time.Duration(1 * time.Millisecond),
		},
		"Invalid baseURL and Invalid Timeout Value": {
			expectedErrorMessage: `parse "wrongURL": invalid URI for request`,
			expectedErrorType:    clientCreationError,
			baseURL:              "wrongURL",
			expectedHttpStatus:   http.StatusBadRequest,
			timeoutValue:         time.Duration(-1 * time.Millisecond),
		},
	}

	for _, tt := range errorCases {

		// Act

		accountClient, err := NewClient(WithBaseURL(tt.baseURL), WithTimeout(time.Duration(tt.timeoutValue)))
		// Assert

		if accountClient != nil {
			t.Errorf("Returned reponse: got %v want %v",
				accountClient, nil)
		}

		assertClientError(err, tt.expectedErrorMessage, t, tt.expectedErrorType, tt.expectedHttpStatus)
	}

}

// newTestServer creates a multiplex server to handle API endpoints
func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)
	return server
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func assertClientError(err error, expectedErrorMessage string, t *testing.T, expectedErrorType error, expectedHttpStatus int) {
	if err != nil {

		if errors.Is(err, expectedErrorType) {

			expectedErrorFinalMessage := fmt.Errorf("%w | %d | %s", expectedErrorType, expectedHttpStatus, expectedErrorMessage)

			if err.Error() != expectedErrorFinalMessage.Error() {
				t.Errorf("Returned error message: got %s want %s",
					err.Error(), expectedErrorFinalMessage.Error())
			}
		} else {
			t.Errorf("Returned error type: got %s want %s",
				err.Error(), expectedErrorType.Error())
		}
	}
}

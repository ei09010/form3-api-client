package accounts_test

import (
	"ei09010/form3-api-client/organisation/accounts"
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
		baseUrl               string
		expectedScheme        string
		expectedPath          string
		expectedHost          string
		argumentTimeoutValue  time.Duration
		expectedTimeoutClient time.Duration
	}{
		"Valid Url": {
			baseUrl:               "http://localhost:8080",
			expectedScheme:        "http",
			expectedPath:          "/v1/organisation/accounts",
			expectedHost:          "localhost:8080",
			expectedTimeoutClient: time.Duration(5 * time.Millisecond),
		},
		"Valid Url and Zero Timeout Value": {
			baseUrl:               "http://localhost:8080",
			expectedScheme:        "http",
			expectedPath:          "/v1/organisation/accounts",
			expectedHost:          "localhost:8080",
			argumentTimeoutValue:  time.Duration(0),
			expectedTimeoutClient: accounts.DefaultTimeOutValue,
		},
		"Valid Url and Invalid Timeout Value": {
			baseUrl:               "http://localhost:8080",
			expectedScheme:        "http",
			expectedPath:          "/v1/organisation/accounts",
			expectedHost:          "localhost:8080",
			argumentTimeoutValue:  time.Duration(-1),
			expectedTimeoutClient: accounts.DefaultTimeOutValue,
		},
	}

	for _, tt := range errorCases {

		accountClient, err := accounts.NewClient(tt.baseUrl, time.Duration(tt.expectedTimeoutClient))

		if err != nil {
			t.Errorf("Returned reponse: got %v want %v",
				accountClient, nil)
		}

		// Assert

		if accountClient.BaseURL.Path != tt.expectedPath {
			t.Errorf("client returned path: got %s want %s",
				accountClient.BaseURL.Path, tt.expectedPath)
		}

		if accountClient.BaseURL.Host != tt.expectedHost {
			t.Errorf("client returned host: got %s want %s",
				accountClient.BaseURL.Host, tt.expectedHost)
		}

		if accountClient.BaseURL.Scheme != tt.expectedScheme {
			t.Errorf("client returned scheme: got %s want %s",
				accountClient.BaseURL.Scheme, tt.expectedScheme)
		}

		if accountClient.Timeout != tt.expectedTimeoutClient {
			t.Errorf("client returned timeout: got %s want %s",
				accountClient.Timeout, tt.expectedTimeoutClient)
		}
	}

}

func TestClientErrorCases(t *testing.T) {

	// Arrange
	errorCases := map[string]struct {
		argumentUrl          string
		expectedErrorMessage string
		expectedErrorType    error
		baseUrl              string
		expectedHttpStatus   int
		timeoutValue         time.Duration
	}{
		"Empty Base Url": {
			expectedErrorMessage: `parse "": empty url`,
			expectedErrorType:    accounts.BaseUrlParsingError,
			baseUrl:              "",
			expectedHttpStatus:   http.StatusBadRequest,
			timeoutValue:         time.Duration(1 * time.Millisecond),
		},
		"Invalid Base Url": {
			expectedErrorMessage: `parse "wrongURL": invalid URI for request`,
			expectedErrorType:    accounts.BaseUrlParsingError,
			baseUrl:              "wrongURL",
			expectedHttpStatus:   http.StatusBadRequest,
			timeoutValue:         time.Duration(1 * time.Millisecond),
		},
		"Invalid BaseUrl and Invalid Timeout Value": {
			expectedErrorMessage: `parse "wrongURL": invalid URI for request`,
			expectedErrorType:    accounts.BaseUrlParsingError,
			baseUrl:              "wrongURL",
			expectedHttpStatus:   http.StatusBadRequest,
			timeoutValue:         time.Duration(-1 * time.Millisecond),
		},
	}

	for _, tt := range errorCases {

		// Act

		accountClient, err := accounts.NewClient(tt.baseUrl, tt.timeoutValue)

		// Assert

		if accountClient != nil {
			t.Errorf("Returned reponse: got %v want %v",
				accountClient, nil)
		}

		assertClientError(err, tt.expectedErrorMessage, t, tt.baseUrl, tt.expectedHttpStatus, tt.expectedErrorType)
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

func assertClientError(err error, expectedErrorMessage string, t *testing.T, expectedCorrectRequest string, expectedhttpStatus int, expectedErrorType error) {
	if err != nil {

		if errors.Is(err, expectedErrorType) {

			expectedErrorFinalMessage := fmt.Errorf("%w | Path: %s returned %d with message %s", expectedErrorType, expectedCorrectRequest, expectedhttpStatus, expectedErrorMessage)

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

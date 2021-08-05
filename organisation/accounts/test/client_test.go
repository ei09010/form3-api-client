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

var (
	validUrl = "http://localhost:8080"
)

func TestNewClient_validUrl_returnsValidClient(t *testing.T) {

	// Arrange

	expectedScheme := "http"
	expectedPath := "/v1/organisation/accounts"
	expectedHost := "localhost:8080"
	expectedTimeoutClient := time.Duration(5 * time.Millisecond)

	// Act

	accountClient, err := accounts.NewClient(validUrl, time.Duration(expectedTimeoutClient))

	if err != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	// Assert

	if accountClient.BaseURL.Path != expectedPath {
		t.Errorf("client returned path: got %s want %s",
			accountClient.BaseURL.Path, expectedPath)
	}

	if accountClient.BaseURL.Host != expectedHost {
		t.Errorf("client returned host: got %s want %s",
			accountClient.BaseURL.Host, expectedHost)
	}

	if accountClient.BaseURL.Scheme != expectedScheme {
		t.Errorf("client returned scheme: got %s want %s",
			accountClient.BaseURL.Scheme, expectedScheme)
	}

	if accountClient.HttpClient.Timeout != expectedTimeoutClient {
		t.Errorf("client returned timeout: got %s want %s",
			accountClient.HttpClient.Timeout, expectedTimeoutClient)
	}
}

func TestNewClient_validUrlAndDefaultTimeoutValue_returnsValidClientWithDefaultTimeoutValue(t *testing.T) {

	// Arrange

	expectedScheme := "http"
	expectedPath := "/v1/organisation/accounts"
	expectedHost := "localhost:8080"
	zeroValueTimeout := time.Duration(0)
	expectedDefaultTimeoutValue := accounts.DefaultTimeOutValue

	// Act

	accountClient, err := accounts.NewClient(validUrl, zeroValueTimeout)

	if err != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	// Assert

	if accountClient.BaseURL.Path != expectedPath {
		t.Errorf("client returned path: got %s want %s",
			accountClient.BaseURL.Path, expectedPath)
	}

	if accountClient.BaseURL.Host != expectedHost {
		t.Errorf("client returned host: got %s want %s",
			accountClient.BaseURL.Host, expectedHost)
	}

	if accountClient.BaseURL.Scheme != expectedScheme {
		t.Errorf("client returned scheme: got %s want %s",
			accountClient.BaseURL.Scheme, expectedScheme)
	}

	if accountClient.HttpClient.Timeout != expectedDefaultTimeoutValue {
		t.Errorf("client returned timeout: got %s want %s",
			accountClient.HttpClient.Timeout, expectedDefaultTimeoutValue)
	}
}

func TestNewClient_emptyBaseUrl_returnsBaseUrlParsingError(t *testing.T) {

	// Arrange

	expectedErrorMessage := `parse "": empty url`

	expectedErrorType := accounts.BaseUrlParsingError

	// Act

	accountClient, err := accounts.NewClient("", time.Duration(1*time.Millisecond))

	// Assert

	if accountClient != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	assertClientInternalError(err, expectedErrorMessage, t, expectedErrorType)

}

func TestNewClient_invalidBaseUrl_returnsBaseUrlParsingError(t *testing.T) {

	// Arrange

	expectedErrorMessage := `parse "wrongURL": invalid URI for request`

	expectedErrorType := accounts.BaseUrlParsingError

	// Act

	accountClient, err := accounts.NewClient("wrongURL", time.Duration(1*time.Millisecond))

	// Assert

	if accountClient != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	assertClientInternalError(err, expectedErrorMessage, t, expectedErrorType)

}

func TestNewClient_invalidTimeoutValue_returnsValidClientWithDefaultValue(t *testing.T) {

	// Arrange

	expectedScheme := "http"
	expectedPath := "/v1/organisation/accounts"
	expectedHost := "localhost:8080"
	invalidValueTimeout := time.Duration(-1 * time.Millisecond)
	expectedDefaultTimeoutValue := accounts.DefaultTimeOutValue

	// Act

	accountClient, err := accounts.NewClient(validUrl, invalidValueTimeout)

	if err != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	// Assert

	if accountClient.BaseURL.Path != expectedPath {
		t.Errorf("client returned path: got %s want %s",
			accountClient.BaseURL.Path, expectedPath)
	}

	if accountClient.BaseURL.Host != expectedHost {
		t.Errorf("client returned host: got %s want %s",
			accountClient.BaseURL.Host, expectedHost)
	}

	if accountClient.BaseURL.Scheme != expectedScheme {
		t.Errorf("client returned scheme: got %s want %s",
			accountClient.BaseURL.Scheme, expectedScheme)
	}

	if accountClient.HttpClient.Timeout != expectedDefaultTimeoutValue {
		t.Errorf("client returned timeout: got %s want %s",
			accountClient.HttpClient.Timeout, expectedDefaultTimeoutValue)
	}
}

func TestNewClient_invalidTimeoutValueInNanoSeconds_returnsValidClientWithDefaultValue(t *testing.T) {

	// Arrange

	expectedScheme := "http"
	expectedPath := "/v1/organisation/accounts"
	expectedHost := "localhost:8080"
	invalidValueTimeout := time.Duration(50 * time.Nanosecond)
	expectedDefaultTimeoutValue := accounts.DefaultTimeOutValue

	// Act

	accountClient, err := accounts.NewClient(validUrl, invalidValueTimeout)

	if err != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	// Assert

	if accountClient.BaseURL.Path != expectedPath {
		t.Errorf("client returned path: got %s want %s",
			accountClient.BaseURL.Path, expectedPath)
	}

	if accountClient.BaseURL.Host != expectedHost {
		t.Errorf("client returned host: got %s want %s",
			accountClient.BaseURL.Host, expectedHost)
	}

	if accountClient.BaseURL.Scheme != expectedScheme {
		t.Errorf("client returned scheme: got %s want %s",
			accountClient.BaseURL.Scheme, expectedScheme)
	}

	if accountClient.HttpClient.Timeout != expectedDefaultTimeoutValue {
		t.Errorf("client returned timeout: got %s want %s",
			accountClient.HttpClient.Timeout, expectedDefaultTimeoutValue)
	}
}

func TestNewClient_invalidBaseUrlAndinvalidTimeout_returnsBaseUrlParsingError(t *testing.T) {

	// Arrange

	expectedErrorMessage := `parse "wrongURL": invalid URI for request`

	expectedErrorType := accounts.BaseUrlParsingError

	// Act

	accountClient, err := accounts.NewClient("wrongURL", time.Duration(-1*time.Millisecond))

	// Assert

	if accountClient != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	assertClientInternalError(err, expectedErrorMessage, t, expectedErrorType)
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

func assertBadStatusError(err error, expectedErrorMessage string, t *testing.T, expectedCorrectRequest string, expectedhttpStatus int) {
	if err != nil {
		if errors.Is(err, accounts.ApiHttpErrorType) {

			tempErr := fmt.Errorf("%w | Path: %s returned %d with message %s", accounts.ApiHttpErrorType, expectedCorrectRequest, expectedhttpStatus, expectedErrorMessage)

			if err.Error() != tempErr.Error() {
				t.Errorf("Returned error message: got %s want %s",
					err.Error(), tempErr.Error())
			}
		}
	}
}

func assertClientInternalError(err error, expectedErrorMessage string, t *testing.T, expectedErrorType error) {
	if err != nil {

		if errors.Is(err, expectedErrorType) {

			testErr := fmt.Errorf("%w | %s", expectedErrorType, expectedErrorMessage)

			if err.Error() != testErr.Error() {
				t.Errorf("Returned error message: got %s want %s",
					err.Error(), testErr.Error())
			}

		}

	}
}

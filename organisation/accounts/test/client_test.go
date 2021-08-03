package accounts_test

import (
	"ei09010/form3-api-client/organisation/accounts"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// test names follow the Given_When_Then naming taxonomy

//New Client Tests

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

	assertClientError(err, expectedErrorMessage, t, expectedErrorType)

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

	assertClientError(err, expectedErrorMessage, t, expectedErrorType)

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

	assertClientError(err, expectedErrorMessage, t, expectedErrorType)
}

// aux

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

func assertClientError(err error, expectedErrorMessage string, t *testing.T, expectedErrorType int) {
	if err != nil {
		if cerr, ok := err.(*accounts.ClientError); ok {

			if cerr.ErrorMessage != expectedErrorMessage {
				t.Errorf("Returned error message: got %s want %s",
					cerr.ErrorMessage, expectedErrorMessage)
			}

			if cerr.ErrorType != expectedErrorType {
				t.Errorf("Returned error type: got %v want %v",
					cerr.ErrorType, expectedErrorType)
			}

			if cerr.BadStatusError != nil {
				t.Errorf("Returned bad status error: got %v want %v",
					cerr.BadStatusError, nil)
			}

		} else {
			t.Errorf("returned error isn't a %T, got %T", err.(*accounts.ClientError), err)
		}

	}
}

func assertBadStatusError(err error, expectedErrorMessage string, t *testing.T, expectedErrorType int, expectedCorrectRequest string, expectedhttpStatus int) {
	if err != nil {
		if cerr, ok := err.(*accounts.ClientError); ok {

			if cerr.ErrorMessage != expectedErrorMessage {
				t.Errorf("Returned error message: got %s want %s",
					cerr.ErrorMessage, expectedErrorMessage)
			}

			if cerr.ErrorType != expectedErrorType {
				t.Errorf("Returned error type: got %v want %v",
					cerr.ErrorType, expectedErrorType)
			}

			if cerr.BadStatusError != nil {

				if cerr.BadStatusError.URL != expectedCorrectRequest {
					t.Errorf("Returned url status for bad status error: got %s want %s",
						cerr.BadStatusError.URL, expectedCorrectRequest)
				}

				if cerr.BadStatusError.HttpCode != expectedhttpStatus {
					t.Errorf("Returned http status code for bad status error: got %d want %d",
						cerr.BadStatusError.HttpCode, expectedhttpStatus)
				}

			} else {
				t.Errorf("Bad status error is %v", cerr.BadStatusError)
			}
		} else {
			t.Errorf("returned error isn't a %T, got %T", err.(*accounts.ClientError), err)
		}
	}
}

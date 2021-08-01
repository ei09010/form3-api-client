package accounts

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
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

	accountClient, err := NewClient(validUrl, time.Duration(expectedTimeoutClient))

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
	expectedDefaultTimeoutValue := DefaultTimeOutValue

	// Act

	accountClient, err := NewClient(validUrl, zeroValueTimeout)

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

	expectedErrorType := BaseUrlParsingError

	// Act

	accountClient, err := NewClient("", time.Duration(1*time.Millisecond))

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

	expectedErrorType := BaseUrlParsingError

	// Act

	accountClient, err := NewClient("wrongURL", time.Duration(1*time.Millisecond))

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
	expectedDefaultTimeoutValue := DefaultTimeOutValue

	// Act

	accountClient, err := NewClient(validUrl, invalidValueTimeout)

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
	expectedDefaultTimeoutValue := DefaultTimeOutValue

	// Act

	accountClient, err := NewClient(validUrl, invalidValueTimeout)

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

	expectedErrorType := BaseUrlParsingError

	// Act

	accountClient, err := NewClient("wrongURL", time.Duration(-1*time.Millisecond))

	// Assert

	if accountClient != nil {
		t.Errorf("Returned reponse: got %v want %v",
			accountClient, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedErrorType)
}

// FETCH TESTS
func TestFetch_validAccountId_returnsAccountsData(t *testing.T) {

	// Arrange

	expectedCorrectResponse := `{"data":{"attributes":{"account_classification":"Personal","alternative_names":["Alternative Names."],"bank_id":"400300","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB22","country":"GB","name":["Name of the account holder, up to four lines possible."]},"created_on":"2021-07-31T22:09:02.680Z","id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc","modified_on":"2021-07-31T22:09:02.680Z","organisation_id":"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c","type":"accounts","version":0},"links":{"self":"/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"}}`
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, expectedCorrectResponse)
	})

	defer ts.Close()

	expectedAccountClassification := "Personal"
	expectedAlternativeNames := []string{"Alternative Names."}
	expectedBankId := "400300"
	expectedBankIdCode := "GBDSC"
	expectedBaseCurrency := "GBP"
	expectedBic := "NWBKGB22"
	expectedCountry := "GB"
	expectedName := []string{"Name of the account holder, up to four lines possible."}
	expectedCreatedOn := "2021-07-31 22:09:02.68 +0000 UTC"
	expectedId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	expectedModifiedOn := "2021-07-31 22:09:02.68 +0000 UTC"
	expectedOrganisationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	expectedType := "accounts"
	expectedVersion := 0
	expectedSelf := "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	accountClient, err := NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	response, err := accountClient.Fetch(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Assert

	// Data

	if response.Data.ID != expectedId {
		t.Errorf("handler returned unexpected Id: got %s want %s",
			response.Data.ID, expectedId)
	}

	if response.Data.OrganisationID != expectedOrganisationId {
		t.Errorf("handler returned unexpected organisationId: got %s want %s",
			response.Data.OrganisationID, expectedOrganisationId)
	}

	if response.Data.Type != expectedType {
		t.Errorf("handler returned unexpected type: got %s want %s",
			response.Data.Type, expectedType)
	}

	if response.Data.Version != expectedVersion {
		t.Errorf("handler returned unexpected version: got %d want %d",
			response.Data.Version, expectedVersion)
	}

	if response.Data.CreatedOn.String() != expectedCreatedOn {
		t.Errorf("handler returned unexpected created on: got %s want %s",
			response.Data.CreatedOn, expectedCreatedOn)
	}

	if response.Data.ModifiedOn.String() != expectedModifiedOn {
		t.Errorf("handler returned unexpected attributes Bic: got %s want %s",
			response.Data.ModifiedOn, expectedModifiedOn)
	}

	// attributes

	if response.Data.Attributes.AccountClassification != expectedAccountClassification {
		t.Errorf("handler returned unexpected attributes account classification: got %s want %s",
			response.Data.Attributes.AccountClassification, expectedAccountClassification)
	}

	if !equal(response.Data.Attributes.AlternativeNames, expectedAlternativeNames) {
		t.Errorf("handler returned unexpecteda attributes alternative names: got %s want %s",
			response.Data.Attributes.AlternativeNames, expectedAlternativeNames)
	}

	if !equal(response.Data.Attributes.Name, expectedName) {
		t.Errorf("handler returned unexpected attributes name: got %s want %s",
			response.Data.Attributes.Name, expectedName)
	}

	if response.Data.Attributes.Country != expectedCountry {
		t.Errorf("handler returned unexpected attributes country: got %s want %s",
			response.Data.Attributes.Country, expectedCountry)
	}

	if response.Data.Attributes.BaseCurrency != expectedBaseCurrency {
		t.Errorf("handler returned unexpected attributes base currency: got %s want %s",
			response.Data.Attributes.BaseCurrency, expectedBaseCurrency)
	}

	if response.Data.Attributes.BankID != expectedBankId {
		t.Errorf("handler returned unexpected attributes bank id: got %s want %s",
			response.Data.Attributes.BankID, expectedBankId)
	}

	if response.Data.Attributes.BankIDCode != expectedBankIdCode {
		t.Errorf("handler returned unexpected attributes BankIdCode: got %s want %s",
			response.Data.Attributes.BankIDCode, expectedBankIdCode)
	}

	if response.Data.Attributes.Bic != expectedBic {
		t.Errorf("handler returned unexpected attributes Bic: got %s want %s",
			response.Data.Attributes.Bic, expectedBic)
	}

	// Links

	if response.Links.Self != expectedSelf {
		t.Errorf("handler returned unexpected Link self: got %s want %s",
			response.Links.Self, expectedSelf)
	}

}

func TestFetch_notFoundAccountId_returnsNotFoundStatusError(t *testing.T) {

	// Arrange

	expectedErrorMessageResponse := `{"error_message":"record ad27e265-9605-4b4b-a0e5-3003ea9cc4dc does not exist"}`
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`

	expectedErrorMessage := "record ad27e265-9605-4b4b-a0e5-3003ea9cc4dc does not exist"
	expectedErrorType := HttResponseStandardError
	expectedhttpStatus := http.StatusNotFound

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedhttpStatus)
		io.WriteString(w, expectedErrorMessageResponse)
	})

	defer ts.Close()

	accountClient, err := NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	response, err := accountClient.Fetch(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))

	// Assert

	if response != nil {
		t.Errorf("Returned reponse: got %v want %v",
			response, nil)
	}

	assertBadStatusError(err, expectedErrorMessage, t, expectedErrorType, expectedCorrectRequest, expectedhttpStatus)
}

func TestFetch_emptyResponseInInternalError_returnsUnmarshallingError(t *testing.T) {

	// Arrange
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`

	expectedErrorType := UnmarshallingError
	expectedErrorMessage := "unexpected end of JSON input"

	expectedhttpStatus := http.StatusInternalServerError

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedhttpStatus)
	})

	defer ts.Close()

	accountClient, err := NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	response, err := accountClient.Fetch(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))

	// Assert

	if response != nil {
		t.Errorf("Returned reponse: got %v want %v",
			response, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedErrorType)
}

func TestFetch_emptyResponseInSuccessResponse_returnsUnmarshallingError(t *testing.T) {

	// Arrange
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`

	expectedErrorType := UnmarshallingError
	expectedErrorMessage := "unexpected end of JSON input"

	expectedhttpStatus := http.StatusOK

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedhttpStatus)
	})

	defer ts.Close()

	accountClient, err := NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	response, err := accountClient.Fetch(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))

	// Assert

	if response != nil {
		t.Errorf("Returned reponse: got %v want %v",
			response, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedErrorType)
}

func TestFetch_invalidResponseBodyInSuccessResponse_returnsUnmarshallingError(t *testing.T) {

	// Arrange
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`
	expectedResponseInvalidJson := `{"data":unt_classion":"Personal222","native_names":["Alternative Names."],"ban":"400300","bank_id_code":"G","base_currency":"GBP","bic":"NWBKGB22","country":"GB","name":["Name of the account holder, up to four lines possible."]},"created_on":"2021-07-31T22:09:02.680Z","id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc","modified_on":"2021-07-31T22:09:02.680Z","organisation_id":"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c","type":"accounts","version":0},"links":{"self":"/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"}}`

	expectedErrorType := UnmarshallingError
	expectedErrorMessage := "invalid character 'u' looking for beginning of value"

	expectedhttpStatus := http.StatusOK

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedhttpStatus)
		io.WriteString(w, expectedResponseInvalidJson)
	})

	defer ts.Close()

	accountClient, err := NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	response, err := accountClient.Fetch(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))

	// Assert

	if response != nil {
		t.Errorf("Returned reponse: got %v want %v",
			response, nil)
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
		if cerr, ok := err.(*ClientError); ok {

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
			t.Errorf("returned error isn't a %T, got %T", err.(*ClientError), err)
		}

	}
}

func assertBadStatusError(err error, expectedErrorMessage string, t *testing.T, expectedErrorType int, expectedCorrectRequest string, expectedhttpStatus int) {
	if err != nil {
		if cerr, ok := err.(*ClientError); ok {

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
			t.Errorf("returned error isn't a %T, got %T", err.(*ClientError), err)
		}
	}
}

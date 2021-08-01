package accounts

import (
	"fmt"
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

	if err != nil {
		if cerr, ok := err.(*ClientError); ok {

			fmt.Println(cerr.ErrorMessage)

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

	if err != nil {
		if cerr, ok := err.(*ClientError); ok {

			fmt.Println(cerr.ErrorMessage)

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

	if err != nil {
		if cerr, ok := err.(*ClientError); ok {

			fmt.Println(cerr.ErrorMessage)

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

// FETCH TESTS
func TestFetch_validAccountId_returnsAccountsData(t *testing.T) {

	// Arrange

	expectedCorrectResponse := `{"data":{"attributes":{"account_classification":"Personal","alternative_names":["Alternative Names."],"bank_id":"400300","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB22","country":"GB","name":["Name of the account holder, up to four lines possible."]},"created_on":"2021-07-31T22:09:02.680Z","id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc","modified_on":"2021-07-31T22:09:02.680Z","organisation_id":"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c","type":"accounts","version":0},"links":{"self":"/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"}}`
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.URL.Path == expectedCorrectRequest {
			io.WriteString(w, expectedCorrectResponse)
		} else {
			io.WriteString(w, "Bad request")
		}
	}))

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

// aux

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

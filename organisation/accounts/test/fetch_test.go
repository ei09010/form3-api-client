package accounts_test

import (
	"bytes"
	"ei09010/form3-api-client/organisation/accounts"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
)

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (mckHt *MockHttpClient) Do(req *http.Request) (*http.Response, error) {

	return mckHt.DoFunc(req)
}

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

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

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
	expectedHttpStatus := http.StatusNotFound
	expectedErrorType := accounts.ApiHttpErrorType

	r := ioutil.NopCloser(bytes.NewReader([]byte(expectedErrorMessageResponse)))

	accountErrClient := &MockHttpClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: expectedHttpStatus,
				Request:    &http.Request{URL: &url.URL{Path: expectedCorrectRequest}},
				Body:       r,
			}, nil
		},
	}

	accountsClient := &accounts.Client{
		BaseURL:    &url.URL{Path: expectedCorrectRequest},
		HttpClient: accountErrClient,
	}

	// Act

	response, err := accountsClient.Fetch(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))

	if response != nil {
		t.Errorf("Returned reponse: got %v want %v",
			response, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

func TestFetch_emptyResponseInInternalError_returnsUnmarshallingError(t *testing.T) {

	// Arrange
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`

	expectedErrorType := accounts.UnmarshallingError
	expectedErrorMessage := "unexpected end of JSON input"

	expectedHttpStatus := http.StatusInternalServerError

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

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

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

func TestFetch_emptyResponseInSuccessResponse_returnsUnmarshallingError(t *testing.T) {

	// Arrange
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`

	expectedErrorType := accounts.UnmarshallingError
	expectedErrorMessage := "unexpected end of JSON input"

	expectedHttpStatus := http.StatusBadRequest

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

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

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

func TestFetch_invalidResponseBodyInSuccessResponse_returnsUnmarshallingError(t *testing.T) {

	// Arrange
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`
	expectedResponseInvalidJson := `{"data":unt_classion":"Personal222","native_names":["Alternative Names."],"ban":"400300","bank_id_code":"G","base_currency":"GBP","bic":"NWBKGB22","country":"GB","name":["Name of the account holder, up to four lines possible."]},"created_on":"2021-07-31T22:09:02.680Z","id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc","modified_on":"2021-07-31T22:09:02.680Z","organisation_id":"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c","type":"accounts","version":0},"links":{"self":"/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"}}`

	expectedErrorType := accounts.UnmarshallingError
	expectedErrorMessage := "invalid character 'u' looking for beginning of value"

	expectedHttpStatus := http.StatusOK

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)
		io.WriteString(w, expectedResponseInvalidJson)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

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

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

// add tests:
// - 404 with un expected error message

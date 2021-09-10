package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestCreate_validAccountData_returnsStoredAccountsData(t *testing.T) {

	// Arrange

	expectedCorrectBody := `{"data":{"attributes":{"account_classification":"Personal","alternative_names":["特别的."],"bank_id":"400300","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB22","country":"GB","name":["Name of the account holder, up to four lines possible."]},"created_on":"2021-07-31T22:09:02Z","id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc","modified_on":"2021-07-31T22:09:02Z","organisation_id":"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c","type":"accounts","version":0},"links":{"self":"/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"}}`
	expectedCorrectRequest := `/v1/organisation/accounts`

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {

		bodyBytes, err := ioutil.ReadAll(r.Body)

		receivedBody := string(bodyBytes)

		if err != nil {
			t.Errorf(err.Error())
		}

		if receivedBody == expectedCorrectBody {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")

			io.WriteString(w, expectedCorrectBody)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	defer ts.Close()

	expectedAccountClassification := "Personal"
	expectedAlternativeNames := []string{"特别的."}
	expectedBankId := "400300"
	expectedBankIdCode := "GBDSC"
	expectedBaseCurrency := "GBP"
	expectedBic := "NWBKGB22"
	expectedCountry := "GB"
	expectedName := []string{"Name of the account holder, up to four lines possible."}

	timeLayout := "2006-01-02 15:04:05 -0700 MST"
	expectedCreatedOn := "2021-07-31 22:09:02 +0000 UTC"
	expectedCreatedOnTime, err := time.Parse(timeLayout, expectedCreatedOn)

	expectedId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	expectedModifiedOn := "2021-07-31 22:09:02 +0000 UTC"
	expectedModifiedOnTime, err := time.Parse(timeLayout, expectedModifiedOn)

	expectedOrganisationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	expectedType := "accounts"
	expectedVersion := 0
	expectedSelf := "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	accountClient, err := NewClient(WithBaseURL(ts.URL), WithTimeout(time.Duration(1000*time.Millisecond)))

	if err != nil {
		t.Errorf(err.Error())
	}

	accountDataToCreate := &AccountData{
		Data: &Data{
			Attributes: &AccountAttributes{
				AccountClassification: expectedAccountClassification,
				AlternativeNames:      expectedAlternativeNames,
				BankID:                expectedBankId,
				BankIDCode:            expectedBankIdCode,
				BaseCurrency:          expectedBaseCurrency,
				Bic:                   expectedBic,
				Country:               expectedCountry,
				Name:                  expectedName,
			},
			CreatedOn:      expectedCreatedOnTime,
			ID:             expectedId,
			ModifiedOn:     expectedModifiedOnTime,
			OrganisationID: expectedOrganisationId,
			Type:           expectedType,
			Version:        expectedVersion,
		},
		Links: &Links{
			Self: expectedSelf,
		},
	}

	ctx := context.Background()

	// Act

	response, err := accountClient.Create(ctx, accountDataToCreate)

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
		t.Errorf("handler returned unexpected attributes alternative names: got %s want %s",
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

func TestCreate_DuplicateAccountData_returns409StatusConflict(t *testing.T) {

	// Arrange

	expectedCorrectResponse := `{"error_message":"Account cannot be created as it violates a duplicate constraint"}`
	expectedCorrectRequest := `/v1/organisation/accounts`
	expectedErrorMessage := `Account cannot be created as it violates a duplicate constraint`
	expectedStatus := http.StatusConflict

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
		io.WriteString(w, expectedCorrectResponse)
	})

	defer ts.Close()

	ctx := context.Background()

	accountClient, err := NewClient(WithBaseURL(ts.URL), WithTimeout(time.Duration(1000*time.Millisecond)))

	// Act

	response, err := accountClient.Create(ctx, generateValidGenericAccountData())

	// Assert

	if response != nil {
		t.Errorf("Returned reponse: got %v want %v",
			response, nil)
	}

	assertClientError(err, expectedErrorMessage, t, ApiHttpErrorType, expectedStatus)

}

func TestCreate_InvalidAccountData_returns400BadRequest(t *testing.T) {

	// Arrange

	expectedCorrectResponse := `{"error_message":"validation failure list:\nvalidation failure list:\nid in body is required"}`
	expectedCorrectRequest := `/v1/organisation/accounts`
	expectedErrorMessage := "validation failure list:\nvalidation failure list:\nid in body is required"
	expectedStatus := http.StatusBadRequest

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
		io.WriteString(w, expectedCorrectResponse)
	})

	defer ts.Close()

	ctx := context.Background()

	accountClient, err := NewClient(WithBaseURL(ts.URL), WithTimeout(time.Duration(1000*time.Millisecond)))

	// Act

	response, err := accountClient.Create(ctx, generateValidGenericAccountData())

	// Assert

	if response != nil {
		t.Errorf("Returned reponse: got %v want %v",
			response, nil)
	}

	assertClientError(err, expectedErrorMessage, t, ApiHttpErrorType, expectedStatus)

}

func TestCreateErrorCases(t *testing.T) {

	// Arrange
	errorCases := map[string]struct {
		accountId            string
		messageResponse      string
		requestPath          string
		expectedErrorMessage string
		expectedHttpStatus   int
		expectedErrorType    error
		doError              error
		accountPayload       *AccountData
	}{
		"Duplicated contraint value": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			messageResponse:      `{"error_message":"Account cannot be created as it violates a duplicate constraint"}`,
			requestPath:          `/v1/organisation/accounts`,
			expectedErrorMessage: "Account cannot be created as it violates a duplicate constraint",
			expectedHttpStatus:   http.StatusConflict,
			expectedErrorType:    ApiHttpErrorType,
			accountPayload:       generateValidGenericAccountData(),
		},
		"Mandatory field missing": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			messageResponse:      `{"error_message":"validation failure list:\nvalidation failure list:\nid in body is required"}`,
			requestPath:          `/v1/organisation/accounts`,
			expectedErrorMessage: "validation failure list:\nvalidation failure list:\nid in body is required",
			expectedHttpStatus:   http.StatusBadRequest,
			expectedErrorType:    ApiHttpErrorType,
			accountPayload:       generateAccountDataWithInvalidID(""),
		},
		"Invalid mandatory field": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			messageResponse:      `{"error_message":"validation failure list:\nvalidation failure list:\nid in body must be of type uuid: \"notValidContent\""}`,
			requestPath:          `/v1/organisation/accounts`,
			expectedErrorMessage: "validation failure list:\nvalidation failure list:\nid in body must be of type uuid: \"notValidContent\"",
			expectedHttpStatus:   http.StatusBadRequest,
			expectedErrorType:    ApiHttpErrorType,
			accountPayload:       generateAccountDataWithInvalidID("notValidContent"),
		},
		"Empty response in internal errror": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			messageResponse:      "",
			requestPath:          `/v1/organisation/accounts`,
			expectedErrorMessage: "EOF",
			expectedHttpStatus:   http.StatusInternalServerError,
			expectedErrorType:    BuildingRequestError,
			accountPayload:       generateValidGenericAccountData(),
		},
		"Handler timeout causes the request to fail": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			messageResponse:      "",
			requestPath:          `/v1/organisation/accounts`,
			expectedErrorMessage: `Requesting "handler url": http: Handler timeout`,
			expectedHttpStatus:   http.StatusInternalServerError,
			expectedErrorType:    BuildingRequestError,
			doError: &url.Error{
				Err: http.ErrHandlerTimeout,
				Op:  "Requesting",
				URL: "handler url",
			},
			accountPayload: generateValidGenericAccountData(),
		},
		"Nil http response": {
			expectedErrorType:    BuildingRequestError,
			expectedHttpStatus:   http.StatusBadRequest,
			expectedErrorMessage: `Unmarshling "handler url": json: error calling MarshalJSON for type string: marshalling error`,
			doError: &url.Error{
				Err: &json.MarshalerError{
					Type: reflect.TypeOf("batata"),
					Err:  errors.New("marshalling error"),
				},
				Op:  "Unmarshling",
				URL: "handler url",
			},
		},
	}

	for _, errCase := range errorCases {

		var mockBehavior func(*http.Request) (*http.Response, error)

		if errCase.accountId != "" {
			mockBehavior = func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: errCase.expectedHttpStatus,
					Request:    &http.Request{URL: &url.URL{Path: errCase.requestPath}},
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(errCase.messageResponse))),
				}, errCase.doError
			}
		} else {
			mockBehavior = func(*http.Request) (*http.Response, error) { return nil, errCase.doError }
		}

		accountErrClient := &MockHttpClient{
			DoFunc: mockBehavior,
		}

		accountsClient := &Client{
			baseURL:    &url.URL{Path: errCase.requestPath},
			httpClient: accountErrClient,
		}

		ctx := context.Background()

		// Act

		response, err := accountsClient.Create(ctx, errCase.accountPayload)

		// Assert

		if response != nil {
			t.Errorf("Returned reponse: got %v want %v",
				response, nil)
		}

		assertClientError(err, errCase.expectedErrorMessage, t, errCase.expectedErrorType, errCase.expectedHttpStatus)

	}

}

func generateAccountDataWithInvalidID(invalidIdContent string) *AccountData {

	tempAccountData := generateValidGenericAccountData()

	tempAccountData.Data.ID = invalidIdContent

	return tempAccountData
}

func generateValidGenericAccountData() *AccountData {

	expectedAccountClassification := "Personal"
	expectedAlternativeNames := []string{"Alternative Names."}
	expectedBankId := "400300"
	expectedBankIdCode := "GBDSC"
	expectedBaseCurrency := "GBP"
	expectedBic := "NWBKGB22"
	expectedCountry := "GB"
	expectedName := []string{"Name of the account holder, up to four lines possible."}

	timeLayout := "2006-01-02 15:04:05 -0700 MST"
	expectedCreatedOn := "2021-07-31 22:09:02 +0000 UTC"
	expectedCreatedOnTime, _ := time.Parse(timeLayout, expectedCreatedOn)

	expectedId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	expectedModifiedOn := "2021-07-31 22:09:02 +0000 UTC"
	expectedModifiedOnTime, _ := time.Parse(timeLayout, expectedModifiedOn)

	expectedOrganisationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	expectedType := "accounts"
	expectedVersion := 0
	expectedSelf := "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	return &AccountData{
		Data: &Data{
			Attributes: &AccountAttributes{
				AccountClassification: expectedAccountClassification,
				AlternativeNames:      expectedAlternativeNames,
				BankID:                expectedBankId,
				BankIDCode:            expectedBankIdCode,
				BaseCurrency:          expectedBaseCurrency,
				Bic:                   expectedBic,
				Country:               expectedCountry,
				Name:                  expectedName,
			},
			CreatedOn:      expectedCreatedOnTime,
			ID:             expectedId,
			ModifiedOn:     expectedModifiedOnTime,
			OrganisationID: expectedOrganisationId,
			Type:           expectedType,
			Version:        expectedVersion,
		},
		Links: &Links{
			Self: expectedSelf,
		},
	}
}

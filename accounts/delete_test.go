package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDelete_ExistingAccountId_Returns200WithNoBody(t *testing.T) {

	// Arrange

	validAccountId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	validVersion := 0
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`
	expectedRawQuery := `version=0`

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {

		if r.URL.RawQuery == expectedRawQuery {
			w.WriteHeader(http.StatusOK)
		}
	})

	defer ts.Close()

	accountClient, err := NewClient(WithBaseURL(ts.URL), WithTimeout(time.Duration(1000*time.Millisecond)))

	if err != nil {
		t.Errorf(err.Error())
	}

	ctx := context.Background()

	// Act

	err = accountClient.Delete(ctx, uuid.MustParse(validAccountId), validVersion)

	// Assert

	if err != nil {
		t.Errorf("delete returned an error: got %v want %v",
			err, nil)
	}
}

func TestDelete_NotExistingAccountId_Returns404WithNoBody(t *testing.T) {

	// Arrange

	validAccountId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	validVersion := 0
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`
	expectedRawQuery := `version=0`

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {

		if r.URL.RawQuery == expectedRawQuery {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	defer ts.Close()

	accountClient, err := NewClient(WithBaseURL(ts.URL), WithTimeout(time.Duration(1000*time.Millisecond)))

	if err != nil {
		t.Errorf(err.Error())
	}

	ctx := context.Background()

	// Act

	err = accountClient.Delete(ctx, uuid.MustParse(validAccountId), validVersion)

	// Assert

	assertClientError(err, "", t, ApiHttpErrorType, http.StatusNotFound)
}

func TestDelete_InternalServerError_Returns500WithNoBody(t *testing.T) {

	// Arrange

	validAccountId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	validVersion := 0
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`
	expectedRawQuery := `version=0`

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery == expectedRawQuery {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	defer ts.Close()

	accountClient, err := NewClient(WithBaseURL(ts.URL), WithTimeout(time.Duration(1000*time.Millisecond)))

	if err != nil {
		t.Errorf(err.Error())
	}

	ctx := context.Background()

	// Act

	err = accountClient.Delete(ctx, uuid.MustParse(validAccountId), validVersion)

	// Assert

	assertClientError(err, "", t, ApiHttpErrorType, http.StatusInternalServerError)
}

func TestDeleteErrorCases(t *testing.T) {

	// Arrange
	errorCases := map[string]struct {
		accountId            string
		version              int
		expectedHttpStatus   int
		requestPath          string
		messageResponse      string
		expectedErrorMessage string
		expectedErrorType    error
		doError              error
		nilResponse          bool
	}{
		"Invalid version": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:              1,
			expectedHttpStatus:   http.StatusConflict,
			requestPath:          `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=1`,
			messageResponse:      `{"error_message": "invalid version"}`,
			expectedErrorMessage: "invalid version",
			expectedErrorType:    ApiHttpErrorType,
		},
		"Invalid UUID": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:              1,
			expectedHttpStatus:   http.StatusBadRequest,
			requestPath:          `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=1`,
			messageResponse:      `{"error_message":"id is not a valid uuid"}`,
			expectedErrorMessage: `id is not a valid uuid`,
			expectedErrorType:    ApiHttpErrorType,
		},
		"No version number in query": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:              0,
			expectedHttpStatus:   http.StatusBadRequest,
			requestPath:          `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`,
			messageResponse:      `{"error_message":"invalid version number"}`,
			expectedErrorMessage: "invalid version number",
			expectedErrorType:    ApiHttpErrorType,
		},
		"Handler timeout causes the request to fail": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:              0,
			expectedHttpStatus:   http.StatusInternalServerError,
			requestPath:          `/v1/organisation/accounts`,
			messageResponse:      "",
			expectedErrorMessage: `Requesting "handler url": http: Handler timeout`,
			expectedErrorType:    BuildingRequestError,
			doError: &url.Error{
				Err: http.ErrHandlerTimeout,
				Op:  "Requesting",
				URL: "handler url",
			},
		},
		"Unmarshalling Error reading error response message": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:              0,
			expectedHttpStatus:   http.StatusBadRequest,
			requestPath:          `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc`,
			messageResponse:      `rror":"invalid version number"}`,
			expectedErrorMessage: "invalid character 'r' looking for beginning of value",
			expectedErrorType:    BuildingRequestError,
		},
		"Nil http response": {
			accountId:            "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:              0,
			nilResponse:          true,
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

		if !errCase.nilResponse {
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

		err := accountsClient.Delete(ctx, uuid.MustParse(errCase.accountId), errCase.version)

		// Assert

		if err == nil {
			t.Errorf("Returned reponse: got %v want %v",
				err, nil)
		}

		assertClientError(err, errCase.expectedErrorMessage, t, errCase.expectedErrorType, errCase.expectedHttpStatus)

	}

}

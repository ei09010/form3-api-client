package accounts_test

import (
	"ei09010/form3-api-client/organisation/accounts"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDelete_existingAccountId_DoesntReturnError(t *testing.T) {

	// Arrange

	validAccountId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	validVersion := 0
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=0`

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(1000*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	err = accountClient.Delete(uuid.MustParse(validAccountId), validVersion)

	// Assert

	if err != nil {
		t.Errorf("delete returned an error: got %v want %v",
			err, nil)
	}
}

func TestDelete_InvalidVersion_returnsConflictStatusError(t *testing.T) {

	// Arrange

	expectedErrorMessageResponse := `{"error_message": "invalid version"}`
	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=1`

	expectedErrorMessage := "invalid version"
	expectedHttpStatus := http.StatusNotFound
	expectedErrorType := accounts.ApiHttpErrorType

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)
		io.WriteString(w, expectedErrorMessageResponse)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	err = accountClient.Delete(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"), 1)

	// Assert

	if err == nil {
		t.Errorf("Returned reponse: got %v want %v",
			err, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

func TestDelete_nonExistingAccountId_returnsNotFoundStatusError(t *testing.T) {

	// Arrange

	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=1`
	expectedErrorMessage := ""
	expectedHttpStatus := http.StatusNotFound
	expectedErrorType := accounts.ApiHttpErrorType

	ts := newTestServer("expectedCorrectRequest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)

	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	err = accountClient.Delete(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"), 1)

	// Assert

	if err == nil {
		t.Errorf("Returned reponse: got %v want %v",
			err, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

func TestDelete_noVersionNumberInQuery_returnsBadRequestError(t *testing.T) {

	// Arrange

	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=1`
	expectedErrorMessage := "invalid version number"
	expectedErrorType := accounts.ApiHttpErrorType
	expectedHttpStatus := http.StatusBadRequest
	expectedResponse := `{"error_message":"invalid version number"}`

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)
		io.WriteString(w, expectedResponse)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	err = accountClient.Delete(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"), 1)

	// Assert

	if err == nil {
		t.Errorf("Returned reponse: got %v want %v",
			err, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

func TestDelete_internalServerError_returnsInternalServerError(t *testing.T) {

	// Arrange

	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=1`
	expectedErrorMessage := ""

	expectedHttpStatus := http.StatusInternalServerError
	expectedErrorType := accounts.ApiHttpErrorType

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	err = accountClient.Delete(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"), 1)

	// Assert

	if err == nil {
		t.Errorf("Returned reponse: got %v want %v",
			err, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

func TestDelete_invalidUUID_returnsBadRequest(t *testing.T) {

	// Arrange

	expectedCorrectRequest := `/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=1`
	expectedErrorMessage := `id is not a valid uuid`

	expectedResponse := `{"error_message":"id is not a valid uuid"}`

	expectedHttpStatus := http.StatusBadRequest
	expectedErrorType := accounts.ApiHttpErrorType

	ts := newTestServer(expectedCorrectRequest, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedHttpStatus)
		io.WriteString(w, expectedResponse)
	})

	defer ts.Close()

	accountClient, err := accounts.NewClient(ts.URL, time.Duration(100*time.Millisecond))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Act

	err = accountClient.Delete(uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"), 1)

	// Assert

	if err == nil {
		t.Errorf("Returned reponse: got %v want %v",
			err, nil)
	}

	assertClientError(err, expectedErrorMessage, t, expectedCorrectRequest, expectedHttpStatus, expectedErrorType)
}

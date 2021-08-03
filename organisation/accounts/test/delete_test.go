package accounts_test

import (
	"ei09010/form3-api-client/organisation/accounts"
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

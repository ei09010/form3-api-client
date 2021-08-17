package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) Create(accountData *AccountData) (*AccountData, error) {

	var err error

	c.baseURL, err = c.baseURL.Parse(c.baseURL.Path)

	if err != nil {

		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.baseURL.Path, http.StatusBadRequest, err.Error())
	}

	return c.postRequest(accountData)

}

func (c *Client) postRequest(accountData *AccountData) (*AccountData, error) {

	accountDataStr, err := json.Marshal(accountData)

	if err != nil {
		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.baseURL.Path, http.StatusBadRequest, err.Error())
	}

	postBody := bytes.NewBuffer(accountDataStr)

	customReq, err := http.NewRequest(http.MethodPost, c.baseURL.String(), postBody)

	if err != nil {

		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.baseURL.Path, http.StatusBadRequest, err.Error())
	}

	customReq.Header.Set("content-encoding", "application/json")
	customReq.Header.Set("user-agent", "golang-sdk")

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {

		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ExecutingRequestError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
	}

	defer httpResponse.Body.Close()

	var accountsData AccountData

	if httpResponse != nil {

		responseBody, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode == http.StatusOK {

			err = json.Unmarshal(responseBody, &accountsData)

			if err != nil {
				return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", UnmarshallingError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
			}

			return &accountsData, nil
		} else {

			apiHttpError := &apiErrorMessage{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", UnmarshallingError, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
			}

			return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
		}

	} else {

		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
	}
}

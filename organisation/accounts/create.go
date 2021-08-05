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

	c.BaseURL, err = c.BaseURL.Parse(c.BaseURL.Path)

	if err != nil {

		//clientErr := handleClientError(FinalUrlParsingError, err.Error())

		return nil, fmt.Errorf("%w | %s", FinalUrlParsingError, err.Error()) //clientErr
	}

	return c.postRequest(accountData)

}

func (c *Client) postRequest(accountData *AccountData) (*AccountData, error) {

	accountDataStr, err := json.Marshal(accountData)

	if err != nil {
		return nil, fmt.Errorf("%w | %s", BuildingRequestError, err.Error()) //handleClientError(BuildingRequestError, err.Error())
	}

	postBody := bytes.NewBuffer(accountDataStr)

	customReq, err := http.NewRequest("POST", c.BaseURL.String(), postBody)

	if err != nil {

		return nil, fmt.Errorf("%w | %s", BuildingRequestError, err.Error()) //handleClientError(BuildingRequestError, err.Error())
	}

	// review headers
	customReq.Header.Set("content-encoding", "application/json")
	customReq.Header.Set("user-agent", "golang-sdk")

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {

		return nil, fmt.Errorf("%w | %s", ExecutingRequestError, err.Error()) //handleClientError(ExecutingRequestError, err.Error())
	}

	defer httpResponse.Body.Close()

	var accountsData AccountData

	if httpResponse != nil {

		responseBody, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return nil, fmt.Errorf("%w | %s", ResponseReadError, err.Error()) // handleClientError(ResponseReadError, err.Error())
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode == http.StatusOK {

			err = json.Unmarshal(responseBody, &accountsData)

			if err != nil {
				return nil, fmt.Errorf("%w | %s", UnmarshallingError, err.Error()) // handleClientError(UnmarshallingError, err.Error())
			}

			return &accountsData, nil
		} else {

			apiHttpError := &ApiHttpError{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return nil, fmt.Errorf("%w | %s", UnmarshallingError, err.Error()) // handleClientError(UnmarshallingError, err.Error())
			}

			return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage) //handleBadStatusError(httpResponse.StatusCode, apiHttpError.ErrorMessage, httpResponse.Request.URL.Path)
		}

	} else {

		return nil, fmt.Errorf("%w | %s", ResponseReadError, err.Error()) //handleClientError(ResponseReadError, err.Error())
	}
}

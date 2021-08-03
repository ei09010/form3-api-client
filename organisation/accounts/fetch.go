package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/google/uuid"
)

func (c *Client) Fetch(accountId uuid.UUID) (*AccountData, error) {

	var err error

	c.BaseURL.Path = path.Join(c.BaseURL.Path, fmt.Sprintf("/%s", accountId.String()))

	c.BaseURL, err = c.BaseURL.Parse(c.BaseURL.Path)

	if err != nil {

		clientErr := handleClientError(FinalUrlParsingError, err.Error())

		return nil, clientErr
	}

	return c.getRequest()

}

func (c *Client) getRequest() (*AccountData, error) {

	customReq, err := http.NewRequest("GET", c.BaseURL.String(), nil)

	if err != nil {

		return nil, handleClientError(BuildingRequestError, err.Error())
	}

	// review headers
	customReq.Header.Set("content-encoding", "application/json; charset=utf-8")
	customReq.Header.Set("user-agent", "golang-sdk")

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {

		return nil, handleClientError(ExecutingRequestError, err.Error())
	}

	defer httpResponse.Body.Close()

	var accountsData AccountData

	if httpResponse != nil {

		responseBody, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return nil, handleClientError(ResponseReadError, err.Error())
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode == http.StatusOK {

			err = json.Unmarshal(responseBody, &accountsData)

			if err != nil {
				return nil, handleClientError(UnmarshallingError, err.Error())
			}

			return &accountsData, nil
		} else {

			apiHttpError := &ApiHttpError{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return nil, handleClientError(UnmarshallingError, err.Error())
			}

			return nil, handleBadStatusError(httpResponse.StatusCode, apiHttpError.ErrorMessage, httpResponse.Request.URL.Path)
		}

	} else {

		return nil, handleClientError(ResponseReadError, err.Error())
	}
}

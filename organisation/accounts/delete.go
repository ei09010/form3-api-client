package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/google/uuid"
)

func (c *Client) Delete(accountId uuid.UUID, version int) error {

	var err error

	c.BaseURL.Path = path.Join(c.BaseURL.Path, fmt.Sprintf("/%s", accountId.String()))

	c.BaseURL, err = c.BaseURL.Parse(c.BaseURL.Path)

	if err != nil {

		clientErr := handleClientError(FinalUrlParsingError, err.Error())

		return clientErr
	}

	return c.deleteRequest()

}

func (c *Client) deleteRequest() error {

	customReq, err := http.NewRequest("DELETE", c.BaseURL.String(), nil)

	if err != nil {

		return handleClientError(BuildingRequestError, err.Error())
	}

	// review headers
	customReq.Header.Set("content-encoding", "application/json")
	customReq.Header.Set("user-agent", "golang-sdk")

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {

		return handleClientError(ExecutingRequestError, err.Error())
	}

	defer httpResponse.Body.Close()

	if httpResponse != nil {

		responseBody, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return handleClientError(ResponseReadError, err.Error())
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode == http.StatusOK {

			return nil

		} else {

			apiHttpError := &ApiHttpError{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return handleClientError(UnmarshallingError, err.Error())
			}

			return handleBadStatusError(httpResponse.StatusCode, apiHttpError.ErrorMessage, httpResponse.Request.URL.Path)
		}

	} else {

		return handleClientError(ResponseReadError, err.Error())
	}
}

package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/google/uuid"
)

func (c *Client) Fetch(accountId uuid.UUID) (*AccountData, error) {

	var err error

	c.BaseURL.Path = path.Join(c.BaseURL.Path, fmt.Sprintf("/%s", accountId.String()))

	c.BaseURL, err = c.BaseURL.Parse(c.BaseURL.Path)

	if err != nil {

		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.BaseURL.Path, http.StatusBadRequest, err.Error())
	}

	return c.getRequest()

}

func (c *Client) getRequest() (*AccountData, error) {

	customReq, err := http.NewRequest(http.MethodGet, c.BaseURL.String(), nil)

	if err != nil {

		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.BaseURL.Path, http.StatusBadRequest, err.Error())
	}

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", TimeoutError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
			} else {
				return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ExecutingRequestError, httpResponse.Request.URL.Path, http.StatusBadRequest, err.Error())
			}
		}
	}

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

			apiHttpError := &ApiHttpError{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", UnmarshallingError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
			}

			return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
		}

	} else {

		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, c.BaseURL.String(), http.StatusInternalServerError, err.Error())
	}

}

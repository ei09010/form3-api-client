package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (c *Client) Delete(accountId uuid.UUID, version int) error {

	var err error

	c.BaseURL.Path = path.Join(c.BaseURL.Path, fmt.Sprintf("/%s", accountId.String()))

	c.BaseURL, err = c.BaseURL.Parse(c.BaseURL.Path)

	if err != nil {

		return fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.BaseURL.Path, http.StatusBadRequest, err.Error())
	}

	return c.deleteRequest(map[string]string{"version": strconv.Itoa(version)})

}

func (c *Client) deleteRequest(queryStringParam map[string]string) error {

	customReq, err := http.NewRequest("DELETE", c.BaseURL.String(), nil)

	q := customReq.URL.Query()

	for k, v := range queryStringParam {
		q.Add(k, v)
	}

	customReq.URL.RawQuery = q.Encode()

	if err != nil {

		return fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.BaseURL.Path, http.StatusBadRequest, err.Error())
	}

	customReq.Header.Set("user-agent", "golang-sdk")

	customReq.URL.Path = strings.Join([]string{customReq.URL.Path, customReq.URL.RawQuery}, "?")

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {

		return fmt.Errorf("%w | Path: %s returned %d with message %s", ExecutingRequestError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
	}

	defer httpResponse.Body.Close()

	if httpResponse != nil {

		responseBody, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode == http.StatusOK {

			return nil

		} else {

			apiHttpError := &ApiHttpError{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
			}

			return fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
		}

	} else {

		return fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
	}
}

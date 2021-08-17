package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (c *Client) Delete(accountId uuid.UUID, version int) error {

	return c.deleteJSON(accountId, map[string]string{"version": strconv.Itoa(version)}, accountsApi)

}
func (c *Client) deleteJSON(accountId uuid.UUID, queryStringParam map[string]string, config *apiConfig) error {
	httpResp, err := c.deleteRequest(accountId, queryStringParam, config)

	if err != nil {
		return fmt.Errorf("%w | %d | %s", ExecutingRequestError, httpResp.StatusCode, err)
	}

	if httpResp.Body == http.NoBody {
		if !isHttpCodeOK(httpResp.StatusCode) {
			return fmt.Errorf("%w | %d | %s", ApiHttpErrorType, httpResp.StatusCode, "")
		}
		return nil
	}

	resp := &apiErrorMessage{}

	err = json.NewDecoder(httpResp.Body).Decode(resp)

	defer httpResp.Body.Close()

	if err != nil {
		return fmt.Errorf("%w | %d | %s", UnmarshallingError, httpResp.StatusCode, err)
	}

	return fmt.Errorf("%w | %d | %s", ApiHttpErrorType, httpResp.StatusCode, resp.ErrorMessage)
}

func (c *Client) deleteRequest(accountId uuid.UUID, queryStringParam map[string]string, config *apiConfig) (*http.Response, error) {

	var err error

	c.baseURL, err = c.baseURL.Parse(config.path + "/" + accountId.String())

	if err != nil {
		return nil, err
	}

	customReq, err := http.NewRequest(http.MethodDelete, c.baseURL.String(), nil)

	if err != nil {
		return nil, err
	}
	urlQuery := customReq.URL.Query()

	for k, v := range queryStringParam {
		urlQuery.Add(k, v)
	}

	customReq.URL.RawQuery = urlQuery.Encode()

	customReq.URL.Path = strings.Join([]string{customReq.URL.Path, customReq.URL.RawQuery}, "?")

	return c.HttpClient.Do(customReq)
}

// func (c *Client) deleteRequest(queryStringParam map[string]string) error {

// 	customReq, err := http.NewRequest(http.MethodDelete, c.baseURL.String(), nil)

// 	q := customReq.URL.Query()

// 	for k, v := range queryStringParam {
// 		q.Add(k, v)
// 	}

// 	customReq.URL.RawQuery = q.Encode()

// 	if err != nil {

// 		return fmt.Errorf("%w | Path: %s returned %d with message %s", BuildingRequestError, c.baseURL.Path, http.StatusBadRequest, err.Error())
// 	}

// 	customReq.URL.Path = strings.Join([]string{customReq.URL.Path, customReq.URL.RawQuery}, "?")

// 	httpResponse, err := c.HttpClient.Do(customReq)

// 	if err != nil {

// 		return fmt.Errorf("%w | Path: %s returned %d with message %s", ExecutingRequestError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
// 	}

// 	defer httpResponse.Body.Close()

// 	if httpResponse != nil {

// 		responseBody, err := ioutil.ReadAll(httpResponse.Body)

// 		if err != nil {
// 			return fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
// 		}

// 		httpResponse.Body.Close()

// 		if httpResponse.StatusCode == http.StatusOK {

// 			return nil

// 		} else {

// 			apiHttpError := &apiErrorMessage{}

// 			err = json.Unmarshal(responseBody, apiHttpError)

// 			if err != nil {
// 				return fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
// 			}

// 			return fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
// 		}

// 	} else {

// 		return fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
// 	}
// }

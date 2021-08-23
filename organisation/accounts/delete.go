package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Delete issues an API request to delete a an account with a given accountId and version number
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

	customReq.Header.Set("Content-Type", "application/json")
	customReq.Header.Set("Accept", "*/*")
	customReq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	customReq.Header.Set("user-agent", "golang-sdk")

	return c.httpClient.Do(customReq)
}

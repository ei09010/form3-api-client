package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// Delete issues an API request to delete a an account with a given accountId and version number
func (c *Client) Delete(ctx context.Context, accountId uuid.UUID, version int) error {

	return c.deleteJSON(ctx, accountId, map[string]string{"version": strconv.Itoa(version)}, AccountsApiDefaultUrl)

}
func (c *Client) deleteJSON(ctx context.Context, accountId uuid.UUID, queryStringParam map[string]string, config *apiConfig) error {
	httpResp, err := c.deleteRequest(ctx, accountId, queryStringParam, config)

	if err != nil {
		if httpResp != nil {
			return fmt.Errorf("%w | %d | %s", BuildingRequestError, httpResp.StatusCode, err)
		} else {
			return fmt.Errorf("%w | %d | %s", BuildingRequestError, http.StatusBadRequest, err)
		}
	}

	if httpResp.Body == http.NoBody {
		if !isHttpCodeOK(httpResp.StatusCode) {
			return fmt.Errorf("%w | %d | %s", ApiHttpErrorType, httpResp.StatusCode, "")
		}
		return nil
	}

	resp := &apiCommonResult{}

	err = json.NewDecoder(httpResp.Body).Decode(resp)

	defer httpResp.Body.Close()

	if err != nil {
		return fmt.Errorf("%w | %d | %s", BuildingRequestError, httpResp.StatusCode, err)
	}

	return fmt.Errorf("%w | %d | %s", ApiHttpErrorType, httpResp.StatusCode, resp.ErrorMessage)
}

func (c *Client) deleteRequest(ctx context.Context, accountId uuid.UUID, queryStringParam map[string]string, config *apiConfig) (*http.Response, error) {

	var err error

	c.baseURL, err = c.baseURL.Parse(config.path + "/" + accountId.String())

	if err != nil {
		return nil, err
	}

	q, err := url.ParseQuery(c.baseURL.RawQuery)

	if err != nil {
		return nil, err
	}

	for k, v := range queryStringParam {
		q.Add(k, v)
	}

	c.baseURL.RawQuery = q.Encode()

	customReq, err := http.NewRequest(http.MethodDelete, c.baseURL.String(), nil)

	if err != nil {
		return nil, err
	}

	customReq.WithContext(ctx)

	addHeaders(customReq)

	return c.httpClient.Do(customReq)
}

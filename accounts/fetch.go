package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Fetch retrieves account related information using an accountId
func (c *Client) Fetch(ctx context.Context, accountId uuid.UUID) (*AccountResponse, error) {

	accountResponse := &AccountResponse{}

	if err := c.getJSON(ctx, accountId, AccountsApiDefaultUrl, accountResponse); err != nil {
		return nil, err
	}

	if err := accountResponse.Error(); err != nil {
		return nil, err
	}

	return accountResponse, nil
}

func (c *Client) getJSON(ctx context.Context, accountId uuid.UUID, config *apiConfig, resp *AccountResponse) error {

	httpResp, err := c.get(ctx, accountId, config)

	if err != nil {
		if httpResp != nil {
			return fmt.Errorf("%w | %d | %s", BuildingRequestError, httpResp.StatusCode, err)
		} else {
			return fmt.Errorf("%w | %d | %s", BuildingRequestError, http.StatusBadRequest, err)
		}
	}

	resp.Status = httpResp.StatusCode

	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(resp)

	if err != nil {
		return fmt.Errorf("%w | %d | %s", BuildingRequestError, httpResp.StatusCode, err)
	}

	return nil

}

func (c *Client) get(ctx context.Context, accountId uuid.UUID, config *apiConfig) (*http.Response, error) {

	var err error

	c.baseURL, err = c.baseURL.Parse(config.path + "/" + accountId.String())

	if err != nil {
		return nil, err
	}
	customReq, err := http.NewRequest(http.MethodGet, c.baseURL.String(), nil)

	if err != nil {
		return nil, err
	}

	customReq.WithContext(ctx)

	addHeaders(customReq)

	return c.httpClient.Do(customReq)

}

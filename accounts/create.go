package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Create issues an API request to store given account related information
func (c *Client) Create(ctx context.Context, accountData *AccountData) (*AccountResponse, error) {

	accountResponse := &AccountResponse{}

	if err := c.postJSON(ctx, AccountsApiDefaultUrl, accountData, accountResponse); err != nil {
		return nil, err
	}

	if err := accountResponse.Error(); err != nil {
		return nil, err
	}

	return accountResponse, nil

}

func (c *Client) postJSON(ctx context.Context, config *apiConfig, apiReqContent interface{}, resp *AccountResponse) error {
	httpResp, err := c.post(ctx, apiReqContent, config)

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

func (c *Client) post(ctx context.Context, apiReq interface{}, config *apiConfig) (*http.Response, error) {

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, err
	}

	if c.baseURL.Host == "" {

		c.baseURL, err = c.baseURL.Parse(config.host + config.path)

	} else {

		c.baseURL, err = c.baseURL.Parse(config.path)
	}

	if err != nil {
		return nil, err
	}

	customReq, err := http.NewRequest(http.MethodPost, c.baseURL.String(), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	addHeaders(customReq)

	customReq.WithContext(ctx)

	return c.httpClient.Do(customReq)
}

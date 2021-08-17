package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var accountsApi = &apiConfig{
	host: "https://api.form3.tech",
	path: "/v1/organisation/accounts",
}

func (c *Client) Fetch(accountId uuid.UUID) (*AccountResponse, error) {

	accountResponse := &AccountResponse{}

	if err := c.getJSON(accountId, accountsApi, accountResponse); err != nil {
		return nil, err
	}

	if err := accountResponse.Error(); err != nil {
		return nil, err
	}

	return accountResponse, nil
}

func (c *Client) getJSON(accountId uuid.UUID, config *apiConfig, resp *AccountResponse) error {

	httpResp, err := c.get(accountId, config)

	if err != nil {
		return fmt.Errorf("%w | %d | %s", ExecutingRequestError, httpResp.StatusCode, err)
	}

	resp.apiErrorMessage.Status = httpResp.StatusCode

	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(resp)

	if err != nil {
		return fmt.Errorf("%w | %d | %s", UnmarshallingError, httpResp.StatusCode, err)
	}

	return nil

}

func (c *Client) get(accountId uuid.UUID, config *apiConfig) (*http.Response, error) {

	var err error

	c.baseURL, err = c.baseURL.Parse(config.path)

	if err != nil {
		return nil, err
	}
	customReq, err := http.NewRequest(http.MethodGet, c.baseURL.String(), nil)

	if err != nil {
		return nil, err
	}

	return c.HttpClient.Do(customReq)

}

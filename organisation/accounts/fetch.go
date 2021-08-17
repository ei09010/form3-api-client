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

	host := config.host

	if c.baseURL.Host != "" {
		host = c.baseURL.Host
	}

	customReq, err := http.NewRequest(http.MethodGet, host+config.path+accountId.String(), nil)

	if err != nil {
		return nil, err
	}

	return c.HttpClient.Do(customReq)

	// if err != nil {
	// 	return nil, err
	// }

	// var accountsData AccountData

	// if httpResponse != nil {

	// 	responseBody, err := ioutil.ReadAll(httpResponse.Body)

	// 	if err != nil {

	// 		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())

	// 	}

	// 	httpResponse.Body.Close()

	// 	if httpResponse.StatusCode == http.StatusOK {

	// 		err = json.Unmarshal(responseBody, &accountsData)

	// 		if err != nil {
	// 			return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", UnmarshallingError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
	// 		}

	// 		return &accountsData, nil
	// 	} else {

	// 		apiHttpError := &ApiHttpError{}

	// 		err = json.Unmarshal(responseBody, apiHttpError)

	// 		if err != nil {
	// 			return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", UnmarshallingError, httpResponse.Request.URL.Path, httpResponse.StatusCode, err.Error())
	// 		}

	// 		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ApiHttpErrorType, httpResponse.Request.URL.Path, httpResponse.StatusCode, apiHttpError.ErrorMessage)
	// 	}

	// } else {

	// 	return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", ResponseReadError, c.baseURL.String(), http.StatusInternalServerError, err.Error())
	// }

}

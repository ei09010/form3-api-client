package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	BaseURL    *url.URL
	HttpClient *http.Client
}

func NewClient(baseUrl string, requestTimeout time.Duration) (*Client, error) {

	parsedUrl, err := url.Parse(baseUrl)

	if err != nil {
		return nil, handleError(UrlParsingError, err.Error(), 0)
	}

	finalUrl, err := parsedUrl.Parse("/v1/organisation/accounts")

	if err != nil {
		return nil, handleError(UrlParsingError, err.Error(), 0)
	}

	return &Client{
		BaseURL:    finalUrl,
		HttpClient: &http.Client{Timeout: requestTimeout},
	}, nil
}

func handleError(errorType int, err string, httpErrorCode int) *ClientError {

	clientErr := &ClientError{
		ErrorType:    errorType,
		ErrorMessage: err,
	}

	if httpErrorCode != 0 {
		clientErr.HttpCode = httpErrorCode
	}

	log.Println(clientErr)

	return clientErr
}

func (c *Client) Fetch(accountId uuid.UUID) (*AccountData, error) {

	var err error

	c.BaseURL.Path = path.Join(c.BaseURL.Path, fmt.Sprintf("/%s", accountId.String()))

	c.BaseURL, err = c.BaseURL.Parse(c.BaseURL.Path)

	if err != nil {

		clientErr := handleError(UrlParsingError, err.Error(), 0)

		log.Fatal(clientErr)

		return nil, clientErr
	}

	return c.getRequest("")

}

func (c *Client) getRequest(queryString string) (*AccountData, error) {

	customReq, err := http.NewRequest("GET", c.BaseURL.String(), nil)

	if err != nil {

		return nil, handleError(RequestError, err.Error(), 0)
	}

	// review headers
	customReq.Header.Set("content-encoding", "application/json; charset=utf-8")
	customReq.Header.Set("user-agent", "golang-sdk")

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {

		return nil, handleError(RequestError, err.Error(), 0)
	}

	defer httpResponse.Body.Close()

	var accountsData AccountData

	if httpResponse != nil {

		responseBody, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return nil, handleError(ResponseError, err.Error(), 0)
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode == http.StatusOK {

			err = json.Unmarshal(responseBody, &accountsData)

			if err != nil {
				handleError(UnmarshallingError, err.Error(), 0)
			}

			return &accountsData, nil
		} else {

			apiHttpError := &ApiHttpError{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return nil, handleError(UnmarshallingError, err.Error(), httpResponse.StatusCode)
			}

			return nil, handleError(HttResponseStandardError, apiHttpError.ErrorMessage, httpResponse.StatusCode)
		}

	} else {

		return nil, nil
	}
}

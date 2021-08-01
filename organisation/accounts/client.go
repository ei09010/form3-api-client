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

	parsedUrl, err := url.ParseRequestURI(baseUrl)

	if err != nil {
		return nil, handleClientError(BaseUrlParsingError, err.Error())
	}

	if requestTimeout.Milliseconds() <= 0 {
		requestTimeout = DefaultTimeOutValue
	}

	finalUrl, err := parsedUrl.Parse(AccountsPath)

	if err != nil {
		return nil, handleClientError(PathParsingError, err.Error())
	}

	return &Client{
		BaseURL:    finalUrl,
		HttpClient: &http.Client{Timeout: requestTimeout},
	}, nil
}

func handleClientError(errorType int, err string) *ClientError {

	clientErr := &ClientError{
		ErrorType:    errorType,
		ErrorMessage: err,
	}

	log.Println(clientErr)

	return clientErr
}

func handleBadStatusError(httpCode int, errorMessage string, url string) *ClientError {

	clientErr := &ClientError{
		ErrorType:    HttResponseStandardError,
		ErrorMessage: errorMessage,
		BadStatusError: &BadStatusError{
			HttpCode: httpCode,
			URL:      url,
		},
	}

	log.Println(clientErr)

	return clientErr
}

func (c *Client) Fetch(accountId uuid.UUID) (*AccountData, error) {

	var err error

	c.BaseURL.Path = path.Join(c.BaseURL.Path, fmt.Sprintf("/%s", accountId.String()))

	c.BaseURL, err = c.BaseURL.Parse(c.BaseURL.Path)

	if err != nil {

		clientErr := handleClientError(FinalUrlParsingError, err.Error())

		return nil, clientErr
	}

	return c.getRequest()

}

func (c *Client) getRequest() (*AccountData, error) {

	customReq, err := http.NewRequest("GET", c.BaseURL.String(), nil)

	if err != nil {

		return nil, handleClientError(BuildingRequestError, err.Error())
	}

	// review headers
	customReq.Header.Set("content-encoding", "application/json; charset=utf-8")
	customReq.Header.Set("user-agent", "golang-sdk")

	httpResponse, err := c.HttpClient.Do(customReq)

	if err != nil {

		return nil, handleClientError(ExecutingRequestError, err.Error())
	}

	defer httpResponse.Body.Close()

	var accountsData AccountData

	if httpResponse != nil {

		responseBody, err := ioutil.ReadAll(httpResponse.Body)

		if err != nil {
			return nil, handleClientError(ResponseError, err.Error())
		}

		httpResponse.Body.Close()

		if httpResponse.StatusCode == http.StatusOK {

			err = json.Unmarshal(responseBody, &accountsData)

			if err != nil {
				handleClientError(UnmarshallingError, err.Error())
			}

			return &accountsData, nil
		} else {

			apiHttpError := &ApiHttpError{}

			err = json.Unmarshal(responseBody, apiHttpError)

			if err != nil {
				return nil, handleClientError(UnmarshallingError, err.Error())
			}

			return nil, handleBadStatusError(HttResponseStandardError, apiHttpError.ErrorMessage, httpResponse.Request.URL.String())
		}

	} else {

		return nil, nil
	}
}

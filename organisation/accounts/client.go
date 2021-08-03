package accounts

import (
	"log"
	"net/http"
	"net/url"
	"time"
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

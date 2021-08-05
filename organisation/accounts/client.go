package accounts

import (
	"fmt"
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
		return nil, fmt.Errorf("%w | %s", BaseUrlParsingError, err.Error())
	}

	if requestTimeout.Milliseconds() <= 0 {
		requestTimeout = DefaultTimeOutValue
	}

	finalUrl, err := parsedUrl.Parse(AccountsPath)

	if err != nil {
		return nil, fmt.Errorf("%w | %s", PathParsingError, err.Error())
	}

	return &Client{
		BaseURL:    finalUrl,
		HttpClient: &http.Client{Timeout: requestTimeout},
	}, nil
}

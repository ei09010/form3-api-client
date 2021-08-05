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
		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", BaseUrlParsingError, baseUrl, http.StatusBadRequest, err.Error())
	}

	if requestTimeout.Milliseconds() <= 0 {
		requestTimeout = DefaultTimeOutValue
	}

	finalUrl, err := parsedUrl.Parse(AccountsPath)

	if err != nil {
		return nil, fmt.Errorf("%w | Path: %s returned %d with message %s", PathParsingError, AccountsPath, http.StatusBadRequest, err.Error())
	}

	return &Client{
		BaseURL:    finalUrl,
		HttpClient: &http.Client{Timeout: requestTimeout},
	}, nil
}

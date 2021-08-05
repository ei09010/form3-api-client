package accounts

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	Timeout    time.Duration
	HttpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
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
		Timeout:    requestTimeout,
		HttpClient: &http.Client{Timeout: requestTimeout},
	}, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	httpClient := c.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: c.Timeout}
	}

	req.Header.Set("content-encoding", "application/json; charset=utf-8")
	req.Header.Set("user-agent", "golang-sdk")

	return httpClient.Do(req)
}

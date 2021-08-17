package accounts

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type apiConfig struct {
	host string
	path string
}

type ClientOption func(*Client) error

type Client struct {
	baseURL    *url.URL
	timeout    time.Duration
	HttpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

func NewClient(clientOptions ...ClientOption) (*Client, error) {

	c := &Client{
		HttpClient: &http.Client{},
	}

	for _, option := range clientOptions {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// WithTimeout sets the client requests timeout with a custom value
func WithTimeout(customRequestTimeout time.Duration) ClientOption {
	return func(c *Client) error {

		if customRequestTimeout.Milliseconds() <= 0 {
			customRequestTimeout = DefaultTimeOutValue
		}

		c.timeout = customRequestTimeout

		return nil
	}
}

// WithBaseURL sets the client with a custom base url
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		parsedUrl, err := url.ParseRequestURI(baseURL)

		if err != nil {
			return fmt.Errorf("%w | %s", clientCreationError, err.Error())
		}

		c.baseURL = parsedUrl
		return nil
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	httpClient := c.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	req.Header.Set("content-encoding", "application/json; charset=utf-8")
	req.Header.Set("user-agent", "golang-sdk")

	return httpClient.Do(req)
}

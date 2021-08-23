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

var accountsApi = &apiConfig{
	host: "https://api.form3.tech",
	path: "/v1/organisation/accounts",
}

// ClientOption is the type of constructor options for NewClient(...)
type ClientOption func(*Client) error

// Client may be used to make requests to the Form3 API
type Client struct {
	baseURL    *url.URL
	timeout    time.Duration
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

// NewClient constructs a new Client which can make requests to the Form3 API
func NewClient(clientOptions ...ClientOption) (*Client, error) {

	c := &Client{
		httpClient: &http.Client{},
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
			return fmt.Errorf("%w | %d | %s", clientCreationError, http.StatusBadRequest, err.Error())
		}

		c.baseURL = parsedUrl
		return nil
	}
}

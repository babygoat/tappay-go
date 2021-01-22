package tappay

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// APIURL is the url of TapPay server for transaction in production server.
	APIURL string = "https://prod.tappaysdk.com/"
)

type client struct {
	partnerKey string
	httpClient *http.Client

	// url is the base URL to use for API paths.
	url string
}

type clientOption func(*client)

// NewClient creates a new TapPay client for transaction
func NewClient(key string, options ...clientOption) (*client, error) {
	var url string
	if url = os.Getenv("TAPPAY_SERVER"); url == "" {
		url = APIURL
	}

	// defaultTimeout is the default timeout on the http.Client used the by the library
	// This is chosen according to the document
	// https://docs.tappaysdk.com/tutorial/zh/back.html
	const defaultTimeout = 30 * time.Second
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	cli := &client{
		partnerKey: key,
		httpClient: httpClient,
		url:        url,
	}

	for _, option := range options {
		option(cli)
	}
	u, err := sanitizeURL(cli.url)
	if err != nil {
		return nil, fmt.Errorf("supplied server %q is not valid: %v", cli.url, err)
	}
	cli.url = u.String()
	return cli, nil
}

// WithServer returns a clientOption to override default base url to be used for the request
func WithServer(server string) clientOption {
	return func(c *client) {
		c.url = server
	}
}

func sanitizeURL(server string) (*url.URL, error) {
	u, err := url.ParseRequestURI(server)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// WithHTTPClient returns a clientOption to override the default http client
func WithHTTPClient(hClient *http.Client) clientOption {
	return func(c *client) {
		c.httpClient = hClient
	}
}

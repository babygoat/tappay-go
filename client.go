package tappay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// SandboxAPIURL is the url of TapPay server for testing.
	SandboxAPIURL string = "https://sandbox.tappaysdk.com/"

	// APIURL is the url of TapPay server for transaction in production server.
	APIURL string = "https://prod.tappaysdk.com/"
)

// type service denotes the the operations provided by TapPay
type service string

const (
	servicePayByPrime service = "pay_by_prime"
	serviceRecord     service = "record"
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

// do is used to issue the http request with client to TapPay server and parse the http.Response
func (c *client) do(req *http.Request) ([]byte, error) {
	rawResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rawResp.Body.Close()

	var b []byte
	buf := bytes.NewBuffer(b)
	if _, err = io.Copy(buf, rawResp.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// newRequest is used to create the http request with the input. Also, appends the common header like
// `content-type`, `x-api-key` and injects the common field `partner_key` into request body.
func (c *client) newRequest(ctx context.Context, method string, svc service, input Marshaler) (*http.Request, error) {
	paramsMap, err := input.MarshalMap()
	if err != nil {
		return nil, err
	}

	var svcPath string
	switch svc {
	case servicePayByPrime:
		svcPath = payByPrimePath
	case serviceRecord:
		svcPath = recordPath
	}

	u, _ := url.Parse(svcPath)
	base, _ := url.Parse(c.url)
	path := base.ResolveReference(u).String()

	paramsMap["partner_key"] = c.partnerKey
	body, _ := json.Marshal(paramsMap)
	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("cannot create a TapPay request: %v", err)
	}
	req.Header.Add("x-api-key", c.partnerKey)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// Marshaler is the interface implemented by the types
// that would marshal themselves into valid map
type Marshaler interface {
	MarshalMap() (map[string]interface{}, error)
}

package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	pathpkg "path"
	"time"

	"github.com/spirifoxy/accountclient/internal/f3"
	"golang.org/x/time/rate"
)

// Client is a wrapper around http.Client with optional timeout
// and rate limiter settings used for accessing accounts API
type Client struct {
	httpClient *http.Client
	apiURL     *url.URL

	basePath    string
	rateLimiter *rate.Limiter
	reqTimeout  time.Duration
}

// Option is function used for applying configurations to client
type Option func(*Client)

// New creates Client with api url and options
func New(baseURL string, options ...Option) (*Client, error) {
	client := &Client{
		httpClient: http.DefaultClient,
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, &ClientInitError{
			Err: fmt.Errorf("not able to parse provided url: %v", err),
		}
	}

	for _, o := range options {
		o(client)
	}

	if client.basePath != "" {
		u.Path = pathpkg.Join(u.Path, client.basePath)
	}
	client.apiURL = u

	return client, nil
}

// WithBasePath might be used for specifying API version.
func WithBasePath(p string) Option {
	return func(c *Client) {
		c.basePath = p
	}
}

// WithRateLimiter allows the client to send up to n requests per s seconds.
func WithRateLimiter(s time.Duration, n int) Option {
	return func(c *Client) {
		c.rateLimiter = rate.NewLimiter(rate.Every(s*time.Second), n)
	}
}

// WithTimeout sets up the timeout for every request client performs.
func WithTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.reqTimeout = t
	}
}

// do is basically a wrapper around httpClient.Do with timeout and rate limiter.
// Sends an HTTP request and returns an HTTP response.
//
// Error is returned in case if either http client returns error
// or requests limits (if set up) exceeded.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	if c.reqTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.reqTimeout)
		defer cancel()
	}

	if c.rateLimiter != nil {
		err := c.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) urlForPath(p string) string {
	return fmt.Sprintf("%s%s", c.apiURL.String(), p)
}

func (c *Client) get(path string, params f3.HTTPParams) (*http.Response, error) {
	req, err := f3.NewRequest("GET", c.urlForPath(path), params, nil, f3.DefaultHeaders())
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) post(path string, body []byte) (*http.Response, error) {
	req, err := f3.NewRequest("POST", c.urlForPath(path), nil, body, f3.ContentHeaders())
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) delete(path string, params f3.HTTPParams) (*http.Response, error) {
	req, err := f3.NewRequest("DELETE", c.urlForPath(path), params, nil, f3.DefaultHeaders())
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

// parseErrorResponse tries to parse body for the response with
// not appropriate status code.
//
// RequestStatusError is returned even if the attempt to parse the
// response fails
func parseErrorResponse(body io.ReadCloser, sc int) error {
	var errorResp *ErrorResponse
	err := json.NewDecoder(body).Decode(&errorResp)
	if err != nil {
		return &RequestStatusError{
			Code: sc,
			Err:  fmt.Errorf("not able to parse error info: %w", err),
		}
	}

	return &RequestStatusError{
		Code: sc,
		Err:  errors.New(errorResp.Message),
	}
}

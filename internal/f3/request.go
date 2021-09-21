package f3

import (
	"bytes"
	"net/http"
	"time"
)

// AccountsEndpoint allows sending POST, GET, DELETE requests
const AccountsEndpoint = "/organisation/accounts"

// HealthEndpoint allows sending GET requests
const HealthEndpoint = "/health"

// HTTPParams is map used for url paramaters
type HTTPParams map[string]string

// Headers is set of headers used for quering F3 api
type Headers map[string]string

// NewRequest is a wrapper over http NewRequest function.
// Adds passed headers and query parameters to the request.
func NewRequest(method, url string, params HTTPParams, body []byte, h Headers) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func timeToRFC1123(t time.Time) string {
	return t.Format(time.RFC1123)
}

// DefaultHeaders is a set of headers required for every f3 request.
func DefaultHeaders() Headers {
	return map[string]string{
		"Date":   timeToRFC1123(time.Now()),
		"Accept": "application/vnd.api+json",
	}
}

// ContentHeaders is a set of headers required for f3 requests that contain body.
func ContentHeaders() Headers {
	h := make(map[string]string)
	for k, v := range DefaultHeaders() {
		h[k] = v
	}
	h["Content-Type"] = "application/vnd.api+json"
	return h
}

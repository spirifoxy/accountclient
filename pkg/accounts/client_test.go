package accounts

import (
	"net/http"
	"os"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const errorValidationJSON = `
		{
			"error_message": "validation failure"
		}
	`

type MockHandler func(req *http.Request) (*http.Response, error)

type MockClient struct {
	DoFunc MockHandler
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{}, nil
}

// getEnv tries to get environment varialbe,
// returns fallback in case of failure
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func baseURL() string {
	return getEnv("API_ADDR", "http://localhost:8080")
}

func getMockClient(h func() MockHandler) *Client {
	c, _ := New(baseURL())
	c.httpClient = &MockClient{DoFunc: h()}
	return c
}

func buildAccount(attr *AccountAttributes) *Account {
	acc := NewAccount(uuid.NewV4().String(), "GB", []string{"Samantha Holder"}, attr)
	return acc
}

func TestWithBasePath(t *testing.T) {
	client, err := New(
		baseURL(),
		WithBasePath("v1"),
	)

	require.Nil(t, err)
	assert.Equal(t, "v1", client.basePath)
}

func TestWithRateLimiter(t *testing.T) {
	client, err := New(
		baseURL(),
		WithRateLimiter(time.Duration(5), 10),
	)

	require.Nil(t, err)
	require.NotNil(t, client.rateLimiter)
}

func TestWithTimeout(t *testing.T) {
	client, err := New(
		baseURL(),
		WithTimeout(time.Duration(5)*time.Second),
	)

	require.Nil(t, err)
	assert.Equal(t, time.Duration(5)*time.Second, client.reqTimeout)
}

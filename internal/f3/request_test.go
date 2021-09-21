package f3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	url        string     = "localhost"
	parameters HTTPParams = HTTPParams{
		"param1": "val1",
		"param2": "val2",
	}
)

func TestTimeToRFC1123(t *testing.T) {
	now, err := time.Parse(time.RFC822, "02 Jan 06 15:04 MST")
	require.NoError(t, err)

	assert.Equal(t, "Mon, 02 Jan 2006 15:04:00 MST", timeToRFC1123(now))
}

func TestDefaultHeaders(t *testing.T) {
	h := DefaultHeaders()

	assert.NotEmpty(t, h["Accept"])
	assert.NotEmpty(t, h["Date"])

	assert.Equal(t, "application/vnd.api+json", h["Accept"])
}

func TestContentHeaders(t *testing.T) {
	h := ContentHeaders()

	assert.NotEmpty(t, h["Accept"])
	assert.NotEmpty(t, h["Date"])
	assert.NotEmpty(t, h["Content-Type"])

	assert.Equal(t, "application/vnd.api+json", h["Content-Type"])
	assert.Equal(t, "application/vnd.api+json", h["Accept"])
}

func TestNewRequestWithDefaultHeaders(t *testing.T) {
	req, err := NewRequest("GET", url, nil, nil, DefaultHeaders())
	require.NoError(t, err)

	assert.NotEmpty(t, req.Header.Get("Accept"))
	assert.NotEmpty(t, req.Header.Get("Date"))

	assert.Equal(t, "GET", req.Method)
}

func TestNewRequestWithContentHeaders(t *testing.T) {
	req, err := NewRequest("POST", url, nil, nil, ContentHeaders())
	require.NoError(t, err)

	assert.NotEmpty(t, req.Header.Get("Accept"))
	assert.NotEmpty(t, req.Header.Get("Date"))
	assert.NotEmpty(t, req.Header.Get("Content-Type"))

	assert.Equal(t, "POST", req.Method)
}

func TestNewRequestWithParameters(t *testing.T) {
	req, err := NewRequest("GET", url, parameters, nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "val1", req.URL.Query().Get("param1"))
	assert.Equal(t, "val2", req.URL.Query().Get("param2"))
}

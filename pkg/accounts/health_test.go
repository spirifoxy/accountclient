package accounts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const statusHealthyJSON = `
		{
			"status": "up"
		}
	`

const errorHealthyJSON = `
	{
		"status": "down"
	}
`

func healthReqAliveHandler() MockHandler {
	return func(req *http.Request) (*http.Response, error) {
		statusCode := 200
		jsonResponse := statusHealthyJSON

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(jsonResponse)),
			StatusCode: statusCode,
		}, nil
	}
}

func TestIsHealthy(t *testing.T) {
	type test struct {
		inHeader string
		outBool  bool
		outErr   error
	}

	client := getMockClient(healthReqAliveHandler)

	isHealthy, err := client.IsHealthy()
	assert.Equal(t, true, isHealthy)
	assert.Equal(t, nil, err)
}

func healthReqBrokenHandler() MockHandler {
	return func(req *http.Request) (*http.Response, error) {
		statusCode := 400
		jsonResponse := errorHealthyJSON

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(jsonResponse)),
			StatusCode: statusCode,
		}, nil
	}
}

func TestIsHealthyBroken(t *testing.T) {
	client := getMockClient(healthReqBrokenHandler)

	expect := &RequestStatusError{
		Code: 400,
		Err:  fmt.Errorf("health returns unexpected status code"),
	}

	isHealthy, err := client.IsHealthy()
	assert.Equal(t, false, isHealthy)
	assert.Equal(t, expect, err)
}

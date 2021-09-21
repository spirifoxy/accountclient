package accounts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

const errorVersionJSON = `
		{
			"error_message": "invalid version"
		}
	`

func deleteReqHandler() MockHandler {
	return func(req *http.Request) (*http.Response, error) {
		var statusCode int
		var jsonResponse string

		versionParam, ok := req.URL.Query()["version"]
		if !ok || len(versionParam[0]) < 1 {
			return &http.Response{}, nil
		}
		version, err := strconv.Atoi(versionParam[0])
		if err != nil {
			return &http.Response{}, err
		}

		switch version {
		case 1:
			statusCode = 404
		case 2:
			statusCode = 409
			jsonResponse = errorVersionJSON
		default:
			statusCode = 204
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(jsonResponse)),
			StatusCode: statusCode,
		}, nil
	}
}

func TestDelete(t *testing.T) {
	type test struct {
		inID  string
		inVer int
		out   error
	}

	client := getMockClient(deleteReqHandler)

	tests := []test{
		{
			inID:  "11111111-1111-1111-1111-111111111111",
			inVer: 0,
			out:   nil,
		},
		{
			inID:  "11111111-1111-1111-1111-111111111111",
			inVer: 1,
			out: &RequestStatusError{
				Code: 404,
				Err:  fmt.Errorf("rejected by server"),
			},
		},
		{
			inID:  "11111111-1111-1111-1111-111111111111",
			inVer: 2,
			out: &RequestStatusError{
				Code: 409,
				Err:  fmt.Errorf("invalid version"),
			},
		},
	}

	for _, tc := range tests {
		u, _ := uuid.FromString(tc.inID)
		err := client.Delete(u, tc.inVer)
		assert.Equal(t, tc.out, err)
	}
}

package accounts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

const foundAccountJSON = `
		{
			"data": {
				"attributes": {
					"alternative_names": null,
					"country": "GB",
					"name": [
						"Samantha Holder"
					]
				},
				"created_on": "2021-01-01T11:22:33.724Z",
				"id": "11111111-1111-1111-1111-111111111111",
				"modified_on": "2021-01-01T11:22:33.724Z",
				"organisation_id": "11111111-1111-1111-1111-111111111111",
				"type": "accounts",
				"version": 0
			},
			"links": {
				"self": "/v1/organisation/accounts/7d2fb82f-33c7-42a0-9174-a054991eb0ca"
			}
		}
	`

const errorNotFoundJSON = `
		{
			"error_message": "record does not exist"
		}
	`

func fetchReqHandler() MockHandler {
	return func(req *http.Request) (*http.Response, error) {
		var statusCode int
		var jsonResponse string

		u := path.Base(req.URL.Path)

		switch {
		case strings.HasPrefix(u, "222"):
			statusCode = 404
			jsonResponse = errorNotFoundJSON
		default:
			statusCode = 200
			jsonResponse = foundAccountJSON

		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(jsonResponse)),
			StatusCode: statusCode,
		}, nil
	}
}

func TestFetch(t *testing.T) {
	type test struct {
		inID   string
		outAcc *Account
		outErr error
	}

	client := getMockClient(fetchReqHandler)

	tests := []test{
		{
			inID: "11111111-1111-1111-1111-111111111111",
			outAcc: (func() *Account {
				acc := buildAccount(nil)
				acc.Data.ID = "11111111-1111-1111-1111-111111111111"
				acc.Data.OrganisationID = "11111111-1111-1111-1111-111111111111"
				return acc
			})(),
			outErr: nil,
		},
		{
			inID:   "22222222-1111-1111-1111-111111111111",
			outAcc: nil,
			outErr: &RequestStatusError{
				Code: 404,
				Err:  fmt.Errorf("record does not exist"),
			},
		},
	}

	for _, tc := range tests {
		u, _ := uuid.FromString(tc.inID)
		acc, err := client.Fetch(u)
		assert.Equal(t, tc.outAcc, acc)
		assert.Equal(t, tc.outErr, err)
	}
}

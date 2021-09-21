package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const errorDuplicateJSON = `
		{
			"error_message": "Account cannot be created as it violates a duplicate constraint"
		}
	`

func createReqHandler() MockHandler {
	return func(req *http.Request) (*http.Response, error) {
		var statusCode int
		var jsonResponse string

		var reqAcc *Account
		json.NewDecoder(req.Body).Decode(&reqAcc)

		switch reqAcc.Data.Attributes.BankID {
		case "corrupted":
			statusCode = 400
			jsonResponse = errorValidationJSON
		case "duplicated":
			statusCode = 409
			jsonResponse = errorDuplicateJSON
		case "unexpected":
			statusCode = 500
		default:
			reqAcc.Data.Version = 0
			jsonBytes, _ := json.Marshal(reqAcc)

			statusCode = 201
			jsonResponse = string(jsonBytes)
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(jsonResponse)),
			StatusCode: statusCode,
		}, nil
	}
}

func TestCreate(t *testing.T) {
	type test struct {
		inAcc  *Account
		outAcc *Account
		outErr error
	}

	client := getMockClient(createReqHandler)
	baseAcc := buildAccount(nil)

	tests := []test{
		{
			inAcc:  buildAccount(&AccountAttributes{BankID: "corrupted"}),
			outAcc: nil,
			outErr: &RequestStatusError{
				Code: 400,
				Err:  fmt.Errorf("validation failure"),
			},
		},
		{
			inAcc:  buildAccount(&AccountAttributes{BankID: "duplicated"}),
			outAcc: nil,
			outErr: &RequestStatusError{
				Code: 409,
				Err:  fmt.Errorf("Account cannot be created as it violates a duplicate constraint"),
			},
		},
		{
			inAcc:  buildAccount(&AccountAttributes{BankID: "unexpected"}),
			outAcc: nil,
			outErr: &RequestStatusError{
				Code: 500,
				Err:  fmt.Errorf("rejected by server"),
			},
		},
		{
			inAcc: baseAcc,
			outAcc: (func() *Account {
				baseAcc.Data.Version = 0
				return baseAcc
			})(),
			outErr: nil,
		},
	}

	for _, tc := range tests {
		acc, err := client.Create(tc.inAcc)
		assert.Equal(t, tc.outAcc, acc)
		assert.Equal(t, tc.outErr, err)
	}
}

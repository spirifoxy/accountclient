package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/spirifoxy/accountclient/internal/f3"
)

// Create creates a new bank account.
//
// Error is returned in case of failure.
func (c *Client) Create(accData *Account) (*Account, error) {
	body, err := json.Marshal(accData)
	if err != nil {
		return nil, err
	}

	resp, err := c.post(f3.AccountsEndpoint, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, parseErrorResponse(resp.Body, resp.StatusCode)
	}

	var acc *Account
	if err = json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}

	return acc, nil
}

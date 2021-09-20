package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/spirifoxy/accountclient/internal/f3"
)

// Fetch gets a single account using the account uuid.
//
// Error is returned in case of failure.
func (c *Client) Fetch(id uuid.UUID) (*Account, error) {
	path := fmt.Sprintf("%s/%s", f3.AccountsEndpoint, id)
	resp, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(resp.Body, resp.StatusCode)
	}

	var acc *Account
	if err = json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}

	return acc, nil
}

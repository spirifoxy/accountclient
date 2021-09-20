package accounts

import (
	"fmt"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/spirifoxy/accountclient/internal/f3"
)

// Delete deletes an account using the account uuid and version.
//
// Error is returned in case of failure.
func (c *Client) Delete(id uuid.UUID, version int) error {
	path := fmt.Sprintf("%s/%s", f3.AccountsEndpoint, id)
	params := map[string]string{"version": strconv.Itoa(version)}
	resp, err := c.delete(path, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return parseErrorResponse(resp.Body, resp.StatusCode)
	}

	return nil
}

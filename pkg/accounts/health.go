package accounts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spirifoxy/accountclient/internal/f3"
)

// IsHealthy checks if API server is functional based on health endpoint.
func (c *Client) IsHealthy() (bool, error) {
	resp, err := c.get(f3.HealthEndpoint, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, &RequestStatusError{
			Code: resp.StatusCode,
			Err:  errors.New("health returns unexpected status code"),
		}
	}

	var h *Health
	if err = json.NewDecoder(resp.Body).Decode(&h); err != nil {
		return false, err
	}

	return h.Status == "up", nil
}

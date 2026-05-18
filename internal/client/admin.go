package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) adminGET(path string, dest any) error {
	endpoint := adminURL + path

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	setBookingsCORSHeaders(req)
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.do(req, bookingsRoot)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := c.readBody(resp)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("decode %s: %w\nbody: %s", path, err, string(body))
	}
	return nil
}

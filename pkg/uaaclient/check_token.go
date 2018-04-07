package uaaclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (c *UAAClient) CheckToken(token string, ctx context.Context) (*User, error) {
	req, err := c.client.PostRequest("/check_token", strings.NewReader("token="+token))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.client_id, c.client_secret)

	resp, err := c.client.Do(req, ctx)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

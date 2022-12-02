package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type CreateAPIKeyRequest struct {
	Name          string `json:"name"`
	Role          string `json:"role"`
	SecondsToLive int64  `json:"secondsToLive,omitempty"`
}

type CreateAPIKeyResponse struct {
	// ID field only returned after Grafana v7.
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type GetAPIKeysResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Expiration time.Time `json:"expiration,omitempty"`
}

type DeleteAPIKeyResponse struct {
	Message string `json:"message"`
}

// CreateAPIKey creates a new Grafana API key.
func (c *Client) CreateAPIKey(request CreateAPIKeyRequest) (CreateAPIKeyResponse, error) {
	return c.CreateAPIKeyContext(context.Background(), request)
}

// CreateAPIKeyContext does the same thing as CreateAPIKey(), but also takes in a context.
func (c *Client) CreateAPIKeyContext(ctx context.Context, request CreateAPIKeyRequest) (CreateAPIKeyResponse, error) {
	response := CreateAPIKeyResponse{}

	data, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	err = c.request(ctx, "POST", "/api/auth/keys", nil, bytes.NewBuffer(data), &response)
	return response, err
}

// GetAPIKeys retrieves a list of all API keys.
func (c *Client) GetAPIKeys(includeExpired bool) ([]*GetAPIKeysResponse, error) {
	return c.GetAPIKeysContext(context.Background(), includeExpired)
}

// GetAPIKeysContext does the same thing as GetAPIKeys(), but also takes in a context.
func (c *Client) GetAPIKeysContext(ctx context.Context, includeExpired bool) ([]*GetAPIKeysResponse, error) {
	response := make([]*GetAPIKeysResponse, 0)

	query := url.Values{}
	query.Add("includeExpired", strconv.FormatBool(includeExpired))

	err := c.request(ctx, "GET", "/api/auth/keys", query, nil, &response)
	return response, err
}

// DeleteAPIKey deletes the Grafana API key with the specified ID.
func (c *Client) DeleteAPIKey(id int64) (DeleteAPIKeyResponse, error) {
	return c.DeleteAPIKeyContext(context.Background(), id)
}

// DeleteAPIKeyContext does the same thing as DeleteAPIKey(), but also takes in a context.
func (c *Client) DeleteAPIKeyContext(ctx context.Context, id int64) (DeleteAPIKeyResponse, error) {
	response := DeleteAPIKeyResponse{}

	path := fmt.Sprintf("/api/auth/keys/%d", id)
	err := c.request(ctx, "DELETE", path, nil, nil, &response)
	return response, err
}

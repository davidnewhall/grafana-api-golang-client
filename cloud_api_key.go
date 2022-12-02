package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type CreateCloudAPIKeyInput struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type ListCloudAPIKeysOutput struct {
	Items []*CloudAPIKey
}

type CloudAPIKey struct {
	ID         int
	Name       string
	Role       string
	Token      string
	Expiration string
}

func (c *Client) CreateCloudAPIKey(org string, input *CreateCloudAPIKeyInput) (*CloudAPIKey, error) {
	return c.CreateCloudAPIKeyContext(context.Background(), org, input)
}

// CreateCloudAPIKeyContext does the same thing as CreateCloudAPIKey(), but also takes in a context.
func (c *Client) CreateCloudAPIKeyContext(ctx context.Context, org string, input *CreateCloudAPIKeyInput) (*CloudAPIKey, error) {
	resp := CloudAPIKey{}
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = c.request(ctx, "POST", fmt.Sprintf("/api/orgs/%s/api-keys", org), nil, bytes.NewBuffer(data), &resp)
	return &resp, err
}

func (c *Client) ListCloudAPIKeys(org string) (*ListCloudAPIKeysOutput, error) {
	return c.ListCloudAPIKeysContext(context.Background(), org)
}

// ListCloudAPIKeysContext does the same thing as ListCloudAPIKeys(), but also takes in a context.
func (c *Client) ListCloudAPIKeysContext(ctx context.Context, org string) (*ListCloudAPIKeysOutput, error) {
	resp := &ListCloudAPIKeysOutput{}
	err := c.request(ctx, "GET", fmt.Sprintf("/api/orgs/%s/api-keys", org), nil, nil, &resp)
	return resp, err
}

func (c *Client) DeleteCloudAPIKey(org string, keyName string) error {
	return c.DeleteCloudAPIKeyContext(context.Background(), org, keyName)
}

// DeleteCloudAPIKeyContext does the same thing as DeleteCloudAPIKey(), but also takes in a context.
func (c *Client) DeleteCloudAPIKeyContext(ctx context.Context, org string, keyName string) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/api/orgs/%s/api-keys/%s", org, keyName), nil, nil, nil)
}

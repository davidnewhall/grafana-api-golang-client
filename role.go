package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type Role struct {
	Version     int64        `json:"version"`
	UID         string       `json:"uid,omitempty"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Global      bool         `json:"global"`
	Group       string       `json:"group"`
	DisplayName string       `json:"displayName"`
	Hidden      bool         `json:"hidden"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Action string `json:"action"`
	Scope  string `json:"scope"`
}

// GetRole gets a role with permissions for the given UID. Available only in Grafana Enterprise 8.+.
func (c *Client) GetRole(uid string) (*Role, error) {
	return c.GetRoleContext(context.Background(), uid)
}

// GetRoleContext does the same thing as GetRole(), but also takes in a context.
func (c *Client) GetRoleContext(ctx context.Context, uid string) (*Role, error) {
	r := &Role{}
	err := c.request(ctx, "GET", buildURL(uid), nil, nil, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NewRole creates a new role with permissions. Available only in Grafana Enterprise 8.+.
func (c *Client) NewRole(role Role) (*Role, error) {
	return c.NewRoleContext(context.Background(), role)
}

// NewRoleContext does the same thing as NewRole(), but also takes in a context.
func (c *Client) NewRoleContext(ctx context.Context, role Role) (*Role, error) {
	data, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	r := &Role{}

	err = c.request(ctx, "POST", "/api/access-control/roles", nil, bytes.NewBuffer(data), &r)
	if err != nil {
		return nil, err
	}

	return r, err
}

// UpdateRole updates the role and permissions. Available only in Grafana Enterprise 8.+.
func (c *Client) UpdateRole(role Role) error {
	return c.UpdateRoleContext(context.Background(), role)
}

// UpdateRoleContext does the same thing as UpdateRole(), but also takes in a context.
func (c *Client) UpdateRoleContext(ctx context.Context, role Role) error {
	data, err := json.Marshal(role)
	if err != nil {
		return err
	}

	err = c.request(ctx, "PUT", buildURL(role.UID), nil, bytes.NewBuffer(data), nil)

	return err
}

// DeleteRole deletes the role with it's permissions. Available only in Grafana Enterprise 8.+.
func (c *Client) DeleteRole(uid string, global bool) error {
	return c.DeleteRoleContext(context.Background(), uid, global)
}

// DeleteRoleContext does the same thing as DeleteRole(), but also takes in a context.
func (c *Client) DeleteRoleContext(ctx context.Context, uid string, global bool) error {
	qp := map[string][]string{
		"global": {fmt.Sprint(global)},
	}
	return c.request(ctx, "DELETE", buildURL(uid), qp, nil, nil)
}

func buildURL(uid string) string {
	const rootURL = "/api/access-control/roles"
	return fmt.Sprintf("%s/%s", rootURL, uid)
}

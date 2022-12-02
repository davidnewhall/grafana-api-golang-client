package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// OrgUser represents a Grafana org user.
type OrgUser struct {
	OrgID  int64  `json:"orgId"`
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}

// OrgUsersCurrent returns all org users within the current organization.
// This endpoint is accessible to users with org admin role.
func (c *Client) OrgUsersCurrent() ([]OrgUser, error) {
	return c.OrgUsersCurrentContext(context.Background())
}

// OrgUsersCurrentContext does the same thing as OrgUsersCurrent(), but also takes in a context.
func (c *Client) OrgUsersCurrentContext(ctx context.Context) ([]OrgUser, error) {
	users := make([]OrgUser, 0)
	err := c.request(ctx, "GET", "/api/org/users", nil, nil, &users)
	if err != nil {
		return nil, err
	}
	return users, err
}

// OrgUsers fetches and returns the users for the org whose ID it's passed.
func (c *Client) OrgUsers(orgID int64) ([]OrgUser, error) {
	return c.OrgUsersContext(context.Background(), orgID)
}

// OrgUsersContext does the same thing as OrgUsers(), but also takes in a context.
func (c *Client) OrgUsersContext(ctx context.Context, orgID int64) ([]OrgUser, error) {
	users := make([]OrgUser, 0)
	err := c.request(ctx, "GET", fmt.Sprintf("/api/orgs/%d/users", orgID), nil, nil, &users)
	if err != nil {
		return users, err
	}

	return users, err
}

// AddOrgUser adds a user to an org with the specified role.
func (c *Client) AddOrgUser(orgID int64, user, role string) error {
	return c.AddOrgUserContext(context.Background(), orgID, user, role)
}

// AddOrgUserContext does the same thing as AddOrgUser(), but also takes in a context.
func (c *Client) AddOrgUserContext(ctx context.Context, orgID int64, user, role string) error {
	dataMap := map[string]string{
		"loginOrEmail": user,
		"role":         role,
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	return c.request(ctx, "POST", fmt.Sprintf("/api/orgs/%d/users", orgID), nil, bytes.NewBuffer(data), nil)
}

// UpdateOrgUser updates and org user.
func (c *Client) UpdateOrgUser(orgID, userID int64, role string) error {
	return c.UpdateOrgUserContext(context.Background(), orgID, userID, role)
}

// UpdateOrgUserContext does the same thing as UpdateOrgUser(), but also takes in a context.
func (c *Client) UpdateOrgUserContext(ctx context.Context, orgID, userID int64, role string) error {
	dataMap := map[string]string{
		"role": role,
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	return c.request(ctx, "PATCH", fmt.Sprintf("/api/orgs/%d/users/%d", orgID, userID), nil, bytes.NewBuffer(data), nil)
}

// RemoveOrgUser removes a user from an org.
func (c *Client) RemoveOrgUser(orgID, userID int64) error {
	return c.RemoveOrgUserContext(context.Background(), orgID, userID)
}

// RemoveOrgUserContext does the same thing as RemoveOrgUser(), but also takes in a context.
func (c *Client) RemoveOrgUserContext(ctx context.Context, orgID, userID int64) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/api/orgs/%d/users/%d", orgID, userID), nil, nil, nil)
}

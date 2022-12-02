package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const baseURL = "/api/access-control/builtin-roles"

type BuiltInRoleAssignment struct {
	BuiltinRole string `json:"builtInRole"`
	RoleUID     string `json:"roleUid"`
	Global      bool   `json:"global"`
}

// GetBuiltInRoleAssignments gets all built-in role assignments. Available only in Grafana Enterprise 8.+.
func (c *Client) GetBuiltInRoleAssignments() (map[string][]*Role, error) {
	return c.GetBuiltInRoleAssignmentsContext(context.Background())
}

// GetBuiltInRoleAssignmentsContext does the same thing as GetBuiltInRoleAssignments(), but also takes in a context.
func (c *Client) GetBuiltInRoleAssignmentsContext(ctx context.Context) (map[string][]*Role, error) {
	br := make(map[string][]*Role)
	err := c.request(ctx, "GET", baseURL, nil, nil, &br)
	if err != nil {
		return nil, err
	}
	return br, nil
}

// NewBuiltInRoleAssignment creates a new built-in role assignment. Available only in Grafana Enterprise 8.+.
func (c *Client) NewBuiltInRoleAssignment(builtInRoleAssignment BuiltInRoleAssignment) (*BuiltInRoleAssignment, error) {
	return c.NewBuiltInRoleAssignmentContext(context.Background(), builtInRoleAssignment)
}

// NewBuiltInRoleAssignmentContext does the same thing as NewBuiltInRoleAssignment(), but also takes in a context.
func (c *Client) NewBuiltInRoleAssignmentContext(ctx context.Context, builtInRoleAssignment BuiltInRoleAssignment) (*BuiltInRoleAssignment, error) {
	body, err := json.Marshal(builtInRoleAssignment)
	if err != nil {
		return nil, err
	}

	br := &BuiltInRoleAssignment{}

	err = c.request(ctx, "POST", baseURL, nil, bytes.NewBuffer(body), &br)
	if err != nil {
		return nil, err
	}

	return br, err
}

// DeleteBuiltInRoleAssignment remove the built-in role assignments. Available only in Grafana Enterprise 8.+.
func (c *Client) DeleteBuiltInRoleAssignment(builtInRole BuiltInRoleAssignment) error {
	return c.DeleteBuiltInRoleAssignmentContext(context.Background(), builtInRole)
}

// DeleteBuiltInRoleAssignmentContext does the same thing as DeleteBuiltInRoleAssignment(), but also takes in a context.
func (c *Client) DeleteBuiltInRoleAssignmentContext(ctx context.Context, builtInRole BuiltInRoleAssignment) error {
	data, err := json.Marshal(builtInRole)
	if err != nil {
		return err
	}

	qp := map[string][]string{
		"global": {fmt.Sprint(builtInRole.Global)},
	}
	url := fmt.Sprintf("%s/%s/roles/%s", baseURL, builtInRole.BuiltinRole, builtInRole.RoleUID)
	err = c.request(ctx, "DELETE", url, qp, bytes.NewBuffer(data), nil)

	return err
}

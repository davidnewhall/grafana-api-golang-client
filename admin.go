package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// PauseAllAlertsResponse represents the response body for a PauseAllAlerts request.
type PauseAllAlertsResponse struct {
	AlertsAffected int64  `json:"alertsAffected,omitempty"`
	State          string `json:"state,omitempty"`
	Message        string `json:"message,omitempty"`
}

// CreateUser creates a Grafana user.
func (c *Client) CreateUser(user User) (int64, error) {
	return c.CreateUserContext(context.Background(), user)
}

// CreateUserContext creates a Grafana user w/ context.
func (c *Client) CreateUserContext(ctx context.Context, user User) (int64, error) {
	id := int64(0)
	data, err := json.Marshal(user)
	if err != nil {
		return id, err
	}

	created := struct {
		ID int64 `json:"id"`
	}{}

	err = c.request(ctx, "POST", "/api/admin/users", nil, bytes.NewBuffer(data), &created)
	if err != nil {
		return id, err
	}

	return created.ID, err
}

// DeleteUser deletes a Grafana user.
func (c *Client) DeleteUser(id int64) error {
	return c.DeleteUserContext(context.Background(), id)
}

// DeleteUserContext deletes a Grafana user w/ context.
func (c *Client) DeleteUserContext(ctx context.Context, id int64) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/api/admin/users/%d", id), nil, nil, nil)
}

// UpdateUserPassword updates a user password.
func (c *Client) UpdateUserPassword(id int64, password string) error {
	return c.UpdateUserPasswordContext(context.Background(), id, password)
}

// UpdateUserPasswordContext updates a user password w/ context.
func (c *Client) UpdateUserPasswordContext(ctx context.Context, id int64, password string) error {
	body := map[string]string{"password": password}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return c.request(ctx, "PUT", fmt.Sprintf("/api/admin/users/%d/password", id), nil, bytes.NewBuffer(data), nil)
}

// UpdateUserPermissions sets a user's admin status.
func (c *Client) UpdateUserPermissions(id int64, isAdmin bool) error {
	return c.UpdateUserPermissionsContext(context.Background(), id, isAdmin)
}

// UpdateUserPermissionsContext sets a user's admin status w/ context.
func (c *Client) UpdateUserPermissionsContext(ctx context.Context, id int64, isAdmin bool) error {
	body := map[string]bool{"isGrafanaAdmin": isAdmin}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return c.request(ctx, "PUT", fmt.Sprintf("/api/admin/users/%d/permissions", id), nil, bytes.NewBuffer(data), nil)
}

// PauseAllAlerts pauses all Grafana alerts.
func (c *Client) PauseAllAlerts() (PauseAllAlertsResponse, error) {
	return c.PauseAllAlertsContext(context.Background())
}

// PauseAllAlertsContext pauses all Grafana alerts w/ context.
func (c *Client) PauseAllAlertsContext(ctx context.Context) (PauseAllAlertsResponse, error) {
	result := PauseAllAlertsResponse{}
	data, err := json.Marshal(PauseAlertRequest{
		Paused: true,
	})
	if err != nil {
		return result, err
	}

	err = c.request(ctx, "POST", "/api/admin/pause-all-alerts", nil, bytes.NewBuffer(data), &result)

	return result, err
}

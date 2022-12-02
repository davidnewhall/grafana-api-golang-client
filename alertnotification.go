package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// AlertNotification represents a Grafana alert notification.
// Deprecated: Grafana Legacy Alerting is deprecated as of 9.0 and will be removed in the future. Use ContactPoint instead.
type AlertNotification struct {
	ID                    int64       `json:"id,omitempty"`
	UID                   string      `json:"uid"`
	Name                  string      `json:"name"`
	Type                  string      `json:"type"`
	IsDefault             bool        `json:"isDefault"`
	DisableResolveMessage bool        `json:"disableResolveMessage"`
	SendReminder          bool        `json:"sendReminder"`
	Frequency             string      `json:"frequency"`
	Settings              interface{} `json:"settings"`
	SecureFields          interface{} `json:"secureFields,omitempty"`
	SecureSettings        interface{} `json:"secureSettings,omitempty"`
}

// AlertNotifications fetches and returns Grafana alert notifications.
// Deprecated: Grafana Legacy Alerting is deprecated as of 9.0 and will be removed in the future. Use ContactPoints instead.
func (c *Client) AlertNotifications() ([]AlertNotification, error) {
	return c.AlertNotificationsContext(context.Background())
}

// AlertNotificationsContext does the same thing as AlertNotifications(), but also takes in a context.
func (c *Client) AlertNotificationsContext(ctx context.Context) ([]AlertNotification, error) {
	alertnotifications := make([]AlertNotification, 0)

	err := c.request(ctx, "GET", "/api/alert-notifications/", nil, nil, &alertnotifications)
	if err != nil {
		return nil, err
	}

	return alertnotifications, err
}

// AlertNotification fetches and returns a Grafana alert notification.
// Deprecated: Grafana Legacy Alerting is deprecated as of 9.0 and will be removed in the future. Use ContactPoint instead.
func (c *Client) AlertNotification(id int64) (*AlertNotification, error) {
	return c.AlertNotificationContext(context.Background(), id)
}

// AlertNotificationContext does the same thing as AlertNotification(), but also takes in a context.
func (c *Client) AlertNotificationContext(ctx context.Context, id int64) (*AlertNotification, error) {
	path := fmt.Sprintf("/api/alert-notifications/%d", id)
	result := &AlertNotification{}
	err := c.request(ctx, "GET", path, nil, nil, result)
	if err != nil {
		return nil, err
	}

	return result, err
}

// NewAlertNotification creates a new Grafana alert notification.
// Deprecated: Grafana Legacy Alerting is deprecated as of 9.0 and will be removed in the future. Use NewContactPoint instead.
func (c *Client) NewAlertNotification(a *AlertNotification) (int64, error) {
	return c.NewAlertNotificationContext(context.Background(), a)
}

// NewAlertNotificationContext does the same thing as NewAlertNotification(), but also takes in a context.
func (c *Client) NewAlertNotificationContext(ctx context.Context, a *AlertNotification) (int64, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return 0, err
	}
	result := struct {
		ID int64 `json:"id"`
	}{}

	err = c.request(ctx, "POST", "/api/alert-notifications", nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return 0, err
	}

	return result.ID, err
}

// UpdateAlertNotification updates a Grafana alert notification.
// Deprecated: Grafana Legacy Alerting is deprecated as of 9.0 and will be removed in the future. Use UpdateContactPoint instead.
func (c *Client) UpdateAlertNotification(a *AlertNotification) error {
	return c.UpdateAlertNotificationContext(context.Background(), a)
}

// UpdateAlertNotificationContext does the same thing as UpdateAlertNotification(), but also takes in a context.
func (c *Client) UpdateAlertNotificationContext(ctx context.Context, a *AlertNotification) error {
	path := fmt.Sprintf("/api/alert-notifications/%d", a.ID)
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	err = c.request(ctx, "PUT", path, nil, bytes.NewBuffer(data), nil)

	return err
}

// DeleteAlertNotification deletes a Grafana alert notification.
// Deprecated: Grafana Legacy Alerting is deprecated as of 9.0 and will be removed in the future. Use DeleteContactPoint instead.
func (c *Client) DeleteAlertNotification(id int64) error {
	return c.DeleteAlertNotificationContext(context.Background(), id)
}

// DeleteAlertNotificationContext does the same thing as DeleteAlertNotification(), but also takes in a context.
func (c *Client) DeleteAlertNotificationContext(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/alert-notifications/%d", id)

	return c.request(ctx, "DELETE", path, nil, nil, nil)
}

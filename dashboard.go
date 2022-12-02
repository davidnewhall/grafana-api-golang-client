package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// DashboardMeta represents Grafana dashboard meta.
type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
	Folder    int64  `json:"folderId"`
	URL       string `json:"url"`
}

// DashboardSaveResponse represents the Grafana API response to creating or saving a dashboard.
type DashboardSaveResponse struct {
	Slug    string `json:"slug"`
	ID      int64  `json:"id"`
	UID     string `json:"uid"`
	Status  string `json:"status"`
	Version int64  `json:"version"`
}

// Dashboard represents a Grafana dashboard.
type Dashboard struct {
	Meta      DashboardMeta          `json:"meta"`
	Model     map[string]interface{} `json:"dashboard"`
	FolderID  int64                  `json:"folderId"`
	FolderUID string                 `json:"folderUid"`
	Overwrite bool                   `json:"overwrite"`

	// This is only used when creating a new dashboard, it will always be empty when getting a dashboard.
	Message string `json:"message"`
}

// SaveDashboard is a deprecated method for saving a Grafana dashboard. Use NewDashboard.
// Deprecated: Use NewDashboard instead.
func (c *Client) SaveDashboard(model map[string]interface{}, overwrite bool) (*DashboardSaveResponse, error) {
	return c.SaveDashboardContext(context.Background(), model, overwrite)
}

// SaveDashboardContext does the same thing as SaveDashboard(), but also takes in a context.
func (c *Client) SaveDashboardContext(ctx context.Context, model map[string]interface{}, overwrite bool) (*DashboardSaveResponse, error) {
	wrapper := map[string]interface{}{
		"dashboard": model,
		"overwrite": overwrite,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	result := &DashboardSaveResponse{}
	err = c.request(ctx, "POST", "/api/dashboards/db", nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

// NewDashboard creates a new Grafana dashboard.
func (c *Client) NewDashboard(dashboard Dashboard) (*DashboardSaveResponse, error) {
	return c.NewDashboardContext(context.Background(), dashboard)
}

// NewDashboardContext does the same thing as NewDashboard(), but also takes in a context.
func (c *Client) NewDashboardContext(ctx context.Context, dashboard Dashboard) (*DashboardSaveResponse, error) {
	data, err := json.Marshal(dashboard)
	if err != nil {
		return nil, err
	}

	result := &DashboardSaveResponse{}
	err = c.request(ctx, "POST", "/api/dashboards/db", nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

// Dashboards fetches and returns all dashboards.
func (c *Client) Dashboards() ([]FolderDashboardSearchResponse, error) {
	return c.DashboardsContext(context.Background())
}

// DashboardsContext does the same thing as Dashboards(), but also takes in a context.
func (c *Client) DashboardsContext(ctx context.Context) ([]FolderDashboardSearchResponse, error) {
	const limit = 1000

	var (
		page          = 0
		newDashboards []FolderDashboardSearchResponse
		dashboards    []FolderDashboardSearchResponse
		query         = make(url.Values)
	)

	query.Set("type", "dash-db")
	query.Set("limit", fmt.Sprint(limit))

	for {
		page++
		query.Set("page", fmt.Sprint(page))

		if err := c.request(ctx, "GET", "/api/search", query, nil, &newDashboards); err != nil {
			return nil, err
		}

		dashboards = append(dashboards, newDashboards...)

		if len(newDashboards) < limit {
			return dashboards, nil
		}
	}
}

// Dashboard will be removed.
// Deprecated: Starting from Grafana v5.0. Use DashboardByUID instead.
func (c *Client) Dashboard(slug string) (*Dashboard, error) {
	return c.DashboardContext(context.Background(), slug)
}

// DashboardContext does the same thing as Dashboard(), but also takes in a context.
func (c *Client) DashboardContext(ctx context.Context, slug string) (*Dashboard, error) {
	return c.dashboard(fmt.Sprintf("/api/dashboards/db/%s", slug))
}

// DashboardByUID gets a dashboard by UID.
func (c *Client) DashboardByUID(uid string) (*Dashboard, error) {
	return c.DashboardByUIDContext(context.Background(), uid)
}

// DashboardByUIDContext does the same thing as DashboardByUID(), but also takes in a context.
func (c *Client) DashboardByUIDContext(ctx context.Context, uid string) (*Dashboard, error) {
	return c.dashboard(fmt.Sprintf("/api/dashboards/uid/%s", uid))
}

// DashboardsByIDs uses the folder and dashboard search endpoint to find
// dashboards by list of dashboard IDs.
func (c *Client) DashboardsByIDs(ids []int64) ([]FolderDashboardSearchResponse, error) {
	return c.DashboardsByIDsContext(context.Background(), ids)
}

// DashboardsByIDsContext does the same thing as DashboardsByIDs(), but also takes in a context.
func (c *Client) DashboardsByIDsContext(ctx context.Context, ids []int64) ([]FolderDashboardSearchResponse, error) {
	dashboardIdsJSON, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"type":         {"dash-db"},
		"dashboardIds": {string(dashboardIdsJSON)},
	}
	return c.FolderDashboardSearch(params)
}

func (c *Client) dashboard(path string) (*Dashboard, error) {
	return c.dashboardContext(context.Background(), path)
}

// dashboardContext does the same thing as dashboard(), but also takes in a context.
func (c *Client) dashboardContext(ctx context.Context, path string) (*Dashboard, error) {
	result := &Dashboard{}
	err := c.request(ctx, "GET", path, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	result.FolderID = result.Meta.Folder

	return result, err
}

// DeleteDashboard will be removed.
// Deprecated: Starting from Grafana v5.0. Use DeleteDashboardByUID instead.
func (c *Client) DeleteDashboard(slug string) error {
	return c.DeleteDashboardContext(context.Background(), slug)
}

// DeleteDashboardContext does the same thing as DeleteDashboard(), but also takes in a context.
func (c *Client) DeleteDashboardContext(ctx context.Context, slug string) error {
	return c.deleteDashboard(fmt.Sprintf("/api/dashboards/db/%s", slug))
}

// DeleteDashboardByUID deletes a dashboard by UID.
func (c *Client) DeleteDashboardByUID(uid string) error {
	return c.DeleteDashboardByUIDContext(context.Background(), uid)
}

// DeleteDashboardByUIDContext does the same thing as DeleteDashboardByUID(), but also takes in a context.
func (c *Client) DeleteDashboardByUIDContext(ctx context.Context, uid string) error {
	return c.deleteDashboard(fmt.Sprintf("/api/dashboards/uid/%s", uid))
}

func (c *Client) deleteDashboard(path string) error {
	return c.deleteDashboardContext(context.Background(), path)
}

// deleteDashboardContext does the same thing as deleteDashboard(), but also takes in a context.
func (c *Client) deleteDashboardContext(ctx context.Context, path string) error {
	return c.request(ctx, "DELETE", path, nil, nil, nil)
}

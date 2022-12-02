package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Plugin struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type CloudPluginInstallation struct {
	ID           int    `json:"id"`
	InstanceID   int    `json:"instanceId"`
	InstanceURL  string `json:"instanceUrl"`
	InstanceSlug string `json:"instanceSlug"`
	PluginID     int    `json:"pluginId"`
	PluginSlug   string `json:"pluginSlug"`
	PluginName   string `json:"pluginName"`
	Version      string `json:"version"`
}

// InstallCloudPlugin installs the specified plugin to the given stack.
func (c *Client) InstallCloudPlugin(stackSlug string, pluginSlug string, pluginVersion string) (*CloudPluginInstallation, error) {
	return c.InstallCloudPluginContext(context.Background(), stackSlug, pluginSlug, pluginVersion)
}

// InstallCloudPluginContext does the same thing as InstallCloudPlugin(), but also takes in a context.
func (c *Client) InstallCloudPluginContext(ctx context.Context, stackSlug string, pluginSlug string, pluginVersion string) (*CloudPluginInstallation, error) {
	installPluginRequest := struct {
		Plugin  string `json:"plugin"`
		Version string `json:"version"`
	}{
		Plugin:  pluginSlug,
		Version: pluginVersion,
	}

	data, err := json.Marshal(installPluginRequest)
	if err != nil {
		return nil, err
	}

	var installation CloudPluginInstallation

	err = c.request(ctx, "POST", fmt.Sprintf("/api/instances/%s/plugins", stackSlug), nil, bytes.NewBuffer(data), &installation)
	if err != nil {
		return nil, err
	}

	return &installation, nil
}

// UninstallCloudPlugin uninstalls the specified plugin to the given stack.
func (c *Client) UninstallCloudPlugin(stackSlug string, pluginSlug string) error {
	return c.UninstallCloudPluginContext(context.Background(), stackSlug, pluginSlug)
}

// UninstallCloudPluginContext does the same thing as UninstallCloudPlugin(), but also takes in a context.
func (c *Client) UninstallCloudPluginContext(ctx context.Context, stackSlug string, pluginSlug string) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/api/instances/%s/plugins/%s", stackSlug, pluginSlug), nil, nil, nil)
}

// IsCloudPluginInstalled returns a boolean if the specified plugin is installed on the stack.
func (c *Client) IsCloudPluginInstalled(stackSlug string, pluginSlug string) (bool, error) {
	return c.IsCloudPluginInstalledContext(context.Background(), stackSlug, pluginSlug)
}

// IsCloudPluginInstalledContext does the same thing as IsCloudPluginInstalled(), but also takes in a context.
func (c *Client) IsCloudPluginInstalledContext(ctx context.Context, stackSlug string, pluginSlug string) (bool, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/api/instances/%s/plugins/%s", stackSlug, pluginSlug), nil, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		bodyContents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		return false, fmt.Errorf("status: %d, body: %v", resp.StatusCode, string(bodyContents))
	}

	return true, nil
}

// GetCloudPluginInstallation returns the cloud plugin installation details for the specified plugin.
func (c *Client) GetCloudPluginInstallation(stackSlug string, pluginSlug string) (*CloudPluginInstallation, error) {
	return c.GetCloudPluginInstallationContext(context.Background(), stackSlug, pluginSlug)
}

// GetCloudPluginInstallationContext does the same thing as GetCloudPluginInstallation(), but also takes in a context.
func (c *Client) GetCloudPluginInstallationContext(ctx context.Context, stackSlug string, pluginSlug string) (*CloudPluginInstallation, error) {
	var installation CloudPluginInstallation

	err := c.request(ctx, "GET", fmt.Sprintf("/api/instances/%s/plugins/%s", stackSlug, pluginSlug), nil, nil, &installation)
	if err != nil {
		return nil, err
	}

	return &installation, nil
}

// PluginBySlug returns the plugin with the given slug.
// An error will be returned given an unknown slug.
func (c *Client) PluginBySlug(slug string) (*Plugin, error) {
	return c.PluginBySlugContext(context.Background(), slug)
}

// PluginBySlugContext does the same thing as PluginBySlug(), but also takes in a context.
func (c *Client) PluginBySlugContext(ctx context.Context, slug string) (*Plugin, error) {
	p := Plugin{}

	err := c.request(ctx, "GET", fmt.Sprintf("/api/plugins/%s", slug), nil, nil, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// PluginByID returns the plugin with the given id.
// An error will be returned given an unknown ID.
func (c *Client) PluginByID(pluginID int64) (*Plugin, error) {
	return c.PluginByIDContext(context.Background(), pluginID)
}

// PluginByIDContext does the same thing as PluginByID(), but also takes in a context.
func (c *Client) PluginByIDContext(ctx context.Context, pluginID int64) (*Plugin, error) {
	p := Plugin{}

	err := c.request(ctx, "GET", fmt.Sprintf("/api/plugins/%d", pluginID), nil, nil, p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

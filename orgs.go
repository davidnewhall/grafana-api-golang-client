package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// Org represents a Grafana org.
type Org struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Orgs fetches and returns the Grafana orgs.
func (c *Client) Orgs() ([]Org, error) {
	return c.OrgsContext(context.Background())
}

// OrgsContext does the same thing as Orgs(), but also takes in a context.
func (c *Client) OrgsContext(ctx context.Context) ([]Org, error) {
	orgs := make([]Org, 0)
	err := c.request(ctx, "GET", "/api/orgs/", nil, nil, &orgs)
	if err != nil {
		return orgs, err
	}

	return orgs, err
}

// OrgByName fetches and returns the org whose name it's passed.
func (c *Client) OrgByName(name string) (Org, error) {
	return c.OrgByNameContext(context.Background(), name)
}

// OrgByNameContext does the same thing as OrgByName(), but also takes in a context.
func (c *Client) OrgByNameContext(ctx context.Context, name string) (Org, error) {
	org := Org{}
	err := c.request(ctx, "GET", fmt.Sprintf("/api/orgs/name/%s", name), nil, nil, &org)
	if err != nil {
		return org, err
	}

	return org, err
}

// Org fetches and returns the org whose ID it's passed.
func (c *Client) Org(id int64) (Org, error) {
	return c.OrgContext(context.Background(), id)
}

// OrgContext does the same thing as Org(), but also takes in a context.
func (c *Client) OrgContext(ctx context.Context, id int64) (Org, error) {
	org := Org{}
	err := c.request(ctx, "GET", fmt.Sprintf("/api/orgs/%d", id), nil, nil, &org)
	if err != nil {
		return org, err
	}

	return org, err
}

// NewOrg creates a new Grafana org.
func (c *Client) NewOrg(name string) (int64, error) {
	return c.NewOrgContext(context.Background(), name)
}

// NewOrgContext does the same thing as NewOrg(), but also takes in a context.
func (c *Client) NewOrgContext(ctx context.Context, name string) (int64, error) {
	id := int64(0)

	dataMap := map[string]string{
		"name": name,
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return id, err
	}
	tmp := struct {
		ID int64 `json:"orgId"`
	}{}

	err = c.request(ctx, "POST", "/api/orgs", nil, bytes.NewBuffer(data), &tmp)
	if err != nil {
		return id, err
	}

	return tmp.ID, err
}

// UpdateOrg updates a Grafana org.
func (c *Client) UpdateOrg(id int64, name string) error {
	return c.UpdateOrgContext(context.Background(), id, name)
}

// UpdateOrgContext does the same thing as UpdateOrg(), but also takes in a context.
func (c *Client) UpdateOrgContext(ctx context.Context, id int64, name string) error {
	dataMap := map[string]string{
		"name": name,
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	return c.request(ctx, "PUT", fmt.Sprintf("/api/orgs/%d", id), nil, bytes.NewBuffer(data), nil)
}

// DeleteOrg deletes the Grafana org whose ID it's passed.
func (c *Client) DeleteOrg(id int64) error {
	return c.DeleteOrgContext(context.Background(), id)
}

// DeleteOrgContext does the same thing as DeleteOrg(), but also takes in a context.
func (c *Client) DeleteOrgContext(ctx context.Context, id int64) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/api/orgs/%d", id), nil, nil, nil)
}

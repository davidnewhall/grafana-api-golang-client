package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// TeamGroup represents a Grafana TeamGroup.
type TeamGroup struct {
	OrgID   int64  `json:"orgId,omitempty"`
	TeamID  int64  `json:"teamId,omitempty"`
	GroupID string `json:"groupID,omitempty"`
}

// TeamGroups fetches and returns the list of Grafana team group whose Team ID it's passed.
func (c *Client) TeamGroups(id int64) ([]TeamGroup, error) {
	return c.TeamGroupsContext(context.Background(), id)
}

// TeamGroupsContext does the same thing as TeamGroups(), but also takes in a context.
func (c *Client) TeamGroupsContext(ctx context.Context, id int64) ([]TeamGroup, error) {
	teamGroups := make([]TeamGroup, 0)
	err := c.request(ctx, "GET", fmt.Sprintf("/api/teams/%d/groups", id), nil, nil, &teamGroups)
	if err != nil {
		return teamGroups, err
	}

	return teamGroups, nil
}

// NewTeamGroup creates a new Grafana Team Group .
func (c *Client) NewTeamGroup(id int64, groupID string) error {
	return c.NewTeamGroupContext(context.Background(), id, groupID)
}

// NewTeamGroupContext does the same thing as NewTeamGroup(), but also takes in a context.
func (c *Client) NewTeamGroupContext(ctx context.Context, id int64, groupID string) error {
	dataMap := map[string]string{
		"groupId": groupID,
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	return c.request(ctx, "POST", fmt.Sprintf("/api/teams/%d/groups", id), nil, bytes.NewBuffer(data), nil)
}

// DeleteTeam deletes the Grafana team whose ID it's passed.
func (c *Client) DeleteTeamGroup(id int64, groupID string) error {
	return c.DeleteTeamGroupContext(context.Background(), id, groupID)
}

// DeleteTeamGroupContext does the same thing as DeleteTeamGroup(), but also takes in a context.
func (c *Client) DeleteTeamGroupContext(ctx context.Context, id int64, groupID string) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/api/teams/%d/groups/%s", id, groupID), nil, nil, nil)
}

package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// PlaylistItem represents a Grafana playlist item.
type PlaylistItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Order int    `json:"order"`
	Title string `json:"title"`
}

// Playlist represents a Grafana playlist.
type Playlist struct {
	ID       int            `json:"id,omitempty"`  // Grafana < 9.0
	UID      string         `json:"uid,omitempty"` // Grafana >= 9.0
	Name     string         `json:"name"`
	Interval string         `json:"interval"`
	Items    []PlaylistItem `json:"items"`
}

// Grafana 9.0+ returns the ID and the UID but uses the UID in the API calls.
// Grafana <9 only returns the ID.
func (p *Playlist) QueryID() string {
	if p.UID != "" {
		return p.UID
	}
	return fmt.Sprintf("%d", p.ID)
}

// Playlist fetches and returns a Grafana playlist.
func (c *Client) Playlist(idOrUID string) (*Playlist, error) {
	return c.PlaylistContext(context.Background(), idOrUID)
}

// PlaylistContext does the same thing as Playlist(), but also takes in a context.
func (c *Client) PlaylistContext(ctx context.Context, idOrUID string) (*Playlist, error) {
	path := fmt.Sprintf("/api/playlists/%s", idOrUID)
	playlist := &Playlist{}
	err := c.request(ctx, "GET", path, nil, nil, playlist)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

// NewPlaylist creates a new Grafana playlist.
func (c *Client) NewPlaylist(playlist Playlist) (string, error) {
	return c.NewPlaylistContext(context.Background(), playlist)
}

// NewPlaylistContext does the same thing as NewPlaylist(), but also takes in a context.
func (c *Client) NewPlaylistContext(ctx context.Context, playlist Playlist) (string, error) {
	data, err := json.Marshal(playlist)
	if err != nil {
		return "", err
	}

	var result Playlist

	err = c.request(ctx, "POST", "/api/playlists", nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return "", err
	}

	return result.QueryID(), nil
}

// UpdatePlaylist updates a Grafana playlist.
func (c *Client) UpdatePlaylist(playlist Playlist) error {
	return c.UpdatePlaylistContext(context.Background(), playlist)
}

// UpdatePlaylistContext does the same thing as UpdatePlaylist(), but also takes in a context.
func (c *Client) UpdatePlaylistContext(ctx context.Context, playlist Playlist) error {
	path := fmt.Sprintf("/api/playlists/%s", playlist.QueryID())
	data, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	return c.request(ctx, "PUT", path, nil, bytes.NewBuffer(data), nil)
}

// DeletePlaylist deletes the Grafana playlist whose ID it's passed.
func (c *Client) DeletePlaylist(idOrUID string) error {
	return c.DeletePlaylistContext(context.Background(), idOrUID)
}

// DeletePlaylistContext does the same thing as DeletePlaylist(), but also takes in a context.
func (c *Client) DeletePlaylistContext(ctx context.Context, idOrUID string) error {
	path := fmt.Sprintf("/api/playlists/%s", idOrUID)

	return c.request(ctx, "DELETE", path, nil, nil, nil)
}

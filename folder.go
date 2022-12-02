package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// Folder represents a Grafana folder.
type Folder struct {
	ID    int64  `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type FolderPayload struct {
	Title     string `json:"title"`
	UID       string `json:"uid,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`
}

// Folders fetches and returns Grafana folders.
func (c *Client) Folders() ([]Folder, error) {
	return c.FoldersContext(context.Background())
}

// FoldersContext does the same thing as Folders(), but also takes in a context.
func (c *Client) FoldersContext(ctx context.Context) ([]Folder, error) {
	folders := make([]Folder, 0)
	err := c.request(ctx, "GET", "/api/folders/", nil, nil, &folders)
	if err != nil {
		return folders, err
	}

	return folders, err
}

// Folder fetches and returns the Grafana folder whose ID it's passed.
func (c *Client) Folder(id int64) (*Folder, error) {
	return c.FolderContext(context.Background(), id)
}

// FolderContext does the same thing as Folder(), but also takes in a context.
func (c *Client) FolderContext(ctx context.Context, id int64) (*Folder, error) {
	folder := &Folder{}
	err := c.request(ctx, "GET", fmt.Sprintf("/api/folders/id/%d", id), nil, nil, folder)
	if err != nil {
		return folder, err
	}

	return folder, err
}

// Folder fetches and returns the Grafana folder whose UID it's passed.
func (c *Client) FolderByUID(uid string) (*Folder, error) {
	return c.FolderByUIDContext(context.Background(), uid)
}

// FolderByUIDContext does the same thing as FolderByUID(), but also takes in a context.
func (c *Client) FolderByUIDContext(ctx context.Context, uid string) (*Folder, error) {
	folder := &Folder{}
	err := c.request(ctx, "GET", fmt.Sprintf("/api/folders/%s", uid), nil, nil, folder)
	if err != nil {
		return folder, err
	}

	return folder, err
}

// NewFolder creates a new Grafana folder.
func (c *Client) NewFolder(title string, uid ...string) (Folder, error) {
	return c.NewFolderContext(context.Background(), title, uid...)
}

// NewFolderContext does the same thing as NewFolder(), but also takes in a context.
func (c *Client) NewFolderContext(ctx context.Context, title string, uid ...string) (Folder, error) {
	if len(uid) > 1 {
		return Folder{}, fmt.Errorf("too many arguments. Expected 1 or 2")
	}

	folder := Folder{}
	payload := FolderPayload{
		Title: title,
	}
	if len(uid) == 1 {
		payload.UID = uid[0]
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return folder, err
	}

	err = c.request(ctx, "POST", "/api/folders", nil, bytes.NewBuffer(data), &folder)
	if err != nil {
		return folder, err
	}

	return folder, err
}

// UpdateFolder updates the folder whose UID it's passed.
func (c *Client) UpdateFolder(uid string, title string, newUID ...string) error {
	return c.UpdateFolderContext(context.Background(), uid, title, newUID...)
}

// UpdateFolderContext does the same thing as UpdateFolder(), but also takes in a context.
func (c *Client) UpdateFolderContext(ctx context.Context, uid string, title string, newUID ...string) error {
	payload := FolderPayload{
		Title:     title,
		Overwrite: true,
	}
	if len(newUID) == 1 {
		payload.UID = newUID[0]
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.request(ctx, "PUT", fmt.Sprintf("/api/folders/%s", uid), nil, bytes.NewBuffer(data), nil)
}

// DeleteFolder deletes the folder whose ID it's passed.
func (c *Client) DeleteFolder(id string) error {
	return c.DeleteFolderContext(context.Background(), id)
}

// DeleteFolderContext does the same thing as DeleteFolder(), but also takes in a context.
func (c *Client) DeleteFolderContext(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/api/folders/%s", id), nil, nil, nil)
}

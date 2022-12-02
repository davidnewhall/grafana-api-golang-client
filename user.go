package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// User represents a Grafana user. It is structured after the UserProfileDTO
// struct in the Grafana codebase.
type User struct {
	ID         int64     `json:"id,omitempty"`
	Email      string    `json:"email,omitempty"`
	Name       string    `json:"name,omitempty"`
	Login      string    `json:"login,omitempty"`
	Theme      string    `json:"theme,omitempty"`
	OrgID      int64     `json:"orgId,omitempty"`
	IsAdmin    bool      `json:"isGrafanaAdmin,omitempty"`
	IsDisabled bool      `json:"isDisabled,omitempty"`
	IsExternal bool      `json:"isExternal,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	AuthLabels []string  `json:"authLabels,omitempty"`
	AvatarURL  string    `json:"avatarUrl,omitempty"`
	Password   string    `json:"password,omitempty"`
}

// UserSearch represents a Grafana user as returned by API endpoints that
// return a collection of Grafana users. This representation of user has
// reduced and differing fields. It is structured after the UserSearchHitDTO
// struct in the Grafana codebase.
type UserSearch struct {
	ID            int64     `json:"id,omitempty"`
	Email         string    `json:"email,omitempty"`
	Name          string    `json:"name,omitempty"`
	Login         string    `json:"login,omitempty"`
	IsAdmin       bool      `json:"isAdmin,omitempty"`
	IsDisabled    bool      `json:"isDisabled,omitempty"`
	LastSeenAt    time.Time `json:"lastSeenAt,omitempty"`
	LastSeenAtAge string    `json:"lastSeenAtAge,omitempty"`
	AuthLabels    []string  `json:"authLabels,omitempty"`
	AvatarURL     string    `json:"avatarUrl,omitempty"`
}

// Users fetches and returns Grafana users.
func (c *Client) Users() (users []UserSearch, err error) {
	return c.UsersContext(context.Background())
}

// UsersContext does the same thing as Users(), but also takes in a context.
func (c *Client) UsersContext(ctx context.Context) (users []UserSearch, err error) {
	var (
		page     = 1
		newUsers []UserSearch
	)
	for len(newUsers) > 0 || page == 1 {
		query := url.Values{}
		query.Add("page", fmt.Sprintf("%d", page))
		if err = c.request(ctx, "GET", "/api/users", query, nil, &newUsers); err != nil {
			return
		}
		users = append(users, newUsers...)
		page++
	}

	return
}

// User fetches a user by ID.
func (c *Client) User(id int64) (user User, err error) {
	return c.UserContext(context.Background(), id)
}

// UserContext does the same thing as User(), but also takes in a context.
func (c *Client) UserContext(ctx context.Context, id int64) (user User, err error) {
	err = c.request(ctx, "GET", fmt.Sprintf("/api/users/%d", id), nil, nil, &user)
	return
}

// UserByEmail fetches a user by email address.
func (c *Client) UserByEmail(email string) (user User, err error) {
	return c.UserByEmailContext(context.Background(), email)
}

// UserByEmailContext does the same thing as UserByEmail(), but also takes in a context.
func (c *Client) UserByEmailContext(ctx context.Context, email string) (user User, err error) {
	query := url.Values{}
	query.Add("loginOrEmail", email)
	err = c.request(ctx, "GET", "/api/users/lookup", query, nil, &user)
	return
}

// UserUpdate updates a user by ID.
func (c *Client) UserUpdate(u User) error {
	return c.UserUpdateContext(context.Background(), u)
}

// UserUpdateContext does the same thing as UserUpdate(), but also takes in a context.
func (c *Client) UserUpdateContext(ctx context.Context, u User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return c.request(ctx, "PUT", fmt.Sprintf("/api/users/%d", u.ID), nil, bytes.NewBuffer(data), nil)
}

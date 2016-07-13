package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

type ContactMethod struct {
	ID             string
	Label          string
	Address        string
	Type           string
	SendShortEmail bool `json:"send_short_email"`
}

type NotificationRule struct {
	ID                  string
	StartDelayInMinutes uint          `json:"start_delay_in_minutes"`
	CreatedAt           string        `json:"created_at"`
	ContactMethod       ContactMethod `json:"contact_method"`
	Urgency             string
	Type                string
}

type User struct {
	APIObject
	Name              string
	Email             string
	Timezone          string
	Color             string
	Role              string
	AvatarURL         string `json:"avatar_url"`
	Description       string
	InvitationSent    bool
	ContactMethods    []ContactMethod    `json:"contact_methods"`
	NotificationRules []NotificationRule `json:"notification_rules"`
	JobTitle          string             `json:"jon_title"`
	Teams             []Team
}

type ListUsersResponse struct {
	APIListObject
	Users []User
}

type ListUserOptions struct {
	APIListObject
	Query    string   `url:"query,omitempty"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	Includes []string `url:"include,omitempty,brackets"`
}

type GetUserOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

func (c *Client) ListUsers(o ListUserOptions) (*ListUsersResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/users?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListUsersResponse
	return &result, c.decodeJson(resp, &result)
}

func (c *Client) CreateUser(u User) error {
	data := make(map[string]User)
	data["user"] = u
	_, err := c.Post("/users", data)
	return err
}

func (c *Client) DeleteUser(id string) error {
	_, err := c.Delete("/users/" + id)
	return err
}

func (c *Client) GetUser(id string, o GetUserOptions) (*User, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/users/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]User
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	u, ok := result["user"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have user field")
	}
	return &u, nil
}

func (c *Client) UpdateUser(u User) error {
	v := make(map[string]User)
	v["user"] = u
	_, err := c.Put("/users/"+u.ID, v)
	return err
}

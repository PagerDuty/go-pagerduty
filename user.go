package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

// ContactMethod is a way of contacting the user.
type ContactMethod struct {
	ID             string
	Label          string
	Address        string
	Type           string
	SendShortEmail bool `json:"send_short_email"`
}

// NotificationRule is a rule for notifying the user.
type NotificationRule struct {
	ID                  string
	StartDelayInMinutes uint          `json:"start_delay_in_minutes"`
	CreatedAt           string        `json:"created_at"`
	ContactMethod       ContactMethod `json:"contact_method"`
	Urgency             string
	Type                string
}

// User is a member of a PagerDuty account that has the ability to interact with incidents and other data on the account.
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

// ListUsersResponse is the data structure returned from calling the ListUsers API endpoint.
type ListUsersResponse struct {
	APIListObject
	Users []User
}

// ListUsersOptions is the data structure used when calling the ListUsers API endpoint.
type ListUsersOptions struct {
	APIListObject
	Query    string   `url:"query,omitempty"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	Includes []string `url:"include,omitempty,brackets"`
}

// GetUserOptions is the data structure used when calling the GetUser API endpoint.
type GetUserOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// ListUsers lists users of your PagerDuty account, optionally filtered by a search query.
func (c *Client) ListUsers(o ListUsersOptions) (*ListUsersResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/users?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListUsersResponse
	return &result, c.decodeJSON(resp, &result)
}

// CreateUser creates a new user.
func (c *Client) CreateUser(u User) error {
	data := make(map[string]User)
	data["user"] = u
	_, err := c.post("/users", data)
	return err
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(id string) error {
	_, err := c.delete("/users/" + id)
	return err
}

// GetUser gets details about an existing user.
func (c *Client) GetUser(id string, o GetUserOptions) (*User, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/users/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]User
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}
	u, ok := result["user"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have user field")
	}
	return &u, nil
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(u User) error {
	v := make(map[string]User)
	v["user"] = u
	_, err := c.put("/users/"+u.ID, v)
	return err
}

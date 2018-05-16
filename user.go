package pagerduty

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
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
	Name              string `json:"name"`
	Email             string `json:"email"`
	Timezone          string `json:"timezone,omitempty"`
	Color             string `json:"color,omitempty"`
	Role              string `json:"role,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	Description       string `json:"description,omitempty"`
	InvitationSent    bool
	ContactMethods    []ContactMethod    `json:"contact_methods"`
	NotificationRules []NotificationRule `json:"notification_rules"`
	JobTitle          string             `json:"job_title,omitempty"`
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
func (pd *PagerdutyClient) ListUsers(o ListUsersOptions) (*ListUsersResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := pd.Get("/users?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListUsersResponse
	return &result, DecodeJSON(resp, &result)
}

// CreateUser creates a new user.
func (pd *PagerdutyClient) CreateUser(u User) (*User, error) {
	data := make(map[string]User)
	data["user"] = u
	resp, err := pd.Post("/users", data)
	return getUserFromResponse(pd, resp, err)
}

// DeleteUser deletes a user.
func (pd *PagerdutyClient) DeleteUser(id string) error {
	_, err := pd.Delete("/users/" + id)
	return err
}

// GetUser gets details about an existing user.
func (pd *PagerdutyClient) GetUser(id string, o GetUserOptions) (*User, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := pd.Get("/users/" + id + "?" + v.Encode())
	return getUserFromResponse(pd, resp, err)
}

// UpdateUser updates an existing user.
func (pd *PagerdutyClient) UpdateUser(u User) (*User, error) {
	v := make(map[string]User)
	v["user"] = u
	resp, err := pd.Put("/users/"+u.ID, v, nil)
	return getUserFromResponse(pd, resp, err)
}

func getUserFromResponse(pd *PagerdutyClient, resp *http.Response, err error) (*User, error) {
	if err != nil {
		return nil, err
	}
	var target map[string]User
	if dErr := DecodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "user"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}

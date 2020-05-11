package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// NotificationRule is a rule for notifying the user.
type NotificationRule struct {
	ID                  string        `json:"id"`
	StartDelayInMinutes uint          `json:"start_delay_in_minutes"`
	CreatedAt           string        `json:"created_at"`
	ContactMethod       ContactMethod `json:"contact_method"`
	Urgency             string        `json:"urgency"`
	Type                string        `json:"type"`
}

// User is a member of a PagerDuty account that has the ability to interact with incidents and other data on the account.
type User struct {
	APIObject
	Type              string             `json:"type"`
	Name              string             `json:"name"`
	Summary           string             `json:"summary"`
	Email             string             `json:"email"`
	Timezone          string             `json:"time_zone,omitempty"`
	Color             string             `json:"color,omitempty"`
	Role              string             `json:"role,omitempty"`
	AvatarURL         string             `json:"avatar_url,omitempty"`
	Description       string             `json:"description,omitempty"`
	InvitationSent    bool               `json:"invitation_sent,omitempty"`
	ContactMethods    []ContactMethod    `json:"contact_methods"`
	NotificationRules []NotificationRule `json:"notification_rules"`
	JobTitle          string             `json:"job_title,omitempty"`
	Teams             []Team
}

// ContactMethod is a way of contacting the user.
type ContactMethod struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Summary        string `json:"summary"`
	Self           string `json:"self"`
	Label          string `json:"label"`
	Address        string `json:"address"`
	SendShortEmail bool   `json:"send_short_email,omitempty"`
	SendHTMLEmail  bool   `json:"send_html_email,omitempty"`
	Blacklisted    bool   `json:"blacklisted,omitempty"`
	CountryCode    int    `json:"country_code,omitempty"`
	Enabled        bool   `json:"enabled,omitempty"`
	HTMLUrl        string `json:"html_url"`
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

// ListContactMethodsResponse is the data structure returned from calling the GetUserContactMethod API endpoint.
type ListContactMethodsResponse struct {
	APIListObject
	ContactMethods []ContactMethod `json:"contact_methods"`
}

// ListUserNotificationRulesResponse the data structure returned from calling the ListNotificationRules API endpoint.
type ListUserNotificationRulesResponse struct {
	APIListObject
	NotificationRules []NotificationRule `json:"notification_rules"`
}

// GetUserOptions is the data structure used when calling the GetUser API endpoint.
type GetUserOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetCurrentUserOptions is the data structure used when calling the GetCurrentUser API endpoint.
type GetCurrentUserOptions struct {
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
func (c *Client) CreateUser(u User) (*User, error) {
	data := make(map[string]User)
	data["user"] = u
	resp, err := c.post("/users", data, nil)
	return getUserFromResponse(c, resp, err)
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
	return getUserFromResponse(c, resp, err)
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(u User) (*User, error) {
	v := make(map[string]User)
	v["user"] = u
	resp, err := c.put("/users/"+u.ID, v, nil)
	return getUserFromResponse(c, resp, err)
}

// GetCurrentUser gets details about the authenticated user when using a user-level API key or OAuth token
func (c *Client) GetCurrentUser(o GetCurrentUserOptions) (*User, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/users/me?" + v.Encode())
	return getUserFromResponse(c, resp, err)
}

func getUserFromResponse(c *Client, resp *http.Response, err error) (*User, error) {
	if err != nil {
		return nil, err
	}
	var target map[string]User
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "user"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}

// ListUserContactMethods fetches contact methods of the existing user.
func (c *Client) ListUserContactMethods(userID string) (*ListContactMethodsResponse, error) {
	resp, err := c.get("/users/" + userID + "/contact_methods")
	if err != nil {
		return nil, err
	}
	var result ListContactMethodsResponse
	return &result, c.decodeJSON(resp, &result)
}

// GetUserContactMethod gets details about a contact method.
func (c *Client) GetUserContactMethod(userID, contactMethodID string) (*ContactMethod, error) {
	resp, err := c.get("/users/" + userID + "/contact_methods/" + contactMethodID)
	return getContactMethodFromResponse(c, resp, err)
}

// DeleteUserContactMethod deletes a user.
func (c *Client) DeleteUserContactMethod(userID, contactMethodID string) error {
	_, err := c.delete("/users/" + userID + "/contact_methods/" + contactMethodID)
	return err
}

// CreateUserContactMethod creates a new contact method for user.
func (c *Client) CreateUserContactMethod(userID string, cm ContactMethod) (*ContactMethod, error) {
	data := make(map[string]ContactMethod)
	data["contact_method"] = cm
	resp, err := c.post("/users/"+userID+"/contact_methods", data, nil)
	return getContactMethodFromResponse(c, resp, err)
}

// UpdateUserContactMethod updates an existing user.
func (c *Client) UpdateUserContactMethod(userID string, cm ContactMethod) (*ContactMethod, error) {
	v := make(map[string]ContactMethod)
	v["contact_method"] = cm
	resp, err := c.put("/users/"+userID+"/contact_methods/"+cm.ID, v, nil)
	return getContactMethodFromResponse(c, resp, err)
}

func getContactMethodFromResponse(c *Client, resp *http.Response, err error) (*ContactMethod, error) {
	if err != nil {
		return nil, err
	}
	var target map[string]ContactMethod
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "contact_method"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}

// GetUserNotificationRule gets details about a notification rule.
func (c *Client) GetUserNotificationRule(userID, ruleID string) (*NotificationRule, error) {
	resp, err := c.get("/users/" + userID + "/notification_rules/" + ruleID)
	return getUserNotificationRuleFromResponse(c, resp, err)
}

// CreateUserNotificationRule creates a new notification rule for a user.
func (c *Client) CreateUserNotificationRule(userID string, rule NotificationRule) (*NotificationRule, error) {
	data := make(map[string]NotificationRule)
	data["notification_rule"] = rule
	resp, err := c.post("/users/"+userID+"/notification_rules", data, nil)
	return getUserNotificationRuleFromResponse(c, resp, err)
}

// UpdateUserNotificationRule updates a notification rule for a user.
func (c *Client) UpdateUserNotificationRule(userID string, rule NotificationRule) (*NotificationRule, error) {
	data := make(map[string]NotificationRule)
	data["notification_rule"] = rule
	resp, err := c.put("/users/"+userID+"/notification_rules/"+rule.ID, data, nil)
	return getUserNotificationRuleFromResponse(c, resp, err)
}

// DeleteUserNotificationRule deletes a notification rule for a user.
func (c *Client) DeleteUserNotificationRule(userID, ruleID string) error {
	_, err := c.delete("/users/" + userID + "/notification_rules/" + ruleID)
	return err
}

// ListUserNotificationRules fetches notification rules of the existing user.
func (c *Client) ListUserNotificationRules(userID string) (*ListUserNotificationRulesResponse, error) {
	resp, err := c.get("/users/" + userID + "/notification_rules")
	if err != nil {
		return nil, err
	}
	var result ListUserNotificationRulesResponse
	return &result, c.decodeJSON(resp, &result)
}

func getUserNotificationRuleFromResponse(c *Client, resp *http.Response, err error) (*NotificationRule, error) {
	if err != nil {
		return nil, err
	}
	var target map[string]NotificationRule
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "notification_rule"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}

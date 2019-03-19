package pagerduty

import (
	"fmt"
	"net/http"
)

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

// ListContactMethodsResponse is the data structure returned from calling the ListContactMethods API endpoint.
type ListContactMethodsResponse struct {
	APIListObject
	ContactMethods []ContactMethod
}

// ListContactMethods lists all contact methods for a user.
// TODO: Unify with `ListUserContactMethods`.
func (c *Client) ListContactMethods(userID string) (*ListContactMethodsResponse, error) {
	resp, err := c.get("/users/" + userID + "/contact_methods")
	if err != nil {
		return nil, err
	}
	var result ListContactMethodsResponse
	return &result, c.decodeJSON(resp, &result)
}

// GetContactMethod gets details about a contact method.
func (c *Client) GetContactMethod(userID, id string) (*ContactMethod, error) {
	resp, err := c.get("/users/" + userID + "/contact_methods/" + id)
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

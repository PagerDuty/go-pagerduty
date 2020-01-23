package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Team is a collection of users and escalation policies that represent a group of people within an organization.
type Team struct {
	APIObject
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ListTeamResponse is the structure used when calling the ListTeams API endpoint.
type ListTeamResponse struct {
	APIListObject
	Teams []Team
}

// ListTeamOptions are the input parameters used when calling the ListTeams API endpoint.
type ListTeamOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

// ListTeams lists teams of your PagerDuty account, optionally filtered by a search query.
func (c *Client) ListTeams(o ListTeamOptions) (*ListTeamResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get("/teams?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListTeamResponse
	return &result, c.decodeJSON(resp, &result)
}

// CreateTeam creates a new team.
func (c *Client) CreateTeam(t *Team) (*Team, error) {
	resp, err := c.post("/teams", t, nil)
	return getTeamFromResponse(c, resp, err)
}

// DeleteTeam removes an existing team.
func (c *Client) DeleteTeam(id string) error {
	_, err := c.delete("/teams/" + id)
	return err
}

// GetTeam gets details about an existing team.
func (c *Client) GetTeam(id string) (*Team, error) {
	resp, err := c.get("/teams/" + id)
	return getTeamFromResponse(c, resp, err)
}

// UpdateTeam updates an existing team.
func (c *Client) UpdateTeam(id string, t *Team) (*Team, error) {
	resp, err := c.put("/teams/"+id, t, nil)
	return getTeamFromResponse(c, resp, err)
}

// RemoveEscalationPolicyFromTeam removes an escalation policy from a team.
func (c *Client) RemoveEscalationPolicyFromTeam(teamID, epID string) error {
	_, err := c.delete("/teams/" + teamID + "/escalation_policies/" + epID)
	return err
}

// AddEscalationPolicyToTeam adds an escalation policy to a team.
func (c *Client) AddEscalationPolicyToTeam(teamID, epID string) error {
	_, err := c.put("/teams/"+teamID+"/escalation_policies/"+epID, nil, nil)
	return err
}

// RemoveUserFromTeam removes a user from a team.
func (c *Client) RemoveUserFromTeam(teamID, userID string) error {
	_, err := c.delete("/teams/" + teamID + "/users/" + userID)
	return err
}

// AddUserToTeam adds a user to a team.
func (c *Client) AddUserToTeam(teamID, userID string) error {
	_, err := c.put("/teams/"+teamID+"/users/"+userID, nil, nil)
	return err
}

func getTeamFromResponse(c *Client, resp *http.Response, err error) (*Team, error) {
	if err != nil {
		return nil, err
	}
	var target map[string]Team
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "team"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}

// Member is a team member.
type Member struct {
	APIObject struct {
		APIObject
	} `json:"user"`
	Role string `json:"role"`
}

// ListMembersOptions are the optional parameters for a members request.
type ListMembersOptions struct {
	APIListObject
}

// ListMembersResponse is the response from the members endpoint.
type ListMembersResponse struct {
	APIListObject
	Members []Member `json:"members"`
}

// ListMembers gets the first page of users associated with the specified team.
func (c *Client) ListMembers(teamID string, o ListMembersOptions) (*ListMembersResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get("/teams/" + teamID + "/members?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListMembersResponse
	return &result, c.decodeJSON(resp, &result)
}

// ListAllMembers gets all members associated with the specified team.
func (c *Client) ListAllMembers(teamID string) ([]Member, error) {
	members := make([]Member, 0)

	// Create a handler closure capable of parsing data from the members endpoint
	// and appending resultant members to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListMembersResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		members = append(members, result.Members...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet("/teams/"+teamID+"/members", responseHandler); err != nil {
		return nil, err
	}

	return members, nil
}

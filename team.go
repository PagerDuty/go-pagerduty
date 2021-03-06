package pagerduty

import (
	"context"
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

// ListTeams lists teams of your PagerDuty account, optionally filtered by a
// search query. It's recommended to use ListTeamsWithContext instead.
func (c *Client) ListTeams(o ListTeamOptions) (*ListTeamResponse, error) {
	return c.ListTeamsWithContext(context.Background(), o)
}

// ListTeamsWithContext lists teams of your PagerDuty account, optionally
// filtered by a search query.
func (c *Client) ListTeamsWithContext(ctx context.Context, o ListTeamOptions) (*ListTeamResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/teams?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListTeamResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateTeam creates a new team. It's recommended to use CreateTeamWithContext
// instead.
func (c *Client) CreateTeam(t *Team) (*Team, error) {
	return c.CreateTeamWithContext(context.Background(), t)
}

// CreateTeamWithContext creates a new team.
func (c *Client) CreateTeamWithContext(ctx context.Context, t *Team) (*Team, error) {
	resp, err := c.post(ctx, "/teams", t, nil)
	return getTeamFromResponse(c, resp, err)
}

// DeleteTeam removes an existing team. It's recommended to use
// DeleteTeamWithContext instead.
func (c *Client) DeleteTeam(id string) error {
	return c.DeleteTeamWithContext(context.Background(), id)
}

// DeleteTeamWithContext removes an existing team.
func (c *Client) DeleteTeamWithContext(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/teams/"+id)
	return err
}

// GetTeam gets details about an existing team. It's recommended to use
// GetTeamWithContext instead.
func (c *Client) GetTeam(id string) (*Team, error) {
	return c.GetTeamWithContext(context.Background(), id)
}

// GetTeamWithContext gets details about an existing team.
func (c *Client) GetTeamWithContext(ctx context.Context, id string) (*Team, error) {
	resp, err := c.get(ctx, "/teams/"+id)
	return getTeamFromResponse(c, resp, err)
}

// UpdateTeam updates an existing team. It's recommended to use
// UpdateTeamWithContext instead.
func (c *Client) UpdateTeam(id string, t *Team) (*Team, error) {
	return c.UpdateTeamWithContext(context.Background(), id, t)
}

// UpdateTeamWithContext updates an existing team.
func (c *Client) UpdateTeamWithContext(ctx context.Context, id string, t *Team) (*Team, error) {
	resp, err := c.put(ctx, "/teams/"+id, t, nil)
	return getTeamFromResponse(c, resp, err)
}

// RemoveEscalationPolicyFromTeam removes an escalation policy from a team. It's
// recommended to use RemoveEscalationPolicyFromTeamWithContext instead.
func (c *Client) RemoveEscalationPolicyFromTeam(teamID, epID string) error {
	return c.RemoveEscalationPolicyFromTeamWithContext(context.Background(), teamID, epID)
}

// RemoveEscalationPolicyFromTeamWithContext removes an escalation policy from a team.
func (c *Client) RemoveEscalationPolicyFromTeamWithContext(ctx context.Context, teamID, epID string) error {
	_, err := c.delete(ctx, "/teams/"+teamID+"/escalation_policies/"+epID)
	return err
}

// AddEscalationPolicyToTeam adds an escalation policy to a team. It's
// recommended to use AddEscalationPolicyToTeamWithContext instead.
func (c *Client) AddEscalationPolicyToTeam(teamID, epID string) error {
	return c.AddEscalationPolicyToTeamWithContext(context.Background(), teamID, epID)
}

// AddEscalationPolicyToTeamWithContext adds an escalation policy to a team.
func (c *Client) AddEscalationPolicyToTeamWithContext(ctx context.Context, teamID, epID string) error {
	_, err := c.put(ctx, "/teams/"+teamID+"/escalation_policies/"+epID, nil, nil)
	return err
}

// RemoveUserFromTeam removes a user from a team. It's recommended to use
// RemoveUserFromTeamWithContext instead.
func (c *Client) RemoveUserFromTeam(teamID, userID string) error {
	return c.RemoveUserFromTeamWithContext(context.Background(), teamID, userID)
}

// RemoveUserFromTeamWithContext removes a user from a team.
func (c *Client) RemoveUserFromTeamWithContext(ctx context.Context, teamID, userID string) error {
	_, err := c.delete(ctx, "/teams/"+teamID+"/users/"+userID)
	return err
}

// AddUserToTeam adds a user to a team.
func (c *Client) AddUserToTeam(teamID, userID string) error {
	return c.AddUserToTeamWithContext(context.Background(), AddUserToTeamOptions{TeamID: teamID, UserID: userID})
}

// TeamUserRole is a named type to represent the different Team Roles supported
// by PagerDuty when adding a user to a team.
//
// For more info: https://support.pagerduty.com/docs/advanced-permissions#team-roles
type TeamUserRole string

const (
	// TeamUserRoleObserver is the obesrver team role, which generally provides
	// read-only access. They gain responder-level permissions on an incident if
	// one is assigned to them.
	TeamUserRoleObserver TeamUserRole = "observer"

	// TeamUserRoleResponder is the responder team role, and they are given the
	// same permissions as the observer plus the ability to respond to
	// incidents, trigger incidents, and manage overrides.
	TeamUserRoleResponder TeamUserRole = "responder"

	// TeamUserRoleManager is the manager team role, and they are given the same
	// permissions as the responder plus the ability to edit and delete the
	// different resources owned by the team.
	TeamUserRoleManager TeamUserRole = "manager"
)

// AddUserToTeamOptions is an option struct for the AddUserToTeamWithContext
// method.
type AddUserToTeamOptions struct {
	TeamID string       `json:"-"`
	UserID string       `json:"-"`
	Role   TeamUserRole `json:"role,omitempty"`
}

// AddUserToTeamWithContext adds a user to a team.
func (c *Client) AddUserToTeamWithContext(ctx context.Context, o AddUserToTeamOptions) error {
	_, err := c.put(ctx, "/teams/"+o.TeamID+"/users/"+o.UserID, o, nil)
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

	const rootNode = "team"

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
// It's recommended to use ListMembersWithContext instead.
func (c *Client) ListMembers(teamID string, o ListMembersOptions) (*ListMembersResponse, error) {
	return c.ListMembersWithContext(context.Background(), teamID, o)
}

// ListMembersWithContext gets the first page of users associated with the specified team.
func (c *Client) ListMembersWithContext(ctx context.Context, teamID string, o ListMembersOptions) (*ListMembersResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/teams/"+teamID+"/members?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListMembersResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListAllMembers gets all members associated with the specified team. It's
// recommended to use ListMembersPaginated instead.
func (c *Client) ListAllMembers(teamID string) ([]Member, error) {
	return c.ListMembersPaginated(context.Background(), teamID)
}

// ListMembersPaginated gets all members associated with the specified team.
func (c *Client) ListMembersPaginated(ctx context.Context, teamID string) ([]Member, error) {
	var members []Member

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
	if err := c.pagedGet(ctx, "/teams/"+teamID+"/members", responseHandler); err != nil {
		return nil, err
	}

	return members, nil
}

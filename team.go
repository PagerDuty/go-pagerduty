package pagerduty

import (
	"github.com/google/go-querystring/query"
)

type Team struct {
	APIObject
	Name        string
	Description string
}

type ListTeamResponse struct {
	APIListObject
	Teams []Team
}

type ListTeamOptions struct {
	Query string `url:"query,omitempty"`
}

func (c *Client) ListTeams(o ListTeamOptions) (*ListTeamResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.Get("/teams?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListTeamResponse
	return &result, c.decodeJson(resp, &result)
}

func (c *Client) CreateTeam(t *Team) error {
	return nil
}

func (c *Client) DeleteTeam(id string) error {
	return nil
}

func (c *Client) GetTeam(id string) error {
	return nil
}

func (c *Client) UpdateTeam(t *Team) error {
	return nil
}

func (c *Client) RmoveEscalationPolicyFromTeam() error {
	return nil
}

func (c *Client) AddEscalationPolicyToTeam() error {
	return nil
}

func (c *Client) RemoveUserFromTeam() error {
	return nil
}

func (c *Client) AddUserToTeam() error {
	return nil
}

package pagerduty

import (
	"fmt"
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
	_, err := c.Post("/teams", t)
	return err
}

func (c *Client) DeleteTeam(id string) error {
	_, err := c.Delete("/teams/" + id)
	return err
}

func (c *Client) GetTeam(id string) (*Team, error) {
	resp, err := c.Get("/teams/" + id)
	if err != nil {
		return nil, err
	}
	var result map[string]Team
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	t, ok := result["team"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have team field")
	}
	return &t, nil
}

func (c *Client) UpdateTeam(t *Team) error {
	//TODO
	return nil
}

func (c *Client) RmoveEscalationPolicyFromTeam() error {
	//TODO
	return nil
}

func (c *Client) AddEscalationPolicyToTeam() error {
	//TODO
	return nil
}

func (c *Client) RemoveUserFromTeam() error {
	return nil
}

func (c *Client) AddUserToTeam() error {
	return nil
}

package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

type EscalationRule struct {
	Id      string `json:"id"`
	Delay   uint   `json:"escalation_delay_in_minutes"`
	Targets []APIObject
}

type EscalationPolicy struct {
	APIObject
	Name            string      `json:"name,omitempty"`
	EscalationRules []APIObject `json:"escalation_rules,omitempty"`
	Services        []APIObject `json:"services,omitempty"`
	NumLoops        uint        `json:"num_loops,omitempty"`
	Teams           []APIObject `json:"teams,omitempty"`
	Description     string      `json:"description,omitempty"`
}

type ListEscalationPolicyResponse struct {
	APIListObject
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
}

type ListEscalationPoliciesOptions struct {
	APIListObject
	Query    string   `url:"query,omitempty"`
	UserIDs  []string `url:"user_ids,omitempty,brackets"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	Includes []string `url:"include,omitempty,brackets"`
	SortBy   string   `url:"sort_by,omitempty"`
}

func (c *Client) ListEscalationPolicies(o ListEscalationPoliciesOptions) (*ListEscalationPolicyResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/escalation_policies?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListEscalationPolicyResponse
	return &result, c.decodeJson(resp, &result)
}

func (c *Client) CreateEscalationPolicy(ep EscalationPolicy) error {
	data := make(map[string]EscalationPolicy)
	data["escalation_policy"] = ep
	_, err := c.Post("/escalation_policies", data)
	return err
}

func (c *Client) DeleteEscalationPolicy(id string) error {
	_, err := c.Delete("/escalation_policies/" + id)
	return err
}

type GetEscalationPolicyOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

func (c *Client) GetEscalationPolicy(id string, o *GetEscalationPolicyOptions) (*EscalationPolicy, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/escalation_policies/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]EscalationPolicy
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	ep, ok := result["escalation_policy"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have escalation_policy field")
	}
	return &ep, nil
}

func (c *Client) UpdateEscalationPolicy(e *EscalationPolicy) error {
	//TODO
	return nil
}

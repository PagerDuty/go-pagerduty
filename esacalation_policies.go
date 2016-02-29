package pagerduty

import (
	"fmt"
	"github.com/google/go-querystring/query"
)

// APIObject represents generic api json response that is shared by most
// domain object (like escalation
type APIObject struct {
	ID      string `json:"id,omitempty"`
	Type    string
	Summary string
	Self    string `json:"omitempty",yaml:"omitempty"`
	HtmlUrl string `json:"html_url,omitempty"`
}

type APIListObject struct {
	Limit  uint
	Offset uint
	More   bool
	Total  uint
}

type EscalationRule struct {
	Id      string `json:"id"`
	Delay   uint   `json:"escalation_delay_in_minutes"`
	Targets []APIObject
}

type EscalationPolicy struct {
	APIObject
	Name            string
	EscalationRules []APIObject `json:"escalation_rules"`
	Services        []APIObject
	NumLoops        uint `json:"num_loops"`
	Teams           []APIObject
	Description     string
}

type ListEscalationPolicyResponse struct {
	APIListObject
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
}

type ListEscalationPoliciesOptions struct {
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

func (c *Client) CreateEscalationPolicy(ep *EscalationPolicy) error {
	_, err := c.Post("/escalation_policies", ep)
	return err
}

func DeleteEscalationPolicy() {}

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
func UpdateEscalationPolicy() {}

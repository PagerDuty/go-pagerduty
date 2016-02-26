package pagerduty

import (
	"github.com/google/go-querystring/query"
)

// APIObject represents generic api json response that is shared by most
// domain object (like escalation
type APIObject struct {
	ID      string `json:"id"`
	Type    string
	Summary string
	Self    string
	HtmlUrl string `json:"html_url"`
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
	Limit              uint
	Offset             uint
	More               bool
	Total              uint
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
}

type ListEscalationPoliciesOptions struct {
	Query   string   `url:"query,omitempty"`
	UserIDs []string `url:"user_ids,omitempty"`
	TeamIDs []string `url:"team_ids,omitempty"`
	Include []string `url:"include,omitempty"`
	SortBy  string   `url:"sort_by,omitempty"`
}

func (c *Client) ListEscalationPolicies(opts ListEscalationPoliciesOptions) (*ListEscalationPolicyResponse, error) {
	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do("GET", "/escalation_policies"+v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListEscalationPolicyResponse
	return &result, c.decodeJson(resp, &result)
}

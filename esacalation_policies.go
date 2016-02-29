package pagerduty

import (
	"github.com/google/go-querystring/query"
)

// APIObject represents generic api json response that is shared by most
// domain object (like escalation
type APIObject struct {
	ID      string `json:"id,omitempty"`
	Type    string
	Summary string
	Self    string `json:"omitempty"`
	HtmlUrl string `json:"html_url,omitempty"`
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
	Query    string   `url:"query,omitempty"`
	UserIDs  []string `url:"user_ids,omitempty,brackets"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	Includes []string `url:"include,omitempty,brackets"`
	SortBy   string   `url:"sort_by,omitempty"`
}

func (c *Client) ListEscalationPolicies(opts ListEscalationPoliciesOptions) (*ListEscalationPolicyResponse, error) {
	v, err := query.Values(opts)
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
func GetEscalationPolicy()    {}
func UpdateEscalationPolicy() {}

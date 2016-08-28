package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// EscalationRule is a rule for an escalation policy to trigger.
type EscalationRule struct {
	ID      string `json:"id"`
	Delay   uint   `json:"escalation_delay_in_minutes"`
	Targets []APIObject
}

// EscalationPolicy is a collection of escalation rules.
type EscalationPolicy struct {
	APIObject
	Name            string      `json:"name,omitempty"`
	EscalationRules []APIObject `json:"escalation_rules,omitempty"`
	Services        []APIObject `json:"services,omitempty"`
	NumLoops        uint        `json:"num_loops,omitempty"`
	Teams           []APIObject `json:"teams,omitempty"`
	Description     string      `json:"description,omitempty"`
}

// ListEscalationPoliciesResponse is the data structure returned from calling the ListEscalationPolicies API endpoint.
type ListEscalationPoliciesResponse struct {
	APIListObject
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
}

// ListEscalationPoliciesOptions is the data structure used when calling the ListEscalationPolicies API endpoint.
type ListEscalationPoliciesOptions struct {
	APIListObject
	Query    string   `url:"query,omitempty"`
	UserIDs  []string `url:"user_ids,omitempty,brackets"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	Includes []string `url:"include,omitempty,brackets"`
	SortBy   string   `url:"sort_by,omitempty"`
}

// ListEscalationPolicies lists all of the existing escalation policies.
func (c *Client) ListEscalationPolicies(o ListEscalationPoliciesOptions) (*ListEscalationPoliciesResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/escalation_policies?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListEscalationPoliciesResponse
	return &result, c.decodeJSON(resp, &result)
}

// CreateEscalationPolicy creates a new escalation policy.
func (c *Client) CreateEscalationPolicy(ep EscalationPolicy) (*EscalationPolicy, error) {
	data := make(map[string]EscalationPolicy)
	data["escalation_policy"] = ep
	resp, err := c.post("/escalation_policies", data)
	return decodeEscalationPolicyFromResponse(c, resp, err)
}

// DeleteEscalationPolicy deletes an existing escalation policy and rules.
func (c *Client) DeleteEscalationPolicy(id string) error {
	_, err := c.delete("/escalation_policies/" + id)
	return err
}

// GetEscalationPolicyOptions is the data structure used when calling the GetEscalationPolicy API endpoint.
type GetEscalationPolicyOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetEscalationPolicy gets information about an existing escalation policy and its rules.
func (c *Client) GetEscalationPolicy(id string, o *GetEscalationPolicyOptions) (*EscalationPolicy, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/escalation_policies/" + id + "?" + v.Encode())
	return decodeEscalationPolicyFromResponse(c, resp, err)
}

// UpdateEscalationPolicy updates an existing escalation policy and its rules.
func (c *Client) UpdateEscalationPolicy(id string, e *EscalationPolicy) (*EscalationPolicy, error) {
	resp, err := c.put("/escalation_policies/"+id, e)
	return decodeEscalationPolicyFromResponse(c, resp, err)
}

func decodeEscalationPolicyFromResponse(c *Client, resp *http.Response, err error) (*EscalationPolicy, error) {
	if err != nil {
		return nil, err
	}
	var result map[string]EscalationPolicy
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}
	t, ok := result["escalation_policy"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have escalation_policy field")
	}
	return &t, nil
}

package pagerduty

import (
	"context"
)

// PriorityProperty is a single priorty object returned from the Priorities endpoint
type PriorityProperty struct {
	APIObject
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Priorities repreents the API response from PagerDuty when listing the
// configured priorities.
type Priorities struct {
	APIListObject
	Priorities []PriorityProperty `json:"priorities"`
}

// ListPriorities lists existing priorities. It's recommended to use
// ListPrioritiesWithContext instead.
func (c *Client) ListPriorities() (*Priorities, error) {
	return c.ListPrioritiesWithContext(context.Background())
}

// ListPrioritiesWithContext lists existing priorities.
func (c *Client) ListPrioritiesWithContext(ctx context.Context) (*Priorities, error) {
	resp, err := c.get(ctx, "/priorities")
	if err != nil {
		return nil, err
	}

	var p Priorities
	if err := c.decodeJSON(resp, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

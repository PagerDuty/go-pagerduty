package pagerduty

import (
	"context"
	"encoding/json"
)

// PriorityProperty is a single priorty object returned from the Priorities endpoint
type PriorityProperty struct {
	APIObject
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Priorities struct {
	APIListObject
	Priorities []PriorityProperty `json:"priorities"`
}

// ListPriorities lists existing priorities
func (c *Client) ListPriorities() (*Priorities, error) {
	resp, err := c.get(context.TODO(), "/priorities")
	if err != nil {
		return nil, err
	}

	// TODO(theckman): make sure we close the resp.Body here

	var p Priorities
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

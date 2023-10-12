package pagerduty

import (
	"context"

	"github.com/google/go-querystring/query"
)

const (
	standardPath = "/standards"
)

// Standard defines a PagerDuty's resource standard.
type Standard struct {
	Active       bool                         `json:"active"`
	Description  string                       `json:"description,omitempty"`
	Exclusions   []StandardInclusionExclusion `json:"exclusions,omitempty"`
	ID           string                       `json:"id,omitempty"`
	Inclusions   []StandardInclusionExclusion `json:"inclusions,omitempty"`
	Name         string                       `json:"name,omitempty"`
	ResourceType string                       `json:"resource_type,omitempty"`
	Type         string                       `json:"type,omitempty"`
}

type StandardInclusionExclusion struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// ListStandardsResponse is the data structure returned from calling the ListStandards API endpoint.
type ListStandardsResponse struct {
	Standards []Standard `json:"standards"`
}

// ListStandardsOptions is the data structure used when calling the ListStandards API endpoint.
type ListStandardsOptions struct {
	Active bool `url:"active,omitempty"`

	// ResourceType query for a specific resource type.
	//  Allowed value: technical_service
	ResourceType string `url:"resource_type,omitempty"`
}

// ListStandardsWithContext lists all the existing standards.
func (c *Client) ListStandardsWithContext(ctx context.Context, o ListStandardsOptions) (*ListStandardsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, standardPath+"?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListStandardsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateStandardWithContext updates an existing standard.
func (c *Client) UpdateStandardWithContext(ctx context.Context, id string, s Standard) (*Standard, error) {
	resp, err := c.put(ctx, standardPath+"/"+id, s, nil)
	if err != nil {
		return nil, err
	}

	var result Standard
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

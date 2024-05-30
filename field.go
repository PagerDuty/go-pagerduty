package pagerduty

import (
	"context"
)

// FieldOptionData TODO
type FieldOptionData struct {
	DataType string `json:"data_type,omitempty"`
	Value    string `json:"value,omitempty"`
}

// FieldOption TODO
type FieldOption struct {
	APIReference
	Data      *FieldOptionData `json:"data,omitempty"`
	CreatedAt string           `json:"created_at,omitempty"`
	UpdatedAt string           `json:"updated_at,omitempty"`
}

// CreateFieldOptionResponse TODO
type FieldOptionResponse struct {
	FieldOption FieldOption `json:"field_option"`
}

// CreateIncidentFieldOptionWithContext creates a new custom field option for custom field with an ID equal to `fieldID`.
func (c *Client) CreateFieldOptionWithContext(ctx context.Context, fieldID string, fieldOption FieldOption) (*FieldOption, error) {
	if fieldOption.Type == "" {
		fieldOption.Type = "incident"
	}

	b := map[string]FieldOption{
		"field_option": fieldOption,
	}

	resp, err := c.post(ctx, "/incidents/custom_fields/"+fieldID+"/field_options", b, nil)
	if err != nil {
		return nil, err
	}

	var response FieldOptionResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.FieldOption, nil
}

// ListFieldOptionsResponse TODO
type ListFieldOptionsResponse struct {
	FieldOptions []FieldOption `json:"field_options"`
}

// ListFieldOptionWithContext TODO
func (c *Client) ListFieldOptionsWithContext(ctx context.Context, fieldID string) (*ListFieldOptionsResponse, error) {
	resp, err := c.get(ctx, "/incidents/custom_fields/"+fieldID+"/field_options", nil)
	if err != nil {
		return nil, err
	}

	var response ListFieldOptionsResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, err
}

// UpdateFieldOptionWithContext TODO
func (c *Client) UpdateFieldOptionWithContext(ctx context.Context, fieldID string, o FieldOption) (*FieldOption, error) {
	b := map[string]FieldOption{
		"field_option": o,
	}

	resp, err := c.put(ctx, "/incidents/custom_fields/"+fieldID+"/field_options"+o.ID, b, nil)
	if err != nil {
		return nil, err
	}

	var response FieldOptionResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.FieldOption, err
}

// DeleteFieldOptionWithContext TODO
func (c *Client) DeleteFieldOptionWithContext(ctx context.Context, fieldID, fieldOptionID string) error {
	_, err := c.delete(ctx, "/incidents/custom_fields/"+fieldID+"/field_options"+fieldOptionID)
	return err
}

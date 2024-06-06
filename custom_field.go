package pagerduty

import (
	"context"

	"github.com/google/go-querystring/query"
)

// CustomField TODO
type CustomField struct {
	APIObject
	Name         string              `json:"name"`
	DisplayName  string              `json:"display_name"`
	Description  string              `json:"description,omitempty"`
	DataType     string              `json:"data_type"`
	FieldType    string              `json:"field_type"`
	DefaultValue interface{}         `json:"default_value"`
	FieldOptions []CustomFieldOption `json:"field_options"`
}

// CustomFieldResponse TODO
type CustomFieldResponse struct {
	Field CustomField `json:"field"`
}

// CreateCustomFieldWithContext create a new Custom Field, along with the
// Custom Field Options if provided. An account may have up to 10 Fields.
func (c *Client) CreateCustomFieldWithContext(ctx context.Context, cf CustomField) (*CustomField, error) {
	b := map[string]CustomField{
		"field": cf,
	}

	resp, err := c.post(ctx, "/incidents/custom_fields", b, nil)
	if err != nil {
		return nil, err
	}

	var response CustomFieldResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Field, nil
}

// ListCustomFieldsOptions TODO
type ListCustomFieldsOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// ListCustomFieldsOptionResponse TODO
type ListCustomFieldsResponse struct {
	Fields []CustomField `json:"fields"`
}

// ListCustomFieldsWithContext shows a list of Custom Fields.
func (c *Client) ListCustomFieldsWithContext(ctx context.Context, o ListCustomFieldsOptions) (*ListCustomFieldsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents/custom_fields/?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var response ListCustomFieldsResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCustomFieldOptions TODO
type GetCustomFieldOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetCustomFieldWithContext shows detailed information about a Custom Field.
func (c *Client) GetCustomFieldWithContext(ctx context.Context, id string, o GetCustomFieldOptions) (*CustomField, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents/custom_fields/"+id+"?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var response CustomFieldResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Field, nil
}

// UpdateCustomFieldWithContext shows detailed information about a Custom Field.
func (c *Client) UpdateCustomFieldWithContext(ctx context.Context, cf CustomField) (*CustomField, error) {
	b := map[string]CustomField{
		"field": cf,
	}

	resp, err := c.put(ctx, "/incidents/custom_fields/"+cf.ID, b, nil)
	if err != nil {
		return nil, err
	}

	var response CustomFieldResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Field, nil
}

// DeleteCustomFieldWithContext deletes a Custom Field. Fields may not be
// deleted if they are used by a Field Schema.
func (c *Client) DeleteCustomFieldWithContext(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/incidents/custom_fields/"+id)
	return err
}

// CustomFieldOptionData TODO
type CustomFieldOptionData struct {
	DataType string `json:"data_type,omitempty"`
	Value    string `json:"value,omitempty"`
}

// CustomFieldOption TODO
type CustomFieldOption struct {
	APIReference
	Data      *CustomFieldOptionData `json:"data,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
}

// CreateCustomFieldOptionResponse TODO
type CustomFieldOptionResponse struct {
	FieldOption CustomFieldOption `json:"field_option"`
}

// CreateCustomFieldOptionWithContext creates a new custom field option for
// custom field with an ID equal to `fieldID`.
func (c *Client) CreateCustomFieldOptionWithContext(ctx context.Context, fieldID string, fieldOption CustomFieldOption) (*CustomFieldOption, error) {
	if fieldOption.Type == "" {
		fieldOption.Type = "incident"
	}

	b := map[string]CustomFieldOption{
		"field_option": fieldOption,
	}

	resp, err := c.post(ctx, "/incidents/custom_fields/"+fieldID+"/field_options", b, nil)
	if err != nil {
		return nil, err
	}

	var response CustomFieldOptionResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.FieldOption, nil
}

// ListCustomFieldOptionsResponse TODO
type ListCustomFieldOptionsResponse struct {
	FieldOptions []CustomFieldOption `json:"field_options"`
}

// ListCustomFieldOptionWithContext TODO
func (c *Client) ListCustomFieldOptionsWithContext(ctx context.Context, fieldID string) (*ListCustomFieldOptionsResponse, error) {
	resp, err := c.get(ctx, "/incidents/custom_fields/"+fieldID+"/field_options", nil)
	if err != nil {
		return nil, err
	}

	var response ListCustomFieldOptionsResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, err
}

// UpdateCustomFieldOptionWithContext TODO
func (c *Client) UpdateCustomFieldOptionWithContext(ctx context.Context, fieldID string, o CustomFieldOption) (*CustomFieldOption, error) {
	b := map[string]CustomFieldOption{
		"field_option": o,
	}

	resp, err := c.put(ctx, "/incidents/custom_fields/"+fieldID+"/field_options/"+o.ID, b, nil)
	if err != nil {
		return nil, err
	}

	var response CustomFieldOptionResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.FieldOption, err
}

// DeleteCustomFieldOptionWithContext TODO
func (c *Client) DeleteCustomFieldOptionWithContext(ctx context.Context, fieldID, fieldOptionID string) error {
	_, err := c.delete(ctx, "/incidents/custom_fields/"+fieldID+"/field_options/"+fieldOptionID)
	return err
}

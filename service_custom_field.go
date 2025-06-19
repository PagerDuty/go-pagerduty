package pagerduty

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

// ServiceCustomFieldDataType represents the data type of a custom field
type ServiceCustomFieldDataType string

const (
	ServiceCustomFieldDataTypeBoolean  ServiceCustomFieldDataType = "boolean"
	ServiceCustomFieldDataTypeInteger  ServiceCustomFieldDataType = "integer"
	ServiceCustomFieldDataTypeFloat    ServiceCustomFieldDataType = "float"
	ServiceCustomFieldDataTypeString   ServiceCustomFieldDataType = "string"
	ServiceCustomFieldDataTypeDatetime ServiceCustomFieldDataType = "datetime"
	ServiceCustomFieldDataTypeURL      ServiceCustomFieldDataType = "url"
)

// ServiceCustomFieldType represents the field type of a custom field
type ServiceCustomFieldType string

const (
	ServiceCustomFieldTypeSingleValue      ServiceCustomFieldType = "single_value"
	ServiceCustomFieldTypeSingleValueFixed ServiceCustomFieldType = "single_value_fixed"
	ServiceCustomFieldTypeMultiValue       ServiceCustomFieldType = "multi_value"
	ServiceCustomFieldTypeMultiValueFixed  ServiceCustomFieldType = "multi_value_fixed"
)

// ServiceCustomFieldDefaultValue represents the default value of a custom field
type ServiceCustomFieldDefaultValue struct {
	Value interface{} `json:"value,omitempty"`
}

// ServiceCustomFieldOptionData represents the data of a field option
type ServiceCustomFieldOptionData struct {
	DataType ServiceCustomFieldDataType `json:"data_type"`
	Value    string                     `json:"value"`
}

// ServiceCustomFieldOption represents an option for a custom field
type ServiceCustomFieldOption struct {
	ID        string                       `json:"id,omitempty"`
	Type      string                       `json:"type,omitempty"`
	CreatedAt string                       `json:"created_at,omitempty"`
	UpdatedAt string                       `json:"updated_at,omitempty"`
	Data      ServiceCustomFieldOptionData `json:"data"`
}

// ServiceCustomField represents a custom field for services
type ServiceCustomField struct {
	APIObject
	Name         string                     `json:"name"`
	DisplayName  string                     `json:"display_name"`
	Description  string                     `json:"description,omitempty"`
	DataType     ServiceCustomFieldDataType `json:"data_type"`
	FieldType    ServiceCustomFieldType     `json:"field_type"`
	DefaultValue interface{}                `json:"default_value,omitempty"`
	Enabled      bool                       `json:"enabled"`
	FieldOptions []ServiceCustomFieldOption `json:"field_options,omitempty"`
	CreatedAt    string                     `json:"created_at,omitempty"`
	UpdatedAt    string                     `json:"updated_at,omitempty"`
}

// ListServiceCustomFieldsResponse is the data structure returned from calling the ListServiceCustomFields API endpoint
type ListServiceCustomFieldsResponse struct {
	Fields []ServiceCustomField `json:"fields"`
}

// ListServiceCustomFieldsOptions is the data structure used when calling the ListServiceCustomFields API endpoint
type ListServiceCustomFieldsOptions struct {
	Include []string `url:"include[],omitempty"`
}

// ListServiceCustomFields lists all custom fields for services
func (c *Client) ListServiceCustomFields(ctx context.Context, o ListServiceCustomFieldsOptions) (*ListServiceCustomFieldsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	resp, err := c.get(ctx, "/services/custom_fields?"+v.Encode(), headers)
	if err != nil {
		return nil, err
	}

	var result ListServiceCustomFieldsResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetServiceCustomField gets a specific custom field for services
func (c *Client) GetServiceCustomField(ctx context.Context, id string, o ListServiceCustomFieldsOptions) (*ServiceCustomField, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	resp, err := c.get(ctx, "/services/custom_fields/"+id+"?"+v.Encode(), headers)
	if err != nil {
		return nil, err
	}

	var result struct {
		Field ServiceCustomField `json:"field"`
	}
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Field, nil
}

// CreateServiceCustomField creates a new custom field for services
func (c *Client) CreateServiceCustomField(ctx context.Context, field *ServiceCustomField) (*ServiceCustomField, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	d := map[string]interface{}{
		"field": field,
	}

	resp, err := c.post(ctx, "/services/custom_fields", d, headers)
	if err != nil {
		return nil, err
	}

	var result struct {
		Field ServiceCustomField `json:"field"`
	}
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result.Field, nil
}

// UpdateServiceCustomField updates an existing custom field for services
func (c *Client) UpdateServiceCustomField(ctx context.Context, field *ServiceCustomField) (*ServiceCustomField, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	d := map[string]interface{}{
		"field": field,
	}

	resp, err := c.put(ctx, "/services/custom_fields/"+field.ID, d, headers)
	if err != nil {
		return nil, err
	}

	var result struct {
		Field ServiceCustomField `json:"field"`
	}
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result.Field, nil
}

// DeleteServiceCustomField deletes a custom field for services
func (c *Client) DeleteServiceCustomField(ctx context.Context, id string) error {
	// Set the required X-EARLY-ACCESS header for this API
	headers := map[string]string{"X-EARLY-ACCESS": "service-custom-fields-preview"}
	_, err := c.deleteWithHeaders(ctx, "/services/custom_fields/"+id, headers)
	return err
}

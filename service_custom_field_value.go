package pagerduty

import (
	"context"
	"fmt"
)

// ServiceCustomFieldValue represents a custom field value for a service
type ServiceCustomFieldValue struct {
	ID          string                     `json:"id,omitempty"`
	Name        string                     `json:"name,omitempty"`
	DisplayName string                     `json:"display_name,omitempty"`
	Description string                     `json:"description,omitempty"`
	DataType    ServiceCustomFieldDataType `json:"data_type"`
	FieldType   ServiceCustomFieldType     `json:"field_type"`
	Type        string                     `json:"type"`
	Value       interface{}                `json:"value"`
}

// ListServiceCustomFieldValuesResponse is the data structure returned from calling the GetServiceCustomFieldValues API endpoint
type ListServiceCustomFieldValuesResponse struct {
	CustomFields []ServiceCustomFieldValue `json:"custom_fields"`
}

// GetServiceCustomFieldValues gets custom field values for a service
func (c *Client) GetServiceCustomFieldValues(ctx context.Context, serviceID string) (*ListServiceCustomFieldValuesResponse, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	resp, err := c.get(ctx, "/services/"+serviceID+"/custom_fields/values", headers)
	if err != nil {
		return nil, err
	}

	var result ListServiceCustomFieldValuesResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result, nil
}

// UpdateServiceCustomFieldValues updates custom field values for a service
func (c *Client) UpdateServiceCustomFieldValues(ctx context.Context, serviceID string, customFields *ListServiceCustomFieldValuesResponse) (*ListServiceCustomFieldValuesResponse, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	resp, err := c.put(ctx, "/services/"+serviceID+"/custom_fields/values", customFields, headers)
	if err != nil {
		return nil, err
	}

	var result ListServiceCustomFieldValuesResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result, nil
}

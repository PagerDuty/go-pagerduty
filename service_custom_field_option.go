package pagerduty

import (
	"context"
	"fmt"
)

// ListServiceCustomFieldOptionsResponse is the data structure returned from calling the ListServiceCustomFieldOptions API endpoint
type ListServiceCustomFieldOptionsResponse struct {
	FieldOptions []ServiceCustomFieldOption `json:"field_options"`
}

// GetServiceCustomFieldOptionResponse is the data structure returned from calling the GetServiceCustomFieldOption API endpoint
type GetServiceCustomFieldOptionResponse struct {
	FieldOption ServiceCustomFieldOption `json:"field_option"`
}

// ListServiceCustomFieldOptions lists all options for a given field
func (c *Client) ListServiceCustomFieldOptions(ctx context.Context, fieldID string) (*ListServiceCustomFieldOptionsResponse, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	resp, err := c.get(ctx, "/services/custom_fields/"+fieldID+"/field_options", headers)
	if err != nil {
		return nil, err
	}

	var result ListServiceCustomFieldOptionsResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result, nil
}

// GetServiceCustomFieldOption gets a specific field option for a given field
func (c *Client) GetServiceCustomFieldOption(ctx context.Context, fieldID string, optionID string) (*ServiceCustomFieldOption, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	resp, err := c.get(ctx, "/services/custom_fields/"+fieldID+"/field_options/"+optionID, headers)
	if err != nil {
		return nil, err
	}

	var result GetServiceCustomFieldOptionResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result.FieldOption, nil
}

// CreateServiceCustomFieldOption creates a new field option for a given field
func (c *Client) CreateServiceCustomFieldOption(ctx context.Context, fieldID string, option *ServiceCustomFieldOption) (*ServiceCustomFieldOption, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	d := map[string]interface{}{
		"field_option": option,
	}

	resp, err := c.post(ctx, "/services/custom_fields/"+fieldID+"/field_options", d, headers)
	if err != nil {
		return nil, err
	}

	var result GetServiceCustomFieldOptionResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result.FieldOption, nil
}

// UpdateServiceCustomFieldOption updates an existing field option for a given field
func (c *Client) UpdateServiceCustomFieldOption(ctx context.Context, fieldID string, option *ServiceCustomFieldOption) (*ServiceCustomFieldOption, error) {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	d := map[string]interface{}{
		"field_option": option,
	}

	resp, err := c.put(ctx, "/services/custom_fields/"+fieldID+"/field_options/"+option.ID, d, headers)
	if err != nil {
		return nil, err
	}

	var result GetServiceCustomFieldOptionResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}

	return &result.FieldOption, nil
}

// DeleteServiceCustomFieldOption deletes a field option for a given field
func (c *Client) DeleteServiceCustomFieldOption(ctx context.Context, fieldID string, optionID string) error {
	headers := map[string]string{
		"X-EARLY-ACCESS": "service-custom-fields-preview",
	}

	_, err := c.deleteWithHeaders(ctx, "/services/custom_fields/"+fieldID+"/field_options/"+optionID, headers)
	return err
}

package pagerduty

import (
	"context"

	"github.com/google/go-querystring/query"
)

// IncidentType is allows to categorize incidents, such as a security incident, a major incident, or a fraud incident.
type IncidentType struct {
	Enabled     bool          `json:"enabled,omitempty"`
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Parent      *APIReference `json:"parent,omitempty"`
	Type        string        `json:"type,omitempty"`
	Description string        `json:"description,omitempty"`
	UpdatedAt   string        `json:"updated_at,omitempty"`
	DisplayName string        `json:"display_name,omitempty"`
}

type incidentTypeResponse struct {
	IncidentType IncidentType `json:"incident_type"`
}

// ListIncidentsTypesOptions is the structure used when passing parameters to the ListIncidentTypes API endpoint.
type ListIncidentTypesOptions struct {
	Filter string `url:"filter,omitempty"` // enabled disabled all
}

// ListIncidentsTypesResponse is the response structure when calling the ListIncidentTypes API endpoint.
type ListIncidentTypesResponse struct {
	IncidentTypes []IncidentType `json:"incident_types"`
}

// ListIncidentTypes list the available incident types.
func (c *Client) ListIncidentTypes(ctx context.Context, o ListIncidentTypesOptions) (*ListIncidentTypesResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents/types?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result ListIncidentTypesResponse
	err = c.decodeJSON(resp, &result)

	return &result, err
}

// CreateIncidentTypeOptions contains the parameters for creating a new incident type.
type CreateIncidentTypeOptions struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	ParentType  string  `json:"parent_type"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CreateIncidentType creates a new incident type.
func (c *Client) CreateIncidentType(ctx context.Context, o CreateIncidentTypeOptions) (*IncidentType, error) {
	d := map[string]CreateIncidentTypeOptions{
		"incident_type": o,
	}

	resp, err := c.post(ctx, "/incidents/types", d, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeResponse
	err = c.decodeJSON(resp, &result)

	return &result.IncidentType, err
}

// GetIncidentTypeOptions contains the parameters for retrieving a specific incident type.
type GetIncidentTypeOptions struct{}

// GetIncidentType retrieves a specific incident type by ID or name.
func (c *Client) GetIncidentType(ctx context.Context, idOrName string, o GetIncidentTypeOptions) (*IncidentType, error) {
	resp, err := c.get(ctx, "/incidents/types/"+idOrName, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeResponse
	err = c.decodeJSON(resp, &result)

	return &result.IncidentType, err
}

// UpdateIncidentTypeOptions contains the parameters for updating an incident type.
type UpdateIncidentTypeOptions struct {
	DisplayName *string `json:"display_name,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateIncidentType updates an existing incident type with the provided options.
func (c *Client) UpdateIncidentType(ctx context.Context, idOrName string, o UpdateIncidentTypeOptions) (*IncidentType, error) {
	d := map[string]UpdateIncidentTypeOptions{
		"incident_type": o,
	}

	resp, err := c.put(ctx, "/incidents/types/"+idOrName, d, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeResponse
	err = c.decodeJSON(resp, &result)

	return &result.IncidentType, err
}

// IncidentTypeField represents a custom field configuration for an incident type.
type IncidentTypeField struct {
	Enabled      bool                      `json:"enabled,omitempty"`
	ID           string                    `json:"id,omitempty"`
	Name         string                    `json:"name,omitempty"`
	Type         string                    `json:"type,omitempty"`
	Self         string                    `json:"self,omitempty"`
	Description  string                    `json:"description,omitempty"`
	FieldType    string                    `json:"field_type,omitempty"` // single_value single_value_fixed multi_value multi_value_fixed
	DataType     string                    `json:"data_type,omitempty"`  // boolean integer float string datetime url
	CreatedAt    string                    `json:"created_at,omitempty"`
	UpdatedAt    string                    `json:"updated_at,omitempty"`
	DisplayName  string                    `json:"display_name,omitempty"`
	DefaultValue interface{}               `json:"default_value,omitempty"`
	IncidentType string                    `json:"incident_type,omitempty"`
	Summary      string                    `json:"summary,omitempty"`
	FieldOptions []IncidentTypeFieldOption `json:"field_options,omitempty"`
}

// IncidentTypeFieldOption represents an option for a custom field.
type IncidentTypeFieldOption struct {
	ID        string                       `json:"id,omitempty"`
	Type      string                       `json:"type,omitempty"`
	CreatedAt string                       `json:"created_at,omitempty"`
	UpdatedAt string                       `json:"updated_at,omitempty"`
	Data      *IncidentTypeFieldOptionData `json:"data,omitempty"`
}

// IncidentTypeFieldOptionData represents the data for a field option.
type IncidentTypeFieldOptionData struct {
	Value    string `json:"value,omitempty"`
	DataType string `json:"data_type,omitempty"`
}

// ListIncidentTypeFieldsOptions contains the parameters for listing incident type fields.
type ListIncidentTypeFieldsOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// ListIncidentTypeFieldsResponse represents the response from listing incident type fields.
type ListIncidentTypeFieldsResponse struct {
	Fields []IncidentTypeField `json:"fields,omitempty"`
}

// ListIncidentTypeFields retrieves all custom fields for a specific incident type.
func (c *Client) ListIncidentTypeFields(ctx context.Context, typeIDOrName string, o ListIncidentTypeFieldsOptions) (*ListIncidentTypeFieldsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result ListIncidentTypeFieldsResponse
	err = c.decodeJSON(resp, &result)

	return &result, err
}

type incidentTypeFieldsResponse struct {
	Field IncidentTypeField `json:"field"`
}

// CreateIncidentTypeFieldOptions contains the parameters for creating a new incident type field.
type CreateIncidentTypeFieldOptions struct {
	Name         string                    `json:"name"`
	DisplayName  string                    `json:"display_name"`
	DataType     string                    `json:"data_type"`  // boolean integer float string datetime url
	FieldType    string                    `json:"field_type"` // single_value single_value_fixed multi_value multi_value_fixed
	DefaultValue interface{}               `json:"default_value,omitempty"`
	Description  *string                   `json:"description,omitempty"`
	Enabled      *bool                     `json:"enabled,omitempty"`
	FieldOptions []IncidentTypeFieldOption `json:"field_options,omitempty"`
}

// CreateIncidentTypeField creates a new custom field for a specific incident type.
func (c *Client) CreateIncidentTypeField(ctx context.Context, typeIDOrName string, o CreateIncidentTypeFieldOptions) (*IncidentTypeField, error) {
	d := map[string]CreateIncidentTypeFieldOptions{
		"field": o,
	}

	resp, err := c.post(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields", d, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeFieldsResponse
	err = c.decodeJSON(resp, &result)

	return &result.Field, err
}

// GetIncidentTypeFieldOptions contains the parameters for retrieving a specific incident type field.
type GetIncidentTypeFieldOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetIncidentTypeField retrieves a specific custom field for an incident type.
func (c *Client) GetIncidentTypeField(ctx context.Context, typeIDOrName string, fieldID string, o GetIncidentTypeFieldOptions) (*IncidentTypeField, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID+"?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeFieldsResponse
	err = c.decodeJSON(resp, &result)

	return &result.Field, err
}

// UpdateIncidentTypeFieldOptions contains the parameters for updating an incident type field.
type UpdateIncidentTypeFieldOptions struct {
	DisplayName  *string      `json:"display_name,omitempty"`
	Enabled      *bool        `json:"enabled,omitempty"`
	DefaultValue *interface{} `json:"default_value,omitempty"`
	Description  *string      `json:"description,omitempty"`
}

// UpdateIncidentTypeField updates an existing custom field for an incident type.
func (c *Client) UpdateIncidentTypeField(ctx context.Context, typeIDOrName, fieldID string, o UpdateIncidentTypeFieldOptions) (*IncidentTypeField, error) {
	d := map[string]UpdateIncidentTypeFieldOptions{
		"field": o,
	}

	resp, err := c.put(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID, d, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeFieldsResponse
	err = c.decodeJSON(resp, &result)

	return &result.Field, err
}

// DeleteIncidentTypeField removes a custom field from an incident type.
func (c *Client) DeleteIncidentTypeField(ctx context.Context, typeIDOrName, fieldID string) error {
	_, err := c.delete(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID)
	return err
}

// ListIncidentTypeFieldOptionsOptions contains the parameters for listing field options.
type ListIncidentTypeFieldOptionsOptions struct{}

// ListIncidentTypeFieldOptionsResponse represents the response from listing field options.
type ListIncidentTypeFieldOptionsResponse struct {
	FieldOptions []IncidentTypeFieldOption `json:"field_options,omitempty"`
}

// ListIncidentTypeFieldOptions retrieves all options for a specific custom field.
func (c *Client) ListIncidentTypeFieldOptions(ctx context.Context, typeIDOrName, fieldID string, o ListIncidentTypeFieldOptionsOptions) (*ListIncidentTypeFieldOptionsResponse, error) {
	resp, err := c.get(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID+"/field_options", nil)
	if err != nil {
		return nil, err
	}

	var result ListIncidentTypeFieldOptionsResponse
	err = c.decodeJSON(resp, &result)

	return &result, err
}

type incidentTypeFieldOptionsResponse struct {
	FieldOption IncidentTypeFieldOption `json:"field_option,omitempty"`
}

// CreateIncidentTypeFieldOptionPayload contains the parameters for creating a new field option.
type CreateIncidentTypeFieldOptionPayload struct {
	Data *IncidentTypeFieldOptionData `json:"data,omitempty"`
}

// CreateIncidentTypeFieldOption creates a new option for a custom field.
func (c *Client) CreateIncidentTypeFieldOption(ctx context.Context, typeIDOrName, fieldID string, o CreateIncidentTypeFieldOptionPayload) (*IncidentTypeFieldOption, error) {
	d := map[string]CreateIncidentTypeFieldOptionPayload{
		"field_option": o,
	}

	resp, err := c.post(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID+"/field_options", d, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeFieldOptionsResponse
	err = c.decodeJSON(resp, &result)

	return &result.FieldOption, err
}

// GetIncidentTypeFieldOptionOptions contains the parameters for retrieving a specific field option.
type GetIncidentTypeFieldOptionOptions struct{}

// GetIncidentTypeFieldOption retrieves a specific option for a custom field.
func (c *Client) GetIncidentTypeFieldOption(ctx context.Context, typeIDOrName, fieldID, fieldOptionID string, o GetIncidentTypeFieldOptionOptions) (*IncidentTypeFieldOption, error) {
	resp, err := c.get(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID+"/field_options/"+fieldOptionID, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeFieldOptionsResponse
	err = c.decodeJSON(resp, &result)

	return &result.FieldOption, err
}

// UpdateIncidentTypeFieldOptionPayload contains the parameters for updating a field option.
type UpdateIncidentTypeFieldOptionPayload struct {
	ID   string                       `json:"id,omitempty"`
	Data *IncidentTypeFieldOptionData `json:"data,omitempty"`
}

// UpdateIncidentTypeFieldOption updates an existing option for a custom field.
func (c *Client) UpdateIncidentTypeFieldOption(ctx context.Context, typeIDOrName, fieldID string, o UpdateIncidentTypeFieldOptionPayload) (*IncidentTypeFieldOption, error) {
	d := map[string]UpdateIncidentTypeFieldOptionPayload{
		"field_option": o,
	}

	resp, err := c.put(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID+"/field_options/"+o.ID, d, nil)
	if err != nil {
		return nil, err
	}

	var result incidentTypeFieldOptionsResponse
	err = c.decodeJSON(resp, &result)

	return &result.FieldOption, err
}

// DeleteIncidentTypeFieldOption removes an option from a custom field.
func (c *Client) DeleteIncidentTypeFieldOption(ctx context.Context, typeIDOrName, fieldID, fieldOptionID string) error {
	_, err := c.delete(ctx, "/incidents/types/"+typeIDOrName+"/custom_fields/"+fieldID+"/field_options/"+fieldOptionID)
	return err
}

package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Integration is an endpoint (like Nagios, email, or an API call) that generates events, which are normalized and de-duplicated by PagerDuty to create incidents.
type Integration struct {
	APIObject
	Name             string     `json:"name,omitempty"`
	Service          *APIObject `json:"service,omitempty"`
	CreatedAt        string     `json:"created_at,omitempty"`
	Vendor           *APIObject `json:"vendor,omitempty"`
	IntegrationKey   string     `json:"integration_key,omitempty"`
	IntegrationEmail string     `json:"integration_email,omitempty"`
}

// InlineModel represents when a scheduled action will occur.
type InlineModel struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// ScheduledAction contains scheduled actions for the service.
type ScheduledAction struct {
	Type      string      `json:"type,omitempty"`
	At        InlineModel `json:"at,omitempty"`
	ToUrgency string      `json:"to_urgency"`
}

// IncidentUrgencyType are the incidents urgency during or outside support hours.
type IncidentUrgencyType struct {
	Type    string `json:"type,omitempty"`
	Urgency string `json:"urgency,omitempty"`
}

// SupportHours are the support hours for the service.
type SupportHours struct {
	Type       string `json:"type,omitempty"`
	Timezone   string `json:"time_zone,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	DaysOfWeek []uint `json:"days_of_week,omitempty"`
}

// IncidentUrgencyRule is the default urgency for new incidents.
type IncidentUrgencyRule struct {
	Type                string               `json:"type,omitempty"`
	Urgency             string               `json:"urgency,omitempty"`
	DuringSupportHours  *IncidentUrgencyType `json:"during_support_hours,omitempty"`
	OutsideSupportHours *IncidentUrgencyType `json:"outside_support_hours,omitempty"`
}

// ListServiceRulesResponse represents a list of rules in a service
type ListServiceRulesResponse struct {
	Offset uint          `json:"offset,omitempty"`
	Limit  uint          `json:"limit,omitempty"`
	More   bool          `json:"more,omitempty"`
	Total  uint          `json:"total,omitempty"`
	Rules  []ServiceRule `json:"rules,omitempty"`
}

// ServiceRule represents a Service rule
type ServiceRule struct {
	ID         string              `json:"id,omitempty"`
	Self       string              `json:"self,omitempty"`
	Disabled   *bool               `json:"disabled,omitempty"`
	Conditions *RuleConditions     `json:"conditions,omitempty"`
	TimeFrame  *RuleTimeFrame      `json:"time_frame,omitempty"`
	Position   *int                `json:"position,omitempty"`
	Actions    *ServiceRuleActions `json:"actions,omitempty"`
}

// ServiceRuleActions represents a rule action
type ServiceRuleActions struct {
	Annotate    *RuleActionParameter   `json:"annotate,omitempty"`
	EventAction *RuleActionParameter   `json:"event_action,omitempty"`
	Extractions []RuleActionExtraction `json:"extractions,omitempty"`
	Priority    *RuleActionParameter   `json:"priority,omitempty"`
	Severity    *RuleActionParameter   `json:"severity,omitempty"`
	Suppress    *RuleActionSuppress    `json:"suppress,omitempty"`
	Suspend     *RuleActionSuspend     `json:"suspend,omitempty"`
}

// Service represents something you monitor (like a web service, email service, or database service).
type Service struct {
	APIObject
	Name                    string                   `json:"name,omitempty"`
	Description             string                   `json:"description,omitempty"`
	AutoResolveTimeout      *uint                    `json:"auto_resolve_timeout"`
	AcknowledgementTimeout  *uint                    `json:"acknowledgement_timeout"`
	CreateAt                string                   `json:"created_at,omitempty"`
	Status                  string                   `json:"status,omitempty"`
	LastIncidentTimestamp   string                   `json:"last_incident_timestamp,omitempty"`
	Integrations            []Integration            `json:"integrations,omitempty"`
	EscalationPolicy        EscalationPolicy         `json:"escalation_policy,omitempty"`
	Teams                   []Team                   `json:"teams,omitempty"`
	IncidentUrgencyRule     *IncidentUrgencyRule     `json:"incident_urgency_rule,omitempty"`
	SupportHours            *SupportHours            `json:"support_hours"`
	ScheduledActions        []ScheduledAction        `json:"scheduled_actions"`
	AlertCreation           string                   `json:"alert_creation,omitempty"`
	AlertGrouping           string                   `json:"alert_grouping,omitempty"`
	AlertGroupingTimeout    *uint                    `json:"alert_grouping_timeout,omitempty"`
	AlertGroupingParameters *AlertGroupingParameters `json:"alert_grouping_parameters,omitempty"`
}

// AlertGroupingParameters defines how alerts on the servicewill be automatically grouped into incidents
type AlertGroupingParameters struct {
	Type   string                 `json:"type"`
	Config AlertGroupParamsConfig `json:"config"`
}

// AlertGroupParamsConfig is the config object on alert_grouping_parameters
type AlertGroupParamsConfig struct {
	Timeout   uint     `json:"timeout,omitempty"`
	Aggregate string   `json:"aggregate,omitempty"`
	Fields    []string `json:"fields,omitempty"`
}

// ListServiceOptions is the data structure used when calling the ListServices API endpoint.
type ListServiceOptions struct {
	APIListObject
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	TimeZone string   `url:"time_zone,omitempty"`
	SortBy   string   `url:"sort_by,omitempty"`
	Query    string   `url:"query,omitempty"`
	Includes []string `url:"include,omitempty,brackets"`
}

// ListServiceResponse is the data structure returned from calling the ListServices API endpoint.
type ListServiceResponse struct {
	APIListObject
	Services []Service
}

// ListServices lists existing services.
//
// Deprecated: Use ListServicesWithContext instead.
func (c *Client) ListServices(o ListServiceOptions) (*ListServiceResponse, error) {
	return c.ListServicesWithContext(context.Background(), o)
}

// ListServicesWithContext lists existing services.
func (c *Client) ListServicesWithContext(ctx context.Context, o ListServiceOptions) (*ListServiceResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/services?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListServiceResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListServicesPaginated lists existing services processing paginated responses
func (c *Client) ListServicesPaginated(ctx context.Context, o ListServiceOptions) ([]Service, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	var services []Service

	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListServiceResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		services = append(services, result.Services...)

		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	if err := c.pagedGet(ctx, "/services?"+v.Encode(), responseHandler); err != nil {
		return nil, err
	}

	return services, nil
}

// GetServiceOptions is the data structure used when calling the GetService API endpoint.
type GetServiceOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// GetService gets details about an existing service.
//
// Deprecated: Use GetServiceWithContext instead.
func (c *Client) GetService(id string, o *GetServiceOptions) (*Service, error) {
	return c.GetServiceWithContext(context.Background(), id, o)
}

// GetServiceWithContext gets details about an existing service.
func (c *Client) GetServiceWithContext(ctx context.Context, id string, o *GetServiceOptions) (*Service, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/services/"+id+"?"+v.Encode())
	return getServiceFromResponse(c, resp, err)
}

// CreateService creates a new service.
//
// Deprecated: Use CreateServiceWithContext instead.
func (c *Client) CreateService(s Service) (*Service, error) {
	return c.CreateServiceWithContext(context.Background(), s)
}

// CreateServiceWithContext creates a new service.
func (c *Client) CreateServiceWithContext(ctx context.Context, s Service) (*Service, error) {
	d := map[string]Service{
		"service": s,
	}

	resp, err := c.post(ctx, "/services", d, nil)
	return getServiceFromResponse(c, resp, err)
}

// UpdateService updates an existing service.
//
// Deprecated: Use UpdateServiceWithContext instead.
func (c *Client) UpdateService(s Service) (*Service, error) {
	return c.UpdateServiceWithContext(context.Background(), s)
}

// UpdateServiceWithContext updates an existing service.
func (c *Client) UpdateServiceWithContext(ctx context.Context, s Service) (*Service, error) {
	d := map[string]Service{
		"service": s,
	}

	resp, err := c.put(ctx, "/services/"+s.ID, d, nil)
	return getServiceFromResponse(c, resp, err)
}

// DeleteService deletes an existing service.
//
// Deprecated: Use DeleteServiceWithContext instead.
func (c *Client) DeleteService(id string) error {
	return c.DeleteServiceWithContext(context.Background(), id)
}

// DeleteServiceWithContext deletes an existing service.
func (c *Client) DeleteServiceWithContext(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/services/"+id)
	return err
}

// CreateIntegration creates a new integration belonging to a service.
//
// Deprecated: Use CreateIntegrationWithContext instead.
func (c *Client) CreateIntegration(id string, i Integration) (*Integration, error) {
	return c.CreateIntegrationWithContext(context.Background(), id, i)
}

// CreateIntegrationWithContext creates a new integration belonging to a service.
func (c *Client) CreateIntegrationWithContext(ctx context.Context, id string, i Integration) (*Integration, error) {
	d := map[string]Integration{
		"integration": i,
	}

	resp, err := c.post(ctx, "/services/"+id+"/integrations", d, nil)
	return getIntegrationFromResponse(c, resp, err)
}

// GetIntegrationOptions is the data structure used when calling the GetIntegration API endpoint.
type GetIntegrationOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetIntegration gets details about an integration belonging to a service.
//
// Deprecated: Use GetIntegrationWithContext instead.
func (c *Client) GetIntegration(serviceID, integrationID string, o GetIntegrationOptions) (*Integration, error) {
	return c.GetIntegrationWithContext(context.Background(), serviceID, integrationID, o)
}

// GetIntegrationWithContext gets details about an integration belonging to a service.
func (c *Client) GetIntegrationWithContext(ctx context.Context, serviceID, integrationID string, o GetIntegrationOptions) (*Integration, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/services/"+serviceID+"/integrations/"+integrationID+"?"+v.Encode())
	return getIntegrationFromResponse(c, resp, err)
}

// UpdateIntegration updates an integration belonging to a service.
//
// Deprecated: Use UpdateIntegrationWithContext instead.
func (c *Client) UpdateIntegration(serviceID string, i Integration) (*Integration, error) {
	return c.UpdateIntegrationWithContext(context.Background(), serviceID, i)
}

// UpdateIntegrationWithContext updates an integration belonging to a service.
func (c *Client) UpdateIntegrationWithContext(ctx context.Context, serviceID string, i Integration) (*Integration, error) {
	resp, err := c.put(ctx, "/services/"+serviceID+"/integrations/"+i.ID, i, nil)
	return getIntegrationFromResponse(c, resp, err)
}

// DeleteIntegration deletes an existing integration.
//
// Deprecated: Use DeleteIntegrationWithContext instead.
func (c *Client) DeleteIntegration(serviceID string, integrationID string) error {
	return c.DeleteIntegrationWithContext(context.Background(), serviceID, integrationID)
}

// DeleteIntegrationWithContext deletes an existing integration.
func (c *Client) DeleteIntegrationWithContext(ctx context.Context, serviceID string, integrationID string) error {
	_, err := c.delete(ctx, "/services/"+serviceID+"/integrations/"+integrationID)
	return err
}

// ListServiceRulesPaginated gets all rules for a service.
func (c *Client) ListServiceRulesPaginated(ctx context.Context, serviceID string) ([]ServiceRule, error) {
	var rules []ServiceRule

	// Create a handler closure capable of parsing data from the Service rules endpoint
	// and appending resultant Service rules to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListServiceRulesResponse

		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		rules = append(rules, result.Rules...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet(ctx, "/services/"+serviceID+"/rules", responseHandler); err != nil {
		return nil, err
	}

	return rules, nil
}

// GetServiceRule gets a service rule.
func (c *Client) GetServiceRule(ctx context.Context, serviceID, ruleID string) (ServiceRule, error) {
	resp, err := c.get(ctx, "/services/"+serviceID+"/rules/"+ruleID)
	return getServiceRuleFromResponse(c, resp, err)
}

// DeleteServiceRule deletes a service rule.
func (c *Client) DeleteServiceRule(ctx context.Context, serviceID, ruleID string) error {
	_, err := c.delete(ctx, "/services/"+serviceID+"/rules/"+ruleID)
	return err
}

// CreateServiceRule creates a service rule.
func (c *Client) CreateServiceRule(ctx context.Context, serviceID string, rule ServiceRule) (ServiceRule, error) {
	d := map[string]ServiceRule{
		"rule": rule,
	}
	resp, err := c.post(ctx, "/services/"+serviceID+"/rules/", d, nil)
	return getServiceRuleFromResponse(c, resp, err)
}

// UpdateServiceRule updates a service rule.
func (c *Client) UpdateServiceRule(ctx context.Context, serviceID, ruleID string, rule ServiceRule) (ServiceRule, error) {
	d := map[string]ServiceRule{
		"rule": rule,
	}
	resp, err := c.put(ctx, "/services/"+serviceID+"/rules/"+ruleID, d, nil)
	return getServiceRuleFromResponse(c, resp, err)
}

func getServiceRuleFromResponse(c *Client, resp *http.Response, err error) (ServiceRule, error) {
	if err != nil {
		return ServiceRule{}, err
	}

	var target map[string]ServiceRule
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return ServiceRule{}, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "rule"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return ServiceRule{}, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return t, nil
}

func getServiceFromResponse(c *Client, resp *http.Response, err error) (*Service, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]Service
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "service"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

func getIntegrationFromResponse(c *Client, resp *http.Response, err error) (*Integration, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]Integration
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", err)
	}

	const rootNode = "integration"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

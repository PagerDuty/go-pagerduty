package pagerduty

import (
	"context"

	"github.com/google/go-querystring/query"
)

// JiraCloudAccountsMapping establishes a connection between a PagerDuty account
// and a Jira Cloud instance, enabling integration and synchronization between
// the two platforms
type JiraCloudAccountsMapping struct {
	UpdatedAt        string           `json:"updated_at,omitempty"`
	CreatedAt        string           `json:"created_at,omitempty"`
	ID               string           `json:"id,omitempty"`
	JiraCloudAccount JiraCloudAccount `json:"jira_cloud_account"`
	PagerDutyAccount PagerDutyAccount `json:"pagerduty_account"`
}

// JiraCloudAccount describes an account from Jira
type JiraCloudAccount struct {
	// The base URL of the Jira Cloud instance, used for API calls and
	// constructing links
	BaseURL string `json:"base_url"`
}

// PagerDutyAccount describes an account from PagerDuty
type PagerDutyAccount struct {
	// The unique subdomain of the PagerDuty account, used to identify and
	// access the account (e.g., acme in https://acme.pagerduty.com)
	Subdomain string `json:"subdomain"`
}

// JiraCloudAccountsMappingRule configures the bidirectional synchronization
// between Jira issues and PagerDuty incidents
type JiraCloudAccountsMappingRule struct {
	AccountsMapping             *APIObject                         `json:"account_mapping,omitempty"`
	AutocreateJqlDisabledReason string                             `json:"autocreate_jql_disabled_reason,omitempty"`
	AutocreateJqlDisabledUntil  string                             `json:"autocreate_jql_disabled_until,omitempty"`
	Config                      JiraCloudAccountsMappingRuleConfig `json:"config"`
	ID                          string                             `json:"id,omitempty"`
	Name                        string                             `json:"name"`
	CreatedAt                   string                             `json:"created_at,omitempty"`
	UpdatedAt                   string                             `json:"updated_at,omitempty"`
}

// JiraCloudAccountsMappingRuleConfig is the configuration for bidirectional
// synchronization between Jira issues and PagerDuty incidents
type JiraCloudAccountsMappingRuleConfig struct {
	Jira    JiraCloudSettings `json:"jira"`
	Service APIObject         `json:"service"`
}

// JiraCloudSettings are settings for the Jira aspect of the synchronization
type JiraCloudSettings struct {
	AutocreateJQL                *string                `json:"autocreate_jql"`
	CreateIssueOnIncidentTrigger bool                   `json:"create_issue_on_incident_trigger"`
	CustomFields                 []JiraCloudCustomField `json:"custom_fields"`
	IssueType                    JiraCloudReference     `json:"issue_type"`
	Priorities                   []JiraCloudPriority    `json:"priorities"`
	Project                      JiraCloudReference     `json:"project"`
	StatusMapping                JiraCloudStatusMapping `json:"status_mapping"`
	SyncNotesUser                *UserJiraCloud         `json:"sync_notes_user"`
}

// JiraCloudCustomField defines how Jira fields are populated when a Jira Issue
// is created from a PagerDuty Incident
type JiraCloudCustomField struct {
	// The PagerDuty incident field from which the value will be extracted
	// (only applicable if type is attribute)
	//
	// Allowed values:
	// incident_number, incident_title, incident_description,
	// incident_status, incident_created_at, incident_service,
	// incident_escalation_policy, incident_impacted_services,
	// incident_html_url, incident_assignees, incident_acknowledgers,
	// incident_last_status_change_at, incident_last_status_change_by,
	// incident_urgency, incident_priority, null
	SourceIncidentField *string `json:"source_incident_field"`

	// The unique identifier key of the Jira field that will be set
	TargetIssueField string `json:"target_issue_field"`

	// The human-readable name of the Jira field
	TargetIssueFieldName string `json:"target_issue_field_name"`

	// The type of the value that will be set
	//
	// Allowed values:
	// attribute, const, jira_value
	Type string `json:"type"`

	// The value to be set for the Jira field (only applicable if type is
	// const or jira_value)
	Value interface{} `json:"value,omitempty"`
}

// JiraCloudReference is a reference pointing to a Jira Cloud object
type JiraCloudReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key,omitempty"`
}

// JiraCloudPriority is an association between a PagerDuty incident priority
// and a Jira issue priority
type JiraCloudPriority struct {
	JiraID      string `json:"jira_id"`
	PagerDutyID string `json:"pagerduty_id"`
}

// JiraCloudStatusMapping is an association between PagerDuty incident statuses
// to their corresponding Jira issue statuses
type JiraCloudStatusMapping struct {
	Acknowledged *JiraCloudReference `json:"acknowledged"`
	Resolved     *JiraCloudReference `json:"resolved"`
	Triggered    *JiraCloudReference `json:"triggered"`
}

// UserJiraCloud is a PagerDuty user for syncing notes and comments between Jira
// issues and PagerDuty incidents. If not provided, note synchronization is
// disabled
type UserJiraCloud struct {
	APIObject
	Email string `json:"email,omitempty"`
}

// ListJiraCloudAccountsMappingsOptions are the options available when calling the ListJiraCloudAccountsMappings API endpoint
type ListJiraCloudAccountsMappingsOptions struct {
	Limit  uint `url:"limit,omitempty"`
	Offset uint `url:"offset,omitempty"`
	Total  bool `url:"total,omitempty"`
}

// ListJiraCloudAccountsMappingsResponse is the response when calling the ListJiraCloudAccountsMappings API endpoint
type ListJiraCloudAccountsMappingsResponse struct {
	APIListObject
	AccountsMappings []JiraCloudAccountsMapping `json:"accounts_mappings"`
}

// ListJiraCloudAccountsMappings lists existing account mappings
func (c *Client) ListJiraCloudAccountsMappings(ctx context.Context, o ListJiraCloudAccountsMappingsOptions) (*ListJiraCloudAccountsMappingsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/integration-jira-cloud/accounts_mappings?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result ListJiraCloudAccountsMappingsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetJiraCloudAccountsMapping lists existing account mappings
func (c *Client) GetJiraCloudAccountsMapping(ctx context.Context, id string) (*JiraCloudAccountsMapping, error) {
	resp, err := c.get(ctx, "/integration-jira-cloud/accounts_mappings/"+id, nil)
	if err != nil {
		return nil, err
	}

	var result JiraCloudAccountsMapping
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListJiraCloudAccountsMappingRulesOptions are the options available when
// calling the ListJiraCloudAccountsMappingRules API endpoint
type ListJiraCloudAccountsMappingRulesOptions struct {
	Limit  uint `url:"limit,omitempty"`
	Offset uint `url:"offset,omitempty"`
	Total  bool `url:"total,omitempty"`
}

// ListJiraCloudAccountsMappingRulesResponse is the response when calling the
// ListJiraCloudAccountsMappingRules API endpoint
type ListJiraCloudAccountsMappingRulesResponse struct {
	APIListObject
	Rules []JiraCloudAccountsMappingRule `json:"rules"`
}

// ListJiraCloudAccountsMappingRules lists existing rules for a specific account
// mapping
func (c *Client) ListJiraCloudAccountsMappingRules(ctx context.Context, id string, o ListJiraCloudAccountsMappingRulesOptions) (*ListJiraCloudAccountsMappingsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/integration-jira-cloud/accounts_mappings/"+id+"/rules?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result ListJiraCloudAccountsMappingsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateJiraCloudAccountsMappingRule creates a new rule in Jira Cloud's integration
func (c *Client) CreateJiraCloudAccountsMappingRule(ctx context.Context, id string, rule JiraCloudAccountsMappingRule) (*JiraCloudAccountsMappingRule, error) {
	resp, err := c.post(ctx, "/integration-jira-cloud/accounts_mappings/"+id+"/rules", rule, nil)
	if err != nil {
		return nil, err
	}

	var result JiraCloudAccountsMappingRule
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetJiraCloudAccountsMappingRule gets detailed information about an existing rule
func (c *Client) GetJiraCloudAccountsMappingRule(ctx context.Context, id, ruleID string) (*JiraCloudAccountsMappingRule, error) {
	resp, err := c.get(ctx, "/integration-jira-cloud/accounts_mappings/"+id+"/rules/"+ruleID, nil)
	if err != nil {
		return nil, err
	}

	var result JiraCloudAccountsMappingRule
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteJiraCloudAccountsMappingRule deletes an existing rule in Jira Cloud's integration
func (c *Client) DeleteJiraCloudAccountsMappingRule(ctx context.Context, id, ruleID string) error {
	_, err := c.delete(ctx, "/integration-jira-cloud/accounts_mappings/"+id+"/rules/"+ruleID)
	return err
}

// updateJiraCloudAccountsMappingRuleBody is the body for the call to the
// UpdateJiraCloudAccountsMappingRule API endpoint
type updateJiraCloudAccountsMappingRuleBody struct {
	Config updateJiraCloudAccountsMappingRuleConfig `json:"config"`
	Name   string                                   `json:"name"`
}

// updateJiraCloudAccountsMappingRuleOptionsConfig is an special representation
// of a configuration used for updating a rule.
type updateJiraCloudAccountsMappingRuleConfig struct {
	Jira updateJiraCloudSettings `json:"jira"`
}

// updateJiraCloudSettings are settings to update the Jira aspect of the
// synchronization
type updateJiraCloudSettings struct {
	AutocreateJQL                *string                `json:"autocreate_jql"`
	CreateIssueOnIncidentTrigger bool                   `json:"create_issue_on_incident_trigger"`
	CustomFields                 []JiraCloudCustomField `json:"custom_fields"`
	IssueType                    JiraCloudReference     `json:"issue_type"`
	Priorities                   []JiraCloudPriority    `json:"priorities"`
	StatusMapping                JiraCloudStatusMapping `json:"status_mapping"`
	SyncNotesUser                *UserJiraCloud         `json:"sync_notes_user"`
}

// UpdateJiraCloudAccountsMappingRule updates an existing rule in Jira Cloud's integration
func (c *Client) UpdateJiraCloudAccountsMappingRule(ctx context.Context, accountMappingID string, rule JiraCloudAccountsMappingRule) (*JiraCloudAccountsMappingRule, error) {
	o := updateJiraCloudAccountsMappingRuleBody{
		Config: updateJiraCloudAccountsMappingRuleConfig{
			Jira: updateJiraCloudSettings{
				AutocreateJQL:                rule.Config.Jira.AutocreateJQL,
				CreateIssueOnIncidentTrigger: rule.Config.Jira.CreateIssueOnIncidentTrigger,
				CustomFields:                 rule.Config.Jira.CustomFields,
				IssueType:                    rule.Config.Jira.IssueType,
				Priorities:                   rule.Config.Jira.Priorities,
				StatusMapping:                rule.Config.Jira.StatusMapping,
				SyncNotesUser:                rule.Config.Jira.SyncNotesUser,
			},
		},
		Name: rule.Name,
	}

	resp, err := c.put(ctx, "/integration-jira-cloud/accounts_mappings/"+accountMappingID+"/rules/"+rule.ID, o, nil)
	if err != nil {
		return nil, err
	}

	var result JiraCloudAccountsMappingRule
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

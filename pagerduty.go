package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiEndpoint = "https://api.pagerduty.com"
)

// APIObject represents generic api json response that is shared by most
// domain object (like escalation
type APIObject struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// APIListObject are the fields used to control pagination when listing objects.
type APIListObject struct {
	Limit  uint `url:"limit,omitempty"`
	Offset uint `url:"offset,omitempty"`
	More   bool `url:"more,omitempty"`
	Total  uint `url:"total,omitempty"`
}

// APIReference are the fields required to reference another API object.
type APIReference struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type errorObject struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PagerdutyClient struct {
	authToken string
}

// Pagerduty all methods that the pagerduty can execute
type Pagerduty interface {
	AuthToken() string
	Delete(path string) (*http.Response, error)
	Put(path string, payload interface{}, headers *map[string]string) (*http.Response, error)
	Post(path string, payload interface{}) (*http.Response, error)
	Get(path string) (*http.Response, error)
	Do(method, path string, body io.Reader, headers *map[string]string) (*http.Response, error)
	CheckResponse(resp *http.Response, err error) (*http.Response, error)
	GetErrorFromResponse(resp *http.Response) (*errorObject, error)
	ListAbilities() (*ListAbilityResponse, error)
	TestAbility(ability string) error
	ListAddons(o ListAddonOptions) (*ListAddonResponse, error)
	InstallAddon(a Addon) (*Addon, error)
	DeleteAddon(id string) error
	GetAddon(id string) (*Addon, error)
	UpdateAddon(id string, a Addon) (*Addon, error)
	ListEscalationPolicies(o ListEscalationPoliciesOptions) (*ListEscalationPoliciesResponse, error)
	CreateEscalationPolicy(e EscalationPolicy) (*EscalationPolicy, error)
	DeleteEscalationPolicy(id string) error
	GetEscalationPolicy(id string, o *GetEscalationPolicyOptions) (*EscalationPolicy, error)
	UpdateEscalationPolicy(id string, e *EscalationPolicy) (*EscalationPolicy, error)
	CreateEscalationRule(escID string, e EscalationRule) (*EscalationRule, error)
	GetEscalationRule(escID string, id string, o *GetEscalationRuleOptions) (*EscalationRule, error)
	DeleteEscalationRule(escID string, id string) error
	UpdateEscalationRule(escID string, id string, e *EscalationRule) (*EscalationRule, error)
	ListEscalationRules(escID string) (*ListEscalationRulesResponse, error)
	ListIncidents(o ListIncidentsOptions) (*ListIncidentsResponse, error)
	ManageIncidents(from string, incidents []Incident) error
	GetIncident(id string) (*Incident, error)
	ListIncidentNotes(id string) ([]IncidentNote, error)
	CreateIncidentNote(id string, note IncidentNote) error
	SnoozeIncident(id string, duration uint) error
	ListIncidentLogEntries(id string, o ListIncidentLogEntriesOptions) (*ListIncidentLogEntriesResponse, error)
	ListLogEntries(o ListLogEntriesOptions) (*ListLogEntryResponse, error)
	GetLogEntry(id string, o GetLogEntryOptions) (*LogEntry, error)
	ListMaintenanceWindows(o ListMaintenanceWindowsOptions) (*ListMaintenanceWindowsResponse, error)
	CreateMaintenanceWindows(m MaintenanceWindow) (*MaintenanceWindow, error)
	DeleteMaintenanceWindow(id string) error
	GetMaintenanceWindow(id string, o GetMaintenanceWindowOptions) (*MaintenanceWindow, error)
	UpdateMaintenanceWindow(m MaintenanceWindow) (*MaintenanceWindow, error)
	ListNotifications(o ListNotificationOptions) (*ListNotificationsResponse, error)
	ListOnCalls(o ListOnCallOptions) (*ListOnCallsResponse, error)
	ListSchedules(o ListSchedulesOptions) (*ListSchedulesResponse, error)
	CreateSchedule(s Schedule) (*Schedule, error)
	PreviewSchedule(s Schedule, o PreviewScheduleOptions) error
	DeleteSchedule(id string) error
	GetSchedule(id string, o GetScheduleOptions) (*Schedule, error)
	UpdateSchedule(id string, s Schedule) (*Schedule, error)
	ListOverrides(id string, o ListOverridesOptions) ([]Override, error)
	CreateOverride(id string, o Override) (*Override, error)
	DeleteOverride(scheduleID, overrideID string) error
	ListOnCallUsers(id string, o ListOnCallUsersOptions) ([]User, error)
	ListServices(o ListServiceOptions) (*ListServiceResponse, error)
	GetService(id string, o *GetServiceOptions) (*Service, error)
	CreateService(s Service) (*Service, error)
	UpdateService(s Service) (*Service, error)
	DeleteService(id string) error
	CreateIntegration(id string, i Integration) (*Integration, error)
	GetIntegration(serviceID, integrationID string, o GetIntegrationOptions) (*Integration, error)
	UpdateIntegration(serviceID string, i Integration) (*Integration, error)
	DeleteIntegration(serviceID string, integrationID string) error
	ListTeams(o ListTeamOptions) (*ListTeamResponse, error)
	CreateTeam(t *Team) (*Team, error)
	DeleteTeam(id string) error
	GetTeam(id string) (*Team, error)
	UpdateTeam(id string, t *Team) (*Team, error)
	RemoveEscalationPolicyFromTeam(teamID, epID string) error
	AddEscalationPolicyToTeam(teamID, epID string) error
	RemoveUserFromTeam(teamID, userID string) error
	AddUserToTeam(teamID, userID string) error
	ListUsers(o ListUsersOptions) (*ListUsersResponse, error)
	CreateUser(u User) (*User, error)
	DeleteUser(id string) error
	GetUser(id string, o GetUserOptions) (*User, error)
	UpdateUser(u User) (*User, error)
	ListVendors(o ListVendorOptions) (*ListVendorResponse, error)
	GetVendor(id string) (*Vendor, error)
}

// NewPagerduty creates a new client instance.
func NewPagerduty(authToken string) Pagerduty {
	pd := PagerdutyClient{
		authToken: authToken,
	}
	return &pd
}

func (pd *PagerdutyClient) AuthToken() string {
	return pd.authToken
}

func (pd *PagerdutyClient) Delete(path string) (*http.Response, error) {
	return pd.Do("DELETE", path, nil, nil)
}

func (pd *PagerdutyClient) Put(path string, payload interface{}, headers *map[string]string) (*http.Response, error) {

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		return pd.Do("PUT", path, bytes.NewBuffer(data), headers)
	}
	return pd.Do("PUT", path, nil, headers)
}

func (pd *PagerdutyClient) Post(path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return pd.Do("POST", path, bytes.NewBuffer(data), nil)
}

func (pd *PagerdutyClient) PostWithHeader(path string, payload interface{}, headers *map[string]string) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return pd.Do("POST", path, bytes.NewBuffer(data), headers)
}

func (pd *PagerdutyClient) Get(path string) (*http.Response, error) {
	return pd.Do("GET", path, nil, nil)
}

func (pd *PagerdutyClient) Do(method, path string, body io.Reader, headers *map[string]string) (*http.Response, error) {
	endpoint := apiEndpoint + path
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token token="+pd.authToken)

	resp, err := http.DefaultClient.Do(req)
	return pd.CheckResponse(resp, err)
}

func (pd *PagerdutyClient) CheckResponse(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, fmt.Errorf("Error calling the API endpoint: %v", err)
	}
	if 199 >= resp.StatusCode || 300 <= resp.StatusCode {
		var eo *errorObject
		var getErr error
		if eo, getErr = pd.GetErrorFromResponse(resp); getErr != nil {
			return resp, fmt.Errorf("Response did not contain formatted error: %s. HTTP response code: %v. Raw response: %+v", getErr, resp.StatusCode, resp)
		}
		return resp, fmt.Errorf("Failed call API endpoint. HTTP response code: %v. Error: %v", resp.StatusCode, eo)
	}
	return resp, nil
}

func (pd *PagerdutyClient) GetErrorFromResponse(resp *http.Response) (*errorObject, error) {
	var result map[string]errorObject
	if err := DecodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", err)
	}
	s, ok := result["error"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have error field")
	}
	return &s, nil
}

func DecodeJSON(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}

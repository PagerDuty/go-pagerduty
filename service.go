package pagerduty

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-querystring/query"
)

type EmailFilter struct {
	SubjectMode    string `json:"subject_mode,omitempty"`
	SubjectRegex   string `json:"subject_regex,omitempty"`
	BodyMode       string `json:"body_mode,omitempty"`
	BodyRegex      string `json:"body_regex,omitempty"`
	FromEmailMode  string `json:"from_email_mode,omitempty"`
	FromEmailRegex string `json:"from_email_regex,omitempty"`
}

type Integration struct {
	APIObject
	Name                  string        `json:"name,omitempty"`
	Service               APIObject     `json:"service,omitempty"`
	CreatedAt             string        `json:"created_at,omitempty"`
	Vendor                APIObject     `json:"vendor,omitempty"`
	IntegrationEmail      string        `json:"integration_email"`
	EmailIncidentCreation string        `json:"email_incident_creation,omitempty"`
	EmailFilterMode       string        `json:"email_filter_mode"`
	EmailFilters          []EmailFilter `json:"email_filters,omitempty"`
}

type NamedTime struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

type ScheduledAction struct {
	Type      string    `json:"type,omitempty"`
	At        NamedTime `json:"at,omitempty"`
	ToUrgency string    `json:"to_urgency"`
}

type SupportHours struct {
	Type    string `json:"type,omitempty"`
	Urgency string `json:"urgency,omitempty"`
}

type SupportHoursDetails struct {
	Type       string `json:"type,omitempty"`
	Timezone   string `json:"time_zone"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	DaysOfWeek []uint `json:"days_of_week"`
}

type IncidentUrgencyRule struct {
	Type                string       `json:"type,omitempty"`
	DuringSupportHours  SupportHours `json:"during_support_hours,omitempty"`
	OutsideSupportHours SupportHours `json:"outside_support_hours,omitempty"`
}

type Service struct {
	APIObject
	Name                   string              `json:"name,omitempty"`
	Description            string              `json:"description,omitempty"`
	AutoResolveTimeout     uint                `json:"auto_resolve_timeout,omitempty"`
	AcknowledgementTimeout uint                `json:"acknowledgement_timeout,omitempty"`
	CreateAt               string              `json:"created_at,omitempty"`
	Status                 string              `json:"status,omitempty"`
	LastIncidentTimestamp  string              `json:"last_incident_timestamp,omitempty"`
	Integrations           []Integration       `json:"integrations,omitempty"`
	EscalationPolicy       EscalationPolicy    `json:"escalation_policy,omitempty"`
	Teams                  []Team              `json:"teams,omitempty"`
	IncidentUrgencyRule    IncidentUrgencyRule `json:"incident_urgency_rule,omitempty"`
	SupportHours           SupportHoursDetails `json:"support_hours,omitempty"`
	ScheduledActions       []ScheduledAction   `json:"scheduled_actions,omitempty"`
}

type ListServiceOptions struct {
	APIListObject
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	TimeZone string   `url:"time_zone,omitempty"`
	SortBy   string   `url:"sort_by,omitempty"`
	Query    string   `url:"query,omitempty"`
	Includes []string `url:"include,omitempty,brackets"`
}

type ListServiceResponse struct {
	APIListObject
	Services []Service
}

func (c *Client) ListServices(o ListServiceOptions) (*ListServiceResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/services?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListServiceResponse
	return &result, c.decodeJson(resp, &result)
}

type GetServiceOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

func (c *Client) GetService(id string, o GetServiceOptions) (*Service, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/services/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]Service
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	s, ok := result["service"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have service field")
	}
	return &s, nil
}

func (c *Client) CreateService(s Service) error {
	data := make(map[string]Service)
	data["service"] = s
	resp, err := c.Post("/services", data)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		ct, rErr := ioutil.ReadAll(resp.Body)
		if rErr == nil {
			log.Debug(string(ct))
		}
		return fmt.Errorf("Failed to create. HTTP Status code: %d", resp.StatusCode)
	}
	return err
}

func (c *Client) UpdateService(s Service) error {
	_, err := c.Put("/services/"+s.ID, s)
	return err
}

func (c *Client) DeleteService(id string) error {
	_, err := c.Delete("/services/" + id)
	return err
}

func (c *Client) CreateIntegration(id string, i Integration) error {
	_, err := c.Post("/services/"+id+"/integrations", i)
	return err
}

type GetIntegrationOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

func (c *Client) GetIntegration(serviceID, integrationID string, o GetIntegrationOptions) (*Integration, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	var result map[string]Integration
	resp, err := c.Get("/services/" + serviceID + "/integrations/" + integrationID + "?" + v.Encode())
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	i, ok := result["integration"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have integration field")
	}
	return &i, nil
}

func (c *Client) UpdateIntegration(serviceID string, i Integration) error {
	_, err := c.Put("/services/"+serviceID+"/integrations/"+i.ID, i)
	return err
}

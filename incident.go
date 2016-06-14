package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

type Acknowledgement struct {
	At           string
	Acknowledger APIObject
}

type PendingAction struct {
	Type string
	At   string
}

type Assignment struct {
	At       string
	Assignee APIObject
}

type Incident struct {
	APIObject
	IncidentNumber       uint              `json:"incident_number,omitempty"`
	CreatedAt            string            `json:"created_at,omitempty"`
	PendingActions       []PendingAction   `json:"pending_actions,omitempty"`
	IncidentKey          string            `json:"incident_key,omitempty"`
	Service              APIObject         `json:"service,omitempty"`
	Assignments          []Assignment      `json:"assignments,omitempty"`
	Acknowledgements     []Acknowledgement `json:"acknowledgements,omitempty"`
	LastStatusChangeAt   string            `json:"last_status_change_at,omitempty"`
	LastStatusChangeBy   APIObject         `json:"last_status_change_by,omitempty"`
	FirstTriggerLogEntry APIObject         `json:"last_trigger_log_entry,omitempty"`
	EscalationPolicy     APIObject         `json:"escalation_policy,omitempty"`
	Teams                []APIObject       `json:"teams,omitempty"`
	Urgency              string            `json:"urgency,omitempty"`
}

type ListIncidentsResponse struct {
	APIListObject
	Incidents []Incident `json:"incidents,omitempty"`
}

type ListIncidentsOptions struct {
	APIListObject
	Since       string   `url:"since,omitempty"`
	Until       string   `url:"until,omitempty"`
	DateRange   string   `url:"date_range,omitempty"`
	Statuses    []string `url:"statuses,omitempty,brackets"`
	IncidentKey string   `url:"incident_key,omitempty"`
	ServiceIDs  []string `url:"service_ids,omitempty,brackets"`
	TeamIDs     []string `url:"team_ids,omitempty,brackets"`
	UserIDs     []string `url:"user_ids,omitempty,brackets"`
	Urgencies   []string `url:"urgencies,omitempty,brackets"`
	TimeZone    string   `url:"time_zone,omitempty"`
	SortBy      string   `url:"sort_by,omitempty"`
	Includes    []string `url:"include,omitempty,brackets"`
}

func (c *Client) ListIncidents(o ListIncidentsOptions) (*ListIncidentsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/incidents?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListIncidentsResponse
	return &result, c.decodeJson(resp, &result)
}

type ManageIncidentsOptions struct {
	From string `url:"from,omitempty"`
}

func (c *Client) ManageIncidents(incidents []Incident, o ManageIncidentsOptions) error {
	v, err := query.Values(o)
	if err != nil {
		return err
	}
	r := make(map[string][]Incident)
	r["incidents"] = incidents
	_, e := c.Put("/incidents?"+v.Encode(), r)
	return e
}

func (c *Client) GetIncident(id string) (*Incident, error) {
	resp, err := c.Get("/incidents/" + id)
	if err != nil {
		return nil, err
	}
	var result map[string]Incident
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	i, ok := result["incident"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have incident field")
	}
	return &i, nil
}

type IncidentNote struct {
	ID        string    `json:"id,omitempty"`
	User      APIObject `json:"user,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt string    `json:"created_at,omitempty"`
}

func (c *Client) ListIncidentNotes(id string) ([]IncidentNote, error) {
	resp, err := c.Get("/incidents/" + id + "/notes")
	if err != nil {
		return nil, err
	}
	var result map[string][]IncidentNote
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	notes, ok := result["notes"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have notes field")
	}
	return notes, nil
}

func (c *Client) CreateIncidentNote(id string, note IncidentNote) error {
	data := make(map[string]IncidentNote)
	data["note"] = note
	_, err := c.Post("/incidents/"+id+"/notes", data)
	return err
}

func (c *Client) SnoozeIncident(id string, duration uint) error {
	data := make(map[string]uint)
	data["duration"] = duration
	_, err := c.Post("/incidents/"+id+"/snooze", data)
	return err
}

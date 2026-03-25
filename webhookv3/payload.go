package webhookv3

import (
	"encoding/json"
	"time"
)

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#webhook-payload
type Payload struct {
	Event Event `json:"event"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#events
type Event struct {
	ID           string          `json:"id"`
	EventType    string          `json:"event_type"`
	ResourceType string          `json:"resource_type"`
	OccurredAt   time.Time       `json:"occurred_at"`
	Agent        Agent           `json:"agent"`
	Client       Client          `json:"client"`
	Data         json.RawMessage `json:"data"`
}

type Agent struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type Client struct {
	Name string `json:"name"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#incident
type IncidentData struct {
	ID               string           `json:"id"`
	Type             string           `json:"type"`
	Self             string           `json:"self"`
	HTMLURL          string           `json:"html_url"`
	Number           int              `json:"number"`
	Status           string           `json:"status"`
	IncidentKey      string           `json:"incident_key"`
	CreatedAt        time.Time        `json:"created_at"`
	Title            string           `json:"title"`
	Service          Service          `json:"service"`
	Assignees        []Assignees      `json:"assignees"`
	EscalationPolicy EscalationPolicy `json:"escalation_policy"`
	Teams            []Teams          `json:"teams"`
	Priority         Priority         `json:"priority"`
	Urgency          string           `json:"urgency"`
	ConferenceBridge ConferenceBridge `json:"conference_bridge"`
	ResolveReason    any              `json:"resolve_reason"`
}

type Service struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type Assignees struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type EscalationPolicy struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type Teams struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type Priority struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type ConferenceBridge struct {
	ConferenceNumber string `json:"conference_number"`
	ConferenceURL    string `json:"conference_url"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#incident_status_update
type IncidentConferenceBridgeData struct {
	Incident          Incident            `json:"incident"`
	ConferenceNumbers []ConferenceNumbers `json:"conference_numbers"`
	ConferenceURL     string              `json:"conference_url"`
	Type              string              `json:"type"`
}

type Incident struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type ConferenceNumbers struct {
	Label  string `json:"label"`
	Number string `json:"number"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#incident_field_values
type IncidentFieldValuesData struct {
	Incident            Incident              `json:"incident"`
	CustomFields        []CustomFields        `json:"custom_fields"`
	ChangedCustomFields []ChangedCustomFields `json:"changed_custom_fields"`
	Type                string                `json:"type"`
}

type CustomFields struct {
	DataType  string `json:"data_type"`
	FieldType string `json:"field_type"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type ChangedCustomFields struct {
	DataType  string `json:"data_type"`
	FieldType string `json:"field_type"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#incident_note
type IncidentNoteData struct {
	Incident Incident `json:"incident"`
	ID       string   `json:"id"`
	Content  string   `json:"content"`
	Trimmed  bool     `json:"trimmed"`
	Type     string   `json:"type"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#incident_status_update
type IncidentStatusUpdateData struct {
	Incident Incident `json:"incident"`
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Trimmed  bool     `json:"trimmed"`
	Type     string   `json:"type"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#incident_responder
type IncidentResponderData struct {
	Incident         Incident         `json:"incident"`
	User             User             `json:"user"`
	EscalationPolicy EscalationPolicy `json:"escalation_policy"`
	Message          string           `json:"message"`
	State            string           `json:"state"`
	Type             string           `json:"type"`
}

type User struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#incident_workflow_instance
type IncidentWorkflowInstanceData struct {
	ID               string           `json:"id"`
	Type             string           `json:"type"`
	Summary          string           `json:"summary"`
	IncidentWorkflow IncidentWorkflow `json:"incident_workflow"`
	WorkflowTrigger  WorkflowTrigger  `json:"workflow_trigger"`
	Incident         Incident         `json:"incident"`
	Service          Service          `json:"service"`
}

type IncidentWorkflow struct {
	HTMLURL string `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type WorkflowTrigger struct {
	HTMLURL any    `json:"html_url"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

// See: https://developer.pagerduty.com/docs/db0fa8c8984fc-overview#service
type ServiceData struct {
	HTMLURL       string  `json:"html_url"`
	ID            string  `json:"id"`
	Self          string  `json:"self"`
	Summary       string  `json:"summary"`
	AlertCreation string  `json:"alert_creation"`
	Teams         []Teams `json:"teams"`
	Type          string  `json:"type"`
}

// See: https://support.pagerduty.com/docs/webhooks#send-a-test-event
type PageyData struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

package pagerduty

import (
	"context"
	"encoding/json"
	"io"
	"time"
)

// IncidentDetails contains a representation of the incident associated with the action that caused this webhook message
type IncidentDetails struct {
	APIObject
	IncidentNumber       int               `json:"incident_number"`
	Title                string            `json:"title"`
	CreatedAt            time.Time         `json:"created_at"`
	Status               string            `json:"status"`
	IncidentKey          *string           `json:"incident_key"`
	PendingActions       []PendingAction   `json:"pending_actions"`
	Service              Service           `json:"service"`
	Assignments          []Assignment      `json:"assignments"`
	Acknowledgements     []Acknowledgement `json:"acknowledgements"`
	LastStatusChangeAt   time.Time         `json:"last_status_change_at"`
	LastStatusChangeBy   APIObject         `json:"last_status_change_by"`
	FirstTriggerLogEntry APIObject         `json:"first_trigger_log_entry"`
	EscalationPolicy     APIObject         `json:"escalation_policy"`
	Teams                []APIObject       `json:"teams"`
	Priority             Priority          `json:"priority"`
	Urgency              string            `json:"urgency"`
	ResolveReason        *string           `json:"resolve_reason"`
	AlertCounts          AlertCounts       `json:"alert_counts"`
	Metadata             interface{}       `json:"metadata"`

	// Alerts is the list of alerts within this incident. Each item in the slice
	// is not fully hydrated, so only the AlertKey field will be set.
	Alerts []IncidentAlert `json:"alerts,omitempty"`

	// Description is deprecated, use Title instead.
	Description string `json:"description"`
}

// WebhookPayloadMessages is the wrapper around the Webhook payloads. The Array may contain multiple message elements if webhook firing actions occurred in quick succession
type WebhookPayloadMessages struct {
	Messages []WebhookPayload `json:"messages"`
}

// WebhookPayload represents the V2 webhook payload
type WebhookPayload struct {
	ID         string          `json:"id"`
	Event      string          `json:"event"`
	CreatedOn  time.Time       `json:"created_on"`
	Incident   IncidentDetails `json:"incident"`
	LogEntries []LogEntry      `json:"log_entries"`
}

// DecodeWebhook decodes a webhook from a response object.
func DecodeWebhook(r io.Reader) (*WebhookPayloadMessages, error) {
	var payload WebhookPayloadMessages
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

type DeliveryMethod struct {
	Url                 string   `json:"url"`
	Type                string   `json:"type"`
	CustomHeaders       []string `json:"custom_headers"`
	TemporarilyDisabled bool     `json:"temporarily_disabled"`
	ID                  string   `json:"id"`
	Secret              string   `json:"secret"`
	ExtensionID         string   `json:"extension_id"`
}

type Filter struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type WebhookSubscription struct {
	APIObject
	ID             string         `json:"id"`
	Type           string         `json:"type"`
	Active         bool           `json:"active"`
	DeliveryMethod DeliveryMethod `json:"delivery_method"`
	Description    string         `json:"description"`
	Events         []string       `json:"events"`
	Filter         Filter         `json:"filter"`
}

type CreateWebhookOptions struct {
	Type           string         `json:"type"`
	Active         bool           `json:"active"`
	DeliveryMethod DeliveryMethod `json:"delivery_method"`
	Description    string         `json:"description"`
	Events         []string       `json:"events"`
	Filter         Filter         `json:"filter"`
}

type WebhookResponse struct {
	WebhookSubscription WebhookSubscription `json:"webhook_subscription"`
}

type UpdateWebhookOptions struct {
	// pointer fields here are used to allow us to omit certain fields when updating
	Active      *bool    `json:"active,omitempty"`
	Description *string  `json:"description,omitempty"`
	Events      []string `json:"events,omitempty"`
	Filter      *Filter  `json:"filter,omitempty"`
}

// CreateWebhookWithContext creates a new webhook.
func (c *Client) CreateWebhookWithContext(ctx context.Context, o *CreateWebhookOptions) (*WebhookSubscription, error) {
	b := map[string]*CreateWebhookOptions{
		"webhook_subscription": o,
	}

	resp, err := c.post(ctx, "/webhook_subscriptions", b, nil)

	if err != nil {
		return nil, err
	}

	var ii WebhookResponse
	if err = c.decodeJSON(resp, &ii); err != nil {
		return nil, err
	}

	return &ii.WebhookSubscription, nil
}

// UpdateWebhookWithContext creates a new webhook.
func (c *Client) UpdateWebhookWithContext(ctx context.Context, id string, o *UpdateWebhookOptions) (*WebhookSubscription, error) {
	b := map[string]*UpdateWebhookOptions{
		"webhook_subscription": o,
	}

	resp, err := c.put(ctx, "/webhook_subscriptions/"+id, b, nil)

	if err != nil {
		return nil, err
	}

	var ii WebhookResponse
	if err = c.decodeJSON(resp, &ii); err != nil {
		return nil, err
	}

	return &ii.WebhookSubscription, nil
}

// GetWebhookWithContext returns information about a specific webhook by ID
func (c *Client) GetWebhookWithContext(ctx context.Context, id string) (*WebhookSubscription, error) {
	resp, err := c.get(ctx, "/webhook_subscriptions/"+id)

	if err != nil {
		return nil, err
	}

	var ii WebhookResponse
	if err = c.decodeJSON(resp, &ii); err != nil {
		return nil, err
	}

	return &ii.WebhookSubscription, nil

}

package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

const changeEventPath = "/v2/change/enqueue"

// ChangeEvent represents a ChangeEvent's request parameters
// https://developer.pagerduty.com/docs/events-api-v2/send-change-events/#parameters
type ChangeEvent struct {
	RoutingKey string             `json:"routing_key"`
	Payload    ChangeEventPayload `json:"payload"`
	Links      []ChangeEventLink  `json:"links"`
}

// ChangeEventPayload ChangeEvent ChangeEventPayload
// https://developer.pagerduty.com/docs/events-api-v2/send-change-events/#example-request-payload
type ChangeEventPayload struct {
	Source        string                 `json:"source"`
	Summary       string                 `json:"summary"`
	Timestamp     string                 `json:"timestamp"`
	CustomDetails map[string]interface{} `json:"custom_details"`
}

// ChangeEventLink represents a single link in a ChangeEvent
// https://developer.pagerduty.com/docs/events-api-v2/send-change-events/#the-links-property
type ChangeEventLink struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

// ChangeEventResponse is the json response body for an event
type ChangeEventResponse struct {
	Status  string   `json:"status,omitempty"`
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// CreateChangeEvent Sends PagerDuty a single ChangeEvent to record
// The v2EventsAPIEndpoint parameter must be set on the client
// Documentation can be found at https://developer.pagerduty.com/docs/events-api-v2/send-change-events
//
// It's recommended to use CreateChangeEventWithContext instead.
func (c *Client) CreateChangeEvent(e ChangeEvent) (*ChangeEventResponse, error) {
	return c.CreateChangeEventWithContext(context.Background(), e)
}

// CreateChangeEventWithContext sends PagerDuty a single ChangeEvent to record
// The v2EventsAPIEndpoint parameter must be set on the client Documentation can
// be found at https://developer.pagerduty.com/docs/events-api-v2/send-change-events
func (c *Client) CreateChangeEventWithContext(ctx context.Context, e ChangeEvent) (*ChangeEventResponse, error) {
	if c.v2EventsAPIEndpoint == "" {
		return nil, errors.New("v2EventsAPIEndpoint field must be set on Client")
	}

	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	resp, err := c.doWithEndpoint(
		ctx,
		c.v2EventsAPIEndpoint,
		http.MethodPost,
		changeEventPath,
		false,
		bytes.NewBuffer(data),
		nil,
	)
	if err != nil {
		return nil, err
	}

	var eventResponse ChangeEventResponse

	if err := c.decodeJSON(resp, &eventResponse); err != nil {
		return nil, err
	}

	return &eventResponse, nil
}

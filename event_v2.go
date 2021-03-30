package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// V2Event includes the incident/alert details
type V2Event struct {
	RoutingKey string        `json:"routing_key"`
	Action     string        `json:"event_action"`
	DedupKey   string        `json:"dedup_key,omitempty"`
	Images     []interface{} `json:"images,omitempty"`
	Links      []interface{} `json:"links,omitempty"`
	Client     string        `json:"client,omitempty"`
	ClientURL  string        `json:"client_url,omitempty"`
	Payload    *V2Payload    `json:"payload,omitempty"`
}

// V2Payload represents the individual event details for an event
type V2Payload struct {
	Summary   string      `json:"summary"`
	Source    string      `json:"source"`
	Severity  string      `json:"severity"`
	Timestamp string      `json:"timestamp,omitempty"`
	Component string      `json:"component,omitempty"`
	Group     string      `json:"group,omitempty"`
	Class     string      `json:"class,omitempty"`
	Details   interface{} `json:"custom_details,omitempty"`
}

// V2EventResponse is the json response body for an event
type V2EventResponse struct {
	Status   string   `json:"status,omitempty"`
	DedupKey string   `json:"dedup_key,omitempty"`
	Message  string   `json:"message,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

const v2eventEndPoint = "https://events.pagerduty.com/v2/enqueue"

// ManageEvent handles the trigger, acknowledge, and resolve methods for an
// event. It's recommended to use ManageEventWithContext instead.
func ManageEvent(e V2Event) (*V2EventResponse, error) {
	return ManageEventWithContext(context.Background(), e)
}

// ManageEventWithContext handles the trigger, acknowledge, and resolve methods for an event.
func ManageEventWithContext(ctx context.Context, e V2Event) (*V2EventResponse, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v2eventEndPoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", "go-pagerduty/"+Version)
	req.Header.Set("Content-Type", "application/json")

	// TODO(theckman): switch to a package-local default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }() // explicitly discard error

	if resp.StatusCode != http.StatusAccepted {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
		}

		return nil, fmt.Errorf("HTTP Status Code: %d, Message: %s", resp.StatusCode, string(bytes))
	}

	var eventResponse V2EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, err
	}

	return &eventResponse, nil
}

// ManageEvent handles the trigger, acknowledge, and resolve methods for an
// event. It's recommended to use ManageEventWithContext instead.
func (c *Client) ManageEvent(e *V2Event) (*V2EventResponse, error) {
	return c.ManageEventWithContext(context.Background(), e)
}

// ManageEventWithContext handles the trigger, acknowledge, and resolve methods for an event.
func (c *Client) ManageEventWithContext(ctx context.Context, e *V2Event) (*V2EventResponse, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	resp, err := c.doWithEndpoint(ctx, c.v2EventsAPIEndpoint, http.MethodPost, "/v2/enqueue", false, bytes.NewBuffer(data), nil)
	if err != nil {
		return nil, err
	}

	var result V2EventResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, err
}

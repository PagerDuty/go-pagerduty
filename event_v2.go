package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Event includes the incident/alert details
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

// Payload represents the individual event details for an event
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

// Response is the json response body for an event
type V2EventResponse struct {
	Status   string   `json:"status,omitempty"`
	DedupKey string   `json:"dedup_key,omitempty"`
	Message  string   `json:"message,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

const v2eventEndPoint = "https://events.pagerduty.com/v2/enqueue"

// ManageEvent handles the trigger, acknowledge, and resolve methods for an event
func ManageEvent(e V2Event) (*V2EventResponse, error) {
	return ManageEventWithContext(context.Background(), e)
}

// ManageEventWithContext is the same as ManageEvent with the addition of
// the ability to pass a context.
//
// Callers can pass a context for use in request cancellation.
func ManageEventWithContext(ctx context.Context, e V2Event) (*V2EventResponse, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, v2eventEndPoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-pagerduty/"+Version)
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("HTTP Status Code: %d, Message: %s", resp.StatusCode, string(msg))
	}
	var eventResponse V2EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, err
	}
	return &eventResponse, nil
}

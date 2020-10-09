package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const EventsEndPoint = "https://events.pagerduty.com"
const ChangeEventPath = "v2/change/enqueue"

type ChangeEvent struct {
	RoutingKey string  `json:"routing_key"`
	Payload    Payload `json:"payload"`
	Links      []Link  `json:"links"`
}

type Payload struct {
	Source          string            `json:"source"`
	Summary         string            `json:"summary"`
	Timestamp       time.Time         `json:"time"`
	TimestampString string            `json:"timestamp"`
	CustomDetails   map[string]string `json:"custom_details"`
}

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

// Response is the json response body for an event
type ChangeEventResponse struct {
	Status  string   `json:"status,omitempty"`
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// ManageEvent handles the trigger, acknowledge, and resolve methods for an event
func (c *Client) SendChangeEvent(e ChangeEvent) (*ChangeEventResponse, error) {
	//PagerDuty expects RFC3339 formatted timestamp so we do the conversion here
	e.Payload.TimestampString = e.Payload.Timestamp.Format(time.RFC3339)
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	//Allows custom endpoints to be passed in from the http client for testing and future customer needs
	endPoint := strings.Join([]string{EventsEndPoint, ChangeEventPath}, "/")
	if c.apiEndpoint != "https://api.pagerduty.com" {
		endPoint = strings.Join([]string{c.apiEndpoint, ChangeEventPath}, "/")
	}
	req, _ := http.NewRequest("POST", endPoint, bytes.NewBuffer(data))
	req.Header.Set("User-Agent", "go-pagerduty/"+Version)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("HTTP Status Code: %d, Message: %s", resp.StatusCode, string(b))
	}
	var eventResponse ChangeEventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, err
	}
	return &eventResponse, nil
}

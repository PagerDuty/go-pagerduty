package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

const EventEndPoint = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"

type Event struct {
	ServiceKey  string        `json:"service_key"`
	Type        string        `json:"event_type"`
	IncidentKey string        `json:"incident_key,omitempty"`
	Description string        `json:"description"`
	Client      string        `json:"client,omitempty"`
	ClientURL   string        `json:"client_url,omitempty"`
	Details     interface{}   `json:"details,omitempty"`
	Contexts    []interface{} `json:"contexts,omitempty"`
}

type EventResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	IncidentKey string `json:"incident_key"`
}

func CreateEvent(e Event) (*EventResponse, error) {
	log.Debugln("Endpoint:", EventEndPoint)
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	log.Debugln(string(data))
	req, _ := http.NewRequest("POST", EventEndPoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
	}
	var eventResponse EventResponse
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, err
	}
	return &eventResponse, nil
}

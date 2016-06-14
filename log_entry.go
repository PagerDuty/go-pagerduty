package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

type Agent APIObject
type Channel struct {
	Type string
}

type LogEntry struct {
	APIObject
	CreatedAt    string `json:"created_at"`
	Agent        Agent
	Channel      Channel
	Incident     Incident
	Teams        []Team
	Contexts     []string
	EventDetails map[string]string
}

type ListLogEntryResponse struct {
	APIListObject
	LogEntries []LogEntry `json:"log_entries"`
}

type ListLogEntriesOptions struct {
	APIListObject
	TimeZone   string   `url:"time_zone"`
	Since      string   `url:"omitempty"`
	Until      string   `url:"omitempty"`
	IsOverview bool     `url:"is_overview,omitempty"`
	Includes   []string `url:"include,omitempty,brackets"`
}

func (c *Client) ListLogEntries(o ListLogEntriesOptions) (*ListLogEntryResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/log_entries?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListLogEntryResponse
	return &result, c.decodeJson(resp, &result)
}

type GetLogEntryOptions struct {
	TimeZone string   `url:"timezone,omitempty"`
	Includes []string `url:"include,omitempty,brackets"`
}

func (c *Client) GetLogEntry(id string, o GetLogEntryOptions) (*LogEntry, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/log_entries/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]LogEntry
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	le, ok := result["log_entry"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have log_entry field")
	}
	return &le, nil
}

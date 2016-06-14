package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

type MaintenanceWindow struct {
	APIObject
	SequenceNumber uint   `json:"sequence_number,omitempty"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Description    string
	Services       []APIObject
	Teams          []APIListObject
	CreatedBy      APIListObject `json:"created_by"`
}

type ListMaintenanceWindowsResponse struct {
	APIListObject
	MaintenanceWindows []MaintenanceWindow `json:"maintenance_windows"`
}

type ListMaintenanceWindowsOptions struct {
	APIListObject
	Query      string   `url:"query,omitempty"`
	Includes   []string `url:"include,omitempty,brackets"`
	TeamIDs    []string `url:"team_ids,omitempty,brackets"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
	Filter     string   `url:"filter,omitempty,brackets"`
}

func (c *Client) ListMaintenanceWindows(o ListMaintenanceWindowsOptions) (*ListMaintenanceWindowsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/maintenance_windows?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListMaintenanceWindowsResponse
	return &result, c.decodeJson(resp, &result)
}

func (c *Client) CreateMaintaienanceWindows(m MaintenanceWindow) error {
	data := make(map[string]MaintenanceWindow)
	data["maintenance_window"] = m
	_, err := c.Post("/mainteance_windows", data)
	return err
}

func (c *Client) DeleteMaintenanceWindow(id string) error {
	_, err := c.Delete("/mainteance_windows/" + id)
	return err
}

type GetMaintenanceWindowOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

func (c *Client) GetMaintenanceWindow(id string, o GetMaintenanceWindowOptions) (*MaintenanceWindow, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/mainteance_windows/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]MaintenanceWindow
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	m, ok := result["maintenance_window"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have maintenance window field")
	}
	return &m, nil
}

func (c *Client) UpdateMaintenanceWindow(m MaintenanceWindow) error {
	_, err := c.Put("/maintenance_windows/"+m.ID, m)
	return err
}

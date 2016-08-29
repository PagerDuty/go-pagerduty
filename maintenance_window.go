package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

// MaintenanceWindow is used to temporarily disable one or more services for a set period of time.
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

// ListMaintenanceWindowsResponse is the data structur returned from calling the ListMaintenanceWindows API endpoint.
type ListMaintenanceWindowsResponse struct {
	APIListObject
	MaintenanceWindows []MaintenanceWindow `json:"maintenance_windows"`
}

// ListMaintenanceWindowsOptions is the data structure used when calling the ListMaintenanceWindows API endpoint.
type ListMaintenanceWindowsOptions struct {
	APIListObject
	Query      string   `url:"query,omitempty"`
	Includes   []string `url:"include,omitempty,brackets"`
	TeamIDs    []string `url:"team_ids,omitempty,brackets"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
	Filter     string   `url:"filter,omitempty,brackets"`
}

// ListMaintenanceWindows lists existing maintenance windows, optionally filtered by service and/or team, or whether they are from the past, present or future.
func (c *Client) ListMaintenanceWindows(o ListMaintenanceWindowsOptions) (*ListMaintenanceWindowsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/maintenance_windows?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListMaintenanceWindowsResponse
	return &result, c.decodeJSON(resp, &result)
}

// CreateMaintaienanceWindows creates a new maintenance window for the specified services.
func (c *Client) CreateMaintaienanceWindows(m MaintenanceWindow) error {
	data := make(map[string]MaintenanceWindow)
	data["maintenance_window"] = m
	_, err := c.post("/mainteance_windows", data)
	return err
}

// DeleteMaintenanceWindow deletes an existing maintenance window if it's in the future, or ends it if it's currently on-going.
func (c *Client) DeleteMaintenanceWindow(id string) error {
	_, err := c.delete("/mainteance_windows/" + id)
	return err
}

// GetMaintenanceWindowOptions is the data structure used when calling the GetMaintenanceWindow API endpoint.
type GetMaintenanceWindowOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetMaintenanceWindow gets an existing maintenance window.
func (c *Client) GetMaintenanceWindow(id string, o GetMaintenanceWindowOptions) (*MaintenanceWindow, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/mainteance_windows/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]MaintenanceWindow
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}
	m, ok := result["maintenance_window"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have maintenance window field")
	}
	return &m, nil
}

// UpdateMaintenanceWindow updates an existing maintenance window.
func (c *Client) UpdateMaintenanceWindow(m MaintenanceWindow) error {
	_, err := c.put("/maintenance_windows/"+m.ID, m)
	return err
}

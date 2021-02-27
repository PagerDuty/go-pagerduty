package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// MaintenanceWindow is used to temporarily disable one or more services for a set period of time.
type MaintenanceWindow struct {
	APIObject
	SequenceNumber uint            `json:"sequence_number,omitempty"`
	StartTime      string          `json:"start_time"`
	EndTime        string          `json:"end_time"`
	Description    string          `json:"description"`
	Services       []APIObject     `json:"services"`
	Teams          []APIListObject `json:"teams"`
	CreatedBy      APIListObject   `json:"created_by"`
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

// ListMaintenanceWindows lists existing maintenance windows, optionally
// filtered by service and/or team, or whether they are from the past, present
// or future. It's recommended to use ListMaintenanceWindowsWithContext instead.
func (c *Client) ListMaintenanceWindows(o ListMaintenanceWindowsOptions) (*ListMaintenanceWindowsResponse, error) {
	return c.ListMaintenanceWindowsWithContext(context.Background(), o)
}

// ListMaintenanceWindowsWithContext lists existing maintenance windows,
// optionally filtered by service and/or team, or whether they are from the
// past, present or future.
func (c *Client) ListMaintenanceWindowsWithContext(ctx context.Context, o ListMaintenanceWindowsOptions) (*ListMaintenanceWindowsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/maintenance_windows?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListMaintenanceWindowsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateMaintenanceWindow creates a new maintenance window for the specified
// services. It's recommended to use CreateMaintenanceWindowWithContext instead.
func (c *Client) CreateMaintenanceWindow(from string, o MaintenanceWindow) (*MaintenanceWindow, error) {
	return c.CreateMaintenanceWindowWithContext(context.Background(), from, o)
}

// CreateMaintenanceWindowWithContext creates a new maintenance window for the specified services.
func (c *Client) CreateMaintenanceWindowWithContext(ctx context.Context, from string, o MaintenanceWindow) (*MaintenanceWindow, error) {
	o.Type = "maintenance_window"

	d := map[string]MaintenanceWindow{
		"maintenance_window": o,
	}

	var h map[string]string
	if from != "" {
		h = map[string]string{
			"From": from,
		}
	}

	resp, err := c.post(ctx, "/maintenance_windows", d, h)
	return getMaintenanceWindowFromResponse(c, resp, err)
}

// CreateMaintenanceWindows creates a new maintenance window for the specified services.
// Deprecated: Use `CreateMaintenanceWindowWithContext` instead.
func (c *Client) CreateMaintenanceWindows(o MaintenanceWindow) (*MaintenanceWindow, error) {
	return c.CreateMaintenanceWindowWithContext(context.Background(), "", o)
}

// DeleteMaintenanceWindow deletes an existing maintenance window if it's in the
// future, or ends it if it's currently on-going. It's recommended to use
// DeleteMaintenanceWindowWithContext instead.
func (c *Client) DeleteMaintenanceWindow(id string) error {
	return c.DeleteMaintenanceWindowWithContext(context.Background(), id)
}

// DeleteMaintenanceWindowWithContext deletes an existing maintenance window if it's in the
// future, or ends it if it's currently on-going.
func (c *Client) DeleteMaintenanceWindowWithContext(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/maintenance_windows/"+id)
	return err
}

// GetMaintenanceWindowOptions is the data structure used when calling the GetMaintenanceWindow API endpoint.
type GetMaintenanceWindowOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetMaintenanceWindow gets an existing maintenance window. It's recommended to
// use GetMaintenanceWindowWithContext instead.
func (c *Client) GetMaintenanceWindow(id string, o GetMaintenanceWindowOptions) (*MaintenanceWindow, error) {
	return c.GetMaintenanceWindowWithContext(context.Background(), id, o)
}

// GetMaintenanceWindowWithContext gets an existing maintenance window.
func (c *Client) GetMaintenanceWindowWithContext(ctx context.Context, id string, o GetMaintenanceWindowOptions) (*MaintenanceWindow, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/maintenance_windows/"+id+"?"+v.Encode())
	return getMaintenanceWindowFromResponse(c, resp, err)
}

// UpdateMaintenanceWindow updates an existing maintenance window. It's
// recommended to use UpdateMaintenanceWindowWithContext instead.
func (c *Client) UpdateMaintenanceWindow(m MaintenanceWindow) (*MaintenanceWindow, error) {
	return c.UpdateMaintenanceWindowWithContext(context.Background(), m)
}

// UpdateMaintenanceWindowWithContext updates an existing maintenance window.
func (c *Client) UpdateMaintenanceWindowWithContext(ctx context.Context, m MaintenanceWindow) (*MaintenanceWindow, error) {
	resp, err := c.put(ctx, "/maintenance_windows/"+m.ID, m, nil)
	return getMaintenanceWindowFromResponse(c, resp, err)
}

func getMaintenanceWindowFromResponse(c *Client, resp *http.Response, err error) (*MaintenanceWindow, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]MaintenanceWindow
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "maintenance_window"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

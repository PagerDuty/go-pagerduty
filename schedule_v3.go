package pagerduty

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-querystring/query"
)

// scheduleV3Headers returns the HTTP headers required by PagerDuty's v3 Schedules API.
// The v3 API requires Accept: application/json and the flexible-schedules early access header.
// The early access value can be overridden via PAGERDUTY_SCHEDULE_V3_EARLY_ACCESS env var.
func scheduleV3Headers() map[string]string {
	earlyAccessVal := "flexible-schedules-early-access"
	if val := os.Getenv("PAGERDUTY_SCHEDULE_V3_EARLY_ACCESS"); val != "" {
		earlyAccessVal = val
	}
	return map[string]string{
		"Accept":         "application/json",
		"X-Early-Access": earlyAccessVal,
	}
}

// ScheduleV3 represents a schedule in PagerDuty's v3 API.
type ScheduleV3 struct {
	ID                 string       `json:"id,omitempty"`
	Type               string       `json:"type,omitempty"`
	Name               string       `json:"name"`
	TimeZone           string       `json:"time_zone"`
	Description        string       `json:"description,omitempty"`
	EscalationPolicies []string     `json:"escalation_policies,omitempty"`
	Rotations          []RotationV3 `json:"rotations,omitempty"`
}

// RotationV3 represents a rotation within a v3 schedule.
type RotationV3 struct {
	ID     string    `json:"id,omitempty"`
	Type   string    `json:"type,omitempty"`
	Events []EventV3 `json:"events,omitempty"`
}

// EventTimeV3 represents a time field in a v3 schedule event.
// The v3 API uses an object {"date_time": "...", "time_zone": "..."} rather than a plain string.
type EventTimeV3 struct {
	DateTime string `json:"date_time"`
	TimeZone string `json:"time_zone,omitempty"`
}

// EventV3 represents an on-call event configuration within a rotation.
type EventV3 struct {
	ID                 string               `json:"id,omitempty"`
	Type               string               `json:"type,omitempty"`
	Name               string               `json:"name"`
	StartTime          EventTimeV3          `json:"start_time"`
	EndTime            EventTimeV3          `json:"end_time"`
	EffectiveSince     string               `json:"effective_since"`
	EffectiveUntil     *string              `json:"effective_until"`
	Recurrence         []string             `json:"recurrence"`
	AssignmentStrategy AssignmentStrategyV3 `json:"assignment_strategy"`
}

// AssignmentStrategyV3 defines how on-call responsibility is assigned within an event.
type AssignmentStrategyV3 struct {
	Type    string     `json:"type"`
	Members []MemberV3 `json:"members"`
}

// MemberV3 represents a member in an assignment strategy.
type MemberV3 struct {
	Type   string  `json:"type"`
	UserID *string `json:"user_id,omitempty"`
}

// ScheduleV3Input contains the mutable fields for creating or updating a v3 schedule.
type ScheduleV3Input struct {
	Name        string `json:"name"`
	TimeZone    string `json:"time_zone"`
	Description string `json:"description,omitempty"`
}

// scheduleV3Payload is the request body shape for v3 schedule create/update.
// The rotations field must be present (even as []) per the v3 API validation.
type scheduleV3Payload struct {
	Name        string       `json:"name"`
	TimeZone    string       `json:"time_zone"`
	Description string       `json:"description,omitempty"`
	Rotations   []RotationV3 `json:"rotations"`
}

type createScheduleV3Request struct {
	Schedule scheduleV3Payload `json:"schedule"`
}

type updateScheduleV3Request struct {
	Schedule scheduleV3Payload `json:"schedule"`
}

type scheduleV3Response struct {
	Schedule ScheduleV3 `json:"schedule"`
}

// ListSchedulesV3Response is the response for listing v3 schedules.
type ListSchedulesV3Response struct {
	Schedules []APIObject `json:"schedules"`
}

// ListSchedulesV3Options are query parameters for listing v3 schedules.
type ListSchedulesV3Options struct {
	Query  string `url:"query,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

type rotationV3Response struct {
	Rotation RotationV3 `json:"rotation"`
}

// eventV3Response wraps the v3 API event response. The v3 API uses "schedule_event" as the key.
type eventV3Response struct {
	Event EventV3 `json:"schedule_event"`
}

// ListSchedulesV3 retrieves a paginated list of v3 schedules.
func (c *Client) ListSchedulesV3(ctx context.Context, o ListSchedulesV3Options) (*ListSchedulesV3Response, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	path := "/v3/schedules"
	if encoded := v.Encode(); encoded != "" {
		path += "?" + encoded
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result ListSchedulesV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetScheduleV3 retrieves a v3 schedule by ID, including its rotations and events.
func (c *Client) GetScheduleV3(ctx context.Context, id string) (*ScheduleV3, error) {
	resp, err := c.get(ctx, "/v3/schedules/"+id, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result scheduleV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Schedule, nil
}

// CreateScheduleV3 creates a new v3 schedule with metadata only.
// Rotations and events must be added via separate API calls.
func (c *Client) CreateScheduleV3(ctx context.Context, s ScheduleV3Input) (*ScheduleV3, error) {
	d := createScheduleV3Request{Schedule: scheduleV3Payload{
		Name:        s.Name,
		TimeZone:    s.TimeZone,
		Description: s.Description,
		Rotations:   []RotationV3{},
	}}

	resp, err := c.post(ctx, "/v3/schedules", d, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 schedule, HTTP status code: %d", resp.StatusCode)
	}

	var result scheduleV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Schedule, nil
}

// UpdateScheduleV3 updates a v3 schedule's metadata (name, time_zone, description).
func (c *Client) UpdateScheduleV3(ctx context.Context, id string, s ScheduleV3Input) (*ScheduleV3, error) {
	d := updateScheduleV3Request{Schedule: scheduleV3Payload{
		Name:        s.Name,
		TimeZone:    s.TimeZone,
		Description: s.Description,
		Rotations:   []RotationV3{},
	}}

	resp, err := c.put(ctx, "/v3/schedules/"+id, d, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result scheduleV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Schedule, nil
}

// DeleteScheduleV3 soft-deletes a v3 schedule and all its rotations and events.
func (c *Client) DeleteScheduleV3(ctx context.Context, id string) error {
	// Use do() directly since delete() does not accept custom headers
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+id, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// CreateRotationV3 creates a new empty rotation for a v3 schedule.
// Events are added to the rotation via CreateEventV3.
func (c *Client) CreateRotationV3(ctx context.Context, scheduleID string) (*RotationV3, error) {
	resp, err := c.post(ctx, "/v3/schedules/"+scheduleID+"/rotations", map[string]interface{}{}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 rotation, HTTP status code: %d", resp.StatusCode)
	}

	var result rotationV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Rotation, nil
}

// GetRotationV3 retrieves a rotation by ID for a given schedule.
func (c *Client) GetRotationV3(ctx context.Context, scheduleID, rotationID string) (*RotationV3, error) {
	resp, err := c.get(ctx, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result rotationV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Rotation, nil
}

// DeleteRotationV3 soft-deletes a rotation from a v3 schedule.
func (c *Client) DeleteRotationV3(ctx context.Context, scheduleID, rotationID string) error {
	// Use do() directly since delete() does not accept custom headers
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// CreateEventV3 creates a new event within a v3 rotation.
func (c *Client) CreateEventV3(ctx context.Context, scheduleID, rotationID string, e EventV3) (*EventV3, error) {
	resp, err := c.post(ctx, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events", e, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 event, HTTP status code: %d", resp.StatusCode)
	}

	var result eventV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Event, nil
}

// UpdateEventV3 updates an event within a v3 rotation.
func (c *Client) UpdateEventV3(ctx context.Context, scheduleID, rotationID, eventID string, e EventV3) (*EventV3, error) {
	resp, err := c.put(ctx, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events/"+eventID, e, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result eventV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Event, nil
}

// DeleteEventV3 deletes an event from a v3 rotation.
func (c *Client) DeleteEventV3(ctx context.Context, scheduleID, rotationID, eventID string) error {
	// Use do() directly since delete() does not accept custom headers
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events/"+eventID, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

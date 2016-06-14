package pagerduty

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

type Restriction struct {
	Type            string
	StartTimeOfDay  string `json:"start_time_of_day,omitempty"`
	DurationSeconds uint   `json:"duration_seconds,omitempty"`
}

type RenderedScheduleEntry struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
	User  APIObject
}

type ScheduleLayer struct {
	APIObject
	Name                       string                  `json:"name,omitempty"`
	Start                      string                  `json:"start,omitempty"`
	End                        string                  `json:"end,omitempty"`
	RotationVirtualStart       string                  `json:"rotation_virtual_start,omitempty"`
	RotationTurnLengthSeconds  uint                    `json:"rotation_virtual_start,omitempty"`
	Users                      []APIObject             `json:"users,omitempty"`
	Restrictions               []Restriction           `json:"restrictions,omitempty"`
	RenderedScheduleEntries    []RenderedScheduleEntry `json:"rendered_schedule_entries,omitempty"`
	RenderedCoveragePercentage float64                 `json:"rendered_coverage_percentage,omitempty"`
}

type Schedule struct {
	APIObject
	Name                 string          `json:"name,omitempty"`
	TimeZone             string          `json:"time_zone,omitempty"`
	Desciption           string          `json:"description,omitempty"`
	EscalationPolicies   []APIObject     `json:"escalation_policies,omitempty"`
	Users                []APIObject     `json:"users,omitempty"`
	ScheduleLayers       []ScheduleLayer `json:"schedule_layers,omitempty"`
	OverridesSubschedule ScheduleLayer   `json:"override_subschedule,omitempty"`
	FinalSchedule        ScheduleLayer   `json:"final_schedule,omitempty"`
}

type ListSchedulesOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

type ListSchedulesResponse struct {
	APIListObject
	Schedules []Schedule
}

func (c *Client) ListSchedules(o ListSchedulesOptions) (*ListSchedulesResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/schedules?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListSchedulesResponse
	return &result, c.decodeJson(resp, &result)
}

func (c *Client) CreateSchedule(s Schedule) error {
	data := make(map[string]Schedule)
	data["schedule"] = s
	_, err := c.Post("/schedules", data)
	return err
}

type PreviewScheduleOptions struct {
	APIListObject
	Since    string `url:"since,omitempty"`
	Until    string `url:"until,omitempty"`
	Overflow bool   `url:"overflow,omitempty"`
}

func (c *Client) PreviewSchedule(s Schedule, o PreviewScheduleOptions) error {
	v, err := query.Values(o)
	if err != nil {
		return err
	}
	var data map[string]Schedule
	data["schedule"] = s
	_, e := c.Post("/schedules/preview?"+v.Encode(), data)
	return e
}

func (c *Client) DeleteSchedule(id string) error {
	_, err := c.Delete("/schedules/" + id)
	return err
}

type GetScheduleOptions struct {
	APIListObject
	TimeZone string `url:"time_zone,omitempty"`
	Since    string `url:"since,omitempty"`
	Until    string `url:"until,omitempty"`
}

func (c *Client) GetSchedule(id string, o GetScheduleOptions) (*Schedule, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/schedules/" + id + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string]Schedule
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	s, ok := result["schedule"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have schedule field")
	}
	return &s, nil
}

type UpdateScheduleOptions struct {
	Overflow bool `url:"overflow,omitempty"`
}

func (c *Client) UpdateSchedule(id string, s Schedule) error {
	v := make(map[string]Schedule)
	v["schedule"] = s
	_, err := c.Put("/schedules/"+id, v)
	return err
}

type ListOverridesOptions struct {
	APIListObject
	Since    string `url:"since,omitempty"`
	Until    string `url:"until,omitempty"`
	Editable bool   `url:"editable,omitempty"`
	Overflow bool   `url:"overflow,omitempty"`
}

type Overrides struct {
	ID    string    `json:"id,omitempty"`
	Start string    `json:"start,omitempty"`
	End   string    `json:"end,omitempty"`
	User  APIObject `json:"user,omitempty"`
}

func (c *Client) ListOverrides(id string, o ListOverridesOptions) ([]Overrides, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/schedules/" + id + "/overrides?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string][]Overrides
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	overrides, ok := result["overrides"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have overrides field")
	}
	return overrides, nil
}

func (c *Client) CreateOverride(id string, o Overrides) error {
	_, err := c.Post("/schedules/"+id+"/overrides", o)
	return err
}

func (c *Client) DeleteOverride(scheduleID, overrideID string) error {
	_, err := c.Delete("/schedules/" + scheduleID + "/overrides/" + overrideID)
	return err
}

type ListOnCallUsersOptions struct {
	APIListObject
	Since string `url:"since,omitempty"`
	Until string `url:"until,omitempty"`
}

func (c *Client) ListOnCallUsers(id string, o ListOnCallUsersOptions) ([]User, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/schedules/" + id + "/users?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result map[string][]User
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	u, ok := result["users"]
	if !ok {
		return nil, fmt.Errorf("JSON responsde does not have users field")
	}
	return u, nil
}

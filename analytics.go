package pagerduty

import (
	"context"
	"fmt"
)

const analyticsBaseURL = "/analytics/metrics/incidents"
const rawDataBaseURL = "/analytics/raw"

// AnalyticsRequest represents the request to be sent to PagerDuty when you want
// aggregated analytics.
type AnalyticsRequest struct {
	Filters       *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit string           `json:"aggregate_unit,omitempty"`
	TimeZone      string           `json:"time_zone,omitempty"`
}

// AnalyticsRawIncidentsRequest represents the request to be sent to PagerDuty when you want
// raw analytics.
type AnalyticsRawIncidentsRequest struct {
	Filters       *AnalyticsFilter `json:"filters,omitempty"`
	StartingAfter string           `json:"starting_after,omitempty"`
	EndingBefore  string           `json:"ending_before,omitempty"`
	Limit         int              `json:"limit,omitempty"`
	Order         string           `json:"order,omitempty"`
	OrderBy       string           `json:"order_by,omitempty"`
	TimeZone      string           `json:"time_zone,omitempty"`
}

// AnalyticsResponse represents the response from the PagerDuty API.
type AnalyticsResponse struct {
	Data          []AnalyticsData  `json:"data,omitempty"`
	Filters       *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit string           `json:"aggregate_unit,omitempty"`
	TimeZone      string           `json:"time_zone,omitempty"`
}

// AnalyticsRawIncidentsResponse represents the response from the PagerDuty API.
type AnalyticsRawIncidentsResponse struct {
	First    string                 `json:"first,omitempty"`
	Last     string                 `json:"last,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	More     bool                   `json:"more,omitempty"`
	Order    string                 `json:"order,omitempty"`
	OrderBy  string                 `json:"order_by,omitempty"`
	Filters  *AnalyticsFilter       `json:"filters,omitempty"`
	TimeZone string                 `json:"time_zone,omitempty"`
	Data     []AnalyticsRawIncident `json:"data,omitempty"`
}

// AnalyticsFilter represents the set of filters as part of the request to PagerDuty when
// requesting analytics.
type AnalyticsFilter struct {
	CreatedAtStart string   `json:"created_at_start,omitempty"`
	CreatedAtEnd   string   `json:"created_at_end,omitempty"`
	Urgency        string   `json:"urgency,omitempty"`
	Major          bool     `json:"major,omitempty"`
	ServiceIDs     []string `json:"service_ids,omitempty"`
	TeamIDs        []string `json:"team_ids,omitempty"`
	PriorityIDs    []string `json:"priority_ids,omitempty"`
	PriorityNames  []string `json:"priority_names,omitempty"`
}

// AnalyticsData represents the structure of the aggregated analytics we have available.
type AnalyticsData struct {
	ServiceID                      string  `json:"service_id,omitempty"`
	ServiceName                    string  `json:"service_name,omitempty"`
	TeamID                         string  `json:"team_id,omitempty"`
	TeamName                       string  `json:"team_name,omitempty"`
	MeanSecondsToResolve           int     `json:"mean_seconds_to_resolve,omitempty"`
	MeanSecondsToFirstAck          int     `json:"mean_seconds_to_first_ack,omitempty"`
	MeanSecondsToEngage            int     `json:"mean_seconds_to_engage,omitempty"`
	MeanSecondsToMobilize          int     `json:"mean_seconds_to_mobilize,omitempty"`
	MeanEngagedSeconds             int     `json:"mean_engaged_seconds,omitempty"`
	MeanEngagedUserCount           int     `json:"mean_engaged_user_count,omitempty"`
	TotalEscalationCount           int     `json:"total_escalation_count,omitempty"`
	MeanAssignmentCount            int     `json:"mean_assignment_count,omitempty"`
	TotalBusinessHourInterruptions int     `json:"total_business_hour_interruptions,omitempty"`
	TotalSleepHourInterruptions    int     `json:"total_sleep_hour_interruptions,omitempty"`
	TotalOffHourInterruptions      int     `json:"total_off_hour_interruptions,omitempty"`
	TotalSnoozedSeconds            int     `json:"total_snoozed_seconds,omitempty"`
	TotalEngagedSeconds            int     `json:"total_engaged_seconds,omitempty"`
	TotalIncidentCount             int     `json:"total_incident_count,omitempty"`
	UpTimePct                      float64 `json:"up_time_pct,omitempty"`
	UserDefinedEffortSeconds       int     `json:"user_defined_effort_seconds,omitempty"`
	RangeStart                     string  `json:"range_start,omitempty"`
}

// AnalyticsRawIncident represents the structure of the raw incident analytics we have available.
type AnalyticsRawIncident struct {
	AssignmentCount           int    `json:"assignment_count,omitempty"`
	BusinessHourInterruptions int    `json:"business_hour_interruptions,omitempty"`
	CreatedAt                 string `json:"created_at,omitempty"`
	Description               string `json:"description,omitempty"`
	EngagedSeconds            int    `json:"engaged_seconds,omitempty"`
	EngagedUserCount          int    `json:"engaged_user_count,omitempty"`
	EscalationCount           int    `json:"escalation_count,omitempty"`
	ID                        string `json:"id,omitempty"`
	IncidentNumber            int    `json:"incident_number,omitempty"`
	IsMajor                   bool   `json:"major,omitempty"`
	OffHourInterruptions      int    `json:"off_hour_interruptions,omitempty"`
	PriorityID                string `json:"priority_id,omitempty"`
	PriorityName              string `json:"priority_name,omitempty"`
	ResolvedAt                string `json:"resolved_at,omitempty"`
	SecondsToEngage           int    `json:"seconds_to_engage,omitempty"`
	SecondsToFirstAck         int    `json:"seconds_to_first_ack,omitempty"`
	SecondsToMobilize         int    `json:"seconds_to_mobilize,omitempty"`
	SecondsToResolve          int    `json:"seconds_to_resolve,omitempty"`
	ServiceID                 string `json:"service_id,omitempty"`
	ServiceName               string `json:"service_name,omitempty"`
	SleepHourInterruptions    int    `json:"sleep_hour_interruptions,omitempty"`
	SnoozedSeconds            int    `json:"snoozed_seconds,omitempty"`
	TeamID                    string `json:"team_id,omitempty"`
	TeamName                  string `json:"team_name,omitempty"`
	Urgency                   string `json:"urgency,omitempty"`
	UserDefinedEffortSeconds  int    `json:"user_defined_effort_seconds,omitempty"`
}

// GetAggregatedIncidentData gets the aggregated incident analytics for the requested data.
func (c *Client) GetAggregatedIncidentData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	return c.getAggregatedData(ctx, analytics, "all")
}

// GetAggregatedServiceData gets the aggregated service analytics for the requested data.
func (c *Client) GetAggregatedServiceData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	return c.getAggregatedData(ctx, analytics, "services")
}

// GetAggregatedTeamData gets the aggregated team analytics for the requested data.
func (c *Client) GetAggregatedTeamData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	return c.getAggregatedData(ctx, analytics, "teams")
}

func (c *Client) getAggregatedData(ctx context.Context, analytics AnalyticsRequest, endpoint string) (AnalyticsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
	}

	u := fmt.Sprintf("%s/%s", analyticsBaseURL, endpoint)
	resp, err := c.post(ctx, u, analytics, h)
	if err != nil {
		return AnalyticsResponse{}, err
	}

	var analyticsResponse AnalyticsResponse
	if err = c.decodeJSON(resp, &analyticsResponse); err != nil {
		return AnalyticsResponse{}, err
	}

	return analyticsResponse, nil
}

// GetAnalyticsIncidentsById gets the raw analytics for the requested incident.
func (c *Client) GetAnalyticsIncidentsById(ctx context.Context, id string) (*AnalyticsRawIncident, error) {
	path := fmt.Sprintf("%s/%s/%s", rawDataBaseURL, "incidents", id)
	resp, err := c.get(ctx, path, nil)

	if err != nil {
		return &AnalyticsRawIncident{}, err
	}

	var rawData AnalyticsRawIncident
	if err = c.decodeJSON(resp, &rawData); err != nil {
		return &AnalyticsRawIncident{}, err
	}

	return &rawData, nil
}

// GetAnalyticsIncidents gets the raw analytics for the requested data.
func (c *Client) GetAnalyticsIncidents(ctx context.Context, rawDataReq AnalyticsRawIncidentsRequest) (*AnalyticsRawIncidentsResponse, error) {
	h := map[string]string{}

	path := fmt.Sprintf("%s/%s", rawDataBaseURL, "incidents")
	resp, err := c.post(ctx, path, rawDataReq, h)

	if err != nil {
		return &AnalyticsRawIncidentsResponse{}, err
	}

	var rawDataResponse AnalyticsRawIncidentsResponse
	if err = c.decodeJSON(resp, &rawDataResponse); err != nil {
		return &AnalyticsRawIncidentsResponse{}, err
	}

	return &rawDataResponse, nil
}

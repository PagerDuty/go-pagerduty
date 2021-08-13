package pagerduty

import (
	"context"
	"fmt"
)

// AnalyticsRequest represents the request to be sent to PagerDuty when you want
// aggregated analytics.
type AnalyticsRequest struct {
	Filters       *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit string           `json:"aggregate_unit,omitempty"`
	TimeZone      string           `json:"time_zone,omitempty"`
}

// AnalyticsResponse represents the response from the PagerDuty API.
type AnalyticsResponse struct {
	Data          []AnalyticsData  `json:"data,omitempty"`
	Filters       *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit string           `json:"aggregate_unit,omitempty"`
	TimeZone      string           `json:"time_zone,omitempty"`
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

// AnalyticsData represents the structure of the analytics we have available.
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

type AnalyticsRawIncidentData struct {
	ID                        string `json:"id,omitempty"`
	TeamID                    string `json:"team_id,omitempty"`
	TeamName                  string `json:"team_name,omitempty"`
	ServiceID                 string `json:"service_id,omitempty"`
	ServiceName               string `json:"service_name,omitempty"`
	CreatedAt                 string `json:"created_at,omitempty"`
	ResolvedAt                string `json:"resolved_at,omitempty"`
	Description               string `json:"description,omitempty"`
	IncidentNumber            uint   `json:"incident_number,omitempty"`
	Urgency                   string `json:"urgency,omitempty"`
	Major                     bool   `json:"major,omitempty"`
	PriorityID                string `json:"priority_id,omitempty"`
	PriorityName              string `json:"priority_name,omitempty"`
	PriorityOrder             uint   `json:"priority_order,omitempty"`
	SecondsToResolve          uint   `json:"seconds_to_resolve,omitempty"`
	SecondsToFirstAck         uint   `json:"seconds_to_first_ack,omitempty"`
	SecondsToEngage           uint   `json:"seconds_to_engage,omitempty"`
	SecondsToMobilize         uint   `json:"seconds_to_mobilize,omitempty"`
	EngagedSeconds            uint   `json:"engaged_seconds,omitempty"`
	EngagedUserCount          uint   `json:"engaged_user_count,omitempty"`
	EscalationCount           uint   `json:"escalation_count,omitempty"`
	AssignmentCount           uint   `json:"assignment_count,omitempty"`
	BusinessHourInterruptions uint   `json:"business_hour_interruptions,omitempty"`
	SleepHourInterruptions    uint   `json:"sleep_hour_interruptions,omitempty"`
	OffHourInterruptions      uint   `json:"off_hour_interruptions,omitempty"`
	SnoozedSeconds            uint   `json:"snoozed_seconds,omitempty"`
	UserDefinedEffortsSeconds uint   `json:"user_defined_effort_seconds,omitempty"`
}

type AnalyticsIncidentResponse struct {
	First    string                     `json:"first,omitempty"`
	Last     string                     `json:"last,omitempty"`
	Limit    uint                       `json:"limit,omitempty"`
	More     bool                       `url:"more,omitempty"`
	Order    string                     `json:"order,omitempty"`
	OrderBy  string                     `json:"order_by,omitempty"`
	Filter   *AnalyticsFilter           `json:"filters,omitempty"`
	TimeZone string                     `json:"time_zone,omitempty"`
	Data     []AnalyticsRawIncidentData `json:"data,omitempty"`
}

type IncidentAnalyticsRequest struct {
	Filter        *AnalyticsFilter `json:"filters,omitempty"`
	StartingAfter string           `json:"starting_after,omitempty"`
	EndingBefore  string           `json:"ending_before,omitempty"`
	Order         string           `json:"order,omitempty"`
	OrderBy       string           `json:"order_by,omitempty"`
	Limit         uint             `json:"limit,omitempty"`
	TimeZone      string           `json:"time_zone,omitempty"`
}

// GetAggregatedIncidentData gets the aggregated incident analytics for the requested data.
func (c *Client) GetAggregatedIncidentData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
	}
	resp, err := c.post(ctx, "/analytics/metrics/incidents/all", analytics, h)
	if err != nil {
		return AnalyticsResponse{}, err
	}

	var analyticsResponse AnalyticsResponse
	if err = c.decodeJSON(resp, &analyticsResponse); err != nil {
		return AnalyticsResponse{}, err
	}

	return analyticsResponse, nil
}

// GetAggregatedServiceData gets the aggregated service analytics for the requested data.
func (c *Client) GetAggregatedServiceData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
	}

	resp, err := c.post(ctx, "/analytics/metrics/incidents/services", analytics, h)
	if err != nil {
		return AnalyticsResponse{}, err
	}

	var analyticsResponse AnalyticsResponse
	if err = c.decodeJSON(resp, &analyticsResponse); err != nil {
		return AnalyticsResponse{}, err
	}

	return analyticsResponse, nil
}

// GetAggregatedTeamData gets the aggregated team analytics for the requested data.
func (c *Client) GetAggregatedTeamData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
	}

	resp, err := c.post(ctx, "/analytics/metrics/incidents/teams", analytics, h)
	if err != nil {
		return AnalyticsResponse{}, err
	}

	var analyticsResponse AnalyticsResponse
	if err = c.decodeJSON(resp, &analyticsResponse); err != nil {
		return AnalyticsResponse{}, err
	}

	return analyticsResponse, nil
}

// Get raw data multiple incidents gets the incidents and its metrics for the requested data.
func (c *Client) GetMultipleRawIncidents(o IncidentAnalyticsRequest) (AnalyticsIncidentResponse, error) {
	return c.GetMultipleRawIncidentsWithContext(context.Background(), o)
}

// Gets raw data multiple incidents with context gets the incidents and its metrics for the requested data.
func (c *Client) GetMultipleRawIncidentsWithContext(ctx context.Context, o IncidentAnalyticsRequest) (AnalyticsIncidentResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
	}

	resp, err := c.post(ctx, "/analytics/raw/incidents", o, h)
	if err != nil {
		return AnalyticsIncidentResponse{}, err
	}

	var target AnalyticsIncidentResponse

	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return AnalyticsIncidentResponse{}, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	return target, nil

}

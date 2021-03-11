package pagerduty

import (
	"context"
)

type AnalyticsRequest struct {
	AnalyticsFilter *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit   string           `json:"aggregate_unit,omitempty"`
	TimeZone        string           `json:"time_zone,omitempty"`
}

type AnalyticsResponse struct {
	Data            []AnalyticsData  `json:"data,omitempty"`
	AnalyticsFilter *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit   string           `json:"aggregate_unit,omitempty"`
	TimeZone        string           `json:"time_zone,omitempty"`
}

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
	UpTimePct                      float32 `json:"up_time_pct,omitempty"`
	UserDefinedEffortSeconds       int     `json:"user_defined_effort_seconds,omitempty"`
	RangeStart                     string  `json:"range_start,omitempty"`
}

func (c *Client) GetAggregatedIncidentData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	var analyticsResponse AnalyticsResponse
	headers := make(map[string]string)
	headers["X-EARLY-ACCESS"] = "analytics-v2"

	resp, err := c.post(ctx, "/analytics/metrics/incidents/all", analytics, headers)
	if err != nil {
		return analyticsResponse, err
	}
	err = c.decodeJSON(resp, &analyticsResponse)
	return analyticsResponse, err
}

func (c *Client) GetAggregatedServiceData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	var analyticsResponse AnalyticsResponse
	headers := make(map[string]string)
	headers["X-EARLY-ACCESS"] = "analytics-v2"

	resp, err := c.post(ctx, "/analytics/metrics/incidents/services", analytics, headers)
	if err != nil {
		return analyticsResponse, err
	}
	err = c.decodeJSON(resp, &analyticsResponse)
	return analyticsResponse, err
}

func (c *Client) GetAggregatedTeamData(ctx context.Context, analytics AnalyticsRequest) (AnalyticsResponse, error) {
	var analyticsResponse AnalyticsResponse
	headers := make(map[string]string)
	headers["X-EARLY-ACCESS"] = "analytics-v2"

	resp, err := c.post(ctx, "/analytics/metrics/incidents/teams", analytics, headers)
	if err != nil {
		return analyticsResponse, err
	}
	err = c.decodeJSON(resp, &analyticsResponse)
	return analyticsResponse, err
}
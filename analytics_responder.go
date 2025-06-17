package pagerduty

import (
	"context"
	"fmt"
)

const analyticsResponderBaseURL = "/analytics/metrics/responders"

type AnalyticsResponderRequest struct {
	Filters  *AnalyticsResponderFilter `json:"filters,omitempty"`
	TimeZone string                    `json:"time_zone,omitempty"`
	Order    string                    `json:"order,omitempty"`
	OrderBy  string                    `json:"order_by,omitempty"`
}

type AnalyticsResponderResponse struct {
	Data     []AnalyticsResponderData  `json:"data,omitempty"`
	Filters  *AnalyticsResponderFilter `json:"filters,omitempty"`
	TimeZone string                    `json:"time_zone,omitempty"`
	Order    string                    `json:"order,omitempty"`
	OrderBy  string                    `json:"order_by,omitempty"`
}

type AnalyticsResponderFilter struct {
	DateRangeStart string   `json:"date_range_start,omitempty"`
	DateRangeEnd   string   `json:"date_range_end,omitempty"`
	Urgency        string   `json:"urgency,omitempty"`
	TeamIDs        []string `json:"team_ids,omitempty"`
	ResponderIDs   []string `json:"responder_ids,omitempty"`
	PriorityIDs    []string `json:"priority_ids,omitempty"`
	PriorityNames  []string `json:"priority_names,omitempty"`
}

type AnalyticsResponderData struct {
	MeanEngagedSeconds                int    `json:"mean_engaged_seconds,omitempty"`
	MeanTimeToAckSeconds              int    `json:"mean_time_to_acknowledge_seconds,omitempty"`
	ResponderID                       string `json:"responder_id,omitempty"`
	ResponderName                     string `json:"responder_name,omitempty"`
	TeamID                            string `json:"team_id,omitempty"`
	TeamName                          string `json:"team_name,omitempty"`
	TotalBusinessHourInterruptions    int    `json:"total_business_hour_interruptions,omitempty"`
	TotalEngagedSeconds               int    `json:"total_engaged_seconds,omitempty"`
	TotalIncidentCount                int    `json:"total_incident_count,omitempty"`
	TotalIncidentAck                  int    `json:"total_incidents_acknowledged,omitempty"`
	TotalIncidentManualEscalatedFrom  int    `json:"total_incidents_manual_escalated_from,omitempty"`
	TotalIncidentManualEscalatedTo    int    `json:"total_incidents_manual_escalated_to,omitempty"`
	TotalIncidentReassignedFrom       int    `json:"total_incidents_reassigned_from,omitempty"`
	TotalIncidentReassignedTo         int    `json:"total_incidents_reassigned_to,omitempty"`
	TotalIncidentTimeoutEscalatedFrom int    `json:"total_incidents_timeout_escalated_from,omitempty"`
	TotalIncidentTimeoutEscalatedTo   int    `json:"total_incidents_timeout_escalated_to,omitempty"`
	TotalInterruptions                int    `json:"total_interruptions,omitempty"`
	TotalNotifications                int    `json:"total_notifications,omitempty"`
	TotalOffHourInterruptions         int    `json:"total_off_hour_interruptions,omitempty"`
	TotalSecondsOnCall                int    `json:"total_seconds_on_call,omitempty"`
	TotalSecondsOnCallLevel1          int    `json:"total_seconds_on_call_level_1,omitempty"`
	TotalSecondsOnCallLevel2Plus      int    `json:"total_seconds_on_call_level_2_plus,omitempty"`
	TotalSleepHourInterruptions       int    `json:"total_sleep_hour_interruptions,omitempty"`
}

func (c *Client) GetAggregatedResponderData(ctx context.Context, analytics AnalyticsResponderRequest) (AnalyticsResponderResponse, error) {
	return c.getAggregatedResponderData(ctx, analytics, "all")
}

func (c *Client) getAggregatedResponderData(ctx context.Context, analytics AnalyticsResponderRequest, endpoint string) (AnalyticsResponderResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
	}

	u := fmt.Sprintf("%s/%s", analyticsResponderBaseURL, endpoint)
	resp, err := c.post(ctx, u, analytics, h)
	if err != nil {
		return AnalyticsResponderResponse{}, err
	}

	var analyticsResponse AnalyticsResponderResponse
	if err = c.decodeJSON(resp, &analyticsResponse); err != nil {
		return AnalyticsResponderResponse{}, err
	}

	return analyticsResponse, nil
}

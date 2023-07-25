package pagerduty

import (
	"context"
)

type RawAnalytics struct {
	IncidentID string     `json:"incident_id,omit_empty"`
	Limit      int        `json:"limit,omit_empty"`
	Order      string     `json:"order,omit_empty"`
	OrderBy    string     `json:"order_by,omit_empty"`
	TimeZone   string     `json:"time_zone,omit_empty"`
	Responses  *Responses `json:"responses,omit_empty"`
}

type Responses struct {
	ResponderName        string `json:"responder_name,omit_empty"`
	ResponderID          string `json:"responder_id,omit_empty"`
	ResponseStatus       string `json:"response_status,omit_empty"`
	ResponderType        string `json:"responder_type,omit_empty"`
	RequestedAt          string `json:"requested_at,omit_empty"`
	RespondedAt          string `json:"responded_at,omit_empty"`
	TimeToRespondSeconds int    `json:"time_to_respond_seconds,omit_empty"`
}

func (c *Client) GetIncidentRawAnalyticsWithContext(ctx context.Context, id string) (*RawAnalytics, error) {
	resp, err := c.get(ctx, "/analytics/raw/incidents/"+id+"/responses", map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
	})
	if err != nil {
		return nil, err
	}

	var result RawAnalytics
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

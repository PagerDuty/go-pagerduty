package pagerduty

import (
	"github.com/google/go-querystring/query"
)

type OnCall struct {
	User             APIObject
	Schedule         APIObject
	EscalationPolicy APIObject
	EscalationLevel  uint
	Start            string
	End              string
}

type ListOnCallsResponse struct {
	OnCalls []OnCall `json:"oncalls"`
}

type ListOnCallOptions struct {
	APIListObject
	TimeZone            string   `url:"time_zone,omitempty"`
	Includes            []string `url:"include,omitempty,brackets"`
	UserIDs             []string `url:"user_ids,omitempty,brackets"`
	EscalationPolicyIDs []string `url:"escalation_policy_ids,omitempty,brackets"`
	ScheduleIDs         []string `url:"schedule_ids,omitempty,brackets"`
	Since               string   `json:"since,omitempty"`
	Until               string   `json:"until,omitempty"`
	Earliest            bool     `json:"earliest,omitempty"`
}

func (c *Client) ListOnCalls(o ListOnCallOptions) (*ListOnCallsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/oncalls?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListOnCallsResponse
	return &result, c.decodeJson(resp, &result)
}

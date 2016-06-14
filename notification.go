package pagerduty

import (
	"github.com/google/go-querystring/query"
)

type Notification struct {
	ID        string `json:"id"`
	Type      string
	StartedAt string `json:"started_at"`
	Address   string
	User      APIObject
}

type ListNotificationOptions struct {
	APIListObject
	TimeZone string   `url:"time_zone,omitempty"`
	Since    string   `url:"since,omitempty"`
	Until    string   `url:"until,omitempty"`
	Filter   string   `url:"filter,omitempty"`
	Includes []string `url:"include,omitempty"`
}

type ListNotificationsResponse struct {
	APIListObject
	Notifications []Notification
}

func (c *Client) ListNotifications(o ListNotificationOptions) (*ListNotificationsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/notifications?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListNotificationsResponse
	return &result, c.decodeJson(resp, &result)
}

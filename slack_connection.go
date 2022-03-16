package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// SlackConnectionConfig is the configuration of a Slack connection as per the documentation
// https://developer.pagerduty.com/api-reference/c2NoOjExMjA5MzMy-slack-connection.
type SlackConnectionConfig struct {
	Events     []string `json:"events"`
	Urgency    *string  `json:"urgency"`
	Priorities []string `json:"priorities"`
}

// SlackConnection is an entity that represents a Slack connections as per the
// documentation https://developer.pagerduty.com/api-reference/c2NoOjExMjA5MzMy-slack-connection.
type SlackConnection struct {
	SourceID   string `json:"source_id"`
	SourceName string `json:"source_name"`
	SourceType string `json:"source_type"`

	ChannelID        string `json:"channel_id"`
	ChannelName      string `json:"channel_name"`
	NotificationType string `json:"notification_type"`

	Config SlackConnectionConfig `json:"config"`
}

// SlackConnectionObject is an API object returned by getter functions.
type SlackConnectionObject struct {
	SlackConnection
	APIObject
}

// ListSlackConnectionsResponse is an API object returned by the list function.
type ListSlackConnectionsResponse struct {
	Connections []SlackConnection `json:"slack_connections"`
	APIListObject
}

// ListSlackConnectionsOptions is the data structure used when calling the ListSlackConnections API endpoint.
type ListSlackConnectionsOptions struct {
	// Limit is the pagination parameter that limits the number of results per
	// page. PagerDuty defaults this value to 50 if omitted, and sets an upper
	// bound of 100.
	Limit uint `url:"limit,omitempty"`

	// Offset is the pagination parameter that specifies the offset at which to
	// start pagination results. When trying to request the next page of
	// results, the new Offset value should be currentOffset + Limit.
	Offset uint `url:"offset,omitempty"`
}

// CreateSlackConnectionWithContext creates a Slack connection.
func (c *Client) CreateSlackConnectionWithContext(ctx context.Context, slackTeamID string, s SlackConnection) (SlackConnectionObject, error) {
	d := map[string]SlackConnection{
		"slack_connection": s,
	}

	resp, err := c.post(ctx, "/workspaces/"+slackTeamID+"/connections", d, nil)
	return getSlackConnectionFromResponse(c, resp, err)
}

// GetSlackConnection gets a Slack connection.
func (c *Client) GetSlackConnectionWithContext(ctx context.Context, slackTeamID, connectionID string) (SlackConnectionObject, error) {
	resp, err := c.get(ctx, "/workspaces/"+slackTeamID+"/connections/"+connectionID)
	return getSlackConnectionFromResponse(c, resp, err)
}

// DeleteSlackConnectionWithContext deletes a Slack connection.
func (c *Client) DeleteSlackConnectionWithContext(ctx context.Context, slackTeamID, connectionID string) error {
	_, err := c.delete(ctx, "/workspaces/"+slackTeamID+"/connections/"+connectionID)
	return err
}

// UpdateSlackConnectionWithContext updates an existing Slack connection.
func (c *Client) UpdateSlackConnectionWithContext(ctx context.Context, slackTeamID, connectionID string, s SlackConnection) (SlackConnectionObject, error) {
	d := map[string]SlackConnection{
		"slack_connection": s,
	}

	resp, err := c.put(ctx, "/workspaces/"+slackTeamID+"/connections/"+connectionID, d, nil)
	return getSlackConnectionFromResponse(c, resp, err)
}

// ListSlackConnectionsWithContext lists Slack connections.
func (c *Client) ListSlackConnectionsWithContext(ctx context.Context, slackTeamID string, o ListSlackConnectionsOptions) (*ListSlackConnectionsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/workspaces/"+slackTeamID+"/connections?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListSlackConnectionsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func getSlackConnectionFromResponse(c *Client, resp *http.Response, err error) (SlackConnectionObject, error) {
	if err != nil {
		return SlackConnectionObject{}, err
	}

	var target map[string]SlackConnectionObject
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return SlackConnectionObject{}, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "slack_connection"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return SlackConnectionObject{}, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return t, nil
}

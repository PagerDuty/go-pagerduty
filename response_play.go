package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type ResponsePlay struct {
	ID                 string          `json:"id,omitempty"`
	Type               string          `json:"type,omitempty"`
	Summary            string          `json:"summary,omitempty"`
	Self               string          `json:"self,omitempty"`
	HTMLURL            string          `json:"html_url,omitempty"`
	Name               string          `json:"name,omitempty"`
	Description        string          `json:"description"`
	Team               *APIReference   `json:"team,omitempty"`
	Subscribers        []*APIReference `json:"subscribers,omitempty"`
	SubscribersMessage string          `json:"subscribers_message"`
	Responders         []*APIReference `json:"responders,omitempty"`
	RespondersMessage  string          `json:"responders_message"`
	Runnability        *string         `json:"runnability"`
	ConferenceNumber   *string         `json:"conference_number"`
	ConferenceUrl      *string         `json:"conference_url"`
	ConferenceType     *string         `json:"conference_type"`
}

type ListResponsePlaysResponse struct {
	ResponsePlays []ResponsePlay `json:"response_plays"`
}

type ListResponsePlaysOptions struct {
	FilterForManualRun bool   `url:"filter_for_manual_run,omitempty"`
	Query              string `url:"query,omitempty"`
}

// ListResponsePlays lists existing response plays.
func (c *Client) ListResponsePlays(ctx context.Context, o ListResponsePlaysOptions) (ListResponsePlaysResponse, error) {
	var result ListResponsePlaysResponse

	v, err := query.Values(o)
	if err != nil {
		return result, err
	}

	resp, err := c.get(ctx, "/response_plays?"+v.Encode())
	if err != nil {
		return result, err
	}

	if err = c.decodeJSON(resp, &result); err != nil {
		return result, err
	}

	return result, nil
}

// CreateResponsePlay creates a new response play.
func (c *Client) CreateResponsePlay(ctx context.Context, rp ResponsePlay) (ResponsePlay, error) {
	d := map[string]ResponsePlay{
		"response_play": rp,
	}

	resp, err := c.post(ctx, "/response_plays", d, nil)
	return getResponsePlayFromResponse(c, resp, err)
}

// GetResponsePlay gets details about an existing response play.
func (c *Client) GetResponsePlay(ctx context.Context, id string) (ResponsePlay, error) {
	resp, err := c.get(ctx, "/response_plays/"+id)
	return getResponsePlayFromResponse(c, resp, err)
}

// UpdateResponsePlay updates an existing response play.
func (c *Client) UpdateResponsePlay(ctx context.Context, rp ResponsePlay) (ResponsePlay, error) {
	d := map[string]ResponsePlay{
		"response_play": rp,
	}

	resp, err := c.put(ctx, "/response_plays/"+rp.ID, d, nil)
	return getResponsePlayFromResponse(c, resp, err)
}

// DeleteResponsePlay deletes an existing response play.
func (c *Client) DeleteResponsePlay(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/response_plays/"+id)
	return err
}

// RunResponsePlay runs a response play on a given incident.
func (c *Client) RunResponsePlay(ctx context.Context, rp string, incident Incident, from string) error {
	d := map[string]APIReference{
		"incident": {
			ID:   incident.ID,
			Type: "incident_reference",
		},
	}

	h := map[string]string{
		"From": from,
	}

	resp, err := c.post(ctx, "/response_plays/"+rp+"/run", d, h)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to run response play %s on incident %s (status code: %d)", rp, incident.ID, resp.StatusCode)
	}

	return nil
}

func getResponsePlayFromResponse(c *Client, resp *http.Response, err error) (ResponsePlay, error) {
	if err != nil {
		return ResponsePlay{}, err
	}

	var target map[string]ResponsePlay
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return ResponsePlay{}, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "response_play"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return ResponsePlay{}, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return t, nil
}

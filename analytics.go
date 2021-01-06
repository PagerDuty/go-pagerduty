package pagerduty

import "io/ioutil"

type Analytics struct {
	AnalyticsFilter *AnalyticsFilter `json:"filters,omitempty"`
	AggregateUnit string `json:"aggregate_unit,omitempty"`
	TimeZone string `json:"time_zone,omitempty"`
}

type AnalyticsFilter struct {
	CreatedAtStart string   `json:"created_at_start,omitempty"`
	CreatedAtEnd   string   `json:"created_at_end,omitempty"`
	Urgency        string   `json:"urgency,omitempty"`
	ServiceIds     []string `json:"service_ids,omitempty"`
	TeamIds     []string `json:"team_ids,omitempty"`
}

func (c *Client) GetAggregatedIncidentData(analytics Analytics) (string,error) {
	headers := make(map[string]string)
	headers["X-EARLY-ACCESS"] = "analytics-v2"
	resp,err:= c.post("/analytics/metrics/incidents/all",analytics,&headers)
	if err != nil {
		return "",err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	bodyString := string(bodyBytes)
	defer resp.Body.Close()
	return bodyString,err
}

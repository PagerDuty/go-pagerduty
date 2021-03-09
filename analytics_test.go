package pagerduty

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAnalytics_GetAggregatedIncidentData(t *testing.T) {
	setup()
	defer teardown()

	analyticsRequest := AnalyticsRequest{
		AnalyticsFilter: &AnalyticsFilter{
			CreatedAtStart: "2021-01-01T15:00:32Z",
			CreatedAtEnd:   "2021-01-08T15:00:32Z",
			TeamIDs:        []string{"PCDYDX0"},
		},
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	analyticsDataWanted := AnalyticsData{MeanSecondsToResolve: 34550, MeanSecondsToFirstAck: 70, MeanEngagedSeconds: 502, MeanAssignmentCount: 1, TotalBusinessHourInterruptions: 1, TotalEngagedSeconds: 2514, TotalIncidentCount: 5, RangeStart: "2021-01-06T00:00:00.000000"}
	analyticsFilterWanted := AnalyticsFilter{CreatedAtStart: "2021-01-06T09:21:41Z", CreatedAtEnd: "2021-01-13T09:21:41Z", TeamIDs: []string{"PCDYDX0"}}
	analyticsResponse := AnalyticsResponse{
		Data:            []AnalyticsData{analyticsDataWanted},
		AnalyticsFilter: &analyticsFilterWanted,
		AggregateUnit:   "day",
		TimeZone:        "Etc/UTC",
	}
	bytesAnalyticsResponse, err := json.Marshal(analyticsResponse)
	mux.HandleFunc("/analytics/metrics/incidents/all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write(bytesAnalyticsResponse)
	})

	client := &Client{apiEndpoint: server.URL,
		authToken:  "foo",
		HTTPClient: defaultHTTPClient,
	}

	res, err := client.GetAggregatedIncidentData(context.Background(), analyticsRequest)
	want := AnalyticsResponse{
		Data:            []AnalyticsData{analyticsDataWanted},
		AnalyticsFilter: &analyticsFilterWanted,
		AggregateUnit:   "day",
		TimeZone:        "Etc/UTC",
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

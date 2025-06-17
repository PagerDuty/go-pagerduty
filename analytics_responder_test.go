package pagerduty

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAnalyticsResponder_GetAggregatedResponderData(t *testing.T) {
	setup()
	defer teardown()

	analyticsResponderRequest := AnalyticsResponderRequest{
		Filters: &AnalyticsResponderFilter{
			DateRangeStart: "2021-01-01T15:00:32Z",
			DateRangeEnd:   "2021-01-08T15:00:32Z",
			TeamIDs:        []string{"PCDYDX0"},
		},
		TimeZone: "Etc/UTC",
		Order:    "desc",
		OrderBy:  "total_incident_count",
	}
	analyticsResponderDataWanted := AnalyticsResponderData{MeanEngagedSeconds: 34550, MeanTimeToAckSeconds: 70, TotalEngagedSeconds: 502, TotalIncidentAck: 1, TotalBusinessHourInterruptions: 1, TotalIncidentCount: 5}
	analyticsResponderFilterWanted := AnalyticsResponderFilter{DateRangeStart: "2021-01-06T09:21:41Z", DateRangeEnd: "2021-01-13T09:21:41Z", TeamIDs: []string{"PCDYDX0"}}
	analyticsResponderResponse := AnalyticsResponderResponse{
		Data:     []AnalyticsResponderData{analyticsResponderDataWanted},
		Filters:  &analyticsResponderFilterWanted,
		TimeZone: "Etc/UTC",
		Order:    "desc",
		OrderBy:  "total_incident_count",
	}

	bytesAnalyticsResponderResponse, err := json.Marshal(analyticsResponderResponse)
	testErrCheck(t, "json.Marshal()", "", err)

	mux.HandleFunc("/analytics/metrics/responders/all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write(bytesAnalyticsResponderResponse)
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetAggregatedResponderData(context.Background(), analyticsResponderRequest)
	want := AnalyticsResponderResponse{
		Data:     []AnalyticsResponderData{analyticsResponderDataWanted},
		Filters:  &analyticsResponderFilterWanted,
		TimeZone: "Etc/UTC",
		Order:    "desc",
		OrderBy:  "total_incident_count",
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

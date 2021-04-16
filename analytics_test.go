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
		Filters: &AnalyticsFilter{
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
		Data:          []AnalyticsData{analyticsDataWanted},
		Filters:       &analyticsFilterWanted,
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}

	bytesAnalyticsResponse, err := json.Marshal(analyticsResponse)
	testErrCheck(t, "json.Marshal()", "", err)

	mux.HandleFunc("/analytics/metrics/incidents/all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write(bytesAnalyticsResponse)
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetAggregatedIncidentData(context.Background(), analyticsRequest)
	want := AnalyticsResponse{
		Data:          []AnalyticsData{analyticsDataWanted},
		Filters:       &analyticsFilterWanted,
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestAnalytics_GetAggregatedServiceData(t *testing.T) {
	setup()
	defer teardown()

	analyticsRequest := AnalyticsRequest{
		Filters: &AnalyticsFilter{
			CreatedAtStart: "2021-01-01T15:00:32Z",
			CreatedAtEnd:   "2021-01-08T15:00:32Z",
			TeamIDs:        []string{"PCDYDX0"},
		},
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	analyticsDataWanted := AnalyticsData{MeanAssignmentCount: 1, MeanEngagedSeconds: 502, MeanEngagedUserCount: 0, MeanSecondsToResolve: 34550, MeanSecondsToFirstAck: 70, TotalBusinessHourInterruptions: 1, TotalEngagedSeconds: 2514, TotalIncidentCount: 5, RangeStart: "2021-01-06T00:00:00.000000", ServiceID: "PSEJLIN", ServiceName: "FooAlerts", TeamID: "PCDYDX0", TeamName: "FooTeam", UpTimePct: 89.86111111111111}
	analyticsFilterWanted := AnalyticsFilter{CreatedAtStart: "2021-01-06T09:21:41Z", CreatedAtEnd: "2021-01-13T09:21:41Z", TeamIDs: []string{"PCDYDX0"}}
	analyticsResponse := AnalyticsResponse{
		Data:          []AnalyticsData{analyticsDataWanted},
		Filters:       &analyticsFilterWanted,
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	bytesAnalyticsResponse, err := json.Marshal(analyticsResponse)
	testErrCheck(t, "json.Marshal()", "", err)

	mux.HandleFunc("/analytics/metrics/incidents/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write(bytesAnalyticsResponse)
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetAggregatedServiceData(context.Background(), analyticsRequest)
	want := AnalyticsResponse{
		Data:          []AnalyticsData{analyticsDataWanted},
		Filters:       &analyticsFilterWanted,
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestAnalytics_GetAggregatedTeamData(t *testing.T) {
	setup()
	defer teardown()

	analyticsRequest := AnalyticsRequest{
		Filters: &AnalyticsFilter{
			CreatedAtStart: "2021-01-01T15:00:32Z",
			CreatedAtEnd:   "2021-01-08T15:00:32Z",
			TeamIDs:        []string{"PCDYDX0"},
		},
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	analyticsDataWanted := AnalyticsData{MeanAssignmentCount: 1, MeanEngagedSeconds: 502, MeanEngagedUserCount: 0, MeanSecondsToResolve: 34550, MeanSecondsToFirstAck: 70, TotalBusinessHourInterruptions: 1, TotalEngagedSeconds: 2514, TotalIncidentCount: 5, RangeStart: "2021-01-06T00:00:00.000000", TeamID: "PCDYDX0", TeamName: "FooTeam", UpTimePct: 89.86111111111111}
	analyticsFilterWanted := AnalyticsFilter{CreatedAtStart: "2021-01-06T09:21:41Z", CreatedAtEnd: "2021-01-13T09:21:41Z", TeamIDs: []string{"PCDYDX0"}}
	analyticsResponse := AnalyticsResponse{
		Data:          []AnalyticsData{analyticsDataWanted},
		Filters:       &analyticsFilterWanted,
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	bytesAnalyticsResponse, err := json.Marshal(analyticsResponse)
	testErrCheck(t, "json.Marshal()", "", err)

	mux.HandleFunc("/analytics/metrics/incidents/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write(bytesAnalyticsResponse)
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetAggregatedTeamData(context.Background(), analyticsRequest)
	want := AnalyticsResponse{
		Data:          []AnalyticsData{analyticsDataWanted},
		Filters:       &analyticsFilterWanted,
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

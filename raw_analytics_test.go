package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

func TestRawAnalytics_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/analytics/raw/incidents/1/responses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"incident_id": "1","limit": 100,"order": "asc","order_by": "requested_at","time_zone": "Etc/UTC","responses": {  "responder_name": "Earline Greenholt",  "responder_id": "PXPGF42",  "response_status": "accepted",  "responder_type": "added_responder",  "requested_at": "2023-01-05T10:15:00",  "responded_at": "2023-01-05T10:18:00",  "time_to_respond_seconds": 180}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	id := "1"
	res, err := client.GetIncidentRawAnalyticsWithContext(context.Background(), id)

	want := &RawAnalytics{
		IncidentID: "1",
		Limit: 100,
		Order: "asc",
		OrderBy: "requested_at",
		TimeZone: "Etc/UTC",
		Responses: &Responses{
			ResponderName: "Earline Greenholt",
			ResponderID: "PXPGF42",
			ResponseStatus: "accepted",
			ResponderType: "added_responder",
			RequestedAt: "2023-01-05T10:15:00",
			RespondedAt: "2023-01-05T10:18:00",
			TimeToRespondSeconds: 180,
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

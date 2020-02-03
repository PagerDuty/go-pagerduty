package pagerduty

import (
	"net/http"
	"testing"
)

// ListOnCalls
func TestOnCall_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/oncalls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"oncalls": [{"escalation_level":2}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListOnCallOptions{
		APIListObject:       listObj,
		TimeZone:            "UTC",
		Includes:            []string{},
		UserIDs:             []string{},
		EscalationPolicyIDs: []string{},
		ScheduleIDs:         []string{},
		Earliest:            false,
		Since:               "bar",
		Until:               "baz",
	}
	res, err := client.ListOnCalls(opts)

	want := &ListOnCallsResponse{
		APIListObject: listObj,
		OnCalls: []OnCall{
			{
				EscalationLevel: 2,
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

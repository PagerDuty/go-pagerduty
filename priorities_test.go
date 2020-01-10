package pagerduty

import (
	"net/http"
	"testing"
)

// ListMaintenanceWindows
func TestPriorities_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/priorities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"priorities": [{"id": "1", "summary": "foo"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	res, err := client.ListPriorities()

	want := &Priorities{
		APIListObject: listObj,
		Priorities: []PriorityProperty{
			{
				APIObject: APIObject{
					ID:      "1",
					Summary: "foo",
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

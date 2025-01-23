package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// ListMaintenanceWindows
func TestPriorities_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/priorities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"priorities": [{"id": "1", "summary": "foo"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListPrioritiesWithContext(context.Background(), ListPrioritiesOptions{})

	want := &ListPrioritiesResponse{
		APIListObject: listObj,
		Priorities: []Priority{
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

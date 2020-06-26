package pagerduty

import (
	"net/http"
	"testing"
)

func TestLogEntry_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/log_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"log_entries": [{"id": "1","summary":"foo"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var entriesOpts = ListLogEntriesOptions{
		APIListObject: listObj,
		Includes:      []string{},
		IsOverview:    true,
		TimeZone:      "UTC",
	}
	res, err := client.ListLogEntries(entriesOpts)

	want := &ListLogEntryResponse{
		APIListObject: listObj,
		LogEntries: []LogEntry{
			{
				CommonLogEntryField: CommonLogEntryField{
					APIObject: APIObject{
						ID:      "1",
						Summary: "foo",
					},
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestLogEntry_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/log_entries/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"log_entry": {"id": "1", "summary": "foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	opts := GetLogEntryOptions{TimeZone: "UTC", Includes: []string{}}
	res, err := client.GetLogEntry(id, opts)

	want := &LogEntry{
		CommonLogEntryField: CommonLogEntryField{
			APIObject: APIObject{
				ID:      "1",
				Summary: "foo",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

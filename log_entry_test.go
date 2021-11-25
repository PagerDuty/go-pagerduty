package pagerduty

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestLogEntry_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/log_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"log_entries": [{"id": "1","summary":"foo"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	entriesOpts := ListLogEntriesOptions{
		Limit:      listObj.Limit,
		Offset:     listObj.Offset,
		Includes:   []string{},
		IsOverview: true,
		TimeZone:   "UTC",
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
		_, _ = w.Write([]byte(`{"log_entry": {"id": "1", "summary": "foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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

func TestChannel_MarhalUnmarshal(t *testing.T) {
	logEntryRaw := []byte(`{
		"id": "1",
		"type": "trigger_log_entry",
		"summary": "foo",
		"channel": {
			"type": "web_trigger",
			"summary": "My new incident",
			"details_omitted": false
		},
		"service": {
			"type": "service_reference",
			"id": "abc123",
			"summary": "xyz890",
			"self": "http://test.local/api/service/abc123",
			"html_url": "http://test.local/service/abc123"
		},
		"user": {
			"type": "user_reference",
			"id": "abc890",
			"summary": "xyz123",
			"self": "http://test.local/api/user/abc890",
			"html_url": "http://test.local/user/abc890"
		}
	}`)
	want := &LogEntry{
		CommonLogEntryField: CommonLogEntryField{
			APIObject: APIObject{
				ID:      "1",
				Type:    "trigger_log_entry",
				Summary: "foo",
			},
			Channel: Channel{
				Type: "web_trigger",
				Raw: map[string]interface{}{
					"type":            "web_trigger",
					"summary":         "My new incident",
					"details_omitted": false,
				},
			},
		},
		Service: APIObject{
			ID:      "abc123",
			Type:    "service_reference",
			Summary: "xyz890",
			Self:    "http://test.local/api/service/abc123",
			HTMLURL: "http://test.local/service/abc123",
		},
		User: APIObject{
			ID:      "abc890",
			Type:    "user_reference",
			Summary: "xyz123",
			Self:    "http://test.local/api/user/abc890",
			HTMLURL: "http://test.local/user/abc890",
		},
	}

	logEntry := &LogEntry{}
	if err := json.Unmarshal(logEntryRaw, logEntry); err != nil {
		t.Fatal(err)
	}

	testEqual(t, want, logEntry)

	newLogEntryRaw, err := json.Marshal(logEntry)
	if err != nil {
		t.Fatal(err)
	}

	newLogEntry := &LogEntry{}
	if err := json.Unmarshal(newLogEntryRaw, newLogEntry); err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, newLogEntry)
}

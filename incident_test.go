package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestIncident_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"incidents": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var opts ListIncidentsOptions
	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListIncidents(opts)

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				APIObject: APIObject{
					ID: "1",
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &CreateIncidentOptions{
		Title:   "foo",
		Urgency: "low",
	}

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		body, err := ioutil.ReadAll(r.Body)
		testErrCheck(t, "ioutil.ReadAll", "", err)

		fmt.Println(string(body))

		got := make(map[string]CreateIncidentOptions)
		testErrCheck(t, "json.Unmarshal()", "", json.Unmarshal(body, &got))

		o, ok := got["incident"]
		if !ok {
			t.Fatal("map does not have an incident key")
		}

		if o.Type != "incident" {
			t.Errorf("o.Type = %q, want %q", o.Type, "incident")
		}

		if o.Title != "foo" {
			t.Errorf("o.Foo = %q, want %q", o.Title, "foo")
		}

		if o.Urgency != "low" {
			t.Errorf("o.Urgency = %q, want %q", o.Urgency, "low")
		}

		_, _ = w.Write([]byte(`{"incident": {"title": "foo", "id": "1", "urgency": "low"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	from := "foo@bar.com"
	res, err := client.CreateIncident(from, input)

	want := &Incident{
		APIObject: APIObject{
			ID: "1",
		},
		Title:   "foo",
		Urgency: "low",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Manage_status(t *testing.T) {
	setup()
	defer teardown()

	wantFrom := "foo@bar.com"

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		if gotFrom := r.Header.Get("From"); gotFrom != wantFrom {
			t.Errorf("From HTTP header = %q, want %q", gotFrom, wantFrom)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		var data map[string][]ManageIncidentsOptions
		testErrCheck(t, "json.Unmarshal()", "", json.Unmarshal(body, &data))

		if len(data["incidents"]) == 0 {
			t.Fatalf("no incidents, expect 1")
		}

		const (
			wantType   = "incident"
			wantID     = "1"
			wantStatus = "acknowledged"
		)

		inc := data["incidents"][0]

		if inc.Type != wantType {
			t.Errorf("inc.Type = %q, want %q", inc.Type, wantType)
		}

		if inc.ID != wantID {
			t.Errorf("inc.ID = %q, want %q", inc.ID, wantID)
		}

		if inc.Status != wantStatus {
			t.Errorf("inc.Status = %q, want %q", inc.Status, wantStatus)
		}

		_, _ = w.Write([]byte(`{"incidents": [{"title": "foo", "id": "1", "status": "acknowledged"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")

	input := []ManageIncidentsOptions{
		{
			ID:     "1",
			Status: "acknowledged",
		},
	}

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				APIObject: APIObject{
					ID: "1",
				},
				Title:  "foo",
				Status: "acknowledged",
			},
		},
	}
	res, err := client.ManageIncidents(wantFrom, input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Manage_priority(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"incidents": [{"title": "foo", "id": "1", "priority": {"id": "PRIORITY_ID_HERE", "type": "priority_reference"}}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	from := "foo@bar.com"

	input := []ManageIncidentsOptions{
		{
			ID:   "1",
			Type: "incident",
			Priority: &APIReference{
				ID:   "PRIORITY_ID_HERE",
				Type: "priority_reference",
			},
		},
	}

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				APIObject: APIObject{
					ID: "1",
				},
				Title: "foo",
				Priority: &Priority{
					APIObject: APIObject{
						ID:   "PRIORITY_ID_HERE",
						Type: "priority_reference",
					},
				},
			},
		},
	}
	res, err := client.ManageIncidents(from, input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Manage_title(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		// test if title is present and correct
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var data map[string][]ManageIncidentsOptions
		if err := json.Unmarshal(body, &data); err != nil {
			t.Fatal(err)
		}
		if len(data["incidents"]) == 0 {
			t.Fatalf("no incidents, expect 1")
		}
		if data["incidents"][0].Title != "bar" {
			t.Fatalf("expected incidents[0].title to be \"bar\" got \"%s\"", data["incidents"][0].Title)
		}
		_, _ = w.Write([]byte(`{"incidents": [{"title": "bar", "id": "1"}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	from := "foo@bar.com"

	input := []ManageIncidentsOptions{
		{
			ID:    "1",
			Type:  "incident",
			Title: "bar",
		},
	}

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				APIObject: APIObject{
					ID: "1",
				},
				Title: "bar",
			},
		},
	}
	res, err := client.ManageIncidents(from, input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Manage_assignments(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"incidents": [{"title": "foo", "id": "1", "assignments": [{"assignee":{"id": "ASSIGNEE_ONE", "type": "user_reference"}},{"assignee":{"id": "ASSIGNEE_TWO", "type": "user_reference"}}]}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	from := "foo@bar.com"

	input := []ManageIncidentsOptions{
		{
			ID:   "1",
			Type: "incident",
			Assignments: []Assignee{
				{
					Assignee: APIObject{
						ID:   "ASSIGNEE_ONE",
						Type: "user_reference",
					},
				},
				{
					Assignee: APIObject{
						ID:   "ASSIGNEE_TWO",
						Type: "user_reference",
					},
				},
			},
		},
	}

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				APIObject: APIObject{
					ID: "1",
				},
				Title: "foo",
				Assignments: []Assignment{
					{
						Assignee: APIObject{
							ID:   "ASSIGNEE_ONE",
							Type: "user_reference",
						},
					},
					{
						Assignee: APIObject{
							ID:   "ASSIGNEE_TWO",
							Type: "user_reference",
						},
					},
				},
			},
		},
	}
	res, err := client.ManageIncidents(from, input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Manage_conference_bridge(t *testing.T) {
	setup()
	defer teardown()

	wantFrom := "foo@bar.com"

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		if gotFrom := r.Header.Get("From"); gotFrom != wantFrom {
			t.Errorf("From HTTP header = %q, want %q", gotFrom, wantFrom)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		var data map[string][]ManageIncidentsOptions
		testErrCheck(t, "json.Unmarshal()", "", json.Unmarshal(body, &data))

		if len(data["incidents"]) == 0 {
			t.Fatalf("no incidents, expect 1")
		}

		const (
			wantNum = "42"
			wantURL = "http://example.org/bridge"
		)

		inc := data["incidents"][0]

		if inc.ConferenceBridge == nil {
			t.Fatalf("inc.ConferenceBridge = <nil>")
		}

		if got := inc.ConferenceBridge.ConferenceNumber; got != wantNum {
			t.Fatalf("inc.ConferenceBridge.ConferenceNumber = %q, want %q", got, wantNum)
		}

		if got := inc.ConferenceBridge.ConferenceURL; got != wantURL {
			t.Fatalf("inc.ConferenceBridge.ConferenceNumber = %q, want %q", got, wantURL)
		}

		_, _ = w.Write([]byte(`{"incidents": [{"title": "foo", "id": "1","conference_bridge":{"conference_number":"42","conference_url":"http://example.org/bridge"}}]}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := []ManageIncidentsOptions{
		{
			ID:   "1",
			Type: "incident",
			ConferenceBridge: &ConferenceBridge{
				ConferenceNumber: "42",
				ConferenceURL:    "http://example.org/bridge",
			},
		},
	}

	want := &ListIncidentsResponse{
		APIListObject: APIListObject{Limit: 0, Offset: 0, More: false, Total: 0},
		Incidents: []Incident{
			{
				APIObject: APIObject{
					ID: "1",
				},
				Title: "foo",
				ConferenceBridge: &ConferenceBridge{
					ConferenceNumber: "42",
					ConferenceURL:    "http://example.org/bridge",
				},
			},
		},
	}

	got, err := client.ManageIncidents(wantFrom, input)
	if err != nil {
		t.Fatal(err)
	}

	testEqual(t, want, got)
}

func TestIncident_Manage_esclation_level(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"incidents": [{"title": "foo", "id": "1"}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	from := "foo@bar.com"

	input := []ManageIncidentsOptions{
		{
			ID:              "1",
			Type:            "incident",
			EscalationLevel: 2,
		},
	}

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				APIObject: APIObject{
					ID: "1",
				},
				Title: "foo",
			},
		},
	}
	res, err := client.ManageIncidents(from, input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Merge(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/merge", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"incident": {"title": "foo", "id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	from := "foo@bar.com"

	input := []MergeIncidentsOptions{{ID: "2", Type: "incident"}}
	want := &Incident{
		APIObject: APIObject{
			ID: "1",
		},
		Title: "foo",
	}

	res, err := client.MergeIncidents(from, "1", input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"incident": {"id": "1", "incidents_responders": [{"incident": {"id": "1"}}]}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	id := "1"
	res, err := client.GetIncident(id)

	want := &Incident{APIObject: APIObject{ID: "1"}, IncidentResponders: []IncidentResponders{{Incident: APIObject{ID: "1"}}}}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_ListIncidentNotes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/notes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"notes": [{"id": "1","content":"foo"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"

	res, err := client.ListIncidentNotes(id)

	want := []IncidentNote{
		{
			ID:      "1",
			Content: "foo",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_ListIncidentAlerts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"alerts": [{"id": "1","summary":"foo"}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	id := "1"

	res, err := client.ListIncidentAlerts(id)

	want := &ListAlertsResponse{
		APIListObject: listObj,
		Alerts: []IncidentAlert{
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

func TestIncident_ListIncidentAlertsWithOpts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"alerts": [{"id": "1","summary":"foo"}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	id := "1"

	alertOpts := ListIncidentAlertsOptions{
		Limit:    listObj.Limit,
		Offset:   listObj.Offset,
		Includes: []string{},
	}

	res, err := client.ListIncidentAlertsWithOpts(id, alertOpts)

	want := &ListAlertsResponse{
		APIListObject: listObj,
		Alerts: []IncidentAlert{
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

// CreateIncidentNote
func TestIncident_CreateIncidentNote(t *testing.T) {
	setup()
	defer teardown()

	input := IncidentNote{
		Content: "foo",
	}

	mux.HandleFunc("/incidents/1/notes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"note": {"id": "1","content": "foo"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	id := "1"
	err := client.CreateIncidentNote(id, input)
	if err != nil {
		t.Fatal(err)
	}
}

// CreateIncidentNoteWithResponse
func TestIncident_CreateIncidentNoteWithResponse(t *testing.T) {
	setup()
	defer teardown()

	input := IncidentNote{
		Content: "foo",
	}

	mux.HandleFunc("/incidents/1/notes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"note": {"id": "1","content": "foo"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	id := "1"
	res, err := client.CreateIncidentNoteWithResponse(id, input)

	want := &IncidentNote{
		ID:      "1",
		Content: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// SnoozeIncident
func TestIncident_SnoozeIncident(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/snooze", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"incident": {"id": "1", "pending_actions": [{"type": "unacknowledge", "at":"2019-12-31T16:58:35Z"}]}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	var duration uint = 3600
	id := "1"

	err := client.SnoozeIncident(id, duration)
	if err != nil {
		t.Fatal(err)
	}
}

// SnoozeIncidentWithResponse
func TestIncident_SnoozeIncidentWithResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/snooze", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"incident": {"id": "1", "pending_actions": [{"type": "unacknowledge", "at":"2019-12-31T16:58:35Z"}]}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	var duration uint = 3600
	id := "1"

	res, err := client.SnoozeIncidentWithResponse(id, duration)

	want := &Incident{
		APIObject: APIObject{
			ID: "1",
		},
		PendingActions: []PendingAction{
			{
				Type: "unacknowledge",
				At:   "2019-12-31T16:58:35Z",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// ListIncidentLogEntries
func TestIncident_ListLogEntries(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/log_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"log_entries": [{"id": "1","summary":"foo"}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	id := "1"
	entriesOpts := ListIncidentLogEntriesOptions{
		Limit:      listObj.Limit,
		Offset:     listObj.Offset,
		Includes:   []string{},
		IsOverview: true,
		TimeZone:   "UTC",
	}
	res, err := client.ListIncidentLogEntries(id, entriesOpts)

	want := &ListIncidentLogEntriesResponse{
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

// ListIncidentLogEntriesSinceUntil
func TestIncident_ListLogEntriesSinceUntil(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/log_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"log_entries": [{"id": "1","summary":"foo"}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	id := "1"
	entriesOpts := ListIncidentLogEntriesOptions{
		Limit:      listObj.Limit,
		Offset:     listObj.Offset,
		Includes:   []string{},
		IsOverview: true,
		TimeZone:   "UTC",
		Since:      "2020-03-27T22:40:00-0700",
		Until:      "2020-03-28T22:50:00-0700",
	}
	res, err := client.ListIncidentLogEntries(id, entriesOpts)

	want := &ListIncidentLogEntriesResponse{
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

func TestIncident_ResponderRequest(t *testing.T) {
	setup()
	defer teardown()

	id := "1"
	mux.HandleFunc("/incidents/"+id+"/responder_requests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{
	"responder_request": {
		"requester": {
			"id": "PL1JMK5",
			"type": "user_reference"
		},
		"message": "Help",
		"responder_request_targets": [{
			"responder_request_target": {
				"id": "PJ25ZYX",
				"type": "user_reference",
				"incidents_responders": [
					{
						"state": "pending",
						"user": {
							"id": "PJ25ZYX",
							"type": "user_reference",
							"summary": "dave"
						}
					}
				]
			}
		}]
	}
}`))
	})
	client := defaultTestClient(server.URL, "foo")
	from := "foo@bar.com"

	request_target := ResponderRequestTarget{}
	request_target.ID = "PJ25ZYX"
	request_target.Type = "user_reference"

	request_target_wrapper := ResponderRequestTargetWrapper{Target: request_target}
	request_targets := []ResponderRequestTargetWrapper{request_target_wrapper}

	input := ResponderRequestOptions{
		From:        from,
		Message:     "Help",
		RequesterID: "PL1JMK5",
		Targets:     request_targets,
	}

	user := User{}
	user.ID = "PL1JMK5"
	user.Type = "user_reference"

	want_target := ResponderRequestTarget{}
	want_target.ID = "PJ25ZYX"
	want_target.Type = "user_reference"

	want_target.Responders = []IncidentResponders{
		{
			State: "pending",
			User: APIObject{
				ID:      "PJ25ZYX",
				Type:    "user_reference",
				Summary: "dave",
			},
		},
	}

	want_target_wrapper := ResponderRequestTargetWrapper{Target: want_target}
	want_targets := []ResponderRequestTargetWrapper{want_target_wrapper}

	want := &ResponderRequestResponse{
		ResponderRequest: ResponderRequest{
			Incident:  Incident{},
			Requester: user,
			Message:   "Help",
			Targets:   want_targets,
		},
	}
	res, err := client.ResponderRequest(id, input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_GetAlert(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/alerts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"alert": {"id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	incidentID := "1"
	alertID := "1"
	res, err := client.GetIncidentAlert(incidentID, alertID)

	want := &IncidentAlertResponse{
		IncidentAlert: &IncidentAlert{
			APIObject: APIObject{
				ID: "1",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_ManageIncidentAlerts(t *testing.T) {
	setup()
	defer teardown()

	from := "pagerduty@example.com"

	mux.HandleFunc("/incidents/1/alerts/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		if hdr := r.Header.Get("From"); hdr != "pagerduty@example.com" {
			t.Errorf("From header = %q, want %q", hdr, from)
		}
		_, _ = w.Write([]byte(`{"alerts": [{"id": "1"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")

	incidentID := "1"

	input := &IncidentAlertList{
		Alerts: []IncidentAlert{
			{
				APIObject: APIObject{
					ID: "1",
				},
			},
		},
	}
	res, err := client.ManageIncidentAlerts(context.Background(), incidentID, from, input)

	want := &ListAlertsResponse{
		Alerts: []IncidentAlert{
			{
				APIObject: APIObject{
					ID: "1",
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// CreateIncidentStatusUpdate
func TestIncident_CreateIncidentStatusUpdate(t *testing.T) {
	setup()
	defer teardown()

	wantID := "1"
	wantFrom := "foo@bar.com"
	wantMessage := "foo"

	mux.HandleFunc("/incidents/1/status_updates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		if gotFrom := r.Header.Get("From"); gotFrom != wantFrom {
			t.Errorf("From HTTP header = %q, want %q", gotFrom, wantFrom)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		var data map[string]string
		if err := json.Unmarshal(body, &data); err != nil {
			t.Fatal(err)
		}
		o, ok := data["message"]
		if !ok {
			t.Fatal("map does not have message key")
		}
		if o != wantMessage {
			t.Errorf("message = %q, want %q", o, wantMessage)
		}

		_, _ = w.Write([]byte(`{"status_update": {"id": "1", "message": "foo", "sender": {"summary": "foo@bar.com", "type": "user_reference"}}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	res, err := client.CreateIncidentStatusUpdate(context.Background(), wantID, wantFrom, wantMessage)
	want := IncidentStatusUpdate{
		ID:      wantID,
		Message: wantMessage,
		Sender: APIObject{
			Summary: wantFrom,
			Type:    "user_reference",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_ListIncidentNotificationSubscribersWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/1/status_updates/subscribers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"subscribers": [ { "subscriber_id": "PD1234", "subscriber_type": "user", "has_indirect_subscription": false, "subscribed_via": null }, { "subscriber_id": "PD1234", "subscriber_type": "team", "has_indirect_subscription": true, "subscribed_via": [ { "id": "PD1234", "type": "business_service" } ] } ], "account_id": "PD1234"}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}

	client := defaultTestClient(server.URL, "foo")
	id := "1"

	res, err := client.ListIncidentNotificationSubscribersWithContext(context.Background(), id)

	want := &ListIncidentNotificationSubscribersResponse{
		APIListObject: listObj,
		Subscribers: []IncidentNotificationSubscriptionWithContext{
			{
				IncidentNotificationSubscriber: IncidentNotificationSubscriber{
					SubscriberID:   "PD1234",
					SubscriberType: "user",
				},
				HasIndirectSubscription: false,
				SubscribedVia:           nil,
			},
			{
				IncidentNotificationSubscriber: IncidentNotificationSubscriber{
					SubscriberID:   "PD1234",
					SubscriberType: "team",
				},
				HasIndirectSubscription: true,
				SubscribedVia: []IncidentNotificationSubscriberVia{
					{
						ID:   "PD1234",
						Type: "business_service",
					},
				},
			},
		},
		AccountID: "PD1234",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_AddIncidentNotificationSubscribersWithContext(t *testing.T) {
	setup()
	defer teardown()

	input := []IncidentNotificationSubscriber{
		{
			SubscriberID:   "PD1234",
			SubscriberType: "team",
		},
	}

	mux.HandleFunc("/incidents/1/status_updates/subscribers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{ "subscriptions": [ { "account_id": "PD1234", "subscribable_id": "PD1234", "subscribable_type": "incident", "subscriber_id": "PD1234", "subscriber_type": "team", "result": "success" } ] }`))
	})
	client := defaultTestClient(server.URL, "foo")
	id := "1"
	res, err := client.AddIncidentNotificationSubscribersWithContext(context.Background(), id, input)
	if err != nil {
		t.Fatal(err)
	}

	want := &AddIncidentNotificationSubscribersResponse{
		Subscriptions: []IncidentNotificationSubscriptionWithContext{
			{
				IncidentNotificationSubscriber: IncidentNotificationSubscriber{
					SubscriberID:   "PD1234",
					SubscriberType: "team",
				},
				SubscribableID:   "PD1234",
				SubscribableType: "incident",
				Result:           "success",
				AccountID:        "PD1234",
			},
		},
	}
	testEqual(t, want, res)
}

func TestIncident_RemoveIncidentNotificationSubscribersWithContext(t *testing.T) {
	setup()
	defer teardown()

	input := []IncidentNotificationSubscriber{
		{
			SubscriberID:   "PD1234",
			SubscriberType: "team",
		},
	}

	mux.HandleFunc("/incidents/1/status_updates/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"deleted_count": 1, "unauthorized_count": 0, "non_existent_count": 0}`))
	})
	client := defaultTestClient(server.URL, "foo")
	id := "1"
	res, err := client.RemoveIncidentNotificationSubscribersWithContext(context.Background(), id, input)
	if err != nil {
		t.Fatal(err)
	}

	want := &RemoveIncidentNotificationSubscribersResponse{
		DeleteCount:       1,
		UnauthorizedCount: 0,
		NonExistentCount:  0,
	}
	testEqual(t, want, res)
}

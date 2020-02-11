package pagerduty

import (
	"net/http"
	"testing"
)

func TestIncident_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"incidents": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var opts ListIncidentsOptions
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	res, err := client.ListIncidents(opts)

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				Id: "1",
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
		Type:    "incident",
		Title:   "foo",
		Urgency: "low",
	}

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"incident": {"title": "foo", "id": "1", "urgency": "low"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	from := "foo@bar.com"
	res, err := client.CreateIncident(from, input)

	want := &Incident{
		Title:   "foo",
		Id:      "1",
		Urgency: "low",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncident_Manage(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"incidents": [{"title": "foo", "id": "1", "status": "acknowledged"}]}`))
	})
	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	from := "foo@bar.com"

	input := []ManageIncidentsOptions{
		{
			ID:     "1",
			Type:   "incident",
			Status: "acknowledged",
		},
	}

	want := &ListIncidentsResponse{
		APIListObject: listObj,
		Incidents: []Incident{
			{
				Id:     "1",
				Title:  "foo",
				Status: "acknowledged",
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
		w.Write([]byte(`{"incident": {"title": "foo", "id": "1"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	from := "foo@bar.com"

	input := []MergeIncidentsOptions{{ID: "2", Type: "incident"}}
	want := &Incident{Id: "1", Title: "foo"}

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
		w.Write([]byte(`{"incident": {"id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	id := "1"
	res, err := client.GetIncident(id)

	want := &Incident{Id: "1"}

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
		w.Write([]byte(`{"notes": [{"id": "1","content":"foo"}]}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
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
		w.Write([]byte(`{"alerts": [{"id": "1","summary":"foo"}]}`))
	})
	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
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

// CreateIncidentNote
func TestIncident_CreateIncidentNote(t *testing.T) {
	setup()
	defer teardown()

	input := IncidentNote{
		Content: "foo",
	}

	mux.HandleFunc("/incidents/1/notes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"note": {"id": "1","content": "foo"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
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
		w.Write([]byte(`{"note": {"id": "1","content": "foo"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
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
		w.Write([]byte(`{"incident": {"id": "1", "pending_actions": [{"type": "unacknowledge", "at":"2019-12-31T16:58:35Z"}]}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
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
		w.Write([]byte(`{"incident": {"id": "1", "pending_actions": [{"type": "unacknowledge", "at":"2019-12-31T16:58:35Z"}]}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var duration uint = 3600
	id := "1"

	res, err := client.SnoozeIncidentWithResponse(id, duration)

	want := &Incident{
		Id: "1",
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
		w.Write([]byte(`{"log_entries": [{"id": "1","summary":"foo"}]}`))
	})
	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	var entriesOpts = ListIncidentLogEntriesOptions{
		APIListObject: listObj,
		Includes:      []string{},
		IsOverview:    true,
		TimeZone:      "UTC",
	}
	res, err := client.ListIncidentLogEntries(id, entriesOpts)

	want := &ListIncidentLogEntriesResponse{
		APIListObject: listObj,
		LogEntries: []LogEntry{
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

func TestIncident_ResponderRequest(t *testing.T) {
	setup()
	defer teardown()

	id := "1"
	mux.HandleFunc("/incidents/"+id+"/responder_requests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{
	"responder_request": {
		"requester": {
			"id": "PL1JMK5",
			"type": "user_reference"
		},
		"message": "Help",
		"responder_request_targets": {
			"responder_request_target": {
				"id": "PJ25ZYX",
				"type": "user_reference",
				"incident_responders": {
					"state": "pending",
					"user": {
						"id": "PJ25ZYX"
					}
				}
			}
		}
	}
}`))

	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	from := "foo@bar.com"

	r := ResponderRequestTarget{}
	r.ID = "PJ25ZYX"
	r.Type = "user_reference"

	input := ResponderRequestOptions{
		From:        from,
		Message:     "help",
		RequesterID: "PL1JMK5",
		Targets:     []ResponderRequestTarget{r},
	}

	user := User{}
	user.ID = "PL1JMK5"
	user.Type = "user_reference"

	target := ResponderRequestTarget{}
	target.ID = "PJ25ZYX"
	target.Type = "user_reference"
	target.Responders.State = "pending"
	target.Responders.User.ID = "PJ25ZYX"

	want := &ResponderRequestResponse{
		ResponderRequest: ResponderRequest{
			Incident:  Incident{},
			Requester: user,
			Message:   "Help",
			Targets:   ResponderRequestTargets{target},
		},
	}
	res, err := client.ResponderRequest(id, input)

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

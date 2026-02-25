package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

var (
	testScheduleV3ID = "SCHED01"
	testRotationV3ID = "ROT01"
	testEventV3ID    = "EVT01"
)

const (
	mockScheduleV3ListResponse = `{
		"schedules": [
			{
				"id": "SCHED01",
				"type": "schedule_reference",
				"summary": "On-Call Schedule",
				"self": "https://api.pagerduty.com/v3/schedules/SCHED01",
				"html_url": "https://example.pagerduty.com/schedules/SCHED01"
			}
		]
	}`

	mockScheduleV3ListEmptyResponse = `{"schedules": []}`

	// No "events" key in rotation — Go decodes Events as nil.
	mockScheduleV3GetResponse = `{
		"schedule": {
			"id": "SCHED01",
			"type": "schedule",
			"name": "On-Call Schedule",
			"time_zone": "UTC",
			"description": "Test schedule",
			"rotations": [
				{"id": "ROT01", "type": "rotation"}
			]
		}
	}`

	// No rotations — used for create / update responses.
	mockScheduleV3MutateResponse = `{
		"schedule": {
			"id": "SCHED01",
			"type": "schedule",
			"name": "On-Call Schedule",
			"time_zone": "UTC",
			"description": "Test schedule"
		}
	}`

	// "events": [] — Go decodes Events as []EventV3{} (non-nil empty slice).
	mockRotationV3Response = `{
		"rotation": {
			"id": "ROT01",
			"type": "rotation",
			"events": []
		}
	}`

	mockEventV3Response = `{
		"schedule_event": {
			"id": "EVT01",
			"type": "schedule_event",
			"name": "On-Call Event",
			"start_time": {"date_time": "2026-02-21T09:00:00Z"},
			"end_time":   {"date_time": "2026-02-21T17:00:00Z"},
			"effective_since": "2026-02-21T09:00:00Z",
			"effective_until": null,
			"recurrence": ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"],
			"assignment_strategy": {
				"type": "user_assignment_strategy",
				"members": [
					{"type": "user_member", "user_id": "USER01"}
				]
			}
		}
	}`

	mockScheduleV3Error400 = `{
		"error": {
			"code": 2001,
			"message": "Invalid Input Provided",
			"errors": ["name is required"]
		}
	}`

	mockScheduleV3Error404 = `{
		"error": {
			"code": 2100,
			"message": "Not Found",
			"errors": ["The specified resource does not exist"]
		}
	}`

	mockScheduleV3Error500 = `{
		"error": {
			"code": 3001,
			"message": "Internal Server Error",
			"errors": ["An unexpected error occurred"]
		}
	}`
)

// testV3EarlyAccessHeader verifies the X-Early-Access header required by
// every v3 endpoint is present on the outgoing request.
func testV3EarlyAccessHeader(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get("X-Early-Access"); got == "" {
		t.Error("X-Early-Access header is missing from v3 request")
	}
}

// ---------------------------------------------------------------------------
// unmarshalApiErrorObject — v3 map[string][]string error format
// ---------------------------------------------------------------------------

func TestUnmarshalApiErrorObject_V3MapFormat(t *testing.T) {
	data := []byte(`{"code":2001,"message":"Unprocessable Entity","errors":{"name":["can't be blank"]}}`)
	aeo, err := unmarshalApiErrorObject(data)
	if err != nil {
		t.Fatal(err)
	}
	if aeo.Code != 2001 {
		t.Errorf("Code = %d, want 2001", aeo.Code)
	}
	if aeo.Message != "Unprocessable Entity" {
		t.Errorf("Message = %q, want %q", aeo.Message, "Unprocessable Entity")
	}
	if len(aeo.Errors) != 1 {
		t.Fatalf("len(Errors) = %d, want 1", len(aeo.Errors))
	}
	if aeo.Errors[0] != "name: can't be blank" {
		t.Errorf("Errors[0] = %q, want %q", aeo.Errors[0], "name: can't be blank")
	}
}

// ---------------------------------------------------------------------------
// ListSchedulesV3
// ---------------------------------------------------------------------------

func TestScheduleV3_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockScheduleV3ListResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListSchedulesV3(context.Background(), ListSchedulesV3Options{})
	if err != nil {
		t.Fatal(err)
	}

	want := &ListSchedulesV3Response{
		Schedules: []APIObject{
			{
				ID:      testScheduleV3ID,
				Type:    "schedule_reference",
				Summary: "On-Call Schedule",
				Self:    "https://api.pagerduty.com/v3/schedules/SCHED01",
				HTMLURL: "https://example.pagerduty.com/schedules/SCHED01",
			},
		},
	}
	testEqual(t, want, res)
}

func TestScheduleV3_ListWithQuery(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if got := r.URL.Query().Get("query"); got != "on-call" {
			t.Errorf("query param = %q, want %q", got, "on-call")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockScheduleV3ListEmptyResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListSchedulesV3(context.Background(), ListSchedulesV3Options{Query: "on-call"})
	if err != nil {
		t.Fatal(err)
	}

	want := &ListSchedulesV3Response{Schedules: []APIObject{}}
	testEqual(t, want, res)
}

func TestScheduleV3_List500Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(mockScheduleV3Error500))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.ListSchedulesV3(context.Background(), ListSchedulesV3Options{})
	if !testErrCheck(t, "ListSchedulesV3", "status code 500", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// GetScheduleV3
// ---------------------------------------------------------------------------

func TestScheduleV3_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockScheduleV3GetResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetScheduleV3(context.Background(), testScheduleV3ID)
	if err != nil {
		t.Fatal(err)
	}

	want := &ScheduleV3{
		ID:          testScheduleV3ID,
		Type:        "schedule",
		Name:        "On-Call Schedule",
		TimeZone:    "UTC",
		Description: "Test schedule",
		// Events key absent in JSON → nil slice
		Rotations: []RotationV3{
			{ID: testRotationV3ID, Type: "rotation"},
		},
	}
	testEqual(t, want, res)
}

func TestScheduleV3_Get404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.GetScheduleV3(context.Background(), testScheduleV3ID)
	if !testErrCheck(t, "GetScheduleV3", "Not Found", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// CreateScheduleV3
// ---------------------------------------------------------------------------

func TestScheduleV3_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockScheduleV3MutateResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.CreateScheduleV3(context.Background(), ScheduleV3Input{
		Name:        "On-Call Schedule",
		TimeZone:    "UTC",
		Description: "Test schedule",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &ScheduleV3{
		ID:          testScheduleV3ID,
		Type:        "schedule",
		Name:        "On-Call Schedule",
		TimeZone:    "UTC",
		Description: "Test schedule",
	}
	testEqual(t, want, res)
}

func TestScheduleV3_Create400Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(mockScheduleV3Error400))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.CreateScheduleV3(context.Background(), ScheduleV3Input{})
	if !testErrCheck(t, "CreateScheduleV3", "Invalid Input", err) {
		return
	}
}

// The v3 API must respond 201; any other 2xx is treated as an error.
func TestScheduleV3_CreateNon201Status(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 instead of 201
		_, _ = w.Write([]byte(mockScheduleV3MutateResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.CreateScheduleV3(context.Background(), ScheduleV3Input{Name: "Test", TimeZone: "UTC"})
	if !testErrCheck(t, "CreateScheduleV3", "failed to create v3 schedule", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// UpdateScheduleV3
// ---------------------------------------------------------------------------

func TestScheduleV3_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockScheduleV3MutateResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateScheduleV3(context.Background(), testScheduleV3ID, ScheduleV3Input{
		Name:        "On-Call Schedule",
		TimeZone:    "UTC",
		Description: "Test schedule",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &ScheduleV3{
		ID:          testScheduleV3ID,
		Type:        "schedule",
		Name:        "On-Call Schedule",
		TimeZone:    "UTC",
		Description: "Test schedule",
	}
	testEqual(t, want, res)
}

func TestScheduleV3_Update404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateScheduleV3(context.Background(), testScheduleV3ID, ScheduleV3Input{})
	if !testErrCheck(t, "UpdateScheduleV3", "Not Found", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// DeleteScheduleV3
// ---------------------------------------------------------------------------

func TestScheduleV3_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testV3EarlyAccessHeader(t, r)
		w.WriteHeader(http.StatusNoContent)
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteScheduleV3(context.Background(), testScheduleV3ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestScheduleV3_Delete404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteScheduleV3(context.Background(), testScheduleV3ID)
	if !testErrCheck(t, "DeleteScheduleV3", "Not Found", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// CreateRotationV3
// ---------------------------------------------------------------------------

func TestRotationV3_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockRotationV3Response))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.CreateRotationV3(context.Background(), testScheduleV3ID)
	if err != nil {
		t.Fatal(err)
	}

	want := &RotationV3{
		ID:     testRotationV3ID,
		Type:   "rotation",
		Events: []EventV3{},
	}
	testEqual(t, want, res)
}

func TestRotationV3_CreateNon201Status(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 instead of 201
		_, _ = w.Write([]byte(mockRotationV3Response))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.CreateRotationV3(context.Background(), testScheduleV3ID)
	if !testErrCheck(t, "CreateRotationV3", "failed to create v3 rotation", err) {
		return
	}
}

func TestRotationV3_Create404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.CreateRotationV3(context.Background(), testScheduleV3ID)
	if !testErrCheck(t, "CreateRotationV3", "Not Found", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// GetRotationV3
// ---------------------------------------------------------------------------

func TestRotationV3_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockRotationV3Response))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetRotationV3(context.Background(), testScheduleV3ID, testRotationV3ID)
	if err != nil {
		t.Fatal(err)
	}

	want := &RotationV3{
		ID:     testRotationV3ID,
		Type:   "rotation",
		Events: []EventV3{},
	}
	testEqual(t, want, res)
}

func TestRotationV3_Get404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.GetRotationV3(context.Background(), testScheduleV3ID, testRotationV3ID)
	if !testErrCheck(t, "GetRotationV3", "Not Found", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// DeleteRotationV3
// ---------------------------------------------------------------------------

func TestRotationV3_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testV3EarlyAccessHeader(t, r)
		w.WriteHeader(http.StatusNoContent)
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteRotationV3(context.Background(), testScheduleV3ID, testRotationV3ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRotationV3_Delete404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteRotationV3(context.Background(), testScheduleV3ID, testRotationV3ID)
	if !testErrCheck(t, "DeleteRotationV3", "Not Found", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// CreateEventV3
// ---------------------------------------------------------------------------

func TestEventV3_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID+"/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockEventV3Response))
	})

	userID := "USER01"
	input := EventV3{
		Name:           "On-Call Event",
		StartTime:      EventTimeV3{DateTime: "2026-02-21T09:00:00Z"},
		EndTime:        EventTimeV3{DateTime: "2026-02-21T17:00:00Z"},
		EffectiveSince: "2026-02-21T09:00:00Z",
		Recurrence:     []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"},
		AssignmentStrategy: AssignmentStrategyV3{
			Type:    "user_assignment_strategy",
			Members: []MemberV3{{Type: "user_member", UserID: &userID}},
		},
	}

	client := defaultTestClient(server.URL, "foo")
	res, err := client.CreateEventV3(context.Background(), testScheduleV3ID, testRotationV3ID, input)
	if err != nil {
		t.Fatal(err)
	}

	want := &EventV3{
		ID:             testEventV3ID,
		Type:           "schedule_event",
		Name:           "On-Call Event",
		StartTime:      EventTimeV3{DateTime: "2026-02-21T09:00:00Z"},
		EndTime:        EventTimeV3{DateTime: "2026-02-21T17:00:00Z"},
		EffectiveSince: "2026-02-21T09:00:00Z",
		EffectiveUntil: nil,
		Recurrence:     []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"},
		AssignmentStrategy: AssignmentStrategyV3{
			Type:    "user_assignment_strategy",
			Members: []MemberV3{{Type: "user_member", UserID: &userID}},
		},
	}
	testEqual(t, want, res)
}

func TestEventV3_CreateNon201Status(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID+"/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 instead of 201
		_, _ = w.Write([]byte(mockEventV3Response))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.CreateEventV3(context.Background(), testScheduleV3ID, testRotationV3ID, EventV3{})
	if !testErrCheck(t, "CreateEventV3", "failed to create v3 event", err) {
		return
	}
}

func TestEventV3_Create400Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID+"/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(mockScheduleV3Error400))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.CreateEventV3(context.Background(), testScheduleV3ID, testRotationV3ID, EventV3{})
	if !testErrCheck(t, "CreateEventV3", "Invalid Input", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// UpdateEventV3
// ---------------------------------------------------------------------------

func TestEventV3_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID+"/events/"+testEventV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testV3EarlyAccessHeader(t, r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockEventV3Response))
	})

	userID := "USER01"
	input := EventV3{
		Name:           "On-Call Event",
		StartTime:      EventTimeV3{DateTime: "2026-02-21T09:00:00Z"},
		EndTime:        EventTimeV3{DateTime: "2026-02-21T17:00:00Z"},
		EffectiveSince: "2026-02-21T09:00:00Z",
		Recurrence:     []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"},
		AssignmentStrategy: AssignmentStrategyV3{
			Type:    "user_assignment_strategy",
			Members: []MemberV3{{Type: "user_member", UserID: &userID}},
		},
	}

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateEventV3(context.Background(), testScheduleV3ID, testRotationV3ID, testEventV3ID, input)
	if err != nil {
		t.Fatal(err)
	}

	want := &EventV3{
		ID:             testEventV3ID,
		Type:           "schedule_event",
		Name:           "On-Call Event",
		StartTime:      EventTimeV3{DateTime: "2026-02-21T09:00:00Z"},
		EndTime:        EventTimeV3{DateTime: "2026-02-21T17:00:00Z"},
		EffectiveSince: "2026-02-21T09:00:00Z",
		EffectiveUntil: nil,
		Recurrence:     []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"},
		AssignmentStrategy: AssignmentStrategyV3{
			Type:    "user_assignment_strategy",
			Members: []MemberV3{{Type: "user_member", UserID: &userID}},
		},
	}
	testEqual(t, want, res)
}

func TestEventV3_Update404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID+"/events/"+testEventV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateEventV3(context.Background(), testScheduleV3ID, testRotationV3ID, testEventV3ID, EventV3{})
	if !testErrCheck(t, "UpdateEventV3", "Not Found", err) {
		return
	}
}

// ---------------------------------------------------------------------------
// DeleteEventV3
// ---------------------------------------------------------------------------

func TestEventV3_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID+"/events/"+testEventV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testV3EarlyAccessHeader(t, r)
		w.WriteHeader(http.StatusNoContent)
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteEventV3(context.Background(), testScheduleV3ID, testRotationV3ID, testEventV3ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEventV3_Delete404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/schedules/"+testScheduleV3ID+"/rotations/"+testRotationV3ID+"/events/"+testEventV3ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockScheduleV3Error404))
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteEventV3(context.Background(), testScheduleV3ID, testRotationV3ID, testEventV3ID)
	if !testErrCheck(t, "DeleteEventV3", "Not Found", err) {
		return
	}
}

package pagerduty

import (
	"net/http"
	"testing"
)

// ListSchedules
func TestSchedule_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"schedules": [{"id": "1","summary":"foo"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	opts := ListSchedulesOptions{
		APIListObject: listObj,
		Query:         "foo",
	}
	res, err := client.ListSchedules(opts)

	want := &ListSchedulesResponse{
		APIListObject: listObj,
		Schedules: []Schedule{
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

// Create a Schedule
func TestSchedule_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"schedule": {"id": "1","summary":"foo"}}`))
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}
	res, err := client.CreateSchedule(input)

	want := &Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TODO: Preview a schedule -- should this function be changed to actually return a preview?

// Delete a schedule
func TestSchedule_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	err := client.DeleteSchedule(id)
	if err != nil {
		t.Fatal(err)
	}
}

// Get a schedule
func TestSchedule_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"schedule": {"id": "1","summary":"foo"}}`))
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}

	input := "1"
	opts := GetScheduleOptions{
		APIListObject: listObj,
		TimeZone:      "UTC",
		Since:         "foo",
		Until:         "bar",
	}
	res, err := client.GetSchedule(input, opts)

	want := &Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update a schedule
func TestSchedule_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"schedule": {"id": "1","summary":"foo"}}`))
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	id := "1"
	sched := Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}
	res, err := client.UpdateSchedule(id, sched)

	want := &Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// List overrides
func TestSchedule_ListOverrides(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules/1/overrides", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"overrides": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	opts := ListOverridesOptions{
		APIListObject: listObj,
		Since:         "foo",
		Until:         "bar",
		Editable:      false,
		Overflow:      false,
	}
	schedID := "1"

	res, err := client.ListOverrides(schedID, opts)

	want := &ListOverridesResponse{
		APIListObject: listObj,
		Overrides: []Override{
			{
				ID: "1",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create an override
func TestSchedule_CreateOverride(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules/1/overrides", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"override": {"id": "1", "start": "foo", "end": "bar"}}`))
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := Override{
		Start: "foo",
		End:   "bar",
	}
	schedID := "1"

	res, err := client.CreateOverride(schedID, input)

	want := &Override{
		ID:    "1",
		Start: "foo",
		End:   "bar",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete an override
func TestSchedule_DeleteOverride(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules/1/overrides/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	schedID := "1"
	overID := "1"
	err := client.DeleteOverride(schedID, overID)
	if err != nil {
		t.Fatal(err)
	}
}

// List users on call
func TestSchedule_ListOnCallUsers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/schedules/1/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"users": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	opts := ListOnCallUsersOptions{
		APIListObject: listObj,
		Since:         "foo",
		Until:         "bar",
	}
	schedID := "1"

	res, err := client.ListOnCallUsers(schedID, opts)

	want := []User{
		{
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

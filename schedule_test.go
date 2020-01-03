package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// ListSchedules
func TestSchedule_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"schedules": [{"id": "1","summary":"foo"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListSchedulesOptions{
		APIListObject: listObj,
		Query:         "foo",
	}
	resp, err := client.ListSchedules(opts)

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

	require.NoError(err)
	require.Equal(want, resp)
}

// Create a Schedule
func TestSchedule_Create(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"schedule": {"id": "1","summary":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}
	resp, err := client.CreateSchedule(input)

	want := &Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// TODO: Preview a schedule -- should this function be changed to actually return a preview?

// Delete a schedule
func TestSchedule_Delete(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	err := client.DeleteSchedule(id)

	require.NoError(err)
}

// Get a schedule
func TestSchedule_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"schedule": {"id": "1","summary":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}

	input := "1"
	opts := GetScheduleOptions{
		APIListObject: listObj,
		TimeZone:      "UTC",
		Since:         "foo",
		Until:         "bar",
	}
	resp, err := client.GetSchedule(input, opts)

	want := &Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Update a schedule
func TestSchedule_Update(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"schedule": {"id": "1","summary":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	id := "1"
	sched := Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}
	resp, err := client.UpdateSchedule(id, sched)

	want := &Schedule{
		APIObject: APIObject{
			ID:      "1",
			Summary: "foo",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// List overrides
func TestSchedule_ListOverrides(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules/1/overrides", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"overrides": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListOverridesOptions{
		APIListObject: listObj,
		Since:         "foo",
		Until:         "bar",
		Editable:      false,
		Overflow:      false,
	}
	schedID := "1"

	resp, err := client.ListOverrides(schedID, opts)

	want := &ListOverridesResponse{
		APIListObject: listObj,
		Overrides: []Override{
			{
				ID: "1",
			},
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Create an override
func TestSchedule_CreateOverride(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules/1/overrides", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"override": {"id": "1", "start": "foo", "end": "bar"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var input = Override{
		Start: "foo",
		End:   "bar",
	}
	schedID := "1"

	resp, err := client.CreateOverride(schedID, input)

	want := &Override{
		ID:    "1",
		Start: "foo",
		End:   "bar",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Delete an override
func TestSchedule_DeleteOverride(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules/1/overrides/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	schedID := "1"
	overID := "1"
	err := client.DeleteOverride(schedID, overID)

	require.NoError(err)
}

// List users on call
func TestSchedule_ListOnCallUsers(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/schedules/1/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"users": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListOnCallUsersOptions{
		APIListObject: listObj,
		Since:         "foo",
		Until:         "bar",
	}
	schedID := "1"

	resp, err := client.ListOnCallUsers(schedID, opts)

	want := []User{
		{
			APIObject: APIObject{
				ID: "1",
			},
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

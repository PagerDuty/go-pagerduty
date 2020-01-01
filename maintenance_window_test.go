package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// ListMaintenanceWindows
func TestMaintenanceWindow_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/maintenance_windows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"maintenance_windows": [{"id": "1", "summary": "foo"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListMaintenanceWindowsOptions{
		APIListObject: listObj,
		Query:         "foo",
		Includes:      []string{},
		TeamIDs:       []string{},
		ServiceIDs:    []string{},
		Filter:        "foo",
	}
	resp, err := client.ListMaintenanceWindows(opts)

	want := &ListMaintenanceWindowsResponse{
		APIListObject: listObj,
		MaintenanceWindows: []MaintenanceWindow{
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

// CreateMaintenanceWindow
func TestMaintenanceWindow_Create(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	input := MaintenanceWindow{Description: "foo"}
	from := "foo@bar.com"

	mux.HandleFunc("/maintenance_windows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	resp, err := client.CreateMaintenanceWindow(from, input)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}
func TestMaintenanceWindow_Create_NoFrom(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	input := MaintenanceWindow{Description: "foo"}
	from := ""

	mux.HandleFunc("/maintenance_windows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	resp, err := client.CreateMaintenanceWindow(from, input)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// DeleteMaintenanceWindows
func TestMaintenanceWindow_Delete(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/maintenance_windows/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	err := client.DeleteMaintenanceWindow("1")

	require.NoError(err)
}

// GetMaintenanceWindow
func TestMaintenanceWindow_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/maintenance_windows/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	opts := GetMaintenanceWindowOptions{Includes: []string{}}
	resp, err := client.GetMaintenanceWindow(id, opts)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// UpdateMaintenanceWindow
func TestMaintenanceWindow_Update(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	input := MaintenanceWindow{
		APIObject: APIObject{
			ID: "1",
		},
		Description: "foo",
	}

	mux.HandleFunc("/maintenance_windows/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	resp, err := client.UpdateMaintenanceWindow(input)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}
	require.NoError(err)
	require.Equal(want, resp)
}

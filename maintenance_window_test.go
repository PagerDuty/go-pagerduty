package pagerduty

import (
	"net/http"
	"testing"
)

// ListMaintenanceWindows
func TestMaintenanceWindow_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/maintenance_windows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"maintenance_windows": [{"id": "1", "summary": "foo"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListMaintenanceWindowsOptions{
		APIListObject: listObj,
		Query:         "foo",
		Includes:      []string{},
		TeamIDs:       []string{},
		ServiceIDs:    []string{},
		Filter:        "foo",
	}
	res, err := client.ListMaintenanceWindows(opts)

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

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// CreateMaintenanceWindow
func TestMaintenanceWindow_Create(t *testing.T) {
	setup()
	defer teardown()

	input := MaintenanceWindow{Description: "foo"}
	from := "foo@bar.com"

	mux.HandleFunc("/maintenance_windows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.CreateMaintenanceWindow(from, input)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestMaintenanceWindow_Create_NoFrom(t *testing.T) {
	setup()
	defer teardown()

	input := MaintenanceWindow{Description: "foo"}
	from := "foo@bar.com"

	mux.HandleFunc("/maintenance_windows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.CreateMaintenanceWindow(from, input)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// DeleteMaintenanceWindows
func TestMaintenanceWindow_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/maintenance_windows/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteMaintenanceWindow("1")
	if err != nil {
		t.Fatal(err)
	}
}

// GetMaintenanceWindow
func TestMaintenanceWindow_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/maintenance_windows/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"
	opts := GetMaintenanceWindowOptions{Includes: []string{}}
	res, err := client.GetMaintenanceWindow(id, opts)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// UpdateMaintenanceWindow
func TestMaintenanceWindow_Update(t *testing.T) {
	setup()
	defer teardown()

	input := MaintenanceWindow{
		APIObject: APIObject{
			ID: "1",
		},
		Description: "foo",
	}

	mux.HandleFunc("/maintenance_windows/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"maintenance_window": {"description": "foo", "id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")

	res, err := client.UpdateMaintenanceWindow(input)

	want := &MaintenanceWindow{
		Description: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

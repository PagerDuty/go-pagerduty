package pagerduty

import (
	"net/http"
	"testing"
)

func TestOrchestration_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"orchestrations": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var opts ListOrchestrationsOptions
	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListOrchestrations(opts)

	want := &ListOrchestrationsResponse{
		APIListObject: listObj,
		Orchestrations: []Orchestration{
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

func TestOrchestration_Create(t *testing.T) {
	setup()
	defer teardown()

	input := Orchestration{Name: "foo"}

	mux.HandleFunc("/event_orchestrations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"orchestration": {"name": "foo", "id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	res, err := client.CreateOrchestration(input)

	want := &Orchestration{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestOrchestration_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteOrchestration("1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrchestration_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"orchestration": {"id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	var opts *GetOrchestrationOptions
	res, err := client.GetOrchestration("1", opts)

	want := &Orchestration{
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestOrchestration_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"orchestration": {"name": "foo", "id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &Orchestration{Name: "foo"}
	want := &Orchestration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	res, err := client.UpdateOrchestration("1", input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestOrchestrationRouter_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/1/router", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"orchestration_path": {"type": "router"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	var opts *GetOrchestrationRouterOptions
	res, err := client.GetOrchestrationRouter("1", opts)

	want := &OrchestrationRouter{
		Type: "router",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestOrchestrationRouter_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/1/router", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"orchestration_path": {"type": "router", "parent": {"id": "1"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &OrchestrationRouter{Type: "router"}
	want := &OrchestrationRouter{
		Type: "router",
		Parent: &APIReference{
			ID: "1",
		},
	}
	res, err := client.UpdateOrchestrationRouter("1", input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestOrchestrationUnrouted_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/1/unrouted", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"orchestration_path": {"type": "unrouted", "parent": {"id": "1"}}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	var opts *GetOrchestrationUnroutedOptions
	res, err := client.GetOrchestrationUnrouted("1", opts)

	want := &OrchestrationUnrouted{
		Type: "unrouted",
		Parent: &APIReference{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestOrchestrationUnrouted_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/1/unrouted", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"orchestration_path": {"type": "unrouted", "parent": {"id": "1"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &OrchestrationUnrouted{Type: "unrouted"}
	want := &OrchestrationUnrouted{
		Type: "unrouted",
		Parent: &APIReference{
			ID: "1",
		},
	}
	res, err := client.UpdateOrchestrationUnrouted("1", input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestServiceOrchestration_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"orchestration_path": {"type": "service", "parent": {"id": "1"}}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	var opts *GetServiceOrchestrationOptions
	res, err := client.GetServiceOrchestration("1", opts)

	want := &ServiceOrchestration{
		Type: "service",
		Parent: &APIReference{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestServiceOrchestration_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"orchestration_path": {"type": "service", "parent": {"id": "1"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ServiceOrchestration{Type: "service"}
	want := &ServiceOrchestration{
		Type: "service",
		Parent: &APIReference{
			ID: "1",
		},
	}
	res, err := client.UpdateServiceOrchestration("1", input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestServiceOrchestrationActive_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/services/1/active", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"active": true}`))
	})
	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetServiceOrchestrationActive("1")

	want := &ServiceOrchestrationActive{
		Active: true,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestServiceOrchestrationActive_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/services/1/active", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"active": true}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ServiceOrchestrationActive{Active: true}
	want := &ServiceOrchestrationActive{
		Active: true,
	}
	res, err := client.UpdateServiceOrchestrationActive("1", input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

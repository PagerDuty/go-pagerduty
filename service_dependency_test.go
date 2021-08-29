package pagerduty

import (
	"net/http"
	"testing"
)

// List BusinessService Dependencies
func TestBusinessServiceDependency_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/service_dependencies/business_services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"relationships": [{"id": "1","dependent_service":{"id":"1"},"supporting_service":{"id":"1"},"type":"service_dependency"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	bServeID := "1"
	res, _, err := client.ListBusinessServiceDependencies(bServeID)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListServiceDependencies{
		Relationships: []*ServiceDependency{
			{
				ID:   "1",
				Type: "service_dependency",
				DependentService: &ServiceObj{
					ID: "1",
				},
				SupportingService: &ServiceObj{
					ID: "1",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// List TechnicalService Dependencies
func TestTechnicalServiceDependency_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/service_dependencies/technical_services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"relationships": [{"id": "1","dependent_service":{"id":"1"},"supporting_service":{"id":"1"},"type":"service_dependency"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	bServeID := "1"
	res, _, err := client.ListTechnicalServiceDependencies(bServeID)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListServiceDependencies{
		Relationships: []*ServiceDependency{
			{
				ID:   "1",
				Type: "service_dependency",
				DependentService: &ServiceObj{
					ID: "1",
				},
				SupportingService: &ServiceObj{
					ID: "1",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// AssociateServiceDependencies
func TestServiceDependency_Associate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/service_dependencies/associate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"relationships": [{"id": "1","dependent_service":{"id":"1"},"supporting_service":{"id":"1"},"type":"service_dependency"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ListServiceDependencies{
		Relationships: []*ServiceDependency{
			{
				ID:   "1",
				Type: "service_dependency",
				DependentService: &ServiceObj{
					ID: "1",
				},
				SupportingService: &ServiceObj{
					ID: "1",
				},
			},
		},
	}
	res, _, err := client.AssociateServiceDependencies(input)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListServiceDependencies{
		Relationships: []*ServiceDependency{
			{
				ID:   "1",
				Type: "service_dependency",
				DependentService: &ServiceObj{
					ID: "1",
				},
				SupportingService: &ServiceObj{
					ID: "1",
				},
			},
		},
	}
	testEqual(t, want, res)
}

// DisassociateServiceDependencies
func TestServiceDependency_Disassociate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/service_dependencies/disassociate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"relationships": [{"id": "1","dependent_service":{"id":"1"},"supporting_service":{"id":"1"},"type":"service_dependency"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ListServiceDependencies{
		Relationships: []*ServiceDependency{
			{
				ID:   "1",
				Type: "service_dependency",
				DependentService: &ServiceObj{
					ID: "1",
				},
				SupportingService: &ServiceObj{
					ID: "1",
				},
			},
		},
	}
	res, _, err := client.DisassociateServiceDependencies(input)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListServiceDependencies{
		Relationships: []*ServiceDependency{
			{
				ID:   "1",
				Type: "service_dependency",
				DependentService: &ServiceObj{
					ID: "1",
				},
				SupportingService: &ServiceObj{
					ID: "1",
				},
			},
		},
	}
	testEqual(t, want, res)
}

package pagerduty

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

// ListServices
func TestService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"services": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListServiceOptions{
		APIListObject: listObj,
		TeamIDs:       []string{},
		TimeZone:      "foo",
		SortBy:        "bar",
		Query:         "baz",
		Includes:      []string{},
	}
	res, err := client.ListServices(opts)

	want := &ListServiceResponse{
		APIListObject: listObj,
		Services: []Service{
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

// ListServices
func TestService_ListPaginated(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		offsetStr := r.URL.Query()["offset"][0]
		offset, _ := strconv.ParseInt(offsetStr, 10, 32)

		var more string
		if offset == 0 {
			more = "true"
		} else {
			more = "false"
		}
		resp := fmt.Sprintf(`{"services": [{"id": "%d"}],
                          "More": %s,
                          "Offset": %d,
                          "Limit": 1}`, offset, more, offset)
		_, _ = w.Write([]byte(resp))
	})

	listObj := APIListObject{Limit: 1, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListServiceOptions{
		APIListObject: listObj,
		TeamIDs:       []string{},
		TimeZone:      "foo",
		SortBy:        "bar",
		Query:         "baz",
		Includes:      []string{},
	}
	res, err := client.ListServicesPaginated(context.Background(), opts)

	want := []Service{
		{
			APIObject: APIObject{
				ID: "0",
			},
		},
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

// Get Service
func TestService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	id := "1"
	opts := &GetServiceOptions{
		Includes: []string{},
	}
	res, err := client.GetService(id, opts)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create Service
func TestService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Service{
		Name: "foo",
	}
	res, err := client.CreateService(input)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create Service with AlertGroupingParameters of type time
func TestService_CreateWithAlertGroupParamsTime(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Service{
		Name: "foo",
		AlertGroupingParameters: &AlertGroupingParameters{
			Type: "time",
			Config: AlertGroupParamsConfig{
				Timeout: 2,
			},
		},
	}
	res, err := client.CreateService(input)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create Service with AlertGroupingParameters of type content_based
func TestService_CreateWithAlertGroupParamsContentBased(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Service{
		Name: "foo",
		AlertGroupingParameters: &AlertGroupingParameters{
			Type: "content_based",
			Config: AlertGroupParamsConfig{
				Aggregate: "any",
				Fields:    []string{"source", "component"},
			},
		},
	}
	res, err := client.CreateService(input)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create Service with AlertGroupingParameters of type intelligent
func TestService_CreateWithAlertGroupParamsIntelligent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Service{
		Name: "foo",
		AlertGroupingParameters: &AlertGroupingParameters{
			Type: "intelligent",
		},
	}
	res, err := client.CreateService(input)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update Service
func TestService_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	res, err := client.UpdateService(input)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete Service
func TestService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"
	err := client.DeleteService(id)
	if err != nil {
		t.Fatal(err)
	}
}

// Create Integration
func TestService_CreateIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Integration{
		Name: "foo",
	}
	servID := "1"

	res, err := client.CreateIntegration(servID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Get Integration
func TestService_GetIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := GetIntegrationOptions{
		Includes: []string{},
	}
	servID := "1"
	intID := "1"

	res, err := client.GetIntegration(servID, intID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update Integration
func TestService_UpdateIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	servID := "1"

	res, err := client.UpdateIntegration(servID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete Integration
func TestService_DeleteIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	servID := "1"
	intID := "1"
	err := client.DeleteIntegration(servID, intID)
	if err != nil {
		t.Fatal(err)
	}
}

// List Service Rules
func TestService_ListRules(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"rules": [{"id": "1"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")

	serviceID := "1"
	res, err := client.ListServiceRulesPaginated(context.Background(), serviceID)
	if err != nil {
		t.Fatal(err)
	}

	want := []ServiceRule{{ID: "1"}}
	testEqual(t, want, res)
}

// Create Service Rule
func TestService_CreateServiceRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/rules/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"rule": {"id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	serviceID := "1"
	rule := ServiceRule{}

	res, err := client.CreateServiceRule(context.Background(), serviceID, rule)
	if err != nil {
		t.Fatal(err)
	}

	want := ServiceRule{
		ID: "1",
	}
	testEqual(t, want, res)
}

// Get Service Rule
func TestService_GetServiceRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"rule": {"id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	serviceID := "1"
	ruleID := "1"
	res, err := client.GetServiceRule(context.Background(), serviceID, ruleID)
	if err != nil {
		t.Fatal(err)
	}

	want := ServiceRule{
		ID: "1",
	}
	testEqual(t, want, res)
}

// Update Service Rule
func TestService_UpdateServiceRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"rule": {"id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	serviceID := "1"
	ruleID := "1"
	rule := ServiceRule{}

	res, err := client.UpdateServiceRule(context.Background(), serviceID, ruleID, rule)
	if err != nil {
		t.Fatal(err)
	}

	want := ServiceRule{
		ID: "1",
	}
	testEqual(t, want, res)
}

// Delete Service Rule
func TestService_DeleteServiceRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	serviceID := "1"
	ruleID := "1"

	err := client.DeleteServiceRule(context.Background(), serviceID, ruleID)
	if err != nil {
		t.Fatal(err)
	}
}

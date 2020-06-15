package pagerduty

import (
	"net/http"
	"testing"
)

// List BusinessServices
func TestBusinessService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/business_services/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"business_services": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListBusinessServiceOptions{
		APIListObject: listObj,
	}
	res, err := client.ListBusinessServices(opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListBusinessServicesResponse{
		BusinessServices: []*BusinessService{
			{
				ID: "1",
			},
		},
	}

	testEqual(t, want, res)
}

// Create BusinessService
func TestBusinessService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/business_services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"business_service": {"id": "1", "name": "foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := &BusinessService{
		Name: "foo",
	}
	res, _, err := client.CreateBusinessService(input)

	want := &BusinessService{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Get BusinessService
func TestBusinessService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/business_services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"business_service": {"id": "1", "name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	ruleSetID := "1"

	res, _, err := client.GetBusinessService(ruleSetID)

	want := &BusinessService{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update BusinessService
func TestBusinessService_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/business_services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"business_service": {"id": "1", "name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := &BusinessService{
		ID:   "1",
		Name: "foo",
	}
	res, _, err := client.UpdateBusinessService(input)

	want := &BusinessService{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete BusinessService
func TestBusinessService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/business_services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	ID := "1"
	err := client.DeleteBusinessService(ID)

	if err != nil {
		t.Fatal(err)
	}
}

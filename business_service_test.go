package pagerduty

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

// List BusinessServices
func TestBusinessService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/business_services/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"business_services": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListBusinessServiceOptions{
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
		_, _ = w.Write([]byte(`{"business_service": {"id": "1", "name": "foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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
		_, _ = w.Write([]byte(`{"business_service": {"id": "1", "name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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
		reqText, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		v := make(map[string]*BusinessService)
		if err := json.Unmarshal(reqText, &v); err != nil {
			t.Fatal(err)
		}
		if v["business_service"].ID != "" {
			t.Fatalf("got ID in the body when we were not supposed to")
		}

		_, _ = w.Write([]byte(`{"business_service": {"id": "1", "name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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

	client := defaultTestClient(server.URL, "foo")
	ID := "1"
	err := client.DeleteBusinessService(ID)
	if err != nil {
		t.Fatal(err)
	}
}

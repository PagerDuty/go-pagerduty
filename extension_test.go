package pagerduty

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestExtension_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extensions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"extensions":[{"id":"1","summary":"foo","config": {"restrict": "any"}, "extension_objects":[{"id":"foo","summary":"foo"}]}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListExtensionOptions{
		APIListObject: listObj,
		Query:         "foo",
	}

	res, err := client.ListExtensions(opts)

	want := &ListExtensionResponse{
		APIListObject: listObj,
		Extensions: []Extension{
			{
				APIObject: APIObject{
					ID:      "1",
					Summary: "foo",
				},
				Config: map[string]interface{}{
					"restrict": "any",
				},
				ExtensionObjects: []APIObject{
					{
						ID:      "foo",
						Summary: "foo",
					},
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestExtension_Create(t *testing.T) {
	setup()
	defer teardown()

	input1 := &Extension{Name: "foo"}
	input2 := &Extension{Name: "bar", EndpointURL: "expected_url"}

	mux.HandleFunc("/extensions", func(w http.ResponseWriter, r *http.Request) {
		var got map[string]interface{}

		err := json.NewDecoder(r.Body).Decode(&got)

		testErrCheck(t, "Extension_Create()", "", err)
		name := got["name"]

		if name == "foo" {
			testNoEndpointURL(t, got)
			testMethod(t, r, "POST")
			_, _ = w.Write([]byte(`{"extension": {"name": "foo", "id": "1"}}`))
		} else {
			testGotExpectedURL(t, "expected_url", got)
			testMethod(t, r, "POST")
			_, _ = w.Write([]byte(`{"extension": {"name": "bar", "id": "2", "endpoint_url": "expected_url"}}`))
		}
	})

	client := defaultTestClient(server.URL, "foo")

	want1 := &Extension{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	want2 := &Extension{
		Name:        "bar",
		EndpointURL: "expected_url",
		APIObject: APIObject{
			ID: "2",
		},
	}

	res1, err := client.CreateExtension(input1)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want1, res1)

	res2, err := client.CreateExtension(input2)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want2, res2)
}

func TestExtension_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extensions/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	client := defaultTestClient(server.URL, "foo")

	if err := client.DeleteExtension("1"); err != nil {
		t.Fatal(err)
	}
}

func TestExtension_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extensions/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"extension": {"name": "foo", "id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetExtension("1")

	want := &Extension{
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

func TestExtension_Update(t *testing.T) {
	setup()
	defer teardown()

	input1 := &Extension{Name: "foo"}
	input2 := &Extension{Name: "foo", EndpointURL: "expected_url"}

	mux.HandleFunc("/extensions/1", func(w http.ResponseWriter, r *http.Request) {
		var got map[string]interface{}

		err := json.NewDecoder(r.Body).Decode(&got)

		testErrCheck(t, "Extension_Update()", "", err)
		testNoEndpointURL(t, got)

		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"extension": {"name": "foo", "id": "1"}}`))
	})

	mux.HandleFunc("/extensions/2", func(w http.ResponseWriter, r *http.Request) {
		var got map[string]interface{}

		err := json.NewDecoder(r.Body).Decode(&got)
		testErrCheck(t, "Extension_Update()", "", err)

		testGotExpectedURL(t, "expected_url", got)

		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"extension": {"name": "foo", "id": "2", "endpoint_url": "expected_url"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	want1 := &Extension{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	want2 := &Extension{
		Name:        "foo",
		EndpointURL: "expected_url",
		APIObject: APIObject{
			ID: "2",
		},
	}

	res1, err := client.UpdateExtension("1", input1)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want1, res1)

	res2, err := client.UpdateExtension("2", input2)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want2, res2)
}

func testNoEndpointURL(t *testing.T, got map[string]interface{}) {
	if _, ok := got["endpoint_url"]; ok {
		t.Errorf(`Expected no url, got: "%v"`, got["endpoint_url"])
	}
}

func testGotExpectedURL(t *testing.T, expected string, got map[string]interface{}) {
	if got["endpoint_url"] != expected {
		t.Errorf(`Expected url: "%v", got: "%v"`, "expected_url", got["endpoint_url"])
	}
}

package pagerduty

import (
	"net/http"
	"testing"
)

func TestExtension_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extensions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"extensions":[{"id":"1","summary":"foo","config": {"restrict": "any"}, "extension_objects":[{"id":"foo","summary":"foo"}]}]}`))

	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListExtensionOptions{
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

	input := &Extension{Name: "foo"}

	mux.HandleFunc("/extensions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"extension": {"name": "foo", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	res, err := client.CreateExtension(input)

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

func TestExtension_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extensions/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	if err := client.DeleteExtension("1"); err != nil {
		t.Fatal(err)
	}
}

func TestExtension_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extensions/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"extension": {"name": "foo", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

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

	input := &Extension{Name: "foo"}

	mux.HandleFunc("/extensions/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"extension": {"name": "foo", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	res, err := client.UpdateExtension("1", input)

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

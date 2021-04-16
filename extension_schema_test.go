package pagerduty

import (
	"net/http"
	"testing"
)

func TestExtensionSchema_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extension_schemas", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"extension_schemas":[{"id":"1","summary":"foo","send_types":["trigger", "acknowledge", "resolve"]}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListExtensionSchemaOptions{
		APIListObject: listObj,
		Query:         "foo",
	}

	res, err := client.ListExtensionSchemas(opts)

	want := &ListExtensionSchemaResponse{
		APIListObject: listObj,
		ExtensionSchemas: []ExtensionSchema{
			{
				APIObject: APIObject{
					ID:      "1",
					Summary: "foo",
				},
				SendTypes: []string{
					"trigger",
					"acknowledge",
					"resolve",
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestExtensionSchema_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/extension_schemas/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"extension_schema": {"name": "foo", "id": "1", "send_types": ["trigger", "acknowledge", "resolve"]}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetExtensionSchema("1")

	want := &ExtensionSchema{
		APIObject: APIObject{
			ID: "1",
		},
		SendTypes: []string{
			"trigger",
			"acknowledge",
			"resolve",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

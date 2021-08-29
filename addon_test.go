package pagerduty

import (
	"net/http"
	"testing"
)

func TestAddon_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/addons", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"addons": [{"name": "Internal Status Page"}]}`))
	})
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var opts ListAddonOptions
	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListAddons(opts)
	want := &ListAddonResponse{
		APIListObject: listObj,
		Addons: []Addon{
			{
				Name: "Internal Status Page",
			},
		},
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestAddon_Install(t *testing.T) {
	setup()
	defer teardown()

	input := Addon{
		Name: "Internal Status Page",
	}

	mux.HandleFunc("/addons", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"addon": {"name": "Internal Status Page", "id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.InstallAddon(input)

	want := &Addon{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "Internal Status Page",
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestAddon_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/addons/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"addon": {"id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetAddon("1")

	want := &Addon{
		APIObject: APIObject{
			ID: "1",
		},
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestAddon_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/addons/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"addon": {"name": "Internal Status Page", "id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")

	input := Addon{
		Name: "Internal Status Page",
	}

	res, err := client.UpdateAddon("1", input)

	want := &Addon{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "Internal Status Page",
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestAddon_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/addons/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteAddon("1")
	if err != nil {
		t.Fatal(err)
	}
}

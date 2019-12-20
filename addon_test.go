package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddon_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/addons", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"addons": [{"name": "Internal Status Page"}]}`))
	})
	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var opts ListAddonOptions
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	resp, err := client.ListAddons(opts)
	want := &ListAddonResponse{
		APIListObject: listObj,
		Addons: []Addon{
			{
				Name: "Internal Status Page",
			},
		},
	}
	require.NoError(err)
	require.Equal(want, resp)
}

func TestAddon_Install(t *testing.T) {
	setup()
	defer teardown()

	input := Addon{
		Name: "Internal Status Page",
	}

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/addons", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"addon": {"name": "Internal Status Page", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	addon, err := client.InstallAddon(input)

	want := &Addon{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "Internal Status Page",
	}
	require.Equal(want, addon)
	require.NoError(err)
}

func TestAddon_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/addons/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"addon": {"id": "1"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	addon, err := client.GetAddon("1")

	want := &Addon{
		APIObject: APIObject{
			ID: "1",
		},
	}

	require.Equal(want, addon)
	require.NoError(err)
}

func TestAddon_Update(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/addons/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"addon": {"name": "Internal Status Page", "id": "1"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	input := Addon{
		Name: "Internal Status Page",
	}

	addon, err := client.UpdateAddon("1", input)

	want := &Addon{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "Internal Status Page",
	}

	require.Equal(want, addon)
	require.NoError(err)
}

func TestAddon_Delete(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/addons/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	err := client.DeleteAddon("1")

	require.NoError(err)
}

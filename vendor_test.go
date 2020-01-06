package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// ListVendors
func TestVendor_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/vendors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"vendors": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListVendorOptions{
		APIListObject: listObj,
		Query:         "foo",
	}
	resp, err := client.ListVendors(opts)

	want := &ListVendorResponse{
		APIListObject: listObj,
		Vendors: []Vendor{
			{
				APIObject: APIObject{
					ID: "1",
				},
			},
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Get Vendor
func TestVendor_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/vendors/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"vendor": {"id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	venID := "1"

	resp, err := client.GetVendor(venID)

	want := &Vendor{
		APIObject: APIObject{
			ID: "1",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

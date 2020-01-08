package pagerduty

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ListServices
func TestService_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"services": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListServiceOptions{
		APIListObject: listObj,
		TeamIDs:       []string{},
		TimeZone:      "foo",
		SortBy:        "bar",
		Query:         "baz",
		Includes:      []string{},
	}
	resp, err := client.ListServices(opts)

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

	require.NoError(err)
	require.Equal(want, resp)
}

// Get Service
func TestService_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	id := "1"
	opts := &GetServiceOptions{
		Includes: []string{},
	}
	resp, err := client.GetService(id, opts)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Create Service
func TestService_Create(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := Service{
		Name: "foo",
	}
	resp, err := client.CreateService(input)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Update Service
func TestService_Update(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"service": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	input := Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	resp, err := client.UpdateService(input)

	want := &Service{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Delete Service
func TestService_Delete(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	err := client.DeleteService(id)

	require.NoError(err)
}

// Create Integration
func TestService_CreateIntegration(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services/1/integrations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var input = Integration{
		Name: "foo",
	}
	servID := "1"

	resp, err := client.CreateIntegration(servID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Get Integration
func TestService_GetIntegration(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var input = GetIntegrationOptions{
		Includes: []string{},
	}
	servID := "1"
	intID := "1"

	resp, err := client.GetIntegration(servID, intID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Update Integration
func TestService_UpdateIntegration(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var input = Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	servID := "1"

	resp, err := client.UpdateIntegration(servID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Delete Integration
func TestService_DeleteIntegration(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	servID := "1"
	intID := "1"
	err := client.DeleteIntegration(servID, intID)

	require.NoError(err)
}

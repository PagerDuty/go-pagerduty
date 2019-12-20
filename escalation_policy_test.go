package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEscalationPolicy_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/escalation_policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"escalation_policies": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var opts ListEscalationPoliciesOptions
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	resp, err := client.ListEscalationPolicies(opts)

	want := &ListEscalationPoliciesResponse{
		APIListObject: listObj,
		EscalationPolicies: []EscalationPolicy{
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

func TestEscalationPolicy_Create(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	input := EscalationPolicy{Name: "foo"}

	mux.HandleFunc("/escalation_policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"escalation_policy": {"name": "foo", "id": "1"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	resp, err := client.CreateEscalationPolicy(input)

	want := &EscalationPolicy{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

func TestEscalationPolicy_Delete(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	err := client.DeleteEscalationPolicy("1")
	require.NoError(err)
}

func TestEscalationPolicy_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"escalation_policy": {"id": "1"}}`))
	})
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts *GetEscalationPolicyOptions
	resp, err := client.GetEscalationPolicy("1", opts)

	want := &EscalationPolicy{
		APIObject: APIObject{
			ID: "1",
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

func TestEscalationPolicy_Update(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"escalation_policy": {"name": "foo", "id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := &EscalationPolicy{Name: "foo"}
	want := &EscalationPolicy{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	resp, err := client.UpdateEscalationPolicy("1", input)

	require.NoError(err)
	require.Equal(want, resp)
}

func TestEscalationPolicy_UpdateTeams(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	input := &EscalationPolicy{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
		Teams: []APIReference{},
	}

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"escalation_policy": {"name": "foo", "id": "1", "teams": []}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	resp, err := client.UpdateEscalationPolicy("1", input)

	want := &EscalationPolicy{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
		Teams: []APIReference{},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

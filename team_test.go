package pagerduty

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ListTeams
func TestTeam_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"teams": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListTeamOptions{
		APIListObject: listObj,
		Query:         "foo",
	}
	resp, err := client.ListTeams(opts)

	want := &ListTeamResponse{
		APIListObject: listObj,
		Teams: []Team{
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

// Create Team
func TestTeam_Create(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"team": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := &Team{
		Name: "foo",
	}
	resp, err := client.CreateTeam(input)

	want := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Delete Team
func TestTeam_Delete(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	err := client.DeleteTeam(id)

	require.NoError(err)
}

// Get Team
func TestTeam_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"team": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	resp, err := client.GetTeam(id)

	want := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Update Team
func TestTeam_Update(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"team": {"id": "1","name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	input := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	id := "1"
	resp, err := client.UpdateTeam(id, input)

	want := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Remove Escalation Policy from Team
func TestTeam_RemoveEscalationPolicyFromTeam(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams/1/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	teamID := "1"
	epID := "1"

	err := client.RemoveEscalationPolicyFromTeam(teamID, epID)

	require.NoError(err)
}

// Add Escalation Policy to Team
func TestTeam_AddEscalationPolicyToTeam(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams/1/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	teamID := "1"
	epID := "1"

	err := client.AddEscalationPolicyToTeam(teamID, epID)

	require.NoError(err)
}

// Remove User from Team
func TestTeam_RemoveUserFromTeam(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams/1/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	teamID := "1"
	userID := "1"

	err := client.RemoveUserFromTeam(teamID, userID)

	require.NoError(err)
}

// Add User to Team
func TestTeam_AddUserToTeam(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	mux.HandleFunc("/teams/1/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	teamID := "1"
	userID := "1"

	err := client.AddUserToTeam(teamID, userID)

	require.NoError(err)
}

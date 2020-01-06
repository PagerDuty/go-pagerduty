package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// ListUsers
func TestUser_List(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"users": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListUsersOptions{
		APIListObject: listObj,
		Query:         "foo",
		TeamIDs:       []string{},
		Includes:      []string{},
	}
	resp, err := client.ListUsers(opts)

	want := &ListUsersResponse{
		APIListObject: listObj,
		Users: []User{
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

// Create User
func TestUser_Create(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := User{
		Email: "foo@bar.com",
	}
	resp, err := client.CreateUser(input)

	want := &User{
		APIObject: APIObject{
			ID: "1",
		},
		Email: "foo@bar.com",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Delete User
func TestUser_Delete(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	err := client.DeleteUser(id)

	require.NoError(err)
}

// Get User
func TestUser_Get(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	userID := "1"
	opts := GetUserOptions{
		Includes: []string{},
	}
	resp, err := client.GetUser(userID, opts)

	want := &User{
		APIObject: APIObject{
			ID: "1",
		},
		Email: "foo@bar.com",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Update
func TestUser_Update(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := User{
		APIObject: APIObject{
			ID: "1",
		},
		Email: "foo@bar.com",
	}
	resp, err := client.UpdateUser(input)

	want := &User{
		APIObject: APIObject{
			ID: "1",
		},
		Email: "foo@bar.com",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// List User Contactmethods
func TestUser_ListContactMethods(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/users/1/contact_methods", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"contact_methods": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	ID := "1"

	resp, err := client.ListUserContactMethods(ID)

	want := &ListContactMethodsResponse{
		APIListObject: listObj,
		ContactMethods: []ContactMethod{
			{
				ID: "1",
			},
		},
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// Get user ContactMethod
func TestUser_GetContactMethod(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/users/1/contact_methods/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"contact_method": {"id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	methodID := "1"
	userID := "1"

	resp, err := client.GetUserContactMethod(userID, methodID)

	want := &ContactMethod{
		ID: "1",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

// TODO: Create user ContactMethod

// TODO: Delete User Contactmethod

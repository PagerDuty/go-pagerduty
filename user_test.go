package pagerduty

import (
	"net/http"
	"testing"
)

// ListUsers
func TestUser_List(t *testing.T) {
	setup()
	defer teardown()

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
	res, err := client.ListUsers(opts)

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

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create User
func TestUser_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := User{
		Email: "foo@bar.com",
	}
	res, err := client.CreateUser(input)

	want := &User{
		APIObject: APIObject{
			ID: "1",
		},
		Email: "foo@bar.com",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete User
func TestUser_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	err := client.DeleteUser(id)

	if err != nil {
		t.Fatal(err)
	}
}

// Get User
func TestUser_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	userID := "1"
	opts := GetUserOptions{
		Includes: []string{},
	}
	res, err := client.GetUser(userID, opts)

	want := &User{
		APIObject: APIObject{
			ID: "1",
		},
		Email: "foo@bar.com",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update
func TestUser_Update(t *testing.T) {
	setup()
	defer teardown()

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
	res, err := client.UpdateUser(input)

	want := &User{
		APIObject: APIObject{
			ID: "1",
		},
		Email: "foo@bar.com",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// List User Contactmethods
func TestUser_ListContactMethods(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/contact_methods", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"contact_methods": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	ID := "1"

	res, err := client.ListUserContactMethods(ID)

	want := &ListContactMethodsResponse{
		APIListObject: listObj,
		ContactMethods: []ContactMethod{
			{
				ID: "1",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Get user ContactMethod
func TestUser_GetContactMethod(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/contact_methods/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"contact_method": {"id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	methodID := "1"
	userID := "1"

	res, err := client.GetUserContactMethod(userID, methodID)

	want := &ContactMethod{
		ID: "1",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create user ContactMethod
func TestUser_CreateContactMethod(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/contact_methods", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"contact_method": {"id": "1", "type": "email_contact_method"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	userID := "1"
	contactMethod := ContactMethod{
		Type: "email_contact_method",
	}
	res, err := client.CreateUserContactMethod(userID, contactMethod)

	want := &ContactMethod{
		ID:   "1",
		Type: "email_contact_method",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete User Contactmethod
func TestUser_DeleteContactMethod(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/contact_methods/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	userID := "1"
	contactMethodID := "1"

	err := client.DeleteUserContactMethod(userID, contactMethodID)

	if err != nil {
		t.Fatal(err)
	}
}

// Update User ContactMethod
func TestUser_UpdateContactMethod(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/contact_methods/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"contact_method": {"id": "1", "type": "email_contact_method"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	userID := "1"
	contactMethod := ContactMethod{
		ID:   "1",
		Type: "email_contact_method",
	}
	res, err := client.UpdateUserContactMethod(userID, contactMethod)

	want := &ContactMethod{
		ID:   "1",
		Type: "email_contact_method",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

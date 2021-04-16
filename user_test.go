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
		_, _ = w.Write([]byte(`{"users": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListUsersOptions{
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
		_, _ = w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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

	client := defaultTestClient(server.URL, "foo")
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
		_, _ = w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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
		_, _ = w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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

// Get Current User
func TestUser_GetCurrent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := GetCurrentUserOptions{
		Includes: []string{},
	}
	res, err := client.GetCurrentUser(opts)

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
		_, _ = w.Write([]byte(`{"contact_methods": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
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
		_, _ = w.Write([]byte(`{"contact_method": {"id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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
		_, _ = w.Write([]byte(`{"contact_method": {"id": "1", "type": "email_contact_method"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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

	client := defaultTestClient(server.URL, "foo")
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
		_, _ = w.Write([]byte(`{"contact_method": {"id": "1", "type": "email_contact_method"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
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

// Get user NotificationRule
func TestUser_GetUserNotificationRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/notification_rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"notification_rule": {"id": "1", "start_delay_in_minutes": 1, "urgency": "low", "contact_method": {"id": "1"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	ruleID := "1"
	userID := "1"

	res, err := client.GetUserNotificationRule(userID, ruleID)

	want := &NotificationRule{
		ID:                  "1",
		StartDelayInMinutes: uint(1),
		Urgency:             "low",
		ContactMethod: ContactMethod{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create user NotificationRule
func TestUser_CreateUserNotificationRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/notification_rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"notification_rule": {"id": "1", "start_delay_in_minutes": 1, "urgency": "low", "contact_method": {"id": "1"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	userID := "1"
	rule := NotificationRule{
		Type: "email_contact_method",
	}
	res, err := client.CreateUserNotificationRule(userID, rule)

	want := &NotificationRule{
		ID:                  "1",
		StartDelayInMinutes: uint(1),
		Urgency:             "low",
		ContactMethod: ContactMethod{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// List User NotificationRules
func TestUser_ListUserNotificationRules(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/notification_rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"notification_rules": [{"id": "1", "start_delay_in_minutes": 1, "urgency": "low", "contact_method": {"id": "1"}}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	ID := "1"

	res, err := client.ListUserNotificationRules(ID)

	want := &ListUserNotificationRulesResponse{
		APIListObject: listObj,
		NotificationRules: []NotificationRule{
			{
				ID:                  "1",
				StartDelayInMinutes: uint(1),
				Urgency:             "low",
				ContactMethod: ContactMethod{
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

// Update user NotificationRule
func TestUser_UpdateUserNotificationRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/notification_rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"notification_rule": {"id": "1", "start_delay_in_minutes": 1, "urgency": "low", "contact_method": {"id": "1"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	userID := "1"
	rule := NotificationRule{
		ID:   "1",
		Type: "email_contact_method",
	}
	res, err := client.UpdateUserNotificationRule(userID, rule)

	want := &NotificationRule{
		ID:                  "1",
		StartDelayInMinutes: uint(1),
		Urgency:             "low",
		ContactMethod: ContactMethod{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete user NotificationRule
func TestUser_DeleteUserNotificationRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1/notification_rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	userID := "1"
	ruleID := "1"

	client := defaultTestClient(server.URL, "foo")
	if err := client.DeleteUserNotificationRule(userID, ruleID); err != nil {
		t.Fatal(err)
	}
}

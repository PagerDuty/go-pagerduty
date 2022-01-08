package pagerduty

import (
	"net/http"
	"testing"
)

// ListTags
func TestTag_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["query"], []string{"MyTag"})
		_, _ = w.Write([]byte(`{"tags": [{"id": "1","label":"MyTag"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListTagOptions{
		Limit:  listObj.Limit,
		Offset: listObj.Offset,
		Query:  "MyTag",
	}
	res, err := client.ListTags(opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListTagResponse{
		Tags: []*Tag{
			{
				APIObject: APIObject{
					ID: "1",
				},
				Label: "MyTag",
			},
		},
	}

	testEqual(t, want, res)
}

// Create Tag
func TestTag_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"tag": {"id": "1","Label":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &Tag{
		Label: "foo",
	}
	res, err := client.CreateTag(input)

	want := &Tag{
		APIObject: APIObject{
			ID: "1",
		},
		Label: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete Tag
func TestTag_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"
	err := client.DeleteTag(id)
	if err != nil {
		t.Fatal(err)
	}
}

// Get Tag
func TestTag_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"tag": {"id": "1","label":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"
	res, err := client.GetTag(id)

	want := &Tag{
		APIObject: APIObject{
			ID: "1",
		},
		Label: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Assign Tags - Add
func TestTag_AssignAdd(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/change_tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	client := defaultTestClient(server.URL, "foo")
	ta := &TagAssignments{
		Add: []*TagAssignment{
			{
				Type:  "tag_reference",
				TagID: "1",
			},
			{
				Type:  "tag",
				Label: "NewTag",
			},
		},
	}
	// this endpoint only returns  an "ok" in the body. no point in testing for it.
	err := client.AssignTags("teams", "1", ta)
	if err != nil {
		t.Fatal(err)
	}
}

// Assign Tags - Remove
func TestTag_AssignRemove(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/change_tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	client := defaultTestClient(server.URL, "foo")
	ta := &TagAssignments{
		Remove: []*TagAssignment{
			{
				Type:  "tag_reference",
				TagID: "1",
			},
		},
	}
	// this endpoint only returns  an "ok" in the body. no point in testing for it.
	err := client.AssignTags("teams", "1", ta)
	if err != nil {
		t.Fatal(err)
	}
}

// GetUsersByTag
func TestTag_GetUsersByTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags/1/users/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"users": [{"id": "1"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	tid := "1"

	res, err := client.GetUsersByTag(tid)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListUserResponse{
		Users: []*APIObject{
			{
				ID: "1",
			},
		},
	}

	testEqual(t, want, res)
}

// GetTeamsByTag
func TestTag_GetTeamsByTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags/1/teams/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"teams": [{"id": "1"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	tid := "1"

	res, err := client.GetTeamsByTag(tid)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListTeamsForTagResponse{
		Teams: []*APIObject{
			{
				ID: "1",
			},
		},
	}

	testEqual(t, want, res)
}

// GetEscalationPoliciesByTag
func TestTag_GetEscalationPoliciesByTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags/1/escalation_policies/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"escalation_policies": [{"id": "1"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	tid := "1"

	res, err := client.GetEscalationPoliciesByTag(tid)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListEPResponse{
		EscalationPolicies: []*APIObject{
			{
				ID: "1",
			},
		},
	}

	testEqual(t, want, res)
}

// GetTagsForEntity
func TestTag_GetTagsForEntity(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/escalation_policies/1/tags/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"tags": [{"id": "1"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	eid := "1"
	e := "escalation_policies"
	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}

	opts := ListTagOptions{
		Limit:  listObj.Limit,
		Offset: listObj.Offset,
	}
	res, err := client.GetTagsForEntity(e, eid, opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListTagResponse{
		Tags: []*Tag{
			{
				APIObject: APIObject{
					ID: "1",
				},
			},
		},
	}

	testEqual(t, want, res)
}

package pagerduty

import (
	"net/http"
	"testing"
)

func TestEscalationPolicy_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/escalation_policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"escalation_policies": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var opts ListEscalationPoliciesOptions
	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListEscalationPolicies(opts)

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
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestEscalationPolicy_Create(t *testing.T) {
	setup()
	defer teardown()

	input := EscalationPolicy{Name: "foo"}

	mux.HandleFunc("/escalation_policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"escalation_policy": {"name": "foo", "id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	res, err := client.CreateEscalationPolicy(input)

	want := &EscalationPolicy{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestEscalationPolicy_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteEscalationPolicy("1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEscalationPolicy_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"escalation_policy": {"id": "1"}}`))
	})
	client := defaultTestClient(server.URL, "foo")
	var opts *GetEscalationPolicyOptions
	res, err := client.GetEscalationPolicy("1", opts)

	want := &EscalationPolicy{
		APIObject: APIObject{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestEscalationPolicy_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"escalation_policy": {"name": "foo", "id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &EscalationPolicy{Name: "foo"}
	want := &EscalationPolicy{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	res, err := client.UpdateEscalationPolicy("1", input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestEscalationPolicy_UpdateTeams(t *testing.T) {
	setup()
	defer teardown()

	input := &EscalationPolicy{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
		Teams: []APIReference{},
	}

	mux.HandleFunc("/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"escalation_policy": {"name": "foo", "id": "1", "teams": []}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateEscalationPolicy("1", input)

	want := &EscalationPolicy{
		Name: "foo",
		APIObject: APIObject{
			ID: "1",
		},
		Teams: []APIReference{},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

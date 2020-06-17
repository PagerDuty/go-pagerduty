package pagerduty

import (
	"net/http"
	"testing"
)

// List Rulesets
func TestRuleset_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"rulesets": [{"id": "1"}]}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	res, err := client.ListRulesets()
	if err != nil {
		t.Fatal(err)
	}
	want := &ListRulesetsResponse{
		Rulesets: []*Ruleset{
			{
				ID: "1",
			},
		},
	}

	testEqual(t, want, res)
}

// List Rulesets
func TestRuleset_List_WithOpts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, "limit=5&offset=3", r.URL.RawQuery)
		w.Write([]byte(`{"rulesets": [{"id": "1"}]}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	o := APIListObject{Limit: 5, Offset: 3}
	res, err := client.ListRulesetsWithOpts(o)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListRulesetsResponse{
		Rulesets: []*Ruleset{
			{
				ID: "1",
			},
		},
	}

	testEqual(t, want, res)
}

// Create Ruleset
func TestRuleset_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"ruleset": {"id": "1", "name": "foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := &Ruleset{
		Name: "foo",
	}
	res, _, err := client.CreateRuleset(input)

	want := &Ruleset{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Get Ruleset
func TestRuleset_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"ruleset": {"id": "1", "name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	ruleSetID := "1"

	res, _, err := client.GetRuleset(ruleSetID)

	want := &Ruleset{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update Ruleset
func TestRuleset_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"ruleset": {"id": "1", "name":"foo"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	input := &Ruleset{
		ID:   "1",
		Name: "foo",
	}
	res, _, err := client.UpdateRuleset(input)

	want := &Ruleset{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete Ruleset
func TestRuleset_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	id := "1"
	err := client.DeleteRuleset(id)

	if err != nil {
		t.Fatal(err)
	}
}

// List Ruleset Rules
func TestRuleset_ListRules(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1/rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"rules": [{"id": "1"}]}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	rulesetID := "1"
	res, err := client.ListRulesetRules(rulesetID)
	if err != nil {
		t.Fatal(err)
	}

	want := &ListRulesetRulesResponse{
		Rules: []*RulesetRule{
			{
				ID: "1",
			},
		},
	}
	testEqual(t, want, res)
}

// Get Ruleset Rule
func TestRuleset_GetRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1/rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"rule": {"id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	rulesetID := "1"
	ruleID := "1"
	res, _, err := client.GetRulesetRule(rulesetID, ruleID)
	if err != nil {
		t.Fatal(err)
	}

	want := &RulesetRule{
		ID: "1",
	}
	testEqual(t, want, res)
}

// Create Ruleset Rule
func TestRuleset_CreateRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1/rules/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Write([]byte(`{"rule": {"id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	rulesetID := "1"
	rule := &RulesetRule{}

	res, _, err := client.CreateRulesetRule(rulesetID, rule)
	if err != nil {
		t.Fatal(err)
	}

	want := &RulesetRule{
		ID: "1",
	}
	testEqual(t, want, res)
}

// Update Ruleset Rule
func TestRuleset_UpdateRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1/rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"rule": {"id": "1"}}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	rulesetID := "1"
	ruleID := "1"
	rule := &RulesetRule{}

	res, _, err := client.UpdateRulesetRule(rulesetID, ruleID, rule)
	if err != nil {
		t.Fatal(err)
	}

	want := &RulesetRule{
		ID: "1",
	}
	testEqual(t, want, res)
}

// Delete Ruleset Rule
func TestRuleset_DeleteRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rulesets/1/rules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	ruleID := "1"
	rulesetID := "1"

	err := client.DeleteRulesetRule(rulesetID, ruleID)

	if err != nil {
		t.Fatal(err)
	}
}

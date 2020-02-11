package pagerduty

import (
	"net/http"
	"testing"
)

func TestAbility_ListAbilities(t *testing.T) {
	t.Parallel()
	setup()
	defer teardown()

	mux.HandleFunc("/abilities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"abilities": ["sso"]}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	want := &ListAbilityResponse{Abilities: []string{"sso"}}

	res, err := client.ListAbilities()

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestAbility_ListAbilitiesFailure(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/abilities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusForbidden)
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	if _, err := client.ListAbilities(); err == nil {
		t.Fatal("expected error; got nil")
	}
}

func TestAbility_TestAbility(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/abilities/sso", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNoContent)
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	if err := client.TestAbility("sso"); err != nil {
		t.Fatal(err)
	}
}

func TestAbility_TestAbilityFailure(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/abilities/sso", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusForbidden)
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	if err := client.TestAbility("sso"); err == nil {
		t.Fatal("expected error; got nil")
	}
}

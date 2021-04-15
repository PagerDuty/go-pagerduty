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
		_, _ = w.Write([]byte(`{"abilities": ["sso"]}`))
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
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

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

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

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

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

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	if err := client.TestAbility("sso"); err == nil {
		t.Fatal("expected error; got nil")
	}
}

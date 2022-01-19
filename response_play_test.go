package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// ListResponsePlays
func TestResponsePlay_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/response_plays", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"response_plays": [{"id": "1"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListResponsePlays(context.TODO(), ListResponsePlaysOptions{})

	want := []ResponsePlay{
		{
			ID: "1",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Get ResponsePlay
func TestResponsePlay_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/response_plays/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"response_play": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	id := "1"
	res, err := client.GetResponsePlay(context.TODO(), id)

	want := ResponsePlay{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create ResponsePlay
func TestResponsePlay_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/response_plays", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"response_play": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := ResponsePlay{
		Name: "foo",
	}
	res, err := client.CreateResponsePlay(context.TODO(), input)

	want := ResponsePlay{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update ResponsePlay
func TestResponsePlay_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/response_plays/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"response_play": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := ResponsePlay{
		ID:   "1",
		Name: "foo",
	}
	res, err := client.UpdateResponsePlay(context.TODO(), input)

	want := ResponsePlay{
		ID:   "1",
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete ResponsePlay
func TestResponsePlay_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/response_plays/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"
	err := client.DeleteResponsePlay(context.TODO(), id)
	if err != nil {
		t.Fatal(err)
	}
}

// Run ResponsePlay
func TestResponsePlay_Run(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/response_plays/1/run", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"incident": {"id": "5","type":"incident_reference"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.RunResponsePlay(context.TODO(), "foo@example.com", "1", "5")
	if err != nil {
		t.Fatal(err)
	}
}

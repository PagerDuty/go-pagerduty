package pagerduty

import (
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	expectedChangeCreatePayload = `{"routing_key":"a0000000aa0000a0a000aa0a0a0aa000","payload":{"source":"Test runner",` +
		`"summary":"Summary can't be blank","timestamp":"2020-10-19T03:06:16.318Z",` +
		`"custom_details":{"DetailKey1":"DetailValue1","DetailKey2":"DetailValue2"}},` +
		`"links":[{"href":"https://acme.pagerduty.dev/build/2","text":"View more details in Acme!"},` +
		`{"href":"https://acme2.pagerduty.dev/build/2","text":"View more details in Acme2!"}]}`
)

func TestChangeEvent_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(
		"/v2/change/enqueue", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"message": "Change event processed", "status": "success"}`))
		},
	)

	client := defaultTestClient(server.URL, "foo")

	want := ChangeEventResponse{
		Status:  "success",
		Message: "Change event processed",
	}

	eventDetails := map[string]interface{}{"DetailKey1": "DetailValue1", "DetailKey2": "DetailValue2"}
	ce := ChangeEvent{
		RoutingKey: "a0000000aa0000a0a000aa0a0a0aa000",
		Payload: ChangeEventPayload{
			Source:        "Test runner",
			Summary:       "Summary can't be blank",
			Timestamp:     "2020-10-19T03:06:16.318Z",
			CustomDetails: eventDetails,
		},
		Links: []ChangeEventLink{
			{
				Href: "https://acme.pagerduty.dev/build/2",
				Text: "View more details in Acme!",
			},
			{
				Href: "https://acme2.pagerduty.dev/build/2",
				Text: "View more details in Acme2!",
			},
		},
	}

	res, err := client.CreateChangeEvent(ce)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, *res)
}

func TestChangeEvent_CreateWithPayloadVerification(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(
		"/v2/change/enqueue", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			testEqual(t, expectedChangeCreatePayload, string(body))
		},
	)

	client := defaultTestClient(server.URL, "foo")

	eventDetails := map[string]interface{}{"DetailKey1": "DetailValue1", "DetailKey2": "DetailValue2"}
	ce := ChangeEvent{
		RoutingKey: "a0000000aa0000a0a000aa0a0a0aa000",
		Payload: ChangeEventPayload{
			Source:        "Test runner",
			Summary:       "Summary can't be blank",
			Timestamp:     "2020-10-19T03:06:16.318Z",
			CustomDetails: eventDetails,
		},
		Links: []ChangeEventLink{
			{
				Href: "https://acme.pagerduty.dev/build/2",
				Text: "View more details in Acme!",
			},
			{
				Href: "https://acme2.pagerduty.dev/build/2",
				Text: "View more details in Acme2!",
			},
		},
	}

	_, _ = client.CreateChangeEvent(ce)
}

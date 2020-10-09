package pagerduty

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var reqBody = `{"routing_key":"a0000000aa0000a0a000aa0a0a0aa000","payload":{"source":"Test runner",` +
	`"summary":"Summary can't be blank","time":"2020-10-09T18:53:40.635779-04:00",` +
	`"timestamp":"2020-10-09T18:53:40-04:00","custom_details":{"DetailKey1":"DetailValue1",` +
	`"DetailKey2":"DetailValue2"}},"links":[{"href":"https://acme.pagerduty.dev/build/2",` +
	`"text":"View more details in Acme!"},{"href":"https://acme2.pagerduty.dev/build/2",` +
	`"text":"View more details in Acme2!"}]}`

func TestClient_SendChangeEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/change/enqueue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "Change event processed", "status": "success"}`))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	want := ChangeEventResponse{
		Status:  "success",
		Message: "Change event processed",
	}

	eventDetails := map[string]string{"DetailKey1": "DetailValue1", "DetailKey2": "DetailValue2"}
	ce := ChangeEvent{
		RoutingKey: "a0000000aa0000a0a000aa0a0a0aa000",
		Payload:    Payload{Source: "Test runner", Summary: "Summary can't be blank", Timestamp: time.Now(), CustomDetails: eventDetails},
		Links: []Link{
			{
				Href: "https://acme.pagerduty.dev/build/2",
				Text: "View more details in Acme!",
			},
			{
				Href: "https://acme2.pagerduty.dev/build/2",
				Text: "View more details in Acme2!",
			}},
	}

	res, err := client.SendChangeEvent(ce)

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, *res)
}

func TestClient_SendChangeEventContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/change/enqueue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		testEqual(t, reqBody, string(body))
	})

	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	eventDetails := map[string]string{"DetailKey1": "DetailValue1", "DetailKey2": "DetailValue2"}
	ts, _ := time.Parse(time.RFC3339, "2020-10-09T18:53:40.635779-04:00")
	ce := ChangeEvent{
		RoutingKey: "a0000000aa0000a0a000aa0a0a0aa000",
		Payload: Payload{Source: "Test runner", Summary: "Summary can't be blank",
			Timestamp: ts, CustomDetails: eventDetails},
		Links: []Link{
			{
				Href: "https://acme.pagerduty.dev/build/2",
				Text: "View more details in Acme!",
			},
			{
				Href: "https://acme2.pagerduty.dev/build/2",
				Text: "View more details in Acme2!",
			}},
	}

	_, _ = client.SendChangeEvent(ce)

}

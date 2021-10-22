package pagerduty

import (
	"net/http"
	"testing"
)

func TestEventV2_ManageEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/enqueue", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"status": "ok", "dedup_key": "yes", "message": "ok"}`))
	})

	client := defaultTestClient(server.URL, "foo")
	evt := &V2Event{
		RoutingKey: "abc123",
	}

	res, err := client.ManageEvent(evt)
	if err != nil {
		t.Fatal(err)
	}

	want := &V2EventResponse{
		Status:   "ok",
		DedupKey: "yes",
		Message:  "ok",
	}

	testEqual(t, want, res)
}

func TestEventsAPIV2Error_BadRequest(t *testing.T) {
	tests := []struct {
		name string
		e    EventsAPIV2Error
		want bool
	}{
		{
			name: "bad_request",
			e: EventsAPIV2Error{
				StatusCode: http.StatusBadRequest,
				EventsAPIV2Error: NullEventsAPIV2ErrorObject{
					Valid: true,
					ErrorObject: EventsAPIV2ErrorObject{
						Status:  "invalid",
						Message: "Event object is invalid",
						Errors:  []string{"Length of 'routing_key' is incorrect (should be 32 characters)", "'event_action' is missing or blank"},
					},
				},
			},
			want: true,
		},
		{
			name: "rate_limited",
			e: EventsAPIV2Error{
				StatusCode: http.StatusTooManyRequests,
				EventsAPIV2Error: NullEventsAPIV2ErrorObject{
					Valid: true,
					ErrorObject: EventsAPIV2ErrorObject{
						Status:  "throttle exceeded",
						Message: "Requests for this service are arriving too quickly.  Please retry later.",
						Errors:  []string{"Enhance Your Calm."},
					},
				},
			},
			want: false,
		},
		{
			name: "InternalServerError",
			e: EventsAPIV2Error{
				StatusCode: http.StatusInternalServerError,
			},
			want: false,
		},
		{
			name: "ServiceUnavailable",
			e: EventsAPIV2Error{
				StatusCode: http.StatusServiceUnavailable,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.BadRequest(); got != tt.want {
				t.Fatalf("tt.e.BadRequest() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestEventsAPIV2Error_RateLimited(t *testing.T) {
	tests := []struct {
		name string
		e    EventsAPIV2Error
		want bool
	}{
		{
			name: "rate_limited",
			e: EventsAPIV2Error{
				StatusCode: http.StatusTooManyRequests,
				EventsAPIV2Error: NullEventsAPIV2ErrorObject{
					Valid: true,
					ErrorObject: EventsAPIV2ErrorObject{
						Status:  "throttle exceeded",
						Message: "Requests for this service are arriving too quickly.  Please retry later.",
						Errors:  []string{"Enhance Your Calm"},
					},
				},
			},
			want: true,
		},
		{
			name: "not_found",
			e: EventsAPIV2Error{
				StatusCode: http.StatusNotFound,
				EventsAPIV2Error: NullEventsAPIV2ErrorObject{
					Valid: true,
					ErrorObject: EventsAPIV2ErrorObject{
						Status:  "Not Found",
						Message: "Not Found",
						Errors:  []string{"Not Found"},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.RateLimited(); got != tt.want {
				t.Fatalf("tt.e.RateLimited() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestEventsAPIV2Error_Temporary(t *testing.T) {
	tests := []struct {
		name string
		e    EventsAPIV2Error
		want bool
	}{
		{
			name: "rate_limited",
			e: EventsAPIV2Error{
				StatusCode: http.StatusTooManyRequests,
				EventsAPIV2Error: NullEventsAPIV2ErrorObject{
					Valid: true,
					ErrorObject: EventsAPIV2ErrorObject{
						Status:  "throttle exceeded",
						Message: "Requests for this service are arriving too quickly.  Please retry later.",
						Errors:  []string{"Enhance Your Calm"},
					},
				},
			},
			want: true,
		},
		{
			name: "InternalServerError",
			e: EventsAPIV2Error{
				StatusCode: http.StatusInternalServerError,
			},
			want: true,
		},
		{
			name: "ServiceUnavailable",
			e: EventsAPIV2Error{
				StatusCode: http.StatusServiceUnavailable,
			},
			want: true,
		},
		{
			name: "not_found",
			e: EventsAPIV2Error{
				StatusCode: http.StatusNotFound,
				EventsAPIV2Error: NullEventsAPIV2ErrorObject{
					Valid: true,
					ErrorObject: EventsAPIV2ErrorObject{
						Status:  "Not Found",
						Message: "Not Found",
						Errors:  []string{"Not Found"},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Temporary(); got != tt.want {
				t.Fatalf("tt.e.Temporary() = %t, want %t", got, tt.want)
			}
		})
	}
}

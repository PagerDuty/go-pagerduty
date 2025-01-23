package pagerduty

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
				APIError: NullEventsAPIV2ErrorObject{
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
				APIError: NullEventsAPIV2ErrorObject{
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
		{
			name: "RequestTimeout",
			e: EventsAPIV2Error{
				StatusCode: http.StatusRequestTimeout,
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

func TestEventsAPIV2Error_APITimeout(t *testing.T) {
	tests := []struct {
		name string
		e    EventsAPIV2Error
		want bool
	}{
		{
			name: "bad_request",
			e: EventsAPIV2Error{
				StatusCode: http.StatusBadRequest,
				APIError: NullEventsAPIV2ErrorObject{
					Valid: true,
					ErrorObject: EventsAPIV2ErrorObject{
						Status:  "invalid",
						Message: "Event object is invalid",
						Errors:  []string{"Length of 'routing_key' is incorrect (should be 32 characters)", "'event_action' is missing or blank"},
					},
				},
			},
			want: false,
		},
		{
			name: "rate_limited",
			e: EventsAPIV2Error{
				StatusCode: http.StatusTooManyRequests,
				APIError: NullEventsAPIV2ErrorObject{
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
		{
			name: "RequestTimeout",
			e: EventsAPIV2Error{
				StatusCode: http.StatusRequestTimeout,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.APITimeout(); got != tt.want {
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
				APIError: NullEventsAPIV2ErrorObject{
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
				APIError: NullEventsAPIV2ErrorObject{
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
				APIError: NullEventsAPIV2ErrorObject{
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
			name: "RequestTimeout",
			e: EventsAPIV2Error{
				StatusCode: http.StatusRequestTimeout,
			},
			want: true,
		},
		{
			name: "not_found",
			e: EventsAPIV2Error{
				StatusCode: http.StatusNotFound,
				APIError: NullEventsAPIV2ErrorObject{
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

func TestEventsAPIV2Error_Error(t *testing.T) {
	t.Run("json_tests", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
		}{
			{
				name:  "no_error",
				input: `{"status": "Calming", "message": "Enhance Your Calm"}`,
				want:  "HTTP response failed with status code 429, status: Calming, message: Enhance Your Calm",
			},
			{
				name:  "one_error",
				input: `{"message": "Enhance Your Calm", "status": "Calming", "errors": ["No Seriously, Enhance Your Calm"]}`,
				want:  "HTTP response failed with status code 429, status: Calming, message: Enhance Your Calm: No Seriously, Enhance Your Calm",
			},
			{
				name:  "two_error",
				input: `{"message": "Enhance Your Calm", "status": "Calming", "errors":["No Seriously, Enhance Your Calm", "Slow Your Roll"]}`,
				want:  "HTTP response failed with status code 429, status: Calming, message: Enhance Your Calm: No Seriously, Enhance Your Calm (and 1 more error...)",
			},
			{
				name:  "three_error",
				input: `{"message": "Enhance Your Calm", "status": "Calming", "errors":["No Seriously, Enhance Your Calm", "Slow Your Roll", "No, really..."]}`,
				want:  "HTTP response failed with status code 429, status: Calming, message: Enhance Your Calm: No Seriously, Enhance Your Calm (and 2 more errors...)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var a EventsAPIV2Error

				if err := json.Unmarshal([]byte(tt.input), &a); err != nil {
					t.Fatalf("failed to unmarshal JSON: %s", err)
				}

				a.StatusCode = 429

				if got := a.Error(); got != tt.want {
					t.Errorf("a.Error() = %q, want %q", got, tt.want)
				}
			})
		}
	})

	tests := []struct {
		name string
		a    EventsAPIV2Error
		want string
	}{
		{
			name: "message",
			a: EventsAPIV2Error{
				message: "test message",
			},
			want: "test message",
		},
		{
			name: "APIError_nil",
			a: EventsAPIV2Error{
				StatusCode: http.StatusServiceUnavailable,
			},
			want: "HTTP response failed with status code 503 and no JSON error object was present",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Error(); got != tt.want {
				fmt.Println(got)
				fmt.Println(tt.want)
				t.Fatalf("tt.a.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEventsAPIV2Error_UnmarshalJSON(t *testing.T) {
	v := &EventsAPIV2Error{}
	if err := v.UnmarshalJSON([]byte(`{`)); !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Error("exepcted error not seen")
	}
}

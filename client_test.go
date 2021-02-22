package pagerduty

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the PagerDuty client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	authToken := "foo"
	client = NewClient(authToken)
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("returned %#v; want %#v", expected, actual)
	}
}

func TestGetBasePrefix(t *testing.T) {
	var testTable = []struct {
		in  string
		out string
	}{
		{"base.com/noparams", "base.com/noparams?"},
		{"base.com/?/noparams", "base.com/?/noparams?"},
		{"base.com/params?value=1", "base.com/params?value=1&"},
		{"base.com/?/params?value=1", "base.com/?/params?value=1&"},
		{"noslashes", "noslashes?"}, // this is what it will do... tbd what it should actually do
	}
	for _, tt := range testTable {
		s := getBasePrefix(tt.in)
		if s != tt.out {
			t.Errorf("got %q, want %q", s, tt.out)
		}
	}
}

func TestAPIError_Error(t *testing.T) {
	const jsonBody = `{"error":{"code": 420, "message": "Enhance Your Calm", "errors":["Enhance Your Calm", "Slow Your Roll"]}}`

	var a APIError

	if err := json.Unmarshal([]byte(jsonBody), &a); err != nil {
		t.Fatalf("failed to unmarshal JSON: %s", err)
	}

	a.StatusCode = 429

	const want = "HTTP response failed with status code 429, message: Enhance Your Calm (code: 420)"

	if got := a.Error(); got != want {
		t.Errorf("a.Error() = %q, want %q", got, want)
	}

	tests := []struct {
		name string
		a    APIError
		want string
	}{
		{
			name: "message",
			a: APIError{
				message: "test message",
			},
			want: "test message",
		},
		{
			name: "APIError_nil",
			a: APIError{
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

func TestAPIError_RateLimited(t *testing.T) {
	tests := []struct {
		name string
		a    APIError
		want bool
	}{
		{
			name: "rate_limited",
			a: APIError{
				StatusCode: http.StatusTooManyRequests,
				APIError: NullAPIErrorObject{
					Valid: true,
					ErrorObject: APIErrorObject{
						Code:    420,
						Message: "Enhance Your Calm",
						Errors:  []string{"Enhance Your Calm"},
					},
				},
			},
			want: true,
		},
		{
			name: "not_found",
			a: APIError{
				StatusCode: http.StatusNotFound,
				APIError: NullAPIErrorObject{
					Valid: true,
					ErrorObject: APIErrorObject{
						Code:    2100,
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
			if got := tt.a.RateLimited(); got != tt.want {
				t.Fatalf("tt.a.RateLimited() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestAPIError_Temporary(t *testing.T) {
	tests := []struct {
		name string
		a    APIError
		want bool
	}{
		{
			name: "rate_limited",
			a: APIError{
				StatusCode: http.StatusTooManyRequests,
				APIError: NullAPIErrorObject{
					Valid: true,
					ErrorObject: APIErrorObject{
						Code:    420,
						Message: "Enhance Your Calm",
						Errors:  []string{"Enhance Your Calm"},
					},
				},
			},
			want: true,
		},
		{
			name: "not_found",
			a: APIError{
				StatusCode: http.StatusNotFound,
				APIError: NullAPIErrorObject{
					Valid: true,
					ErrorObject: APIErrorObject{
						Code:    2100,
						Message: "Not Found",
						Errors:  []string{"Not Found"},
					},
				},
			},
			want: false,
		},
		{
			name: "InternalServerError",
			a: APIError{
				StatusCode: http.StatusInternalServerError,
			},
			want: true,
		},
		{
			name: "ServiceUnavailable",
			a: APIError{
				StatusCode: http.StatusServiceUnavailable,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Temporary(); got != tt.want {
				t.Fatalf("tt.a.Temporary() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestAPIError_NotFound(t *testing.T) {
	tests := []struct {
		name string
		a    APIError
		want bool
	}{
		{
			name: "rate_limited",
			a: APIError{
				StatusCode: http.StatusTooManyRequests,
				APIError: NullAPIErrorObject{
					Valid: true,
					ErrorObject: APIErrorObject{
						Code:    420,
						Message: "Enhance Your Calm",
						Errors:  []string{"Enhance Your Calm"},
					},
				},
			},
			want: false,
		},
		{
			name: "not_found",
			a: APIError{
				StatusCode: http.StatusNotFound,
				APIError: NullAPIErrorObject{
					Valid: true,
					ErrorObject: APIErrorObject{
						Code:    2100,
						Message: "Not Found",
						Errors:  []string{"Not Found"},
					},
				},
			},
			want: true,
		},
		{
			name: "not_found_weird_status",
			a: APIError{
				StatusCode: http.StatusBadRequest,
				APIError: NullAPIErrorObject{
					Valid: true,
					ErrorObject: APIErrorObject{
						Code:    2100,
						Message: "Not Found",
						Errors:  []string{"Not Found"},
					},
				},
			},
			want: true,
		},
		{
			name: "not_found_weird_error_code",
			a: APIError{
				StatusCode: http.StatusNotFound,
				APIError: NullAPIErrorObject{
					ErrorObject: APIErrorObject{
						Code:    2101,
						Message: "Not Found",
						Errors:  []string{"Not Found"},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.NotFound(); got != tt.want {
				t.Fatalf("tt.a.NotFound() = %t, want %t", got, tt.want)
			}
		})
	}
}

package pagerduty

import (
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
	var authToken = "foo"
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

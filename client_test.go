package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func defaultTestClient(serverURL, authToken string) *Client {
	return &Client{
		v2EventsAPIEndpoint: serverURL,
		apiEndpoint:         serverURL,
		authToken:           authToken,
		HTTPClient:          defaultHTTPClient,
		debugFlag:           new(uint64),
		lastRequest:         &atomic.Value{},
		lastResponse:        &atomic.Value{},
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()

	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testEqual(t *testing.T, want interface{}, got interface{}) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("values not equal (-want / +got):\n%s", diff)
	}
}

// testErrCheck looks to see if errContains is a substring of err.Error(). If
// not, this calls t.Fatal(). It also calls t.Fatal() if there was an error, but
// errContains is empty. Returns true if you should continue running the test,
// or false if you should stop the test.
func testErrCheck(t *testing.T, name string, errContains string, err error) bool {
	t.Helper()

	if len(errContains) > 0 {
		if err == nil {
			t.Fatalf("%s error = <nil>, should contain %q", name, errContains)
			return false
		}

		if errStr := err.Error(); !strings.Contains(errStr, errContains) {
			t.Fatalf("%s error = %q, should contain %q", name, errStr, errContains)
			return false
		}

		return false
	}

	if err != nil && len(errContains) == 0 {
		t.Fatalf("%s unexpected error: %v", name, err)
		return false
	}

	return true
}

func TestGetBasePrefix(t *testing.T) {
	testTable := []struct {
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
	t.Run("json_tests", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
		}{
			{
				name:  "one_error",
				input: `{"error":{"code": 420, "message": "Enhance Your Calm", "errors":["No Seriously, Enhance Your Calm"]}}`,
				want:  "HTTP response failed with status code 429, message: Enhance Your Calm (code: 420): No Seriously, Enhance Your Calm",
			},
			{
				name:  "two_error",
				input: `{"error":{"code": 420, "message": "Enhance Your Calm", "errors":["No Seriously, Enhance Your Calm", "Slow Your Roll"]}}`,
				want:  "HTTP response failed with status code 429, message: Enhance Your Calm (code: 420): No Seriously, Enhance Your Calm (and 1 more error...)",
			},
			{
				name:  "three_error",
				input: `{"error":{"code": 420, "message": "Enhance Your Calm", "errors":["No Seriously, Enhance Your Calm", "Slow Your Roll", "No, really..."]}}`,
				want:  "HTTP response failed with status code 429, message: Enhance Your Calm (code: 420): No Seriously, Enhance Your Calm (and 2 more errors...)",
			},
			{
				name:  "issue_478",
				input: `{"error":["links should have at most 50 item(s)"]}`,
				want:  "HTTP response failed with status code 429, message: none (code: 0): links should have at most 50 item(s)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var a APIError
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

func TestClient_SetDebugFlag(t *testing.T) {
	c := defaultTestClient("", "")
	c.SetDebugFlag(42)

	tests := []struct {
		name string
		flag DebugFlag
	}{
		{
			name: "zero_flag",
		},

		{
			name: "capture_response_flag",
			flag: DebugCaptureLastResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.SetDebugFlag(tt.flag)

			got := atomic.LoadUint64(c.debugFlag)
			if got != uint64(tt.flag) {
				t.Fatalf("got = %64b, want = %64b", got, tt.flag)
			}
		})
	}
}

func TestClient_LastAPIRequest(t *testing.T) {
	t.Run("unit", func(t *testing.T) {
		c := defaultTestClient("", "")
		got, ok := c.LastAPIRequest()
		if ok {
			t.Fatal("new client ok = true, want false")
		}

		if got != nil {
			t.Fatal("got != nil")
		}

		tests := []struct {
			name string
			req  *http.Request
			ok   bool
		}{
			{
				name: "nil_response",
			},
			{
				name: "non-nil_response",
				req:  &http.Request{},
				ok:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c.lastRequest.Store(tt.req)

				got, ok = c.LastAPIRequest()
				if ok != tt.ok {
					t.Fatalf("ok = %t, want %t", ok, tt.ok)
				}

				if !ok {
					if got != nil {
						t.Fatal("got != nil")
					}
					return
				}

				if got != tt.req {
					t.Fatalf("got = %v, want %v", got, tt.req)
				}
			})
		}
	})

	t.Run("integration", func(t *testing.T) {
		const requestBody = `{"user":{"id":"1","name":"","email":"foo@bar.com"}}`

		setup()
		defer teardown()

		mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPut)
			_, _ = w.Write([]byte(`{"user": {"id": "1", "email":"foo@bar.com"}}`))
		})

		c := defaultTestClient(server.URL, "foo")
		c.SetDebugFlag(DebugCaptureLastRequest)

		_, err := c.UpdateUserWithContext(context.Background(), User{
			APIObject: APIObject{
				ID: "1",
			},
			Email: "foo@bar.com",
		})
		testErrCheck(t, "c.UpdateUserWithContext()", "", err)

		got, ok := c.LastAPIRequest()
		if !ok {
			t.Fatal("ok = false, want true")
		}

		if got == nil {
			t.Fatal("got == nil")
		}

		if got.Method != http.MethodPut {
			t.Fatalf("got.Method = %s, want %s", got.Method, http.MethodPut)
		}

		body, err := ioutil.ReadAll(got.Body)
		testErrCheck(t, "ioutil.ReadAll()", "", err)

		if jb := string(body); jb != requestBody {
			t.Fatalf("got.Body = %q, want %q", jb, requestBody)
		}
	})
}

func TestClient_LastAPIResponse(t *testing.T) {
	t.Run("unit", func(t *testing.T) {
		c := defaultTestClient("", "")
		got, ok := c.LastAPIResponse()
		if ok {
			t.Fatal("new client ok = true, want false")
		}

		if got != nil {
			t.Fatal("got != nil")
		}

		tests := []struct {
			name string
			resp *http.Response
			ok   bool
		}{
			{
				name: "nil_response",
			},
			{
				name: "non-nil_response",
				resp: &http.Response{},
				ok:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c.lastResponse.Store(tt.resp)

				got, ok = c.LastAPIResponse()
				if ok != tt.ok {
					t.Fatalf("ok = %t, want %t", ok, tt.ok)
				}

				if !ok {
					if got != nil {
						t.Fatal("got != nil")
					}
					return
				}

				if got != tt.resp {
					t.Fatalf("got = %v, want %v", got, tt.resp)
				}
			})
		}
	})

	t.Run("integration", func(t *testing.T) {
		const responseBody = `{"user": {"id": "1", "email":"foo@bar.com"}}`

		setup()
		defer teardown()

		mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			_, _ = w.Write([]byte(responseBody))
		})

		c := defaultTestClient(server.URL, "foo")
		c.SetDebugFlag(DebugCaptureLastResponse)

		_, err := c.GetUser("1", GetUserOptions{})
		testErrCheck(t, "c.GetUser()", "", err)

		got, ok := c.LastAPIResponse()
		if !ok {
			t.Fatal("ok = false, want true")
		}

		if got == nil {
			t.Fatal("got == nil")
		}

		if got.StatusCode != 200 {
			t.Errorf("got.StatusCode = %d, want 200", got.StatusCode)
		}

		body, err := ioutil.ReadAll(got.Body)
		testErrCheck(t, "ioutil.ReadAll()", "", err)

		if jb := string(body); jb != responseBody {
			t.Fatalf("got.Body = %q, want %q", jb, responseBody)
		}
	})
}

func clientDoHandler(t *testing.T, needsAuth bool) func(w http.ResponseWriter, r *http.Request) {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		t.Helper()
		testMethod(t, r, http.MethodPost)

		auth := r.Header.Get("Authorization")
		if needsAuth && auth != "Token token=foo" {
			_, _ = w.Write([]byte("badAuth"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !needsAuth && len(auth) > 0 {
			_, _ = w.Write([]byte("Authentication header should not be provided"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if accept := r.Header.Get("Accept"); accept != acceptHeader {
			_, _ = w.Write([]byte(fmt.Sprintf("%q Accept unexpected", accept)))
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		if ua := r.Header.Get("User-Agent"); ua != userAgentHeader {
			_, _ = w.Write([]byte(fmt.Sprintf("%q User-Agent unexpected", ua)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if ct := r.Header.Get("Content-Type"); ct != contentTypeHeader {
			_, _ = w.Write([]byte(fmt.Sprintf("%q Content-Type unexpected", ct)))
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		_, _ = w.Write([]byte("ok"))
	}
}

func TestClient_Do(t *testing.T) {
	c := defaultTestClient(server.URL, "foo")

	tests := []struct {
		name string
		auth bool
	}{
		{
			name: "no_auth",
			auth: false,
		},
		{
			name: "auth",
			auth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			mux.HandleFunc("/test", clientDoHandler(t, tt.auth))

			req, err := http.NewRequest(http.MethodPost, server.URL+"/test", strings.NewReader(`{"empty":"object"}`))
			testErrCheck(t, "http.NewRequest()", "", err)

			resp, err := c.Do(req, tt.auth)
			testErrCheck(t, "c.Do()", "", err)

			defer func() {
				_, _ = io.Copy(ioutil.Discard, resp.Body)
				_ = resp.Body.Close()
			}()

			body, err := ioutil.ReadAll(resp.Body)
			testErrCheck(t, "ioutil.ReadAll()", "", err)

			if resp.StatusCode != 200 {
				t.Fatalf("request failed with status %q: %s", resp.Status, string(body))
			}

			if bs := string(body); bs != "ok" {
				t.Fatalf("body = %s, want ok", bs)
			}
		})
	}
}

func TestClient_UserAgentDefault(t *testing.T) {
	setup()
	defer teardown()

	defaultUserAgent := userAgentHeader
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		userAgentHeader := r.Header.Get("User-Agent")
		if userAgentHeader != defaultUserAgent {
			t.Fatalf("want %q, but got %q", defaultUserAgent, userAgentHeader)
		}
		w.WriteHeader(http.StatusOK)
	})

	client := defaultTestClient(server.URL, "foo")

	_, err := client.do(context.Background(), "GET", "/foo", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_UserAgentOverwrite(t *testing.T) {
	setup()
	defer teardown()

	terraformVersion := "terraform-version-for-testing"
	newUserAgent := fmt.Sprintf("(%s %s) Terraform/%s", runtime.GOOS, runtime.GOARCH, terraformVersion)
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		userAgentHeader := r.Header.Get("User-Agent")
		if userAgentHeader != newUserAgent {
			t.Fatalf("want %q, but got %q", newUserAgent, userAgentHeader)
		}
		w.WriteHeader(http.StatusOK)
	})

	client := NewClient("foo",
		WithAPIEndpoint(server.URL),
		WithTerraformProvider(terraformVersion),
	)

	_, err := client.do(context.Background(), "GET", "/foo", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullAPIErrorObject_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
		err  string
		want *NullAPIErrorObject
	}{
		{
			name: "error_per_api_spec",
			json: `{"code":42,"message":"test message","errors":["first error","second error"]}`,
			want: &NullAPIErrorObject{
				Valid: true,
				ErrorObject: APIErrorObject{
					Code:    42,
					Message: "test message",
					Errors: []string{
						"first error",
						"second error",
					},
				},
			},
		},
		{
			name: "issue_339",
			json: `{"code":84,"message":"other message","errors":"only error"}`,
			want: &NullAPIErrorObject{
				Valid: true,
				ErrorObject: APIErrorObject{
					Code:    84,
					Message: "other message",
					Errors: []string{
						"only error",
					},
				},
			},
		},
		{
			name: "returns_type_errors",
			json: `{"code":"42","message":"test message","errors":"first error"}`,
			err:  "json: cannot unmarshal string into Go struct field APIErrorObject.code of type int",
		},
		{
			name: "returns_syntax_errors",
			json: `}`,
			err:  "invalid character '}' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &NullAPIErrorObject{}

			err := json.Unmarshal([]byte(tt.json), &got)
			if !testErrCheck(t, "*NullAPIErrorObject.UnmarshalJSON()", tt.err, err) {
				return
			}

			testEqual(t, tt.want, got)
		})
	}
}

func Test_dupeRequest(t *testing.T) {
	reqA, _ := http.NewRequest(http.MethodGet, "http://localhost/api/v1/foo", strings.NewReader("some data\n"))
	reqB, _ := http.NewRequest(http.MethodGet, "http://localhost/api/v1/foo", nil)

	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			name: "ok",
			req:  reqA,
		},
		{
			name: "ok_nil_body",
			req:  reqB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotr, err := dupeRequest(tt.req)
			testErrCheck(t, "dupeRequest()", "", err)

			var gotb, wantb []byte
			if gotr.Body != nil {
				gotb, err = ioutil.ReadAll(gotr.Body)
				testErrCheck(t, "ioutil.ReadAll(gotr.Body)", "", err)
			}

			if tt.req.Body != nil {
				wantb, err = ioutil.ReadAll(tt.req.Body)
				testErrCheck(t, "ioutil.ReadAll(tt.req.Body)", "", err)
			}

			if !bytes.Equal(gotb, wantb) {
				t.Fatalf("gotb = %q, want %q", gotb, wantb)
			}
		})
	}
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func TestIncidentShow(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	homedir.Reset()
	t.Cleanup(homedir.Reset)

	logger := log.StandardLogger()
	output, formatter, level := logger.Out, logger.Formatter, logger.Level
	t.Cleanup(func() {
		log.SetOutput(output)
		log.SetFormatter(formatter)
		log.SetLevel(level)
	})

	const incident = `{"incident":{"id":"P123ABC","title":"Example incident","status":"triggered"}}`
	for _, tc := range []struct {
		name          string
		args          []string
		config        string
		token         string
		requestURI    string
		response      string
		status        int
		wantSuccess   bool
		wantIncident  bool
		wantHelp      bool
		closedServer  bool
		stallResponse bool
		unwritableOut bool
	}{
		{
			name: "configuration token", args: []string{"-id", "P123ABC"},
			requestURI: "/incidents/P123ABC", response: incident,
			wantSuccess: true, wantIncident: true,
		},
		{
			name: "flag token overrides configuration", args: []string{"-id", "P123ABC", "-authtoken", "flag-token"},
			token: "flag-token", requestURI: "/incidents/P123ABC", response: incident,
			wantSuccess: true, wantIncident: true,
		},
		{
			name: "escaped incident ID", args: []string{"-id", "P/id?x=1#fragment"},
			requestURI: "/incidents/P%2Fid%3Fx=1%23fragment", response: incident,
			wantSuccess: true, wantIncident: true,
		},
		{name: "missing ID"},
		{name: "blank ID", args: []string{"-id", " \t "}},
		{name: "unexpected argument", args: []string{"-id", "P123ABC", "extra"}},
		{name: "unknown flag", args: []string{"-unknown"}},
		{name: "missing flag value", args: []string{"-id"}},
		{name: "missing token", args: []string{"-id", "P123ABC"}, config: "{}\n"},
		{name: "help", args: []string{"-h"}, config: "{}\n", wantSuccess: true, wantHelp: true},
		{
			name: "API error", args: []string{"-id", "P123ABC"},
			requestURI: "/incidents/P123ABC", status: http.StatusNotFound,
			response: `{"error":{"code":2100,"message":"Not Found"}}`,
		},
		{
			name: "invalid JSON", args: []string{"-id", "P123ABC"},
			requestURI: "/incidents/P123ABC", response: `{`,
		},
		{
			name: "missing incident", args: []string{"-id", "P123ABC"},
			requestURI: "/incidents/P123ABC", response: `{}`,
		},
		{
			name: "request timeout", args: []string{"-id", "P123ABC"},
			requestURI: "/incidents/P123ABC", response: incident, stallResponse: true,
		},
		{name: "connection failure", args: []string{"-id", "P123ABC"}, closedServer: true},
		{
			name: "output failure", args: []string{"-id", "P123ABC"},
			requestURI: "/incidents/P123ABC", response: incident, unwritableOut: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := tc.config
			if config == "" {
				config = "authtoken: test-token\n"
			}
			if err := os.WriteFile(filepath.Join(home, ".pd.yml"), []byte(config), 0600); err != nil {
				t.Fatal(err)
			}

			var requests int32
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				if r.Method != http.MethodGet || r.RequestURI != tc.requestURI {
					t.Errorf("request = %s %s, want GET %s", r.Method, r.RequestURI, tc.requestURI)
				}
				token := tc.token
				if token == "" {
					token = "test-token"
				}
				if got := r.Header.Get("Authorization"); got != "Token token="+token {
					t.Errorf("Authorization = %q, want token %q", got, token)
				}
				if tc.stallResponse {
					timer := time.NewTimer(35 * time.Second)
					defer timer.Stop()
					select {
					case <-r.Context().Done():
						return
					case <-timer.C:
					}
				}
				w.Header().Set("Content-Type", "application/json")
				if tc.status != 0 {
					w.WriteHeader(tc.status)
				}
				fmt.Fprint(w, tc.response)
			}))
			t.Cleanup(server.Close)

			// SDK clients share this HTTP client. Keep these tests sequential and
			// route its real TLS connections exclusively to the test server.
			client := pagerduty.NewClient("").HTTPClient.(*http.Client)
			transport := server.Client().Transport.(*http.Transport).Clone()
			transport.Proxy = nil
			transport.TLSClientConfig.ServerName = "127.0.0.1"
			transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				if address != "api.pagerduty.com:443" {
					return nil, fmt.Errorf("unexpected destination: %s", address)
				}
				return (&net.Dialer{}).DialContext(ctx, network, server.Listener.Addr().String())
			}
			originalTransport := client.Transport
			client.Transport = transport
			t.Cleanup(func() {
				client.Transport = originalTransport
				transport.CloseIdleConnections()
			})
			if (&Meta{}).Client().HTTPClient != client {
				t.Fatal("command client does not share the test transport")
			}
			if tc.closedServer {
				server.Close()
			}

			stdout, err := os.CreateTemp(t.TempDir(), "stdout")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { stdout.Close() })
			originalStdout := os.Stdout
			os.Stdout = stdout
			t.Cleanup(func() { os.Stdout = originalStdout })
			if tc.unwritableOut {
				readOnly, err := os.Open(os.DevNull)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { readOnly.Close() })
				os.Stdout = readOnly
			}
			stderr, err := os.CreateTemp(t.TempDir(), "stderr")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { stderr.Close() })
			originalStderr := os.Stderr
			os.Stderr = stderr
			t.Cleanup(func() { os.Stderr = originalStderr })
			var diagnostics bytes.Buffer
			log.SetOutput(&diagnostics)

			code := invokeCLI(append([]string{"incident", "show"}, tc.args...))
			if (code == 0) != tc.wantSuccess {
				t.Errorf("exit code = %d, want success %t; diagnostics: %s", code, tc.wantSuccess, diagnostics.String())
			}
			if tc.stallResponse && !strings.Contains(diagnostics.String(), "context deadline exceeded") {
				t.Errorf("missing timeout diagnostic: %s", diagnostics.String())
			}
			wantRequests := int32(0)
			if tc.requestURI != "" {
				wantRequests = 1
			}
			if got := atomic.LoadInt32(&requests); got != wantRequests {
				t.Errorf("requests = %d, want %d", got, wantRequests)
			}
			data, err := os.ReadFile(stdout.Name())
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantIncident {
				var got pagerduty.Incident
				if err := yaml.Unmarshal(data, &got); err != nil {
					t.Fatalf("invalid YAML output: %v", err)
				}
				if got.ID != "P123ABC" || got.Title != "Example incident" || got.Status != "triggered" {
					t.Errorf("incident output = %s", data)
				}
			} else if tc.wantHelp {
				data, err := os.ReadFile(stderr.Name())
				if err != nil {
					t.Fatal(err)
				}
				for _, option := range []string{"-id", "-authtoken", "-loglevel"} {
					if !strings.Contains(string(data), option) {
						t.Errorf("help does not describe %s: %s", option, data)
					}
				}
			} else if tc.requestURI != "" && len(data) != 0 {
				t.Errorf("unexpected output on failure: %s", data)
			}
		})
	}
}

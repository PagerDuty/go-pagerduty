package webhookv3

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

const (
	secret      = "lDQHScfUeXUKaQRNF+8XIiDKZ7XX3itBAYzwU0TARw8lJqRnkKl2iB1anSb0Z+IK" /* #nosec */
	defaultBody = `{"event":{"id":"01BWDWL3NYY7LUFPZCC28QUCMK","event_type":"incident.priority_updated","resource_type":"incident","occurred_at":"2021-04-26T17:36:27.458Z","agent":{"html_url":"https://acme.pagerduty.com/users/PLH1HKV","id":"PLH1HKV","self":"https://api.pagerduty.com/users/PLH1HKV","summary":"Tenex Engineer","type":"user_reference"},"client":null,"data":{"id":"PGR0VU2","type":"incident","self":"https://api.pagerduty.com/incidents/PGR0VU2","html_url":"https://acme.pagerduty.com/incidents/PGR0VU2","number":2,"status":"triggered","title":"A little bump in the road","service":{"html_url":"https://acme.pagerduty.com/services/PF9KMXH","id":"PF9KMXH","self":"https://api.pagerduty.com/services/PF9KMXH","summary":"API Service","type":"service_reference"},"assignees":[{"html_url":"https://acme.pagerduty.com/users/PTUXL6G","id":"PTUXL6G","self":"https://api.pagerduty.com/users/PTUXL6G","summary":"User 123","type":"user_reference"}],"escalation_policy":{"html_url":"https://acme.pagerduty.com/escalation_policies/PUS0KTE","id":"PUS0KTE","self":"https://api.pagerduty.com/escalation_policies/PUS0KTE","summary":"Default","type":"escalation_policy_reference"},"teams":[{"html_url":"https://acme.pagerduty.com/teams/PFCVPS0","id":"PFCVPS0","self":"https://api.pagerduty.com/teams/PFCVPS0","summary":"Engineering","type":"team_reference"}],"priority":{"html_url":"https://acme.pagerduty.com/account/incident_priorities","id":"PSO75BM","self":"https://api.pagerduty.com/priorities/PSO75BM","summary":"P1","type":"priority_reference"},"urgency":"high","conference_bridge":{"conference_number":1000,"conference_url":"https://example.com"},"resolve_reason":null}}}`
)

func TestVerifySignature(t *testing.T) {
	tests := []struct {
		name string
		sig  string
		body string
		err  error
	}{
		{
			name: "valid",
			sig:  "v1=0c0b9495b893a39e70d1fea2fe11fbe0a825f88b9f67846f6cc07dd2bc5476cd",
			body: defaultBody,
		},
		{
			name: "mismatch",
			sig:  "v1=7020c8a7ec668a9b7012bc3dd82e483394b038f4230acc6785efbf2a7d8bcaf5",
			body: defaultBody,
			err:  ErrNoValidSignatures,
		},
		{
			name: "malformed_header",
			body: defaultBody,
			err:  ErrMalformedHeader,
		},
		{
			name: "malformed_body",
			sig:  "v1=0c0b9495b893a39e70d1fea2fe11fbe0a825f88b9f67846f6cc07dd2bc5476cd",
			err:  ErrMalformedBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:80/test", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("failed to generate new request: %s", err.Error())
			}

			req.Header.Set("X-PagerDuty-Signature", tt.sig)

			testErrIs(t, "VerifySignature", tt.err, VerifySignature(req, secret))
		})
	}
}

// testErrIs looks to see if wantErr is gotErr. If not, this calls t.Fatal(). It
// also calls t.Fatal() if there gotErr is not nil, but wantErr is. Returns true
// if you should continue running the test, or false if you should stop the
// test.
func testErrIs(t *testing.T, name string, wantErr, gotErr error) bool {
	t.Helper()

	if wantErr != nil {
		if gotErr == nil {
			t.Fatalf("%s error = <nil>, should be %v", name, wantErr)
			return false
		}

		if !errors.Is(gotErr, wantErr) {
			t.Fatalf("error %v is not %v", gotErr, wantErr)
			return false
		}
	}

	if wantErr == nil && gotErr != nil {
		t.Fatalf("%s unexpected error: %v", name, gotErr)
		return false
	}

	return true
}

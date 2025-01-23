package pagerduty

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// Create Integration
func TestClient_CreateIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Integration{
		Name: "foo",
	}
	servID := "1"

	res, err := client.CreateIntegration(servID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Get Integration
func TestClient_GetIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := GetIntegrationOptions{
		Includes: []string{},
	}
	servID := "1"
	intID := "1"

	res, err := client.GetIntegration(servID, intID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update Integration
func TestClient_UpdateIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"integration": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	servID := "1"

	res, err := client.UpdateIntegration(servID, input)

	want := &Integration{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete Integration
func TestClient_DeleteIntegration(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/1/integrations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	servID := "1"
	intID := "1"
	err := client.DeleteIntegration(servID, intID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIntegrationEmailFilterMode_String(t *testing.T) {
	tests := []struct {
		name  string
		input IntegrationEmailFilterMode
		want  string
	}{
		{name: "unknown", input: ^IntegrationEmailFilterMode(0), want: "invalid"},
		{name: "invalid", input: EmailFilterModeInvalid, want: "invalid"},
		{name: "all", input: EmailFilterModeAll, want: "all-email"},
		{name: "or", input: EmailFilterModeOr, want: "or-rules-email"},
		{name: "and", input: EmailFilterModeAnd, want: "and-rules-email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIntegrationEmailFilterMode_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input IntegrationEmailFilterMode
		want  string
	}{
		{name: "unknown", input: ^IntegrationEmailFilterMode(0), want: "invalid"},
		{name: "invalid", input: EmailFilterModeInvalid, want: "invalid"},
		{name: "all", input: EmailFilterModeAll, want: "all-email"},
		{name: "or", input: EmailFilterModeOr, want: "or-rules-email"},
		{name: "and", input: EmailFilterModeAnd, want: "and-rules-email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json, err := tt.input.MarshalJSON()
			testErrCheck(t, "tt.input.MarshalJSON()", "", err)

			want := fmt.Sprintf("%q", tt.want)

			if got := string(json); got != want {
				t.Fatalf("MarshalJSON() = `%s`, want `%s`", got, want)
			}
		})
	}
}

func TestIntegrationEmailFilterMode_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  IntegrationEmailFilterMode
		err   string
	}{
		{name: "invalid", input: `"invalid"`, err: `unknown value "invalid"`},
		{name: "null", input: `null`, err: `value cannot be null`},
		{name: "not_string", input: `42`, err: `json: cannot unmarshal number into Go value of type string`},
		{name: "all", input: `"all-email"`, want: EmailFilterModeAll},
		{name: "or", input: `"or-rules-email"`, want: EmailFilterModeOr},
		{name: "and", input: `"and-rules-email"`, want: EmailFilterModeAnd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(IntegrationEmailFilterMode)

			if cont := testErrCheck(t, "got.UnmarshalJSON()", tt.err, got.UnmarshalJSON([]byte(tt.input))); !cont {
				return
			}

			if *got != tt.want {
				t.Fatalf("got = %d (%s), want = %d (%s)", *got, got.String(), tt.want, tt.want.String())
			}
		})
	}
}

func TestIntegrationEmailFilterRuleMode_String(t *testing.T) {
	tests := []struct {
		name  string
		input IntegrationEmailFilterRuleMode
		want  string
	}{
		{name: "unknown", input: ^IntegrationEmailFilterRuleMode(0), want: "invalid"},
		{name: "invalid", input: EmailFilterRuleModeInvalid, want: "invalid"},
		{name: "always", input: EmailFilterRuleModeAlways, want: "always"},
		{name: "match", input: EmailFilterRuleModeMatch, want: "match"},
		{name: "no_match", input: EmailFilterRuleModeNoMatch, want: "no-match"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIntegrationEmailFilterRuleMode_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input IntegrationEmailFilterRuleMode
		want  string
	}{
		{name: "unknown", input: ^IntegrationEmailFilterRuleMode(0), want: "invalid"},
		{name: "invalid", input: EmailFilterRuleModeInvalid, want: "invalid"},
		{name: "always", input: EmailFilterRuleModeAlways, want: "always"},
		{name: "match", input: EmailFilterRuleModeMatch, want: "match"},
		{name: "no_match", input: EmailFilterRuleModeNoMatch, want: "no-match"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json, err := tt.input.MarshalJSON()
			testErrCheck(t, "tt.input.MarshalJSON()", "", err)

			want := fmt.Sprintf("%q", tt.want)

			if got := string(json); got != want {
				t.Fatalf("MarshalJSON() = `%s`, want `%s`", got, want)
			}
		})
	}
}

func TestIntegrationEmailFilterRuleMode_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  IntegrationEmailFilterRuleMode
		err   string
	}{
		{name: "invalid", input: `"invalid"`, err: `unknown value "invalid"`},
		{name: "null", input: `null`, err: `value cannot be null`},
		{name: "not_string", input: `42`, err: `json: cannot unmarshal number into Go value of type string`},
		{name: "always", input: `"always"`, want: EmailFilterRuleModeAlways},
		{name: "match", input: `"match"`, want: EmailFilterRuleModeMatch},
		{name: "no_match", input: `"no-match"`, want: EmailFilterRuleModeNoMatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(IntegrationEmailFilterRuleMode)

			if cont := testErrCheck(t, "got.UnmarshalJSON()", tt.err, got.UnmarshalJSON([]byte(tt.input))); !cont {
				return
			}

			if *got != tt.want {
				t.Fatalf("got = %d (%s), want = %d (%s)", *got, got.String(), tt.want, tt.want.String())
			}
		})
	}
}

func TestIntegrationEmailFilterRule(t *testing.T) {
	t.Run("zero_value/omitempty", func(t *testing.T) {
		var iefr IntegrationEmailFilterRule

		j, err := json.Marshal(iefr)
		testErrCheck(t, "json.Marshal()", "", err)

		if string(j) != "{}" {
			t.Fatalf("expected empty object, got %q", string(j))
		}
	})
}

func TestIntegrationEmailFilterRule_UnmarshalJSON(t *testing.T) {
	subjectRegex := ""
	bodyRegex := "testbody"
	fromEmailRegex := "testform"

	tests := []struct {
		name  string
		input string
		want  IntegrationEmailFilterRule
		err   string
	}{
		{
			name:  "full",
			input: fmt.Sprintf(`{"subject_mode":"always", "subject_regex":"%s", "body_mode":"match", "body_regex":"%s", "from_email_mode":"no-match", "from_email_regex":"%s"}`, subjectRegex, bodyRegex, fromEmailRegex),
			want: IntegrationEmailFilterRule{
				SubjectMode:    EmailFilterRuleModeAlways,
				SubjectRegex:   &subjectRegex,
				BodyMode:       EmailFilterRuleModeMatch,
				BodyRegex:      &bodyRegex,
				FromEmailMode:  EmailFilterRuleModeNoMatch,
				FromEmailRegex: &fromEmailRegex,
			},
		},

		{
			name:  "empty",
			input: `{}`,
			want: IntegrationEmailFilterRule{
				SubjectMode:    EmailFilterRuleModeInvalid,
				SubjectRegex:   new(string),
				BodyMode:       EmailFilterRuleModeInvalid,
				BodyRegex:      new(string),
				FromEmailMode:  EmailFilterRuleModeInvalid,
				FromEmailRegex: new(string),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(IntegrationEmailFilterRule)

			err := got.UnmarshalJSON([]byte(tt.input))

			if !testErrCheck(t, "got.UnmarshalJSON()", tt.err, err) {
				return
			}

			if got.SubjectRegex == nil {
				t.Error("got.SubjectRegex = <nil>")
			}

			if got.BodyRegex == nil {
				t.Error("got.BodyRegex = <nil>")
			}

			if got.FromEmailRegex == nil {
				t.Error("got.FromEmailRegex = <nil>")
			}

			testEqual(t, tt.want, *got)
		})
	}
}

package main

import (
	"os"
	"testing"
)

func TestInvokeCLIVersion(t *testing.T) {
	args := []string{"-v"}
	if invokeCLI(args) != 0 {
		t.Errorf("`pd -v` return code is non zero")
	}
}

func TestMetaAuthtokenFromEnv(t *testing.T) {
	t.Setenv("PAGERDUTY_TOKEN", "test-token-from-env")

	m := &Meta{}
	// loadConfig may warn about missing config file; that's expected in tests
	_ = m.loadConfig()

	if m.Authtoken != "test-token-from-env" {
		t.Errorf("expected authtoken from PAGERDUTY_TOKEN env var, got %q", m.Authtoken)
	}
}

func TestMetaAuthtokenFlagTakesPrecedenceOverEnv(t *testing.T) {
	t.Setenv("PAGERDUTY_TOKEN", "env-token")

	m := &Meta{Authtoken: "flag-token"}
	_ = m.loadConfig()

	if m.Authtoken != "flag-token" {
		t.Errorf("expected flag authtoken to take precedence, got %q", m.Authtoken)
	}
}

func TestMetaAuthtokenEnvNotSetWhenEmpty(t *testing.T) {
	os.Unsetenv("PAGERDUTY_TOKEN")

	m := &Meta{}
	_ = m.loadConfig()

	if m.Authtoken != "" {
		t.Errorf("expected empty authtoken when env var not set, got %q", m.Authtoken)
	}
}

package pagerduty

import (
	"strings"
	"testing"
)

// DecodeWebhook
func TestWebhook_DecodeWebhook(t *testing.T) {
	setup()
	defer teardown()

	jsonData := strings.NewReader(`{"id": "1"},{"id": "2"}`)
	res, err := DecodeWebhook(jsonData)

	want := &WebhookPayload{
		ID: "1",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

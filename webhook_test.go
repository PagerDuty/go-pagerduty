package pagerduty

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// DecodeWebhook
func TestWebhook_DecodeWebhook(t *testing.T) {
	setup()
	defer teardown()

	require := require.New(t)

	jsonData := strings.NewReader(`{"id": "1"},{"id": "2"}`)
	resp, err := DecodeWebhook(jsonData)

	want := &WebhookPayload{
		ID: "1",
	}

	require.NoError(err)
	require.Equal(want, resp)
}

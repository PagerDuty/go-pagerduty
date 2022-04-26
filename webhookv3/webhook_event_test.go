package webhookv3

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookPayload_UnmarshallJSON(t *testing.T) {
	var wp WebhookPayload

	data := `{ "event": { "id": "5ac64822-4adc-4fda-ade0-410becf0de4f", "event_type": "incident.priority_updated", "resource_type": "incident", "occurred_at": "2020-10-02T18:45:22.169Z", "agent": { "html_url": "https://acme.pagerduty.com/users/PLH1HKV", "id": "PLH1HKV", "self": "https://api.pagerduty.com/users/PLH1HKV", "summary": "Tenex Engineer", "type": "user_reference" }, "client": { "name": "PagerDuty" }, "data": { "id": "PGR0VU2", "type": "incident", "self": "https://api.pagerduty.com/incidents/PGR0VU2", "html_url": "https://acme.pagerduty.com/incidents/PGR0VU2", "number": 2, "status": "triggered", "title": "A little bump in the road", "service": { "html_url": "https://acme.pagerduty.com/services/PF9KMXH", "id": "PF9KMXH", "self": "https://api.pagerduty.com/services/PF9KMXH", "summary": "API Service", "type": "service_reference" }, "assignees": [ { "html_url": "https://acme.pagerduty.com/users/PTUXL6G", "id": "PTUXL6G", "self": "https://api.pagerduty.com/users/PTUXL6G", "summary": "User 123", "type": "user_reference" } ], "escalation_policy": { "html_url": "https://acme.pagerduty.com/escalation_policies/PUS0KTE", "id": "PUS0KTE", "self": "https://api.pagerduty.com/escalation_policies/PUS0KTE", "summary": "Default", "type": "escalation_policy_reference" }, "teams": [ { "html_url": "https://acme.pagerduty.com/teams/PFCVPS0", "id": "PFCVPS0", "self": "https://api.pagerduty.com/teams/PFCVPS0", "summary": "Engineering", "type": "team_reference" } ], "priority": { "html_url": "https://acme.pagerduty.com/account/incident_priorities", "id": "PSO75BM", "self": "https://api.pagerduty.com/priorities/PSO75BM", "summary": "P1", "type": "priority_reference" }, "urgency": "high", "conference_bridge": { "conference_number": "+1 1234123412,,987654321#", "conference_url": "https://example.com" }, "resolve_reason": null } } }`

	err := json.Unmarshal([]byte(data), &wp)
	assert.NoError(t, err)

	oe := wp.Event

	assert.Equal(t, "5ac64822-4adc-4fda-ade0-410becf0de4f", oe.ID)
	assert.Equal(t, "incident.priority_updated", oe.EventType)
	assert.Equal(t, "incident", oe.ResourceType)
	assert.Equal(t, "2020-10-02T18:45:22.169Z", oe.OccurredAt)

	assert.Equal(t, "PLH1HKV", oe.Agent.ID)
	assert.Equal(t, "user_reference", oe.Agent.Type)
}

func TestWebhookEvent_GetEventDataValue(t *testing.T) {
	var wp WebhookPayload

	data := `{ "event": { "id": "5ac64822-4adc-4fda-ade0-410becf0de4f", "event_type": "incident.priority_updated", "resource_type": "incident", "occurred_at": "2020-10-02T18:45:22.169Z", "agent": { "html_url": "https://acme.pagerduty.com/users/PLH1HKV", "id": "PLH1HKV", "self": "https://api.pagerduty.com/users/PLH1HKV", "summary": "Tenex Engineer", "type": "user_reference" }, "client": { "name": "PagerDuty" }, "data": { "id": "PGR0VU2", "type": "incident", "self": "https://api.pagerduty.com/incidents/PGR0VU2", "html_url": "https://acme.pagerduty.com/incidents/PGR0VU2", "number": 2, "status": "triggered", "title": "A little bump in the road", "service": { "html_url": "https://acme.pagerduty.com/services/PF9KMXH", "id": "PF9KMXH", "self": "https://api.pagerduty.com/services/PF9KMXH", "summary": "API Service", "type": "service_reference" }, "assignees": [ { "html_url": "https://acme.pagerduty.com/users/PTUXL6G", "id": "PTUXL6G", "self": "https://api.pagerduty.com/users/PTUXL6G", "summary": "User 123", "type": "user_reference" } ], "escalation_policy": { "html_url": "https://acme.pagerduty.com/escalation_policies/PUS0KTE", "id": "PUS0KTE", "self": "https://api.pagerduty.com/escalation_policies/PUS0KTE", "summary": "Default", "type": "escalation_policy_reference" }, "teams": [ { "html_url": "https://acme.pagerduty.com/teams/PFCVPS0", "id": "PFCVPS0", "self": "https://api.pagerduty.com/teams/PFCVPS0", "summary": "Engineering", "type": "team_reference" } ], "priority": { "html_url": "https://acme.pagerduty.com/account/incident_priorities", "id": "PSO75BM", "self": "https://api.pagerduty.com/priorities/PSO75BM", "summary": "P1", "type": "priority_reference" }, "urgency": "high", "conference_bridge": { "conference_number": "+1 1234123412,,987654321#", "conference_url": "https://example.com" }, "resolve_reason": null } } }`

	err := json.Unmarshal([]byte(data), &wp)
	assert.NoError(t, err)

	oe := wp.Event

	value, _ := oe.GetEventDataValue("type")
	assert.Equal(t, "incident", value)

	value, _ = oe.GetEventDataValue("title")
	assert.Equal(t, "A little bump in the road", value)

	value, _ = oe.GetEventDataValue("service", "summary")
	assert.Equal(t, "API Service", value)

	value, err = oe.GetEventDataValue("not_a_field")
	assert.Equal(t, "", value)
	assert.Error(t, err)

	value, _ = oe.GetEventDataValue("assignees", "0", "summary")
	assert.Equal(t, "User 123", value)

	value, err = oe.GetEventDataValue("assignees", "not_an_integer")
	assert.Equal(t, "", value)
	assert.Error(t, err)
}

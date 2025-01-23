package pagerduty

import (
	"strings"
	"testing"
)

const webhookPayload = `{"messages":[{"event":"incident.trigger","log_entries":[{"id":"R2XGXEI3W0FHMSDXHDIBQGBQ5E","type":"trigger_log_entry","summary":"Triggered through the website","self":"https://api.pagerduty.com/log_entries/R2XGXEI3W0FHMSDXHDIBQGBQ5E","html_url":"https://webdemo.pagerduty.com/incidents/PRORDTY/log_entries/R2XGXEI3W0FHMSDXHDIBQGBQ5E","created_at":"2017-09-26T15:14:36Z","agent":{"id":"P553OPV","type":"user_reference","summary":"Laura Haley","self":"https://api.pagerduty.com/users/P553OPV","html_url":"https://webdemo.pagerduty.com/users/P553OPV"},"channel":{"type":"web_trigger","summary":"My new incident","subject":"My new incident","details":"Oh my gosh","details_omitted":false},"service":{"id":"PN49J75","type":"service_reference","summary":"Production XDB Cluster","self":"https://api.pagerduty.com/services/PN49J75","html_url":"https://webdemo.pagerduty.com/services/PN49J75"},"incident":{"id":"PRORDTY","type":"incident_reference","summary":"[#33] My new incident","self":"https://api.pagerduty.com/incidents/PRORDTY","html_url":"https://webdemo.pagerduty.com/incidents/PRORDTY"},"teams":[{"id":"P4SI59S","type":"team_reference","summary":"Engineering","self":"https://api.pagerduty.com/teams/P4SI59S","html_url":"https://webdemo.pagerduty.com/teams/P4SI59S"}],"contexts":[],"event_details":{"description":"My new incident"}}],"webhook":{"endpoint_url":"https://requestb.in/18ao6fs1","name":"V2 wabhook","description":null,"webhook_object":{"id":"PN49J75","type":"service_reference","summary":"Production XDB Cluster","self":"https://api.pagerduty.com/services/PN49J75","html_url":"https://webdemo.pagerduty.com/services/PN49J75"},"config":{},"outbound_integration":{"id":"PJFWPEP","type":"outbound_integration_reference","summary":"Generic V2 Webhook","self":"https://api.pagerduty.com/outbound_integrations/PJFWPEP","html_url":null},"accounts_addon":null,"id":"PKT9NNX","type":"webhook","summary":"V2 wabhook","self":"https://api.pagerduty.com/webhooks/PKT9NNX","html_url":null},"incident":{"incident_number":33,"title":"My new incident","description":"My new incident","created_at":"2017-09-26T15:14:36Z","status":"triggered","pending_actions":[{"type":"escalate","at":"2017-09-26T15:44:36Z"},{"type":"resolve","at":"2017-09-26T19:14:36Z"}],"incident_key":null,"service":{"id":"PN49J75","name":"Production XDB Cluster","description":"This service was created during onboarding on July 5, 2017.","auto_resolve_timeout":14400,"acknowledgement_timeout":1800,"created_at":"2017-07-05T17:33:09Z","status":"critical","last_incident_timestamp":"2017-09-26T15:14:36Z","teams":[{"id":"P4SI59S","type":"team_reference","summary":"Engineering","self":"https://api.pagerduty.com/teams/P4SI59S","html_url":"https://webdemo.pagerduty.com/teams/P4SI59S"}],"incident_urgency_rule":{"type":"constant","urgency":"high"},"scheduled_actions":[],"support_hours":null,"escalation_policy":{"id":"PINYWEF","type":"escalation_policy_reference","summary":"Default","self":"https://api.pagerduty.com/escalation_policies/PINYWEF","html_url":"https://webdemo.pagerduty.com/escalation_policies/PINYWEF"},"addons":[],"privilege":null,"alert_creation":"create_alerts_and_incidents","integrations":[{"id":"PUAYF96","type":"generic_events_api_inbound_integration_reference","summary":"API","self":"https://api.pagerduty.com/services/PN49J75/integrations/PUAYF96","html_url":"https://webdemo.pagerduty.com/services/PN49J75/integrations/PUAYF96"},{"id":"P90GZUH","type":"generic_email_inbound_integration_reference","summary":"Email","self":"https://api.pagerduty.com/services/PN49J75/integrations/P90GZUH","html_url":"https://webdemo.pagerduty.com/services/PN49J75/integrations/P90GZUH"}],"metadata":{},"type":"service","summary":"Production XDB Cluster","self":"https://api.pagerduty.com/services/PN49J75","html_url":"https://webdemo.pagerduty.com/services/PN49J75"},"assignments":[{"at":"2017-09-26T15:14:36Z","assignee":{"id":"P553OPV","type":"user_reference","summary":"Laura Haley","self":"https://api.pagerduty.com/users/P553OPV","html_url":"https://webdemo.pagerduty.com/users/P553OPV"}}],"acknowledgements":[],"last_status_change_at":"2017-09-26T15:14:36Z","last_status_change_by":{"id":"PN49J75","type":"service_reference","summary":"Production XDB Cluster","self":"https://api.pagerduty.com/services/PN49J75","html_url":"https://webdemo.pagerduty.com/services/PN49J75"},"first_trigger_log_entry":{"id":"R2XGXEI3W0FHMSDXHDIBQGBQ5E","type":"trigger_log_entry_reference","summary":"Triggered through the website","self":"https://api.pagerduty.com/log_entries/R2XGXEI3W0FHMSDXHDIBQGBQ5E","html_url":"https://webdemo.pagerduty.com/incidents/PRORDTY/log_entries/R2XGXEI3W0FHMSDXHDIBQGBQ5E"},"escalation_policy":{"id":"PINYWEF","type":"escalation_policy_reference","summary":"Default","self":"https://api.pagerduty.com/escalation_policies/PINYWEF","html_url":"https://webdemo.pagerduty.com/escalation_policies/PINYWEF"},"privilege":null,"teams":[{"id":"P4SI59S","type":"team_reference","summary":"Engineering","self":"https://api.pagerduty.com/teams/P4SI59S","html_url":"https://webdemo.pagerduty.com/teams/P4SI59S"}],"alert_counts":{"all":0,"triggered":0,"resolved":0},"impacted_services":[{"id":"PN49J75","type":"service_reference","summary":"Production XDB Cluster","self":"https://api.pagerduty.com/services/PN49J75","html_url":"https://webdemo.pagerduty.com/services/PN49J75"}],"is_mergeable":true,"basic_alert_grouping":null,"alert_grouping":null,"metadata":{},"external_references":[],"importance":null,"incidents_responders":[],"responder_requests":[],"subscriber_requests":[],"urgency":"high","id":"PRORDTY","type":"incident","summary":"[#33] My new incident","self":"https://api.pagerduty.com/incidents/PRORDTY","html_url":"https://webdemo.pagerduty.com/incidents/PRORDTY","alerts":[{"alert_key":"c24117fc42e44b44b4d6876190583378"}]},"id":"69a7ced0-a2cd-11e7-a799-22000a15839c","created_on":"2017-09-26T15:14:36Z"}]}`

// DecodeWebhook
func TestWebhook_DecodeWebhook(t *testing.T) {
	setup()
	defer teardown()

	jsonData := strings.NewReader(webhookPayload)
	res, err := DecodeWebhook(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Messages) != 1 {
		t.Fatal("Expect 1 message")
	}

	incidentDetails := res.Messages[0].Incident
	if incidentDetails.IncidentNumber != 33 {
		t.Fatal("Unexpected Incident Number")
	}

	if len(incidentDetails.PendingActions) != 2 {
		t.Fatal("Expected 2 pending actions")
	}

	if incidentDetails.Service.ID != "PN49J75" {
		t.Fatal("Unexpected Service ID")
	}

	if len(incidentDetails.Assignments) != 1 {
		t.Fatal("Expected 1 Assignment")
	}
}

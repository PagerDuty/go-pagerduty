package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// Delete Slack Connection test.
func TestSlackConnection_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-slack/workspaces/foo/connections/connectionid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteSlackConnectionWithContext(context.Background(), "foo", "connectionid")
	if err != nil {
		t.Fatal(err)
	}
}

// List Slack Connections test.
func TestSlackConnection_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-slack/workspaces/foo/connections", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`
		{
			"slack_connections":[
			   {
				  "id":"A12BCDE",
				  "source_id":"A1234B5",
				  "source_name":"test_service",
				  "source_type":"service_reference",
				  "channel_id":"A123B456C7D",
				  "channel_name":"random",
				  "notification_type":"responder",
				  "config":{
					 "events":[
						"incident.acknowledged",
						"incident.annotated",
						"incident.delegated",
						"incident.escalated",
						"incident.reassigned",
						"incident.resolved",
						"incident.triggered",
						"incident.unacknowledged",
						"incident.priority_updated",
						"incident.responder.added",
						"incident.responder.replied",
						"incident.status_update_published",
						"incident.reopened",
						"incident.action_invocation.created",
						"incident.action_invocation.updated",
						"incident.action_invocation.terminated"
					 ],
					 "priorities":[
						"ABCDEF1",
						"AB1CDE2"
					 ],
					 "urgency":null
				  }
			   }
			],
			"limit":1,
			"offset":0,
			"more":true,
			"total":99
		 }`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListSlackConnectionsWithContext(context.Background(), "foo", ListSlackConnectionsOptions{})

	want := &ListSlackConnectionsResponse{
		APIListObject: APIListObject{
			Limit: 1,
			More:  true,
			Total: 99,
		},
		Connections: []SlackConnectionObject{
			{
				APIObject: APIObject{
					ID: "A12BCDE",
				},
				SlackConnection: SlackConnection{
					SourceID:         "A1234B5",
					SourceName:       "test_service",
					SourceType:       "service_reference",
					ChannelID:        "A123B456C7D",
					ChannelName:      "random",
					NotificationType: "responder",
					Config: SlackConnectionConfig{
						Events: []string{
							"incident.acknowledged",
							"incident.annotated",
							"incident.delegated",
							"incident.escalated",
							"incident.reassigned",
							"incident.resolved",
							"incident.triggered",
							"incident.unacknowledged",
							"incident.priority_updated",
							"incident.responder.added",
							"incident.responder.replied",
							"incident.status_update_published",
							"incident.reopened",
							"incident.action_invocation.created",
							"incident.action_invocation.updated",
							"incident.action_invocation.terminated",
						},
						Priorities: []string{"ABCDEF1", "AB1CDE2"},
						Urgency:    nil,
					},
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update Slack Connection test.
func TestSlackConnection_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-slack/workspaces/foo/connections/connectionid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`
		{
			"slack_connection": {
			  "id": "A12BCDE",
			  "source_id": "A1234B5",
			  "source_name": "test_service",
			  "source_type": "service_reference",
			  "channel_id": "FOOBAR",
			  "channel_name": "random",
			  "notification_type": "responder",
			  "config": {
				"events": [
				  "incident.acknowledged",
				  "incident.annotated",
				  "incident.delegated",
				  "incident.escalated",
				  "incident.reassigned",
				  "incident.resolved",
				  "incident.triggered",
				  "incident.unacknowledged",
				  "incident.priority_updated",
				  "incident.responder.added",
				  "incident.responder.replied",
				  "incident.status_update_published",
				  "incident.reopened",
				  "incident.action_invocation.created",
				  "incident.action_invocation.updated",
				  "incident.action_invocation.terminated"
				],
				"priorities": [
					"ABCDEF1",
					"AB1CDE2"
				],
				"urgency": null
			  }
			}
		  }`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.UpdateSlackConnectionWithContext(context.Background(), "foo", "connectionid", SlackConnection{
		SourceID:         "A1234B5",
		SourceName:       "test_service",
		SourceType:       "service_reference",
		ChannelID:        "FOOBAR",
		ChannelName:      "random",
		NotificationType: "responder",
		Config: SlackConnectionConfig{
			Events: []string{
				"incident.acknowledged",
				"incident.annotated",
				"incident.delegated",
				"incident.escalated",
				"incident.reassigned",
				"incident.resolved",
				"incident.triggered",
				"incident.unacknowledged",
				"incident.priority_updated",
				"incident.responder.added",
				"incident.responder.replied",
				"incident.status_update_published",
				"incident.reopened",
				"incident.action_invocation.created",
				"incident.action_invocation.updated",
				"incident.action_invocation.terminated",
			},
			Priorities: []string{"ABCDEF1", "AB1CDE2"},
			Urgency:    nil,
		},
	})

	want := SlackConnectionObject{
		APIObject: APIObject{
			ID: "A12BCDE",
		},
		SlackConnection: SlackConnection{
			SourceID:         "A1234B5",
			SourceName:       "test_service",
			SourceType:       "service_reference",
			ChannelID:        "FOOBAR",
			ChannelName:      "random",
			NotificationType: "responder",
			Config: SlackConnectionConfig{
				Events: []string{
					"incident.acknowledged",
					"incident.annotated",
					"incident.delegated",
					"incident.escalated",
					"incident.reassigned",
					"incident.resolved",
					"incident.triggered",
					"incident.unacknowledged",
					"incident.priority_updated",
					"incident.responder.added",
					"incident.responder.replied",
					"incident.status_update_published",
					"incident.reopened",
					"incident.action_invocation.created",
					"incident.action_invocation.updated",
					"incident.action_invocation.terminated",
				},
				Priorities: []string{"ABCDEF1", "AB1CDE2"},
				Urgency:    nil,
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Get Slack Connection test.
func TestSlackConnection_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-slack/workspaces/foo/connections/connectionid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`
		{
			"slack_connection": {
			  "id": "A12BCDE",
			  "source_id": "A1234B5",
			  "source_name": "test_service",
			  "source_type": "service_reference",
			  "channel_id": "AABBCC",
			  "channel_name": "random",
			  "notification_type": "responder",
			  "config": {
				"events": [
				  "incident.acknowledged",
				  "incident.annotated",
				  "incident.delegated",
				  "incident.escalated",
				  "incident.reassigned",
				  "incident.resolved",
				  "incident.triggered",
				  "incident.unacknowledged",
				  "incident.priority_updated",
				  "incident.responder.added",
				  "incident.responder.replied",
				  "incident.status_update_published",
				  "incident.reopened",
				  "incident.action_invocation.created",
				  "incident.action_invocation.updated",
				  "incident.action_invocation.terminated"
				],
				"priorities": [
					"ABCDEF1",
					"AB1CDE2"
				],
				"urgency": null
			  }
			}
		  }`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetSlackConnectionWithContext(context.Background(), "foo", "connectionid")

	want := SlackConnectionObject{
		APIObject: APIObject{
			ID: "A12BCDE",
		},
		SlackConnection: SlackConnection{
			SourceID:         "A1234B5",
			SourceName:       "test_service",
			SourceType:       "service_reference",
			ChannelID:        "AABBCC",
			ChannelName:      "random",
			NotificationType: "responder",
			Config: SlackConnectionConfig{
				Events: []string{
					"incident.acknowledged",
					"incident.annotated",
					"incident.delegated",
					"incident.escalated",
					"incident.reassigned",
					"incident.resolved",
					"incident.triggered",
					"incident.unacknowledged",
					"incident.priority_updated",
					"incident.responder.added",
					"incident.responder.replied",
					"incident.status_update_published",
					"incident.reopened",
					"incident.action_invocation.created",
					"incident.action_invocation.updated",
					"incident.action_invocation.terminated",
				},
				Priorities: []string{"ABCDEF1", "AB1CDE2"},
				Urgency:    nil,
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

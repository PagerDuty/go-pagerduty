package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

func TestAudit_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/audit/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = w.Write([]byte(`{"records":[{"id":"PDRECORDID4_UPDATED_USERS_NOTIFICATION_RULE","execution_time":"2020-06-04T15:30:16.272Z","execution_context":{"request_id":"222lDEOIH-534-4ljhLHJjh222","remote_address":"201.19.20.19"},"actors":[{"id":"PDUSER","summary":"John Snow","type":"user_reference"}],"method":{"type":"api_token","truncated_token":"2adm"},"root_resource":{"id":"PDUSER","type":"user_reference","summary":"John Snow"},"action":"update","details":{"resource":{"id":"PXOGWUS","type":"assignment_notification_rule_reference","summary":"0 minutes: channel P1IAAPZ"},"fields":[{"name":"start_delay_in_minutes","before_value":"0","value":"2"}],"references":[{"name":"contact_method","removed":[{"id":"POE6L88","type":"push_notification_contact_method_reference","summary":"Pixel 3"}],"added":[{"id":"P4GTUMK","type":"sms_contact_method_reference","summary":"Mobile"}]}]}}],"next_cursor":null,"limit":10}`))
	})

	client := defaultTestClient(server.URL, "foo")
	recordOpts := ListAuditRecordsOptions{
		Limit:  10,
		Cursor: "",
	}

	resp, err := client.ListAuditRecords(context.Background(), recordOpts)
	if err != nil {
		t.Fatal(err)
	}

	want := ListAuditRecordsResponse{
		Records: []AuditRecord{
			{
				ID:            "PDRECORDID4_UPDATED_USERS_NOTIFICATION_RULE",
				ExecutionTime: "2020-06-04T15:30:16.272Z",
				ExecutionContext: ExecutionContext{
					RequestID:     "222lDEOIH-534-4ljhLHJjh222",
					RemoteAddress: "201.19.20.19",
				},
				Actors: []APIObject{
					{
						ID:      "PDUSER",
						Summary: "John Snow",
						Type:    "user_reference",
					},
				},
				Method: Method{
					Type:           "api_token",
					TruncatedToken: "2adm",
				},
				RootResource: APIObject{
					ID:      "PDUSER",
					Summary: "John Snow",
					Type:    "user_reference",
				},
				Action: "update",
				Details: Details{
					Resource: APIObject{
						ID:      "PXOGWUS",
						Type:    "assignment_notification_rule_reference",
						Summary: "0 minutes: channel P1IAAPZ",
					},
					Fields: []Field{
						{
							Name:        "start_delay_in_minutes",
							BeforeValue: "0",
							Value:       "2",
						},
					},
					References: []Reference{
						{
							Name: "contact_method",
							Removed: []APIObject{
								{
									ID:      "POE6L88",
									Type:    "push_notification_contact_method_reference",
									Summary: "Pixel 3",
								},
							},
							Added: []APIObject{
								{
									ID:      "P4GTUMK",
									Type:    "sms_contact_method_reference",
									Summary: "Mobile",
								},
							},
						},
					},
				},
			},
		},
		NextCursor: nil,
		Limit:      10,
	}

	testEqual(t, want, resp)
}

package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

func TestJiraCloudAccountsMapping_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-jira-cloud/accounts_mappings/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
  "accounts_mappings": [
    {
      "created_at": "2024-09-12T12:34:55.000Z",
      "id": "PIJ90N7",
      "jira_cloud_account": {
        "base_url": "https://atlassian-subdomain.atlassian.net"
      },
      "pagerduty_account": {
        "subdomain": "pagerduty-subdomain"
      },
      "updated_at": "2024-09-13T10:11:12.000Z"
    }
  ],
  "limit": 25,
  "more": false,
  "offset": 0,
  "total": 1
}`))
	})
	client := defaultTestClient(server.URL, "foo")

	want := &ListJiraCloudAccountsMappingsResponse{
		APIListObject: APIListObject{
			Limit:  25,
			More:   false,
			Offset: 0,
			Total:  1,
		},
		AccountsMappings: []JiraCloudAccountsMapping{
			{
				ID: "PIJ90N7",
				JiraCloudAccount: JiraCloudAccount{
					BaseURL: "https://atlassian-subdomain.atlassian.net",
				},
				PagerDutyAccount: PagerDutyAccount{
					Subdomain: "pagerduty-subdomain",
				},
				CreatedAt: "2024-09-12T12:34:55.000Z",
				UpdatedAt: "2024-09-13T10:11:12.000Z",
			},
		},
	}

	ctx := context.Background()
	res, err := client.ListJiraCloudAccountsMappings(ctx, ListJiraCloudAccountsMappingsOptions{})
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestJiraCloudAccountsMapping_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-jira-cloud/accounts_mappings/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
  "created_at": "2024-09-12T12:34:55.000Z",
  "id": "PIJ90N7",
  "jira_cloud_account": {
    "base_url": "https://atlassian-subdomain.atlassian.net"
  },
  "pagerduty_account": {
    "subdomain": "pagerduty-subdomain"
  },
  "updated_at": "2024-09-13T10:11:12.000Z"
}`))
	})
	client := defaultTestClient(server.URL, "foo")

	want := &JiraCloudAccountsMapping{
		ID: "PIJ90N7",
		JiraCloudAccount: JiraCloudAccount{
			BaseURL: "https://atlassian-subdomain.atlassian.net",
		},
		PagerDutyAccount: PagerDutyAccount{
			Subdomain: "pagerduty-subdomain",
		},
		CreatedAt: "2024-09-12T12:34:55.000Z",
		UpdatedAt: "2024-09-13T10:11:12.000Z",
	}

	ctx := context.Background()
	res, err := client.GetJiraCloudAccountsMapping(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestJiraCloudAccountsMappingRule_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-jira-cloud/accounts_mappings/1/rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"id": "2"}`))
	})
	client := defaultTestClient(server.URL, "foo")

	want := JiraCloudAccountsMappingRule{ID: "2"}

	ctx := context.Background()
	res, err := client.CreateJiraCloudAccountsMappingRule(ctx, "1", want)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, &want, res)
}

func TestJiraCloudAccountsMappingRule_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-jira-cloud/accounts_mappings/1/rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"rules": [], "limit": 25, "more": false, "offset": 0, "total": 0}`))
	})
	client := defaultTestClient(server.URL, "foo")

	want := &ListJiraCloudAccountsMappingsResponse{
		APIListObject: APIListObject{
			Limit:  25,
			More:   false,
			Offset: 0,
			Total:  0,
		},
		AccountsMappings: nil,
	}

	ctx := context.Background()
	res, err := client.ListJiraCloudAccountsMappingRules(ctx, "1", ListJiraCloudAccountsMappingRulesOptions{})
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestJiraCloudAccountsMappingRule_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-jira-cloud/accounts_mappings/1/rules/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"id": "PF9KMXH"}`))
	})
	client := defaultTestClient(server.URL, "foo")

	want := &JiraCloudAccountsMappingRule{ID: "PF9KMXH"}

	ctx := context.Background()
	res, err := client.GetJiraCloudAccountsMappingRule(ctx, "1", "2")
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestJiraCloudAccountsMappingRule_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-jira-cloud/accounts_mappings/1/rules/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	client := defaultTestClient(server.URL, "foo")

	ctx := context.Background()
	err := client.DeleteJiraCloudAccountsMappingRule(ctx, "1", "2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestJiraCloudAccountsMappingRule_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration-jira-cloud/accounts_mappings/1/rules/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"id": "2"}`))
	})
	client := defaultTestClient(server.URL, "foo")

	want := JiraCloudAccountsMappingRule{ID: "2"}

	ctx := context.Background()
	res, err := client.UpdateJiraCloudAccountsMappingRule(ctx, "1", want)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, &want, res)
}

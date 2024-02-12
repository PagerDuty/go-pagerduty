package pagerduty

import (
	"net/http"
	"testing"
)

// ListStatusPages
func TestStatusPage_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["status_page_type"], []string{"public"})
		_, _ = w.Write([]byte(`{"status_pages": [{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := ListStatusPageOptions{
		StatusPageType: "public",
	}
	res, err := client.ListStatusPages(opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPagesResponse{
		APIListObject: APIListObject{},
		StatusPages: []StatusPage{
			{
				ID:             "1",
				Name:           "MyStatusPage",
				PublishedAt:    "2024-02-12T09:23:23Z",
				StatusPageType: "public",
				URL:            "https://mypagerduty",
			},
		},
	}

	testEqual(t, want, res)
}

// ListStatusPageImpacts
func TestStatusPage_ListImpacts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/impacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["post_type"], []string{"incident"})
		_, _ = w.Write([]byte(`{"impacts": [{"id": "1","description":"Extreme","post_type":"incident","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"}}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := ListStatusPageImpactOptions{
		PostType: "incident",
	}
	res, err := client.ListStatusPageImpacts("1", opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPageImpactsResponse{
		APIListObject: APIListObject{},
		StatusPageImpacts: []StatusPageImpact{
			{
				ID:          "1",
				Description: "Extreme",
				PostType:    "incident",
				StatusPage: StatusPage{
					ID:             "1",
					Name:           "MyStatusPage",
					PublishedAt:    "2024-02-12T09:23:23Z",
					StatusPageType: "public",
					URL:            "https://mypagerduty",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// GetStatusPageImpact
func TestStatusPage_GetImpact(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/impacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"impact": {"id": "1","description":"Extreme","post_type":"incident","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetStatusPageImpact("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPageImpact{
		ID:          "1",
		Description: "Extreme",
		PostType:    "incident",
		StatusPage: StatusPage{
			ID:             "1",
			Name:           "MyStatusPage",
			PublishedAt:    "2024-02-12T09:23:23Z",
			StatusPageType: "public",
			URL:            "https://mypagerduty",
		},
	}

	testEqual(t, want, res)
}

// ListStatusPageServices
func TestStatusPage_ListServices(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"services": [{"id": "1","name":"MyService","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"},"business_service":{"name":"MyService"}}]}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListStatusPageServices("1")
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPageServicesResponse{
		APIListObject: APIListObject{},
		StatusPageServices: []StatusPageService{
			{
				ID:   "1",
				Name: "MyService",
				StatusPage: StatusPage{
					ID:             "1",
					Name:           "MyStatusPage",
					PublishedAt:    "2024-02-12T09:23:23Z",
					StatusPageType: "public",
					URL:            "https://mypagerduty",
				},
				BusinessService: Service{
					Name: "MyService",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// GetStatusPageService
func TestStatusPage_GetServices(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"service": {"id": "1","name":"MyService","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"},"business_service":{"name":"MyService"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetStatusPageService("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPageService{
		ID:   "1",
		Name: "MyService",
		StatusPage: StatusPage{
			ID:             "1",
			Name:           "MyStatusPage",
			PublishedAt:    "2024-02-12T09:23:23Z",
			StatusPageType: "public",
			URL:            "https://mypagerduty",
		},
		BusinessService: Service{
			Name: "MyService",
		},
	}

	testEqual(t, want, res)
}

// ListStatusPageSeverities
func TestStatusPage_ListSeverities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/severities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["post_type"], []string{"incident"})
		_, _ = w.Write([]byte(`{"severities": [{"id": "1","description":"Extreme","post_type":"incident","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"}}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := ListStatusPageSeveritiesOptions{
		PostType: "incident",
	}

	res, err := client.ListStatusPageSeverities("1", opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPageSeveritiesResponse{
		APIListObject: APIListObject{},
		StatusPageSeverities: []StatusPageSeverity{
			{
				ID:          "1",
				Description: "Extreme",
				PostType:    "incident",
				StatusPage: StatusPage{
					ID:             "1",
					Name:           "MyStatusPage",
					PublishedAt:    "2024-02-12T09:23:23Z",
					StatusPageType: "public",
					URL:            "https://mypagerduty",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// GetStatusPageSeverity
func TestStatusPage_GetSeverity(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/severities/1/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"severity": {"id": "1","description":"Extreme","post_type":"incident","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetStatusPageSeverity("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPageSeverity{
		ID:          "1",
		Description: "Extreme",
		PostType:    "incident",
		StatusPage: StatusPage{
			ID:             "1",
			Name:           "MyStatusPage",
			PublishedAt:    "2024-02-12T09:23:23Z",
			StatusPageType: "public",
			URL:            "https://mypagerduty",
		},
	}

	testEqual(t, want, res)
}

// ListStatusPageStatuses
func TestStatusPage_ListStatuses(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/statuses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["post_type"], []string{"incident"})
		_, _ = w.Write([]byte(`{"statuses": [{"id": "1","description":"Extreme","post_type":"incident","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"}}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := ListStatusPageStatusesOptions{
		PostType: "incident",
	}

	res, err := client.ListStatusPageStatuses("1", opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPageStatusesResponse{
		APIListObject: APIListObject{},
		StatusPageStatuses: []StatusPageStatus{
			{
				ID:          "1",
				Description: "Extreme",
				PostType:    "incident",
				StatusPage: StatusPage{
					ID:             "1",
					Name:           "MyStatusPage",
					PublishedAt:    "2024-02-12T09:23:23Z",
					StatusPageType: "public",
					URL:            "https://mypagerduty",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// GetStatusPageStatus
func TestStatusPage_GetStatus(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/statuses/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"status": {"id": "1","description":"Extreme","post_type":"incident","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetStatusPageStatus("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPageStatus{
		ID:          "1",
		Description: "Extreme",
		PostType:    "incident",
		StatusPage: StatusPage{
			ID:             "1",
			Name:           "MyStatusPage",
			PublishedAt:    "2024-02-12T09:23:23Z",
			StatusPageType: "public",
			URL:            "https://mypagerduty",
		},
	}

	testEqual(t, want, res)
}

// ListStatusPagePosts
func TestStatusPage_ListPosts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["post_type"], []string{"incident"})
		testEqual(t, r.URL.Query()["reviewed_status"], []string{"approved"})
		testEqual(t, r.URL.Query()["status"], []string{"status"})
		_, _ = w.Write([]byte(`{"posts": [{"id": "1","post_type":"incident","status_page":{"id": "1","name":"MyStatusPage","published_at":"2024-02-12T09:23:23Z","status_page_type":"public","url":"https://mypagerduty"},"title":"MyPost","starts_at":"2024-02-12T09:23:23Z","ends_at":"2024-02-12T09:23:23Z"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := ListStatusPagePostOptions{
		PostType:       "incident",
		ReviewedStatus: "approved",
		Status:         []string{"status"},
	}

	res, err := client.ListStatusPagePosts("1", opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPagePostsResponse{
		APIListObject: APIListObject{},
		StatusPagePosts: []StatusPagePost{
			{
				ID:       "1",
				PostType: "incident",
				StatusPage: StatusPage{
					ID:             "1",
					Name:           "MyStatusPage",
					PublishedAt:    "2024-02-12T09:23:23Z",
					StatusPageType: "public",
					URL:            "https://mypagerduty",
				},
				Title:    "MyPost",
				StartsAt: "2024-02-12T09:23:23Z",
				EndsAt:   "2024-02-12T09:23:23Z",
			},
		},
	}

	testEqual(t, want, res)
}

// CreateStatusPagePost
func TestStatusPage_CreatePost(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"post": {"id": "1","post_type":"incident","status_page":{"id": "1","type":"status_page"},"title":"MyPost","starts_at":"2024-02-12T09:23:23Z","ends_at":"2024-02-12T09:23:23Z"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := StatusPagePost{
		PostType: "incident",
		StatusPage: StatusPage{
			ID:             "1",
			Name:           "MyStatusPage",
			PublishedAt:    "2024-02-12T09:23:23Z",
			StatusPageType: "public",
			URL:            "https://mypagerduty",
		},
		Title:    "MyPost",
		StartsAt: "2024-02-12T09:23:23Z",
		EndsAt:   "2024-02-12T09:23:23Z",
	}
	res, err := client.CreateStatusPagePost("1", input)
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPagePost{
		ID:       "1",
		PostType: "incident",
		StatusPage: StatusPage{
			ID:   "1",
			Type: "status_page",
		},
		Title:    "MyPost",
		StartsAt: "2024-02-12T09:23:23Z",
		EndsAt:   "2024-02-12T09:23:23Z",
	}

	testEqual(t, want, res)
}

// GetStatusPagePost
func TestStatusPage_GetPost(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"post": {"id": "1","post_type":"incident","status_page":{"id": "1","type":"status_page"},"title":"MyPost","starts_at":"2024-02-12T09:23:23Z","ends_at":"2024-02-12T09:23:23Z"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetStatusPagePost("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPagePost{
		ID:       "1",
		PostType: "incident",
		StatusPage: StatusPage{
			ID:   "1",
			Type: "status_page",
		},
		Title:    "MyPost",
		StartsAt: "2024-02-12T09:23:23Z",
		EndsAt:   "2024-02-12T09:23:23Z",
	}

	testEqual(t, want, res)
}

// UpdateStatusPagePost
func TestStatusPage_UpdatePost(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"post": {"id": "1","post_type":"incident","status_page":{"id": "1","type":"status_page"},"title":"MyPost","starts_at":"2024-02-12T09:23:23Z","ends_at":"2024-02-12T09:23:23Z"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := StatusPagePost{
		PostType: "incident",
		StatusPage: StatusPage{
			ID:             "1",
			Name:           "MyStatusPage",
			PublishedAt:    "2024-02-12T09:23:23Z",
			StatusPageType: "public",
			URL:            "https://mypagerduty",
		},
		Title:    "MyPost",
		StartsAt: "2024-02-12T09:23:23Z",
		EndsAt:   "2024-02-12T09:23:23Z",
	}
	res, err := client.UpdateStatusPagePost("1", "1", input)
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPagePost{
		ID:       "1",
		PostType: "incident",
		StatusPage: StatusPage{
			ID:   "1",
			Type: "status_page",
		},
		Title:    "MyPost",
		StartsAt: "2024-02-12T09:23:23Z",
		EndsAt:   "2024-02-12T09:23:23Z",
	}

	testEqual(t, want, res)
}

// DeleteStatusPagePost
func TestStatusPage_DeletePost(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")

	err := client.DeleteStatusPagePost("1", "1")
	if err != nil {
		t.Fatal(err)
	}
}

// ListStatusPagePostUpdates
func TestStatusPage_ListPostUpdates(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/post_updates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["reviewed_status"], []string{"approved"})
		_, _ = w.Write([]byte(`{
			"post_updates": [
				{
					"id":"1", "message":"Hello world", "reviewed_status":"approved", "notify_subscribers":false, "impacted_services": [
					{
						"impact":{
							"id":"1"," type":"status_page_impact"
						},
						"service":{
							"id":"1", "type":"status_page_service"
						}
					}],
					"post": {
						"id":"1", "type":"status_page_post"
					}
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := ListStatusPagePostUpdateOptions{
		ReviewedStatus: "approved",
	}

	res, err := client.ListStatusPagePostUpdates("1", "1", opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPagePostUpdatesResponse{
		APIListObject: APIListObject{},
		StatusPagePostUpdates: []StatusPagePostUpdate{
			{
				ID:             "1",
				Message:        "Hello world",
				ReviewedStatus: "approved",
				ImpactedServices: []StatusPagePostUpdateImpact{
					{
						Service: Service{
							APIObject: APIObject{
								ID:   "1",
								Type: "status_page_service",
							},
						},
					},
				},
				NotifySubscribers: false,
				Post: StatusPagePost{
					ID:   "1",
					Type: "status_page_post",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// CreateStatusPagePostUpdate
func TestStatusPage_CreatePostUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/post_updates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{
			"post_update": {
				"id": "1",
				"message": "Hello world",
				"reviewed_status": "approved",
				"notify_subscribers": false,
				"update_frequency_ms": 30000,
				"reported_at": "2024-02-12T09:23:23Z",
				"impacted_services": [
					{
						"severity": {
							"id": "1",
							"type": "status_page_severity"
						},
						"service": {
							"id": "1",
							"type": "status_page_service"
						}
					}
				],
				"status": {
					"id": "1",
					"type": "status_page_status"
				},
				"severity": {
					"id": "1",
					"type": "status_page_severity"
				},
				"post": {
					"id": "1",
					"type": "status_page_post"
				}
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := StatusPagePostUpdate{
		Message:           "Hello world",
		NotifySubscribers: false,
		ReportedAt:        "2024-02-12T09:23:23Z",
		UpdateFrequencyMS: 30000,
		ReviewedStatus:    "approved",
		Post: StatusPagePost{
			ID:   "1",
			Type: "status_page_post",
		},
		Status: StatusPageStatus{
			ID:   "1",
			Type: "status_page_status",
		},
		Severity: StatusPageSeverity{
			ID:   "1",
			Type: "status_page_severity",
		},
		ImpactedServices: []StatusPagePostUpdateImpact{
			{
				Service: Service{
					APIObject: APIObject{
						ID:   "1",
						Type: "status_page_service",
					},
				},
				Severity: StatusPageSeverity{
					ID:   "1",
					Type: "status_page_severity",
				},
			},
		},
	}
	res, err := client.CreateStatusPagePostUpdate("1", "1", input)
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPagePostUpdate{
		ID:                "1",
		Message:           "Hello world",
		NotifySubscribers: false,
		ReportedAt:        "2024-02-12T09:23:23Z",
		UpdateFrequencyMS: 30000,
		ReviewedStatus:    "approved",
		Post: StatusPagePost{
			ID:   "1",
			Type: "status_page_post",
		},
		Status: StatusPageStatus{
			ID:   "1",
			Type: "status_page_status",
		},
		Severity: StatusPageSeverity{
			ID:   "1",
			Type: "status_page_severity",
		},
		ImpactedServices: []StatusPagePostUpdateImpact{
			{
				Service: Service{
					APIObject: APIObject{
						ID:   "1",
						Type: "status_page_service",
					},
				},
				Severity: StatusPageSeverity{
					ID:   "1",
					Type: "status_page_severity",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// GetStatusPagePostUpdate
func TestStatusPage_GetPostUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/post_updates/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
			"post_update": {
				"id": "1",
				"message": "Hello world",
				"reviewed_status": "approved",
				"notify_subscribers": false,
				"update_frequency_ms": 30000,
				"reported_at": "2024-02-12T09:23:23Z",
				"impacted_services": [
					{
						"severity": {
							"id": "1",
							"type": "status_page_severity"
						},
						"service": {
							"id": "1",
							"type": "status_page_service"
						}
					}
				],
				"status": {
					"id": "1",
					"type": "status_page_status"
				},
				"severity": {
					"id": "1",
					"type": "status_page_severity"
				},
				"post": {
					"id": "1",
					"type": "status_page_post"
				}
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetStatusPagePostUpdate("1", "1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPagePostUpdate{
		ID:                "1",
		Message:           "Hello world",
		NotifySubscribers: false,
		ReportedAt:        "2024-02-12T09:23:23Z",
		UpdateFrequencyMS: 30000,
		ReviewedStatus:    "approved",
		Post: StatusPagePost{
			ID:   "1",
			Type: "status_page_post",
		},
		Status: StatusPageStatus{
			ID:   "1",
			Type: "status_page_status",
		},
		Severity: StatusPageSeverity{
			ID:   "1",
			Type: "status_page_severity",
		},
		ImpactedServices: []StatusPagePostUpdateImpact{
			{
				Service: Service{
					APIObject: APIObject{
						ID:   "1",
						Type: "status_page_service",
					},
				},
				Severity: StatusPageSeverity{
					ID:   "1",
					Type: "status_page_severity",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// UpdateStatusPagePostUpdate
func TestStatusPage_UpdatePostUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/post_updates/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{
			"post_update": {
				"id": "1",
				"message": "Hello world",
				"reviewed_status": "approved",
				"notify_subscribers": false,
				"update_frequency_ms": 30000,
				"reported_at": "2024-02-12T09:23:23Z",
				"impacted_services": [
					{
						"severity": {
							"id": "1",
							"type": "status_page_severity"
						},
						"service": {
							"id": "1",
							"type": "status_page_service"
						}
					}
				],
				"status": {
					"id": "1",
					"type": "status_page_status"
				},
				"severity": {
					"id": "1",
					"type": "status_page_severity"
				},
				"post": {
					"id": "1",
					"type": "status_page_post"
				}
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := StatusPagePostUpdate{
		Message:           "Hello world",
		NotifySubscribers: false,
		ReportedAt:        "2024-02-12T09:23:23Z",
		UpdateFrequencyMS: 30000,
		ReviewedStatus:    "approved",
		Post: StatusPagePost{
			ID:   "1",
			Type: "status_page_post",
		},
		Status: StatusPageStatus{
			ID:   "1",
			Type: "status_page_status",
		},
		Severity: StatusPageSeverity{
			ID:   "1",
			Type: "status_page_severity",
		},
		ImpactedServices: []StatusPagePostUpdateImpact{
			{
				Service: Service{
					APIObject: APIObject{
						ID:   "1",
						Type: "status_page_service",
					},
				},
				Severity: StatusPageSeverity{
					ID:   "1",
					Type: "status_page_severity",
				},
			},
		},
	}
	res, err := client.UpdateStatusPagePostUpdate("1", "1", "1", input)
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPagePostUpdate{
		ID:                "1",
		Message:           "Hello world",
		NotifySubscribers: false,
		ReportedAt:        "2024-02-12T09:23:23Z",
		UpdateFrequencyMS: 30000,
		ReviewedStatus:    "approved",
		Post: StatusPagePost{
			ID:   "1",
			Type: "status_page_post",
		},
		Status: StatusPageStatus{
			ID:   "1",
			Type: "status_page_status",
		},
		Severity: StatusPageSeverity{
			ID:   "1",
			Type: "status_page_severity",
		},
		ImpactedServices: []StatusPagePostUpdateImpact{
			{
				Service: Service{
					APIObject: APIObject{
						ID:   "1",
						Type: "status_page_service",
					},
				},
				Severity: StatusPageSeverity{
					ID:   "1",
					Type: "status_page_severity",
				},
			},
		},
	}

	testEqual(t, want, res)
}

// DeleteStatusPagePostUpdate
func TestStatusPage_DeletePostUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/post_updates/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")

	err := client.DeleteStatusPagePostUpdate("1", "1", "1")
	if err != nil {
		t.Fatal(err)
	}
}

// GetStatusPagePostPostmortem
func TestStatusPage_GetPostPostmortem(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/postmortem", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
			"postmortem": {
				"id": "1",
				"notify_subscribers": false,
				"message": "Hello world",
				"reported_at": "2024-02-12T09:23:23Z",
				"type": "status_page_post_postmortem"
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetStatusPagePostPostmortem("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &Postmortem{
		ID:                "1",
		NotifySubscribers: false,
		Message:           "Hello world",
		ReportedAt:        "2024-02-12T09:23:23Z",
		Type:              "status_page_post_postmortem",
	}

	testEqual(t, want, res)
}

// CreateStatusPagePostPostmortem
func TestStatusPage_CreatePostPostmortem(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/postmortem", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"postmortem": {"id": "1","message":"Hello world","notify_subscribers":false,"post":{"id": "1","type":"status_page_post"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := Postmortem{
		Message:           "Hello world",
		NotifySubscribers: false,
		Post: ShortPostType{
			ID:   "1",
			Type: "status_page_post",
		},
	}
	res, err := client.CreateStatusPagePostPostmortem("1", "1", input)
	if err != nil {
		t.Fatal(err)
	}
	want := &Postmortem{
		ID:                "1",
		Message:           "Hello world",
		NotifySubscribers: false,
		Post: ShortPostType{
			ID:   "1",
			Type: "status_page_post",
		},
	}

	testEqual(t, want, res)
}

// UpdateStatusPagePostPostmortem
func TestStatusPage_UpdatePostPostmortem(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/postmortem", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"postmortem": {"id": "1","message":"Hello world","notify_subscribers":false,"post":{"id": "1","type":"status_page_post"}}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := Postmortem{
		Message:           "Hello world",
		NotifySubscribers: false,
		Post: ShortPostType{
			ID:   "1",
			Type: "status_page_post",
		},
	}
	res, err := client.UpdateStatusPagePostPostmortem("1", "1", input)
	if err != nil {
		t.Fatal(err)
	}
	want := &Postmortem{
		ID:                "1",
		Message:           "Hello world",
		NotifySubscribers: false,
		Post: ShortPostType{
			ID:   "1",
			Type: "status_page_post",
		},
	}

	testEqual(t, want, res)
}

// DeleteStatusPagePostPostmortem
func TestStatusPage_DeletePostPostmortem(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/posts/1/postmortem", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")

	err := client.DeleteStatusPagePostPostmortem("1", "1")
	if err != nil {
		t.Fatal(err)
	}
}

// ListStatusPageSubscriptions
func TestStatusPage_ListSubscriptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testEqual(t, r.URL.Query()["channel"], []string{"email"})
		testEqual(t, r.URL.Query()["status"], []string{"active"})
		_, _ = w.Write([]byte(`{"subscriptions": [{"id": "1","channel":"email","contact":"address@email.example","status":"active","status_page":{"id": "1","type":"status_page"},"subscribable_object":{"id": "1","type":"status_page"},"type":"status_page_subscription"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	opts := ListStatusPageSubscriptionsOptions{
		Channel: "email",
		Status:  "active",
	}
	res, err := client.ListStatusPageSubscriptions("1", opts)
	if err != nil {
		t.Fatal(err)
	}
	want := &ListStatusPageSubscriptionsResponse{
		APIListObject: APIListObject{},
		StatusPageSubscriptions: []StatusPageSubscription{
			{
				ID:      "1",
				Channel: "email",
				Contact: "address@email.example",
				Status:  "active",
				StatusPage: StatusPage{
					ID:   "1",
					Type: "status_page",
				},
				SubscribableObject: SubscribableObject{
					ID:   "1",
					Type: "status_page",
				},
				Type: "status_page_subscription",
			},
		},
	}

	testEqual(t, want, res)
}

// CreateStatusPageSubscription
func TestStatusPage_CreateSubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"subscription": {"id": "1","channel":"email","contact":"address@email.example","status":"active","status_page":{"id": "1","type":"status_page"},"subscribable_object":{"id": "1","type":"status_page_service"},"type":"status_page_subscription"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := StatusPageSubscription{
		Channel: "email",
		Contact: "address@email.example",
		Status:  "active",
		StatusPage: StatusPage{
			ID:   "1",
			Type: "status_page",
		},
		SubscribableObject: SubscribableObject{
			ID:   "1",
			Type: "status_page",
		},
		Type: "status_page_subscription",
	}
	res, err := client.CreateStatusPageSubscription("1", input)
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPageSubscription{
		ID:      "1",
		Channel: "email",
		Contact: "address@email.example",
		Status:  "active",
		StatusPage: StatusPage{
			ID:   "1",
			Type: "status_page",
		},
		SubscribableObject: SubscribableObject{
			ID:   "1",
			Type: "status_page_service",
		},
		Type: "status_page_subscription",
	}

	testEqual(t, want, res)
}

// GetStatusPageSubscription
func TestStatusPage_GetSubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/status_pages/1/subscriptions/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"subscription": {"id": "1","channel":"email","contact":"address@email.example","status":"active","status_page":{"id": "1","type":"status_page"},"subscribable_object":{"id": "1","type":"status_page_service"},"type":"status_page_subscription"}}`))

	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetStatusPageSubscription("1", "1")
	if err != nil {
		t.Fatal(err)
	}
	want := &StatusPageSubscription{
		ID:      "1",
		Channel: "email",
		Contact: "address@email.example",
		Status:  "active",
		StatusPage: StatusPage{
			ID:   "1",
			Type: "status_page",
		},
		SubscribableObject: SubscribableObject{
			ID:   "1",
			Type: "status_page_service",
		},
		Type: "status_page_subscription",
	}

	testEqual(t, want, res)
}

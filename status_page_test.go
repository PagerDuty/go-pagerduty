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

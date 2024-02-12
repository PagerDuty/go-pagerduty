package pagerduty

import (
	"net/http"
	"testing"
)

// ListTags
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

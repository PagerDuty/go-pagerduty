package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

func TestStandard_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/standards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"standards":[{"active":true,"description":"A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.","id":"1","inclusions":[{"type":"technical_service_reference","id":"1"}],"name":"Service has a description","resource_type":"technical_service","type":"has_technical_service_description"}]}`))
	})

	var opts ListStandardsOptions
	client := defaultTestClient(server.URL, standardPath)

	res, err := client.ListStandardsWithContext(context.Background(), opts)

	want := &ListStandardsResponse{
		Standards: []Standard{
			{
				Active:      true,
				ID:          "1",
				Description: "A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.",
				Inclusions: []StandardInclusionExclusion{
					{
						ID:   "1",
						Type: "technical_service_reference",
					},
				},
				Name:         "Service has a description",
				ResourceType: "technical_service",
				Type:         "has_technical_service_description",
			},
		},
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestStandard_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/standards/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"active":true,"description":"A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.","id":"1","inclusions":[{"type":"technical_service_reference","id":"1"}],"name":"Service has a description","resource_type":"technical_service","type":"has_technical_service_description"}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &Standard{Active: true}
	want := &Standard{
		Active:      true,
		ID:          "1",
		Description: "A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.",
		Inclusions: []StandardInclusionExclusion{
			{
				ID:   "1",
				Type: "technical_service_reference",
			},
		},
		Name:         "Service has a description",
		ResourceType: "technical_service",
		Type:         "has_technical_service_description",
	}
	res, err := client.UpdateStandardWithContext(context.Background(), "1", *input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

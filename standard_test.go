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

	res, err := client.ListStandards(context.Background(), opts)

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
	res, err := client.UpdateStandard(context.Background(), "1", *input)
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestStandard_ListTechniServiceStandardScores(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/standards/scores/technical_services/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"resource_id":"1","resource_type":"technical_service","score":{"passing":1,"total":1},"standards":[{"active":true,"description":"A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.","id":"01CXX38Q0U8XKHO4LNKXUJTBFG","pass":true,"name":"Service has a description","type":"has_technical_service_description"}]}`))
	})

	client := defaultTestClient(server.URL, standardPath)

	res, err := client.ListResourceStandardScores(context.Background(), "1", "technical_services")

	want := &ResourceStandardScore{
		ResourceID:   "1",
		ResourceType: "technical_service",
		Score: &ResourceScore{
			Passing: 1,
			Total:   1,
		},
		Standards: []ResourceStandard{
			{
				Active:      true,
				Description: "A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.",
				ID:          "01CXX38Q0U8XKHO4LNKXUJTBFG",
				Pass:        true,
				Name:        "Service has a description",
				Type:        "has_technical_service_description",
			},
		},
	}
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestStandard_ListManyTechniServiceStandardScores(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/standards/scores/technical_services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"resources":[{"resource_id":"1","resource_type":"technical_service","score":{"passing":1,"total":1},"standards":[{"active":true,"description":"A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.","id":"01CXX38Q0U8XKHO4LNKXUJTBFG","pass":true,"name":"Service has a description","type":"has_technical_service_description"}]}]}`))
	})

	opts := ListMultiResourcesStandardScoresOptions{
		IDs: []string{"1"},
	}
	client := defaultTestClient(server.URL, standardPath)

	res, err := client.ListMultiResourcesStandardScores(context.Background(), "technical_services", opts)

	want := &ListMultiResourcesStandardScoresResponse{
		Resources: []ResourceStandardScore{
			{
				ResourceID:   "1",
				ResourceType: "technical_service",
				Score: &ResourceScore{
					Passing: 1,
					Total:   1,
				},
				Standards: []ResourceStandard{
					{
						Active:      true,
						Description: "A description provides critical context about what a service represents or is used for to inform team members and responders. The description should be kept concise and understandable by those without deep knowledge of the service.",
						ID:          "01CXX38Q0U8XKHO4LNKXUJTBFG",
						Pass:        true,
						Name:        "Service has a description",
						Type:        "has_technical_service_description",
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

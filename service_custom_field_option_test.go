package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// TestServiceCustomFieldOptions_List tests the ListServiceCustomFieldOptions method
func TestServiceCustomFieldOptions_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42/field_options", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"field_options": [
				{
					"created_at": "2023-01-01T00:00:00Z",
					"data": {
						"data_type": "string",
						"value": "production"
					},
					"id": "OPT1",
					"type": "field_option",
					"updated_at": "2023-01-01T00:00:00Z"
				},
				{
					"created_at": "2023-01-01T00:00:00Z",
					"data": {
						"data_type": "string",
						"value": "staging"
					},
					"id": "OPT2",
					"type": "field_option",
					"updated_at": "2023-01-01T00:00:00Z"
				},
				{
					"created_at": "2023-01-01T00:00:00Z",
					"data": {
						"data_type": "string",
						"value": "development"
					},
					"id": "OPT3",
					"type": "field_option",
					"updated_at": "2023-01-01T00:00:00Z"
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListServiceCustomFieldOptions(context.Background(), "PXPGF42")

	want := &ListServiceCustomFieldOptionsResponse{
		FieldOptions: []ServiceCustomFieldOption{
			{
				ID:        "OPT1",
				Type:      "field_option",
				CreatedAt: "2023-01-01T00:00:00Z",
				UpdatedAt: "2023-01-01T00:00:00Z",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "production",
				},
			},
			{
				ID:        "OPT2",
				Type:      "field_option",
				CreatedAt: "2023-01-01T00:00:00Z",
				UpdatedAt: "2023-01-01T00:00:00Z",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "staging",
				},
			},
			{
				ID:        "OPT3",
				Type:      "field_option",
				CreatedAt: "2023-01-01T00:00:00Z",
				UpdatedAt: "2023-01-01T00:00:00Z",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "development",
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFieldOptions_Get tests the GetServiceCustomFieldOption method
func TestServiceCustomFieldOptions_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42/field_options/OPT1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"field_option": {
				"created_at": "2023-01-01T00:00:00Z",
				"data": {
					"data_type": "string",
					"value": "production"
				},
				"id": "OPT1",
				"type": "field_option",
				"updated_at": "2023-01-01T00:00:00Z"
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetServiceCustomFieldOption(context.Background(), "PXPGF42", "OPT1")

	want := &ServiceCustomFieldOption{
		ID:        "OPT1",
		Type:      "field_option",
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-01T00:00:00Z",
		Data: ServiceCustomFieldOptionData{
			DataType: "string",
			Value:    "production",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFieldOptions_Create tests the CreateServiceCustomFieldOption method
func TestServiceCustomFieldOptions_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42/field_options", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"field_option": {
				"created_at": "2023-01-01T00:00:00Z",
				"data": {
					"data_type": "string",
					"value": "production"
				},
				"id": "OPT1",
				"type": "field_option",
				"updated_at": "2023-01-01T00:00:00Z"
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ServiceCustomFieldOption{
		Data: ServiceCustomFieldOptionData{
			DataType: "string",
			Value:    "production",
		},
	}
	res, err := client.CreateServiceCustomFieldOption(context.Background(), "PXPGF42", input)

	want := &ServiceCustomFieldOption{
		ID:        "OPT1",
		Type:      "field_option",
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-01T00:00:00Z",
		Data: ServiceCustomFieldOptionData{
			DataType: "string",
			Value:    "production",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFieldOptions_Update tests the UpdateServiceCustomFieldOption method
func TestServiceCustomFieldOptions_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42/field_options/OPT1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"field_option": {
				"created_at": "2023-01-01T00:00:00Z",
				"data": {
					"data_type": "string",
					"value": "prod"
				},
				"id": "OPT1",
				"type": "field_option",
				"updated_at": "2023-01-02T00:00:00Z"
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ServiceCustomFieldOption{
		ID: "OPT1",
		Data: ServiceCustomFieldOptionData{
			DataType: "string",
			Value:    "prod",
		},
	}
	res, err := client.UpdateServiceCustomFieldOption(context.Background(), "PXPGF42", input)

	want := &ServiceCustomFieldOption{
		ID:        "OPT1",
		Type:      "field_option",
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-02T00:00:00Z",
		Data: ServiceCustomFieldOptionData{
			DataType: "string",
			Value:    "prod",
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFieldOptions_Delete tests the DeleteServiceCustomFieldOption method
func TestServiceCustomFieldOptions_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42/field_options/OPT1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}

		w.WriteHeader(http.StatusNoContent)
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteServiceCustomFieldOption(context.Background(), "PXPGF42", "OPT1")

	if err != nil {
		t.Fatal(err)
	}
}

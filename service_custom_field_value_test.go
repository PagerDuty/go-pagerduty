package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// TestServiceCustomFieldValues_Get tests the GetServiceCustomFieldValues method
func TestServiceCustomFieldValues_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/PXPGF42/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		
		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}
		
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"custom_fields": [
				{
					"data_type": "string",
					"display_name": "Environment",
					"field_type": "single_value_fixed",
					"id": "PT4KHEE",
					"name": "environment",
					"type": "field_value",
					"value": "production"
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetServiceCustomFieldValues(context.Background(), "PXPGF42")

	want := &ListServiceCustomFieldValuesResponse{
		CustomFields: []ServiceCustomFieldValue{
			{
				ID:          "PT4KHEE",
				Name:        "environment",
				DisplayName: "Environment",
				DataType:    "string",
				FieldType:   "single_value_fixed",
				Type:        "field_value",
				Value:       "production",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFieldValues_GetMultiValue tests the GetServiceCustomFieldValues method with multi-value fields
func TestServiceCustomFieldValues_GetMultiValue(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/PXPGF42/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		
		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}
		
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"custom_fields": [
				{
					"data_type": "string",
					"display_name": "Environment",
					"field_type": "multi_value_fixed",
					"id": "PT4KHEE",
					"name": "environment",
					"type": "field_value",
					"value": ["production", "staging"]
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetServiceCustomFieldValues(context.Background(), "PXPGF42")

	want := &ListServiceCustomFieldValuesResponse{
		CustomFields: []ServiceCustomFieldValue{
			{
				ID:          "PT4KHEE",
				Name:        "environment",
				DisplayName: "Environment",
				DataType:    "string",
				FieldType:   "multi_value_fixed",
				Type:        "field_value",
				Value:       []interface{}{"production", "staging"},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFieldValues_Update tests the UpdateServiceCustomFieldValues method
func TestServiceCustomFieldValues_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/PXPGF42/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		
		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"custom_fields": [
				{
					"data_type": "string",
					"display_name": "Environment",
					"field_type": "single_value_fixed",
					"id": "PT4KHEE",
					"name": "environment",
					"type": "field_value",
					"value": "staging"
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ListServiceCustomFieldValuesResponse{
		CustomFields: []ServiceCustomFieldValue{
			{
				ID:    "PT4KHEE",
				Value: "staging",
			},
		},
	}
	res, err := client.UpdateServiceCustomFieldValues(context.Background(), "PXPGF42", input)

	want := &ListServiceCustomFieldValuesResponse{
		CustomFields: []ServiceCustomFieldValue{
			{
				ID:          "PT4KHEE",
				Name:        "environment",
				DisplayName: "Environment",
				DataType:    "string",
				FieldType:   "single_value_fixed",
				Type:        "field_value",
				Value:       "staging",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFieldValues_UpdateMultiValue tests the UpdateServiceCustomFieldValues method with multi-value fields
func TestServiceCustomFieldValues_UpdateMultiValue(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/PXPGF42/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		
		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"custom_fields": [
				{
					"data_type": "string",
					"display_name": "Environment",
					"field_type": "multi_value_fixed",
					"id": "PT4KHEE",
					"name": "environment",
					"type": "field_value",
					"value": ["staging", "development"]
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ListServiceCustomFieldValuesResponse{
		CustomFields: []ServiceCustomFieldValue{
			{
				ID:    "PT4KHEE",
				Value: []string{"staging", "development"},
			},
		},
	}
	res, err := client.UpdateServiceCustomFieldValues(context.Background(), "PXPGF42", input)

	want := &ListServiceCustomFieldValuesResponse{
		CustomFields: []ServiceCustomFieldValue{
			{
				ID:          "PT4KHEE",
				Name:        "environment",
				DisplayName: "Environment",
				DataType:    "string",
				FieldType:   "multi_value_fixed",
				Type:        "field_value",
				Value:       []interface{}{"staging", "development"},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

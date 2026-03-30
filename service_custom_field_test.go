package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// TestServiceCustomFields_List tests the ListServiceCustomFields method
func TestServiceCustomFields_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"fields": [
				{
					"id": "PXPGF42",
					"type": "service_custom_field",
					"name": "environment",
					"display_name": "Environment",
					"description": "The environment of the service",
					"data_type": "string",
					"field_type": "single_value_fixed",
					"default_value": {
						"value": "production"
					},
					"field_options": [
						{
							"id": "PXPGF42",
							"type": "service_custom_field_option",
							"data": {
								"data_type": "string",
								"value": "production"
							}
						},
						{
							"id": "PXPGF43",
							"type": "service_custom_field_option",
							"data": {
								"data_type": "string",
								"value": "staging"
							}
						}
					],
					"created_at": "2023-07-11T16:42:33Z",
					"updated_at": "2023-07-11T16:42:33Z"
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	o := ListServiceCustomFieldsOptions{}
	res, err := client.ListServiceCustomFields(context.Background(), o)

	want := &ListServiceCustomFieldsResponse{
		Fields: []ServiceCustomField{
			{
				APIObject: APIObject{
					ID:   "PXPGF42",
					Type: "service_custom_field",
				},
				Name:        "environment",
				DisplayName: "Environment",
				Description: "The environment of the service",
				DataType:    "string",
				FieldType:   "single_value_fixed",
				DefaultValue: map[string]interface{}{
					"value": "production",
				},
				FieldOptions: []ServiceCustomFieldOption{
					{
						ID:   "PXPGF42",
						Type: "service_custom_field_option",
						Data: ServiceCustomFieldOptionData{
							DataType: "string",
							Value:    "production",
						},
					},
					{
						ID:   "PXPGF43",
						Type: "service_custom_field_option",
						Data: ServiceCustomFieldOptionData{
							DataType: "string",
							Value:    "staging",
						},
					},
				},
				CreatedAt: "2023-07-11T16:42:33Z",
				UpdatedAt: "2023-07-11T16:42:33Z",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFields_Get tests the GetServiceCustomField method
func TestServiceCustomFields_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"field": {
				"id": "PXPGF42",
				"type": "service_custom_field",
				"name": "environment",
				"display_name": "Environment",
				"description": "The environment of the service",
				"data_type": "string",
				"field_type": "single_value_fixed",
				"default_value": {
					"value": "production"
				},
				"field_options": [
					{
						"id": "PXPGF42",
						"type": "service_custom_field_option",
						"data": {
							"data_type": "string",
							"value": "production"
						}
					},
					{
						"id": "PXPGF43",
						"type": "service_custom_field_option",
						"data": {
							"data_type": "string",
							"value": "staging"
						}
					}
				],
				"created_at": "2023-07-11T16:42:33Z",
				"updated_at": "2023-07-11T16:42:33Z"
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	o := ListServiceCustomFieldsOptions{}
	res, err := client.GetServiceCustomField(context.Background(), "PXPGF42", o)

	want := &ServiceCustomField{
		APIObject: APIObject{
			ID:   "PXPGF42",
			Type: "service_custom_field",
		},
		Name:        "environment",
		DisplayName: "Environment",
		Description: "The environment of the service",
		DataType:    "string",
		FieldType:   "single_value_fixed",
		DefaultValue: map[string]interface{}{
			"value": "production",
		},
		FieldOptions: []ServiceCustomFieldOption{
			{
				ID:   "PXPGF42",
				Type: "service_custom_field_option",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "production",
				},
			},
			{
				ID:   "PXPGF43",
				Type: "service_custom_field_option",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "staging",
				},
			},
		},
		CreatedAt: "2023-07-11T16:42:33Z",
		UpdatedAt: "2023-07-11T16:42:33Z",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFields_Create tests the CreateServiceCustomField method
func TestServiceCustomFields_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"field": {
				"id": "PXPGF42",
				"type": "service_custom_field",
				"name": "environment",
				"display_name": "Environment",
				"description": "The environment of the service",
				"data_type": "string",
				"field_type": "single_value_fixed",
				"default_value": {
					"value": "production"
				},
				"field_options": [
					{
						"id": "PXPGF42",
						"type": "service_custom_field_option",
						"data": {
							"data_type": "string",
							"value": "production"
						}
					},
					{
						"id": "PXPGF43",
						"type": "service_custom_field_option",
						"data": {
							"data_type": "string",
							"value": "staging"
						}
					}
				],
				"created_at": "2023-07-11T16:42:33Z",
				"updated_at": "2023-07-11T16:42:33Z"
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ServiceCustomField{
		Name:        "environment",
		DisplayName: "Environment",
		Description: "The environment of the service",
		DataType:    "string",
		FieldType:   "single_value_fixed",
		DefaultValue: map[string]interface{}{
			"value": "production",
		},
		FieldOptions: []ServiceCustomFieldOption{
			{
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "production",
				},
			},
			{
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "staging",
				},
			},
		},
	}
	res, err := client.CreateServiceCustomField(context.Background(), input)

	want := &ServiceCustomField{
		APIObject: APIObject{
			ID:   "PXPGF42",
			Type: "service_custom_field",
		},
		Name:        "environment",
		DisplayName: "Environment",
		Description: "The environment of the service",
		DataType:    "string",
		FieldType:   "single_value_fixed",
		DefaultValue: map[string]interface{}{
			"value": "production",
		},
		FieldOptions: []ServiceCustomFieldOption{
			{
				ID:   "PXPGF42",
				Type: "service_custom_field_option",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "production",
				},
			},
			{
				ID:   "PXPGF43",
				Type: "service_custom_field_option",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "staging",
				},
			},
		},
		CreatedAt: "2023-07-11T16:42:33Z",
		UpdatedAt: "2023-07-11T16:42:33Z",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFields_Update tests the UpdateServiceCustomField method
func TestServiceCustomFields_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"field": {
				"id": "PXPGF42",
				"type": "service_custom_field",
				"name": "environment",
				"display_name": "Environment Updated",
				"description": "The environment of the service (updated)",
				"data_type": "string",
				"field_type": "single_value_fixed",
				"default_value": {
					"value": "staging"
				},
				"field_options": [
					{
						"id": "PXPGF42",
						"type": "service_custom_field_option",
						"data": {
							"data_type": "string",
							"value": "production"
						}
					},
					{
						"id": "PXPGF43",
						"type": "service_custom_field_option",
						"data": {
							"data_type": "string",
							"value": "staging"
						}
					}
				],
				"created_at": "2023-07-11T16:42:33Z",
				"updated_at": "2023-07-12T10:15:20Z"
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &ServiceCustomField{
		APIObject: APIObject{
			ID:   "PXPGF42",
			Type: "service_custom_field",
		},
		Name:        "environment",
		DisplayName: "Environment Updated",
		Description: "The environment of the service (updated)",
		DataType:    "string",
		FieldType:   "single_value_fixed",
		DefaultValue: map[string]interface{}{
			"value": "staging",
		},
	}
	res, err := client.UpdateServiceCustomField(context.Background(), input)

	want := &ServiceCustomField{
		APIObject: APIObject{
			ID:   "PXPGF42",
			Type: "service_custom_field",
		},
		Name:        "environment",
		DisplayName: "Environment Updated",
		Description: "The environment of the service (updated)",
		DataType:    "string",
		FieldType:   "single_value_fixed",
		DefaultValue: map[string]interface{}{
			"value": "staging",
		},
		FieldOptions: []ServiceCustomFieldOption{
			{
				ID:   "PXPGF42",
				Type: "service_custom_field_option",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "production",
				},
			},
			{
				ID:   "PXPGF43",
				Type: "service_custom_field_option",
				Data: ServiceCustomFieldOptionData{
					DataType: "string",
					Value:    "staging",
				},
			},
		},
		CreatedAt: "2023-07-11T16:42:33Z",
		UpdatedAt: "2023-07-12T10:15:20Z",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// TestServiceCustomFields_Delete tests the DeleteServiceCustomField method
func TestServiceCustomFields_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/custom_fields/PXPGF42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		// Verify that the X-EARLY-ACCESS header is present with the correct value
		if r.Header.Get("X-EARLY-ACCESS") != "service-custom-fields-preview" {
			t.Errorf("Expected X-EARLY-ACCESS header to be 'service-custom-fields-preview', got %s", r.Header.Get("X-EARLY-ACCESS"))
		}

		w.WriteHeader(204)
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteServiceCustomField(context.Background(), "PXPGF42")

	if err != nil {
		t.Fatal(err)
	}
}

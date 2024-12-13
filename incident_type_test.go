package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// Incident Type

func TestIncidentType_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"incident_types": [{"id": "PT1234", "name": "Test Type"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")
	want := &ListIncidentTypesResponse{
		IncidentTypes: []IncidentType{
			{
				ID:   "PT1234",
				Name: "Test Type",
			},
		},
	}

	res, err := client.ListIncidentTypes(context.Background(), ListIncidentTypesOptions{})
	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestIncidentType_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"incident_type": {"id": "PT1234", "name": "New Type", "enabled": true}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	enabled := true
	input := CreateIncidentTypeOptions{
		Name:        "New Type",
		DisplayName: "New Type Display",
		ParentType:  "incident",
		Enabled:     &enabled,
	}

	res, err := client.CreateIncidentType(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentType{
		ID:      "PT1234",
		Name:    "New Type",
		Enabled: true,
	}

	testEqual(t, want, res)
}

func TestIncidentType_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"incident_type": {"id": "PT1234", "name": "Test Type"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetIncidentType(context.Background(), "PT1234", GetIncidentTypeOptions{})
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentType{
		ID:   "PT1234",
		Name: "Test Type",
	}

	testEqual(t, want, res)
}

func TestIncidentType_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"incident_type": {"id": "PT1234", "name": "Updated Type", "enabled": false}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	displayName := "Updated Type"
	enabled := false
	input := UpdateIncidentTypeOptions{
		DisplayName: &displayName,
		Enabled:     &enabled,
	}

	res, err := client.UpdateIncidentType(context.Background(), "PT1234", input)
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentType{
		ID:      "PT1234",
		Name:    "Updated Type",
		Enabled: false,
	}

	testEqual(t, want, res)
}

// Incident Type Field

func TestIncidentTypeField_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"fields": [{"id": "PF1234", "name": "Test Field"}]}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListIncidentTypeFields(context.Background(), "PT1234", ListIncidentTypeFieldsOptions{})
	if err != nil {
		t.Fatal(err)
	}

	want := &ListIncidentTypeFieldsResponse{
		Fields: []IncidentTypeField{
			{
				ID:   "PF1234",
				Name: "Test Field",
			},
		},
	}

	testEqual(t, want, res)
}

func TestIncidentTypeField_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{
			"field": {
				"id": "PF1234",
				"name": "Test Field",
				"display_name": "Test Field Display",
				"data_type": "string",
				"field_type": "single_value",
				"enabled": true
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	enabled := true
	input := CreateIncidentTypeFieldOptions{
		Name:        "Test Field",
		DisplayName: "Test Field Display",
		DataType:    "string",
		FieldType:   "single_value",
		Enabled:     &enabled,
	}

	res, err := client.CreateIncidentTypeField(context.Background(), "PT1234", input)
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentTypeField{
		ID:          "PF1234",
		Name:        "Test Field",
		DisplayName: "Test Field Display",
		DataType:    "string",
		FieldType:   "single_value",
		Enabled:     true,
	}

	testEqual(t, want, res)
}

func TestIncidentTypeField_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
			"field": {
				"id": "PF1234",
				"name": "Test Field",
				"display_name": "Test Field Display",
				"data_type": "string",
				"field_type": "single_value",
				"enabled": true
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetIncidentTypeField(context.Background(), "PT1234", "PF1234", GetIncidentTypeFieldOptions{})
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentTypeField{
		ID:          "PF1234",
		Name:        "Test Field",
		DisplayName: "Test Field Display",
		DataType:    "string",
		FieldType:   "single_value",
		Enabled:     true,
	}

	testEqual(t, want, res)
}

func TestIncidentTypeField_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{
			"field": {
				"id": "PF1234",
				"name": "Test Field",
				"display_name": "Updated Field Display",
				"data_type": "string",
				"field_type": "single_value",
				"enabled": false
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	displayName := "Updated Field Display"
	enabled := false
	input := UpdateIncidentTypeFieldOptions{
		DisplayName: &displayName,
		Enabled:     &enabled,
	}

	res, err := client.UpdateIncidentTypeField(context.Background(), "PT1234", "PF1234", input)
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentTypeField{
		ID:          "PF1234",
		Name:        "Test Field",
		DisplayName: "Updated Field Display",
		DataType:    "string",
		FieldType:   "single_value",
		Enabled:     false,
	}

	testEqual(t, want, res)
}

func TestIncidentTypeField_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	client := defaultTestClient(server.URL, "foo")

	err := client.DeleteIncidentTypeField(context.Background(), "PT1234", "PF1234")
	if err != nil {
		t.Fatal(err)
	}
}

// Incident Type Field Option

func TestIncidentTypeFieldOption_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234/field_options", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
			"field_options": [
				{
					"id": "PFO123",
					"type": "field_option",
					"data": {
						"value": "Option 1",
						"data_type": "string"
					}
				}
			]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListIncidentTypeFieldOptions(context.Background(), "PT1234", "PF1234", ListIncidentTypeFieldOptionsOptions{})
	if err != nil {
		t.Fatal(err)
	}

	want := &ListIncidentTypeFieldOptionsResponse{
		FieldOptions: []IncidentTypeFieldOption{
			{
				ID:   "PFO123",
				Type: "field_option",
				Data: &IncidentTypeFieldOptionData{
					Value:    "Option 1",
					DataType: "string",
				},
			},
		},
	}

	testEqual(t, want, res)
}

func TestIncidentTypeFieldOption_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234/field_options", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{
			"field_option": {
				"id": "PFO123",
				"type": "field_option",
				"data": {
					"value": "New Option",
					"data_type": "string"
				}
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := CreateIncidentTypeFieldOptionPayload{
		Data: &IncidentTypeFieldOptionData{
			Value:    "New Option",
			DataType: "string",
		},
	}

	res, err := client.CreateIncidentTypeFieldOption(context.Background(), "PT1234", "PF1234", input)
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentTypeFieldOption{
		ID:   "PFO123",
		Type: "field_option",
		Data: &IncidentTypeFieldOptionData{
			Value:    "New Option",
			DataType: "string",
		},
	}

	testEqual(t, want, res)
}

func TestIncidentTypeFieldOption_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234/field_options/PFO123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
			"field_option": {
				"id": "PFO123",
				"type": "field_option",
				"data": {
					"value": "Option 1",
					"data_type": "string"
				}
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.GetIncidentTypeFieldOption(context.Background(), "PT1234", "PF1234", "PFO123", GetIncidentTypeFieldOptionOptions{})
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentTypeFieldOption{
		ID:   "PFO123",
		Type: "field_option",
		Data: &IncidentTypeFieldOptionData{
			Value:    "Option 1",
			DataType: "string",
		},
	}

	testEqual(t, want, res)
}

func TestIncidentTypeFieldOption_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234/field_options/PFO123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{
			"field_option": {
				"id": "PFO123",
				"type": "field_option",
				"data": {
					"value": "Updated Option",
					"data_type": "string"
				}
			}
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := UpdateIncidentTypeFieldOptionPayload{
		ID: "PFO123",
		Data: &IncidentTypeFieldOptionData{
			Value:    "Updated Option",
			DataType: "string",
		},
	}

	res, err := client.UpdateIncidentTypeFieldOption(context.Background(), "PT1234", "PF1234", input)
	if err != nil {
		t.Fatal(err)
	}

	want := &IncidentTypeFieldOption{
		ID:   "PFO123",
		Type: "field_option",
		Data: &IncidentTypeFieldOptionData{
			Value:    "Updated Option",
			DataType: "string",
		},
	}

	testEqual(t, want, res)
}

func TestIncidentTypeFieldOption_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/incidents/types/PT1234/custom_fields/PF1234/field_options/PFO123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	client := defaultTestClient(server.URL, "foo")

	err := client.DeleteIncidentTypeFieldOption(context.Background(), "PT1234", "PF1234", "PFO123")
	if err != nil {
		t.Fatal(err)
	}
}

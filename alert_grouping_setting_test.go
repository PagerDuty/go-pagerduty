package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// Create AlertGroupingSetting
func TestAlertGroupingSetting_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alert_grouping_settings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{
			"alert_grouping_setting": {
				"id": "PZC4OM1",
				"name": "Example of Alert Grouping Setting",
				"description": "This Alert Grouping Setting is an example",
				"type": "content_based",
				"config": {
					"time_window": 900,
					"aggregate": "all",
					"fields": [
						"summary",
						"component",
						"custom_details.host",
						"custom_details.field1.field2"
					]
				},
				"services": [
					{"id": "P0KJZ0A"},
					{"id": "PA15YRT"}
				]
			}
		}`))
	})

	ctx := context.Background()
	client := defaultTestClient(server.URL, "foo")
	input := AlertGroupingSetting{
		Name: "foo",
	}
	res, err := client.CreateAlertGroupingSetting(ctx, input)

	want := &AlertGroupingSetting{
		ID:          "PZC4OM1",
		Name:        "Example of Alert Grouping Setting",
		Description: "This Alert Grouping Setting is an example",
		Type:        "content_based",
		Config: AlertGroupingSettingConfigContentBased{
			TimeWindow: uint(900),
			Aggregate:  "all",
			Fields: []string{
				"summary",
				"component",
				"custom_details.host",
				"custom_details.field1.field2",
			},
		},
		Services: []AlertGroupingSettingService{
			{ID: "P0KJZ0A"},
			{ID: "PA15YRT"},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// ListAlertGroupingSettings
func TestAlertGroupingSettings_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alert_grouping_settings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
			"alert_grouping_settings": [{
				"id": "PZC4OM1",
				"name": "Example of Alert Grouping Setting",
				"description": "This Alert Grouping Setting is an example",
				"type": "time",
				"config": {"timeout": 60},
				"services": [
					{"id": "P0KJZ0A"},
					{"id": "PA15YRT"}
				]
			}]
		}`))
	})

	client := defaultTestClient(server.URL, "foo")
	o := ListAlertGroupingSettingsOptions{}
	res, err := client.ListAlertGroupingSettings(context.Background(), o)

	want := &ListAlertGroupingSettingsResponse{
		AlertGroupingSettings: []AlertGroupingSetting{
			{
				ID:          "PZC4OM1",
				Name:        "Example of Alert Grouping Setting",
				Description: "This Alert Grouping Setting is an example",
				Type:        "time",
				Config: AlertGroupingSettingConfigTime{
					Timeout: uint(60),
				},
				Services: []AlertGroupingSettingService{
					{ID: "P0KJZ0A"},
					{ID: "PA15YRT"},
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// GetAlertGroupingSettings
func TestAlertGroupingSettings_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alert_grouping_settings/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"alert_grouping_setting": {"id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.GetAlertGroupingSetting(context.Background(), "1")

	want := &AlertGroupingSetting{
		ID: "1",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// UpdateAlertGroupingSettings
func TestAlertGroupingSettings_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alert_grouping_settings/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"alert_grouping_setting": {"id": "1"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	want := AlertGroupingSetting{ID: "1"}
	res, err := client.UpdateAlertGroupingSetting(context.Background(), want)

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, &want, res)
}

// DeleteAlertGroupingSettings
func TestAlertGroupingSettings_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alert_grouping_settings/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(204)
	})

	client := defaultTestClient(server.URL, "foo")
	err := client.DeleteAlertGroupingSetting(context.Background(), "1")

	if err != nil {
		t.Fatal(err)
	}
}

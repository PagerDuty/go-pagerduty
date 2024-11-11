package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"
)

// AlertGroupingSettingConfigTime is the configuration content for a
// AlertGroupingSetting of type "time"
type AlertGroupingSettingConfigTime struct {
	Timeout uint `json:"timeout"`
}

// AlertGroupingSettingConfigIntelligent is the configuration content for a
// AlertGroupingSetting of type "intelligent"
type AlertGroupingSettingConfigIntelligent struct {
	TimeWindow            uint  `json:"time_window"`
	RecommendedTimeWindow *uint `json:"recommended_time_window,omitempty"`
}

// AlertGroupingSettingConfigContentBased is the configuration content for a
// AlertGroupingSetting of type "content_based" or "content_based_intelligent"
type AlertGroupingSettingConfigContentBased struct {
	TimeWindow            uint     `json:"time_window"`
	RecommendedTimeWindow *uint    `json:"recommended_time_window,omitempty"`
	Aggregate             string   `json:"aggregate"`
	Fields                []string `json:"fields"`
}

type AlertGroupingSettingType string

const (
	AlertGroupingSettingContentBasedType            AlertGroupingSettingType = "content_based"
	AlertGroupingSettingContentBasedIntelligentType AlertGroupingSettingType = "content_based_intelligent"
	AlertGroupingSettingIntelligentType             AlertGroupingSettingType = "intelligent"
	AlertGroupingSettingTimeType                    AlertGroupingSettingType = "time"
)

// AlertGroupingSetting is a configuration used during the grouping of the
// alerts easier to reuse and share between many services
type AlertGroupingSetting struct {
	ID          string                        `json:"id,omitempty"`
	Name        string                        `json:"name,omitempty"`
	Description string                        `json:"description,omitempty"`
	Type        AlertGroupingSettingType      `json:"type,omitempty"`
	Config      interface{}                   `json:"config,omitempty"`
	Services    []AlertGroupingSettingService `json:"services"`
}

// AlertGroupingSettingService is a reference to services associated with an
// alert grouping setting.
type AlertGroupingSettingService struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// CreateAlertGroupingSetting creates an instance of AlertGroupingSettings for
// either one service or many services that are in the alert group setting
func (c *Client) CreateAlertGroupingSetting(ctx context.Context, a AlertGroupingSetting) (*AlertGroupingSetting, error) {
	d := map[string]AlertGroupingSetting{"alert_grouping_setting": a}

	resp, err := c.post(ctx, "/alert_grouping_settings", d, nil)
	if err != nil {
		return nil, err
	}

	var resultRaw struct {
		AlertGroupingSetting alertGroupingSettingRaw `json:"alert_grouping_setting"`
	}
	if err := c.decodeJSON(resp, &resultRaw); err != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", err)
	}

	value := resultRaw.AlertGroupingSetting
	return getAlertGroupingSettingFromRaw(value)
}

// ListAlertGroupingSettingsOptions is the data structure used when calling the
// ListAlertGroupingSettings API endpoint.
type ListAlertGroupingSettingsOptions struct {
	After      string   `url:"after,omitempty"`
	Before     string   `url:"before,omitempty"`
	Limit      uint     `url:"limit,omitempty"`
	Total      bool     `url:"total,omitempty"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
}

// ListAlertGroupingSettingsResponse is the data structure returned from calling
// the ListAlertGroupingSettingsResponse API endpoint.
type ListAlertGroupingSettingsResponse struct {
	After                 string                 `json:"after,omitempty"`
	Before                string                 `json:"before,omitempty"`
	Limit                 uint                   `json:"limit,omitempty"`
	Total                 bool                   `json:"total,omitempty"`
	AlertGroupingSettings []AlertGroupingSetting `json:"alert_grouping_settings"`
}

// ListAlertGroupingSettings lists all of your alert grouping settings including
// both single service settings and global content based settings.
func (c *Client) ListAlertGroupingSettings(ctx context.Context, o ListAlertGroupingSettingsOptions) (*ListAlertGroupingSettingsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/alert_grouping_settings?"+v.Encode(), nil)

	// If there are no alert grouping settings, return an empty response.
	if resp.StatusCode == 404 {
		return &ListAlertGroupingSettingsResponse{}, nil
	}

	if err != nil {
		return nil, err
	}

	var resultRaw struct {
		ListAlertGroupingSettingsResponse
		AlertGroupingSettings []alertGroupingSettingRaw `json:"alert_grouping_settings"`
	}

	if err = c.decodeJSON(resp, &resultRaw); err != nil {
		return nil, err
	}

	settings := make([]AlertGroupingSetting, 0, len(resultRaw.AlertGroupingSettings))
	for _, rawItem := range resultRaw.AlertGroupingSettings {
		v, _ := getAlertGroupingSettingFromRaw(rawItem)
		if v != nil {
			settings = append(settings, *v)
		}
	}

	result := &resultRaw.ListAlertGroupingSettingsResponse
	result.AlertGroupingSettings = settings
	return result, nil
}

// GetAlertGroupingSetting get an existing Alert Grouping Setting.
func (c *Client) GetAlertGroupingSetting(ctx context.Context, id string) (*AlertGroupingSetting, error) {
	resp, err := c.get(ctx, "/alert_grouping_settings/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resultRaw struct {
		AlertGroupingSetting alertGroupingSettingRaw `json:"alert_grouping_setting"`
	}
	if err := c.decodeJSON(resp, &resultRaw); err != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", err)
	}

	value := resultRaw.AlertGroupingSetting
	return getAlertGroupingSettingFromRaw(value)
}

// DeleteAlertGroupingSetting deletes an existing Alert Grouping Setting.
func (c *Client) DeleteAlertGroupingSetting(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/alert_grouping_settings/"+id)
	return err
}

// UpdateAlertGroupingSetting updates an Alert Grouping Setting.
func (c *Client) UpdateAlertGroupingSetting(ctx context.Context, a AlertGroupingSetting) (*AlertGroupingSetting, error) {
	d := map[string]AlertGroupingSetting{"alert_grouping_setting": a}

	resp, err := c.put(ctx, "/alert_grouping_settings/"+a.ID, d, nil)
	if err != nil {
		return nil, err
	}

	var resultRaw struct {
		AlertGroupingSetting alertGroupingSettingRaw `json:"alert_grouping_setting"`
	}
	if err := c.decodeJSON(resp, &resultRaw); err != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", err)
	}

	value := resultRaw.AlertGroupingSetting
	return getAlertGroupingSettingFromRaw(value)
}

// alertGroupingSettingRaw is an AlertGroupingSetting that overrides its Config with a json raw
// message in order to parse it later.
type alertGroupingSettingRaw struct {
	AlertGroupingSetting
	Config json.RawMessage `json:"config,omitempty"`
}

// getAlertGroupingSettingFromRaw transform the content of a Alert Grouping
// Setting "config" field from a json raw message into the data structure
// corresponding to its "type".
func getAlertGroupingSettingFromRaw(raw alertGroupingSettingRaw) (*AlertGroupingSetting, error) {
	result := &raw.AlertGroupingSetting

	switch raw.Type {
	case AlertGroupingSettingContentBasedType, AlertGroupingSettingContentBasedIntelligentType:
		var cfg AlertGroupingSettingConfigContentBased
		if err := json.Unmarshal(raw.Config, &cfg); err != nil {
			return nil, err
		}
		result.Config = cfg
	case AlertGroupingSettingIntelligentType:
		var cfg AlertGroupingSettingConfigIntelligent
		if err := json.Unmarshal(raw.Config, &cfg); err != nil {
			return nil, err
		}
		result.Config = cfg
	case AlertGroupingSettingTimeType:
		var cfg AlertGroupingSettingConfigTime
		if err := json.Unmarshal(raw.Config, &cfg); err != nil {
			return nil, err
		}
		result.Config = cfg
	}

	return result, nil
}

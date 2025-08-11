package pagerduty

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Enablement represents a feature enablement for a PagerDuty entity
type Enablement struct {
	Feature string `json:"feature,omitempty"`
	Enabled bool   `json:"enabled"`
}

// EnablementRequest represents a request to update an enablement
type EnablementRequest struct {
	Enablement Enablement `json:"enablement"`
}

// EnablementResponse represents the response from an enablement API call
type EnablementResponse struct {
	Enablement Enablement `json:"enablement"`
}

// ListEnablementsResponse represents the response from listing enablements
type ListEnablementsResponse struct {
	Enablements []Enablement `json:"enablements,omitempty"`
	Warnings    []string     `json:"warnings,omitempty"`
}

// EnablementWithWarnings represents an enablement with associated API warnings
type EnablementWithWarnings struct {
	Enablement *Enablement
	Warnings   []string
}

// EnablementsWithWarnings represents multiple enablements with associated API warnings
type EnablementsWithWarnings struct {
	Enablements []Enablement
	Warnings    []string
}

// validateEntityID validates that the entity ID is not empty
func validateEntityID(entityID string) error {
	if strings.TrimSpace(entityID) == "" {
		return fmt.Errorf("entity ID cannot be empty")
	}
	return nil
}

// validateFeatureName validates the feature name
func validateFeatureName(feature string) error {
	if strings.TrimSpace(feature) == "" {
		return fmt.Errorf("feature name cannot be empty")
	}

	// Currently only "aiops" is supported
	validFeatures := []string{"aiops"}
	for _, validFeature := range validFeatures {
		if strings.ToLower(strings.TrimSpace(feature)) == validFeature {
			return nil
		}
	}

	return fmt.Errorf("unsupported feature: %s. Supported features: %s",
		feature, strings.Join(validFeatures, ", "))
}

// validateEntityType validates the entity type
func validateEntityType(entityType string) error {
	if strings.TrimSpace(entityType) == "" {
		return fmt.Errorf("entity type cannot be empty")
	}

	validTypes := []string{"service", "event_orchestration"}
	for _, validType := range validTypes {
		if strings.ToLower(strings.TrimSpace(entityType)) == validType {
			return nil
		}
	}

	return fmt.Errorf("unsupported entity type: %s. Supported types: %s",
		entityType, strings.Join(validTypes, ", "))
}

// validateEnablementRequest validates the enablement request structure
func validateEnablementRequest(req *EnablementRequest) error {
	if req == nil {
		return fmt.Errorf("enablement request cannot be nil")
	}

	return validateFeatureName(req.Enablement.Feature)
}


// getEnablementPath constructs the API path for enablements based on entity type
func getEnablementPath(entityType, entityID string) (string, error) {
	if err := validateEntityType(entityType); err != nil {
		return "", err
	}
	if err := validateEntityID(entityID); err != nil {
		return "", err
	}

	switch strings.ToLower(strings.TrimSpace(entityType)) {
	case "service":
		return fmt.Sprintf("/services/%s/enablements", entityID), nil
	case "event_orchestration":
		return fmt.Sprintf("/event_orchestrations/%s/enablements", entityID), nil
	default:
		return "", fmt.Errorf("unsupported entity type: %s", entityType)
	}
}

// getEnablementFeaturePath constructs the API path for a specific enablement feature
func getEnablementFeaturePath(entityType, entityID, feature string) (string, error) {
	basePath, err := getEnablementPath(entityType, entityID)
	if err != nil {
		return "", err
	}

	if err := validateFeatureName(feature); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", basePath, strings.ToLower(strings.TrimSpace(feature))), nil
}

// handleEnablementError provides enhanced error handling for enablement-specific scenarios
func handleEnablementError(err error, operation, entityType, entityID string) error {
	if err == nil {
		return nil
	}

	// Check if it's an APIError for more specific handling
	if apiErr, ok := err.(APIError); ok {
		switch apiErr.StatusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("bad request when %s enablement for %s %s: %w",
				operation, entityType, entityID, err)
		case http.StatusForbidden:
			return fmt.Errorf("access forbidden when %s enablement for %s %s (check permissions and account entitlements): %w",
				operation, entityType, entityID, err)
		case http.StatusNotFound:
			return fmt.Errorf("%s %s not found when %s enablement: %w",
				entityType, entityID, operation, err)
		case http.StatusTooManyRequests:
			return fmt.Errorf("rate limited when %s enablement for %s %s (retry after rate limit reset): %w",
				operation, entityType, entityID, err)
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return fmt.Errorf("server error when %s enablement for %s %s (temporary issue, retry recommended): %w",
				operation, entityType, entityID, err)
		default:
			return fmt.Errorf("API error when %s enablement for %s %s: %w",
				operation, entityType, entityID, err)
		}
	}

	// Handle network and other errors
	if strings.Contains(err.Error(), "connection") || strings.Contains(err.Error(), "timeout") {
		return fmt.Errorf("network error when %s enablement for %s %s (check connectivity): %w",
			operation, entityType, entityID, err)
	}

	if strings.Contains(err.Error(), "JSON") || strings.Contains(err.Error(), "unmarshal") {
		return fmt.Errorf("JSON parsing error when %s enablement for %s %s (malformed API response): %w",
			operation, entityType, entityID, err)
	}

	// Generic error handling
	return fmt.Errorf("error when %s enablement for %s %s: %w",
		operation, entityType, entityID, err)
}


// getEnablementsFromResponseWithWarnings processes the response from list enablements API and returns warnings
func getEnablementsFromResponseWithWarnings(c *Client, resp *http.Response, err error, entityType, entityID string) ([]Enablement, []string, error) {
	if err != nil {
		return nil, nil, handleEnablementError(err, "listing", entityType, entityID)
	}

	var target ListEnablementsResponse
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, nil, handleEnablementError(
			fmt.Errorf("could not decode JSON response: %w", dErr),
			"listing", entityType, entityID)
	}

	// Return warnings to caller instead of logging them
	return target.Enablements, target.Warnings, nil
}


// getEnablementFromResponseWithWarnings processes the response from update enablement API and returns warnings
func getEnablementFromResponseWithWarnings(c *Client, resp *http.Response, err error, entityType, entityID, feature string) (*Enablement, []string, error) {
	if err != nil {
		return nil, nil, handleEnablementError(err, "updating", entityType, entityID)
	}

	// Different API endpoints use different response formats:
	// - Service endpoints return {"enablements": [...]}
	// - Event orchestration endpoints return {"enablement": {...}}

	var dErr error
	// Based on entity type, try the appropriate parsing method first
	if entityType == "event_orchestration" {
		// Try parsing as single enablement response (event orchestrations)
		var singleTarget EnablementResponse
		dErr = c.decodeJSON(resp, &singleTarget)
		if dErr == nil && singleTarget.Enablement.Feature == feature {
			// Event orchestration responses don't include warnings in the single format
			return &singleTarget.Enablement, []string{}, nil
		}
	}

	if entityType == "service" {
		// Try parsing as array response (services)
		var listTarget ListEnablementsResponse
		dErr = c.decodeJSON(resp, &listTarget)
		if dErr == nil {
			// Find the matching enablement in the list
			for _, enablement := range listTarget.Enablements {
				if enablement.Feature == feature {
					// Return warnings to caller instead of logging them
					return &enablement, listTarget.Warnings, nil
				}
			}
		}
	}

	if dErr != nil {
		return nil, nil, handleEnablementError(
			fmt.Errorf("could not decode JSON response: %w", dErr),
			"updating", entityType, entityID)
	}

	return nil, nil, handleEnablementError(
		fmt.Errorf("enablement %s not found in API response", feature),
		"updating", entityType, entityID)
}

// ListServiceEnablementsWithContext lists all enablements for a service.
func (c *Client) ListServiceEnablementsWithContext(ctx context.Context, serviceID string) (*EnablementsWithWarnings, error) {
	if err := validateEntityID(serviceID); err != nil {
		return nil, handleEnablementError(err, "listing", "service", serviceID)
	}

	path, err := getEnablementPath("service", serviceID)
	if err != nil {
		return nil, handleEnablementError(err, "listing", "service", serviceID)
	}

	resp, err := c.get(ctx, path, nil)
	enablements, warnings, err := getEnablementsFromResponseWithWarnings(c, resp, err, "service", serviceID)
	if err != nil {
		return nil, err
	}
	
	return &EnablementsWithWarnings{
		Enablements: enablements,
		Warnings:    warnings,
	}, nil
}

// ListEventOrchestrationEnablementsWithContext lists all enablements for an event orchestration.
func (c *Client) ListEventOrchestrationEnablementsWithContext(ctx context.Context, orchestrationID string) (*EnablementsWithWarnings, error) {
	if err := validateEntityID(orchestrationID); err != nil {
		return nil, handleEnablementError(err, "listing", "event_orchestration", orchestrationID)
	}

	path, err := getEnablementPath("event_orchestration", orchestrationID)
	if err != nil {
		return nil, handleEnablementError(err, "listing", "event_orchestration", orchestrationID)
	}

	resp, err := c.get(ctx, path, nil)
	enablements, warnings, err := getEnablementsFromResponseWithWarnings(c, resp, err, "event_orchestration", orchestrationID)
	if err != nil {
		return nil, err
	}
	
	return &EnablementsWithWarnings{
		Enablements: enablements,
		Warnings:    warnings,
	}, nil
}

// UpdateServiceEnablementWithContext updates a specific enablement for a service.
func (c *Client) UpdateServiceEnablementWithContext(ctx context.Context, serviceID, feature string, enabled bool) (*EnablementWithWarnings, error) {
	if err := validateEntityID(serviceID); err != nil {
		return nil, handleEnablementError(err, "updating", "service", serviceID)
	}

	if err := validateFeatureName(feature); err != nil {
		return nil, handleEnablementError(err, "updating", "service", serviceID)
	}

	path, err := getEnablementFeaturePath("service", serviceID, feature)
	if err != nil {
		return nil, handleEnablementError(err, "updating", "service", serviceID)
	}

	req := EnablementRequest{
		Enablement: Enablement{
			Feature: feature,
			Enabled: enabled,
		},
	}

	if err := validateEnablementRequest(&req); err != nil {
		return nil, handleEnablementError(err, "updating", "service", serviceID)
	}

	resp, err := c.put(ctx, path, req, nil)
	enablement, warnings, err := getEnablementFromResponseWithWarnings(c, resp, err, "service", serviceID, feature)
	if err != nil {
		return nil, err
	}
	
	return &EnablementWithWarnings{
		Enablement: enablement,
		Warnings:   warnings,
	}, nil
}

// UpdateEventOrchestrationEnablementWithContext updates a specific enablement for an event orchestration.
func (c *Client) UpdateEventOrchestrationEnablementWithContext(ctx context.Context, orchestrationID, feature string, enabled bool) (*EnablementWithWarnings, error) {
	if err := validateEntityID(orchestrationID); err != nil {
		return nil, handleEnablementError(err, "updating", "event_orchestration", orchestrationID)
	}

	if err := validateFeatureName(feature); err != nil {
		return nil, handleEnablementError(err, "updating", "event_orchestration", orchestrationID)
	}

	path, err := getEnablementFeaturePath("event_orchestration", orchestrationID, feature)
	if err != nil {
		return nil, handleEnablementError(err, "updating", "event_orchestration", orchestrationID)
	}

	req := EnablementRequest{
		Enablement: Enablement{
			Feature: feature,
			Enabled: enabled,
		},
	}

	if err := validateEnablementRequest(&req); err != nil {
		return nil, handleEnablementError(err, "updating", "event_orchestration", orchestrationID)
	}

	resp, err := c.put(ctx, path, req, nil)
	enablement, warnings, err := getEnablementFromResponseWithWarnings(c, resp, err, "event_orchestration", orchestrationID, feature)
	if err != nil {
		return nil, err
	}
	
	return &EnablementWithWarnings{
		Enablement: enablement,
		Warnings:   warnings,
	}, nil
}

package pagerduty

import (
	"context"
	"net/http"
)

// ServiceDependency represents a relationship between a business and technical service
type ServiceDependency struct {
	ID                string      `json:"id,omitempty"`
	Type              string      `json:"type,omitempty"`
	SupportingService *ServiceObj `json:"supporting_service,omitempty"`
	DependentService  *ServiceObj `json:"dependent_service,omitempty"`
}

// ServiceObj represents a service object in service relationship
type ServiceObj struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// ListServiceDependencies represents a list of dependencies for a service
type ListServiceDependencies struct {
	Relationships []*ServiceDependency `json:"relationships,omitempty"`
}

// ListBusinessServiceDependencies lists dependencies of a business service.
// It's recommended to use ListBusinessServiceDependenciesWithContext instead.
func (c *Client) ListBusinessServiceDependencies(businessServiceID string) (*ListServiceDependencies, *http.Response, error) {
	return c.listBusinessServiceDependenciesWithContext(context.Background(), businessServiceID)
}

// ListBusinessServiceDependenciesWithContext lists dependencies of a business service.
func (c *Client) ListBusinessServiceDependenciesWithContext(ctx context.Context, businessServiceID string) (*ListServiceDependencies, error) {
	lsd, _, err := c.listBusinessServiceDependenciesWithContext(ctx, businessServiceID)
	return lsd, err
}

func (c *Client) listBusinessServiceDependenciesWithContext(ctx context.Context, businessServiceID string) (*ListServiceDependencies, *http.Response, error) {
	resp, err := c.get(ctx, "/service_dependencies/business_services/"+businessServiceID)
	if err != nil {
		return nil, nil, err
	}

	var result ListServiceDependencies
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListTechnicalServiceDependencies lists dependencies of a technical service.
// It's recommended to use ListTechnicalServiceDependenciesWithContext instead.
func (c *Client) ListTechnicalServiceDependencies(serviceID string) (*ListServiceDependencies, *http.Response, error) {
	return c.listTechnicalServiceDependenciesWithContext(context.Background(), serviceID)
}

// ListTechnicalServiceDependenciesWithContext lists dependencies of a technical service.
func (c *Client) ListTechnicalServiceDependenciesWithContext(ctx context.Context, serviceID string) (*ListServiceDependencies, error) {
	lsd, _, err := c.listTechnicalServiceDependenciesWithContext(ctx, serviceID)
	return lsd, err
}

func (c *Client) listTechnicalServiceDependenciesWithContext(ctx context.Context, serviceID string) (*ListServiceDependencies, *http.Response, error) {
	resp, err := c.get(ctx, "/service_dependencies/technical_services/"+serviceID)
	if err != nil {
		return nil, nil, err
	}

	var result ListServiceDependencies
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// AssociateServiceDependencies Create new dependencies between two services.
// It's recommended to use AssociateServiceDependenciesWithContext instead.
func (c *Client) AssociateServiceDependencies(dependencies *ListServiceDependencies) (*ListServiceDependencies, *http.Response, error) {
	return c.associateServiceDependenciesWithContext(context.Background(), dependencies)
}

// AssociateServiceDependenciesWithContext Create new dependencies between two services.
func (c *Client) AssociateServiceDependenciesWithContext(ctx context.Context, dependencies *ListServiceDependencies) (*ListServiceDependencies, error) {
	lsd, _, err := c.associateServiceDependenciesWithContext(ctx, dependencies)
	return lsd, err
}

func (c *Client) associateServiceDependenciesWithContext(ctx context.Context, dependencies *ListServiceDependencies) (*ListServiceDependencies, *http.Response, error) {
	resp, err := c.post(ctx, "/service_dependencies/associate", dependencies, nil)
	if err != nil {
		return nil, nil, err
	}

	var result ListServiceDependencies
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// DisassociateServiceDependencies Disassociate dependencies between two services.
func (c *Client) DisassociateServiceDependencies(dependencies *ListServiceDependencies) (*ListServiceDependencies, *http.Response, error) {
	return c.disassociateServiceDependenciesWithContext(context.Background(), dependencies)
}

// DisassociateServiceDependenciesWithContext Disassociate dependencies between two services.
func (c *Client) DisassociateServiceDependenciesWithContext(ctx context.Context, dependencies *ListServiceDependencies) (*ListServiceDependencies, error) {
	lsd, _, err := c.disassociateServiceDependenciesWithContext(ctx, dependencies)
	return lsd, err
}

// DisassociateServiceDependencies Disassociate dependencies between two services.
func (c *Client) disassociateServiceDependenciesWithContext(ctx context.Context, dependencies *ListServiceDependencies) (*ListServiceDependencies, *http.Response, error) {
	resp, err := c.post(ctx, "/service_dependencies/disassociate", dependencies, nil)
	if err != nil {
		return nil, nil, err
	}

	var result ListServiceDependencies
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

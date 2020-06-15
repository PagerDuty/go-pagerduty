package pagerduty

import (
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
func (c *Client) ListBusinessServiceDependencies(businessServiceID string) (*ListServiceDependencies, *http.Response, error) {
	resp, err := c.get("/service_dependencies/business_services/" + businessServiceID)
	if err != nil {
		return nil, nil, err
	}
	var result ListServiceDependencies
	return &result, resp, c.decodeJSON(resp, &result)
}

// ListTechnicalServiceDependencies lists dependencies of a technical service.
func (c *Client) ListTechnicalServiceDependencies(serviceID string) (*ListServiceDependencies, *http.Response, error) {
	resp, err := c.get("/service_dependencies/technical_services/" + serviceID)
	if err != nil {
		return nil, nil, err
	}
	var result ListServiceDependencies
	return &result, resp, c.decodeJSON(resp, &result)
}

// AssociateServiceDependencies Create new dependencies between two services.
func (c *Client) AssociateServiceDependencies(dependencies *ListServiceDependencies) (*ListServiceDependencies, *http.Response, error) {
	data := make(map[string]*ListServiceDependencies)
	data["relationships"] = dependencies
	resp, err := c.post("/service_dependencies/associate", data, nil)
	if err != nil {
		return nil, nil, err
	}
	var result ListServiceDependencies
	return &result, resp, c.decodeJSON(resp, &result)
}

// DisassociateServiceDependencies Disassociate dependencies between two services.
func (c *Client) DisassociateServiceDependencies(dependencies *ListServiceDependencies) (*ListServiceDependencies, *http.Response, error) {
	data := make(map[string]*ListServiceDependencies)
	data["relationships"] = dependencies
	resp, err := c.post("/service_dependencies/disassociate", data, nil)
	if err != nil {
		return nil, nil, err
	}
	var result ListServiceDependencies
	return &result, resp, c.decodeJSON(resp, &result)
}

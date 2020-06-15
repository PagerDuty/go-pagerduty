package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// BusinessService represents a business service.
type BusinessService struct {
	ID             string               `json:"id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Type           string               `json:"type,omitempty"`
	Summary        string               `json:"summary,omitempty"`
	Self           string               `json:"self,omitempty"`
	PointOfContact string               `json:"point_of_contact,omitempty"`
	HTMLUrl        string               `json:"html_url,omitempty"`
	Description    string               `json:"description,omitempty"`
	Team           *BusinessServiceTeam `json:"team,omitempty"`
}

// BusinessServiceTeam represents a team object in a business service
type BusinessServiceTeam struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Self string `json:"self,omitempty"`
}

// BusinessServicePayload represents payload with a business service object
type BusinessServicePayload struct {
	BusinessService *BusinessService `json:"business_service,omitempty"`
}

// ListBusinessServicesResponse represents a list response of business services.
type ListBusinessServicesResponse struct {
	Total            uint               `json:"total,omitempty"`
	BusinessServices []*BusinessService `json:"business_services,omitempty"`
	Offset           uint               `json:"offset,omitempty"`
	More             bool               `json:"more,omitempty"`
	Limit            uint               `json:"limit,omitempty"`
}

// ListBusinessServiceOptions is the data structure used when calling the ListBusinessServices API endpoint.
type ListBusinessServiceOptions struct {
	APIListObject
}

// ListBusinessServices lists existing business services.
func (c *Client) ListBusinessServices(o ListBusinessServiceOptions) (*ListBusinessServicesResponse, error) {
	queryParms, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	businessServiceResponse := new(ListBusinessServicesResponse)
	businessServices := make([]*BusinessService, 0)

	// Create a handler closure capable of parsing data from the business_services endpoint
	// and appending resultant business_services to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListBusinessServicesResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		businessServices = append(businessServices, result.BusinessServices...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet("/business_services"+queryParms.Encode(), responseHandler); err != nil {
		return nil, err
	}
	businessServiceResponse.BusinessServices = businessServices

	return businessServiceResponse, nil
}

// CreateBusinessService creates a new business service.
func (c *Client) CreateBusinessService(b *BusinessService) (*BusinessService, *http.Response, error) {
	data := make(map[string]*BusinessService)
	data["business_service"] = b
	resp, err := c.post("/business_services", data, nil)
	return getBusinessServiceFromResponse(c, resp, err)
}

// GetBusinessService gets details about a business service.
func (c *Client) GetBusinessService(ID string) (*BusinessService, *http.Response, error) {
	resp, err := c.get("/business_services/" + ID)
	return getBusinessServiceFromResponse(c, resp, err)
}

// DeleteBusinessService deletes a business_service.
func (c *Client) DeleteBusinessService(ID string) error {
	_, err := c.delete("/business_services/" + ID)
	return err
}

// UpdateBusinessService updates a business_service.
func (c *Client) UpdateBusinessService(b *BusinessService) (*BusinessService, *http.Response, error) {
	v := make(map[string]*BusinessService)
	v["business_service"] = b
	resp, err := c.put("/business_services/"+b.ID, v, nil)
	return getBusinessServiceFromResponse(c, resp, err)
}

func getBusinessServiceFromResponse(c *Client, resp *http.Response, err error) (*BusinessService, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}
	var target map[string]BusinessService
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	t, nodeOK := target["business_service"]
	if !nodeOK {
		return nil, nil, fmt.Errorf("JSON response does not have business_service field")
	}
	return &t, resp, nil
}

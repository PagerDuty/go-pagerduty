package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Vendor represents a specific type of integration. AWS Cloudwatch, Splunk, Datadog, etc are all examples of vendors that can be integrated in PagerDuty by making an integration.
type Vendor struct {
	APIObject
	Name                string `json:"name,omitempty"`
	LogoURL             string `json:"logo_url,omitempty"`
	LongName            string `json:"long_name,omitempty"`
	WebsiteURL          string `json:"website_url,omitempty"`
	Description         string `json:"description,omitempty"`
	Connectable         bool   `json:"connectable,omitempty"`
	ThumbnailURL        string `json:"thumbnail_url,omitempty"`
	GenericServiceType  string `json:"generic_service_type,omitempty"`
	IntegrationGuideURL string `json:"integration_guide_url,omitempty"`
}

// ListVendorResponse is the data structure returned from calling the ListVendors API endpoint.
type ListVendorResponse struct {
	APIListObject
	Vendors []Vendor
}

// ListVendorOptions is the data structure used when calling the ListVendors API endpoint.
type ListVendorOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

// ListVendors lists existing vendors.
func (pd *PagerdutyClient) ListVendors(o ListVendorOptions) (*ListVendorResponse, error) {
	v, err := query.Values(o)

	if err != nil {
		return nil, err
	}

	resp, err := pd.Get("/vendors?" + v.Encode())

	if err != nil {
		return nil, err
	}

	var result ListVendorResponse
	return &result, DecodeJSON(resp, &result)
}

// GetVendor gets details about an existing vendor.
func (pd *PagerdutyClient) GetVendor(id string) (*Vendor, error) {
	resp, err := pd.Get("/vendors/" + id)
	return getVendorFromResponse(pd, resp, err)
}

func getVendorFromResponse(pd *PagerdutyClient, resp *http.Response, err error) (*Vendor, error) {
	if err != nil {
		return nil, err
	}
	var target map[string]Vendor
	if dErr := DecodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "vendor"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}

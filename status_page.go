package pagerduty

import (
	"context"
)

type StatusPage struct {
	ID             string
	Name           string
	PublishedAt    string
	StatusPageType string
	URL            string
	Type           string
}

type StatusPageImpact struct {
	ID          string
	Self        string
	Description string
	PostType    string
	StatusPage  StatusPage
	Type        string
}

// ListStatusPagesResponse is the data structure returned from calling the ListStatusPages API endpoint.
type ListStatusPagesResponse struct {
	APIListObject
	StatusPages []StatusPage `json:"status_pages"`
}

// ListStatusPageImpactsResponse is the data structure returned from calling the ListStatusPagesImpacts API endpoint.
type ListStatusPageImpactsResponse struct {
	APIListObject
	StatusPageImpacts []StatusPageImpact `json:"impacts"`
}

// ListStatusPages lists the given types of status pages
func (c *Client) ListStatusPages(statusPageType string) (*ListStatusPagesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages?status_page_type="+statusPageType, h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPagesResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListStatusPageImpacts lists the given types of impacts for the specific status page
func (c *Client) ListStatusPageImpacts(id string, postType string) (*ListStatusPageImpactsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/impacts?post_type="+postType, h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPageImpactsResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStatusPageImpact gets the specified status page impact
func (c *Client) GetStatusPageImpact(statusPageID string, impactID string) (*StatusPageImpact, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/impacts/"+impactID, h)
	if err != nil {
		return nil, err
	}

	var result StatusPageImpact
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

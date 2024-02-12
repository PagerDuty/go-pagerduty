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

// ListStatusPagesResponse is the data structure returned from calling the ListStatusPages API endpoint.
type ListStatusPagesResponse struct {
	APIListObject
	StatusPages []StatusPage `json:"status_pages"`
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

package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
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

type StatusPageService struct {
	ID              string
	Self            string
	Name            string
	StatusPage      StatusPage
	BusinessService Service
	Type            string
}

type StatusPageSeverity struct {
	ID          string
	Self        string
	Description string
	PostType    string
	StatusPage  StatusPage
	Type        string
}

type StatusPageStatus struct {
	ID          string
	Self        string
	Description string
	PostType    string
	StatusPage  StatusPage
	Type        string
}

type StatusPagePost struct {
	ID             string
	Self           string
	Type           string
	PostType       string
	StatusPage     StatusPage
	LinkedResource LinkedResource
	PostMortem     PostMortem
	Title          string
	StartsAt       string
	EndsAt         string
	Updates        []StatusPageUpdate
}

type LinkedResource struct {
	ID   string
	Self string
}

type PostMortem struct {
	ID   string
	Self string
}

type StatusPageUpdate struct {
	ID                string
	Self              string
	Message           string
	ReviewedStatus    string
	Status            StatusPageStatus
	Severity          StatusPageSeverity
	ImpactedServices  []StatusPagePostUpdateImpact
	UpdateFrequencyMS uint
	NotifySubscribers bool
	ReportedAt        string
	Type              string
}

type StatusPagePostUpdateImpact struct {
	Service  Service
	Severity StatusPageSeverity
}

type ListStatusPagePostOptions struct {
	PostType       string             `url:"post_type,omitempty"`
	ReviewedStatus string             `url:"reviewed_status,omitempty"`
	Statuses       []StatusPageStatus `url:"statuses,omitempty"`
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

// ListStatusPageServicesResponse is the data structure returned from calling the ListStatusPagesServices API endpoint.
type ListStatusPageServicesResponse struct {
	APIListObject
	StatusPageServices []StatusPageService `json:"services"`
}

// ListStatusPageSeveritiesResponse is the data structure returned from calling the ListStatusPageSeverities API endpoint.
type ListStatusPageSeveritiesResponse struct {
	APIListObject
	StatusPageSeverities []StatusPageSeverity `json:"severities"`
}

// ListStatusPageStatusesResponse is the data structure returned from calling the ListStatusPageStatuses API endpoint.
type ListStatusPageStatusesResponse struct {
	APIListObject
	StatusPageStatuses []StatusPageStatus `json:"statuses"`
}

// ListStatusPagePostsResponse is the data structure returned from calling the ListStatusPagePosts API endpoint.
type ListStatusPagePostsResponse struct {
	APIListObject
	StatusPagePosts []StatusPagePost `json:"posts"`
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

// ListStatusPageImpacts lists the given types of impacts for the specified status page
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

// ListStatusPageServices lists the services for the specified status page
func (c *Client) ListStatusPageServices(id string) (*ListStatusPageServicesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/services", h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPageServicesResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStatusPageService gets the specified status page service
func (c *Client) GetStatusPageService(statusPageID string, serviceID string) (*StatusPageService, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/services/"+serviceID, h)
	if err != nil {
		return nil, err
	}

	var result StatusPageService
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListStatusPageSeverities lists the severities for the specified status page
func (c *Client) ListStatusPageSeverities(id string, postType string) (*ListStatusPageSeveritiesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/severities?post_type="+postType, h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPageSeveritiesResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStatusPageSeverity gets the specified status page severity
func (c *Client) GetStatusPageSeverity(statusPageID string, severityID string) (*StatusPageSeverity, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/severities/"+severityID, h)
	if err != nil {
		return nil, err
	}

	var result StatusPageSeverity
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListStatusPageStatuses lists the statuses for the specified status page
func (c *Client) ListStatusPageStatuses(id string, postType string) (*ListStatusPageStatusesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/statuses?post_type="+postType, h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPageStatusesResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStatusPageStatus gets the specified status page status
func (c *Client) GetStatusPageStatus(statusPageID string, statusID string) (*StatusPageStatus, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/statuses/"+statusID, h)
	if err != nil {
		return nil, err
	}

	var result StatusPageStatus
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListStatusPagePosts lists the posts for the specified status page
func (c *Client) ListStatusPagePosts(id string, o ListStatusPagePostOptions) (*ListStatusPagePostsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/posts?"+v.Encode(), h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPagePostsResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateStatusPagePost create a Post for a Status Page by Status Page ID
func (c *Client) CreateStatusPagePost(id string, p StatusPagePost) (*StatusPagePost, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	d := map[string]StatusPagePost{
		"post": p,
	}
	resp, err := c.post(context.Background(), "/status_pages/"+id+"/posts", d, h)
	return getStatusPagePostFromResponse(c, resp, err)
}

func getStatusPagePostFromResponse(c *Client, resp *http.Response, err error) (*StatusPagePost, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]StatusPagePost
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "post"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

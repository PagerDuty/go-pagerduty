package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type StatusPage struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	PublishedAt    string `json:"published_at,omitempty"`
	StatusPageType string `json:"status_page_type,omitempty"`
	URL            string `json:"url,omitempty"`
	Type           string `json:"type,omitempty"`
}

type StatusPageImpact struct {
	ID          string     `json:"id,omitempty"`
	Self        string     `json:"self,omitempty"`
	Description string     `json:"description,omitempty"`
	PostType    string     `json:"post_type,omitempty"`
	StatusPage  StatusPage `json:"status_page,omitempty"`
	Type        string     `json:"type,omitempty"`
}

type StatusPageService struct {
	ID              string     `json:"id,omitempty"`
	Self            string     `json:"self,omitempty"`
	Name            string     `json:"name,omitempty"`
	StatusPage      StatusPage `json:"status_page,omitempty"`
	BusinessService Service    `json:"business_service,omitempty"`
	Type            string     `json:"type,omitempty"`
}

type StatusPageSeverity struct {
	ID          string     `json:"id,omitempty"`
	Self        string     `json:"self,omitempty"`
	Description string     `json:"description,omitempty"`
	PostType    string     `json:"post_type,omitempty"`
	StatusPage  StatusPage `json:"status_page,omitempty"`
	Type        string     `json:"type,omitempty"`
}

type StatusPageStatus struct {
	ID          string     `json:"id,omitempty"`
	Self        string     `json:"self,omitempty"`
	Description string     `json:"description,omitempty"`
	PostType    string     `json:"post_type,omitempty"`
	StatusPage  StatusPage `json:"status_page,omitempty"`
	Type        string     `json:"type,omitempty"`
}

type StatusPagePost struct {
	ID             string                 `json:"id,omitempty"`
	Self           string                 `json:"self,omitempty"`
	Type           string                 `json:"type,omitempty"`
	PostType       string                 `json:"post_type,omitempty"`
	StatusPage     StatusPage             `json:"status_page,omitempty"`
	LinkedResource LinkedResource         `json:"linked_resource,omitempty"`
	Postmortem     Postmortem             `json:"postmortem,omitempty"`
	Title          string                 `json:"title,omitempty"`
	StartsAt       string                 `json:"starts_at,omitempty"`
	EndsAt         string                 `json:"ends_at,omitempty"`
	Updates        []StatusPagePostUpdate `json:"updates,omitempty"`
}

type LinkedResource struct {
	ID   string `json:"id,omitempty"`
	Self string `json:"self,omitempty"`
}

type Postmortem struct {
	ID                string        `json:"id,omitempty"`
	Self              string        `json:"self,omitempty"`
	NotifySubscribers bool          `json:"notify_subscribers,omitempty"`
	ReportedAt        string        `json:"reported_at,omitempty"`
	Type              string        `json:"type,omitempty"`
	Message           string        `json:"message,omitempty"`
	Post              ShortPostType `json:"post,omitempty"`
}

type ShortPostType struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type StatusPagePostUpdate struct {
	ID                string                       `json:"id,omitempty"`
	Self              string                       `json:"self,omitempty"`
	Message           string                       `json:"message,omitempty"`
	ReviewedStatus    string                       `json:"reviewed_status,omitempty"`
	Status            StatusPageStatus             `json:"status,omitempty"`
	Severity          StatusPageSeverity           `json:"severity,omitempty"`
	ImpactedServices  []StatusPagePostUpdateImpact `json:"impacted_services,omitempty"`
	UpdateFrequencyMS uint                         `json:"update_frequency_ms,omitempty"`
	Post              StatusPagePost               `json:"post,omitempty"`
	NotifySubscribers bool                         `json:"notify_subscribers,omitempty"`
	ReportedAt        string                       `json:"reported_at,omitempty"`
	Type              string                       `json:"type,omitempty"`
}

type StatusPagePostUpdateImpact struct {
	Service  Service            `json:"service,omitempty"`
	Severity StatusPageSeverity `json:"severity,omitempty"`
}

type StatusPageSubscription struct {
	Channel            string             `json:"channel,omitempty"`
	Contact            string             `json:"contact,omitempty"`
	ID                 string             `json:"id,omitempty"`
	Self               string             `json:"self,omitempty"`
	Status             string             `json:"status,omitempty"`
	StatusPage         StatusPage         `json:"status_page,omitempty"`
	SubscribableObject SubscribableObject `json:"subscribable_object,omitempty"`
	Type               string             `json:"type,omitempty"`
}

type SubscribableObject struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type ListStatusPageOptions struct {
	StatusPageType string `url:"status_page_type,omitempty"`
}

type ListStatusPageImpactOptions struct {
	PostType string `url:"post_type,omitempty"`
}

type ListStatusPageSeveritiesOptions struct {
	PostType string `url:"post_type,omitempty"`
}

type ListStatusPageStatusesOptions struct {
	PostType string `url:"post_type,omitempty"`
}

type ListStatusPagePostOptions struct {
	PostType       string   `url:"post_type,omitempty"`
	ReviewedStatus string   `url:"reviewed_status,omitempty"`
	Status         []string `url:"status,omitempty"`
}

type ListStatusPagePostUpdateOptions struct {
	ReviewedStatus string `url:"reviewed_status,omitempty"`
}

type ListStatusPageSubscriptionsOptions struct {
	Channel string `url:"channel,omitempty"`
	Status  string `url:"status,omitempty"`
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

// ListStatusPagePostUpdatesResponse is the data structure returned from calling the ListStatusPagePostUpdates API endpoint.
type ListStatusPagePostUpdatesResponse struct {
	APIListObject
	StatusPagePostUpdates []StatusPagePostUpdate `json:"post_updates"`
}

// ListStatusPageSubscriptionsResponse is the data structure returned from calling the ListStatusPageSubscriptions API endpoint.
type ListStatusPageSubscriptionsResponse struct {
	APIListObject
	StatusPageSubscriptions []StatusPageSubscription `json:"subscriptions"`
}

// ListStatusPages lists the given types of status pages
func (c *Client) ListStatusPages(o ListStatusPageOptions) (*ListStatusPagesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(context.Background(), "/status_pages?"+v.Encode(), h)
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
func (c *Client) ListStatusPageImpacts(id string, o ListStatusPageImpactOptions) (*ListStatusPageImpactsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/impacts?"+v.Encode(), h)
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
	return getStatusPageImpactFromResponse(c, resp, err)
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
	return getStatusPageServiceFromResponse(c, resp, err)
}

// ListStatusPageSeverities lists the severities for the specified status page
func (c *Client) ListStatusPageSeverities(id string, o ListStatusPageSeveritiesOptions) (*ListStatusPageSeveritiesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/severities?"+v.Encode(), h)
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
	return getStatusPageSeverityFromResponse(c, resp, err)
}

// ListStatusPageStatuses lists the statuses for the specified status page
func (c *Client) ListStatusPageStatuses(id string, o ListStatusPageStatusesOptions) (*ListStatusPageStatusesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/statuses?"+v.Encode(), h)
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
	/* Note: Does not currently support the include query parameter */
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/statuses/"+statusID, h)
	return getStatusPageStatusFromResponse(c, resp, err)
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

// CreateStatusPagePost creates a Post for a Status Page by Status Page ID
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

// GetStatusPagePost gets the specified status page status
func (c *Client) GetStatusPagePost(statusPageID string, postID string) (*StatusPagePost, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID, h)
	return getStatusPagePostFromResponse(c, resp, err)
}

// UpdateStatusPagePost updates a Post for a Status Page by Status Page ID
func (c *Client) UpdateStatusPagePost(statusPageID string, postID string, p StatusPagePost) (*StatusPagePost, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	d := map[string]StatusPagePost{
		"post": p,
	}
	resp, err := c.put(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID, d, h)
	return getStatusPagePostFromResponse(c, resp, err)
}

// DeleteStatusPagePost deletes a Post for a Status Page by Status Page ID
func (c *Client) DeleteStatusPagePost(statusPageID string, postID string) error {
	/* Note: The API requires sending in the below header, but the client does not support headers for the delete() function, so we have to use do() */
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	_, err := c.do(context.Background(), http.MethodDelete, "/status_pages/"+statusPageID+"/posts/"+postID, nil, h)
	return err
}

// ListStatusPagePostUpdates lists the post updates for the specified status page and post
func (c *Client) ListStatusPagePostUpdates(statusPageID string, postID string, o ListStatusPagePostUpdateOptions) (*ListStatusPagePostUpdatesResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID+"/post_updates?"+v.Encode(), h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPagePostUpdatesResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateStatusPagePostUpdate creates an Post Update for a Status Page by Status Page ID and Post ID
func (c *Client) CreateStatusPagePostUpdate(statusPageID string, postID string, u StatusPagePostUpdate) (*StatusPagePostUpdate, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	d := map[string]StatusPagePostUpdate{
		"post_update": u,
	}
	resp, err := c.post(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID+"/post_updates", d, h)
	return getStatusPagePostUpdateFromResponse(c, resp, err)
}

// GetStatusPagePostUpdate gets the specified status page post update
func (c *Client) GetStatusPagePostUpdate(statusPageID string, postID string, postUpdateID string) (*StatusPagePostUpdate, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID+"/post_updates/"+postUpdateID, h)
	return getStatusPagePostUpdateFromResponse(c, resp, err)
}

// UpdateStatusPagePostUpdate updates a Post Update for a Status Page by Status Page ID and Post ID
func (c *Client) UpdateStatusPagePostUpdate(statusPageID string, postID string, postUpdateID string, u StatusPagePostUpdate) (*StatusPagePostUpdate, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	d := map[string]StatusPagePostUpdate{
		"post_update": u,
	}
	resp, err := c.put(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID+"/post_updates/"+postUpdateID, d, h)
	return getStatusPagePostUpdateFromResponse(c, resp, err)
}

// DeleteStatusPagePostUpdate deletes a Post Update for a Status Page by Status Page ID and Post ID
func (c *Client) DeleteStatusPagePostUpdate(statusPageID string, postID string, postUpdateID string) error {
	/* Note: The API requires sending in the below header, but the client does not support headers for the delete() function, so we have to use do() */
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	_, err := c.do(context.Background(), http.MethodDelete, "/status_pages/"+statusPageID+"/posts/"+postID+"/post_updates/"+postUpdateID, nil, h)
	return err
}

// GetStatusPagePostPostmortem gets the specified status page post post-mortem
func (c *Client) GetStatusPagePostPostmortem(statusPageID string, postID string) (*Postmortem, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID+"/postmortem", h)
	return getStatusPagePostPostmortemFromResponse(c, resp, err)
}

// CreateStatusPagePostPostmortem creates a post-mortem for a Status Page by Status Page ID and Post ID
func (c *Client) CreateStatusPagePostPostmortem(statusPageID string, postID string, p Postmortem) (*Postmortem, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	d := map[string]Postmortem{
		"postmortem": p,
	}
	resp, err := c.post(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID+"/postmortem", d, h)
	return getStatusPagePostPostmortemFromResponse(c, resp, err)
}

// UpdateStatusPagePostPostmortem updates a post-mortem for a Status Page by Status Page ID and Post ID
func (c *Client) UpdateStatusPagePostPostmortem(statusPageID string, postID string, p Postmortem) (*Postmortem, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	d := map[string]Postmortem{
		"postmortem": p,
	}
	resp, err := c.put(context.Background(), "/status_pages/"+statusPageID+"/posts/"+postID+"/postmortem", d, h)
	return getStatusPagePostPostmortemFromResponse(c, resp, err)
}

// DeleteStatusPagePostPostmortem deletes a post-mortem for a Status Page by Status Page ID and Post ID
func (c *Client) DeleteStatusPagePostPostmortem(statusPageID string, postID string) error {
	/* Note: The API requires sending in the below header, but the client does not support headers for the delete() function, so we have to use do() */
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	_, err := c.do(context.Background(), http.MethodDelete, "/status_pages/"+statusPageID+"/posts/"+postID+"/postmortem", nil, h)
	return err
}

// ListStatusPageSubscriptions lists the subscriptions for the specified status page
func (c *Client) ListStatusPageSubscriptions(id string, o ListStatusPageSubscriptionsOptions) (*ListStatusPageSubscriptionsResponse, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(context.Background(), "/status_pages/"+id+"/subscriptions?"+v.Encode(), h)
	if err != nil {
		return nil, err
	}

	var result ListStatusPageSubscriptionsResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateStatusPageSubscription creates a Subscription for a Status Page by Status Page ID
func (c *Client) CreateStatusPageSubscription(statusPageID string, s StatusPageSubscription) (*StatusPageSubscription, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	d := map[string]StatusPageSubscription{
		"subscription": s,
	}
	resp, err := c.post(context.Background(), "/status_pages/"+statusPageID+"/subscriptions", d, h)
	return getStatusPageSubscriptionFromResponse(c, resp, err)
}

// GetStatusPageSubscription gets the Subscription for a Status Page by Status Page ID and Subscription ID.
func (c *Client) GetStatusPageSubscription(statusPageID string, subscriptionID string) (*StatusPageSubscription, error) {
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	resp, err := c.get(context.Background(), "/status_pages/"+statusPageID+"/subscriptions/"+subscriptionID, h)
	return getStatusPageSubscriptionFromResponse(c, resp, err)
}

// DeleteStatusPageSubscription deletes a Subscription for a Status Page by Status Page ID and Subscription ID.
func (c *Client) DeleteStatusPageSubscription(statusPageID string, subscriptionID string) error {
	/* Note: The API requires sending in the below header, but the client does not support headers for the delete() function, so we have to use do() */
	h := map[string]string{
		"X-EARLY-ACCESS": "status-pages-early-access",
	}
	_, err := c.do(context.Background(), http.MethodDelete, "/status_pages/"+statusPageID+"/subscriptions/"+subscriptionID, nil, h)
	return err
}

func getStatusPageImpactFromResponse(c *Client, resp *http.Response, err error) (*StatusPageImpact, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]StatusPageImpact
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "impact"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

func getStatusPageServiceFromResponse(c *Client, resp *http.Response, err error) (*StatusPageService, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]StatusPageService
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "service"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

func getStatusPageStatusFromResponse(c *Client, resp *http.Response, err error) (*StatusPageStatus, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]StatusPageStatus
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "status"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

func getStatusPageSeverityFromResponse(c *Client, resp *http.Response, err error) (*StatusPageSeverity, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]StatusPageSeverity
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "severity"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
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

func getStatusPagePostUpdateFromResponse(c *Client, resp *http.Response, err error) (*StatusPagePostUpdate, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]StatusPagePostUpdate
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "post_update"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

func getStatusPagePostPostmortemFromResponse(c *Client, resp *http.Response, err error) (*Postmortem, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]Postmortem
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "postmortem"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

func getStatusPageSubscriptionFromResponse(c *Client, resp *http.Response, err error) (*StatusPageSubscription, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]StatusPageSubscription
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	const rootNode = "subscription"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

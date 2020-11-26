package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Tag is a way to label user, team and escalation policies in PagerDuty
type Tag struct {
	APIObject
	Label string `json:"label,omitempty"`
}

// ListTagResponse is the structure used when calling the ListTags API endpoint.
type ListTagResponse struct {
	APIListObject
	Tags []*Tag `json:"tags"`
}

// ListUserResponse is the structure used to list users assigned a given tag
type ListUserResponse struct {
	APIListObject
	Users []*APIObject `json:"users,omitempty"`
}

// ListTeamsForTagResponse is the structure used to list teams assigned a given tag
type ListTeamsForTagResponse struct {
	APIListObject
	Teams []*APIObject `json:"teams,omitempty"`
}

// ListEPResponse is the structure used to list escalation policies assigned a given tag
type ListEPResponse struct {
	APIListObject
	EscalationPolicies []*APIObject `json:"escalation_policies,omitempty"`
}

// ListTagOptions are the input parameters used when calling the ListTags API endpoint.
type ListTagOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

// TagAssignments can be applied teams, users and escalation policies
type TagAssignments struct {
	Add    []*TagAssignment `json:"add,omitempty"`
	Remove []*TagAssignment `json:"remove,omitempty"`
}

// TagAssignment is the structure for assigning tags to an entity
type TagAssignment struct {
	Type  string `json:"type"`
	TagID string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

// ListTags lists tags of your PagerDuty account, optionally filtered by a search query.
func (c *Client) ListTags(o ListTagOptions) (*ListTagResponse, error) {
	return getTagList(c, "", "", o)
}

// CreateTag creates a new tag.
func (c *Client) CreateTag(t *Tag) (*Tag, *http.Response, error) {
	data := make(map[string]*Tag)
	data["tag"] = t
	resp, err := c.post("/tags", data, nil)
	return getTagFromResponse(c, resp, err)
}

// DeleteTag removes an existing tag.
func (c *Client) DeleteTag(id string) error {
	_, err := c.delete("/tags/" + id)
	return err
}

// GetTag gets details about an existing tag.
func (c *Client) GetTag(id string) (*Tag, *http.Response, error) {
	resp, err := c.get("/tags/" + id)
	return getTagFromResponse(c, resp, err)
}

// AssignTags adds and removes tag assignments with entities
func (c *Client) AssignTags(e, eid string, a *TagAssignments) (*http.Response, error) {
	if resp, err := c.post("/"+e+"/"+eid+"/change_tags", a, nil); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetUsersByTag get related Users for the Tag.
func (c *Client) GetUsersByTag(tid string) (*ListUserResponse, error) {
	userResponse := new(ListUserResponse)
	users := make([]*APIObject, 0)

	// Create a handler closure capable of parsing data from the business_services endpoint
	// and appending resultant business_services to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListUserResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		users = append(users, result.Users...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet("/tags/"+tid+"/users/", responseHandler); err != nil {
		return nil, err
	}
	userResponse.Users = users
	fmt.Println()
	return userResponse, nil
}

// GetTeamsByTag get related Users for the Tag.
func (c *Client) GetTeamsByTag(tid string) (*ListTeamsForTagResponse, error) {
	teamsResponse := new(ListTeamsForTagResponse)
	teams := make([]*APIObject, 0)

	// Create a handler closure capable of parsing data from the business_services endpoint
	// and appending resultant business_services to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListTeamsForTagResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		teams = append(teams, result.Teams...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet("/tags/"+tid+"/teams/", responseHandler); err != nil {
		return nil, err
	}
	teamsResponse.Teams = teams

	return teamsResponse, nil
}

// GetEscalationPoliciesByTag get related Users for the Tag.
func (c *Client) GetEscalationPoliciesByTag(tid string) (*ListEPResponse, error) {
	epResponse := new(ListEPResponse)
	eps := make([]*APIObject, 0)

	// Create a handler closure capable of parsing data from the business_services endpoint
	// and appending resultant business_services to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListEPResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		eps = append(eps, result.EscalationPolicies...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet("/tags/"+tid+"/escalation_policies/", responseHandler); err != nil {
		return nil, err
	}
	epResponse.EscalationPolicies = eps

	return epResponse, nil
}

// GetTagsForEntity Get related tags for Users, Teams or Escalation Policies.
func (c *Client) GetTagsForEntity(e, eid string, o ListTagOptions) (*ListTagResponse, error) {
	return getTagList(c, e, eid, o)
}

func getTagFromResponse(c *Client, resp *http.Response, err error) (*Tag, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}
	var target map[string]Tag
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "tag"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, resp, nil
}

// getTagList  is a utility function that processes all pages of a ListTagResponse
func getTagList(c *Client, e, eid string, o ListTagOptions) (*ListTagResponse, error) {
	queryParms, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	tagResponse := new(ListTagResponse)
	tags := make([]*Tag, 0)

	// Create a handler closure capable of parsing data from the business_services endpoint
	// and appending resultant business_services to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListTagResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		tags = append(tags, result.Tags...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}
	path := "/tags"
	if e != "" && eid != "" {
		path = "/" + e + "/" + eid + "/tags"
	}
	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet(path+queryParms.Encode(), responseHandler); err != nil {
		return nil, err
	}
	tagResponse.Tags = tags

	return tagResponse, nil
}

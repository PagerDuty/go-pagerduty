package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type Extension struct {
	APIObject
	Name             string      `json:"name"`
	EndpointURL      string      `json:"endpoint_url"`
	ExtensionObjects []APIObject `json:"extension_objects"`
	ExtensionSchema  APIObject   `json:"extension_schema"`
	Config           interface{} `json:"config"`
}

type ListExtensionResponse struct {
	APIListObject
	Extensions []Extension `json:"extensions"`
}

type ListExtensionOptions struct {
	APIListObject
	ExtensionObjectID string `url:"extension_object_id,omitempty"`
	ExtensionSchemaID string `url:"extension_schema_id,omitempty"`
	Query             string `url:"query,omitempty"`
}

func (c *Client) ListExtensions(o ListExtensionOptions) (*ListExtensionResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get("/extensions?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListExtensionResponse

	return &result, c.decodeJSON(resp, &result)
}

func (c *Client) CreateExtension(e *Extension) (*Extension, error) {
	resp, err := c.post("/extensions", e, nil)
	return getExtensionFromResponse(c, resp, err)
}

func (c *Client) DeleteExtension(id string) error {
	_, err := c.delete("/extensions/" + id)
	return err
}

func (c *Client) GetExtension(id string) (*Extension, error) {
	resp, err := c.get("/extensions/" + id)
	return getExtensionFromResponse(c, resp, err)
}

func (c *Client) UpdateExtension(id string, e *Extension) (*Extension, error) {
	resp, err := c.put("/extensions/"+id, e, nil)
	return getExtensionFromResponse(c, resp, err)
}

func getExtensionFromResponse(c *Client, resp *http.Response, err error) (*Extension, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]Extension
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	rootNode := "extension"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type ExtensionSchema struct {
	APIObject
	IconURL     string   `json:"icon_url"`
	LogoURL     string   `json:"logo_url"`
	Label       string   `json:"label"`
	Key         string   `json:"key"`
	Description string   `json:"description"`
	GuideURL    string   `json:"guide_url"`
	SendTypes   []string `json:"send_types"`
	URL         string   `json:"url"`
}

type ListExtensionSchemaResponse struct {
	APIListObject
	ExtensionSchemas []ExtensionSchema `json:"extension_schemas"`
}

type ListExtensionSchemaOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

func (c *Client) ListExtensionSchemas(o ListExtensionSchemaOptions) (*ListExtensionSchemaResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get("/extension_schemas?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListExtensionSchemaResponse

	return &result, c.decodeJSON(resp, &result)
}

func (c *Client) GetExtensionSchema(id string) (*ExtensionSchema, error) {
	resp, err := c.get("/extension_schemas/" + id)
	return getExtensionSchemaFromResponse(c, resp, err)
}

func getExtensionSchemaFromResponse(c *Client, resp *http.Response, err error) (*ExtensionSchema, error) {
	if err != nil {
		return nil, err
	}

	var target map[string]ExtensionSchema
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	rootNode := "extension_schema"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

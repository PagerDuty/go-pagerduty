package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// ExtensionSchema represnts the object presented by the API for each extension
// schema.
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

// ListExtensionSchemaResponse is the object presented in response to the
// request to list all extension schemas.
type ListExtensionSchemaResponse struct {
	APIListObject
	ExtensionSchemas []ExtensionSchema `json:"extension_schemas"`
}

// ListExtensionSchemaOptions are the options to send with the
// ListExtensionSchema reques(s).
type ListExtensionSchemaOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

// ListExtensionSchemas lists all of the extension schemas. Each schema
// represents a specific type of outbound extension. It's recommended to use
// ListExtensionSchemasWithContext instead.
func (c *Client) ListExtensionSchemas(o ListExtensionSchemaOptions) (*ListExtensionSchemaResponse, error) {
	return c.ListExtensionSchemasWithContext(context.Background(), o)
}

// ListExtensionSchemasWithContext lists all of the extension schemas. Each
// schema represents a specific type of outbound extension.
func (c *Client) ListExtensionSchemasWithContext(ctx context.Context, o ListExtensionSchemaOptions) (*ListExtensionSchemaResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/extension_schemas?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListExtensionSchemaResponse
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExtensionSchema gets a single extension schema. It's recommended to use
// GetExtensionSchemaWithContext instead.
func (c *Client) GetExtensionSchema(id string) (*ExtensionSchema, error) {
	return c.GetExtensionSchemaWithContext(context.Background(), id)
}

// GetExtensionSchemaWithContext gets a single extension schema.
func (c *Client) GetExtensionSchemaWithContext(ctx context.Context, id string) (*ExtensionSchema, error) {
	resp, err := c.get(ctx, "/extension_schemas/"+id)
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

	const rootNode = "extension_schema"

	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}

	return &t, nil
}

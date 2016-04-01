package pagerduty

import (
	"github.com/google/go-querystring/query"
)

type Addon struct {
	APIObject
	Name     string
	Src      string
	Services []APIObject
}

type ListAddonOptions struct {
	Includes   []string `url:"include,omitempty,brackets"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
	Filter     string   `url:"filter,omitempty"`
}

type ListAddonResponse struct {
	APIListObject
	Addons []Addon `json:"addons"`
}

func (c *Client) ListAddons(o ListAddonOptions) (*ListAddonResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/addons?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListAddonResponse
	return &result, c.decodeJson(resp, &result)
}

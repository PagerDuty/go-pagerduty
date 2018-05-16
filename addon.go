package pagerduty

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
)

// Addon is a third-party add-on to PagerDuty's UI.
type Addon struct {
	APIObject
	Name     string      `json:"name,omitempty"`
	Src      string      `json:"src,omitempty"`
	Services []APIObject `json:"services,omitempty"`
}

// ListAddonOptions are the options available when calling the ListAddons API endpoint.
type ListAddonOptions struct {
	APIListObject
	Includes   []string `url:"include,omitempty,brackets"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
	Filter     string   `url:"filter,omitempty"`
}

// ListAddonResponse is the response when calling the ListAddons API endpoint.
type ListAddonResponse struct {
	APIListObject
	Addons []Addon `json:"addons"`
}

// ListAddons lists all of the add-ons installed on your account.
func (c *PagerdutyClient) ListAddons(o ListAddonOptions) (*ListAddonResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get("/addons?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListAddonResponse
	return &result, DecodeJSON(resp, &result)
}

// InstallAddon installs an add-on for your account.
func (c *PagerdutyClient) InstallAddon(a Addon) (*Addon, error) {
	data := make(map[string]Addon)
	data["addon"] = a
	resp, err := c.Post("/addons", data)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create. HTTP Status code: %d", resp.StatusCode)
	}
	return getAddonFromResponse(c, resp)
}

// DeleteAddon deletes an add-on from your account.
func (c *PagerdutyClient) DeleteAddon(id string) error {
	_, err := c.Delete("/addons/" + id)
	return err
}

// GetAddon gets details about an existing add-on.
func (c *PagerdutyClient) GetAddon(id string) (*Addon, error) {
	resp, err := c.Get("/addons/" + id)
	if err != nil {
		return nil, err
	}
	return getAddonFromResponse(c, resp)
}

// UpdateAddon updates an existing add-on.
func (c *PagerdutyClient) UpdateAddon(id string, a Addon) (*Addon, error) {
	v := make(map[string]Addon)
	v["addon"] = a
	resp, err := c.Put("/addons/"+id, v, nil)
	if err != nil {
		return nil, err
	}
	return getAddonFromResponse(c, resp)
}

func getAddonFromResponse(c *PagerdutyClient, resp *http.Response) (*Addon, error) {
	var result map[string]Addon
	if err := DecodeJSON(resp, &result); err != nil {
		return nil, err
	}
	a, ok := result["addon"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have 'addon' field")
	}
	return &a, nil
}

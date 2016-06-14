package pagerduty

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-querystring/query"
)

type Addon struct {
	APIObject
	Name     string      `json:"name,omitempty"`
	Src      string      `json:"src,omitempty"`
	Services []APIObject `json:"services,omitempty"`
}

type ListAddonOptions struct {
	APIListObject
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

func (c *Client) InstallAddon(a Addon) error {
	data := make(map[string]Addon)
	data["addon"] = a
	resp, err := c.Post("/addons", data)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		ct, rErr := ioutil.ReadAll(resp.Body)
		if rErr == nil {
			log.Debug(string(ct))
		}
		return fmt.Errorf("Failed to create. HTTP Status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) DeleteAddon(id string) error {
	_, err := c.Delete("/addons/" + id)
	return err
}

func (c *Client) GetAddon(id string) (*Addon, error) {
	resp, err := c.Get("/addons/" + id)
	if err != nil {
		return nil, err
	}
	var result map[string]Addon
	if err := c.decodeJson(resp, &result); err != nil {
		return nil, err
	}
	a, ok := result["addon"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have 'addon' field")
	}
	return &a, nil
}

func (c *Client) UpdateAddon(id string, a Addon) error {
	v := make(map[string]Addon)
	v["addon"] = a
	_, err := c.Put("/addons/"+id, v)
	return err
}

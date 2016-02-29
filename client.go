package pagerduty

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io"
	"net/http"
)

// Client wraps http client
type Client struct {
	Subdomain string
	authToken string
}

// NewClient creates an API client
func NewClient(subdomain, authToken string) *Client {
	return &Client{
		Subdomain: subdomain,
		authToken: authToken,
	}
}

func (c *Client) Delete(path string) (*http.Response, error) {
	return c.Do("DELETE", path, nil)
}

func (c *Client) Post(path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.Do("POST", path, bytes.NewBuffer(data))
}

func (c *Client) Get(path string) (*http.Response, error) {
	return c.Do("GET", path, nil)
}

func (c *Client) Do(method, path string, body io.Reader) (*http.Response, error) {
	endpoint := "https://" + c.Subdomain + ".pagerduty.com/api/v1" + path
	log.Debugf("Endpoint", endpoint)
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token token="+c.authToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) decodeJson(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}

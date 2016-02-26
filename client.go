package pagerduty

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
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

// Do makes an http API call method is http method (GET, PUT etc..)
// path is rest endpoint (/incidents)
func (c *Client) Do(method, path string) (*http.Response, error) {
	endpoint := "https://" + c.Subdomain + ".pagerduty.com/api/v1" + path
	log.Debugf("Endpoint", endpoint)
	req, _ := http.NewRequest("GET", endpoint, nil)
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

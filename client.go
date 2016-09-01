package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const (
	apiEndpoint = "https://api.pagerduty.com"
)

// APIObject represents generic api json response that is shared by most
// domain object (like escalation
type APIObject struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// APIListObject are the fields used to control pagination when listing objects.
type APIListObject struct {
	Limit  uint `url:"limit,omitempty"`
	Offset uint `url:"offset,omitempty"`
	More   bool `url:"more,omitempty"`
	Total  uint `url:"total,omitempty"`
}

// APIReference are the fields required to reference another API object.
type APIReference struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type errorObject struct {
	Code   int      `json:"code,omitempty"`
	Mesage string   `json:"message,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// Client wraps http client
type Client struct {
	authToken string
}

// NewClient creates an API client
func NewClient(authToken string) *Client {
	return &Client{
		authToken: authToken,
	}
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.do("DELETE", path, nil)
}

func (c *Client) put(path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.do("PUT", path, bytes.NewBuffer(data))
}

func (c *Client) post(path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	log.Debugln(string(data))
	return c.do("POST", path, bytes.NewBuffer(data))
}

func (c *Client) get(path string) (*http.Response, error) {
	return c.do("GET", path, nil)
}

func (c *Client) do(method, path string, body io.Reader) (*http.Response, error) {
	endpoint := apiEndpoint + path
	log.Debugln("Endpoint:", endpoint)
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token token="+c.authToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if 199 >= resp.StatusCode || 300 <= resp.StatusCode {
		return resp, fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
	}
	return resp, nil
}

func (c *Client) decodeJSON(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}

func (c *Client) decodeObjectFromResponse(resp *http.Response, err error, data interface{}, target interface{}, rootNode string) (interface{}, error) {
	if err != nil {
		return nil, fmt.Errorf("Error calling the API endpoint: %v", err)
	}
	defer resp.Body.Close()
	if 199 >= resp.StatusCode || 300 <= resp.StatusCode {
		var eo *errorObject
		var getErr error
		if eo, getErr = c.getErrorFromResponse(resp); getErr != nil {
			return nil, fmt.Errorf("Response did not contain formatted error: %s. HTTP response code: %v. Raw response: %+v", getErr, resp.StatusCode, resp)
		}
		d, marshalErr := json.Marshal(data)
		if marshalErr != nil {
			return nil, fmt.Errorf("Could not marshal data '%+v': %s", data, marshalErr)
		}
		return nil, fmt.Errorf("Failed call API endpoint. Data: %v. Error: %v", string(d), eo)
	}
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	m, castOK := target.(map[string]interface{})
	if !castOK {
		return nil, fmt.Errorf("Could not cast target as map[string]interface{}")
	}
	t, nodeOK := m[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return t, nil
}

func (c *Client) getErrorFromResponse(resp *http.Response) (*errorObject, error) {
	var result map[string]errorObject
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", err)
	}
	s, ok := result["error"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have error field")
	}
	return &s, nil
}

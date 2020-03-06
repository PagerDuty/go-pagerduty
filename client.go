package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"time"
)

const (
	apiEndpoint = "https://api.pagerduty.com"
)

// The type of authentication to use with the API client
type authType int

const (
	// Account/user API token authentication
	apiToken authType = iota

	// OAuth token authentication
	oauthToken
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

type APIDetails struct {
	Type    string `json:"type,omitempty"`
	Details string `json:"details,omitempty"`
}

type errorObject struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}

// HTTPClient is an interface which declares the functionality we need from an
// HTTP client. This is to allow consumers to provide their own HTTP client as
// needed, without restricting them to only using *http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// defaultHTTPClient is our own default HTTP client. We use this, instead of
// http.DefaultClient, to avoid other packages tweaks to http.DefaultClient
// causing issues with our HTTP calls. This also allows us to tweak the
// transport values to be more resilient without making changes to the
// http.DefaultClient.
//
// Keep this unexported so consumers of the package can't make changes to it.
var defaultHTTPClient HTTPClient = newDefaultHTTPClient()

// Client wraps http client
type Client struct {
	authToken   string
	apiEndpoint string

	// Authentication type to use for API
	authType authType

	// HTTPClient is the HTTP client used for making requests against the
	// PagerDuty API. You can use either *http.Client here, or your own
	// implementation.
	HTTPClient HTTPClient
}

// NewClient creates an API client using an account/user API token
func NewClient(authToken string, options ...ClientOptions) *Client {
	client := Client{
		authToken:   authToken,
		apiEndpoint: apiEndpoint,
		authType:    apiToken,
		HTTPClient:  defaultHTTPClient,
	}

	for _, opt := range options {
		opt(&client)
	}

	return &client
}

// NewOAuthClient creates an API client using an OAuth token
func NewOAuthClient(authToken string, options ...ClientOptions) *Client {
	return NewClient(authToken, WithOAuth())
}

// ClientOptions allows for options to be passed into the Client for customization
type ClientOptions func(*Client)

// WithAPIEndpoint allows for a custom API endpoint to be passed into the the client
func WithAPIEndpoint(endpoint string) ClientOptions {
	return func(c *Client) {
		c.apiEndpoint = endpoint
	}
}

// WithOAuth allows for an OAuth token to be passed into the the client
func WithOAuth() ClientOptions {
	return func(c *Client) {
		c.authType = oauthToken
	}
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.do("DELETE", path, nil, nil)
}

func (c *Client) put(path string, payload interface{}, headers *map[string]string) (*http.Response, error) {

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		return c.do("PUT", path, bytes.NewBuffer(data), headers)
	}
	return c.do("PUT", path, nil, headers)
}

func (c *Client) post(path string, payload interface{}, headers *map[string]string) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.do("POST", path, bytes.NewBuffer(data), headers)
}

func (c *Client) get(path string) (*http.Response, error) {
	return c.do("GET", path, nil, nil)
}

func (c *Client) do(method, path string, body io.Reader, headers *map[string]string) (*http.Response, error) {
	endpoint := c.apiEndpoint + path
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("User-Agent", "go-pagerduty/"+Version)
	req.Header.Set("Content-Type", "application/json")

	switch c.authType {
	case oauthToken:
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	default:
		req.Header.Set("Authorization", "Token token="+c.authToken)
	}

	resp, err := c.HTTPClient.Do(req)
	return c.checkResponse(resp, err)
}

func (c *Client) decodeJSON(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}

func (c *Client) checkResponse(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, fmt.Errorf("Error calling the API endpoint: %v", err)
	}
	if 199 >= resp.StatusCode || 300 <= resp.StatusCode {
		var eo *errorObject
		var getErr error
		if eo, getErr = c.getErrorFromResponse(resp); getErr != nil {
			return resp, fmt.Errorf("Response did not contain formatted error: %s. HTTP response code: %v. Raw response: %+v", getErr, resp.StatusCode, resp)
		}
		return resp, fmt.Errorf("Failed call API endpoint. HTTP response code: %v. Error: %v", resp.StatusCode, eo)
	}
	return resp, nil
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

// responseHandler is capable of parsing a response. At a minimum it must
// extract the page information for the current page. It can also execute
// additional necessary handling; for example, if a closure, it has access
// to the scope in which it was defined, and can be used to append data to
// a specific slice. The responseHandler is responsible for closing the response.
type responseHandler func(response *http.Response) (APIListObject, error)

func (c *Client) pagedGet(basePath string, handler responseHandler) error {
	// Indicates whether there are still additional pages associated with request.
	var stillMore bool

	// Offset to set for the next page request.
	var nextOffset uint

	// While there are more pages, keep adjusting the offset to get all results.
	for stillMore, nextOffset = true, 0; stillMore; {
		response, err := c.do("GET", fmt.Sprintf("%s?offset=%d", basePath, nextOffset), nil, nil)
		if err != nil {
			return err
		}

		// Call handler to extract page information and execute additional necessary handling.
		pageInfo, err := handler(response)
		if err != nil {
			return err
		}

		// Bump the offset as necessary and set whether more results exist.
		nextOffset = pageInfo.Offset + pageInfo.Limit
		stillMore = pageInfo.More
	}

	return nil
}

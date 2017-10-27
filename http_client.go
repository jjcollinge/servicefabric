package servicefabric

import "net/http"

// HTTPClient is an interface for HTTP clients
// to implement. The client only requires
// read-only access to the Service Fabric API,
// thus only the HTTP GET method needs to be implemented
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Transport(transport http.RoundTripper)
	AddBasicAuth(username, password string)
}

// HTTPClientImpl is an implementation of HTTPClient
// that wraps the net/http HTTP client
type httpClientImpl struct {
	client   http.Client
	username string
	password string
}

// NewHTTPClient creates a new HTTPClient instance
func NewHTTPClient(client http.Client) HTTPClient {
	return &httpClientImpl{client: client}
}

// Get is a method that implements a HTTP GET request
func (c *httpClientImpl) Get(url string) (resp *http.Response, err error) {
	if c.username != "" && c.password != "" {
		req, _ := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(c.username, c.password)
		return c.client.Do(req)
	}
	return c.client.Get(url)
}

// Transport sets the HTTP client transport property
func (c *httpClientImpl) Transport(transport http.RoundTripper) {
	c.client.Transport = transport
}

func (c *httpClientImpl) AddBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

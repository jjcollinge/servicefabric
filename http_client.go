package servicefabric

import "net/http"

// HTTPClient is an interface for HTTP clients
// to implement. The client only requires
// read-only access to the Service Fabric API,
// thus only the HTTP GET method needs to be implemented
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Transport(transport *http.Transport)
}

// HTTPClientImpl is an implementation of HTTPClient
// that wraps the net/http HTTP client
type HTTPClientImpl struct {
	client http.Client
}

// Get is a method that implements a HTTP GET request
func (c *HTTPClientImpl) Get(url string) (resp *http.Response, err error) {
	return c.client.Get(url)
}

// Transport sets the HTTP client transport property
func (c *HTTPClientImpl) Transport(transport *http.Transport) {
	c.client.Transport = transport
}

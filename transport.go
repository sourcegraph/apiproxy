package apiproxy

import (
	"net/http"
)

// RevalidateFunc is the signature of functions used to determine whether the
// HTTP resource indicated by req should be revalidated by contacting the
// destination server.
type RevalidateFunc func(req *http.Request) bool

// RevalidationTransport is an implementation of net/http.RoundTripper that
// permits custom behavior with respect to cache entry revalidation for
// resources on the target server.
type RevalidationTransport struct {
	// Revalidate is called on each request in RoundTrip. If it returns true,
	// RoundTrip synthesizes and returns an HTTP 304 Not Modified response.
	// Otherwise, the request is passed through to the underlying transport.
	Revalidate RevalidateFunc

	// Transport is the underlying transport. If nil, net/http.DefaultTransport is used.
	Transport http.RoundTripper
}

// RoundTrip takes a Request and returns a Response.
func (t *RevalidationTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if t.Revalidate != nil && !t.Revalidate(req) {
		resp = &http.Response{
			Request:          req,
			TransferEncoding: req.TransferEncoding,
			StatusCode:       http.StatusNotModified,
		}
		return
	}

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	return transport.RoundTrip(req)
}

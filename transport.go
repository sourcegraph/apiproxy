package apiproxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// RevalidateFunc is the signature of functions used to determine whether the
// HTTP resource indicated by req should be revalidated by contacting the
// destination server.
type RevalidateFunc func(req *http.Request) bool

// RevalidationTransport is an implementation of net/http.RoundTripper that
// permits custom behavior with respect to cache entry revalidation for
// resources on the target server.
//
// If the request contains cache validators (an If-None-Match or
// If-Modified-Since header), then Revalidate is called to determine whether the
// cache entry should be revalidated (by being passed to the underlying
// transport). In this way, the Revalidate function can effectively extend or shorten cache
// age limits.
//
// If the request does not contain cache validators, then it is passed to the
// underlying transport.
type RevalidationTransport struct {
	// Revalidate is called on each request in RoundTrip. If it returns true,
	// RoundTrip synthesizes and returns an HTTP 304 Not Modified response.
	// Otherwise, the request is passed through to the underlying transport.
	Revalidate RevalidateFunc

	// Transport is the underlying transport. If nil, net/http.DefaultTransport is used.
	Transport http.RoundTripper
}

// NeverRevalidate is a RevalidateFunc for use with RevalidationTransport that
// causes HTTP requests to never revalidate cache entries. If a cache entry
// exists, it will always be used, even if it is expired.
func NeverRevalidate(_ *http.Request) bool {
	return false
}

// RoundTrip takes a Request and returns a Response.
func (t *RevalidationTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if t.Revalidate != nil && hasCacheValidator(req.Header) && !t.Revalidate(req) {
		resp = &http.Response{
			Request:          req,
			TransferEncoding: req.TransferEncoding,
			StatusCode:       http.StatusNotModified,
			Body:             ioutil.NopCloser(bytes.NewReader([]byte(""))),
		}
		return
	}

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	return transport.RoundTrip(req)
}

// hasCacheValidator returns true if the headers contain cache validators. See
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.3 for more
// information.
func hasCacheValidator(headers http.Header) bool {
	return headers.Get("if-none-match") != "" || headers.Get("if-modified-since") != ""
}

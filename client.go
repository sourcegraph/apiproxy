package apiproxy

import (
	"net/http"
	"net/textproto"
	"regexp"
	"strings"
	"sync"
)

// RequestModifyingTransport is an implementation of net/http.RoundTripper that
// allows headers to be overwritten on requests matching certain predicates.
//
// It gives more control over HTTP requests (e.g., caching) when using libraries
// whose only HTTP configuration point is a http.Client or http.RoundTripper.
type RequestModifyingTransport struct {
	overrides   map[*regexp.Regexp]requestOverride
	overridesMu sync.Mutex

	// Transport is the underlying transport. If nil, net/http.DefaultTransport
	// is used.
	Transport http.RoundTripper
}

// RoundTrip implements net/http.RoundTripper.
func (t *RequestModifyingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req = t.applyOverrides(req)

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	return transport.RoundTrip(req)
}

// requestOverride represents how a request should be modified by
// RequestModifyingTransport.
type requestOverride struct {
	setHeaders  http.Header
	runOnlyOnce bool
}

// Override instructs the transport to set the specified headers (overwriting
// existing headers of the same name) on a GET or HEAD request whose request URI
// matches the regexp. If runOnlyOnce is true, the override will be deleted
// after execution (and won't affect any future requests); otherwise, it will
// remain in effect for the lifetime of the transport.
func (t *RequestModifyingTransport) Override(requestURI *regexp.Regexp, setHeaders http.Header, runOnlyOnce bool) {
	t.overridesMu.Lock()
	defer t.overridesMu.Unlock()
	if t.overrides == nil {
		t.overrides = make(map[*regexp.Regexp]requestOverride)
	}
	t.overrides[requestURI] = requestOverride{setHeaders, runOnlyOnce}
}

// applyOverrides applies the transport's request overrides to req. If any
// overrides apply, req is cloned and the overrides are applied to the clone.
func (t *RequestModifyingTransport) applyOverrides(req *http.Request) *http.Request {
	// Only override GET and HEAD requests, just to be safe. We may want to
	// revisit this constraint later.
	if method := strings.ToUpper(req.Method); method != "GET" && method != "HEAD" {
		return req
	}

	requestURI := req.URL.RequestURI()

	t.overridesMu.Lock()
	defer t.overridesMu.Unlock()

	cloned := false
	for requestURIRegexp, override := range t.overrides {
		if requestURIRegexp.MatchString(requestURI) {
			if !cloned {
				req = cloneRequest(req)
				cloned = true
			}

			for name, val := range override.setHeaders {
				req.Header[textproto.CanonicalMIMEHeaderKey(name)] = val
			}

			if override.runOnlyOnce {
				delete(t.overrides, requestURIRegexp)
			}
		}
	}

	return req
}

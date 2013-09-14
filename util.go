package apiproxy

import (
	"net/http"
)

// cloneRequest returns a clone of the provided *http.Request. The clone is a
// shallow copy of the struct and its Header map. (This function copyright (c)
// goauth2 authors: https://code.google.com/p/goauth2.)
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

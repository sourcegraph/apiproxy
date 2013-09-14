package apiproxy

import (
	"github.com/gregjones/httpcache"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewCachingSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := NewSingleHostReverseProxy(target)
	proxy.Transport = httpcache.NewMemoryCacheTransport()
	return proxy
}

// NewSingleHostReverseProxy wraps net/http/httputil.NewSingleHostReverseProxy
// and sets the Host header.
func NewSingleHostReverseProxy(url *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(url)
	oldDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		oldDirector(r)
		r.Host = url.Host
	}
	return proxy
}

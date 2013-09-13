package apiproxy

import (
	"github.com/gregjones/httpcache"
	"net/http/httputil"
	"net/url"
)

func NewCachingSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = httpcache.NewMemoryCacheTransport()
	return proxy
}

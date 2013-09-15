package apiproxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewCachingSingleHostReverseProxy(t *testing.T) {
	targetRequestCount := 0
	targetResponseBody := []byte("qux")

	// Start the target server.
	targetMux := http.NewServeMux()
	targetMux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		targetRequestCount++
		w.Header().Add("Cache-Control", "max-age=60")
		w.Write(targetResponseBody)
	})
	target := httptest.NewServer(targetMux)
	defer target.Close()
	targetURL := mustParseURL(t, target.URL)

	// Start the reverse proxy.
	proxyMux := http.NewServeMux()
	proxyMux.Handle("/", NewCachingSingleHostReverseProxy(targetURL, nil))
	proxy := httptest.NewServer(proxyMux)
	defer proxy.Close()
	proxyURL := mustParseURL(t, proxy.URL)
	proxiedFooURL := proxyURL.ResolveReference(&url.URL{Path: "/foo"})

	// First request will hit target because the response has not been cached yet.
	res := httpGet(t, proxiedFooURL)
	resBody := readAll(t, res.Body)
	if want := 1; targetRequestCount != want {
		t.Errorf("want targetRequestCount == %d, got %d", want, targetRequestCount)
	}
	if !bytes.Equal(targetResponseBody, resBody) {
		t.Errorf("want response body == %q, got %q", targetResponseBody, resBody)
	}

	// Subsequent requests (within max-age) will hit cache.
	res = httpGet(t, proxiedFooURL)
	resBody = readAll(t, res.Body)
	if want := 1; targetRequestCount != want {
		t.Errorf("want targetRequestCount == %d, got %d", want, targetRequestCount)
	}
	if !bytes.Equal(targetResponseBody, resBody) {
		t.Errorf("want response body == %q, got %q", targetResponseBody, resBody)
	}
}

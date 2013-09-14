package apiproxy

import (
	"net/http"
	"regexp"
	"testing"
)

func TestRequestModifyingTransport_NonOverridden(t *testing.T) {
	mockTransport := newMockTransport()
	mockTransport.defaultResponse = &http.Response{}

	transport := &RequestModifyingTransport{
		Transport: mockTransport,
	}

	_, err := transport.RoundTrip(newHTTPGETRequest(t, "http://example.com"))
	if err != nil {
		t.Error("RoundTrip", err)
	}
	if numRequests := len(mockTransport.requests); numRequests != 1 {
		t.Errorf("want numRequests == %d, got %d", 1, numRequests)
	}
}

func TestRequestModifyingTransport_Overridden(t *testing.T) {
	mockTransport := newMockTransport()
	mockTransport.defaultResponse = &http.Response{}

	transport := &RequestModifyingTransport{Transport: mockTransport}
	transport.Override(regexp.MustCompile(`^/foo$`), http.Header{"X-Foo": []string{"bar"}}, false)

	req := newHTTPGETRequest(t, "http://example.com/foo")

	// The override will be applied for all matching requests (not just the
	// first) because runOnlyOnce is false.
	for i := 1; i <= 2; i++ {
		_, err := transport.RoundTrip(req)
		if err != nil {
			t.Error("RoundTrip", err)
		}
		if numRequests := len(mockTransport.requests); numRequests != i {
			t.Errorf("want numRequests == %d, got %d", i, numRequests)
		}
		if want, got := "bar", mockTransport.requests[i-1].Header.Get("X-Foo"); want != got {
			t.Errorf("want X-Foo header %q, got %q", want, got)
		}
	}
}

func TestRequestModifyingTransport_Overridden_RunOnlyOnce(t *testing.T) {
	mockTransport := newMockTransport()
	mockTransport.defaultResponse = &http.Response{}

	transport := &RequestModifyingTransport{Transport: mockTransport}
	transport.Override(regexp.MustCompile(`^/foo$`), http.Header{"X-Foo": []string{"bar"}}, true)

	req := newHTTPGETRequest(t, "http://example.com/foo")

	// The override will be applied the first time.
	_, err := transport.RoundTrip(req)
	if err != nil {
		t.Error("RoundTrip", err)
	}
	if numRequests := len(mockTransport.requests); numRequests != 1 {
		t.Errorf("want numRequests == %d, got %d", 1, numRequests)
	}
	if want, got := "bar", mockTransport.requests[0].Header.Get("X-Foo"); want != got {
		t.Errorf("want X-Foo header %q, got %q", want, got)
	}

	// The override will NOT be applied the second time because runOnlyOnce is true.
	_, err = transport.RoundTrip(req)
	if err != nil {
		t.Error("RoundTrip", err)
	}
	if numRequests := len(mockTransport.requests); numRequests != 2 {
		t.Errorf("want numRequests == %d, got %d", 2, numRequests)
	}
	if got := mockTransport.requests[1].Header.Get("X-Foo"); got != "" {
		t.Errorf("want overrides not to be applied, but they were (want X-Foo header empty, got %q)", got)
	}
}

func TestRequestModifyingTransport_GET_HEAD_Only(t *testing.T) {
	tests := []struct {
		method    string
		overrides bool
	}{
		{"GET", true}, {"HEAD", true},
		{"POST", false}, {"PUT", false}, {"DELETE", false},
	}
	for _, test := range tests {
		mockTransport := newMockTransport()
		mockTransport.defaultResponse = &http.Response{}

		transport := &RequestModifyingTransport{Transport: mockTransport}
		transport.Override(regexp.MustCompile(`^/foo$`), http.Header{"X-Foo": []string{"bar"}}, true)

		req, err := http.NewRequest(test.method, "http://example.com/foo", nil)
		if err != nil {
			t.Fatal("http.NewRequest", err)
		}

		_, err = transport.RoundTrip(req)
		var wantXFooHeader string
		if test.overrides {
			wantXFooHeader = "bar"
		}
		if got := mockTransport.requests[0].Header.Get("X-Foo"); wantXFooHeader != got {
			t.Errorf("%s: want overrides not to be applied, but they were (want X-Foo header %q, got %q)", test.method, wantXFooHeader, got)
		}
	}
}

package apiproxy

import (
	"errors"
	"net/http"
	"testing"
)

func TestRevalidationTransport_NilRevalidateFunc(t *testing.T) {
	mockTransport := newMockTransport()
	mockTransport.defaultResponse = &http.Response{}

	transport := &RevalidationTransport{
		Revalidate: nil,
		Transport:  mockTransport,
	}

	resp, err := transport.RoundTrip(&http.Request{})
	if err != nil {
		t.Error("RoundTrip", err)
	}
	if mockTransport.defaultResponse != resp {
		t.Errorf("want resp == %+v, got %+v", mockTransport.defaultResponse, resp)
	}
	if numRequests := len(mockTransport.requests); numRequests != 1 {
		t.Errorf("want numRequests == %d, got %d", 1, numRequests)
	}
}

func TestRevalidationTransport_RevalidateFunc(t *testing.T) {
	mockTransport := newMockTransport()
	mockTransport.defaultResponse = &http.Response{}

	transport := &RevalidationTransport{
		Revalidate: func(req *http.Request) bool {
			return false
		},
		Transport: mockTransport,
	}

	req := newHTTPGETRequest(t, "")
	req.Header.Add("if-none-match", `"foo"`)

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Error("RoundTrip", err)
	}
	if want := http.StatusNotModified; want != resp.StatusCode {
		t.Errorf("want resp.StatusCode == %d, got %d", want, resp.StatusCode)
	}
	if numRequests := len(mockTransport.requests); numRequests != 0 {
		t.Errorf("want numRequests == %d, got %d", 0, numRequests)
	}
}

func newMockTransport() *mockTransport {
	return &mockTransport{
		responses: make(map[*http.Request]*http.Response),
	}
}

type mockTransport struct {
	requests        []*http.Request
	responses       map[*http.Request]*http.Response
	defaultResponse *http.Response
}

func (t *mockTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.requests = append(t.requests, req)

	var present bool
	resp, present = t.responses[req]
	if present {
		return
	}

	resp = t.defaultResponse

	if resp == nil {
		err = errors.New("no mocked response")
	}
	return
}

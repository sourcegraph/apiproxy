package apiproxy

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func mustParseURL(t *testing.T, urlstr string) *url.URL {
	u, err := url.Parse(urlstr)
	if err != nil {
		t.Fatalf("mustParseURL %q: %s", urlstr, err)
	}
	return u
}

func httpGet(t *testing.T, url *url.URL) *http.Response {
	resp, err := http.Get(url.String())
	if err != nil {
		t.Fatal("http.Get", err)
	}
	return resp
}

func newHTTPGETRequest(t *testing.T, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal("http.NewRequest", err)
	}
	return req
}

func readAll(t *testing.T, rdr io.ReadCloser) []byte {
	defer rdr.Close()
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		t.Fatal("ReadAll", err)
	}
	return data
}

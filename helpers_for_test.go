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

func httpGet(t *testing.T, url *url.URL) (res *http.Response) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		t.Fatal("http.NewRequest", err)
	}

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("http.DefaultClient.Do", err)
	}
	return
}

func readAll(t *testing.T, rdr io.ReadCloser) []byte {
	defer rdr.Close()
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		t.Fatal("ReadAll", err)
	}
	return data
}

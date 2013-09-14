package apiproxy

import (
	"net/url"
	"regexp"
	"testing"
	"time"
)

func TestPathMatchValidator(t *testing.T) {
	v := &PathMatchValidator{
		regexp.MustCompile(`^/foo$`): 5 * time.Second,
		regexp.MustCompile(`^/qux$`): 10 * time.Second,
	}
	tests := []struct {
		path     string
		cacheAge time.Duration
		valid    bool
	}{
		{"/foo", 10 * time.Second, false}, {"/foo", 5 * time.Second, true},
		{"/qux", 15 * time.Second, false}, {"/qux", 10 * time.Second, true},
		{"/xyz", 0 * time.Second, false}, {"/xyz", 5 * time.Second, false},
	}
	for _, test := range tests {
		valid := v.Valid(&url.URL{Path: test.path}, test.cacheAge)
		if test.valid != valid {
			t.Errorf("path %s age %d: want valid == %v, got %v", test.path, test.cacheAge, test.valid, valid)
		}
	}
}

package apiproxy

import (
	"net/url"
	"regexp"
	"time"
)

// Validators are used to determine whether a cache entry for a URL is still
// valid at a certain age. They are used to set custom max-ages for cache
// entries (e.g., to extend the max-age of a certain API resource that an
// application knows to be rarely updated).
type Validator interface {
	Valid(url *url.URL, age time.Duration) bool
}

// ValidatorFunc is an adapter type to allow the use of ordinary functions as
// validators. If f is a function with the appropriate signature,
// ValidatorFunc(f) is a Validator object that calls f.
type ValidatorFunc func(url *url.URL, age time.Duration) bool

// Valid implements Validator.
func (f ValidatorFunc) Valid(url *url.URL, age time.Duration) bool {
	return f(url, age)
}

// NeverRevalidate is a Validator for use with RevalidationTransport that causes
// HTTP requests to never revalidate cache entries. If a cache entry exists, it
// will always be used, even if it is expired.
var NeverRevalidate = ValidatorFunc(func(_ *url.URL, _ time.Duration) bool {
	return true
})

// PathMatchValidator is a map of path regexps to the maximum age of resources
// matching one of those regexps.
type PathMatchValidator map[*regexp.Regexp]time.Duration

// Valid implements Validator.
func (v PathMatchValidator) Valid(url *url.URL, age time.Duration) bool {
	for re, maxAge := range v {
		if re.MatchString(url.Path) {
			return age <= maxAge
		}
	}
	return false
}

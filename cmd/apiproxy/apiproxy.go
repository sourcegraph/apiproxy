package main

import (
	"flag"
	"fmt"
	"github.com/sourcegraph/apiproxy"
	"github.com/sourcegraph/httpcache"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var bindAddr = flag.String("http", ":8080", "HTTP bind address for proxy")
var neverRevalidate = flag.Bool("never-revalidate", false, "never revalidate cached responses (use them regardless of age)")
var onlyRevalOlderThanStr = flag.String("only-revalidate-older-than", "", "only revalidate cached responses older than this duration (extends cache duration)")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "apiproxy proxies and mocks HTTP APIs.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "\tapiproxy [options] url\n\n")
		fmt.Fprintf(os.Stderr, "url is the base URL of the HTTP server to proxy.\n\n")
		fmt.Fprintf(os.Stderr, "The options are:\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Example usage:\n\n")
		fmt.Fprintf(os.Stderr, "\tTo run a caching proxy at http://localhost:8080 with target http://example.com:\n")
		fmt.Fprintf(os.Stderr, "\t    $ apiproxy http://example.com\n\n")
		fmt.Fprintf(os.Stderr, "\t... and never revalidate cached responses:\n")
		fmt.Fprintf(os.Stderr, "\t    $ apiproxy -never-revalidate http://example.com\n\n")
		fmt.Fprintf(os.Stderr, "\t... and only revalidate cached responses older than an hour:\n")
		fmt.Fprintf(os.Stderr, "\t    $ apiproxy -only-revalidate-older-than=1h http://example.com\n\n")
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
	}

	targetURL, err := url.Parse(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URL %q: %s\n", flag.Arg(0), err)
		os.Exit(1)
	}

	var onlyRevalOlderThan time.Duration
	if *onlyRevalOlderThanStr != "" {
		onlyRevalOlderThan, err = time.ParseDuration(*onlyRevalOlderThanStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse duration %q: %s\n", *onlyRevalOlderThanStr, err)
			os.Exit(1)
		}
	}

	proxy := apiproxy.NewCachingSingleHostReverseProxy(targetURL)
	cachingTransport := proxy.Transport.(*httpcache.Transport)
	cachingTransport.Transport = &apiproxy.RevalidationTransport{
		Check: apiproxy.ValidatorFunc(func(url *url.URL, age time.Duration) bool {
			if *neverRevalidate {
				return true
			}
			if *onlyRevalOlderThanStr != "" {
				return age <= onlyRevalOlderThan
			}
			return false
		}),
	}

	http.Handle("/", proxy)

	fmt.Fprintf(os.Stderr, "Starting proxy on %s with target %s\n", *bindAddr, targetURL.String())
	err = http.ListenAndServe(*bindAddr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}

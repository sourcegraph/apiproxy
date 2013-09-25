apiproxy
========

apiproxy is a proxy for HTTP/REST APIs with configurable cache timeouts, etc.

**Documentation:** <https://sourcegraph.com/github.com/sourcegraph/apiproxy/tree>

[![Build Status](https://travis-ci.org/sourcegraph/apiproxy.png?branch=master)](https://travis-ci.org/sourcegraph/apiproxy)
[![status](https://sourcegraph.com/api/repos/github.com/sourcegraph/apiproxy/badges/status.png)](https://sourcegraph.com/github.com/sourcegraph/apiproxy)
[![xrefs](https://sourcegraph.com/api/repos/github.com/sourcegraph/apiproxy/badges/xrefs.png)](https://sourcegraph.com/github.com/sourcegraph/apiproxy)
[![funcs](https://sourcegraph.com/api/repos/github.com/sourcegraph/apiproxy/badges/funcs.png)](https://sourcegraph.com/github.com/sourcegraph/apiproxy)
[![top func](https://sourcegraph.com/api/repos/github.com/sourcegraph/apiproxy/badges/top-func.png)](https://sourcegraph.com/github.com/sourcegraph/apiproxy)


Installation
------------

```bash
go get github.com/sourcegraph/apiproxy
```


Usage
-----

apiproxy supports 3 modes of usage: as a standalone server, as a Go client, and
as a Go HTTP server handler.


### As a standalone HTTP proxy server

Running apiproxy as a standalone HTTP proxy server lets you access it from any
HTTP client on any host.

The included `apiproxy` program runs a proxy server with a specified target URL:

```bash
$ go install github.com/sourcegraph/apiproxy/cmd/apiproxy
$ apiproxy http://api.example.com
2013/09/13 21:19:57 Starting proxy on :8080 with target http://api.example.com
```

Once launched, HTTP requests to http://localhost:8080 will be proxied to
http://api.example.com and the responses cached according to the HTTP standard.

See `apiproxy -h` for more information.


### As a Go client [`http.RoundTripper`](https://sourcegraph.com/code.google.com/p/go/symbols/go/code.google.com/p/go/src/pkg/net/http/RoundTripper:type)

Clients can use [`apiproxy.RevalidationTransport`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/RevalidationTransport:type) to modify the caching behavior of HTTP requests by setting a custom [`apiproxy.Validator`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/Validator:type) on the transport. The [`Validator`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/Validator:type) is used to determine whether a cache entry for a URL is still valid at a certain age.

A [`Validator`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/Validator:type) can be created by wrapping a `func(url *url.URL, age time.Duration) bool`
function with [`apiproxy.ValidatorFunc(...)`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/ValidatorFunc:type) or by using the built-in GitHub API implementation, [`MaxAge`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/service/github/MaxAge:type).

The [`RevalidationTransport`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/RevalidationTransport:type) can be used in an [`http.Client`](https://sourcegraph.com/code.google.com/p/go/symbols/go/code.google.com/p/go/src/pkg/net/http/Client:type) that is passed to external libraries, to give control over HTTP requests when using libraries whose only configuration point is an [`http.Client`](https://sourcegraph.com/code.google.com/p/go/symbols/go/code.google.com/p/go/src/pkg/net/http/Client:type).

The file [`service/github/client_test.go`](https://github.com/sourcegraph/apiproxy/blob/master/service/github/client_test.go)
contains a full example using the [go-github library](https://github.com/google/go-github), summarized here:

```go
transport := &apiproxy.RevalidationTransport{
  Transport: httpcache.NewMemoryCacheTransport(),
  Check: (&githubproxy.MaxAge{
    User:         time.Hour * 24,
    Repository:   time.Hour * 24,
    Repositories: time.Hour * 24,
    Activity:     time.Hour * 12,
  }).Validator(),
}
httpClient := &http.Client{Transport: transport}

client := github.NewClient(httpClient)
```

Now HTTP requests initiated by go-github will be subject to the caching policy set by the custom [`RevalidationTransport`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/RevalidationTransport:type).


### As a Go server [`http.Handler`](https://sourcegraph.com/code.google.com/p/go/symbols/go/code.google.com/p/go/src/pkg/net/http/Handler:type)

The function [`apiproxy.NewCachingSingleHostReverseProxy(target *url.URL, cache
Cache)
*httputil.ReverseProxy`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/NewCachingSingleHostReverseProxy)
returns a simple caching reverse proxy that you can use as an
`http.Handler`.

You can wrap the handler's `Transport` in an
['apiproxy.RevalidationTransport`](https://sourcegraph.com/github.com/sourcegraph/apiproxy/symbols/go/github.com/sourcegraph/apiproxy/RevalidationTransport:type)
to specify custom cache timeout behavior.

The file `cmd/apiproxy/apiproxy.go` contains a full example, summarized here:

```go
proxy := apiproxy.NewCachingSingleHostReverseProxy("https://api.github.com", httpcache.NewMemoryCache())
cachingTransport := proxy.Transport.(*httpcache.Transport)
cachingTransport.Transport = &apiproxy.RevalidationTransport{
  Check: apiproxy.ValidatorFunc(func(url *url.URL, age time.Duration) bool {
    // only revalidate expired cache entries older than 30 minutes
    return age > 30 * time.Minute
  }),
}
http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, proxy))
http.ListenAndServe(":8080", nil)
```


Examples
--------

The included `cmd/chirper/chirper.go` example program helps demonstrate
apiproxy's features. It returns a constantly updating JSON array of "chirps" at
the path `/chirps`.

1. Run `go run example/chirper/chirper.go` in one terminal window.
1. Install apiproxy: `go install github.com/sourcegraph/apiproxy/cmd/apiproxy`
1. Run `apiproxy -http=:8080 -never-revalidate http://localhost:9090` in another terminal window.

Now, let's make a request to the chirper API via apiproxy. Since this is our
first request, apiproxy will fetch the response from the chirper HTTP server.

```
$ curl http://localhost:8080/chirps
```

Notice that apiproxy hit the chirper HTTP server: `chirper` logged the message "Listed chirps".

But next time we make the same request, apiproxy won't need to hit chirper,
because the response has been cached and we are using the `-never-revalidate`
option to treat cache entries as though they never expire.

```
$ curl http://localhost:8080/chirps
```

Note that we didn't hit chirper (it didn't log "Listed chirps").

However, if we pass a `Cache-Control: no-cache` header in our request, apiproxy
will ignore the cache and hit chirper:

```
$ curl -H 'Cache-Control: no-cache' q:8080/chirps
```

Note that chirper logs "Listed chirps" after this request.


Cache backends
--------------

Any cache backend that implements
[`httpcache.Cache`](https://sourcegraph.com/github.com/gregjones/httpcache/symbols/go/github.com/gregjones/httpcache/Cache:type)
suffices, including:

* [`httpcache.MemoryCache`](https://sourcegraph.com/github.com/gregjones/httpcache/symbols/go/github.com/gregjones/httpcache/MemoryCache:type),
  instantiated with [`NewMemoryCache() *MemoryCache`](https://sourcegraph.com/github.com/gregjones/httpcache/symbols/go/github.com/gregjones/httpcache/NewMemoryCache)
* [`diskcache.Cache`](https://sourcegraph.com/github.com/gregjones/httpcache/symbols/go/github.com/gregjones/httpcache/diskcache/Cache:type), instantiated with [`diskcache.New(basePath string) *Cache`](https://sourcegraph.com/github.com/gregjones/httpcache/symbols/go/github.com/gregjones/httpcache/diskcache/New)
* [`s3cache.Cache`](https://sourcegraph.com/github.com/sourcegraph/s3cache/symbols/go/github.com/sourcegraph/s3cache/Cache:type), instantiated with [`s3cache.New(bucketURL string) *Cache`](https://sourcegraph.com/github.com/sourcegraph/s3cache/symbols/go/github.com/sourcegraph/s3cache/New) (requires env vars `S3_ACCESS_KEY` and `S3_SECRET_KEY`)

Contributing
------------

Patches and bug reports welcomed! Report issues and submit pull requests using
GitHub.

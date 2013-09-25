apiproxy
========

apiproxy is a proxy for HTTP/REST APIs with configurable cache timeouts, etc.

**Documentation:** <https://sourcegraph.com/github.com/sourcegraph/apiproxy>

[![Build Status](https://travis-ci.org/sourcegraph/apiproxy.png?branch=master)](https://travis-ci.org/sourcegraph/apiproxy)
[![status](https://sourcegraph.com/api/repos/github.com/sourcegraph/apiproxy/badges/status.png)](https://sourcegraph.com/github.com/sourcegraph/apiproxy)


Installation
------------

```bash
go get github.com/sourcegraph/apiproxy
```


Usage
-----

apiproxy supports 2 modes of usage: as a standalone server and as a Go library.


### As a standalone HTTP proxy server

Running apiproxy as a standalone HTTP proxy server lets you access it from any
HTTP client on any host.

The included `apiproxy` program runs a proxy server with a specified target URL:

```bash
$ apiproxy http://api.example.com
2013/09/13 21:19:57 Starting proxy on :8080 with target http://api.example.com
```

Once launched, HTTP requests to http://localhost:8080 will be proxied to
http://api.example.com and the responses cached according to the HTTP standard.

See `apiproxy -h` for more information.


### As a Go `http.RoundTripper`

Using apiproxy's
[`http.RoundTripper`](https://sourcegraph.com/code.google.com/p/go/symbols/go/code.google.com/p/go/src/pkg/net/http/RoundTripper:type)
lets you create HTTP clients and servers that proxy access to APIs. See
`cmd/apiproxy/apiproxy.go` for a usage example.


Examples
--------

The included `cmd/chirper/chirper.go` example program helps demonstrate
apiproxy's features. It returns a constantly updating JSON array of "chirps" at
the path `/chirps`.

1. Run `chirper` in one terminal window.
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

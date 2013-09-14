package githubproxy

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/sourcegraph/httpcache"
	"github.com/sourcegraph/apiproxy"
	"net/http"
	"regexp"
)

func ExampleGitHubClient_repositories() {
	cachingTransport := httpcache.NewMemoryCacheTransport()
	reqModifyingTransport := &apiproxy.RequestModifyingTransport{Transport: cachingTransport}
	httpClient := &http.Client{Transport: reqModifyingTransport}

	client := github.NewClient(httpClient)
	if exceededRateLimit(client) {
		return
	}

	// Get a repository.
	repo, resp, err := client.Repositories.Get("sourcegraph", "apiproxy")
	if err != nil {
		fmt.Printf("Error getting repository: %s\n", err)
		return
	}
	fmt.Printf("1st time: got repository %s info (from cache: %v).\n", *repo.Name, resp.Header.Get("X-From-Cache") == "1")

	// Get the same repository again. This request will hit the cache instead of
	// GitHub's API.
	repo, resp, err = client.Repositories.Get("sourcegraph", "apiproxy")
	if err != nil {
		fmt.Printf("Error getting repository: %s\n", err)
		return
	}
	fmt.Printf("2nd time: got repository %s info (from cache: %v).\n", *repo.Name, resp.Header.Get("X-From-Cache") == "1")

	// Once again, get the same repository, but override the request so we
	// bypass the cache and hit GitHub's API.
	reqModifyingTransport.Override(regexp.MustCompile(`^/repos/sourcegraph/apiproxy$`), apiproxy.NoCache, true)
	repo, resp, err = client.Repositories.Get("sourcegraph", "apiproxy")
	if err != nil {
		fmt.Printf("Error getting repository: %s\n", err)
		return
	}
	fmt.Printf("3nd time: got repository %s info (from cache: %v).\n", *repo.Name, resp.Header.Get("X-From-Cache") == "1")

	// Subsequent requests will hit the cache.
	repo, resp, err = client.Repositories.Get("sourcegraph", "apiproxy")
	if err != nil {
		fmt.Printf("Error getting repository: %s\n", err)
		return
	}
	fmt.Printf("4th time: got repository %s info (from cache: %v).\n", *repo.Name, resp.Header.Get("X-From-Cache") == "1")

	// Output:
	// 1st time: got repository apiproxy info (from cache: false).
	// 2nd time: got repository apiproxy info (from cache: true).
	// 3nd time: got repository apiproxy info (from cache: false).
	// 4th time: got repository apiproxy info (from cache: true).
}

func ExampleGitHubClient_events() {
	cachingTransport := httpcache.NewMemoryCacheTransport()
	reqModifyingTransport := &apiproxy.RequestModifyingTransport{Transport: cachingTransport}
	httpClient := &http.Client{Transport: reqModifyingTransport}

	client := github.NewClient(httpClient)
	if exceededRateLimit(client) {
		return
	}

	// List a user's events.
	events, resp, err := client.Activity.ListEventsPerformedByUser("sqs", true, nil)
	if err != nil {
		fmt.Printf("Error listing events: %s\n", err)
		return
	}
	fmt.Printf("1st time: got user %s events (from cache: %v).\n", *events[0].Actor.Login, resp.Header.Get("X-From-Cache") == "1")

	// List the same user's events again. This request will hit the cache
	// instead of GitHub's API.
	events, resp, err = client.Activity.ListEventsPerformedByUser("sqs", true, nil)
	if err != nil {
		fmt.Printf("Error listing events: %s\n", err)
		return
	}
	fmt.Printf("2nd time: got user %s events (from cache: %v).\n", *events[0].Actor.Login, resp.Header.Get("X-From-Cache") == "1")

	// Once again, list the same user's events, but override the request so we
	// bypass the cache and hit GitHub's API.
	reqModifyingTransport.Override(regexp.MustCompile(`^/users/sqs/events/public$`), apiproxy.NoCache, true)
	events, resp, err = client.Activity.ListEventsPerformedByUser("sqs", true, nil)
	if err != nil {
		fmt.Printf("Error listing events: %s\n", err)
		return
	}
	fmt.Printf("3nd time: got user %s events (from cache: %v).\n", *events[0].Actor.Login, resp.Header.Get("X-From-Cache") == "1")

	// Subsequent requests will hit the cache.
	events, resp, err = client.Activity.ListEventsPerformedByUser("sqs", true, nil)
	if err != nil {
		fmt.Printf("Error listing events: %s\n", err)
		return
	}
	fmt.Printf("4th time: got user %s events (from cache: %v).\n", *events[0].Actor.Login, resp.Header.Get("X-From-Cache") == "1")

	// Output:
	// 1st time: got user sqs events (from cache: false).
	// 2nd time: got user sqs events (from cache: true).
	// 3nd time: got user sqs events (from cache: false).
	// 4th time: got user sqs events (from cache: true).
}

func exceededRateLimit(client *github.Client) bool {
	rate, _, err := client.RateLimit()
	if err != nil {
		fmt.Printf("Error checking rate limit: %s\n", err)
		return false
	}
	// Check for a margin sufficient to run both examples.
	if rate.Remaining < 4 {
		fmt.Printf("Exceeded (or almost exceeded) GitHub API rate limit: %s. Try again later.\n", rate)
		return true
	}
	return false
}

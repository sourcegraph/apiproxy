package githubproxy

import (
	_ "github.com/google/go-github/github"
	"github.com/sourcegraph/apiproxy"
	"regexp"
	"time"
)

// MaxAge represents custom cache max-ages for GitHub API resources. It
// implements the apiproxy.Validator interface and is intended for use with
// RevalidationTransport. TODO(sqs): add fields for all GitHub API resources.
type MaxAge struct {
	User         time.Duration
	Repository   time.Duration
	Repositories time.Duration
	Activity     time.Duration
}

// Regexps for GitHub API resource paths. TODO(sqs): add regexps for all GitHub
// API resources.
var (
	publicRepos      = regexp.MustCompile(`^/repositories$`)
	repoPath         = regexp.MustCompile(`^/repos/[^/]+/[^/]+`)
	userPath         = regexp.MustCompile(`^/user(s/[^/]+)?$`)
	userPublicEvents = regexp.MustCompile(`^/users/[^/]+/events/public$`)
	userReposPath    = regexp.MustCompile(`^/user(s/[^/]+)?/repos$`)
)

// Validator returns an apiproxy.Validator that implements the MaxAge cache
// aging logic.
func (a *MaxAge) Validator() apiproxy.Validator {
	return &apiproxy.PathMatchValidator{
		publicRepos:      a.Repositories,
		repoPath:         a.Repository,
		userPath:         a.User,
		userPublicEvents: a.Activity,
		userReposPath:    a.Repositories,
	}
}

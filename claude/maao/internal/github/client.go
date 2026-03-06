package github

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	gh "github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

// GitHubClient wraps go-github with authentication and rate limit handling.
type GitHubClient struct {
	Client *gh.Client
	owner  string
	repo   string
}

// NewClient creates an authenticated GitHub client using an OAuth2 token.
func NewClient(token string) *GitHubClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	tc.Transport = &rateLimitTransport{base: tc.Transport}
	return &GitHubClient{
		Client: gh.NewClient(tc),
	}
}

// SetRepo configures the default owner and repo for operations.
func (c *GitHubClient) SetRepo(owner, repo string) {
	c.owner = owner
	c.repo = repo
}

// Owner returns the configured repository owner.
func (c *GitHubClient) Owner() string { return c.owner }

// Repo returns the configured repository name.
func (c *GitHubClient) Repo() string { return c.repo }

const (
	maxRetries     = 5
	baseRetryDelay = 1 * time.Second
)

type rateLimitTransport struct {
	base http.RoundTripper
}

func (t *rateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err = t.base.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			rateLimitErr := gh.CheckResponse(resp)
			if rateLimitErr != nil {
				if rlErr, ok := rateLimitErr.(*gh.RateLimitError); ok {
					waitUntil := rlErr.Rate.Reset.Time
					delay := time.Until(waitUntil)
					if delay <= 0 {
						delay = baseRetryDelay
					}
					select {
					case <-req.Context().Done():
						return nil, req.Context().Err()
					case <-time.After(delay):
						continue
					}
				}
				if _, ok := rateLimitErr.(*gh.AbuseRateLimitError); ok {
					delay := time.Duration(math.Pow(2, float64(attempt))) * baseRetryDelay
					select {
					case <-req.Context().Done():
						return nil, req.Context().Err()
					case <-time.After(delay):
						continue
					}
				}
			}
		}

		return resp, nil
	}

	return nil, fmt.Errorf("exceeded max retries due to rate limiting")
}

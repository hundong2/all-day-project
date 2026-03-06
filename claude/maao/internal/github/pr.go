package github

import (
	"context"

	gh "github.com/google/go-github/v68/github"
)

// PR wraps a GitHub pull request.
type PR = gh.PullRequest

// PRStatus represents the combined status of a PR.
type PRStatus struct {
	State       string // "open", "closed", "merged"
	Mergeable   bool
	ChecksPass  bool
	ReviewState string // "approved", "changes_requested", "pending"
}

// CreatePR creates a new pull request.
func (c *GitHubClient) CreatePR(ctx context.Context, title, body, head, base string) (*PR, error) {
	pr, _, err := c.Client.PullRequests.Create(ctx, c.owner, c.repo, &gh.NewPullRequest{
		Title: gh.String(title),
		Body:  gh.String(body),
		Head:  gh.String(head),
		Base:  gh.String(base),
	})
	return pr, err
}

// MergePR merges a pull request using the specified method ("merge", "squash", "rebase").
func (c *GitHubClient) MergePR(ctx context.Context, number int, method string) error {
	opts := &gh.PullRequestOptions{
		MergeMethod: method,
	}
	_, _, err := c.Client.PullRequests.Merge(ctx, c.owner, c.repo, number, "", opts)
	return err
}

// RequestReview requests reviews from specified users.
func (c *GitHubClient) RequestReview(ctx context.Context, number int, reviewers []string) error {
	_, _, err := c.Client.PullRequests.RequestReviewers(ctx, c.owner, c.repo, number, gh.ReviewersRequest{
		Reviewers: reviewers,
	})
	return err
}

// GetPRStatus returns the combined status of a pull request.
func (c *GitHubClient) GetPRStatus(ctx context.Context, number int) (*PRStatus, error) {
	pr, _, err := c.Client.PullRequests.Get(ctx, c.owner, c.repo, number)
	if err != nil {
		return nil, err
	}

	status := &PRStatus{
		State:     pr.GetState(),
		Mergeable: pr.GetMergeable(),
	}

	if pr.GetMerged() {
		status.State = "merged"
	}

	// Check combined status
	ref := pr.GetHead().GetSHA()
	if ref != "" {
		combined, _, err := c.Client.Repositories.GetCombinedStatus(ctx, c.owner, c.repo, ref, nil)
		if err == nil {
			status.ChecksPass = combined.GetState() == "success"
		}
	}

	// Check reviews
	reviews, _, err := c.Client.PullRequests.ListReviews(ctx, c.owner, c.repo, number, nil)
	if err == nil {
		status.ReviewState = "pending"
		for _, r := range reviews {
			switch r.GetState() {
			case "APPROVED":
				status.ReviewState = "approved"
			case "CHANGES_REQUESTED":
				status.ReviewState = "changes_requested"
			}
		}
	}

	return status, nil
}

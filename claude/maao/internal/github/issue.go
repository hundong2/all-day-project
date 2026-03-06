package github

import (
	"context"

	gh "github.com/google/go-github/v68/github"
)

// Issue wraps a GitHub issue.
type Issue = gh.Issue

// CreateIssue creates a new issue in the configured repository.
func (c *GitHubClient) CreateIssue(ctx context.Context, title, body string, labels []string, assignee string) (*Issue, error) {
	req := &gh.IssueRequest{
		Title:  gh.String(title),
		Body:   gh.String(body),
		Labels: &labels,
	}
	if assignee != "" {
		req.Assignee = gh.String(assignee)
	}
	issue, _, err := c.Client.Issues.Create(ctx, c.owner, c.repo, req)
	return issue, err
}

// GetIssue retrieves an issue by number.
func (c *GitHubClient) GetIssue(ctx context.Context, number int) (*Issue, error) {
	issue, _, err := c.Client.Issues.Get(ctx, c.owner, c.repo, number)
	return issue, err
}

// CloseIssue closes an issue by number.
func (c *GitHubClient) CloseIssue(ctx context.Context, number int) error {
	state := "closed"
	_, _, err := c.Client.Issues.Edit(ctx, c.owner, c.repo, number, &gh.IssueRequest{
		State: &state,
	})
	return err
}

// AddLabel adds a label to an issue.
func (c *GitHubClient) AddLabel(ctx context.Context, number int, label string) error {
	_, _, err := c.Client.Issues.AddLabelsToIssue(ctx, c.owner, c.repo, number, []string{label})
	return err
}

// ListIssues lists issues by state ("open", "closed", "all").
func (c *GitHubClient) ListIssues(ctx context.Context, state string) ([]*Issue, error) {
	var allIssues []*Issue
	opts := &gh.IssueListByRepoOptions{
		State: state,
		ListOptions: gh.ListOptions{
			PerPage: 50,
		},
	}

	for {
		issues, resp, err := c.Client.Issues.ListByRepo(ctx, c.owner, c.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, issue := range issues {
			if !issue.IsPullRequest() {
				allIssues = append(allIssues, issue)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}

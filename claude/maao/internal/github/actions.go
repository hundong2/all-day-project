package github

import (
	"context"
	"fmt"
	"time"

	gh "github.com/google/go-github/v68/github"
)

// ActionsManager manages GitHub Actions check runs for PRs.
type ActionsManager struct {
	client *GitHubClient
}

// NewActionsManager creates a new ActionsManager.
func NewActionsManager(client *GitHubClient) *ActionsManager {
	return &ActionsManager{client: client}
}

// WaitForChecks polls check runs until all required checks succeed or the timeout is reached.
func (am *ActionsManager) WaitForChecks(ctx context.Context, prNum int, requiredChecks []string, timeout time.Duration) error {
	deadline := time.After(timeout)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline:
			return fmt.Errorf("timeout waiting for checks on PR #%d", prNum)
		case <-ticker.C:
			done, err := am.checkStatus(ctx, prNum, requiredChecks)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}
}

func (am *ActionsManager) checkStatus(ctx context.Context, prNum int, requiredChecks []string) (bool, error) {
	pr, _, err := am.client.Client.PullRequests.Get(ctx, am.client.owner, am.client.repo, prNum)
	if err != nil {
		return false, fmt.Errorf("getting PR #%d: %w", prNum, err)
	}

	ref := pr.GetHead().GetSHA()
	if ref == "" {
		return false, fmt.Errorf("PR #%d has no head SHA", prNum)
	}

	checks, _, err := am.client.Client.Checks.ListCheckRunsForRef(ctx, am.client.owner, am.client.repo, ref, &gh.ListCheckRunsOptions{})
	if err != nil {
		return false, fmt.Errorf("listing check runs: %w", err)
	}

	results := make(map[string]string)
	for _, cr := range checks.CheckRuns {
		results[cr.GetName()] = cr.GetConclusion()
	}

	for _, required := range requiredChecks {
		conclusion, found := results[required]
		if !found {
			return false, nil // not yet reported
		}
		if conclusion == "failure" || conclusion == "cancelled" || conclusion == "timed_out" {
			return false, fmt.Errorf("check %q failed with conclusion: %s", required, conclusion)
		}
		if conclusion != "success" {
			return false, nil // still in progress
		}
	}

	return true, nil
}

package github

import (
	"context"
	"log"
	"time"

	gh "github.com/google/go-github/v68/github"
)

// EventType represents the type of change detected by the poller.
type EventType int

const (
	PlanMDChanged EventType = iota
	NewIssue
	NewComment
)

// PollEvent represents a change detected during polling.
type PollEvent struct {
	Type    EventType
	Payload interface{}
}

// Poller watches a GitHub repository for changes.
type Poller struct {
	client      *GitHubClient
	owner       string
	repo        string
	interval    time.Duration
	lastChecked time.Time

	lastPlanSHA  string
	lastIssueID  int64
	lastCommentT time.Time
}

// NewPoller creates a new repository poller.
func NewPoller(client *GitHubClient, owner, repo string, interval time.Duration) *Poller {
	return &Poller{
		client:       client,
		owner:        owner,
		repo:         repo,
		interval:     interval,
		lastChecked:  time.Now(),
		lastCommentT: time.Now(),
	}
}

// Start begins polling and returns a channel of detected events.
func (p *Poller) Start(ctx context.Context) <-chan PollEvent {
	ch := make(chan PollEvent, 32)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.poll(ctx, ch)
			}
		}
	}()
	return ch
}

func (p *Poller) poll(ctx context.Context, ch chan<- PollEvent) {
	p.checkPlanMD(ctx, ch)
	p.checkNewIssues(ctx, ch)
	p.checkNewComments(ctx, ch)
	p.lastChecked = time.Now()
}

func (p *Poller) checkPlanMD(ctx context.Context, ch chan<- PollEvent) {
	fc, _, resp, err := p.client.Client.Repositories.GetContents(
		ctx, p.owner, p.repo, "plan.md",
		&gh.RepositoryContentGetOptions{},
	)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return
		}
		log.Printf("poller: error checking plan.md: %v", err)
		return
	}
	if fc == nil {
		return
	}

	sha := fc.GetSHA()
	if p.lastPlanSHA == "" {
		p.lastPlanSHA = sha
		return
	}
	if sha != p.lastPlanSHA {
		p.lastPlanSHA = sha
		content, _ := fc.GetContent()
		ch <- PollEvent{Type: PlanMDChanged, Payload: content}
	}
}

func (p *Poller) checkNewIssues(ctx context.Context, ch chan<- PollEvent) {
	issues, _, err := p.client.Client.Issues.ListByRepo(ctx, p.owner, p.repo, &gh.IssueListByRepoOptions{
		State:     "open",
		Sort:      "created",
		Direction: "desc",
		Since:     p.lastChecked,
		ListOptions: gh.ListOptions{
			PerPage: 20,
		},
	})
	if err != nil {
		log.Printf("poller: error listing issues: %v", err)
		return
	}

	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}
		if issue.GetID() > p.lastIssueID {
			p.lastIssueID = issue.GetID()
			ch <- PollEvent{Type: NewIssue, Payload: issue}
		}
	}
}

func (p *Poller) checkNewComments(ctx context.Context, ch chan<- PollEvent) {
	comments, _, err := p.client.Client.Issues.ListComments(ctx, p.owner, p.repo, 0, &gh.IssueListCommentsOptions{
		Sort:      gh.String("created"),
		Direction: gh.String("desc"),
		Since:     &p.lastCommentT,
		ListOptions: gh.ListOptions{
			PerPage: 30,
		},
	})
	if err != nil {
		log.Printf("poller: error listing comments: %v", err)
		return
	}

	for _, comment := range comments {
		created := comment.GetCreatedAt().Time
		if created.After(p.lastCommentT) {
			p.lastCommentT = created
			ch <- PollEvent{Type: NewComment, Payload: comment}
		}
	}
}

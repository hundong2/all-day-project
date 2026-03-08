package token

import (
	"sync"
	"time"

	"github.com/maao/internal/store"
)

// TokenUsage represents token counts from a single API call.
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	CachedTokens int
}

// Tracker records per-execution token usage and caches daily totals.
type Tracker struct {
	store      *store.TokenStore
	mu         sync.Mutex
	dailyUsage map[string]int // agentName -> tokens used today
}

// NewTracker creates a new Tracker, pre-loading today's usage from the store.
func NewTracker(ts *store.TokenStore) *Tracker {
	return &Tracker{
		store:      ts,
		dailyUsage: make(map[string]int),
	}
}

// Record persists a token usage record and updates the in-memory daily cache.
func (t *Tracker) Record(agentName string, issueNum int, usage TokenUsage) error {
	record := store.TokenRecord{
		AgentName:    agentName,
		IssueNum:     issueNum,
		InputTokens:  usage.InputTokens,
		OutputTokens: usage.OutputTokens,
		CachedTokens: usage.CachedTokens,
	}

	if err := t.store.Record(record); err != nil {
		return err
	}

	t.mu.Lock()
	t.dailyUsage[agentName] += usage.InputTokens + usage.OutputTokens
	t.mu.Unlock()

	return nil
}

// GetDailyUsage returns the cached daily token usage for an agent.
// If the cache is empty, it fetches from the store.
func (t *Tracker) GetDailyUsage(agentName string) int {
	t.mu.Lock()
	usage, ok := t.dailyUsage[agentName]
	t.mu.Unlock()

	if ok {
		return usage
	}

	// Fall back to store query for today's usage.
	total, err := t.store.GetDailyUsage(agentName, time.Now())
	if err != nil {
		return 0
	}

	t.mu.Lock()
	t.dailyUsage[agentName] = total
	t.mu.Unlock()

	return total
}

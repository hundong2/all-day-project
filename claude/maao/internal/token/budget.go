package token

// TokenBudget represents the budget status for an agent.
type TokenBudget struct {
	AgentName  string
	DailyLimit int
	Used       int
	Remaining  int
}

// BudgetManager tracks daily token budgets per agent.
type BudgetManager struct {
	tracker *Tracker
	limits  map[string]int // agentName -> daily limit
}

// NewBudgetManager creates a BudgetManager with the given per-agent limits.
func NewBudgetManager(tracker *Tracker, limits map[string]int) *BudgetManager {
	return &BudgetManager{
		tracker: tracker,
		limits:  limits,
	}
}

// GetBudget returns the current budget status for an agent.
func (bm *BudgetManager) GetBudget(agentName string) TokenBudget {
	limit := bm.limits[agentName]
	used := bm.tracker.GetDailyUsage(agentName)
	remaining := limit - used
	if remaining < 0 {
		remaining = 0
	}
	return TokenBudget{
		AgentName:  agentName,
		DailyLimit: limit,
		Used:       used,
		Remaining:  remaining,
	}
}

// HasBudget checks whether an agent has enough remaining budget for the estimated cost.
func (bm *BudgetManager) HasBudget(agentName string, estimated int) bool {
	budget := bm.GetBudget(agentName)
	return budget.Remaining >= estimated
}

// GetStatus returns budget status for all configured agents.
func (bm *BudgetManager) GetStatus() map[string]TokenBudget {
	status := make(map[string]TokenBudget, len(bm.limits))
	for name := range bm.limits {
		status[name] = bm.GetBudget(name)
	}
	return status
}

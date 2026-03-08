package token

import "strings"

// Estimator estimates token usage for tasks based on heuristics.
type Estimator struct {
	smallLOC  int // threshold for small tasks
	mediumLOC int // threshold for medium tasks
	smallEst  int // estimated tokens for small tasks
	mediumEst int // estimated tokens for medium tasks
	largeEst  int // estimated tokens for large tasks
}

// NewEstimator creates an Estimator with default thresholds.
func NewEstimator() *Estimator {
	return &Estimator{
		smallLOC:  100,
		mediumLOC: 500,
		smallEst:  10000,
		mediumEst: 30000,
		largeEst:  80000,
	}
}

// EstimateByLOC returns an estimated token count based on lines of code.
func (e *Estimator) EstimateByLOC(loc int) int {
	switch {
	case loc < e.smallLOC:
		return e.smallEst
	case loc <= e.mediumLOC:
		return e.mediumEst
	default:
		return e.largeEst
	}
}

// EstimateByTitle returns an estimated token count based on task title keywords.
func (e *Estimator) EstimateByTitle(title string) int {
	lower := strings.ToLower(title)

	largeKeywords := []string{"refactor", "migration", "redesign", "rewrite", "architecture"}
	for _, kw := range largeKeywords {
		if strings.Contains(lower, kw) {
			return e.largeEst
		}
	}

	smallKeywords := []string{"typo", "fix", "bump", "rename", "lint", "format"}
	for _, kw := range smallKeywords {
		if strings.Contains(lower, kw) {
			return e.smallEst
		}
	}

	return e.mediumEst
}

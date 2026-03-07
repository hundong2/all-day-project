package store

import "time"

// TokenRecord represents a single token usage entry.
type TokenRecord struct {
	AgentName    string
	IssueNum     int
	InputTokens  int
	OutputTokens int
	CachedTokens int
	RecordedAt   time.Time
}

// TokenStore provides operations on token usage records.
type TokenStore struct {
	db *DB
}

// NewTokenStore creates a new TokenStore.
func NewTokenStore(db *DB) *TokenStore {
	return &TokenStore{db: db}
}

// Record inserts a new token usage record.
func (s *TokenStore) Record(record TokenRecord) error {
	_, err := s.db.db.Exec(`
		INSERT INTO token_usage (agent_name, issue_num, input_tokens, output_tokens, cached_tokens)
		VALUES (?, ?, ?, ?, ?)
	`, record.AgentName, record.IssueNum, record.InputTokens, record.OutputTokens, record.CachedTokens)
	return err
}

// GetDailyUsage returns total tokens used by an agent on a given date.
func (s *TokenStore) GetDailyUsage(agentName string, date time.Time) (int, error) {
	dateStr := date.Format("2006-01-02")
	var total int
	err := s.db.db.QueryRow(`
		SELECT COALESCE(SUM(input_tokens + output_tokens), 0)
		FROM token_usage
		WHERE agent_name = ? AND DATE(recorded_at) = ?
	`, agentName, dateStr).Scan(&total)
	return total, err
}

// GetHistory returns the most recent token usage records for an agent.
func (s *TokenStore) GetHistory(agentName string, limit int) ([]TokenRecord, error) {
	rows, err := s.db.db.Query(`
		SELECT agent_name, issue_num, input_tokens, output_tokens, cached_tokens, recorded_at
		FROM token_usage
		WHERE agent_name = ?
		ORDER BY recorded_at DESC
		LIMIT ?
	`, agentName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []TokenRecord
	for rows.Next() {
		var r TokenRecord
		if err := rows.Scan(&r.AgentName, &r.IssueNum, &r.InputTokens, &r.OutputTokens, &r.CachedTokens, &r.RecordedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

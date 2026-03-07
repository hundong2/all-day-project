package store

import (
	"database/sql"
	"encoding/json"
	"time"
)

// WorkflowState represents a persisted workflow state for a repository.
type WorkflowState struct {
	RepoURL   string
	Phase     string
	StateData map[string]interface{}
	UpdatedAt time.Time
}

// WorkflowStore provides operations on workflow states.
type WorkflowStore struct {
	db *DB
}

// NewWorkflowStore creates a new WorkflowStore.
func NewWorkflowStore(db *DB) *WorkflowStore {
	return &WorkflowStore{db: db}
}

// Save inserts or updates a workflow state.
func (s *WorkflowStore) Save(state WorkflowState) error {
	data, err := json.Marshal(state.StateData)
	if err != nil {
		return err
	}

	_, err = s.db.db.Exec(`
		INSERT INTO workflow_states (repo_url, phase, state_data, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(repo_url) DO UPDATE SET
			phase = excluded.phase,
			state_data = excluded.state_data,
			updated_at = CURRENT_TIMESTAMP
	`, state.RepoURL, state.Phase, string(data))
	return err
}

// Load retrieves a workflow state by repository URL.
func (s *WorkflowStore) Load(repoURL string) (*WorkflowState, error) {
	row := s.db.db.QueryRow(`
		SELECT repo_url, phase, state_data, updated_at
		FROM workflow_states WHERE repo_url = ?
	`, repoURL)

	var ws WorkflowState
	var raw string
	if err := row.Scan(&ws.RepoURL, &ws.Phase, &raw, &ws.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(raw), &ws.StateData); err != nil {
		return nil, err
	}
	return &ws, nil
}

// Delete removes a workflow state by repository URL.
func (s *WorkflowStore) Delete(repoURL string) error {
	_, err := s.db.db.Exec("DELETE FROM workflow_states WHERE repo_url = ?", repoURL)
	return err
}

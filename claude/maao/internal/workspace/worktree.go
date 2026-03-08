package workspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// WorktreeManager manages git worktrees for agents.
type WorktreeManager struct {
	repoPath    string
	worktreeDir string // e.g. ".worktrees"
	prefix      string // e.g. "agent/"
}

// NewWorktreeManager creates a new WorktreeManager.
func NewWorktreeManager(repoPath, worktreeDir, prefix string) *WorktreeManager {
	return &WorktreeManager{
		repoPath:    repoPath,
		worktreeDir: worktreeDir,
		prefix:      prefix,
	}
}

// CreateForAgent creates a worktree and branch for an agent working on a specific issue.
// Returns the absolute path to the new worktree.
func (wm *WorktreeManager) CreateForAgent(agentName string, issueNum int) (string, error) {
	branch := BranchName(wm.prefix, agentName, issueNum)
	wtPath := wm.worktreePath(agentName, issueNum)

	if err := os.MkdirAll(filepath.Dir(wtPath), 0o755); err != nil {
		return "", fmt.Errorf("creating worktree parent dir: %w", err)
	}

	cmd := exec.Command("git", "worktree", "add", wtPath, "-b", branch, "main")
	cmd.Dir = wm.repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git worktree add: %s: %w", string(output), err)
	}

	return wtPath, nil
}

// CleanupAgent removes the worktree and branch for an agent's issue.
func (wm *WorktreeManager) CleanupAgent(agentName string, issueNum int) error {
	wtPath := wm.worktreePath(agentName, issueNum)
	branch := BranchName(wm.prefix, agentName, issueNum)

	removeCmd := exec.Command("git", "worktree", "remove", wtPath, "--force")
	removeCmd.Dir = wm.repoPath
	if output, err := removeCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git worktree remove: %s: %w", string(output), err)
	}

	branchCmd := exec.Command("git", "branch", "-D", branch)
	branchCmd.Dir = wm.repoPath
	if output, err := branchCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git branch -D: %s: %w", string(output), err)
	}

	return nil
}

// WorktreePath returns the path for an agent's worktree.
func (wm *WorktreeManager) worktreePath(agentName string, issueNum int) string {
	dirName := fmt.Sprintf("%s-issue-%d", agentName, issueNum)
	return filepath.Join(wm.repoPath, wm.worktreeDir, dirName)
}

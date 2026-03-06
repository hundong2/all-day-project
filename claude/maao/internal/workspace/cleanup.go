package workspace

import (
	"fmt"
	"os/exec"
	"strings"
)

// CleanupMerged removes all worktrees and branches for merged issues.
// It lists worktrees, identifies ones matching the agent prefix, and removes those
// whose branches have been merged into main.
func (wm *WorktreeManager) CleanupMerged() ([]string, error) {
	// List merged branches matching the prefix
	cmd := exec.Command("git", "branch", "--merged", "main")
	cmd.Dir = wm.repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("listing merged branches: %s: %w", string(output), err)
	}

	var cleaned []string
	for _, line := range strings.Split(string(output), "\n") {
		branch := strings.TrimSpace(line)
		if branch == "" || branch == "main" || branch == "* main" {
			continue
		}
		if !strings.HasPrefix(branch, wm.prefix) {
			continue
		}

		// Remove associated worktree first
		pruneCmd := exec.Command("git", "worktree", "prune")
		pruneCmd.Dir = wm.repoPath
		_ = pruneCmd.Run()

		branchCmd := exec.Command("git", "branch", "-d", branch)
		branchCmd.Dir = wm.repoPath
		if branchOutput, err := branchCmd.CombinedOutput(); err != nil {
			return cleaned, fmt.Errorf("deleting branch %s: %s: %w", branch, string(branchOutput), err)
		}

		cleaned = append(cleaned, branch)
	}

	return cleaned, nil
}

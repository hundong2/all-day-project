package workspace

import "fmt"

// BranchName returns the standard branch name for an agent working on an issue.
// Format: {prefix}{agentName}/issue-{N}
func BranchName(prefix, agentName string, issueNum int) string {
	return fmt.Sprintf("%s%s/issue-%d", prefix, agentName, issueNum)
}

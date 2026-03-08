# Implement Issue

You are working on a GitHub issue as part of the MAAO orchestrated workflow.

## Issue Details

- **Issue:** #{{.IssueNumber}}
- **Title:** {{.IssueTitle}}
- **Body:**

{{.IssueBody}}

## Repository Context

- **Repository:** {{.RepoURL}}
- **Branch:** {{.WorkBranch}}
- **Base Branch:** {{.BaseBranch}}

{{if .AdditionalContext}}
## Additional Context

{{.AdditionalContext}}
{{end}}

## Instructions

1. Read and understand the issue requirements fully before writing code
2. Explore the existing codebase to understand conventions and patterns
3. Implement the changes described in the issue
4. Write tests for your implementation
5. Ensure all existing tests still pass
6. Keep your changes focused on this issue only

## Constraints

- Do NOT modify files unrelated to this issue
- Follow the existing code style and conventions
- Do NOT introduce new dependencies without clear justification
- Commit messages should reference the issue number: "fix #{{.IssueNumber}}: description"
- If you encounter blockers, document them clearly in your output

## Output

When complete, provide:
1. A summary of changes made
2. List of files modified/created
3. Any concerns or follow-up items

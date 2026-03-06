# Context Setup

Set up your working context for the assigned task.

## Repository

- **URL:** {{.RepoURL}}
- **Local Path:** {{.LocalPath}}
- **Branch:** {{.WorkBranch}}

## Task

You are about to work on Issue #{{.IssueNumber}}: {{.IssueTitle}}

## Instructions

1. Verify you are on the correct branch: `{{.WorkBranch}}`
2. Ensure the branch is up to date with `{{.BaseBranch}}`
3. Read the following key files to understand the project:
{{range .KeyFiles}}
   - `{{.}}`
{{end}}
4. Review any related recent changes:
{{range .RelatedIssues}}
   - Issue #{{.Number}}: {{.Title}} ({{.Status}})
{{end}}

## Project Structure

{{.ProjectStructure}}

## Conventions

{{.Conventions}}

## Output

Confirm that you have:
- [ ] Verified branch state
- [ ] Read and understood relevant source files
- [ ] Identified the scope of changes needed
- [ ] Noted any potential conflicts with ongoing work

Then provide a brief implementation plan (3-5 bullet points).

# Code Review

You are reviewing a pull request as part of the MAAO workflow.

## Pull Request

- **PR:** #{{.PRNumber}}
- **Title:** {{.PRTitle}}
- **Author Agent:** {{.AuthorAgent}}
- **Branch:** {{.HeadBranch}} -> {{.BaseBranch}}
- **Related Issue:** #{{.IssueNumber}}

## Diff

```diff
{{.Diff}}
```

## Instructions

Review the code changes against the following criteria:

### Correctness
- Does the code correctly implement the requirements from Issue #{{.IssueNumber}}?
- Are there any logic errors or edge cases not handled?
- Do the tests adequately cover the changes?

### Code Quality
- Does the code follow project conventions and style?
- Are there any code smells or anti-patterns?
- Is the code readable and well-structured?

### Security
- Are there any security vulnerabilities (injection, XSS, etc.)?
- Are secrets or credentials properly handled?
- Are inputs validated at system boundaries?

### Performance
- Are there any obvious performance issues?
- Are database queries or API calls efficient?

## Output Format

```json
{
  "verdict": "approve|request_changes|comment",
  "summary": "One-line summary of review",
  "comments": [
    {
      "file": "path/to/file.go",
      "line": 42,
      "body": "Review comment text",
      "severity": "critical|warning|suggestion|nitpick"
    }
  ],
  "blocking_issues": ["List of issues that must be fixed before merge"],
  "suggestions": ["Optional improvements that are not blocking"]
}
```

## Guidelines

- Be constructive and specific
- Distinguish between blocking issues and nice-to-haves
- If the PR looks good, approve it promptly
- Focus on substantive issues, not style preferences already handled by linters

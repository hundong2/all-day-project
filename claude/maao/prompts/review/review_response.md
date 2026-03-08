# Review Response

You received code review feedback on your pull request. Address the feedback and update your code.

## Review Feedback

**Reviewer:** {{.ReviewerAgent}}
**Verdict:** {{.Verdict}}
**Round:** {{.ReviewRound}} of {{.MaxRounds}}

### Comments

{{range .Comments}}
#### {{.File}}:{{.Line}} [{{.Severity}}]
{{.Body}}

{{end}}

{{if .BlockingIssues}}
### Blocking Issues
{{range .BlockingIssues}}
- {{.}}
{{end}}
{{end}}

## Instructions

1. Address ALL blocking issues (critical and warning severity)
2. Consider suggestions but use your judgment on whether to implement them
3. For each comment, either:
   - Fix the issue and note what you changed
   - Explain why no change is needed (with justification)
4. Run tests after making changes to ensure nothing is broken
5. Keep changes focused on addressing the review feedback

## Output

Provide a response for each review comment:

```json
{
  "responses": [
    {
      "file": "path/to/file.go",
      "line": 42,
      "action": "fixed|acknowledged|declined",
      "explanation": "What was changed or why no change was made"
    }
  ],
  "files_modified": ["list of files changed"],
  "tests_passed": true
}
```

## Guidelines

- Fix blocking issues first
- Do not introduce new features while addressing review feedback
- If a fix requires significant refactoring, note it and discuss before proceeding
- Commit with message: "review: address feedback on PR #{{.PRNumber}}"

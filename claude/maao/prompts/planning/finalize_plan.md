# Finalize Plan

You are the PM agent finalizing the project plan after {{.TotalRounds}} discussion rounds.

## Discussion Summary

{{.DiscussionSummary}}

## All Proposed Issues

{{.ProposedIssues}}

## Instructions

1. Incorporate feedback from all discussion rounds
2. Resolve any conflicting opinions (prefer consensus, otherwise use your judgment)
3. Produce the final list of GitHub issues to create
4. Assign a definitive priority order
5. Assign each issue to the most suitable agent
6. Verify the dependency graph has no cycles

## Output Format

Return the finalized JSON array:

```json
{
  "issues": [
    {
      "title": "Issue title",
      "body": "Final description with acceptance criteria",
      "labels": ["enhancement"],
      "complexity": "small|medium|large",
      "assigned_agent": "claude|gemini|codex|copilot",
      "depends_on": [],
      "priority": 1
    }
  ],
  "milestones": [
    {
      "title": "Milestone name",
      "description": "Milestone description",
      "issues": [1, 2, 3]
    }
  ],
  "estimated_total_tokens": 150000,
  "notes": "Any important notes for execution"
}
```

## Guidelines

- Ensure every issue has clear, testable acceptance criteria
- Balance workload across agents considering their token budgets
- Place quick wins and unblocking tasks at highest priority
- Include a "notes" field for any caveats or special instructions

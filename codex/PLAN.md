# All Day Project Plan

## Assumption
Because the repository had no existing product code and `PLAN.md` was empty, this plan proceeds under an explicit assumption: build a small, dependency-free Python CLI called **All Day** for managing a daily task list.

## Goals
- Support `add`, `list`, `done`, `remove`, and `summary` commands.
- Support `update` for correcting task details without deleting and re-adding.
- Support tags and simple text search for faster task filtering.
- Store data locally in JSON with no third-party dependencies.
- Keep the project easy to run and easy to test.

## Scope
- Single-user local CLI
- JSON-backed persistence
- Simple terminal output
- Automated tests with `unittest`

## Out of Scope
- Web UI
- Cloud sync
- Multi-user collaboration
- External database or API integrations

## Acceptance Criteria
1. Users can add tasks with a title, optional note, priority, and plan date.
2. Users can optionally attach tags to tasks.
3. Users can list tasks filtered by date, status, tag, and simple text query.
4. Users can update an existing task's title, note, priority, date, and tags.
5. Users can mark tasks complete and remove tasks by ID.
6. Users can view a per-day summary with total, pending, completed, and high-priority pending counts.
7. The application works with only the Python standard library.
8. Automated tests pass and usage is documented.

## Implementation Steps
1. Create package structure and CLI contract.
2. Implement JSON storage and task lifecycle rules.
3. Implement CLI commands and formatted output.
4. Add automated tests for normal flows and edge cases.
5. Document usage and verification steps.

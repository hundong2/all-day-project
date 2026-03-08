# All Day CLI Implementation Plan

## Context
Brand-new repository. Build a dependency-free Python CLI named **All Day** that stores task data in a local JSON file and supports:
- adding tasks
- listing tasks
- marking tasks done
- removing tasks
- showing a daily summary

## Work Objectives
- Deliver a small, dependency-free CLI runnable with the Python standard library only.
- Keep data persistence local, transparent, and easy to inspect.
- Establish a minimal project structure that is easy to test and extend.

## Guardrails
### Must Have
- Python standard library only
- Local JSON-backed persistence
- Clear CLI commands for add, list, done, remove, and summary
- Basic validation and user-friendly error messages
- Lightweight automated test coverage for core flows

### Must NOT Have
- Third-party packages
- Network/database dependencies
- Multi-user sync or cloud features
- Over-engineered architecture for a single-user local CLI

## Task Flow
1. Set up the minimal project skeleton and CLI entrypoint.
2. Implement JSON storage and task domain behavior.
3. Wire commands and output formatting.
4. Add tests and usage documentation.

## Detailed TODOs

### 1. Create the project skeleton and CLI contract
**Acceptance criteria:**
- Repository contains a clear Python package/module layout and a documented entrypoint.
- Command surface for `add`, `list`, `done`, `remove`, and `summary` is defined before implementation begins.
- Storage file location strategy is documented (for example, default local JSON path with predictable creation behavior).

### 2. Implement local JSON persistence and task lifecycle rules
**Acceptance criteria:**
- Tasks can be created, loaded, updated, and deleted using a local JSON file only.
- Stored task records include enough fields to support status tracking and daily summary behavior.
- Corrupt/missing file handling has explicit behavior (initialize new store, fail gracefully, or both as designed).

### 3. Implement CLI commands and user-facing output
**Acceptance criteria:**
- `add` creates a new task from CLI input.
- `list` shows tasks in a readable format with done/pending state.
- `done` marks a selected task complete.
- `remove` deletes a selected task.
- `summary` reports at least today's total, completed, and remaining tasks.
- Invalid command usage returns helpful guidance instead of raw tracebacks.

### 4. Add tests for core workflows and persistence edge cases
**Acceptance criteria:**
- Automated tests cover add/list/done/remove/summary flows against temporary JSON files.
- Tests cover at least one error-path or edge case (missing store, empty list, invalid task id, or malformed JSON strategy).
- Test execution is documented and passes locally with the standard Python test runner approach used by the project.

### 5. Document usage and delivery checks
**Acceptance criteria:**
- README or equivalent usage doc explains installation-free execution, command examples, and data file behavior.
- Final verification checklist includes manual smoke steps for each command.
- A new contributor can run the CLI and tests from repository instructions without guessing missing setup.

## Success Criteria
- A user can manage a daily task list entirely from the terminal with no external dependencies.
- Data survives across CLI runs via the JSON file.
- Core flows are covered by tests and documented with clear examples.

## Risks
- **Task identity ambiguity:** If IDs are not stable and human-friendly, `done`/`remove` may be awkward or error-prone.
- **Date boundary confusion:** Daily summary behavior can become inconsistent if “today” is not clearly defined and tested.
- **JSON corruption handling:** A malformed store file could block the CLI unless recovery behavior is designed up front.
- **Output usability drift:** Even a small CLI can feel cumbersome if list/summary formatting is not kept simple and readable.

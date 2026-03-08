# Autopilot Spec: All Day CLI

## Product Summary
All Day is a dependency-free Python CLI for planning a day's tasks in a local JSON file.

## User Story
As a single user working in the terminal, I want to add and manage my day plan quickly so I can track pending and completed work without setting up any external service.

## Core Features
- Add a task with title, optional note, priority, and plan date
- List tasks by date and status
- Mark a task done
- Remove a task
- Show a daily summary

## Technical Constraints
- Python standard library only
- Local JSON persistence
- Deterministic tests using temporary storage files

## Key Assumptions
- The repo has no grounded pre-existing product direction
- A small CLI is the safest autonomous default for an empty repo
- `plan_date` defaults to today's local date

# All Day

All Day is a dependency-free Python CLI for managing a daily task list from the terminal.

## Features
- Add tasks with title, priority, note, and plan date
- Add one or more tags to tasks
- Update an existing task without recreating it
- List tasks by date, status, tag, and simple text query
- Mark tasks complete
- Remove tasks
- View a daily summary

## Requirements
- Python 3.11+

## Quick Start

```bash
python3 -m all_day --storage ./.demo/tasks.json add "Write report" --priority high --tag work --tag writing
python3 -m all_day --storage ./.demo/tasks.json add "Take a walk" --note "15 minutes" --tag health
python3 -m all_day --storage ./.demo/tasks.json list --status all
python3 -m all_day --storage ./.demo/tasks.json list --tag work --query report
python3 -m all_day --storage ./.demo/tasks.json update 1 --note "First draft done" --tag work --tag review
python3 -m all_day --storage ./.demo/tasks.json done 1
python3 -m all_day --storage ./.demo/tasks.json summary
```

If installed from the built wheel, you can run the same commands with:

```bash
all-day list --status all
```

If `--storage` is omitted, All Day stores data at:

```text
~/.all_day/tasks.json
```

You can also override the storage path with:

```bash
ALL_DAY_STORAGE=/path/to/tasks.json python3 -m all_day list
```

## Commands

### Add

```bash
python3 -m all_day add "Buy groceries" --priority high --note "Milk and fruit" --tag errands --tag home
```

Options:
- `--priority {low,medium,high}` (default: `medium`)
- `--note TEXT`
- `--date YYYY-MM-DD` (default: today)
- `--tag TAG` (repeatable)

### List

```bash
python3 -m all_day list --status pending --date 2026-03-06 --tag work --query draft
```

Options:
- `--status {pending,done,all}` (default: `pending`)
- `--date YYYY-MM-DD` (default: today)
- `--tag TAG` (filter by a single tag)
- `--query TEXT` (matches title, note, or tags)

### Done

```bash
python3 -m all_day done 1
```

### Update

```bash
python3 -m all_day update 1 --title "Buy groceries and fruit" --tag errands --tag home
```

Options:
- `--title TEXT`
- `--priority {low,medium,high}`
- `--note TEXT` (pass empty string to clear)
- `--date YYYY-MM-DD`
- `--tag TAG` (replaces tags; repeatable)
- `--clear-tags`

### Remove

```bash
python3 -m all_day remove 1
```

### Summary

```bash
python3 -m all_day summary --date 2026-03-06
```

## Verification

Run:

```bash
python3 -m compileall all_day tests
python3 -m unittest discover -s tests -v
```

## Make Targets

You can also use `make` for common project commands:

```bash
make install
make run
make run ARGS='add "Write report" --priority high --tag work'
make run ARGS='update 1 --title "Final report" --clear-tags'
make run ARGS='list --status all'
make add TITLE='Write report' ADD_ARGS='--priority high --tag work'
make list LIST_ARGS='--status all'
make update TASK_ID=1 UPDATE_ARGS='--title "Final report" --clear-tags'
make done TASK_ID=1
make remove TASK_ID=1
make summary
make stop
make test
make build
```

Notes:
- `make run` uses `STORAGE=./.demo/tasks.json` by default.
- Override the storage file with `make run STORAGE=./.demo/my-tasks.json ARGS='list --status all'`.
- Because All Day is a one-shot CLI, `make stop` clears the demo storage file instead of stopping a background server.

## Manual Smoke Checklist
- Add two tasks
- Update one task
- List pending tasks
- Filter by tag and query
- Mark one task done
- Confirm `list --status done` shows it
- Run `summary` and check counts
- Remove a task and verify it no longer appears

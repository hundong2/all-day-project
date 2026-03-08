from __future__ import annotations

import argparse
from typing import Sequence

from .service import (
    MISSING,
    TaskNotFoundError,
    TaskService,
    validate_plan_date,
    validate_priority,
    validate_status,
)
from .storage import StoreError, TaskStore, resolve_storage_path


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="all-day",
        description="Manage a daily task list from the terminal.",
    )
    parser.add_argument(
        "--storage",
        help="Path to the JSON storage file. Defaults to ~/.all_day/tasks.json.",
    )

    subparsers = parser.add_subparsers(dest="command", required=True)

    add_parser = subparsers.add_parser("add", help="Add a new task.")
    add_parser.add_argument("title", help="Task title.")
    add_parser.add_argument(
        "--priority",
        default="medium",
        choices=["low", "medium", "high"],
        help="Task priority.",
    )
    add_parser.add_argument("--note", default="", help="Optional task note.")
    add_parser.add_argument(
        "--tag",
        action="append",
        default=[],
        help="Tag to attach to the task. Repeat to add multiple tags.",
    )
    add_parser.add_argument(
        "--date",
        default=None,
        help="Plan date in YYYY-MM-DD format. Defaults to today.",
    )

    list_parser = subparsers.add_parser("list", help="List tasks for a date.")
    list_parser.add_argument(
        "--status",
        default="pending",
        choices=["pending", "done", "all"],
        help="Filter by status.",
    )
    list_parser.add_argument(
        "--tag",
        default=None,
        help="Filter by a single tag.",
    )
    list_parser.add_argument(
        "--query",
        default=None,
        help="Search title, note, and tags for a text match.",
    )
    list_parser.add_argument(
        "--date",
        default=None,
        help="Plan date in YYYY-MM-DD format. Defaults to today.",
    )

    done_parser = subparsers.add_parser("done", help="Mark a task as complete.")
    done_parser.add_argument("task_id", type=int, help="Task ID.")

    update_parser = subparsers.add_parser(
        "update",
        help="Update an existing task.",
    )
    update_parser.add_argument("task_id", type=int, help="Task ID.")
    update_parser.add_argument("--title", default=None, help="New task title.")
    update_parser.add_argument(
        "--priority",
        default=None,
        choices=["low", "medium", "high"],
        help="New task priority.",
    )
    update_parser.add_argument(
        "--note",
        default=None,
        help="New note. Pass an empty string to clear it.",
    )
    tag_group = update_parser.add_mutually_exclusive_group()
    tag_group.add_argument(
        "--tag",
        action="append",
        default=None,
        help="Replace task tags. Repeat to set multiple tags.",
    )
    tag_group.add_argument(
        "--clear-tags",
        action="store_true",
        help="Remove all tags from the task.",
    )
    update_parser.add_argument(
        "--date",
        default=None,
        help="New plan date in YYYY-MM-DD format.",
    )

    remove_parser = subparsers.add_parser("remove", help="Remove a task.")
    remove_parser.add_argument("task_id", type=int, help="Task ID.")

    summary_parser = subparsers.add_parser(
        "summary",
        help="Show a summary for a date.",
    )
    summary_parser.add_argument(
        "--date",
        default=None,
        help="Plan date in YYYY-MM-DD format. Defaults to today.",
    )

    return parser


def create_service(storage_path: str | None) -> TaskService:
    return TaskService(TaskStore(resolve_storage_path(storage_path)))


def format_task(task) -> str:
    status_icon = "x" if task.status == "done" else " "
    parts = [
        f"#{task.id}",
        f"[{status_icon}]",
        task.title,
        f"({task.priority})",
    ]
    line = " ".join(parts)
    if task.tags:
        line += f"\n    tags: {', '.join(task.tags)}"
    if task.note:
        line += f"\n    note: {task.note}"
    return line


def main(argv: Sequence[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)

    try:
        service = create_service(args.storage)

        if args.command == "add":
            task = service.add_task(
                args.title,
                priority=validate_priority(args.priority),
                tags=args.tag,
                note=args.note,
                plan_date=validate_plan_date(args.date),
            )
            print(
                f"Added task #{task.id} for {task.plan_date}: {task.title}"
            )
            return 0

        if args.command == "list":
            plan_date = validate_plan_date(args.date)
            tasks = service.list_tasks(
                plan_date=plan_date,
                status=validate_status(args.status),
                tag=args.tag,
                query=args.query,
            )
            if not tasks:
                print(f"No tasks found for {plan_date}.")
                return 0
            print(f"Tasks for {plan_date}:")
            for task in tasks:
                print(format_task(task))
            return 0

        if args.command == "done":
            task = service.mark_done(args.task_id)
            print(f"Completed task #{task.id}: {task.title}")
            return 0

        if args.command == "update":
            if not any(
                value is not None or key == "clear_tags" and value
                for key, value in {
                    "title": args.title,
                    "priority": args.priority,
                    "note": args.note,
                    "tag": args.tag,
                    "date": args.date,
                    "clear_tags": args.clear_tags,
                }.items()
            ):
                raise ValueError("Provide at least one field to update.")

            tags = args.tag
            if args.clear_tags:
                tags = []

            task = service.update_task(
                args.task_id,
                title=args.title if args.title is not None else MISSING,
                priority=(
                    validate_priority(args.priority)
                    if args.priority is not None
                    else MISSING
                ),
                note=args.note if args.note is not None else MISSING,
                tags=tags if tags is not None else MISSING,
                plan_date=(
                    validate_plan_date(args.date)
                    if args.date is not None
                    else MISSING
                ),
            )
            print(f"Updated task #{task.id}: {task.title}")
            return 0

        if args.command == "remove":
            task = service.remove_task(args.task_id)
            print(f"Removed task #{task.id}: {task.title}")
            return 0

        if args.command == "summary":
            summary = service.summary(plan_date=validate_plan_date(args.date))
            print(f"Summary for {summary['plan_date']}:")
            print(f"- Total: {summary['total']}")
            print(f"- Pending: {summary['pending']}")
            print(f"- Completed: {summary['completed']}")
            print(
                "- High priority pending: "
                f"{summary['high_priority_pending']}"
            )
            return 0

        parser.print_help()
        return 1
    except (StoreError, TaskNotFoundError, ValueError) as exc:
        print(f"Error: {exc}", file=__import__("sys").stderr)
        return 1

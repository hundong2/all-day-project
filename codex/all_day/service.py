from __future__ import annotations

from dataclasses import replace
from datetime import date

from .models import Task, current_timestamp, default_plan_date
from .storage import TaskStore

VALID_PRIORITIES = {"low", "medium", "high"}
VALID_STATUSES = {"pending", "done", "all"}
MISSING = object()


class TaskNotFoundError(RuntimeError):
    """Raised when a task ID does not exist."""


class TaskService:
    def __init__(self, store: TaskStore):
        self.store = store

    def add_task(
        self,
        title: str,
        *,
        priority: str = "medium",
        tags: list[str] | None = None,
        note: str = "",
        plan_date: str | None = None,
    ) -> Task:
        normalized_title = title.strip()
        if not normalized_title:
            raise ValueError("Task title cannot be empty.")
        tasks, next_id = self.store.load_tasks()
        task = Task.create(
            next_id,
            title=normalized_title,
            priority=priority,
            tags=normalize_tags(tags or []),
            note=note.strip(),
            plan_date=plan_date or default_plan_date(),
        )
        tasks.append(task)
        self.store.save_tasks(tasks, next_id + 1)
        return task

    def list_tasks(
        self,
        *,
        plan_date: str | None = None,
        status: str = "pending",
        tag: str | None = None,
        query: str | None = None,
    ) -> list[Task]:
        tasks, _ = self.store.load_tasks()
        selected_date = plan_date or default_plan_date()
        normalized_tag = normalize_tag(tag) if tag else None
        normalized_query = query.strip().lower() if query else ""
        return [
            task
            for task in tasks
            if task.plan_date == selected_date
            and (status == "all" or task.status == status)
            and (normalized_tag is None or normalized_tag in task.tags)
            and task_matches_query(task, normalized_query)
        ]

    def mark_done(self, task_id: int) -> Task:
        tasks, next_id = self.store.load_tasks()
        updated: list[Task] = []
        matched: Task | None = None
        for task in tasks:
            if task.id == task_id:
                matched = replace(
                    task,
                    status="done",
                    completed_at=current_timestamp(),
                )
                updated.append(matched)
            else:
                updated.append(task)
        if matched is None:
            raise TaskNotFoundError(f"Task #{task_id} was not found.")
        self.store.save_tasks(updated, next_id)
        return matched

    def update_task(
        self,
        task_id: int,
        *,
        title: object = MISSING,
        priority: object = MISSING,
        tags: object = MISSING,
        note: object = MISSING,
        plan_date: object = MISSING,
    ) -> Task:
        tasks, next_id = self.store.load_tasks()
        updated: list[Task] = []
        matched: Task | None = None

        for task in tasks:
            if task.id != task_id:
                updated.append(task)
                continue

            next_title = task.title
            if title is not MISSING:
                normalized_title = str(title).strip()
                if not normalized_title:
                    raise ValueError("Task title cannot be empty.")
                next_title = normalized_title

            next_priority = task.priority
            if priority is not MISSING:
                next_priority = str(priority)

            next_tags = task.tags
            if tags is not MISSING:
                next_tags = normalize_tags(list(tags))

            next_note = task.note
            if note is not MISSING:
                next_note = str(note).strip()

            next_plan_date = task.plan_date
            if plan_date is not MISSING:
                next_plan_date = str(plan_date)

            matched = replace(
                task,
                title=next_title,
                priority=next_priority,
                tags=next_tags,
                note=next_note,
                plan_date=next_plan_date,
            )
            updated.append(matched)

        if matched is None:
            raise TaskNotFoundError(f"Task #{task_id} was not found.")

        self.store.save_tasks(updated, next_id)
        return matched

    def remove_task(self, task_id: int) -> Task:
        tasks, next_id = self.store.load_tasks()
        remaining: list[Task] = []
        removed: Task | None = None
        for task in tasks:
            if task.id == task_id:
                removed = task
            else:
                remaining.append(task)
        if removed is None:
            raise TaskNotFoundError(f"Task #{task_id} was not found.")
        self.store.save_tasks(remaining, next_id)
        return removed

    def summary(self, *, plan_date: str | None = None) -> dict[str, int | str]:
        selected_date = plan_date or default_plan_date()
        tasks = self.list_tasks(plan_date=selected_date, status="all")
        pending = [task for task in tasks if task.status == "pending"]
        completed = [task for task in tasks if task.status == "done"]
        high_priority_pending = [
            task for task in pending if task.priority == "high"
        ]
        return {
            "plan_date": selected_date,
            "total": len(tasks),
            "pending": len(pending),
            "completed": len(completed),
            "high_priority_pending": len(high_priority_pending),
        }


def validate_priority(priority: str) -> str:
    value = priority.lower()
    if value not in VALID_PRIORITIES:
        raise ValueError(
            f"Invalid priority '{priority}'. Choose from: "
            + ", ".join(sorted(VALID_PRIORITIES))
        )
    return value


def validate_status(status: str) -> str:
    value = status.lower()
    if value not in VALID_STATUSES:
        raise ValueError(
            f"Invalid status '{status}'. Choose from: "
            + ", ".join(sorted(VALID_STATUSES))
        )
    return value


def validate_plan_date(raw_value: str | None) -> str:
    if not raw_value:
        return default_plan_date()
    try:
        return date.fromisoformat(raw_value).isoformat()
    except ValueError as exc:
        raise ValueError(
            f"Invalid date '{raw_value}'. Use YYYY-MM-DD format."
        ) from exc


def normalize_tag(tag: str | None) -> str:
    normalized = (tag or "").strip().lower()
    if not normalized:
        raise ValueError("Tags must not be empty.")
    return normalized


def normalize_tags(tags: list[str]) -> list[str]:
    deduped: list[str] = []
    seen: set[str] = set()
    for raw_tag in tags:
        normalized = normalize_tag(raw_tag)
        if normalized not in seen:
            seen.add(normalized)
            deduped.append(normalized)
    return deduped


def task_matches_query(task: Task, query: str) -> bool:
    if not query:
        return True
    haystacks = [task.title, task.note, " ".join(task.tags)]
    lowered = [value.lower() for value in haystacks]
    return any(query in value for value in lowered)

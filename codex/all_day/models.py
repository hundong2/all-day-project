from __future__ import annotations

from dataclasses import asdict, dataclass
from datetime import date, datetime


def current_timestamp() -> str:
    return datetime.now().replace(microsecond=0).isoformat()


def default_plan_date() -> str:
    return date.today().isoformat()


@dataclass(slots=True)
class Task:
    id: int
    title: str
    priority: str
    plan_date: str
    tags: list[str]
    note: str = ""
    status: str = "pending"
    created_at: str = ""
    completed_at: str | None = None

    def to_dict(self) -> dict[str, object]:
        return asdict(self)

    @classmethod
    def create(
        cls,
        task_id: int,
        title: str,
        priority: str = "medium",
        tags: list[str] | None = None,
        note: str = "",
        plan_date: str | None = None,
    ) -> "Task":
        return cls(
            id=task_id,
            title=title,
            priority=priority,
            plan_date=plan_date or default_plan_date(),
            tags=list(tags or []),
            note=note,
            status="pending",
            created_at=current_timestamp(),
            completed_at=None,
        )

    @classmethod
    def from_dict(cls, payload: dict[str, object]) -> "Task":
        raw_tags = payload.get("tags", [])
        if not isinstance(raw_tags, list):
            raise TypeError("Task tags must be stored as a list.")
        return cls(
            id=int(payload["id"]),
            title=str(payload["title"]),
            priority=str(payload.get("priority", "medium")),
            plan_date=str(payload.get("plan_date", default_plan_date())),
            tags=[str(tag) for tag in raw_tags],
            note=str(payload.get("note", "")),
            status=str(payload.get("status", "pending")),
            created_at=str(payload.get("created_at", "")),
            completed_at=(
                str(payload["completed_at"])
                if payload.get("completed_at") is not None
                else None
            ),
        )

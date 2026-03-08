from __future__ import annotations

import json
import os
import tempfile
from pathlib import Path
from typing import Any

from .models import Task


class StoreError(RuntimeError):
    """Raised when the JSON store cannot be read or updated safely."""


def resolve_storage_path(explicit_path: str | None = None) -> Path:
    if explicit_path:
        return Path(explicit_path).expanduser()
    if os.getenv("ALL_DAY_STORAGE"):
        return Path(os.environ["ALL_DAY_STORAGE"]).expanduser()
    return Path.home() / ".all_day" / "tasks.json"


class TaskStore:
    def __init__(self, path: Path):
        self.path = path

    def load(self) -> dict[str, Any]:
        if not self.path.exists():
            return {"next_id": 1, "tasks": []}

        try:
            payload = json.loads(self.path.read_text(encoding="utf-8"))
        except json.JSONDecodeError as exc:
            raise StoreError(
                f"Storage file is corrupted: {self.path}"
            ) from exc
        except OSError as exc:
            raise StoreError(
                f"Storage file could not be read: {self.path}"
            ) from exc

        if not isinstance(payload, dict):
            raise StoreError(f"Storage file has invalid format: {self.path}")

        next_id = payload.get("next_id", 1)
        tasks = payload.get("tasks", [])

        if not isinstance(next_id, int) or next_id < 1:
            raise StoreError(f"Storage file has invalid next_id: {self.path}")
        if not isinstance(tasks, list):
            raise StoreError(f"Storage file has invalid tasks list: {self.path}")

        return {"next_id": next_id, "tasks": tasks}

    def save(self, state: dict[str, Any]) -> None:
        target = self.path.expanduser()
        if target.exists() and target.is_symlink():
            raise StoreError(f"Refusing to write to symlink: {target}")
        temporary_name: str | None = None
        try:
            target.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
            payload = json.dumps(state, indent=2, sort_keys=True) + "\n"
            descriptor, temporary_name = tempfile.mkstemp(dir=str(target.parent))
            os.fchmod(descriptor, 0o600)
        except OSError as exc:
            raise StoreError(
                f"Storage file could not be written: {target}"
            ) from exc

        try:
            os.fchmod(descriptor, 0o600)
            with os.fdopen(descriptor, "w", encoding="utf-8") as handle:
                handle.write(payload)
            os.replace(temporary_name, target)
        except OSError as exc:
            raise StoreError(
                f"Storage file could not be written: {target}"
            ) from exc
        finally:
            if temporary_name and os.path.exists(temporary_name):
                os.unlink(temporary_name)

    def load_tasks(self) -> tuple[list[Task], int]:
        state = self.load()
        try:
            tasks = [Task.from_dict(item) for item in state["tasks"]]
        except (KeyError, TypeError, ValueError) as exc:
            raise StoreError(
                f"Storage file has invalid task records: {self.path}"
            ) from exc
        return tasks, int(state["next_id"])

    def save_tasks(self, tasks: list[Task], next_id: int) -> None:
        self.save(
            {
                "next_id": next_id,
                "tasks": [task.to_dict() for task in tasks],
            }
        )

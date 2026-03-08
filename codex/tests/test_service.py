from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

from all_day.service import TaskNotFoundError, TaskService
from all_day.storage import TaskStore


class TaskServiceTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temp_dir = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp_dir.cleanup)
        self.storage_path = Path(self.temp_dir.name) / "tasks.json"
        self.service = TaskService(TaskStore(self.storage_path))

    def test_add_and_list_tasks_for_specific_date(self) -> None:
        self.service.add_task(
            "Write unit tests",
            priority="high",
            plan_date="2026-03-06",
        )
        self.service.add_task(
            "Plan tomorrow",
            priority="low",
            plan_date="2026-03-07",
        )

        tasks = self.service.list_tasks(
            plan_date="2026-03-06",
            status="pending",
        )

        self.assertEqual(len(tasks), 1)
        self.assertEqual(tasks[0].title, "Write unit tests")
        self.assertEqual(tasks[0].priority, "high")

    def test_mark_done_and_summary(self) -> None:
        task = self.service.add_task(
            "Finish report",
            priority="high",
            tags=["work"],
            plan_date="2026-03-06",
        )
        self.service.add_task(
            "Take a walk",
            priority="low",
            tags=["health"],
            plan_date="2026-03-06",
        )

        completed = self.service.mark_done(task.id)
        summary = self.service.summary(plan_date="2026-03-06")

        self.assertEqual(completed.status, "done")
        self.assertEqual(summary["total"], 2)
        self.assertEqual(summary["completed"], 1)
        self.assertEqual(summary["pending"], 1)
        self.assertEqual(summary["high_priority_pending"], 0)

    def test_remove_missing_task_raises(self) -> None:
        with self.assertRaises(TaskNotFoundError):
            self.service.remove_task(999)

    def test_blank_title_is_rejected(self) -> None:
        with self.assertRaises(ValueError):
            self.service.add_task("   ")

    def test_list_tasks_can_filter_by_tag_and_query(self) -> None:
        self.service.add_task(
            "Write launch draft",
            priority="high",
            tags=["work", "writing"],
            note="Focus on the intro",
            plan_date="2026-03-06",
        )
        self.service.add_task(
            "Stretch break",
            priority="low",
            tags=["health"],
            plan_date="2026-03-06",
        )

        tagged_tasks = self.service.list_tasks(
            plan_date="2026-03-06",
            status="pending",
            tag="work",
        )
        queried_tasks = self.service.list_tasks(
            plan_date="2026-03-06",
            status="pending",
            query="intro",
        )

        self.assertEqual(len(tagged_tasks), 1)
        self.assertEqual(tagged_tasks[0].tags, ["work", "writing"])
        self.assertEqual(len(queried_tasks), 1)
        self.assertEqual(queried_tasks[0].title, "Write launch draft")

    def test_update_task_changes_core_fields(self) -> None:
        task = self.service.add_task(
            "Draft report",
            priority="medium",
            tags=["work"],
            note="first pass",
            plan_date="2026-03-07",
        )

        updated = self.service.update_task(
            task.id,
            title="Draft final report",
            priority="high",
            tags=["work", "review"],
            note="",
            plan_date="2026-03-08",
        )

        self.assertEqual(updated.title, "Draft final report")
        self.assertEqual(updated.priority, "high")
        self.assertEqual(updated.tags, ["work", "review"])
        self.assertEqual(updated.note, "")
        self.assertEqual(updated.plan_date, "2026-03-08")

from __future__ import annotations

import io
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from pathlib import Path

from all_day.cli import main


class CliTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temp_dir = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp_dir.cleanup)
        self.storage_path = Path(self.temp_dir.name) / "tasks.json"

    def run_cli(self, *args: str) -> tuple[int, str, str]:
        stdout = io.StringIO()
        stderr = io.StringIO()
        with redirect_stdout(stdout), redirect_stderr(stderr):
            code = main(["--storage", str(self.storage_path), *args])
        return code, stdout.getvalue(), stderr.getvalue()

    def test_add_list_done_remove_flow(self) -> None:
        code, stdout, stderr = self.run_cli(
            "add",
            "Write docs",
            "--priority",
            "high",
            "--tag",
            "work",
            "--tag",
            "writing",
            "--date",
            "2026-03-06",
        )
        self.assertEqual(code, 0)
        self.assertIn("Added task #1", stdout)
        self.assertEqual(stderr, "")

        code, stdout, _ = self.run_cli(
            "list",
            "--status",
            "pending",
            "--date",
            "2026-03-06",
        )
        self.assertEqual(code, 0)
        self.assertIn("Tasks for 2026-03-06", stdout)
        self.assertIn("Write docs", stdout)
        self.assertIn("tags: work, writing", stdout)

        code, stdout, _ = self.run_cli("done", "1")
        self.assertEqual(code, 0)
        self.assertIn("Completed task #1", stdout)

        code, stdout, _ = self.run_cli("summary", "--date", "2026-03-06")
        self.assertEqual(code, 0)
        self.assertIn("- Completed: 1", stdout)
        self.assertIn("- Pending: 0", stdout)

        code, stdout, _ = self.run_cli("remove", "1")
        self.assertEqual(code, 0)
        self.assertIn("Removed task #1", stdout)

    def test_list_can_filter_by_tag_and_query(self) -> None:
        code, _, _ = self.run_cli(
            "add",
            "Draft outline",
            "--tag",
            "work",
            "--note",
            "launch memo",
            "--date",
            "2026-03-06",
        )
        self.assertEqual(code, 0)
        code, _, _ = self.run_cli(
            "add",
            "Take a walk",
            "--tag",
            "health",
            "--date",
            "2026-03-06",
        )
        self.assertEqual(code, 0)

        code, stdout, stderr = self.run_cli(
            "list",
            "--status",
            "pending",
            "--date",
            "2026-03-06",
            "--tag",
            "work",
            "--query",
            "memo",
        )

        self.assertEqual(code, 0)
        self.assertEqual(stderr, "")
        self.assertIn("Draft outline", stdout)
        self.assertNotIn("Take a walk", stdout)

    def test_update_can_replace_fields_and_clear_tags(self) -> None:
        code, _, _ = self.run_cli(
            "add",
            "Draft outline",
            "--tag",
            "work",
            "--note",
            "rough version",
            "--date",
            "2026-03-06",
        )
        self.assertEqual(code, 0)

        code, stdout, stderr = self.run_cli(
            "update",
            "1",
            "--title",
            "Final outline",
            "--priority",
            "high",
            "--note",
            "",
            "--clear-tags",
            "--date",
            "2026-03-07",
        )
        self.assertEqual(code, 0)
        self.assertEqual(stderr, "")
        self.assertIn("Updated task #1", stdout)

        code, stdout, stderr = self.run_cli(
            "list",
            "--status",
            "pending",
            "--date",
            "2026-03-07",
        )
        self.assertEqual(code, 0)
        self.assertEqual(stderr, "")
        self.assertIn("Final outline", stdout)
        self.assertIn("(high)", stdout)
        self.assertNotIn("tags:", stdout)
        self.assertNotIn("note:", stdout)

    def test_invalid_date_returns_error(self) -> None:
        code, stdout, stderr = self.run_cli(
            "add",
            "Broken date",
            "--date",
            "03-06-2026",
        )
        self.assertEqual(code, 1)
        self.assertEqual(stdout, "")
        self.assertIn("Invalid date", stderr)

    def test_corrupted_store_returns_error(self) -> None:
        self.storage_path.write_text("{not json}", encoding="utf-8")

        code, stdout, stderr = self.run_cli("list")

        self.assertEqual(code, 1)
        self.assertEqual(stdout, "")
        self.assertIn("Storage file is corrupted", stderr)

    def test_directory_storage_path_is_rejected_cleanly(self) -> None:
        code, stdout, stderr = self.run_cli("list")
        self.assertEqual(code, 0)
        self.assertEqual(stderr, "")
        self.assertIn("No tasks found", stdout)

        stdout = io.StringIO()
        stderr = io.StringIO()
        with redirect_stdout(stdout), redirect_stderr(stderr):
            code = main(["--storage", self.temp_dir.name, "list"])

        self.assertEqual(code, 1)
        self.assertEqual(stdout.getvalue(), "")
        self.assertIn("Storage file could not be read", stderr.getvalue())

    def test_invalid_task_schema_returns_error(self) -> None:
        self.storage_path.write_text(
            '{"next_id": 2, "tasks": [{"title": "Missing id"}]}',
            encoding="utf-8",
        )

        code, stdout, stderr = self.run_cli("list", "--status", "all")

        self.assertEqual(code, 1)
        self.assertEqual(stdout, "")
        self.assertIn("Storage file has invalid task records", stderr)

    def test_invalid_tags_schema_returns_error(self) -> None:
        self.storage_path.write_text(
            '{"next_id": 2, "tasks": [{"id": 1, "title": "Oops", "tags": "work"}]}',
            encoding="utf-8",
        )

        code, stdout, stderr = self.run_cli("list", "--status", "all")

        self.assertEqual(code, 1)
        self.assertEqual(stdout, "")
        self.assertIn("Storage file has invalid task records", stderr)

    def test_symlink_storage_path_is_rejected(self) -> None:
        real_path = Path(self.temp_dir.name) / "real.json"
        real_path.write_text("{}", encoding="utf-8")
        symlink_path = Path(self.temp_dir.name) / "link.json"
        symlink_path.symlink_to(real_path)

        stdout = io.StringIO()
        stderr = io.StringIO()
        with redirect_stdout(stdout), redirect_stderr(stderr):
            code = main(
                [
                    "--storage",
                    str(symlink_path),
                    "add",
                    "Unsafe path",
                ]
            )

        self.assertEqual(code, 1)
        self.assertEqual(stdout.getvalue(), "")
        self.assertIn("Refusing to write to symlink", stderr.getvalue())

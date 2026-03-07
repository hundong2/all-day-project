import os
import subprocess
from pathlib import Path

class GitManager:
    def __init__(self, repo_path="."):
        self.repo_path = Path(repo_path).absolute()

    def run_command(self, cmd, cwd=None):
        work_dir = cwd if cwd else self.repo_path
        print(f"Running: {' '.join(cmd)} in {work_dir}")
        result = subprocess.run(cmd, cwd=work_dir, capture_output=True, text=True)
        if result.returncode != 0:
            print(f"Error executing {' '.join(cmd)}: {result.stderr}")
            return False, result.stderr
        return True, result.stdout

    def create_worktree(self, branch_name, worktree_path):
        """Creates a new branch and worktree."""
        # Ensure latest main
        self.run_command(['git', 'fetch', 'origin'])
        self.run_command(['git', 'checkout', 'main'])
        self.run_command(['git', 'pull', 'origin', 'main'])
        
        wt_path = self.repo_path.parent / worktree_path
        if wt_path.exists():
            print(f"Worktree path {wt_path} already exists.")
            return False, str(wt_path)
            
        success, out = self.run_command(['git', 'worktree', 'add', '-b', branch_name, str(wt_path)])
        return success, str(wt_path)

    def commit_and_push(self, worktree_path, branch_name, message):
        """Commits changes in worktree and pushes to remote."""
        success, _ = self.run_command(['git', 'add', '.'], cwd=worktree_path)
        if not success: return False
        
        success, out = self.run_command(['git', 'commit', '-m', message], cwd=worktree_path)
        if not success and "nothing to commit" not in out:
            return False
            
        success, _ = self.run_command(['git', 'push', '-u', 'origin', branch_name], cwd=worktree_path)
        return success

    def remove_worktree(self, worktree_path):
        """Removes the worktree and cleans up."""
        success, _ = self.run_command(['git', 'worktree', 'remove', '-f', str(worktree_path)])
        return success

from github import Github
from config import GITHUB_TOKEN, REPO_NAME

class GitHubManager:
    def __init__(self):
        if not GITHUB_TOKEN or not REPO_NAME:
            raise ValueError("GITHUB_TOKEN and REPO_NAME must be set in .env")
        self.gh = Github(GITHUB_TOKEN)
        self.repo = self.gh.get_repo(REPO_NAME)

    def get_file_content(self, filepath, branch="main"):
        try:
            file_content = self.repo.get_contents(filepath, ref=branch)
            return file_content.decoded_content.decode('utf-8')
        except Exception as e:
            return None

    def create_issue(self, title, body):
        issue = self.repo.create_issue(title=title, body=body)
        print(f"Created issue #{issue.number}: {title}")
        return issue

    def create_issue_comment(self, issue_number, body):
        issue = self.repo.get_issue(number=issue_number)
        comment = issue.create_comment(body)
        print(f"Added comment to issue #{issue_number}")
        return comment

    def get_issue_comments(self, issue_number):
        issue = self.repo.get_issue(number=issue_number)
        return [{"user": c.user.login, "body": c.body} for c in issue.get_comments()]

    def update_file(self, filepath, message, content, branch="main"):
        try:
            file = self.repo.get_contents(filepath, ref=branch)
            self.repo.update_file(filepath, message, content, file.sha, branch=branch)
            print(f"Updated {filepath} in branch {branch}")
        except Exception as e:
            # File doesn't exist, create it
            self.repo.create_file(filepath, message, content, branch=branch)
            print(f"Created {filepath} in branch {branch}")

    def create_pull_request(self, title, body, head, base="main"):
        try:
            pr = self.repo.create_pull(title=title, body=body, head=head, base=base)
            print(f"Created PR #{pr.number}: {title}")
            return pr
        except Exception as e:
            print(f"Failed to create PR: {e}")
            return None

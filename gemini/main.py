import os
import time
from apscheduler.schedulers.background import BackgroundScheduler
from github_manager import GitHubManager
from git_manager import GitManager
from agent_manager import AgentManager
from config import POLLING_INTERVAL_SECONDS

def check_for_new_plans():
    try:
        gh_mgr = GitHubManager()
        agent_mgr = AgentManager()
        git_mgr = GitManager()
        
        # In a real scenario, this would check commit history or webhooks for plan.md updates
        # For simplicity, we check if final_plan.md exists
        plan_content = gh_mgr.get_file_content("gemini/plan.md")
        final_plan_content = gh_mgr.get_file_content("gemini/final_plan.md")

        if plan_content and not final_plan_content:
            print("Detected new plan.md without final_plan.md. Initiating Planning Phase...")
            
            # 1. Create Discussion Issue
            issue = gh_mgr.create_issue("Planning Discussion: AI Architecture", f"Initial Plan:\n\n{plan_content}")
            
            # 2. PM Agent analyzes and discusses
            print("PM Agent is analyzing the plan...")
            final_plan_output = agent_mgr.plan_project(plan_content)
            
            # 3. Post summary and update final_plan.md
            gh_mgr.create_issue_comment(issue.number, f"**Final Plan Summary**\n\n{final_plan_output}")
            gh_mgr.update_file("gemini/final_plan.md", "docs: update final plan from PM Agent", str(final_plan_output))
            print("final_plan.md created successfully!")

        elif final_plan_content:
            print("final_plan.md exists. Checking for open task issues...")
            # For demonstration, this is where git_worktree logic and issue allocation runs.
            # E.g. creating issues based on final_plan.md sections, parsing them, assigning to dev_agent.
            print("Task assignment and execution logic goes here.")
            
    except Exception as e:
        print(f"Error during polling: {e}")

if __name__ == "__main__":
    print("Starting AI PM Scheduler...")
    
    # Run once at startup
    check_for_new_plans()

    # Schedule periodic checks
    scheduler = BackgroundScheduler()
    scheduler.add_job(check_for_new_plans, 'interval', seconds=POLLING_INTERVAL_SECONDS)
    scheduler.start()

    try:
        # Keep the main thread alive
        while True:
            time.sleep(1)
    except (KeyboardInterrupt, SystemExit):
        scheduler.shutdown()
        print("Scheduler shut down successfully.")

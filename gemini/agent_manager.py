import os
from langchain_google_genai import ChatGoogleGenerativeAI
from crewai import Agent, Task, Crew, Process
from config import GEMINI_API_KEY

class AgentManager:
    def __init__(self):
        # We initialize Gemini as the PM agent
        self.llm = ChatGoogleGenerativeAI(
            model="gemini-1.5-pro",
            verbose=True,
            temperature=0.7,
            google_api_key=GEMINI_API_KEY
        )

        self.pm_agent = Agent(
            role='Project Manager',
            goal='Analyze plan.md, initiate discussions, and create final_plan.md. Allocate tasks to other agents.',
            backstory='An expert AI PM who excels at software architecture and delegating tasks to developer agents.',
            verbose=True,
            allow_delegation=False,
            llm=self.llm
        )
        
        # In a full implementation, you'd add claude for review, codex for code.
        self.dev_agent = Agent(
            role='Senior Developer',
            goal='Write clean, efficient, and well-documented code based on assigned issues.',
            backstory='A 10x developer AI capable of building any software component.',
            verbose=True,
            allow_delegation=False,
            llm=self.llm  # Using Gemini as fallback for now
        )

        self.reviewer_agent = Agent(
            role='Code Reviewer',
            goal='Review code for bugs, style guidelines, and performance issues.',
            backstory='A meticulous senior engineer who ensures no bad code enters the main branch.',
            verbose=True,
            allow_delegation=False,
            llm=self.llm
        )

    def plan_project(self, plan_content):
        # We start a discussion and summarize
        task = Task(
            description=f"Analyze the following project plan and create a detailed sequence of actionable implementation steps. Plan: \n\n{plan_content}",
            agent=self.pm_agent,
            expected_output="A markdown formatted final plan outlining the sub-tasks to be created as GitHub issues."
        )

        crew = Crew(
            agents=[self.pm_agent, self.dev_agent, self.reviewer_agent],
            tasks=[task],
            verbose=True,
            process=Process.sequential
        )

        result = crew.kickoff()
        return result

    def execute_task(self, issue_title, issue_body):
        task = Task(
            description=f"Implement the requested feature. Title: {issue_title}\n\nDetails: {issue_body}",
            agent=self.dev_agent,
            expected_output="The source code to resolve the issue with an explanation of changes."
        )

        crew = Crew(
            agents=[self.dev_agent],
            tasks=[task],
            verbose=True
        )

        result = crew.kickoff()
        return result

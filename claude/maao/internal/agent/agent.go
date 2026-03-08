package agent

import (
	"context"
	"time"
)

type TaskType string

const (
	TaskArchitecture        TaskType = "architecture"
	TaskRefactoring         TaskType = "refactoring"
	TaskComplexLogic        TaskType = "complex_logic"
	TaskPlanning            TaskType = "planning"
	TaskDocumentation       TaskType = "documentation"
	TaskSpecWriting         TaskType = "spec_writing"
	TaskUIDesign            TaskType = "ui_design"
	TaskResearch            TaskType = "research"
	TaskDebugging           TaskType = "debugging"
	TaskTesting             TaskType = "testing"
	TaskCICD                TaskType = "ci_cd"
	TaskShellScripting      TaskType = "shell_scripting"
	TaskInfrastructure      TaskType = "infrastructure"
	TaskCodeReview          TaskType = "code_review"
	TaskGitHubIntegration   TaskType = "github_integration"
	TaskCodebaseExploration TaskType = "codebase_exploration"
	TaskPRManagement        TaskType = "pr_management"
	TaskDependencyUpgrade   TaskType = "dependency_upgrade"
	TaskGeneral             TaskType = "general"
)

type Agent interface {
	Name() string
	Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error)
	IsAvailable() bool
	Specialties() []TaskType
	TokenBudget() TokenBudget
	UpdateTokenUsage(usage TokenUsage)
}

type ExecuteRequest struct {
	WorkDir   string
	Prompt    string
	Context   []string
	SessionID string
	Timeout   time.Duration
}

type ExecuteResponse struct {
	Output    string
	Tokens    TokenUsage
	ExitCode  int
	Duration  time.Duration
	SessionID string
}

type TokenUsage struct {
	Input  int
	Output int
	Total  int
	Cached int
}

type TokenBudget struct {
	DailyLimit int
	UsedToday  int
	Remaining  int
}

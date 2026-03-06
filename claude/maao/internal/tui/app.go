package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	program *tea.Program
	model   *dashboardModel
}

func New() *App {
	m := newDashboardModel()
	return &App{
		model: &m,
	}
}

func (a *App) Run() error {
	a.program = tea.NewProgram(*a.model, tea.WithAltScreen())
	_, err := a.program.Run()
	return err
}

func (a *App) UpdateAgentStatus(name string, status AgentStatus) {
	if a.program != nil {
		a.program.Send(AgentStatusMsg{Name: name, Status: status})
	}
}

func (a *App) AddLogEntry(entry LogEntry) {
	if a.program != nil {
		a.program.Send(LogEntryMsg{Entry: entry})
	}
}

func (a *App) SetWorkflow(repo, phase string) {
	if a.program != nil {
		a.program.Send(setWorkflowMsg{repo: repo, phase: phase})
	}
}

func (a *App) SetTokenBudgets(budgets []TokenBudget) {
	if a.program != nil {
		a.program.Send(setTokenBudgetsMsg{budgets: budgets})
	}
}

type setWorkflowMsg struct {
	repo  string
	phase string
}

type setTokenBudgetsMsg struct {
	budgets []TokenBudget
}


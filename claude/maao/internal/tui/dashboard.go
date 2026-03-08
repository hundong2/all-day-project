package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type AgentStatusMsg struct {
	Name   string
	Status AgentStatus
}

type LogEntryMsg struct {
	Entry LogEntry
}

type dashboardModel struct {
	width     int
	height    int
	repo      string
	phase     string
	startTime time.Time
	agents    []AgentStatus
	tokens    []TokenBudget
	logs      *logViewer
	quitting  bool
	paused    bool
}

func newDashboardModel() dashboardModel {
	return dashboardModel{
		width:     80,
		height:    24,
		repo:      "",
		phase:     "idle",
		startTime: time.Now(),
		agents:    nil,
		tokens:    nil,
		logs:      newLogViewer(8),
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m dashboardModel) Init() tea.Cmd {
	return tickCmd()
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "p":
			m.paused = !m.paused
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		return m, tickCmd()
	case AgentStatusMsg:
		m.updateAgent(msg.Name, msg.Status)
	case LogEntryMsg:
		m.logs.add(msg.Entry)
	case setWorkflowMsg:
		m.repo = msg.repo
		m.phase = msg.phase
	case setTokenBudgetsMsg:
		m.tokens = msg.budgets
	}
	return m, nil
}

func (m *dashboardModel) updateAgent(name string, status AgentStatus) {
	for i, a := range m.agents {
		if a.Name == name {
			m.agents[i] = status
			return
		}
	}
	m.agents = append(m.agents, status)
}

func (m dashboardModel) View() string {
	if m.quitting {
		return "Shutting down...\n"
	}

	w := m.width
	if w < 40 {
		w = 40
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("4")).
		Width(w).
		Padding(0, 1)

	elapsed := time.Since(m.startTime).Truncate(time.Second)
	pauseIndicator := ""
	if m.paused {
		pauseIndicator = " [PAUSED]"
	}
	header := fmt.Sprintf("MAAO Dashboard | Workflow: %s | Phase: %s | Elapsed: %s%s",
		m.repo, m.phase, elapsed, pauseIndicator)
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n\n")

	// Agents panel
	b.WriteString(renderAgentPanel(m.agents, w))
	b.WriteString("\n\n")

	// Token budget panel
	b.WriteString(renderTokenPanel(m.tokens, w))
	b.WriteString("\n\n")

	// Activity log
	b.WriteString(m.logs.render(w))
	b.WriteString("\n\n")

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))
	b.WriteString(footerStyle.Render("[q] Quit  [p] Pause/Resume"))
	b.WriteString("\n")

	return b.String()
}

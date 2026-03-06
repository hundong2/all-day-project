package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type AgentStatusType int

const (
	StatusIdle AgentStatusType = iota
	StatusRunning
	StatusWaiting
	StatusDone
	StatusError
)

type AgentStatus struct {
	Name       string
	Status     AgentStatusType
	Progress   int // 0-100
	Issue      string
	Specialty  string
	QueuedItem string
	StartedAt  time.Time
}

func (s AgentStatusType) Icon() string {
	switch s {
	case StatusRunning:
		return "[RUN]"
	case StatusWaiting:
		return "[WAIT]"
	case StatusDone:
		return "[DONE]"
	case StatusError:
		return "[ERR]"
	default:
		return "[IDLE]"
	}
}

func (s AgentStatusType) Style() lipgloss.Style {
	switch s {
	case StatusRunning:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
	case StatusWaiting:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow
	case StatusDone:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // blue
	case StatusError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // red
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // gray
	}
}

func progressBar(percent int, width int) string {
	filled := width * percent / 100
	empty := width - filled
	return "[" + strings.Repeat("#", filled) + strings.Repeat(".", empty) + "]"
}

func renderAgentRow(a AgentStatus) string {
	icon := a.Status.Style().Render(a.Status.Icon())
	bar := progressBar(a.Progress, 10)
	detail := a.Issue
	if a.Status == StatusWaiting && a.QueuedItem != "" {
		detail = fmt.Sprintf("Waiting   (queued: %s)", a.QueuedItem)
	}
	return fmt.Sprintf("  %s %-8s %s %3d%%  %-10s %-20s",
		icon, a.Name, bar, a.Progress, detail, a.Specialty)
}

func renderAgentPanel(agents []AgentStatus, width int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(width - 4).
		Padding(0, 1)

	var rows []string
	for _, a := range agents {
		rows = append(rows, renderAgentRow(a))
	}

	title := lipgloss.NewStyle().Bold(true).Render(" Agents ")
	content := strings.Join(rows, "\n")
	return title + "\n" + border.Render(content)
}

type TokenBudget struct {
	Agent     string
	Used      int
	Total     int
	Remaining int
}

func renderTokenPanel(budgets []TokenBudget, width int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(width - 4).
		Padding(0, 1)

	var rows []string
	for _, b := range budgets {
		pct := 0
		if b.Total > 0 {
			pct = b.Remaining * 10 / b.Total
		}
		bar := strings.Repeat("#", pct) + strings.Repeat(".", 10-pct)
		rows = append(rows, fmt.Sprintf("  %-8s %s  %s / %s remaining",
			b.Agent+":", bar, formatTokens(b.Remaining), formatTokens(b.Total)))
	}

	title := lipgloss.NewStyle().Bold(true).Render(" Token Budget (Today) ")
	content := strings.Join(rows, "\n")
	return title + "\n" + border.Render(content)
}

func formatTokens(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.0fM", float64(n)/1000000)
	}
	return fmt.Sprintf("%dK", n/1000)
}

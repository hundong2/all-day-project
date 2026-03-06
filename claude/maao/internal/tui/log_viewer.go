package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogWarn
	LogError
	LogDebug
)

type LogEntry struct {
	Timestamp time.Time
	Agent     string
	Message   string
	Level     LogLevel
}

const maxLogEntries = 100

type logViewer struct {
	entries []LogEntry
	maxShow int
}

func newLogViewer(maxShow int) *logViewer {
	return &logViewer{
		maxShow: maxShow,
	}
}

func (lv *logViewer) add(entry LogEntry) {
	lv.entries = append(lv.entries, entry)
	if len(lv.entries) > maxLogEntries {
		lv.entries = lv.entries[len(lv.entries)-maxLogEntries:]
	}
}

func (lv *logViewer) render(width int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(width - 4).
		Padding(0, 1)

	start := 0
	if len(lv.entries) > lv.maxShow {
		start = len(lv.entries) - lv.maxShow
	}

	var rows []string
	for _, e := range lv.entries[start:] {
		ts := e.Timestamp.Format("15:04:05")
		agentStyle := lipgloss.NewStyle().Foreground(agentColor(e.Agent))
		agent := agentStyle.Render(fmt.Sprintf("[%s]", e.Agent))
		rows = append(rows, fmt.Sprintf("  %s  %s  %s", ts, agent, e.Message))
	}

	if len(rows) == 0 {
		rows = append(rows, "  No activity yet.")
	}

	title := lipgloss.NewStyle().Bold(true).Render(" Activity Log ")
	content := strings.Join(rows, "\n")
	return title + "\n" + border.Render(content)
}

func agentColor(name string) lipgloss.Color {
	switch name {
	case "claude":
		return lipgloss.Color("208") // orange
	case "gemini":
		return lipgloss.Color("12") // blue
	case "codex":
		return lipgloss.Color("10") // green
	case "copilot":
		return lipgloss.Color("13") // magenta
	default:
		return lipgloss.Color("7") // white
	}
}

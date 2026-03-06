package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

type CodexEvent struct {
	Type    string        `json:"type"`
	Message *CodexMessage `json:"message,omitempty"`
	Usage   *CodexUsage   `json:"usage,omitempty"`
}

type CodexMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CodexUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

func ParseCodexJSONL(data []byte) (string, TokenUsage, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var outputs []string
	var tokens TokenUsage

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event CodexEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		switch event.Type {
		case "message":
			if event.Message != nil && event.Message.Role == "assistant" && event.Message.Content != "" {
				outputs = append(outputs, event.Message.Content)
			}
		case "usage":
			if event.Usage != nil {
				tokens.Input += event.Usage.InputTokens
				tokens.Output += event.Usage.OutputTokens
				tokens.Total += event.Usage.TotalTokens
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", TokenUsage{}, fmt.Errorf("parse codex jsonl: %w", err)
	}

	return strings.Join(outputs, "\n"), tokens, nil
}

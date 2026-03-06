package agent

import (
	"context"
	"fmt"

	"github.com/maao/internal/agent/parser"
)

type ClaudeAgent struct {
	BaseAgent
}

func NewClaudeAgent() *ClaudeAgent {
	return &ClaudeAgent{
		BaseAgent: NewBaseAgent("claude", "claude", []TaskType{
			TaskArchitecture,
			TaskRefactoring,
			TaskComplexLogic,
			TaskDebugging,
			TaskCodeReview,
			TaskCodebaseExploration,
		}, 1000000),
	}
}

func (c *ClaudeAgent) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	args := []string{
		"-p", req.Prompt,
		"--output-format", "json",
		"--allowedTools", "Read,Write,Edit,Bash",
		"--dangerously-skip-permissions",
		"--max-turns", "30",
	}

	if req.SessionID != "" {
		args = append(args, "--session-id", req.SessionID)
	}

	stdout, stderr, exitCode, duration, err := c.RunCommand(ctx, req, args)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 && len(stdout) == 0 {
		return &ExecuteResponse{
			Output:   string(stderr),
			ExitCode: exitCode,
			Duration: duration,
		}, fmt.Errorf("claude exited with code %d: %s", exitCode, string(stderr))
	}

	output, ptokens, parseErr := parser.ParseClaudeJSON(stdout)
	if parseErr != nil {
		return &ExecuteResponse{
			Output:   string(stdout),
			ExitCode: exitCode,
			Duration: duration,
		}, nil
	}

	tokens := convertTokenUsage(ptokens)
	c.UpdateTokenUsage(tokens)

	return &ExecuteResponse{
		Output:    output,
		Tokens:    tokens,
		ExitCode:  exitCode,
		Duration:  duration,
		SessionID: req.SessionID,
	}, nil
}

func init() {
	RegisterFactory("claude", func() (Agent, error) {
		return NewClaudeAgent(), nil
	})
}

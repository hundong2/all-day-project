package agent

import (
	"context"
	"fmt"

	"github.com/maao/internal/agent/parser"
)

type CodexAgent struct {
	BaseAgent
	model string
}

func NewCodexAgent() *CodexAgent {
	return &CodexAgent{
		BaseAgent: NewBaseAgent("codex", "codex", []TaskType{
			TaskDebugging,
			TaskTesting,
			TaskShellScripting,
			TaskDependencyUpgrade,
		}, 1500000),
		model: "gpt-5.3-codex",
	}
}

func (c *CodexAgent) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	args := []string{
		"exec", req.Prompt,
		"--yolo",
		"--jsonl",
		"-m", c.model,
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
		}, fmt.Errorf("codex exited with code %d: %s", exitCode, string(stderr))
	}

	output, ptokens, parseErr := parser.ParseCodexJSONL(stdout)
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
	RegisterFactory("codex", func() (Agent, error) {
		return NewCodexAgent(), nil
	})
}

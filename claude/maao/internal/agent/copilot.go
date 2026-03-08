package agent

import (
	"context"
	"fmt"

	"github.com/maao/internal/agent/parser"
)

type CopilotAgent struct {
	BaseAgent
}

func NewCopilotAgent() *CopilotAgent {
	return &CopilotAgent{
		BaseAgent: NewBaseAgent("copilot", "copilot", []TaskType{
			TaskCICD,
			TaskGitHubIntegration,
			TaskPRManagement,
			TaskInfrastructure,
		}, 1000000),
	}
}

func (c *CopilotAgent) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	args := []string{
		"--allow-all-tools",
		"-p", req.Prompt,
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
		}, fmt.Errorf("copilot exited with code %d: %s", exitCode, string(stderr))
	}

	output, ptokens := parser.ParseCopilotText(stdout)

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
	RegisterFactory("copilot", func() (Agent, error) {
		return NewCopilotAgent(), nil
	})
}

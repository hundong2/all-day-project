package agent

import (
	"context"
	"fmt"

	"github.com/maao/internal/agent/parser"
)

type GeminiAgent struct {
	BaseAgent
}

func NewGeminiAgent() *GeminiAgent {
	return &GeminiAgent{
		BaseAgent: NewBaseAgent("gemini", "gemini", []TaskType{
			TaskPlanning,
			TaskDocumentation,
			TaskSpecWriting,
			TaskResearch,
			TaskUIDesign,
		}, 2000000),
	}
}

func (g *GeminiAgent) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	args := []string{
		"-p", req.Prompt,
		"--output-format", "json",
	}

	stdout, stderr, exitCode, duration, err := g.RunCommand(ctx, req, args)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 && len(stdout) == 0 {
		return &ExecuteResponse{
			Output:   string(stderr),
			ExitCode: exitCode,
			Duration: duration,
		}, fmt.Errorf("gemini exited with code %d: %s", exitCode, string(stderr))
	}

	output, ptokens, parseErr := parser.ParseGeminiJSON(stdout)
	if parseErr != nil {
		return &ExecuteResponse{
			Output:   string(stdout),
			ExitCode: exitCode,
			Duration: duration,
		}, nil
	}

	tokens := convertTokenUsage(ptokens)
	g.UpdateTokenUsage(tokens)

	return &ExecuteResponse{
		Output:    output,
		Tokens:    tokens,
		ExitCode:  exitCode,
		Duration:  duration,
		SessionID: req.SessionID,
	}, nil
}

func init() {
	RegisterFactory("gemini", func() (Agent, error) {
		return NewGeminiAgent(), nil
	})
}

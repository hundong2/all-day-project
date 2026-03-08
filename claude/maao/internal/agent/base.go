package agent

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/maao/internal/agent/parser"
)

type CommandRunner func(ctx context.Context, name string, args ...string) *exec.Cmd

var DefaultCommandRunner CommandRunner = exec.CommandContext

type BaseAgent struct {
	name        string
	binary      string
	specialties []TaskType
	budget      TokenBudget
	mu          sync.Mutex
	cmdRunner   CommandRunner
}

func NewBaseAgent(name, binary string, specialties []TaskType, dailyLimit int) BaseAgent {
	return BaseAgent{
		name:        name,
		binary:      binary,
		specialties: specialties,
		budget: TokenBudget{
			DailyLimit: dailyLimit,
			Remaining:  dailyLimit,
		},
		cmdRunner: DefaultCommandRunner,
	}
}

func (b *BaseAgent) Name() string {
	return b.name
}

func (b *BaseAgent) IsAvailable() bool {
	_, err := exec.LookPath(b.binary)
	return err == nil
}

func (b *BaseAgent) Specialties() []TaskType {
	return b.specialties
}

func (b *BaseAgent) TokenBudget() TokenBudget {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.budget
}

func (b *BaseAgent) UpdateTokenUsage(usage TokenUsage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.budget.UsedToday += usage.Total
	b.budget.Remaining = b.budget.DailyLimit - b.budget.UsedToday
	if b.budget.Remaining < 0 {
		b.budget.Remaining = 0
	}
}

func (b *BaseAgent) SetCommandRunner(runner CommandRunner) {
	b.cmdRunner = runner
}

func (b *BaseAgent) RunCommand(ctx context.Context, req ExecuteRequest, args []string) ([]byte, []byte, int, time.Duration, error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	runner := b.cmdRunner
	if runner == nil {
		runner = DefaultCommandRunner
	}

	cmd := runner(ctx, b.binary, args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else if ctx.Err() == context.DeadlineExceeded {
			return nil, nil, -1, duration, fmt.Errorf("agent %s timed out after %s", b.name, timeout)
		} else {
			return nil, nil, -1, duration, fmt.Errorf("agent %s execution failed: %w", b.name, err)
		}
	}

	return stdout.Bytes(), stderr.Bytes(), exitCode, duration, nil
}

func convertTokenUsage(p parser.TokenUsage) TokenUsage {
	return TokenUsage{
		Input:  p.Input,
		Output: p.Output,
		Total:  p.Total,
		Cached: p.Cached,
	}
}

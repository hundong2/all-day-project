package agent

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func mockCommandRunner(stdout, stderr string, exitCode int) CommandRunner {
	return func(ctx context.Context, name string, args ...string) *exec.Cmd {
		// Use echo to simulate output; we override Stdout/Stderr in RunCommand,
		// so we use a helper script approach.
		cmd := exec.CommandContext(ctx, "sh", "-c",
			`printf '%s' "$MOCK_STDOUT"; printf '%s' "$MOCK_STDERR" >&2; exit `+itoa(exitCode))
		cmd.Env = append(cmd.Environ(),
			"MOCK_STDOUT="+stdout,
			"MOCK_STDERR="+stderr,
		)
		return cmd
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

func TestClaudeAgent_Execute(t *testing.T) {
	claudeJSON := `{
		"type": "result",
		"role": "assistant",
		"model": "claude-sonnet-4-6",
		"content": [{"type": "text", "text": "Hello from Claude"}],
		"stop_reason": "end_turn",
		"usage": {"input_tokens": 100, "output_tokens": 50, "cache_creation_input_tokens": 0, "cache_read_input_tokens": 10}
	}`

	agent := NewClaudeAgent()
	agent.SetCommandRunner(mockCommandRunner(claudeJSON, "", 0))

	resp, err := agent.Execute(context.Background(), ExecuteRequest{
		Prompt:  "test prompt",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Output != "Hello from Claude" {
		t.Errorf("expected 'Hello from Claude', got %q", resp.Output)
	}
	if resp.Tokens.Input != 100 {
		t.Errorf("expected input tokens 100, got %d", resp.Tokens.Input)
	}
	if resp.Tokens.Output != 50 {
		t.Errorf("expected output tokens 50, got %d", resp.Tokens.Output)
	}
	if resp.Tokens.Cached != 10 {
		t.Errorf("expected cached tokens 10, got %d", resp.Tokens.Cached)
	}
}

func TestGeminiAgent_Execute(t *testing.T) {
	geminiJSON := `{
		"candidates": [{"content": {"parts": [{"text": "Hello from Gemini"}], "role": "model"}}],
		"usageMetadata": {"promptTokenCount": 200, "candidatesTokenCount": 80, "totalTokenCount": 280, "cachedContentTokenCount": 0}
	}`

	agent := NewGeminiAgent()
	agent.SetCommandRunner(mockCommandRunner(geminiJSON, "", 0))

	resp, err := agent.Execute(context.Background(), ExecuteRequest{
		Prompt:  "test prompt",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Output != "Hello from Gemini" {
		t.Errorf("expected 'Hello from Gemini', got %q", resp.Output)
	}
	if resp.Tokens.Total != 280 {
		t.Errorf("expected total tokens 280, got %d", resp.Tokens.Total)
	}
}

func TestCodexAgent_Execute(t *testing.T) {
	codexJSONL := `{"type":"message","message":{"role":"assistant","content":"Hello from Codex"}}
{"type":"usage","usage":{"input_tokens":150,"output_tokens":60,"total_tokens":210}}`

	agent := NewCodexAgent()
	agent.SetCommandRunner(mockCommandRunner(codexJSONL, "", 0))

	resp, err := agent.Execute(context.Background(), ExecuteRequest{
		Prompt:  "test prompt",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Output != "Hello from Codex" {
		t.Errorf("expected 'Hello from Codex', got %q", resp.Output)
	}
	if resp.Tokens.Total != 210 {
		t.Errorf("expected total tokens 210, got %d", resp.Tokens.Total)
	}
}

func TestCopilotAgent_Execute(t *testing.T) {
	agent := NewCopilotAgent()
	agent.SetCommandRunner(mockCommandRunner("Hello from Copilot", "", 0))

	resp, err := agent.Execute(context.Background(), ExecuteRequest{
		Prompt:  "test prompt",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Output != "Hello from Copilot" {
		t.Errorf("expected 'Hello from Copilot', got %q", resp.Output)
	}
}

func TestRegistry(t *testing.T) {
	reg := NewRegistry()

	claude := NewClaudeAgent()
	gemini := NewGeminiAgent()

	if err := reg.Register(claude); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := reg.Register(gemini); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// duplicate registration
	if err := reg.Register(claude); err == nil {
		t.Error("expected error for duplicate registration")
	}

	// get
	a, err := reg.Get("claude")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "claude" {
		t.Errorf("expected 'claude', got %q", a.Name())
	}

	// not found
	_, err = reg.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent agent")
	}

	// list
	all := reg.List()
	if len(all) != 2 {
		t.Errorf("expected 2 agents, got %d", len(all))
	}

	// find by specialty
	architects := reg.FindBySpecialty(TaskArchitecture)
	if len(architects) != 1 || architects[0].Name() != "claude" {
		t.Errorf("expected claude for architecture specialty")
	}
}

func TestTokenBudget(t *testing.T) {
	agent := NewClaudeAgent()
	budget := agent.TokenBudget()
	if budget.DailyLimit != 1000000 {
		t.Errorf("expected daily limit 1000000, got %d", budget.DailyLimit)
	}

	agent.UpdateTokenUsage(TokenUsage{Total: 500})
	budget = agent.TokenBudget()
	if budget.UsedToday != 500 {
		t.Errorf("expected used 500, got %d", budget.UsedToday)
	}
	if budget.Remaining != 999500 {
		t.Errorf("expected remaining 999500, got %d", budget.Remaining)
	}
}

func TestFactory(t *testing.T) {
	// factories are registered in init() of each adapter file
	names := RegisteredFactories()
	if len(names) < 4 {
		t.Errorf("expected at least 4 factories, got %d", len(names))
	}

	a, err := CreateAgent("claude")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name() != "claude" {
		t.Errorf("expected 'claude', got %q", a.Name())
	}

	_, err = CreateAgent("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent factory")
	}
}

func TestAgentError(t *testing.T) {
	agent := NewClaudeAgent()
	agent.SetCommandRunner(mockCommandRunner("", "auth failed", 1))

	_, err := agent.Execute(context.Background(), ExecuteRequest{
		Prompt:  "test",
		Timeout: 10 * time.Second,
	})
	if err == nil {
		t.Error("expected error for non-zero exit code")
	}
}

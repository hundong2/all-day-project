package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/maao/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "maao",
	Short: "Multi-Agent AI Orchestrator",
	Long:  "MAAO orchestrates multiple AI CLI agents (Claude, Gemini, Codex, Copilot) to collaborate on GitHub repositories.",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// init subcommand
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize MAAO configuration",
		RunE:  runInit,
	})
	// register
	rootCmd.AddCommand(&cobra.Command{
		Use:   "register [repo-url]",
		Short: "Register a repository to monitor",
		Args:  cobra.ExactArgs(1),
		RunE:  runRegister,
	})
	// agents
	agentsCmd := &cobra.Command{Use: "agents", Short: "Manage agents"}
	agentsCmd.AddCommand(&cobra.Command{Use: "check", Short: "Check agent availability", RunE: runAgentsCheck})
	agentsCmd.AddCommand(&cobra.Command{Use: "status", Short: "Show token usage", RunE: runAgentsStatus})
	rootCmd.AddCommand(agentsCmd)
	// run
	runCmd := &cobra.Command{Use: "run", Short: "Start orchestration", RunE: runRun}
	runCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "Verbose output")
	runCmd.Flags().BoolVar(&flagDashboard, "dashboard", false, "TUI dashboard mode")
	rootCmd.AddCommand(runCmd)
	// status
	rootCmd.AddCommand(&cobra.Command{Use: "status", Short: "Show workflow status", RunE: runStatus})
	// logs
	rootCmd.AddCommand(&cobra.Command{Use: "logs", Short: "Show logs", RunE: runLogs})
	// config
	configCmd := &cobra.Command{Use: "config", Short: "Manage configuration"}
	configCmd.AddCommand(&cobra.Command{Use: "show", Short: "Show current config", RunE: runConfigShow})
	configCmd.AddCommand(&cobra.Command{Use: "set [key] [value]", Short: "Set config value", Args: cobra.ExactArgs(2), RunE: runConfigSet})
	rootCmd.AddCommand(configCmd)
}

var (
	flagVerbose   bool
	flagDashboard bool
)

const configDir = ".maao"
const configFile = "config.yaml"

func runInit(cmd *cobra.Command, args []string) error {
	dir := configDir
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	path := filepath.Join(dir, configFile)
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stderr, "Config already exists at %s\n", path)
		return nil
	}

	var cfg config.Config
	cfg.SetDefaults()
	cfg.Agents = map[string]config.AgentConfig{
		"claude": {Path: "/usr/local/bin/claude", Enabled: true, DailyTokenBudget: 500000, Specialties: []string{"architecture", "refactoring"}},
		"gemini": {Path: "/usr/local/bin/gemini", Enabled: true, DailyTokenBudget: 1000000, Specialties: []string{"planning", "documentation"}},
		"codex":  {Path: "/usr/local/bin/codex", Enabled: true, DailyTokenBudget: 500000, Specialties: []string{"debugging", "testing"}},
	}
	cfg.Repositories = []config.RepoConfig{
		{URL: "https://github.com/user/project", LocalPath: ".", PollInterval: "60s", DefaultBranch: "main"},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	fmt.Printf("Initialized MAAO config at %s\n", path)
	fmt.Println("Edit the config file to customize your setup.")
	return nil
}

func runRegister(cmd *cobra.Command, args []string) error {
	path := filepath.Join(configDir, configFile)
	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("loading config: %w (run 'maao init' first)", err)
	}

	repoURL := args[0]
	for _, r := range cfg.Repositories {
		if r.URL == repoURL {
			fmt.Fprintf(os.Stderr, "Repository %s is already registered\n", repoURL)
			return nil
		}
	}

	cfg.Repositories = append(cfg.Repositories, config.RepoConfig{
		URL:           repoURL,
		LocalPath:     ".",
		PollInterval:  "60s",
		DefaultBranch: "main",
	})

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	fmt.Printf("Registered repository: %s\n", repoURL)
	return nil
}

func runAgentsCheck(cmd *cobra.Command, args []string) error {
	path := filepath.Join(configDir, configFile)
	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("loading config: %w (run 'maao init' first)", err)
	}

	allOk := true
	for name, agent := range cfg.Agents {
		if !agent.Enabled {
			fmt.Printf("  %-10s DISABLED\n", name)
			continue
		}

		binPath, lookErr := exec.LookPath(agent.Path)
		if lookErr != nil {
			fmt.Printf("  %-10s NOT FOUND (%s)\n", name, agent.Path)
			allOk = false
			continue
		}

		// Try to get version
		version := "unknown"
		out, verErr := exec.Command(binPath, "--version").CombinedOutput()
		if verErr == nil && len(out) > 0 {
			version = string(out)
			if len(version) > 60 {
				version = version[:60]
			}
		}
		fmt.Printf("  %-10s OK  %s  (%s)\n", name, binPath, version)
	}

	if !allOk {
		fmt.Println("\nSome agents are not available. Install them or disable in config.")
	}
	return nil
}

func runAgentsStatus(cmd *cobra.Command, args []string) error {
	path := filepath.Join(configDir, configFile)
	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("loading config: %w (run 'maao init' first)", err)
	}

	fmt.Println("Agent Token Budgets (Daily):")
	fmt.Println("----------------------------")
	for name, agent := range cfg.Agents {
		if !agent.Enabled {
			continue
		}
		// Token usage tracking is placeholder until store layer is integrated
		used := 0
		remaining := agent.DailyTokenBudget - used
		fmt.Printf("  %-10s  %7d / %7d remaining  (specialties: %v)\n",
			name, remaining, agent.DailyTokenBudget, agent.Specialties)
	}
	return nil
}

// runRun, runStatus, runLogs depend on orchestrator - placeholder for now
func runRun(cmd *cobra.Command, args []string) error    { fmt.Println("TODO: run (requires orchestrator)"); return nil }
func runStatus(cmd *cobra.Command, args []string) error  { fmt.Println("TODO: status (requires orchestrator)"); return nil }
func runLogs(cmd *cobra.Command, args []string) error    { fmt.Println("TODO: logs (requires orchestrator)"); return nil }

func runConfigShow(cmd *cobra.Command, args []string) error {
	path := filepath.Join(configDir, configFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading config: %w (run 'maao init' first)", err)
	}
	fmt.Print(string(data))
	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	fmt.Printf("TODO: config set %s=%s (use direct YAML editing for now)\n", args[0], args[1])
	return nil
}

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	var verbose, dashboard bool
	runCmd := &cobra.Command{Use: "run", Short: "Start orchestration", RunE: runRun}
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	runCmd.Flags().BoolVar(&dashboard, "dashboard", false, "TUI dashboard mode")
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

func runInit(cmd *cobra.Command, args []string) error           { fmt.Println("TODO: init"); return nil }
func runRegister(cmd *cobra.Command, args []string) error       { fmt.Println("TODO: register", args[0]); return nil }
func runAgentsCheck(cmd *cobra.Command, args []string) error    { fmt.Println("TODO: agents check"); return nil }
func runAgentsStatus(cmd *cobra.Command, args []string) error   { fmt.Println("TODO: agents status"); return nil }
func runRun(cmd *cobra.Command, args []string) error            { fmt.Println("TODO: run"); return nil }
func runStatus(cmd *cobra.Command, args []string) error         { fmt.Println("TODO: status"); return nil }
func runLogs(cmd *cobra.Command, args []string) error           { fmt.Println("TODO: logs"); return nil }
func runConfigShow(cmd *cobra.Command, args []string) error     { fmt.Println("TODO: config show"); return nil }
func runConfigSet(cmd *cobra.Command, args []string) error      { fmt.Println("TODO: config set"); return nil }

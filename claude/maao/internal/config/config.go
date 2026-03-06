package config

import "time"

type Config struct {
	Agents       map[string]AgentConfig `yaml:"agents"`
	PM           PMConfig               `yaml:"pm"`
	Repositories []RepoConfig           `yaml:"repositories"`
	GitHub       GitHubConfig           `yaml:"github"`
	Notification NotificationConfig     `yaml:"notification"`
	Workspace    WorkspaceConfig        `yaml:"workspace"`
	CI           CIConfig               `yaml:"ci"`
}

type AgentConfig struct {
	Path             string         `yaml:"path"`
	Enabled          bool           `yaml:"enabled"`
	DailyTokenBudget int            `yaml:"daily_token_budget"`
	Specialties      []string       `yaml:"specialties"`
	Headless         HeadlessConfig `yaml:"headless"`
}

type HeadlessConfig struct {
	Flags     []string `yaml:"flags"`
	MaxTurns  int      `yaml:"max_turns"`
	Timeout   string   `yaml:"timeout"`
	Model     string   `yaml:"model"`
	DenyTools []string `yaml:"deny_tools"`
}

type PMConfig struct {
	DefaultAgent     string           `yaml:"default_agent"`
	DiscussionRounds int              `yaml:"discussion_rounds"`
	ReviewMaxRounds  int              `yaml:"review_max_rounds"`
	TaskEstimation   EstimationConfig `yaml:"task_estimation"`
}

type EstimationConfig struct {
	SmallLOCThreshold  int `yaml:"small_loc_threshold"`
	MediumLOCThreshold int `yaml:"medium_loc_threshold"`
	LargeTokenEstimate int `yaml:"large_token_estimate"`
}

type RepoConfig struct {
	URL           string `yaml:"url"`
	LocalPath     string `yaml:"local_path"`
	PollInterval  string `yaml:"poll_interval"`
	DefaultBranch string `yaml:"default_branch"`
}

type GitHubConfig struct {
	Token  string `yaml:"token"`
	APIURL string `yaml:"api_url"`
}

type NotificationConfig struct {
	Email EmailConfig `yaml:"email"`
}

type EmailConfig struct {
	To       string `yaml:"to"`
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type WorkspaceConfig struct {
	WorktreeDir  string `yaml:"worktree_dir"`
	BranchPrefix string `yaml:"branch_prefix"`
	AutoCleanup  bool   `yaml:"auto_cleanup"`
}

type CIConfig struct {
	Enabled           bool     `yaml:"enabled"`
	WaitForCompletion bool     `yaml:"wait_for_completion"`
	Timeout           string   `yaml:"timeout"`
	RequiredChecks    []string `yaml:"required_checks"`
}

func (c *Config) SetDefaults() {
	if c.PM.DefaultAgent == "" {
		c.PM.DefaultAgent = "gemini"
	}
	if c.PM.DiscussionRounds == 0 {
		c.PM.DiscussionRounds = 3
	}
	if c.PM.ReviewMaxRounds == 0 {
		c.PM.ReviewMaxRounds = 3
	}
	if c.Workspace.WorktreeDir == "" {
		c.Workspace.WorktreeDir = ".worktrees"
	}
	if c.Workspace.BranchPrefix == "" {
		c.Workspace.BranchPrefix = "agent/"
	}
	if c.GitHub.APIURL == "" {
		c.GitHub.APIURL = "https://api.github.com"
	}
	_ = time.Second // ensure time import used
}

package config

import "fmt"

func Validate(cfg *Config) error {
	if cfg.GitHub.Token == "" {
		return fmt.Errorf("github.token is required")
	}
	if len(cfg.Repositories) == 0 {
		return fmt.Errorf("at least one repository must be configured")
	}
	for name, agent := range cfg.Agents {
		if agent.Enabled && agent.Path == "" {
			return fmt.Errorf("agent %s: path is required when enabled", name)
		}
	}
	return nil
}

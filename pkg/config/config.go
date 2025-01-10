package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/sethvargo/go-githubactions"
)

func NewFromInputs(a *githubactions.Action) (*Config, error) {
	cfg := Config{}

    command := os.Getenv("GAS_COMMAND")
    if command == "" {
        cfg.Command = a.GetInput(commandInput)
        cfg.AccessToken = a.GetInput(accessTokenInput)
        cfg.ProjectId = a.GetInput(projectIdInput)
        cfg.ScriptDir = a.GetInput(scriptDirInput)
        cfg.Function = a.GetInput(functionInput)
    } else {
        cfg.Command = command
        cfg.AccessToken = os.Getenv("GAS_ACCESS_TOKEN")
        cfg.ProjectId = os.Getenv("GAS_PROJECT_ID")
        cfg.ScriptDir = os.Getenv("GAS_SCRIPT_DIR")
        cfg.Function = os.Getenv("GAS_FUNCTION")
    }

	// Validate
	if cfg.AccessToken == "" || cfg.ProjectId == "" {
		return nil, fmt.Errorf("at least one of the '%s' and '%s' inputs are not set", accessTokenInput, projectIdInput)
	}

	if !slices.Contains(supportedCommands, cfg.Command) {
		return nil, fmt.Errorf("submitted command '%s' is not supported", cfg.Command)
	}

	return &cfg, nil
}

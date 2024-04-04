package config

import (
	"fmt"
	"slices"

	"github.com/sethvargo/go-githubactions"
)

func NewFromInputs(a *githubactions.Action) (*Config, error) {
	cfg := Config{}

	cfg.Command = a.GetInput(commandInput)
	cfg.AccessToken = a.GetInput(accessTokenInput)
	cfg.ProjectId = a.GetInput(projectIdInput)
	cfg.ScriptDir = a.GetInput(scriptDirInput)
	cfg.Entrypoint = a.GetInput(entrypointInput)

	// Validate
	if cfg.AccessToken == "" || cfg.ProjectId == "" {
		return nil, fmt.Errorf("at least one of the '%s' and '%s' inputs are not set", accessTokenInput, projectIdInput)
	}

	if !slices.Contains(supportedCommands, cfg.Command) {
		return nil, fmt.Errorf("submitted command '%s' is not supported", cfg.Command)
	}

	return &cfg, nil
}

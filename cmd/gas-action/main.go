package main

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-githubactions"
	"github.com/shiftavenue/gas-action/pkg/commands"
	"github.com/shiftavenue/gas-action/pkg/config"
)

func run() error {
	ctx := context.Background()
	action := githubactions.New()

	cfg, err := config.NewFromInputs(action)
	if err != nil {
		return err
	}

	switch cfg.Command {
	case "push":
		return commands.Push(ctx, cfg)
	case "deploy":
		return commands.Deploy(ctx, cfg)
	case "run":
		return commands.Run(ctx, cfg)
	default:
		return fmt.Errorf("unsupported command '%s' specified", cfg.Command)
	}
}

// Entrypoint
func main() {
	err := run()
	if err != nil {
		githubactions.Fatalf("%v", err)
	}
}

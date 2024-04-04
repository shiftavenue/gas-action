package commands

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/shiftavenue/gas-action/pkg/client"
	"github.com/shiftavenue/gas-action/pkg/config"
	"google.golang.org/api/script/v1"
)

// Run triggers an execution of an Apps Script
// NOTE: only works if the script is deployed as an API executable, other deployment types are not supported
func Run(ctx context.Context, cfg *config.Config) error {
	// create client
	client, err := client.NewWithAccessToken(ctx, cfg.AccessToken)
	if err != nil {
		return fmt.Errorf("could not create Apps Script client: %v", err)
	}

	// Get latest versioned deployment
	version := int64(1)
	deploymentID := ""
	existingDeployments, err := client.Projects.Deployments.List(cfg.ScriptDir).Do()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("error while running script: %s", err))
	}

	for _, d := range existingDeployments.Deployments {
		if d.DeploymentConfig.VersionNumber >= version {
			version = d.DeploymentConfig.VersionNumber
			deploymentID = d.DeploymentId
		}
	}

	if deploymentID == "" {
		return fmt.Errorf("no active versioned deployment of script project %s found; aborting", cfg.ProjectId)
	}

	// Trigger execution
	req := &script.ExecutionRequest{
		Function: cfg.Entrypoint,
	}

	run, err := client.Scripts.Run(deploymentID, req).Do()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("error while running script: %s", err))
	}

	if run.Error != nil {
		log.Fatal().Msg(fmt.Sprintf("script execution failed with code %d: %s", run.Error.Code, run.Error.Message))
	}

	log.Info().Msg(fmt.Sprintf("Deployment %s of script project %s executed successfully", deploymentID, cfg.ProjectId))

	return nil
}

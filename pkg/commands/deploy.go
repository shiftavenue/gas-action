package commands

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/shiftavenue/gas-action/pkg/client"
	"github.com/shiftavenue/gas-action/pkg/config"
	"google.golang.org/api/script/v1"
)

// Deploy creates a new versioned deployment
// NOTE: Archiving or deleting old deployments cannot be automated (via API or anything) and is a manual task
func Deploy(ctx context.Context, cfg *config.Config) error {
	// create client
	client, err := client.NewWithAccessToken(ctx, cfg.AccessToken)
	if err != nil {
		return fmt.Errorf("could not create Apps Script client: %v", err)
	}

	// Determine and increment version
	version := int64(1)

	existingDeployments, err := client.Projects.Deployments.List(cfg.ProjectId).Do()
	if err != nil {
		return fmt.Errorf("error while listing all existing deployments: %v", err)
	}

	// Determine latest version
	for _, d := range existingDeployments.Deployments {
		if d.DeploymentConfig.VersionNumber > version {
			version = d.DeploymentConfig.VersionNumber
		}
	}

	versionConf := &script.Version{
		VersionNumber: version,
		Description:   fmt.Sprintf("Deployment version %d (created by GitHub Actions)", version),
	}

	// Create Apps Script version
	versionResp, err := client.Projects.Versions.Create(cfg.ProjectId, versionConf).Do()
	if err != nil {
		return fmt.Errorf("error while creating new script version: %v", err)
	}

	log.Info().Msg(fmt.Sprintf("Script version %d successfully created", versionConf.VersionNumber))

	deploymentConf := &script.DeploymentConfig{
		ScriptId:      cfg.ProjectId,
		Description:   "Script deployment (deployed by GitHub Actions)",
		VersionNumber: versionResp.VersionNumber,
	}

	// Deploy Apps Script as API
	deployResp, err := client.Projects.Deployments.Create(cfg.ProjectId, deploymentConf).Do()
	if err != nil {
		return fmt.Errorf("error while deploying script: %v", err)
	}

	log.Info().Msg(fmt.Sprintf("Script successfully deployed; deployment ID is %s", deployResp.DeploymentId))

	return nil
}

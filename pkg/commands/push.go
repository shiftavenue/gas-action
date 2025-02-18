package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/shiftavenue/gas-action/pkg/client"
	"github.com/shiftavenue/gas-action/pkg/config"
	"google.golang.org/api/script/v1"
)

// Push uses the code in GitHub and pushes it to the Apps Script portal
func Push(ctx context.Context, cfg *config.Config) error {
	// create client
	client, err := client.NewWithAccessToken(ctx, cfg.AccessToken)
	if err != nil {
		return fmt.Errorf("could not create Apps Script client: %v", err)
	}

	// Get all files that need to be deployed
	scriptFiles := []*script.File{}

	// Get all files in script directory
	scriptDir, err := os.Open(cfg.ScriptDir)
	if err != nil {
		return err
	}

	defer scriptDir.Close()

	files, err := scriptDir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Get file content
		fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", cfg.ScriptDir, file))
		if err != nil {
			return fmt.Errorf("unable to read single script file: %v", err)
		}

		// Set file type
		fileType := "ENUM_TYPE_UNSPECIFIED"

		switch {
		case strings.HasSuffix(file, ".js"):
			fileType = "SERVER_JS"
		case strings.HasSuffix(file, ".json"):
			fileType = "JSON"
		case strings.HasSuffix(file, ".html"):
			fileType = "HTML"
		}

		// Add file
		scriptFiles = append(scriptFiles, &script.File{
			Name:   strings.Split(file, ".")[0],
			Type:   fileType,
			Source: string(fileContent),
		})
	}

	content := &script.Content{
		ScriptId: cfg.ProjectId,
		Files:    scriptFiles,
	}

	// Push code to Apps Script
	updateResp, err := client.Projects.UpdateContent(cfg.ProjectId, content).Do()
	if err != nil {
		return fmt.Errorf("error while pushing and updating script: %v", err)
	}

	log.Info().Msg(fmt.Sprintf("Script %s successfully pushed", updateResp.ScriptId))

	return nil
}

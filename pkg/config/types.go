package config

// String identifiers of all inputs
const (
	commandInput     = "command"
	accessTokenInput = "gcp-access-token"
	projectIdInput   = "project-id"
	scriptDirInput   = "script-dir"
	entrypointInput  = "entrypoint"
)

// Allowed/supported commands
var (
	supportedCommands = []string{"push", "deploy", "run"}
)

// Config holds all possible Action inputs
type Config struct {
	Command     string
	AccessToken string
	ProjectId   string
	ScriptDir   string
	Entrypoint  string
}
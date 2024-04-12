package config

// String identifiers of all inputs
const (
	commandInput     = "command"
	accessTokenInput = "access-token"
	projectIdInput   = "project-id"
	scriptDirInput   = "script-dir"
	functionInput    = "function"
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
	Function    string
}

package mcsm

// Config represents PufferPanel provider configuration.
type Config struct {
	BaseUrl  string `validate:"required"` // BaseUrl of PufferPanel.
	ApiKey   string `validate:"required"` // ClientId of PufferPanel service account.
	DaemonId string `validate:"required"` // Daemon Id to watch
}

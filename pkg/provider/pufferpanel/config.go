package pufferpanel

// Config represents PufferPanel provider configuration.
type Config struct {
	BaseUrl        string `validate:"required"` // BaseUrl of PufferPanel.
	ClientId       string `validate:"required"` // ClientId of PufferPanel service account.
	ClientSecret   string `validate:"required"` // ClientSecret of PufferPanel service account.
	ConfigFilePath string // Path of config file.
}

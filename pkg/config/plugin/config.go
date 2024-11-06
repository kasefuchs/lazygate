package plugin

// Config represents plugin config.
type Config struct {
	Namespace string // Namespace of this proxy.
}

// DefaultConfig returns default config.
func DefaultConfig() *Config {
	return &Config{
		Namespace: "default",
	}
}

package plugin

type Config struct {
	Provider  string `validate:"required"` // Provider to use.
	Namespace string `validate:"required"` // Namespace of this proxy.
}

func DefaultConfig() *Config {
	return &Config{
		Provider:  "docker",
		Namespace: "default",
	}
}

package static

import (
	"os"

	"github.com/traefik/paerser/env"
)

// Prefix of configuration environment variables.
const prefix = "LAZYGATE_"

// ParseEnv parses environment variables to static configuration.
func ParseEnv() (*Config, error) {
	cfg := &Config{}

	vars := env.FindPrefixedEnvVars(os.Environ(), prefix, cfg)
	if err := env.Decode(vars, prefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

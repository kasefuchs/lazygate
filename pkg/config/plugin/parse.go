package plugin

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/traefik/paerser/env"
)

// Root env prefix.
const rootPrefix = "LAZYGATE_"

// Config validator.
var validate = validator.New()

// ParseEnv parses plugin configuration from environment.
func ParseEnv() (*Config, error) {
	cfg := DefaultConfig()

	vars := env.FindPrefixedEnvVars(os.Environ(), rootPrefix, cfg)
	if err := env.Decode(vars, rootPrefix, cfg); err != nil {
		return nil, err
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

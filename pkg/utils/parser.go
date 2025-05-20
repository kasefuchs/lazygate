package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/traefik/paerser/env"
	"github.com/traefik/paerser/parser"
)

// Root label name.
const (
	rootLabel         = "lazygate"
	internalLabel     = "lazygateinternal"
	EnvPrefix         = "LAZYGATE_"
	internalEnvPrefix = "LAZYGATEINTERNAL_"
)

// Structure validator.
var validate = validator.New()

// ParseLabels parses labels to dynamic config.
func ParseLabels(labels map[string]string, cfg interface{}, rootName string) (interface{}, error) {
	transformedLabels := TransformLabels(labels, rootName, internalLabel)

	if err := parser.Decode(transformedLabels, cfg, internalLabel, internalLabel); err != nil {
		return nil, err
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ParseTags parses tags to dynamic config.
func ParseTags(tags []string, cfg interface{}, rootName string) (interface{}, error) {
	labels := TagsToLabels(tags)

	return ParseLabels(labels, cfg, rootName)
}

// ParseEnv parses plugin configuration from environment.
func ParseEnv(cfg interface{}, prefix string) (interface{}, error) {
	vars := env.FindPrefixedEnvVars(os.Environ(), prefix, cfg)
	transformedVars := TransformEnvVars(vars, prefix, internalEnvPrefix)
	if err := env.Decode(transformedVars, internalEnvPrefix, cfg); err != nil {
		return nil, err
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ChildEnvPrefix(names ...string) string {
	joinedName := strings.ToUpper(strings.Join(names, "_"))
	return fmt.Sprintf("%s%s_", EnvPrefix, joinedName)
}

func ChildLabel(names ...string) string {
	joinedName := strings.Join(names, ".")
	return fmt.Sprintf("%s.%s", rootLabel, joinedName)
}

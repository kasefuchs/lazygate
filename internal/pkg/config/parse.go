package config

import (
	"github.com/kasefuchs/lazygate/internal/pkg/util"
	"github.com/traefik/paerser/parser"
)

func ParseLabels(labels map[string]string) (*Config, error) {
	cfg := &Config{}

	err := parser.Decode(labels, cfg, "lazygate")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ParseTags(tags []string) (*Config, error) {
	labels := util.TagsToLabels(tags, "lazygate")

	return ParseLabels(labels)
}

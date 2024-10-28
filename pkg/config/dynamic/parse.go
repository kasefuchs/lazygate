package dynamic

import (
	"github.com/kasefuchs/lazygate/pkg/utils"
	"github.com/traefik/paerser/parser"
)

// Root label name.
const rootName = "plugin"

// ParseLabels parses labels to dynamic config.
func ParseLabels(labels map[string]string) (*Config, error) {
	cfg := &Config{}

	err := parser.Decode(labels, cfg, rootName, rootName)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// ParseTags parses tags to dynamic config.
func ParseTags(tags []string) (*Config, error) {
	labels := utils.TagsToLabels(tags)

	return ParseLabels(labels)
}
